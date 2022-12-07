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
	"os"
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

func (i *DeleteInjector) Validator() error {
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

	return i.BaseInjector.Validator()
}

func (i *DeleteInjector) getBackupDir() string {
	return fmt.Sprintf("%s%s", BackUpDir, i.Info.Uid)
}

func (i *DeleteInjector) Inject() error {
	backupDir := i.getBackupDir()
	if err := filesys.MkdirP(backupDir); err != nil {
		return fmt.Errorf("create backup dir[%s] error: %s", backupDir, err.Error())
	}

	return os.Rename(i.Args.Path, fmt.Sprintf("%s/%s", backupDir, filepath.Base(i.Args.Path)))
}

func (i *DeleteInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}
	backupDir := i.getBackupDir()

	isExist, err := filesys.ExistPath(i.Args.Path)
	if err != nil {
		return fmt.Errorf("check path[%s] exist error: %s", i.Args.Path, err.Error())
	}
	if !isExist {
		backupFile := fmt.Sprintf("%s/%s", backupDir, filepath.Base(i.Args.Path))
		if err := os.Rename(backupFile, i.Args.Path); err != nil {
			return fmt.Errorf("mv from[%s] to[%s] error: %s", backupFile, i.Args.Path, err.Error())
		}
	}

	return os.Remove(backupDir)
}
