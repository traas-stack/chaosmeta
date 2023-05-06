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
	"os"
)

func init() {
	injector.Register(TargetJVM, FaultMethodReturn, func() injector.IInjector { return &MethodReturnInjector{} })
}

type MethodReturnInjector struct {
	injector.BaseInjector
	Args    MethodReturnArgs
	Runtime MethodReturnRuntime
}

type MethodReturnArgs struct {
	Pid        int    `json:"pid,omitempty"`
	Key        string `json:"key,omitempty"`
	MethodList string `json:"method"` // class@method@"ok",
}

type MethodReturnRuntime struct {
	AttackPids []int `json:"attack_pids"`
}

func (i *MethodReturnInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *MethodReturnInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *MethodReturnInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&i.Args.Pid, "pid", "p", 0, "target process's pid")
	cmd.Flags().StringVarP(&i.Args.Key, "key", "k", "", "the key used to grep to get target process, the effect is equivalent to \"ps -ef | grep [key]\". if \"pid\" provided, \"key\" will be ignored")
	cmd.Flags().StringVarP(&i.Args.MethodList, "method", "m", "", "target method of the process, format: \"class1@method1@return_value,class1@method2@return_value\", eg: com.test.Client@sayHello@\"ok\",com.test.Client@sayHello@5")
}

func (i *MethodReturnInjector) Validator(ctx context.Context) error {
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

	_, err = getMethodList(i.Args.MethodList, FaultMethodReturn)
	if err != nil {
		return fmt.Errorf("\"method\" is invalid: %s", err.Error())
	}

	return nil
}

func (i *MethodReturnInjector) Inject(ctx context.Context) error {
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

	methodListMap, _ := getMethodList(i.Args.MethodList, FaultMethodReturn)
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

func (i *MethodReturnInjector) Recover(ctx context.Context) error {
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

func getMethodReturnRule(methodName, valueStr string) (*MethodJVMRule, error) {
	return &MethodJVMRule{
		Method:  methodName,
		Fault:   SetBodyInject,
		Content: fmt.Sprintf("{ return %s; }", valueStr),
	}, nil
}
