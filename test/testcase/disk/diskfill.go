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

package disk

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/traas-stack/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmetad/test/common"
	"math"
	"os/exec"
	"strconv"
	"time"
)

var (
	diskFillSleepTime = 1 * time.Second
	diskUsageOffset   = 3
	diskBytesKbOffset = 10000
	tmpDir, dataDir   = "/tmp", "/data"
	diskFillFile      = "chaosmeta_fill"

	tmpUsedKb, dataUsedKb uint64
	tmpUsagePer, lowerTmpPer, upperTmpPer,
	dataUsagePer, lowerDataPer, upperDataPer int
)

func updateDiskUsage() (err error) {
	tmpUsedKb, tmpUsagePer, lowerTmpPer, upperTmpPer, err = getDirUsage(tmpDir)
	if err != nil {
		return
	}

	dataUsedKb, dataUsagePer, lowerDataPer, upperDataPer, err = getDirUsage(dataDir)

	return
}

func GetDiskFillTest() []common.TestCase {
	_ = exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s; mkdir -p %s", tmpDir, dataDir)).Run()
	_ = updateDiskUsage()
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
			Args:  "-p 101",
			Error: true,
		},
		{
			Args:  "-p 0",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d", lowerTmpPer),
			Error: true,
		},
		{
			Args:  "-d /abcnotexist",
			Error: true,
		},
		{
			Args:  "-d /data -b -1",
			Error: true,
		},
		{
			Args:  "-d /data -b 1000B",
			Error: true,
		},
		{
			Args:  "-d /data -b s1000B",
			Error: true,
		},
		{
			Args:  "-d /data -b MB",
			Error: true,
		},
		{
			Args:  "-d /data -b 0.5gb",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-d /data -b %dtb", math.MaxInt),
			Error: true,
		},
		{
			Args:  "-d /data -b 1tb",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d", upperTmpPer),
			Error: false,
			Check: func() error {
				return checkDiskFillPer(tmpDir, getTmpFile(), upperTmpPer, false)
			},
			CheckRecover: func() error {
				return checkDiskFillPer(tmpDir, getTmpFile(), tmpUsagePer, true)
			},
		},
		{
			Args:  "-p 100",
			Error: false,
			Check: func() error {
				return checkDiskFillPer(tmpDir, getTmpFile(), 100, false)
			},
			CheckRecover: func() error {
				return checkDiskFillPer(tmpDir, getTmpFile(), tmpUsagePer, true)
			},
		},
		{
			Args:  fmt.Sprintf("-d /data -p %d", upperDataPer),
			Error: false,
			Check: func() error {
				return checkDiskFillPer(dataDir, getDataFile(), upperDataPer, false)
			},
			CheckRecover: func() error {
				return checkDiskFillPer(dataDir, getDataFile(), dataUsagePer, true)
			},
		},
		{
			Args:  fmt.Sprintf("-d /data -p %d -b 10Mb", upperDataPer),
			Error: false,
			Check: func() error {
				return checkDiskFillPer(dataDir, getDataFile(), upperDataPer, false)
			},
			CheckRecover: func() error {
				return checkDiskFillPer(dataDir, getDataFile(), dataUsagePer, true)
			},
		},
		{
			Args:  "-d /data -b 500000",
			Error: false,
			Check: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb)+500000, false)
			},
			CheckRecover: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb), true)
			},
		},
		{
			Args:  "-d /data -b 500mB",
			Error: false,
			Check: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb)+500*1024, false)
			},
			CheckRecover: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb), true)
			},
		},
		{
			Args:  "-d /data -b 500Mb",
			Error: false,
			Check: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb)+500*1024, false)
			},
			CheckRecover: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb), true)
			},
		},
		{
			Args:  "-d /data -b 500000kB",
			Error: false,
			Check: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb)+500000, false)
			},
			CheckRecover: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb), true)
			},
		},
		{
			Args:  "-d /data -b 2gb",
			Error: false,
			Check: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb)+2*1024*1024, false)
			},
			CheckRecover: func() error {
				return checkDiskFillByte(dataDir, getDataFile(), int(dataUsedKb), true)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "disk"
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

func checkDiskFillPer(dir, file string, targetPer int, recover bool) error {
	time.Sleep(diskFillSleepTime)
	u, err := disk.Usage(dir)
	if err != nil {
		return fmt.Errorf("get disk usage of %s error: %s", dir, err.Error())
	}
	now := int(u.UsedPercent)
	fmt.Printf("now per: %d, target per: %d\n", now, targetPer)
	if now < targetPer-diskUsageOffset || now > targetPer+diskUsageOffset {
		return fmt.Errorf("usage is unexpected: %d", now)
	}

	exist, err := filesys.ExistPath(file)
	if err != nil {
		return fmt.Errorf("check file exsit error: %s", err.Error())
	}

	if !recover {
		if !exist {
			return fmt.Errorf("fill file not exist:%s", file)
		}
	} else {
		if exist {
			return fmt.Errorf("fill file exist:%s", file)
		}

		time.Sleep(diskFillSleepTime)
		if err := updateDiskUsage(); err != nil {
			return fmt.Errorf("update disk usage error: %s", err.Error())
		}
	}

	return nil
}

func checkDiskFillByte(dir, file string, targetByte int, recover bool) error {
	time.Sleep(diskFillSleepTime)
	u, err := disk.Usage(dir)
	if err != nil {
		return fmt.Errorf("get disk usage of %s error: %s", dir, err.Error())
	}
	now := int(u.Used) / 1024
	fmt.Printf("now bytes: %dKB, target bytes: %dKB\n", now, targetByte)
	if now < targetByte-diskBytesKbOffset || now > targetByte+diskBytesKbOffset {
		return fmt.Errorf("usage is unexpected: %d", now)
	}

	exist, err := filesys.ExistPath(file)
	if err != nil {
		return fmt.Errorf("check file exsit error: %s", err.Error())
	}

	if !recover {
		if !exist {
			return fmt.Errorf("fill file not exist:%s", file)
		}
	} else {
		if exist {
			return fmt.Errorf("fill file exist:%s", file)
		}

		time.Sleep(diskFillSleepTime)
		if err := updateDiskUsage(); err != nil {
			return fmt.Errorf("update disk usage error: %s", err.Error())
		}
	}

	return nil
}

func getTmpFile() string {
	return fmt.Sprintf("%s/%s%s.dat", tmpDir, diskFillFile, common.UID)
}

func getDataFile() string {
	return fmt.Sprintf("%s/%s%s.dat", dataDir, diskFillFile, common.UID)
}

func getDirUsage(dir string) (uint64, int, int, int, error) {
	tmpUsage, err := disk.Usage(dir)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("get usage of disk[%s] error: %s", dir, err.Error())
	}

	tmpUsedKb := tmpUsage.Used / 1024
	tmpUsagePer := int(tmpUsage.UsedPercent)
	lowerTmpPer := tmpUsagePer - 10
	if lowerTmpPer < 0 {
		lowerTmpPer = 0
	}

	upperTmpPer := tmpUsagePer + 10
	if upperTmpPer > 100 {
		upperTmpPer = 100
	}
	fmt.Printf("updated dir: %s, now per: %d, bytes: %dKB\n", dir, tmpUsagePer, tmpUsedKb)

	return tmpUsedKb, tmpUsagePer, lowerTmpPer, upperTmpPer, nil
}
