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
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/cpu"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/disk"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/diskio"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/file"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/kernel"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/mem"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/network"
	"github.com/ChaosMetaverse/chaosmetad/test/testcase/process"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//var testCaseMap = map[string]func() []common.TestCase{
//	"mem-fill":       mem.GetMemFillTest,
//	"cpu-burn":       cpu.GetCpuBurnTest,
//	"cpu-load":       cpu.GetCpuLoadTest,
//	"disk-fill":      disk.GetDiskFillTest,
//	"diskio-burn":    diskio.GetDiskIOBurnTest,
//	"network-occupy": network.GetNetOccupyTest,
//	"kernel-nproc":   kernel.GetNProcTest,
//	"kernel-fd":      kernel.GetFdTest,
//	"file-add":       file.GetFileAddTest,
//	"file-append":    file.GetFileAppendTest,
//	"file-chmod":     file.GetFileChmodTest,
//	"file-mv":        file.GetFileMvTest,
//	"file-delete":    file.GetFileDeleteTest,
//	"process-stop":   process.GetProStopTest,
//	"process-kill":   process.GetProKillTest,
//}

func getTestCases() []common.TestCase {
	var testCases []common.TestCase
	testCases = append(testCases, mem.GetMemFillTest()...)
	testCases = append(testCases, cpu.GetCpuBurnTest()...)
	testCases = append(testCases, cpu.GetCpuLoadTest()...)
	testCases = append(testCases, disk.GetDiskFillTest()...)
	testCases = append(testCases, diskio.GetDiskIOBurnTest()...)
	testCases = append(testCases, network.GetNetOccupyTest()...)
	testCases = append(testCases, kernel.GetNProcTest()...)
	testCases = append(testCases, kernel.GetFdTest()...)
	testCases = append(testCases, file.GetFileAddTest()...)
	testCases = append(testCases, file.GetFileAppendTest()...)
	testCases = append(testCases, file.GetFileChmodTest()...)
	testCases = append(testCases, file.GetFileMvTest()...)
	testCases = append(testCases, file.GetFileDeleteTest()...)
	testCases = append(testCases, process.GetProStopTest()...)
	testCases = append(testCases, process.GetProKillTest()...)

	return testCases
}

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

	var testCases = getTestCases()
	runTestCases(tool, testCases)
	//for count, t := range testCases {
	//	fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@ %d %s %s @@@@@@@@@@@@@@@@@@@@@@@\n", count, t.Target, t.Fault)
	//	runTestCases(tool, )
	//}
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
