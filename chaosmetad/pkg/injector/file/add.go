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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
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

func (i *AddInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" can not be empty")
	}

	if !filesys.IfPathAbs(ctx, i.Args.Path) {
		return fmt.Errorf("\"path\" must provide absolute path")
	}

	if i.Args.Permission != "" {
		if err := filesys.CheckPermission(i.Args.Permission); err != nil {
			return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
		}
	}

	isPathExist, err := filesys.ExistPath(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check exist path[%s] error: %s", i.Args.Path, err.Error())
	}

	dir := filepath.Dir(i.Args.Path)
	isDir, err := filesys.CheckDir(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, dir)
	if err != nil {
		return fmt.Errorf("check dir[%s] error: %s", dir, err.Error())
	}

	if isPathExist {
		checkDir, err := filesys.CheckDir(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
		if err != nil {
			return fmt.Errorf("check file[%s] error: %s", i.Args.Path, err.Error())
		}

		if checkDir {
			return fmt.Errorf("\"path\"[%s] is an existed dir", i.Args.Path)
		}
	}

	if !i.Args.Force {
		if isPathExist {
			return fmt.Errorf("file[%s] exist, if want to force to overwrite, please provide [-f] or [--force] args", i.Args.Path)
		}

		if !isDir {
			return fmt.Errorf("dir[%s] is not exist, if want to auto create, please provide [-f] or [--force] args", dir)
		}
	}

	return nil
}

func (i *AddInjector) Inject(ctx context.Context) error {
	dir := filepath.Dir(i.Args.Path)
	isDirExist, err := filesys.CheckDir(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, dir)
	if err != nil {
		return fmt.Errorf("check dir[%s] error: %s", i.Args.Path, err.Error())
	}

	if !isDirExist {
		if err := filesys.MkdirForce(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, dir); err != nil {
			return fmt.Errorf("mkdir dir[%s] error: %s", dir, err.Error())
		}
	}

	if err := filesys.OverWriteFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path, i.Args.Content); err != nil {
		return fmt.Errorf("add content to %s error: %s", i.Args.Path, err.Error())
	}

	if i.Args.Permission != "" {
		if err := filesys.Chmod(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path, i.Args.Permission); err != nil {
			if err := i.Recover(ctx); err != nil {
				log.GetLogger(ctx).Warnf("undo error: %s", err.Error())
			}

			return fmt.Errorf("chmod file[%s] to[%s] error: %s", i.Args.Path, i.Args.Permission, err.Error())
		}
	}

	return nil
}

func (i *AddInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	isExist, err := filesys.CheckFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check exist file[%s] error: %s", i.Args.Path, err.Error())
	}

	if isExist {
		return filesys.RemoveFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	}

	return nil
}
