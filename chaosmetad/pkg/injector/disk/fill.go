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

package disk

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
	"path/filepath"
)

func init() {
	injector.Register(TargetDisk, FaultDiskFill, func() injector.IInjector { return &FillInjector{} })
}

type FillInjector struct {
	injector.BaseInjector
	Args    FillArgs
	Runtime FillRuntime
}

type FillArgs struct {
	Percent int    `json:"percent,omitempty"`
	Bytes   string `json:"bytes,omitempty"`
	Dir     string `json:"dir,omitempty"`
}

type FillRuntime struct {
}

func (i *FillInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *FillInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *FillInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Dir == "" {
		i.Args.Dir = DefaultDir
	}
}

func (i *FillInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "disk fill target percent, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")
	cmd.Flags().StringVarP(&i.Args.Bytes, "bytes", "b", "", "disk fill bytes to add, support unit: KB/MB/GB/TB（default KB）")
	cmd.Flags().StringVarP(&i.Args.Dir, "dir", "d", "", fmt.Sprintf("disk fill target dir（default %s）", DefaultDir))
}

func (i *FillInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT},
		ToolKey:          DiskFillExec,
		Method:           method,
		Fault:            FaultDiskFill,
		Args:             args,
	}
}

func (i *FillInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Info.ContainerRuntime != "" {
		if !filepath.IsAbs(i.Args.Dir) {
			return fmt.Errorf("\"dir\" must provide absolute path")
		}
	} else {
		var err error
		i.Args.Dir, err = filesys.GetAbsPath(i.Args.Dir)
		if err != nil {
			return fmt.Errorf("\"dir\"[%s] get absolute path error: %s", i.Args.Dir, err.Error())
		}
	}

	return i.getCmdExecutor(utils.MethodValidator, fmt.Sprintf("%d '%s' %s", i.Args.Percent, i.Args.Bytes, i.Args.Dir)).ExecTool(ctx)
}

func (i *FillInjector) Inject(ctx context.Context) error {
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%d '%s' %s %s", i.Args.Percent, i.Args.Bytes, i.Args.Dir, i.Info.Uid)).ExecTool(ctx)
}

func (i *FillInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%s %s", i.Args.Dir, i.Info.Uid)).ExecTool(ctx)
}
