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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/test/common"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	fileAddFileName = "/tmp/chaosmeta_file.test"
	fileAddContent  = "wdvew\nregewc24\ncf234c\neg34"
	fileAddPerm     = "762"
)

func GetFileAddTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  "-p /notexistpath/temp.log",
			Error: true,
		},
		{
			Args:  "-p tempdir",
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, "mkdir tempdir")
			},
			PostProcessor: func() error {
				return os.Remove("tempdir")
			},
		},
		{
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, "touch ./tempfile")
			},
			Args:  "-p tempfile/temp.log",
			Error: true,
			PostProcessor: func() error {
				return os.RemoveAll("tempfile")
			},
		},
		{
			Args:  fmt.Sprintf("-p %s -c \"%s\" -P 759", fileAddFileName, fileAddContent),
			Error: true,
		},
		{
			Args: fmt.Sprintf("-p %s -c \"%s\" -P %s", fileAddFileName, fileAddContent, fileAddPerm),
			Check: func() error {
				return checkFileAdd(fileAddFileName, fileAddPerm, fileAddContent)
			},
			CheckRecover: func() error {
				return checkFileNotExist(fileAddFileName)
			},
		},
		{
			Args: fmt.Sprintf("-p %s -c \"%s\" -P %s -f", "/notexist/abc/efg/temp.log", fileAddContent, fileAddPerm),
			Check: func() error {
				return checkFileAdd("/notexist/abc/efg/temp.log", fileAddPerm, fileAddContent)
			},
			CheckRecover: func() error {
				return checkFileNotExist("/notexist/abc/efg/temp.log")
			},
			PostProcessor: func() error {
				return os.RemoveAll("/notexist")
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "file"
		tempCaseList[i].Fault = "add"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkFileAdd(fileName, addPerm, addContent string) error {
	exist, err := filesys.ExistFile(fileName)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", fileName, err.Error())
	}

	if !exist {
		return fmt.Errorf("add file[%s] failed: not exist", fileName)
	}

	perm, err := filesys.GetPerm(context.Background(), "", "", fileName)
	if err != nil {
		return fmt.Errorf("get perm of file[%s] error: %s", fileName, err.Error())
	}

	if perm != addPerm {
		return fmt.Errorf("expected perm: %s, actually %s", addPerm, perm)
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file[%s] error: %s", fileName, err.Error())
	}

	if string(data) != addContent {
		return fmt.Errorf("expected content: [%s], actually: [%s]", addContent, string(data))
	}

	return nil
}

func checkFileNotExist(file string) error {
	fileName := fmt.Sprintf("%s/%s", utils.GetRunPath(), file)
	exist, err := filesys.ExistFile(fileName)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", fileName, err.Error())
	}

	if exist {
		return fmt.Errorf("recover file[%s] failed: still exist", fileName)
	}

	return nil
}
