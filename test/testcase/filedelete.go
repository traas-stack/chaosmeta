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
	file2 "github.com/ChaosMetaverse/chaosmetad/pkg/injector/file"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"os"
	"strconv"
)

var (
	fileDeleteFileName = "chaosmeta_file.test"
)

func GetFileDeleteTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  "-p /fg3g",
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
			Args: fmt.Sprintf("-p %s", fileDeleteFileName),
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileDeleteFileName))
			},
			PostProcessor: func() error {
				return os.Remove(fileDeleteFileName)
			},
			Check: func() error {
				return checkDelete(fileDeleteFileName, false)
			},
			CheckRecover: func() error {
				return checkDelete(fileDeleteFileName, true)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "file"
		tempCaseList[i].Fault = "del"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkDelete(file string, ifExist bool) error {
	exist, err := filesys.ExistFile(file)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", file, err.Error())
	}

	if exist != ifExist {
		return fmt.Errorf("expected exist status: %v, actually: %v", ifExist, exist)
	}

	backupFile := fmt.Sprintf("%s%s/%s", file2.BackUpDir, common.UID, file)
	exist, err = filesys.ExistFile(backupFile)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", backupFile, err.Error())
	}

	if exist == ifExist {
		return fmt.Errorf("expected backup exist status: %v, actually: %v", !ifExist, exist)
	}

	return nil
}
