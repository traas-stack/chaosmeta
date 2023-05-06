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

package kernel

import (
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/test/common"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	nProcSleepTime   = 15 * time.Second
	nProcUser        = "chaosmeta_mkpcswa"
	nProcCountOffset = 300
)

func GetNProcTest() []common.TestCase {
	if err := delUser(nProcUser); err != nil {
		fmt.Printf("del user[%s] error: %s\n", nProcUser, err.Error())
	}

	if err := addUser(nProcUser); err != nil {
		panic(any(fmt.Sprintf("add user[%s] error: %s", nProcUser, err.Error())))
	}

	maxPro, err := getMaxNproc(nProcUser)
	if err != nil {
		panic(any(fmt.Sprintf("get max proc of user[%s] error: %s", nProcUser, err.Error())))
	}

	fmt.Println(maxPro)

	var tempCaseList = []common.TestCase{
		{
			Args:  "",
			Error: true,
		},
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  "-u root",
			Error: true,
		},
		{
			Args:  "-u notexist",
			Error: true,
		},
		{
			Args:  "-c -1",
			Error: true,
		},
		{
			Args:  "-c -1",
			Error: true,
		},
		{
			Args: fmt.Sprintf("-u %s -c 1000", nProcUser),
			Check: func() error {
				targetCount := 1000
				if maxPro < 1000 {
					targetCount = maxPro
				}

				return checkNproc(nProcUser, targetCount)
			},
			CheckRecover: func() error {
				return checkNproc(nProcUser, 0)
			},
		},
		{
			Args: fmt.Sprintf("-u %s", nProcUser),
			Check: func() error {
				return checkNproc(nProcUser, maxPro)
			},
			CheckRecover: func() error {
				return checkNproc(nProcUser, 0)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "kernel"
		tempCaseList[i].Fault = "nproc"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkNproc(user string, count int) error {
	time.Sleep(nProcSleepTime)
	nproc, err := getNowProc(user)
	if err != nil {
		return fmt.Errorf("get nproc of user[%s] error: %s", user, err.Error())
	}

	fmt.Printf("user: %s, target nproc: %d, now nproc: %d\n", user, count, nproc)
	if nproc < count-nProcCountOffset || nproc > count+nProcCountOffset {
		return fmt.Errorf("expected nproc: %d, now nproc: %d", count, nproc)
	}

	return nil
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

func getMaxNproc(user string) (int, error) {
	re, err := exec.Command("runuser", "-l", user, "-c", "ulimit -u").CombinedOutput()
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

func addUser(user string) error {
	return exec.Command("/bin/bash", "-c", fmt.Sprintf("useradd %s", user)).Run()
}

func delUser(user string) error {
	return exec.Command("/bin/bash", "-c", fmt.Sprintf("userdel %s", user)).Run()
}
