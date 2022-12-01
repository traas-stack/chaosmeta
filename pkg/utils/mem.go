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

package utils

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/shirou/gopsutil/mem"
)

// CalculateFillKBytes The calculation of memory usage is consistent with the calculation method of the top command: Available/Total.
// Because whether oom is calculated according to this
func CalculateFillKBytes(percent int, fillBytes string) (int64, error) {
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
		fillKBytes, _ = GetKBytes(fillBytes)
	}

	// prevent overflow
	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}

func FillCache(fillKBytes int64, dir string, filename string) error {
	if err := MkdirP(dir); err != nil {
		return fmt.Errorf("create tmpfs dir[%s] error: %s", dir, err.Error())
	}

	file := fmt.Sprintf("%s/%s", dir, filename)

	if err := RunBashCmdWithoutOutput(fmt.Sprintf("mount -t tmpfs tmpfs %s -o size=%dk", dir, fillKBytes)); err != nil {
		return fmt.Errorf("mount tmpfs[%s] error: %s", dir, err.Error())
	}

	if err := RunFillDisk(fillKBytes, file); err != nil {
		UndoTmpfs(dir)
		return fmt.Errorf("fill file[%s] error: %s", file, err.Error())
	}

	return nil
}

func UndoTmpfs(dir string) error {
	logger := log.GetLogger()

	if err := RunBashCmdWithoutOutput(fmt.Sprintf("umount %s", dir)); err != nil {
		logger.Warnf("umount %s error: %s", dir, err.Error())
	}

	if err := RunBashCmdWithoutOutput(fmt.Sprintf("rm -rf %s", dir)); err != nil {
		logger.Warnf("rm %s error: %s", dir, err.Error())
		return fmt.Errorf("rm %s error: %s", dir, err.Error())
	}

	return nil
}

func FillRam(memFillKey string, fillKBytes int64, uid string, timeout int64) (int, error) {
	return StartBashCmdAndWaitPid(fmt.Sprintf("%s %s %d %dkb %d", GetToolPath(memFillKey), uid, -999, fillKBytes, timeout))
}
