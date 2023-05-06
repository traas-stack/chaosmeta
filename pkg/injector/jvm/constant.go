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

const (
	TargetJVM = "jvm"

	FaultMethodDelay     = "methoddelay"
	FaultMethodException = "methodexception"
	FaultMethodReturn    = "methodreturn"

	JVMRuleDir  = "jvm_rule"
	JVMRuleFile = "jvm_rule.json"

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
