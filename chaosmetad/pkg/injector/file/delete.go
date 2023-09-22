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

func (i *DeleteInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	if !filesys.IfPathAbs(ctx, i.Args.Path) {
		return fmt.Errorf("\"path\" must provide absolute path")
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

func (i *DeleteInjector) Inject(ctx context.Context) error {
	backupDir := getBackupDir(i.Info.Uid)
	if err := filesys.MkdirForce(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, backupDir); err != nil {
		return fmt.Errorf("create backup dir[%s] error: %s", backupDir, err.Error())
	}

	return filesys.MoveFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path, fmt.Sprintf("%s/%s", backupDir, filepath.Base(i.Args.Path)))
}

func (i *DeleteInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	backupFile := fmt.Sprintf("%s/%s", getBackupDir(i.Info.Uid), filepath.Base(i.Args.Path))
	isExist, err := filesys.CheckFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check exist file[%s] error: %s", i.Args.Path, err.Error())
	}

	if !isExist {
		if err := filesys.MoveFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, backupFile, i.Args.Path); err != nil {
			return fmt.Errorf("mv from[%s] to[%s] error: %s", backupFile, i.Args.Path, err.Error())
		}
	}

	return filesys.RemoveRF(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, getBackupDir(i.Info.Uid))
}
