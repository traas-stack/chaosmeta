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
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/spf13/cobra"
	"os"
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
	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	var err error
	i.Args.Path, err = filesys.GetAbsPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("get absolute path of path[%s] error: %s", i.Args.Path, err.Error())
	}

	isPathExist, err := filesys.ExistPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", i.Args.Path, err.Error())
	}

	dir := filepath.Dir(i.Args.Path)
	isDirExist, err := filesys.ExistPath(dir)
	if err != nil {
		return fmt.Errorf("check dir[%s] exist error: %s", dir, err.Error())
	}

	if isPathExist {
		isFile, _ := filesys.ExistFile(i.Args.Path)
		if !isFile {
			return fmt.Errorf("\"path\"[%s] is an existed dir", i.Args.Path)
		}
	}

	if !i.Args.Force {
		if isPathExist {
			return fmt.Errorf("file[%s] exist, if want to force to overwrite, please provide [-f] or [--force] args", i.Args.Path)
		}

		if !isDirExist {
			return fmt.Errorf("dir[%s] is not exist, if want to auto create, please provide [-f] or [--force] args", dir)
		}
	}

	if i.Args.Permission != "" {
		if err := filesys.CheckPermission(i.Args.Permission); err != nil {
			return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
		}
	}

	return i.BaseInjector.Validator(ctx)
}

func (i *AddInjector) Inject(ctx context.Context) error {
	logger := log.GetLogger(ctx)

	dir := filepath.Dir(i.Args.Path)
	isDirExist, _ := filesys.ExistPath(dir)

	if !isDirExist {
		if err := filesys.MkdirP(ctx, dir); err != nil {
			return fmt.Errorf("mkdir dir[%s] error: %s", dir, err.Error())
		}
	}

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo -en \"%s\" > %s", i.Args.Content, i.Args.Path)); err != nil {
		return fmt.Errorf("add content to %s error: %s", i.Args.Path, err.Error())
	}

	if i.Args.Permission != "" {
		if err := filesys.Chmod(ctx, i.Args.Path, i.Args.Permission); err != nil {
			if err := i.Recover(ctx); err != nil {
				logger.Warnf("undo error: %s", err.Error())
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

	isExist, err := filesys.ExistPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", i.Args.Path, err.Error())
	}

	if isExist {
		return os.Remove(i.Args.Path)
	}

	return nil
}
