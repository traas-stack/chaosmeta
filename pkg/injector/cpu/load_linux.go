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
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
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

func (i *LoadInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Count, "count", "c", 0, "cpu load value（default 0, mean: cpu core num * 4）")
}

func (i *LoadInjector) getCmdExecutor() *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.PID},
	}
}

func (i *LoadInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	cpuList, err := getAllCpuList(ctx, i.Info.ContainerRuntime, i.Info.ContainerId)
	if err != nil {
		return fmt.Errorf("get all available cpu list error: %s", err.Error())
	}

	if i.Args.Count == 0 {
		i.Args.Count = len(cpuList) * 4
	}

	if i.Args.Count < 0 {
		return fmt.Errorf("\"count\"[%d] can not less than 0", i.Args.Count)
	}

	return nil
}

func (i *LoadInjector) Inject(ctx context.Context) error {
	cmd := fmt.Sprintf("%s %s %d", utils.GetToolPath(CpuLoadKey), i.Info.Uid, i.Args.Count)
	if err :=i.getCmdExecutor().StartCmd(ctx, cmd); err != nil {
		if err := i.Recover(ctx); err != nil {
			log.GetLogger(ctx).Warnf("undo error: %s", err.Error())
		}

		return err
	}

	return nil
}

func (i *LoadInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return process.CheckExistAndKillByKey(ctx, fmt.Sprintf("%s %s", CpuLoadKey, i.Info.Uid))
}
