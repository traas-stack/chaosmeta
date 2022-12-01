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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"strconv"
	"time"
)

var (
	fdSleepTime = 2 * time.Second
	fdOffset    = 1200
)

func GetFdTest() []common.TestCase {
	nowFd, maxFd, err := utils.GetKernelFdStatus()
	if err != nil {
		panic(any(fmt.Sprintf("get kernel max fd count error: %s", err.Error())))
	}

	var tempCaseList = []common.TestCase{
		{
			Args:  "awvgv",
			Error: true,
		},
		{
			Args:  "-m test",
			Error: true,
		},
		{
			Args:  "",
			Error: false,
			Check: func() error {
				return checkFd(nowFd, nowFd-2000)
			},
			CheckRecover: func() error {
				return checkFd(nowFd, maxFd)
			},
		},
		{
			Args:  "-m fill -c 100000",
			Error: false,
			Check: func() error {
				addCount := 100000
				if maxFd < 100000 {
					addCount = maxFd
				}

				return checkFd(nowFd+addCount, maxFd)
			},
			CheckRecover: func() error {
				return checkFd(nowFd, maxFd)
			},
		},
		{
			Args:  "-m fill",
			Error: false,
			Check: func() error {
				return checkFd(maxFd, maxFd)
			},
			CheckRecover: func() error {
				return checkFd(nowFd, maxFd)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "kernel"
		tempCaseList[i].Fault = "fdfull"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkFd(targetNow, targetMax int) error {
	time.Sleep(fdSleepTime)
	nowFd, maxFd, err := utils.GetKernelFdStatus()
	if err != nil {
		return fmt.Errorf("get kernel max fd count error: %s", err.Error())
	}

	fmt.Printf("nowFd: %d, maxFd: %d\n", nowFd, maxFd)
	fmt.Printf("targetNow: %d, targetMax: %d\n", targetNow, targetMax)

	if nowFd > maxFd {
		return nil
	}

	if nowFd < targetNow-fdOffset || nowFd > targetNow+fdOffset {
		return fmt.Errorf("nowFd: %d, target: %d", nowFd, targetNow)
	}

	if maxFd < targetMax-fdOffset || maxFd > targetMax+fdOffset {
		return fmt.Errorf("maxFd: %d, target: %d", maxFd, targetMax)
	}

	return nil
}
