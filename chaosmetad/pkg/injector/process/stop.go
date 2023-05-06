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

package process

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
)

func init() {
	injector.Register(TargetProcess, FaultProcessStop, func() injector.IInjector { return &StopInjector{} })
}

type StopInjector struct {
	injector.BaseInjector
	Args    StopArgs
	Runtime StopRuntime
}

type StopArgs struct {
	Pid int    `json:"pid,omitempty"`
	Key string `json:"key,omitempty"`
}

type StopRuntime struct {
}

func (i *StopInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *StopInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *StopInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Pid, "pid", "p", 0, "target process's pid")
	cmd.Flags().StringVarP(&i.Args.Key, "key", "k", "", "the key used to grep to get target process, the effect is equivalent to \"ps -ef | grep [key]\". if \"pid\" provided, \"key\" will be ignored")
}

func (i *StopInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT, namespace.PID},
		ToolKey:          ProcessExec,
		Method:           method,
		Fault:            FaultProcessStop,
		Args:             args,
	}
}

func (i *StopInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}
	return i.getCmdExecutor(utils.MethodValidator, fmt.Sprintf("%d '%s'", i.Args.Pid, i.Args.Key)).ExecTool(ctx)
}

func (i *StopInjector) Inject(ctx context.Context) error {
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%d '%s'", i.Args.Pid, i.Args.Key)).ExecTool(ctx)
}

func (i *StopInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%d '%s'", i.Args.Pid, i.Args.Key)).ExecTool(ctx)
}
