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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
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

func (i *ChmodInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	if i.Args.Permission == "" {
		return fmt.Errorf("\"permission\" is empty")
	}

	if !filesys.IfPathAbs(ctx, i.Args.Path) {
		return fmt.Errorf("\"path\" must provide absolute path")
	}

	if err := filesys.CheckPermission(i.Args.Permission); err != nil {
		return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
	}

	exist, err := filesys.CheckFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check exist file[%s] error: %s", i.Args.Path, err.Error())
	}

	if !exist {
		return fmt.Errorf("file[%s] is not exist", i.Args.Path)
	}

	return nil
}

func (i *ChmodInjector) Inject(ctx context.Context) error {
	perm, err := filesys.GetPerm(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("get perm of path[%s] error: %s", i.Args.Path, err.Error())
	}

	i.Runtime.Permission = perm
	return filesys.Chmod(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path, i.Args.Permission)
}

func (i *ChmodInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	exist, err := filesys.CheckFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check exist file[%s] error: %s", i.Args.Path, err.Error())
	}

	if exist {
		return filesys.Chmod(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path, i.Runtime.Permission)
	}

	return nil
}
