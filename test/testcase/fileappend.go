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

package testcase

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	fileAppendFileName = "chaosmeta_file.test"
	fileInitContent    = "init\n123\ng34g"
	fileAppendContent  = "wdvew\nregewc24\ncf234c\neg34"
)

func GetFileAppendTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p /notexistpath/temp.log -c \"%s\"", fileAppendContent),
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p tempdir -c \"%s\"", fileAppendContent),
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, "mkdir tempdir")
			},
			PostProcessor: func() error {
				return os.Remove("tempdir")
			},
		},
		{
			Args: fmt.Sprintf("-p %s -c \"%s\"", fileAppendFileName, fileAppendContent),
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo -en \"%s\" > %s", fileInitContent, fileAppendFileName))
			},
			Check: func() error {
				return checkAppend(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileAppendFileName), 3, 4, true)
			},
			CheckRecover: func() error {
				return checkAppend(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileAppendFileName), 3, 0, false)
			},
			PostProcessor: func() error {
				return os.Remove(fileAppendFileName)
			},
		},
		{
			Args: fmt.Sprintf("-p %s -c \"%s\" -r", fileAppendFileName, fileAppendContent),
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo -en \"%s\" > %s", fileInitContent, fileAppendFileName))
			},
			Check: func() error {
				return checkAppend(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileAppendFileName), 3, 4, false)
			},
			CheckRecover: func() error {
				return checkAppend(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileAppendFileName), 3, 4, false)
			},
			PostProcessor: func() error {
				return os.Remove(fileAppendFileName)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "file"
		tempCaseList[i].Fault = "append"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkAppend(fileName string, initCount, appendCount int, flag bool) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file[%s] error: %s", fileName, err.Error())
	}

	contentStr := string(content)
	fmt.Println(contentStr)
	arr := strings.Split(contentStr, "\n")

	expectLine := initCount + appendCount + 1
	if len(arr) != expectLine {
		return fmt.Errorf("expected line: %d, actually: %d", expectLine, len(arr))
	}

	if flag {
		flagCount := strings.Count(contentStr, "chaosmetad-")
		if flagCount != appendCount {
			return fmt.Errorf("expected flag count: %d, actually: %d", appendCount, flagCount)
		}
	}

	return nil
}
