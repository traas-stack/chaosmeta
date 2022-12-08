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
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"os"
	"strconv"
)

var (
	fileChmodFileName = "chaosmeta_file.test"
	fileChmodPerm     = "762"
	fileChmodOldPerm  = ""
)

func GetFileChmodTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p /notexistpath/temp.log -P %s", fileChmodPerm),
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p tempdir -P %s", fileChmodPerm),
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, "mkdir tempdir")
			},
			PostProcessor: func() error {
				return os.Remove("tempdir")
			},
		},
		{
			Args:  fmt.Sprintf("-p %s -P %s", fileChmodFileName, "799"),
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileChmodFileName))
			},
			PostProcessor: func() error {
				return os.Remove(fileChmodFileName)
			},
		},
		{
			Args:  fmt.Sprintf("-p %s -P %s", fileChmodFileName, "reg34t"),
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileChmodFileName))
			},
			PostProcessor: func() error {
				return os.Remove(fileChmodFileName)
			},
		},
		{
			Args: fmt.Sprintf("-p %s -P %s", fileChmodFileName, fileChmodPerm),
			PreProcessor: func() error {
				if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileChmodFileName)); err != nil {
					return err
				}

				var err error
				fileChmodOldPerm, err = filesys.GetPermission(fileChmodFileName)
				if err != nil {
					return err
				}

				fmt.Println(fileChmodOldPerm)
				return nil
			},
			Check: func() error {
				return checkPerm(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileChmodFileName), fileChmodPerm)
			},
			CheckRecover: func() error {
				return checkPerm(fmt.Sprintf("%s/%s", utils.GetRunPath(), fileChmodFileName), fileChmodOldPerm)
			},
			PostProcessor: func() error {
				return os.Remove(fileChmodFileName)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "file"
		tempCaseList[i].Fault = "chmod"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkPerm(fileName, targetPerm string) error {
	perm, err := filesys.GetPermission(fileName)
	if err != nil {
		return fmt.Errorf("get perm error: %s", err.Error())
	}

	if perm != targetPerm {
		return fmt.Errorf("expected perm: %s, actually: %s", targetPerm, perm)
	}

	return nil
}
