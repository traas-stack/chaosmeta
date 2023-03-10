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
	"github.com/traas-stack/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmetad/pkg/utils/namespace"
	"path/filepath"
)

func init() {
	injector.Register(TargetFile, FaultFileDelete, func() injector.IInjector { return &DeleteInjector{} })
}

type DeleteInjector struct {
	injector.BaseInjector
	Args    DeleteArgs
	Runtime DeleteRuntime
}

type DeleteArgs struct {
	Path string `json:"path"`
}

type DeleteRuntime struct {
}

func (i *DeleteInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *DeleteInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *DeleteInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Path, "path", "p", "", "file path, include dir and file name")
}

func (i *DeleteInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT},
		ToolKey:          FileExec,
		Method:           method,
		Fault:            FaultFileDelete,
		Args:             args,
	}
}

func (i *DeleteInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	if i.Info.ContainerRuntime != "" {
		if !filepath.IsAbs(i.Args.Path) {
			return fmt.Errorf("\"path\" must provide absolute path")
		}
	} else {
		var err error
		i.Args.Path, err = filesys.GetAbsPath(i.Args.Path)
		if err != nil {
			return fmt.Errorf("\"path\"[%s] get absolute path error: %s", i.Args.Path, err.Error())
		}
	}

	return i.getCmdExecutor(utils.MethodValidator, i.Args.Path).ExecTool(ctx)
}

func (i *DeleteInjector) Inject(ctx context.Context) error {
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%s %s", i.Args.Path, i.Info.Uid)).ExecTool(ctx)
}

func (i *DeleteInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%s %s", i.Args.Path, i.Info.Uid)).ExecTool(ctx)
}
