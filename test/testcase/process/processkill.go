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

package process

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"strconv"
	"time"
)

var (
	proKillSleepTime = 200 * time.Millisecond
	proKillCmd       = "sleep 1579"
	proKillPid       int
)

func startDaemonCmd(ctx context.Context, cmd string) (int, error) {
	err := cmdexec.RunBashCmdWithoutOutput(ctx, cmd+"&")
	if err != nil {
		return -1, fmt.Errorf("start simple process error: %s", err.Error())
	}

	pid, err := process.GetPidByKeyWithoutRunUser(ctx, cmd)
	if err != nil {
		return -1, fmt.Errorf("get pid by key error: %s", err.Error())
	}

	return pid, nil
}

func GetProKillTest() []common.TestCase {
	ctx := context.Background()
	var err error
	proKillPid, err = startDaemonCmd(ctx, proKillCmd)
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
			Args:  fmt.Sprintf("-p %d -r '%s'", proKillPid, proKillCmd),
			Error: false,
			Check: func() error {
				time.Sleep(proKillSleepTime)
				fmt.Printf("target pid: %d\n", proKillPid)
				proExist, err := process.ExistPid(ctx, proKillPid)
				if err != nil {
					return fmt.Errorf("check pid[%d] exist error: %s", proKillPid, err.Error())
				}

				if proExist {
					return fmt.Errorf("pid[%d] still exist", proKillPid)
				}

				return nil
			},
			CheckRecover: func() error {
				return checkProExistByKey(ctx, "", "", proKillCmd, true)
			},
		},
		{
			Args:  fmt.Sprintf("-k '%s' -r '%s'", proKillCmd, proKillCmd),
			Error: false,
			Check: func() error {
				return checkProExistByKey(ctx, "", "", proKillCmd, false)
			},
			CheckRecover: func() error {
				return checkProExistByKey(ctx, "", "", proKillCmd, true)
			},
		},
		{
			Args:  fmt.Sprintf("-k '%s' -r '%s'", proKillCmd, "chaosfalsed"),
			Error: false,
			Check: func() error {
				return checkProExistByKey(ctx, "", "", proKillCmd, false)
			},
			CheckRecover: func() error {
				return checkProExistByKey(ctx, "", "", proKillCmd, false)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "process"
		tempCaseList[i].Fault = "kill"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkProExistByKey(ctx context.Context, cr, cId, key string, expected bool) error {
	fmt.Printf("expected exist status: %v\n", expected)
	time.Sleep(proKillSleepTime)
	exist, err := process.ExistProcessByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("check process exist by key[%s] error: %s", key, err.Error())
	}

	if exist != expected {
		return fmt.Errorf("process exist status is unexpected")
	}

	return nil
}
