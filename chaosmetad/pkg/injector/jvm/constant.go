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
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
	"os"
	"strings"
)

const (
	TargetJVM = "jvm"

	FaultMethodDelay     = "methoddelay"
	FaultMethodException = "methodexception"
	FaultMethodReturn    = "methodreturn"

	JVMRuleDir      = "jvm_rule"
	JVMContainerDir = "/tmp/chaosmeta_jvm"

	InsertAtInject     = "insertAt"
	InsertBeforeInject = "insertBefore"
	SetBodyInject      = "setBody"
	//InsertAfterInject  = "insertAfter"

	MethodRuleSplit  = ","
	ClassMethodSplit = "@"

	AttacherTool = "ChaosMetaJVMAttacher"
	JVMAgentTool = "ChaosMetaJVMAgent.jar"

	TimeoutSec = 2
)

type JVMRuleConfig struct {
	Duration  int64           `json:"Duration"`
	ClassList []*ClassJVMRule `json:"ClassList"`
}

type ClassJVMRule struct {
	Class      string           `json:"Class"`
	MethodList []*MethodJVMRule `json:"MethodList"`
}

type MethodJVMRule struct {
	Method    string `json:"Method"`
	Fault     string `json:"Fault"`
	Content   string `json:"Content"`
	ImportPkg string `json:"ImportPkg,omitempty"`
	LineNum   int    `json:"LineNum"`
}

func getRuleDir(cId string) string {
	if cId == "" {
		return fmt.Sprintf("%s/%s", utils.GetRunPath(), JVMRuleDir)
	} else {
		return fmt.Sprintf("%s/%s/%s", utils.GetRunPath(), JVMRuleDir, cId)
	}
}

// getRuleFile chaosmetad/jvm_rule/123.json or chaosmetad/jvm_rule/[containerId]/123.json
func getRuleFile(cId string, pid int) string {
	return fmt.Sprintf("%s/%d.json", getRuleDir(cId), pid)
}

func getContainerRuleDir() string {
	return JVMContainerDir
}

func getContainerRuleFile(pid int) string {
	return fmt.Sprintf("%s/%d.json", JVMContainerDir, pid)
}

func writeRule(ctx context.Context, cId string, pid int, ruleBytes []byte) error {
	dir := getRuleDir(cId)
	if err := filesys.MkdirP(ctx, dir); err != nil {
		return fmt.Errorf("create rule dir[%s] error: %s", dir, err.Error())
	}

	// write rule to fileName
	fileName := getRuleFile(cId, pid)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open file[%s] fail: %s", fileName, err.Error())
	}

	if _, err := f.Write(ruleBytes); err != nil {
		return fmt.Errorf("write jvm rule error: %s", err.Error())
	}

	return nil
}

// TODO: It is currently impossible to determine accurately because the ExecContainer function does not completely determine whether it is successful based on the return code.
func checkJavaCmd(ctx context.Context, cr, cId string) error {
	cmd := "java -help"
	if cr != "" {
		_, err := cmdexec.ExecContainer(ctx, cr, cId, []string{namespace.MNT, namespace.ENV}, cmd, cmdexec.ExecRun)
		return err
	} else {
		return cmdexec.RunBashCmdWithoutOutput(ctx, cmd)
	}
}

func doInject(ctx context.Context, cr, cId string, pidList []int, ruleBytes []byte) error {
	// create rule file
	for _, pid := range pidList {
		if err := writeRule(ctx, cId, pid, ruleBytes); err != nil {
			return fmt.Errorf("write rule for process[%d] error: %s", pid, err.Error())
		}
	}

	// cp tool to container
	if cr != "" {
		// cp rule file and jvm inject tool
		if err := filesys.MkdirPInContainer(ctx, cr, cId, getContainerRuleDir()); err != nil {
			return fmt.Errorf("mkdir tool dir in container error: %s", err.Error())
		}

		for _, pid := range pidList {
			src, dst := getRuleFile(cId, pid), getContainerRuleFile(pid)
			if err := cmdexec.CpContainerFile(ctx, cr, cId, src, dst); err != nil {
				return fmt.Errorf("cp file[%s] to [%s] in container[%s] error: %s", src, dst, cId, err.Error())
			}
		}

		jvmToolList := []string{"json-20190722.jar", "tools.jar", "javassist.jar", "MANIFEST.MF", "ChaosMetaJVMAttacher.class", "ChaosMetaJVMAgent.jar"}
		for _, unitTool := range jvmToolList {
			src, dst := utils.GetToolPath(unitTool), fmt.Sprintf("%s/%s", JVMContainerDir, unitTool)
			if err := cmdexec.CpContainerFile(ctx, cr, cId, src, dst); err != nil {
				return fmt.Errorf("cp file[%s] to [%s] in container[%s] error: %s", src, dst, cId, err.Error())
			}
		}

		for _, pid := range pidList {
			cmd := fmt.Sprintf("cd %s && java -cp .:tools.jar %s %d %s %s", JVMContainerDir, AttacherTool, pid,
				fmt.Sprintf("%s/%s", JVMContainerDir, JVMAgentTool), getContainerRuleFile(pid))
			if _, err := cmdexec.ExecContainer(ctx, cr, cId, []string{namespace.MNT, namespace.PID, namespace.IPC, namespace.ENV}, cmd, cmdexec.ExecStart); err != nil {
				return fmt.Errorf("execute fault for process[%d] error: %s", pid, err.Error())
			}
		}
	} else {
		for _, pid := range pidList {
			cmd := fmt.Sprintf("cd %s && java -cp .:tools.jar %s %d %s %s", utils.GetToolDir(), AttacherTool, pid, utils.GetToolPath(JVMAgentTool), getRuleFile("", pid))
			if _, err := cmdexec.StartBashCmdAndWaitPid(ctx, cmd, TimeoutSec); err != nil {
				return fmt.Errorf("execute fault for process[%d] error: %s", pid, err.Error())
			}
		}
	}

	return nil
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
