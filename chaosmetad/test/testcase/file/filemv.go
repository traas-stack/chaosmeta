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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/test/common"
	"os"
	"strconv"
)

var (
	fileMvSrcFileName = "chaosmeta_src.test"
	fileMvDstFileName = "chaosmeta_dst.test"
)

func GetFileMvTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  "-s /fg3g -d wfq",
			Error: true,
		},
		{
			Args:  "-s tempdir",
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, "mkdir tempdir")
			},
			PostProcessor: func() error {
				return os.Remove("tempdir")
			},
		},
		{
			Args:  fmt.Sprintf("-s %s", fileMvSrcFileName),
			Error: true,
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileMvSrcFileName))
			},
			PostProcessor: func() error {
				return os.Remove(fileMvSrcFileName)
			},
		},
		{
			Args: fmt.Sprintf("-s %s -d %s", fileMvSrcFileName, fileMvDstFileName),
			PreProcessor: func() error {
				return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("touch %s", fileMvSrcFileName))
			},
			PostProcessor: func() error {
				return os.Remove(fileMvSrcFileName)
			},
			Check: func() error {
				return checkMv(fileMvSrcFileName, fileMvDstFileName)
			},
			CheckRecover: func() error {
				return checkMv(fileMvDstFileName, fileMvSrcFileName)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "file"
		tempCaseList[i].Fault = "mv"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkMv(src, dst string) error {
	exist, err := filesys.ExistFile(dst)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", dst, err.Error())
	}

	if !exist {
		return fmt.Errorf("dst file[%s] is not exist", dst)
	}

	exist, err = filesys.ExistFile(src)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", src, err.Error())
	}

	if exist {
		return fmt.Errorf("src file[%s] is still exist", src)
	}

	return nil
}
