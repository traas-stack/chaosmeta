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

package main

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: Currently, only the P0 use case (single fault execution) is supported, the non-root user and concurrent faults are not added, the use case of automatic recovery, and the use case of manually recovering and then performing recover
func main() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		common.ExitErr(fmt.Sprintf("get path error: %s", err.Error()))
	}

	pathArr := strings.Split(path, "/")
	rootPath := strings.Join(pathArr[:len(pathArr)-1], "/")
	tool := fmt.Sprintf("%s/build/chaosmetad/chaosmetad", rootPath)
	fmt.Println(tool)

	var testCases []common.TestCase
	fmt.Println("@@@@@@@@@@@@@@@CPU BURN@@@@@@@@@@@@@@@")
	testCases = testcase.GetCpuBurnTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@CPU LOAD@@@@@@@@@@@@@@@")
	testCases = testcase.GetCpuLoadTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@DISK FILL@@@@@@@@@@@@@@@")
	testCases = testcase.GetDiskFillTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@MEM FILL@@@@@@@@@@@@@@@")
	testCases = testcase.GetMemFillTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@PROCESS STOP@@@@@@@@@@@@@@@")
	testCases = testcase.GetProStopTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@PROCESS KILL@@@@@@@@@@@@@@@")
	testCases = testcase.GetProKillTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@KERNEL NPROC@@@@@@@@@@@@@@@")
	testCases = testcase.GetNProcTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@KERNEL FD@@@@@@@@@@@@@@@")
	testCases = testcase.GetFdTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@FILE ADD@@@@@@@@@@@@@@@")
	testCases = testcase.GetFileAddTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@FILE APPEND@@@@@@@@@@@@@@@")
	testCases = testcase.GetFileAppendTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@FILE CHMOD@@@@@@@@@@@@@@@")
	testCases = testcase.GetFileChmodTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@FILE MV@@@@@@@@@@@@@@@")
	testCases = testcase.GetFileMvTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@FILE DELETE@@@@@@@@@@@@@@@")
	testCases = testcase.GetFileDeleteTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@DISKIO BURN@@@@@@@@@@@@@@@")
	testCases = testcase.GetDiskIOBurnTest()
	runTestCases(tool, testCases)

	fmt.Println("@@@@@@@@@@@@@@@NETWORK OCCUPY@@@@@@@@@@@@@@@")
	testCases = testcase.GetNetOccupyTest()
	runTestCases(tool, testCases)
}

func runTestCases(tool string, testCases []common.TestCase) {
	for _, t := range testCases {
		fmt.Printf("===============CASE %s==============\n", t.Name)
		fmt.Println("***********PRE PROCESS***********")
		if t.PreProcessor != nil {
			if err := t.PreProcessor(); err != nil {
				common.ExitErr(fmt.Sprintf("pre process error: %s", err.Error()))
			}
		}
		fmt.Println("***********INJECT***********")
		injectCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s inject %s %s %s", tool, t.Target, t.Fault, t.Args))
		fmt.Printf("inject cmd: %s\n", injectCmd.Args)
		re, err := injectCmd.CombinedOutput()
		if err != nil {
			if !t.Error {
				common.ExitErr(fmt.Sprintf("exec unexpected error: %s, output: %s", err.Error(), string(re)))
			} else {
				fmt.Printf("exec error: %s, output: %s\n", err.Error(), string(re))
			}
		} else {
			fmt.Println(string(re))
			fmt.Println("***********CHECK***********")
			common.UID = common.GetUid(string(re))
			checkInfo := t.Check()
			fmt.Println("***********RECOVER***********")
			recoverCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s recover %s", tool, common.UID))
			fmt.Printf("recover cmd: %s\n", recoverCmd.Args)
			if re, err := recoverCmd.CombinedOutput(); err != nil {
				common.ExitErr(fmt.Sprintf("recover error: %s, outpur: %s", err.Error(), string(re)))
			}

			if err := t.CheckRecover(); err != nil {
				common.ExitErr(fmt.Sprintf("recover check error: %s", err.Error()))
			}

			if checkInfo != nil {
				common.ExitErr(fmt.Sprintf("check unexpected: %s", checkInfo.Error()))
			}
		}

		if t.PostProcessor != nil {
			if err := t.PostProcessor(); err != nil {
				common.ExitErr(fmt.Sprintf("post process error: %s", err.Error()))
			}
		}
	}
}
