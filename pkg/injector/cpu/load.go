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

package cpu

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/spf13/cobra"
	"runtime"
)

// Register
func init() {
	injector.Register(TargetCpu, FaultCpuLoad, func() injector.IInjector { return &LoadInjector{} })
}

type LoadInjector struct {
	injector.BaseInjector
	Args    LoadArgs
	Runtime LoadRuntime
}

type LoadArgs struct {
	Count int `json:"count,omitempty"`
}

type LoadRuntime struct {
}

func (i *LoadInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *LoadInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *LoadInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Count == 0 {
		i.Args.Count = runtime.NumCPU() * 4
	}
}

func (i *LoadInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Count, "count", "c", 0, "cpu load value（default 0, mean: cpu core num * 4）")
}

func (i *LoadInjector) Validator() error {
	if i.Args.Count < 0 {
		return fmt.Errorf("\"count\"[%d] can not less than 0", i.Args.Count)
	}

	return i.BaseInjector.Validator()
}

func (i *LoadInjector) Inject() error {
	return utils.StartBashCmd(fmt.Sprintf("%s %s %d", utils.GetToolPath(CpuLoadKey), i.Info.Uid, i.Args.Count))
}

func (i *LoadInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	processKey := fmt.Sprintf("%s %s", CpuLoadKey, i.Info.Uid)
	isExist, err := utils.ExistProcessByKey(processKey)
	if err != nil {
		return fmt.Errorf("check process exist by key[%s] error: %s", processKey, err.Error())
	}

	if isExist {
		if err := utils.KillProcessByKey(processKey, utils.SIGKILL); err != nil {
			return fmt.Errorf("kill process by key[%s] error: %s", processKey, err.Error())
		}
	}

	return nil
}
