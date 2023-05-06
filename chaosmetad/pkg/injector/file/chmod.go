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
	"strings"
)

func init() {
	injector.Register(TargetFile, FaultFileChmod, func() injector.IInjector { return &ChmodInjector{} })
}

type ChmodInjector struct {
	injector.BaseInjector
	Args    ChmodArgs
	Runtime ChmodRuntime
}

type ChmodArgs struct {
	Path       string `json:"path"`
	Permission string `json:"permission,omitempty"`
	Force      bool   `json:"force,omitempty"`
}

type ChmodRuntime struct {
	Permission string `json:"permission,omitempty"`
}

func (i *ChmodInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *ChmodInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *ChmodInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Path, "path", "p", "", "file path, include dir and file name")
	cmd.Flags().StringVarP(&i.Args.Permission, "permission", "P", "", "file's permission, compose format: three number in [0,7], example: 777")
}

func (i *ChmodInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT},
		ToolKey:          FileExec,
		Method:           method,
		Fault:            FaultFileChmod,
		Args:             args,
	}
}

func (i *ChmodInjector) Validator(ctx context.Context) error {
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

	return i.getCmdExecutor(utils.MethodValidator, i.Args.Path).ExecTool(ctx)
}

func (i *ChmodInjector) Inject(ctx context.Context) error {
	perm, err := filesys.GetPerm(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("get perm of path[%s] error: %s", i.Args.Path, err.Error())
	}

	i.Runtime.Permission = strings.TrimSpace(perm)
	return i.getCmdExecutor(utils.MethodInject, fmt.Sprintf("%s %s", i.Args.Path, i.Args.Permission)).ExecTool(ctx)
}

func (i *ChmodInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%s %s", i.Args.Path, i.Runtime.Permission)).ExecTool(ctx)
}
