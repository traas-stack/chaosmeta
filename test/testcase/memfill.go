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
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"time"
)

var (
	memFillSleepTime                             = 1 * time.Second
	memPerOffset                                 = 3
	memBytesKbOffset                             = 16000
	memFillDir                                   = "/tmp/chaosmeta_mem_tmpfs"
	memBytesKb, memPer, upperMemPer, lowerMemPer int
)

func getMemFillDir() string {
	return fmt.Sprintf("%s%s", memFillDir, common.UID)
}

func updateMemUsage() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get mem usage error: %s", err.Error())
	}

	memBytesKb = int((v.Total - v.Available) / 1024)
	memPer = int(float64(v.Total-v.Available) / float64(v.Total) * 100)
	fmt.Printf("updated percent: %d%%, bytes: %dKB\n", memPer, memBytesKb)

	upperMemPer = memPer + 20
	if upperMemPer > 100 {
		upperMemPer = 100
	}

	lowerMemPer = memPer - 20
	if lowerMemPer < 0 {
		lowerMemPer = 0
	}

	return nil
}

func GetMemFillTest() []common.TestCase {
	_ = updateMemUsage()
	var tempCaseList = []common.TestCase{
		{
			Args:  "",
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
			Args:  "-b 100m",
			Error: true,
		},
		{
			Args:  "-b 100b",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d", lowerMemPer),
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d", upperMemPer),
			Error: false,
			Check: func() error {
				return checkMemPer(upperMemPer, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemPer(memPer, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  fmt.Sprintf("-p %d -b 5000kb", 100),
			Error: false,
			Check: func() error {
				return checkMemPer(100, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemPer(memPer, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 200000kB",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+200000, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 200000Kb",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+200000, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 200000",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+200000, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 200MB",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+200*1024, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 1gb",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+1024*1024, getMemFillDir(), false, "cache")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "cache")
			},
		},
		{
			Args:  "-b 200mb -m ram",
			Error: false,
			Check: func() error {
				return checkMemByteKb(memBytesKb+200*1024, getMemFillDir(), false, "ram")
			},
			CheckRecover: func() error {
				return checkMemByteKb(memBytesKb, getMemFillDir(), true, "ram")
			},
		},

		{
			Args:  fmt.Sprintf("-p %d -m ram", upperMemPer),
			Error: false,
			Check: func() error {
				return checkMemPer(upperMemPer, getMemFillDir(), false, "ram")
			},
			CheckRecover: func() error {
				return checkMemPer(memPer, getMemFillDir(), true, "ram")
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "mem"
		tempCaseList[i].Fault = "fill"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return fmt.Errorf("unexpected: exec recover fun")
			}
		}
	}

	return tempCaseList
}

func checkMemPer(targetPer int, fillDir string, recover bool, mode string) error {
	time.Sleep(memFillSleepTime)
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get virtual memory error: %s", err.Error())
	}

	now := int(float64(v.Total-v.Available) / float64(v.Total) * 100)

	fmt.Printf("now: %d, target: %d\n", now, targetPer)

	if now < targetPer-memPerOffset || now > targetPer+memPerOffset {
		return fmt.Errorf("unexpected mem percent, now: %d, target: %d", now, targetPer)
	}

	if mode == "ram" {
		return nil
	}

	exist, err := utils.ExistPath(fillDir)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", err.Error())
	}

	if !recover {
		if !exist {
			return fmt.Errorf("file[%s] not exist", fillDir)
		}
	} else {
		if exist {
			return fmt.Errorf("file[%s] exist", fillDir)
		}

		time.Sleep(memFillSleepTime)
		if err := updateMemUsage(); err != nil {
			return fmt.Errorf("update mem usage error: %s", err.Error())
		}
	}

	return nil
}

func checkMemByteKb(targetByte int, fillDir string, recover bool, mode string) error {
	time.Sleep(memFillSleepTime)
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get virtual memory error: %s", err.Error())
	}

	now := int(v.Total-v.Available) / 1024

	fmt.Printf("now: %dkb, target: %dkb\n", now, targetByte)

	if now < targetByte-memBytesKbOffset || now > targetByte+memBytesKbOffset {
		return fmt.Errorf("unexpected mem used byte, now: %dkb, target: %dkb", now, targetByte)
	}

	if mode == "ram" {
		return nil
	}

	exist, err := utils.ExistPath(fillDir)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", err.Error())
	}

	if !recover {
		if !exist {
			return fmt.Errorf("file[%s] not exist", fillDir)
		}
	} else {
		if exist {
			return fmt.Errorf("file[%s] exist", fillDir)
		}

		time.Sleep(memFillSleepTime)
		if err := updateMemUsage(); err != nil {
			return fmt.Errorf("update mem usage error: %s", err.Error())
		}
	}

	return nil
}
