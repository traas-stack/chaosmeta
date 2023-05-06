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
	injector.Register(TargetFile, FaultFileAdd, func() injector.IInjector { return &AddInjector{} })
}

type AddInjector struct {
	injector.BaseInjector
	Args    AddArgs
	Runtime AddRuntime
}

type AddArgs struct {
	Path       string `json:"path"`
	Content    string `json:"content,omitempty"`
	Permission string `json:"permission,omitempty"`
	Force      bool   `json:"force,omitempty"`
}

type AddRuntime struct {
}

func (i *AddInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *AddInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *AddInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Path, "path", "p", "", "file path, include dir and file name")
	cmd.Flags().StringVarP(&i.Args.Content, "content", "c", "", "add content to the new file")
	cmd.Flags().StringVarP(&i.Args.Permission, "permission", "P", "", "file's permission, compose format: three number in [0,7], example: 777")
	cmd.Flags().BoolVarP(&i.Args.Force, "force", "f", false, "if target dir not exist, will create. if target file exist, will overwrite")
}

func (i *AddInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT},
		ToolKey:          FileExec,
		Method:           method,
		Fault:            FaultFileAdd,
		Args:             args,
	}
}

func (i *AddInjector) Validator(ctx context.Context) error {
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

	if i.Args.Permission != "" {
		if err := filesys.CheckPermission(i.Args.Permission); err != nil {
			return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
		}
	}

	return i.getCmdExecutor(utils.MethodValidator, fmt.Sprintf("%s %v", i.Args.Path, i.Args.Force)).ExecTool(ctx)
}

func (i *AddInjector) Inject(ctx context.Context) error {
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%s '%s' '%s'", i.Args.Path, i.Args.Permission, i.Args.Content)).ExecTool(ctx)
}

func (i *AddInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, i.Args.Path).ExecTool(ctx)
}
