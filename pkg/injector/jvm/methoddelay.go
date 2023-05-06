/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jvm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmetad/pkg/utils/process"
	"os"
	"strconv"
	"strings"
)

func init() {
	injector.Register(TargetJVM, FaultMethodDelay, func() injector.IInjector { return &MethodDelayInjector{} })
}

type MethodDelayInjector struct {
	injector.BaseInjector
	Args    MethodDelayArgs
	Runtime MethodDelayRuntime
}

type MethodDelayArgs struct {
	Pid        int    `json:"pid,omitempty"`
	Key        string `json:"key,omitempty"`
	MethodList string `json:"method"` // class@method@3000,
}

type MethodDelayRuntime struct {
	AttackPids []int `json:"attack_pids"`
}

func (i *MethodDelayInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *MethodDelayInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *MethodDelayInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&i.Args.Pid, "pid", "p", 0, "target process's pid")
	cmd.Flags().StringVarP(&i.Args.Key, "key", "k", "", "the key used to grep to get target process, the effect is equivalent to \"ps -ef | grep [key]\". if \"pid\" provided, \"key\" will be ignored")
	cmd.Flags().StringVarP(&i.Args.MethodList, "method", "m", "", "target method of the process, format: \"class1@method1@delay_ms,class1@method2@delay_ms\", eg: \"com.test.Client@sayHello@3000\"")
}

func (i *MethodDelayInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	pidList, err := getPidList(ctx, i.Args.Pid, i.Args.Key, i.Info.ContainerRuntime, i.Info.ContainerId)
	if err != nil {
		return fmt.Errorf("get target process's pid error: %s", err.Error())
	}

	for _, pid := range pidList {
		if filesys.CheckDir(getRuleDir(pid)) == nil {
			return fmt.Errorf("has jvm experiment running in process[%d]", pid)
		}
	}

	_, err = getMethodList(i.Args.MethodList, FaultMethodDelay)
	if err != nil {
		return fmt.Errorf("\"method\" is invalid: %s", err.Error())
	}

	return nil
}

func (i *MethodDelayInjector) Inject(ctx context.Context) error {
	var (
		pidList     []int
		err         error
		needRecover bool
		errMsg      string
		logger      = log.GetLogger(ctx)
	)

	pidList, _ = getPidList(ctx, i.Args.Pid, i.Args.Key, i.Info.ContainerRuntime, i.Info.ContainerId)

	// save target
	i.Runtime.AttackPids = pidList

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	methodListMap, _ := getMethodList(i.Args.MethodList, FaultMethodDelay)
	ruleBytes, err := json.Marshal(getRuleConfig(methodListMap, timeout))
	if err != nil {
		return fmt.Errorf("get rule file bytes error: %s", err.Error())
	}
	logger.Debugf("rule json: %s", string(ruleBytes))

	// create rule file
	for _, pid := range pidList {
		if err := writeRule(ctx, pid, ruleBytes); err != nil {
			needRecover = true
			errMsg = fmt.Sprintf("write rule for process[%d] error: %s", pid, err.Error())
			break
		}
	}

	if !needRecover {
		// execute
		for _, pid := range pidList {
			if _, err := cmdexec.StartBashCmdAndWaitPid(ctx, getCmd(pid), TimeoutSec); err != nil {
				needRecover = true
				errMsg = fmt.Sprintf("execute fault for process[%d] error: %s", pid, err.Error())
				break
			}
		}
	}

	if needRecover {
		// undo recover
		if err := i.Recover(ctx); err != nil {
			logger.Warnf("undo error: %s", err.Error())
		}
	}

	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}

	return nil
}

func (i *MethodDelayInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	logger := log.GetLogger(ctx)
	for _, pid := range i.Runtime.AttackPids {
		targetDir := getRuleDir(pid)
		if filesys.CheckDir(targetDir) == nil {
			if err := os.RemoveAll(targetDir); err != nil {
				logger.Errorf("remove dir[%s] error: %s", targetDir, err.Error())
			}
		}
	}

	return nil
}

func getPidList(ctx context.Context, pid int, key string, cr, cId string) (pidList []int, err error) {
	if pid > 0 {
		if cr != "" {
			return nil, fmt.Errorf("not support args \"pid\" in container runtime: \"%s\"", cr)
		}
		pidList = append(pidList, pid)
	} else if key != "" {
		pidList, err = process.GetPidListByKey(ctx, cr, cId, key)
		if err != nil {
			return nil, fmt.Errorf("get pid list by grep key[%s] error: %s", key, err.Error())
		}
	} else {
		return nil, fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if len(pidList) == 0 {
		return nil, fmt.Errorf("target pid is empty")
	}

	return
}

func getMethodList(methodStr string, fault string) (map[string][]*MethodJVMRule, error) {
	var result = make(map[string][]*MethodJVMRule)
	if methodStr == "" {
		return nil, fmt.Errorf("is empty")
	}

	tmpList := strings.Split(methodStr, MethodRuleSplit)
	for _, unitMethod := range tmpList {
		kv := strings.Split(unitMethod, ClassMethodSplit)
		if len(kv) != 3 {
			return nil, fmt.Errorf("\"%s\" is not a valid format", kv)
		}

		var (
			className, methodName, valueStr = kv[0], kv[1], kv[2]
			rule                            *MethodJVMRule
			err                             error
		)

		switch fault {
		case FaultMethodDelay:
			rule, err = getMethodDelayRule(methodName, valueStr)
		case FaultMethodReturn:
			rule, err = getMethodReturnRule(methodName, valueStr)
		case FaultMethodException:
			rule, err = getMethodExceptionRule(methodName, valueStr)
		default:
			return nil, fmt.Errorf("not support fault: %s", fault)
		}

		if err != nil {
			return nil, fmt.Errorf("get rule of[%s] error: %s", unitMethod, err.Error())
		}

		result[className] = append(result[className], rule)
	}

	return result, nil
}

func getRuleConfig(methodListMap map[string][]*MethodJVMRule, timeout int64) *JVMRuleConfig {
	var classRuleList []*ClassJVMRule
	for className, methods := range methodListMap {
		classRuleList = append(classRuleList, &ClassJVMRule{
			Class:      className,
			MethodList: methods,
		})
	}

	return &JVMRuleConfig{
		Duration:  timeout,
		ClassList: classRuleList,
	}
}

func getCmd(pid int) string {
	return fmt.Sprintf("cd %s && java %s %d %s %s", utils.GetToolDir(), AttacherTool, pid, utils.GetToolPath(JVMAgentTool), getRuleFile(pid))
}

func writeRule(ctx context.Context, pid int, ruleBytes []byte) error {
	dir := getRuleDir(pid)
	if err := filesys.MkdirP(ctx, dir); err != nil {
		return fmt.Errorf("create rule dir[%s] error: %s", dir, err.Error())
	}

	// write rule to fileName
	fileName := getRuleFile(pid)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open file[%s] fail: %s", fileName, err.Error())
	}

	if _, err := f.Write(ruleBytes); err != nil {
		return fmt.Errorf("write jvm rule error: %s", err.Error())
	}

	return nil
}

func getRuleDir(pid int) string {
	return fmt.Sprintf("%s/%s/%d", utils.GetRunPath(), JVMRuleDir, pid)
}

func getRuleFile(pid int) string {
	return fmt.Sprintf("%s/%s/%d/%s", utils.GetRunPath(), JVMRuleDir, pid, JVMRuleFile)
}

func getMethodDelayRule(methodName, delayMsStr string) (*MethodJVMRule, error) {
	delayMs, err := strconv.Atoi(delayMsStr)
	if err != nil {
		return nil, fmt.Errorf("is not a valid delay ms: %s", err.Error())
	}
	if delayMs <= 0 {
		return nil, fmt.Errorf("delay ms must larger than 0")
	}

	return &MethodJVMRule{
		Method:  methodName,
		Fault:   InsertBeforeInject,
		Content: fmt.Sprintf("try {Thread.sleep((long)%d);} catch (Exception e) {System.out.println(\"ChaosMeta Delay Inject Failed: \" + e.getMessage());}", delayMs),
	}, nil
}
