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
	"github.com/traas-stack/chaosmeta/chaosmetad/tools/common"
	"os"
	"os/exec"
	"strconv"
)

// dd if=/dev/zero of=/root/testfile bs=512 count=100000 oflag=dsync
// uid, file_path mode bs count flag timeout
var notFirst bool

func main() {
	args := os.Args
	if len(args) < 7 {
		common.ExitWithErr("args must at lease 7")
	}
	argsFile, argsMode, argsBs, argsCount, argsFlag, timeStr := args[2], args[3], args[4], args[5], args[6], args[7]

	timeout, err := strconv.Atoi(timeStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("args timeout is not a num: %s", err.Error()))
	}

	if argsMode == "write" {
		go burnWriteDisk(argsFile, argsBs, argsCount, argsFlag)
	} else if argsMode == "read" {
		go burnReadDisk(argsFile, argsBs, argsCount, argsFlag)
	} else {
		common.ExitWithErr("only support one of read or write flag")
	}

	common.SleepWait(timeout)
}

func burnReadDisk(file, argsBs, argsCount, flagArgs string) {
	//if err := utils.RunBashCmdWithoutOutput(fmt.Sprintf("dd if=/dev/zero of=%s bs=1M count=1024 oflag=direct", file)); err != nil {
	//	common.ExitWithErr(fmt.Sprintf("dd fill initial file error: %s", err.Error()))
	//}

	if err := exec.Command("dd", "if=/dev/zero", fmt.Sprintf("of=%s", file), "bs=1M", "count=1024", "oflag=direct").Run(); err != nil {
		common.ExitWithErr(fmt.Sprintf("dd fill initial file error: %s", err.Error()))
	}

	argsIf, argsOf, argsBs, argsCount, flagArgs := fmt.Sprintf("if=%s", file), "of=/dev/null", fmt.Sprintf("bs=%s", argsBs), fmt.Sprintf("count=%s", argsCount), fmt.Sprintf("iflag=%s", flagArgs)
	for {
		ddBurn(argsIf, argsOf, argsBs, argsCount, flagArgs)
	}
}

func burnWriteDisk(file, argsBs, argsCount, flagArgs string) {
	argsIf, argsOf, argsBs, argsCount, flagArgs := "if=/dev/zero", fmt.Sprintf("of=%s", file), fmt.Sprintf("bs=%s", argsBs), fmt.Sprintf("count=%s", argsCount), fmt.Sprintf("oflag=%s", flagArgs)
	for {
		ddBurn(argsIf, argsOf, argsBs, argsCount, flagArgs)
		if err := os.Remove(file); err != nil {
			common.ExitWithErr(fmt.Sprintf("rm file error: %s", err.Error()))
		}
	}
}

func ddBurn(argsIf, argsOf, argsBs, argsCount, flagArgs string) {
	cmd := exec.Command("dd", argsIf, argsOf, argsBs, argsCount, flagArgs)
	if err := cmd.Start(); err != nil {
		common.ExitWithErr(fmt.Sprintf("dd cmd start error: %s", err.Error()))
	}

	if err := cmd.Wait(); err != nil {
		common.ExitWithErr(fmt.Sprintf("dd cmd wait error: %s", err.Error()))
	}

	if !notFirst {
		//fmt.Println(cmd.Args)
		fmt.Println("[success]inject success")
		notFirst = true
	}
}
