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
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
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

func (i *AddInjector) Validator() error {
	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	isPathExist, err := utils.ExistPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", i.Args.Path, err.Error())
	}

	dir := filepath.Dir(i.Args.Path)

	isDirExist, err := utils.ExistPath(dir)
	if err != nil {
		return fmt.Errorf("check dir[%s] exist error: %s", dir, err.Error())
	}

	if isPathExist {
		isFile, _ := utils.ExistFile(i.Args.Path)
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
		if err := utils.CheckPermission(i.Args.Permission); err != nil {
			return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
		}
	}

	return i.BaseInjector.Validator()
}

func (i *AddInjector) Inject() error {
	dir := filepath.Dir(i.Args.Path)
	isDirExist, _ := utils.ExistPath(dir)

	if !isDirExist {
		if err := utils.RunBashCmdWithoutOutput(fmt.Sprintf("mkdir -p %s", dir)); err != nil {
			return fmt.Errorf("mkdir [%s] error: %s", dir, err.Error())
		}
	}

	var chmodCmd string
	if i.Args.Permission != "" {
		chmodCmd = fmt.Sprintf(" && chmod %s %s", i.Args.Permission, i.Args.Path)
	}
	if err := utils.RunBashCmdWithoutOutput(fmt.Sprintf("echo -e \"%s\" > %s %s", i.Args.Content, i.Args.Path, chmodCmd)); err != nil {
		return fmt.Errorf("add content to %s error: %s", i.Args.Path, err.Error())
	}

	return nil
}

func (i *AddInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	isExist, err := utils.ExistPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", i.Args.Path, err.Error())
	}

	if isExist {
		return utils.RunBashCmdWithoutOutput(fmt.Sprintf("rm -rf %s", i.Args.Path))
	}

	return nil
}
