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

func (i *MvInjector) Validator() error {
	if i.Args.Src == "" {
		return fmt.Errorf("\"src\" is empty")
	}

	var err error
	i.Args.Src, err = filepath.Abs(i.Args.Src)
	if err != nil {
		return fmt.Errorf("get absolute path of src[%s] error: %s", i.Args.Src, err.Error())
	}

	isPathExist, err := filesys.ExistFile(i.Args.Src)
	if err != nil {
		return fmt.Errorf("\"src\"[%s] check exist error: %s", i.Args.Src, err.Error())
	}

	if !isPathExist {
		return fmt.Errorf("\"src\"[%s] is not an existed file", i.Args.Src)
	}

	if i.Args.Dst == "" {
		return fmt.Errorf("\"dst\" is empty")
	}

	i.Args.Dst, err = filepath.Abs(i.Args.Dst)
	if err != nil {
		return fmt.Errorf("get absolute path of dst[%s] error: %s", i.Args.Dst, err.Error())
	}

	isPathExist, err = filesys.ExistPath(i.Args.Dst)
	if err != nil {
		return fmt.Errorf("\"dst\"[%s] check exist error: %s", i.Args.Dst, err.Error())
	}

	if isPathExist {
		return fmt.Errorf("\"dst\"[%s] is existed", i.Args.Dst)
	}

	return i.BaseInjector.Validator()
}

// Inject TODO: Consider whether to add a backup operation, copy first and then move
func (i *MvInjector) Inject() error {
	return os.Rename(i.Args.Src, i.Args.Dst)
}

func (i *MvInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	isExist, err := filesys.ExistPath(i.Args.Src)
	if err != nil {
		return fmt.Errorf("check src[%s] exist error: %s", i.Args.Src, err.Error())
	}

	if isExist {
		return nil
	}

	return os.Rename(i.Args.Dst, i.Args.Src)
}
