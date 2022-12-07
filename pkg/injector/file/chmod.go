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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/spf13/cobra"
	"path/filepath"
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

func (i *ChmodInjector) Validator() error {
	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	var err error
	i.Args.Path, err = filepath.Abs(i.Args.Path)
	if err != nil {
		return fmt.Errorf("get absolute path of path[%s] error: %s", i.Args.Path, err.Error())
	}

	isPathExist, err := filesys.ExistFile(i.Args.Path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", i.Args.Path, err.Error())
	}

	if !isPathExist {
		return fmt.Errorf("\"path\"[%s] is not an existed file", i.Args.Path)
	}

	if i.Args.Permission != "" {
		if err := filesys.CheckPermission(i.Args.Permission); err != nil {
			return fmt.Errorf("\"permission\" is invalid: %s", err.Error())
		}
	}

	return i.BaseInjector.Validator()
}

func (i *ChmodInjector) Inject() error {
	perm, err := filesys.GetPermission(i.Args.Path)
	if err != nil {
		return fmt.Errorf("get perm of path[%s] error: %s", i.Args.Path, err.Error())
	}

	i.Runtime.Permission = perm
	return filesys.Chmod(i.Args.Path, i.Args.Permission)
}

func (i *ChmodInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	isExist, err := filesys.ExistFile(i.Args.Path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", i.Args.Path, err.Error())
	}

	if isExist {
		return filesys.Chmod(i.Args.Path, i.Runtime.Permission)
	}

	return nil
}
