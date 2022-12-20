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
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/disk"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/shirou/gopsutil/mem"
	"os"
)

// CalculateFillKBytes The calculation of memory usage is consistent with the calculation method of the top command: Available/Total.
// Because whether oom is calculated according to this
func CalculateFillKBytes(ctx context.Context, percent int, fillBytes string) (int64, error) {
	var fillKBytes int64
	if percent != 0 {
		v, err := mem.VirtualMemory()
		if err != nil {
			return -1, fmt.Errorf("check vm error: %s", err.Error())
		}

		usedPercent := float64(v.Total-v.Available) / float64(v.Total) * 100

		if float64(percent) < usedPercent {
			return -1, fmt.Errorf("current mem usage is %.2f%%, no need to fill any mem", usedPercent)
		}

		fillKBytes = int64((float64(percent) - usedPercent) / 100 * (float64(v.Total) / 1024))
	} else {
		fillKBytes, _ = utils.GetKBytes(fillBytes)
	}

	// prevent overflow
	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}

func FillCache(ctx context.Context, fillKBytes int64, dir string, filename string) error {
	if err := filesys.MkdirP(ctx, dir); err != nil {
		return fmt.Errorf("create tmpfs dir[%s] error: %s", dir, err.Error())
	}

	file := fmt.Sprintf("%s/%s", dir, filename)

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("mount -t tmpfs tmpfs %s -o size=%dk", dir, fillKBytes)); err != nil {
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

	if err := os.RemoveAll(dir); err != nil {
		logger.Warnf("rm %s error: %s", dir, err.Error())
		return fmt.Errorf("rm %s error: %s", dir, err.Error())
	}

	return nil
}
