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
	process2 "github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"github.com/shirou/gopsutil/process"
	"strconv"
	"time"
)

var (
	proStopSleepTime = 200 * time.Millisecond
	proStopCmd       = "sleep 23"
	proStopPid       int
)

func GetProStopTest() []common.TestCase {
	ctx := context.Background()
	process2.KillProcessByKey(ctx, proStopCmd, process2.SIGKILL)

	var err error
	proStopPid, err = startDaemonCmd(ctx, proStopCmd)
	if err != nil {
		panic(any(fmt.Sprintf("start test process error: %s", err.Error())))
	}

	var tempCaseList = []common.TestCase{
		{
			Args:  "",
			Error: true,
		},
		{
			Args:  "-p 29999",
			Error: true,
		},
		{
			Args:  "-k fbigrw3",
			Error: true,
		},
		{
			Args:  "-p -1",
			Error: true,
		},
		{
			Args:  "-p 2h",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d", proStopPid),
			Error: false,
			Check: func() error {

				return checkProStatusByPid(proStopPid, "T")
			},
			CheckRecover: func() error {
				startDaemonCmd(ctx, proStopCmd)
				startDaemonCmd(ctx, proStopCmd)
				return checkProStatusByPid(proStopPid, "S")
			},
		},
		{
			Args:  fmt.Sprintf("-k '%s'", proStopCmd),
			Error: false,
			Check: func() error {

				return checkProStatusByKey(ctx, proStopCmd, "T", 3)
			},
			CheckRecover: func() error {
				return checkProStatusByKey(ctx, proStopCmd, "S", 3)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "process"
		tempCaseList[i].Fault = "stop"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkProStatusByKey(ctx context.Context, key string, expectedStatus string, expectedCount int) error {
	time.Sleep(proStopSleepTime)
	fmt.Printf("key: %s, expected status: %s, expected count: %d\n", key, expectedStatus, expectedCount)

	pidList, err := process2.GetPidListByKey(ctx, "", "", key)
	if err != nil {
		return fmt.Errorf("get pid list by key[%s] error: %s", key, err.Error())
	}

	if len(pidList) != expectedCount {
		return fmt.Errorf("now count: %d, expected: %d", len(pidList), expectedCount)
	}

	for _, p := range pidList {
		if err := checkProStatusByPid(p, expectedStatus); err != nil {
			return err
		}
	}

	return nil
}

func checkProStatusByPid(pid int, expectedStatus string) error {
	time.Sleep(proStopSleepTime)
	fmt.Printf("pid[%d]'s expected status %s\n", pid, expectedStatus)
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("get process by pid[%d] error: %s", pid, err.Error())
	}

	status, err := p.Status()
	if err != nil {
		return fmt.Errorf("get process's status error: %s", err.Error())
	}

	if status != expectedStatus {
		return fmt.Errorf("now status: %s, expected status: %s", status, expectedStatus)
	}

	return nil
}
