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
	cpu2 "github.com/ChaosMetaverse/chaosmetad/pkg/injector/cpu"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"strconv"
	"time"
)

var cpuBurnSleepTime = 5 * time.Second
var cpuUsageOffset = 5

func GetCpuBurnTest() []common.TestCase {
	var tempCaseList = []common.TestCase{
		{
			Args:  "",
			Error: true,
		},
		{
			Args:  "greogh3wg",
			Error: true,
		},
		{
			Args:  "-p -1",
			Error: true,
		},
		{
			Args:  "-p 0",
			Error: true,
		},
		{
			Args:  "-p 101",
			Error: true,
		},
		{
			Args:  "-p reg",
			Error: true,
		},
		{
			Args:  "-p 1hs",
			Error: true,
		},
		{
			Args:  "-p 20%",
			Error: true,
		},
		{
			Args:  "-l 1,2,3",
			Error: true,
		},
		{
			Args:  "-p 20%",
			Error: true,
		},
		{
			Args: "-p 50",
			Check: func() error {
				return ifCpuUsageCount(50, runtime.NumCPU())
			},
		},
		{
			Args:  "-p 20 -l t43",
			Error: true,
		},
		{
			Args:  "-p 20 -l 1-2-3",
			Error: true,
		},
		{
			Args:  "-p 20 -l s-3",
			Error: true,
		},
		{
			Args:  "-p 20 -l g",
			Error: true,
		},
		{
			Args:  "-p 20 -c greg",
			Error: true,
		},
		{
			Args:  "-p 20 -c -1",
			Error: true,
		},
		{
			Args: "-p 50 -c 999999",
			Check: func() error {
				return ifCpuUsageCount(50, runtime.NumCPU())
			},
		},
		{
			Args: "-p 50 -c 2",
			Check: func() error {
				if runtime.NumCPU() < 2 {
					return nil
				}

				return ifCpuUsageCount(50, 2)
			},
		},
		{
			Args: "-p 50 -c 1 -l 0,2",
			Check: func() error {
				if runtime.NumCPU() < 3 {
					return nil
				}

				return ifCpuUsageFit(map[int]bool{
					2: true,
					0: true,
				},
					50,
				)
			},
		},
		{
			Args: "-p 50 -c 3 -l 0,2",
			Check: func() error {
				if runtime.NumCPU() < 3 {
					return nil
				}

				return ifCpuUsageFit(map[int]bool{
					2: true,
					0: true,
				},
					50,
				)
			},
		},
		{
			Args: "-p 50 -l 0,1-2",
			Check: func() error {
				if runtime.NumCPU() < 3 {
					return nil
				}

				return ifCpuUsageFit(map[int]bool{
					1: true,
					2: true,
					0: true,
				},
					50,
				)
			},
		},
		{
			Args: "-p 50 -l 2,1-3",
			Check: func() error {
				if runtime.NumCPU() < 4 {
					return nil
				}

				return ifCpuUsageFit(map[int]bool{
					1: true,
					2: true,
					3: true,
				},
					50,
				)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "cpu"
		tempCaseList[i].Fault = "burn"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		tempCaseList[i].CheckRecover = ifCpuBurnRecover
	}

	return tempCaseList
}

func ifCpuUsageCount(percent int, count int) error {
	time.Sleep(cpuBurnSleepTime)
	p, err := cpu.Percent(time.Second, true)
	if err != nil {
		return fmt.Errorf("get cpu usage error: %s", err.Error())
	}

	var c int
	for i, unitP := range p {
		fmt.Printf("%d: %f\n", i, unitP)
		nowPer := int(unitP)
		if nowPer >= percent-cpuUsageOffset && nowPer <= percent+cpuUsageOffset {
			c++
		}
	}

	if c != count {
		return fmt.Errorf("expected count[%d], now count[%d]", count, c)
	}

	return nil
}

func ifCpuUsageFit(coreMap map[int]bool, percent int) error {
	time.Sleep(cpuBurnSleepTime)
	p, err := cpu.Percent(time.Second, true)
	if err != nil {
		return fmt.Errorf("get cpu usage error: %s", err.Error())
	}

	for i, unitP := range p {
		fmt.Printf("%d: %f\n", i, unitP)
		nowPer := int(unitP)
		if coreMap[i] {
			if nowPer < percent-cpuUsageOffset || nowPer > percent+cpuUsageOffset {
				return fmt.Errorf("target core[%d] is not fit, target[> %d], now[%f]", i, percent, unitP)
			}
		}
	}

	return nil
}

func ifCpuBurnRecover() error {
	isExist, err := utils.ExistProcessByKey(cpu2.CpuBurnKey)
	if err != nil {
		return fmt.Errorf("check process exist error: %s", err.Error())
	}

	if isExist {
		return fmt.Errorf("process is running, please kill")
	}

	return nil
}
