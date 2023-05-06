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

package file

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
	injector.Register(TargetFile, FaultFileMv, func() injector.IInjector { return &MvInjector{} })
}

type MvInjector struct {
	injector.BaseInjector
	Args    MvArgs
	Runtime MvRuntime
}

type MvArgs struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type MvRuntime struct {
}

func (i *MvInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *MvInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *MvInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Src, "src", "s", "", "source file path, include dir and file name")
	cmd.Flags().StringVarP(&i.Args.Dst, "dst", "d", "", "destination file path, include dir and file name")
}

func (i *MvInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT},
		ToolKey:          FileExec,
		Method:           method,
		Fault:            FaultFileMv,
		Args:             args,
	}
}

func (i *MvInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Src == "" {
		return fmt.Errorf("\"src\" is empty")
	}

	if i.Args.Dst == "" {
		return fmt.Errorf("\"dst\" is empty")
	}

	if i.Info.ContainerRuntime != "" {
		if !filepath.IsAbs(i.Args.Src) {
			return fmt.Errorf("\"src\" must provide absolute path")
		}

		if !filepath.IsAbs(i.Args.Dst) {
			return fmt.Errorf("\"dst\" must provide absolute path")
		}
	} else {
		var err error
		i.Args.Src, err = filesys.GetAbsPath(i.Args.Src)
		if err != nil {
			return fmt.Errorf("\"src\"[%s] get absolute path error: %s", i.Args.Src, err.Error())
		}

		i.Args.Dst, err = filesys.GetAbsPath(i.Args.Dst)
		if err != nil {
			return fmt.Errorf("\"dst\"[%s] get absolute path error: %s", i.Args.Dst, err.Error())
		}
	}

	return i.getCmdExecutor(utils.MethodValidator, fmt.Sprintf("%s %s", i.Args.Src, i.Args.Dst)).ExecTool(ctx)
}

// Inject TODO: Consider whether to add a backup operation, copy first and then move
func (i *MvInjector) Inject(ctx context.Context) error {
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%s %s", i.Args.Src, i.Args.Dst)).ExecTool(ctx)
}

func (i *MvInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%s %s", i.Args.Src, i.Args.Dst)).ExecTool(ctx)
}
