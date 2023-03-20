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
	"context"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmetad/pkg/utils/namespace"
	"strings"
)

func GetDevList(ctx context.Context, cr, cId string, devStr string) ([]string, error) {
	if devStr == "" {
		return nil, fmt.Errorf("args dev-list is empty")
	}

	devStrList := strings.Split(devStr, ",")
	for _, unit := range devStrList {
		isExist, err := existDev(ctx, cr, cId, unit)
		if err != nil {
			return nil, fmt.Errorf("check dev[%s] exist error: %s", unit, err.Error())
		}

		if !isExist {
			return nil, fmt.Errorf("dev[%s] is not exist", unit)
		}
	}

	return devStrList, nil
}

func existDev(ctx context.Context, cr, cId string, devNum string) (bool, error) {
	cmd := fmt.Sprintf("lsblk -a | grep disk | awk '{print $2}' | grep \"%s\" | wc -l", devNum)
	var (
		re  string
		err error
	)

	if cr == "" {
		re, err = cmdexec.RunBashCmdWithOutput(ctx, cmd)
	} else {
		re, err = cmdexec.ExecContainer(ctx, cr, cId, []string{namespace.MNT}, cmd, cmdexec.ExecRun)
	}

	if err != nil {
		return false, err
	}

	if strings.TrimSpace(re) == "1" {
		return true, nil
	}

	return false, nil
}

func RunFillDisk(ctx context.Context, size int64, file string) error {
	unit := "K"
	if size/1024 >= 100 {
		unit = "M"
		size /= 1024
	}

	if cmdexec.SupportCmd("fallocate") {
		return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("fallocate -l %d%s %s", size, unit, file))
	}

	if cmdexec.SupportCmd("dd") {
		return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("dd if=/dev/zero of=%s bs=1%s count=%d iflag=fullblock", file, unit, size))
	}

	return fmt.Errorf("not support \"fallocate\" and \"dd\"")
}

func GetFillKBytes(dir string, percent int, bytes string) (int64, error) {
	var fillKBytes int64
	usage, err := disk.Usage(dir)
	if err != nil {
		return -1, fmt.Errorf("get disk info error: %s", err.Error())
	}

	if percent != 0 {
		if float64(percent) < usage.UsedPercent {
			return -1, fmt.Errorf("target path current disk usage is %.2f%%, no need to fill", usage.UsedPercent)
		}

		fillKBytes = int64(uint64((float64(percent) - usage.UsedPercent) / 100 * (float64(usage.Total) / 1024)))
	} else {
		fillKBytes, err = utils.GetKBytes(bytes)
		if err != nil {
			return -1, fmt.Errorf("\"bytes\" is invalid: %s", err.Error())
		}
	}

	freeKb := int64(usage.Free / 1024)
	if fillKBytes > freeKb {
		return -1, fmt.Errorf("space not enough, fill: %dKB, free: %dKB", fillKBytes, freeKb)
	}

	// fix bug: If it is the disk where the database file is located, the database cannot be read when it is full.
	// so temporarily free up 10k space to solve the problem
	if fillKBytes == freeKb {
		fillKBytes -= 10
	}

	// prevent overflow
	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}
