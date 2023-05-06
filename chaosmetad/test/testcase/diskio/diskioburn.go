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

package diskio

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/process"
	"github.com/traas-stack/chaosmeta/chaosmetad/test/common"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	diskIOBurnSleepSec  int = 3
	diskIOBurnSleepTime     = time.Duration(diskIOBurnSleepSec) * time.Second

	diskioReadBytes, diskioWriteBytes, diskioReadCount, diskioWriteCount uint64

	diskioBurnMBOffset    int = 5000
	diskioBurnCountOffset int = 5000
)

func updateDiskioStat(dir string) error {
	dev, err := getDevByDir(dir)
	if err != nil {
		return fmt.Errorf("get dev from path[%s] error: %s", dir, err.Error())
	}

	stat, err := disk.IOCounters(dev)
	if err != nil {
		return fmt.Errorf("get disk[%s] io stat error: %s", dev, err.Error())
	}

	dev = filepath.Base(dev)
	diskioReadBytes, diskioWriteBytes, diskioReadCount, diskioWriteCount =
		stat[dev].ReadBytes, stat[dev].WriteBytes, stat[dev].ReadCount, stat[dev].WriteCount

	return nil
}

func GetDiskIOBurnTest() []common.TestCase {
	var tempCaseList = []common.TestCase{
		{
			Args:  "egbs",
			Error: true,
		},
		{
			Args:  "-m auto",
			Error: true,
		},
		{
			Args:  "-m all",
			Error: true,
		},
		{
			Args:  "-b 100B",
			Error: true,
		},
		{
			Args:  "-b 1GB",
			Error: true,
		},
		{
			Args:  "-b 1025MB",
			Error: true,
		},
		{
			Args:  "-d /notexist",
			Error: true,
		},
		{
			Args:  "-b 1GB",
			Error: true,
		},
		{
			Args: "",
			PreProcessor: func() error {
				return updateDiskioStat("/tmp")
			},
			Check: func() error {
				if err := checkExistPro(true, "/tmp"); err != nil {
					return err
				}

				return checkDiskioBurn(10000, 1000, 10000, 1000, "/tmp")
			},

			CheckRecover: func() error {
				time.Sleep(diskIOBurnSleepTime)
				return checkExistPro(false, "/tmp")
			},
		},
		{
			Args: "-m write",
			PreProcessor: func() error {
				return updateDiskioStat("/tmp")
			},
			Check: func() error {
				if err := checkExistPro(true, "/tmp"); err != nil {
					return err
				}

				return checkDiskioBurn(0, 5000, 0, 5000, "/tmp")
			},

			CheckRecover: func() error {
				time.Sleep(diskIOBurnSleepTime)
				return checkExistPro(false, "/tmp")
			},
		},
		{
			Args: "-m write -b 5m -d /data/testdiskio_chaosmeta",
			PreProcessor: func() error {
				dir := "/data/testdiskio_chaosmeta"
				if err := filesys.MkdirP(context.Background(), dir); err != nil {
					return fmt.Errorf("create dir[%s] error: %s", err.Error())
				}
				return updateDiskioStat(dir)
			},
			Check: func() error {
				if err := checkExistPro(true, "/data/testdiskio_chaosmeta"); err != nil {
					return err
				}

				return checkDiskioBurn(0, 5000, 0, 5000, "/data/testdiskio_chaosmeta")
			},
			PostProcessor: func() error {
				return os.Remove("/data/testdiskio_chaosmeta")
			},
			CheckRecover: func() error {
				time.Sleep(diskIOBurnSleepTime)
				return checkExistPro(false, "/data/testdiskio_chaosmeta")
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "diskio"
		tempCaseList[i].Fault = "burn"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return fmt.Errorf("unexpected: exec recover fun")
			}
		}
	}

	return tempCaseList
}

func getDevByDir(dir string) (string, error) {
	re, err := cmdexec.RunBashCmdWithOutput(context.Background(), fmt.Sprintf("df -h %s | sed '1d' | awk '{print $1}'", dir))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(re), nil
}

func checkExistPro(ifExist bool, dir string) error {
	k := "chaosmeta_diskburn"
	exist, err := process.ExistProcessByKey(context.Background(), k)
	if err != nil {
		return fmt.Errorf("check pro[%s] exist error: %s", k, err.Error())
	}
	if exist != ifExist {
		return fmt.Errorf("pro[%s]'s expected exist status is %v, actually: %v", k, ifExist, exist)
	}

	//k = "dd"
	//exist, err = process.ExistProcessByKey(context.Background(), k)
	//if err != nil {
	//	return fmt.Errorf("check pro[%s] exist error: %s", k, err.Error())
	//}
	//if exist != ifExist {
	//	return fmt.Errorf("pro[%s]'s expected exist status is %v, actually: %v", k, exist, ifExist)
	//}

	file := fmt.Sprintf("%s/chaosmeta_diskburn_%s", dir, common.UID)
	exist, err = filesys.ExistPath(file)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", file, err.Error())
	}

	if exist != ifExist {
		return fmt.Errorf("file[%s]'s expected exist status is %v, actually: %v", file, ifExist, exist)
	}

	return nil
}

func checkDiskioBurn(rMB, wMB, rC, wC int, dir string) error {
	time.Sleep(diskIOBurnSleepTime)
	tmpRBytes, tmpWBytes, tmpRCount, tmpWCount := diskioReadBytes, diskioWriteBytes, diskioReadCount, diskioWriteCount

	if err := updateDiskioStat(dir); err != nil {
		return err
	}

	rBMInterval := int((diskioReadBytes - tmpRBytes) / 1024 / 1024)
	rCInterval := int(diskioReadCount - tmpRCount)
	wBMInterval := int((diskioWriteBytes - tmpWBytes) / 1024 / 1024)
	wCInterval := int(diskioWriteCount - tmpWCount)

	fmt.Printf("write MB: %d, read MB: %d, write count: %d, read count: %d\n", wBMInterval, rBMInterval, wCInterval, rCInterval)

	if wBMInterval < wMB {
		return fmt.Errorf("expected write MB: %d, actually: %d", wMB, wBMInterval)
	}

	if rBMInterval < rMB {
		return fmt.Errorf("expected read MB: %d, actually: %d", rMB, rBMInterval)
	}

	if wCInterval < wC {
		return fmt.Errorf("expected write count: %d, actually: %d", wC, wCInterval)
	}

	if rCInterval < rC {
		return fmt.Errorf("expected read count: %d, actually: %d", rC, rCInterval)
	}

	return nil
}
