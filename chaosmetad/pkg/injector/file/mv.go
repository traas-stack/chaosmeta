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
	injector.Register(TargetFile, FaultFileMv, func() injector.IInjector { return &MvInjector{} })
}

type MvInjector struct {
	injector.BaseInjector
	Args    MvArgs
	Runtime MvRuntime
}

type MvArgs struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type MvRuntime struct {
}

func (i *MvInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *MvInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *MvInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Src, "src", "s", "", "source file path, include dir and file name")
	cmd.Flags().StringVarP(&i.Args.Dst, "dst", "d", "", "destination file path, include dir and file name")
}

func (i *MvInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Src == "" {
		return fmt.Errorf("\"src\" is empty")
	}

	if i.Args.Dst == "" {
		return fmt.Errorf("\"dst\" is empty")
	}

	if !filesys.IfPathAbs(ctx, i.Args.Src) {
		return fmt.Errorf("\"src\" must provide absolute path")
	}

	if !filesys.IfPathAbs(ctx, i.Args.Dst) {
		return fmt.Errorf("\"dst\" must provide absolute path")
	}

	srcExist, err := filesys.CheckFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Src)
	if err != nil {
		return fmt.Errorf("check exist file[%s] error: %s", i.Args.Src, err.Error())
	}

	dstExist, err := filesys.ExistPath(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Dst)
	if err != nil {
		return fmt.Errorf("check exist path[%s] error: %s", i.Args.Dst, err.Error())
	}

	if !srcExist {
		return fmt.Errorf("source file[%s] is not exist", i.Args.Src)
	}

	if dstExist {
		return fmt.Errorf("dst path[%s] is exist", i.Args.Dst)
	}

	return nil
}

// Inject TODO: Consider whether to add a backup operation, copy first and then move
func (i *MvInjector) Inject(ctx context.Context) error {
	return filesys.MoveFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Src, i.Args.Dst)
}

func (i *MvInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	exist, err := filesys.ExistPath(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Src)
	if err != nil {
		return fmt.Errorf("check exist path[%s] error: %s", i.Args.Src, err.Error())
	}

	if !exist {
		return filesys.MoveFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Dst, i.Args.Src)
	}

	return nil
}
