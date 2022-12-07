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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
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

func (i *StopInjector) Validator() error {
	if i.Args.Pid < 0 {
		return fmt.Errorf("\"pid\" can not less than 0")
	}

	if i.Args.Pid == 0 && i.Args.Key == "" {
		return fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if i.Args.Pid > 0 {
		exist, err := process.ExistPid(i.Args.Pid)
		if err != nil {
			return fmt.Errorf("check pid[%d] exist error: %s", i.Args.Pid, err.Error())
		}

		if !exist {
			return fmt.Errorf("pid[%d] not exist", i.Args.Pid)
		}
	} else {
		exist, err := process.ExistProcessByKey(i.Args.Key)
		if err != nil {
			return fmt.Errorf("check pid by key[%s] error: %s", i.Args.Key, err.Error())
		}

		if !exist {
			return fmt.Errorf("no process grep by key[%s]", i.Args.Key)
		}
	}

	return i.BaseInjector.Validator()
}

func (i *StopInjector) Inject() error {
	if i.Args.Pid > 0 {
		if err := process.KillPidWithSignal(i.Args.Pid, process.SIGSTOP); err != nil {
			return err
		}
	} else {
		if err := process.KillProcessByKey(i.Args.Key, process.SIGSTOP); err != nil {
			return err
		}
	}

	return nil
}

func (i *StopInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	if i.Args.Pid > 0 {
		if err := process.KillPidWithSignal(i.Args.Pid, process.SIGCONT); err != nil {
			return err
		}
	} else {
		if err := process.KillProcessByKey(i.Args.Key, process.SIGCONT); err != nil {
			return err
		}
	}

	return nil
}
