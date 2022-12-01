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
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

func init() {
	injector.Register(TargetFile, FaultFileAppend, func() injector.IInjector { return &AppendInjector{} })
}

type AppendInjector struct {
	injector.BaseInjector
	Args    AppendArgs
	Runtime AppendRuntime
}

type AppendArgs struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
	Raw     bool   `json:"raw,omitempty"`
}

type AppendRuntime struct {
}

func (i *AppendInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *AppendInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *AppendInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Path, "path", "p", "", "file path, include dir and file name")
	cmd.Flags().StringVarP(&i.Args.Content, "content", "c", "", "append content to the existed file")
	cmd.Flags().BoolVarP(&i.Args.Raw, "raw", "r", false, "if raw content, raw content can not recover")
}

func (i *AppendInjector) Validator() error {
	if i.Args.Path == "" {
		return fmt.Errorf("\"path\" is empty")
	}

	var err error
	i.Args.Path, err = filepath.Abs(i.Args.Path)
	if err != nil {
		return fmt.Errorf("get absolute path of path[%s] error: %s", i.Args.Path, err.Error())
	}

	isFileExist, err := utils.ExistFile(i.Args.Path)
	if err != nil {
		return fmt.Errorf("\"path\"[%s] check exist error: %s", i.Args.Path, err.Error())
	}

	if !isFileExist {
		return fmt.Errorf("file[%s] is not exist", i.Args.Path)
	}

	if i.Args.Content == "" {
		return fmt.Errorf("\"content\" is empty")
	}

	return i.BaseInjector.Validator()
}

func getAppendFlag(uid string) string {
	return fmt.Sprintf(" %s-%s", utils.RootName, uid)
}

func (i *AppendInjector) Inject() error {
	content := i.Args.Content
	flag := getAppendFlag(i.Info.Uid)

	if !i.Args.Raw {
		content = strings.ReplaceAll(content, "\\n", "\n")
		content = fmt.Sprintf("%s%s", strings.ReplaceAll(content, "\n", fmt.Sprintf("%s\n", flag)), flag)
	}

	content = fmt.Sprintf("\n%s", content)
	if err := utils.RunBashCmdWithoutOutput(fmt.Sprintf("echo -e \"%s\" >> %s", content, i.Args.Path)); err != nil {
		return fmt.Errorf("append content to %s error: %s", i.Args.Path, err.Error())
	}

	return nil
}

func (i *AppendInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	if i.Args.Raw {
		return nil
	}

	fileExist, err := utils.ExistFile(i.Args.Path)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", i.Args.Path, err.Error())
	}

	if !fileExist {
		return nil
	}

	flag := getAppendFlag(i.Info.Uid)
	isExist, err := utils.HasFileLineByKey(flag, i.Args.Path)
	if err != nil {
		return fmt.Errorf("check file[%s] line exist key[%s] error: %s", i.Args.Path, flag, err.Error())
	}

	if isExist {
		return utils.RunBashCmdWithoutOutput(fmt.Sprintf("sed -i '/%s/d' %s", getAppendFlag(i.Info.Uid), i.Args.Path))
	}

	return nil
}
