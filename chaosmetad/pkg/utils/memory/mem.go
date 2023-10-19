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

package memory

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/disk"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
)

// CalculateFillKBytes The calculation of memory usage is consistent with the calculation method of the top command: Available/Total.
// Because whether oom is calculated according to this
func CalculateFillKBytes(ctx context.Context, cr, cId string, percent int, fillBytes string) (int64, error) {
	var fillKBytes int64
	if percent != 0 {
		total, err := getMemTotal(ctx, cr, cId)
		if err != nil {
			return -1, fmt.Errorf("get total mem error: %s", err.Error())
		}

		avail, err := getMemAvailable(ctx, cr, cId)
		if err != nil {
			return -1, fmt.Errorf("get avail mem error: %s", err.Error())
		}

		usedPercent := (total - avail) / total * 100
		if float64(percent) < usedPercent {
			return -1, fmt.Errorf("current mem usage is %.2f%%, no need to fill any mem", usedPercent)
		}

		fillKBytes = int64((float64(percent) - usedPercent) / 100 * total)
	} else {
		fillKBytes, _ = utils.GetKBytes(fillBytes)
	}

	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}

func FillCache(ctx context.Context, cr, cId string, percent int, bytes string, dir string, filename string) error {
	fillKBytes, err := CalculateFillKBytes(ctx, cr, cId, percent, bytes)
	if err != nil {
		return err
	}

	if err := filesys.MkdirP(ctx, dir); err != nil {
		return fmt.Errorf("create tmpfs dir[%s] error: %s", dir, err.Error())
	}

	file := fmt.Sprintf("%s/%s", dir, filename)

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("mount -t tmpfs tmpfs %s -o size=%dk", dir, fillKBytes)); err != nil {
		UndoTmpfs(ctx, dir)
		return fmt.Errorf("mount tmpfs[%s] error: %s", dir, err.Error())
	}

	if err := disk.RunFillDisk(ctx, fillKBytes, file); err != nil {
		UndoTmpfs(ctx, dir)
		return fmt.Errorf("fill file[%s] error: %s", file, err.Error())
	}

	return nil
}

func UndoTmpfs(ctx context.Context, dir string) error {
	logger := log.GetLogger(ctx)

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("umount %s", dir)); err != nil {
		logger.Warnf("umount %s error: %s", dir, err.Error())
	}

	time.Sleep(500 * time.Millisecond)

	if err := os.RemoveAll(dir); err != nil {
		logger.Warnf("rm %s error: %s", dir, err.Error())
		return fmt.Errorf("rm %s error: %s", dir, err.Error())
	}

	return nil
}

func GetContainerMemTotal(ctx context.Context, cr, cId string) (memTotal float64, err error) {
	path, err := cgroup.GetContainerCgroupPath(ctx, cr, cId, cgroup.MEMORY)
	if err != nil {
		return 0, err
	}
	memTotalStr, err := cgroup.ReadCgroupFileStr(ctx, path, cgroup.MEMORY, cgroup.MemoryLimitInBytesFile)
	if err != nil {
		return 0, fmt.Errorf("read total mem error: %s", err.Error())
	}

	return getMemByStr(memTotalStr)
}

//func GetContainerMemCache(ctx context.Context, cr, cId string) (memTotal float64, err error) {
//	path, err := cgroup.GetContainerCgroupPath(ctx, cr, cId, cgroup.MEMORY)
//	if err != nil {
//		return 0, err
//	}
//	memStatStr, err := cgroup.ReadCgroupFileStr(ctx, path, cgroup.MEMORY, cgroup.MemoryStatFile)
//	if err != nil {
//		return 0, fmt.Errorf("read total mem error: %s", err.Error())
//	}
//
//	memStatStrList := strings.Split(memStatStr, "\n")
//
//	var cacheStr, rssStr string
//	for _, unitStr := range memStatStrList {
//		kvList := strings.Split(unitStr, " ")
//		if kvList[0] == "cache" {
//			cacheStr = kvList[1]
//		}
//		if kvList[0] == "rss" {
//			rssStr = kvList[1]
//		}
//	}
//
//	if cacheStr == "" || rssStr == "" {
//		return 0, fmt.Errorf("not found cache or rss data")
//	}
//
//	cache, err := getMemByStr(cacheStr)
//	if err != nil {
//		return 0, fmt.Errorf("cache[%s] is not a valid float: %s", cacheStr, err.Error())
//	}
//
//	rss, err := getMemByStr(rssStr)
//	if err != nil {
//		return 0, fmt.Errorf("rss[%s] is not a valid float: %s", cacheStr, err.Error())
//	}
//
//	return cache - rss, nil
//}

func GetContainerMemAvailable(ctx context.Context, cr, cId string) (memAvailable float64, err error) {
	path, err := cgroup.GetContainerCgroupPath(ctx, cr, cId, cgroup.MEMORY)
	if err != nil {
		return 0, err
	}
	memUsageStr, err := cgroup.ReadCgroupFileStr(ctx, path, cgroup.MEMORY, cgroup.MemoryUsageInBytesFile)
	if err != nil {
		return 0, fmt.Errorf("read usage mem error: %s", err.Error())
	}

	memUsage, err := getMemByStr(memUsageStr)
	if err != nil {
		return 0, fmt.Errorf("get value from mem string[%s] error: %s", memUsageStr, err.Error())
	}

	//memCache, err := GetContainerMemCache(ctx, cr, cId)
	//if err != nil {
	//	return 0, fmt.Errorf("get cache mem error: %s", err.Error())
	//}

	memTotal, err := GetContainerMemTotal(ctx, cr, cId)
	if err != nil {
		return 0, fmt.Errorf("get total mem error: %s", err.Error())
	}

	//return memTotal - memUsage + memCache, nil
	return memTotal - memUsage, nil
}

func GetHostMemTotal(ctx context.Context, cr, cId string) (float64, error) {
	cmd := fmt.Sprintf("grep -m1 MemTotal /proc/meminfo | sed 's/[^0-9]*//g'")
	totalStr, err := cmdexec.ExecCommonWithNS(ctx, cr, cId, cmd, []string{namespace.MNT})
	totalStr = strings.TrimSpace(totalStr)
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return -1, fmt.Errorf("get total mem[%s] error: %s", totalStr, err.Error())
	}

	return total, err
}

func GetHostMemAvailable(ctx context.Context, cr, cId string) (float64, error) {
	cmd := fmt.Sprintf("grep -m1 MemAvailable /proc/meminfo | sed 's/[^0-9]*//g'")
	availStr, err := cmdexec.ExecCommonWithNS(ctx, cr, cId, cmd, []string{namespace.MNT})
	availStr = strings.TrimSpace(availStr)
	avail, err := strconv.ParseFloat(availStr, 64)
	if err != nil {
		return -1, fmt.Errorf("get avail mem[%s] error: %s", availStr, err.Error())
	}

	return avail, err
}

func getMemTotal(ctx context.Context, cr, cId string) (float64, error) {
	if cr == "" {
		return GetHostMemTotal(ctx, cr, cId)
	} else {
		return GetContainerMemTotal(ctx, cr, cId)
	}
}

func getMemAvailable(ctx context.Context, cr, cId string) (float64, error) {
	if cr == "" {
		return GetHostMemAvailable(ctx, cr, cId)
	} else {
		return GetContainerMemAvailable(ctx, cr, cId)
	}
}

func getMemByStr(content string) (float64, error) {
	value, err := strconv.ParseFloat(content, 64)
	// Not allowed to use --percent to inject pressure into containers without setting memory limits
	if value == cgroup.MemUnLimit {
		return -1, fmt.Errorf("Container has not set a memory limit and should be filled with a fixed size of pressure using the --bytes.")
	}

	if err != nil {
		return 0, err
	}
	return value / 1024, nil
}
