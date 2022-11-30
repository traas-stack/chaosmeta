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
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ps h -Led -o user | sort | uniq -c | grep temp | awk '{print $1}'
// uid user count timeout
func main() {
	args := os.Args
	if len(args) < 5 {
		common.ExitWithErr(fmt.Sprintf("args must provide: uid, user, count, timeout"))
	}
	user, countStr, timeStr := args[2], args[3], args[4]
	timeout, err := strconv.Atoi(timeStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("args timeout is not a num"))
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("args count is not a num"))
	}

	nproc, err := getNproc()
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("get nproc error: %s", err.Error()))
	}

	pCount, err := getNowProc(user)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("get proc count error: %s", err.Error()))
	}

	needP := nproc - pCount + 100
	if count <= 0 || count >= needP {
		count = needP
	}

	//fmt.Printf("nproc: %d, now: %d, need: %d\n", nproc, pCount, count)

	cmd := fmt.Sprintf("sleep %ds", timeout)
	if timeout == 0 {
		cmd = fmt.Sprintf("while true; do sleep %ds; done", 15)
		count /= 2
	}

	for i := 0; i < count; i++ {
		go simpleProc(cmd)
	}

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}

func simpleProc(cmd string) {
	for {
		if err := exec.Command("/bin/bash", "-c", cmd).Start(); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func getNowProc(user string) (int, error) {
	re, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("ps h -Led -o user | sort | uniq -c | grep %s | awk '{print $1}'", user)).CombinedOutput()
	if err != nil {
		return -1, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	reStr := strings.TrimSpace(string(re))
	if reStr == "" {
		return 0, nil
	}

	nproc, err := strconv.Atoi(reStr)
	if err != nil {
		return -1, fmt.Errorf("%s is not a num", reStr)
	}

	return nproc, nil
}

func getNproc() (int, error) {
	re, err := exec.Command("/bin/bash", "-c", "ulimit -u").CombinedOutput()
	if err != nil {
		return -1, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	reStr := strings.TrimSpace(string(re))
	nproc, err := strconv.Atoi(reStr)
	if err != nil {
		return -1, fmt.Errorf("%s is not a num", reStr)
	}

	return nproc, nil
}
