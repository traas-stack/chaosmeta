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
	"strconv"
	"strings"
)

const (
	BlkioPath       = "/sys/fs/cgroup/blkio"
	WriteBytesFile  = "blkio.throttle.write_bps_device"
	ReadBytesFile   = "blkio.throttle.read_bps_device"
	WriteIOFile     = "blkio.throttle.write_iops_device"
	ReadIOFile      = "blkio.throttle.read_iops_device"
	BlkioCgroupName = "chaosmeta_blkio"
)

func NewCgroup(cgroupPath string, configCmdStr string) error {
	if err := RunBashCmdWithoutOutput(fmt.Sprintf("mkdir %s%s%s", cgroupPath, CmdSplit, configCmdStr)); err != nil {
		return err
	}

	return nil
}

func GetBlkioCPath(uid string) string {
	return fmt.Sprintf("%s/%s_%s", BlkioPath, BlkioCgroupName, uid)
}

func GetBlkioConfig(devList []string, rBytes, wBytes string, rIO, wIO int64, cgroupPath string) string {
	var re = ""
	if rBytes != "" {
		b, _ := GetBytes(rBytes)
		re += getThrottleDeviceCmdStr(devList, b, fmt.Sprintf("%s/%s", cgroupPath, ReadBytesFile))
	}

	if wBytes != "" {
		b, _ := GetBytes(wBytes)
		re += getThrottleDeviceCmdStr(devList, b, fmt.Sprintf("%s/%s", cgroupPath, WriteBytesFile))
	}

	if rIO != 0 {
		re += getThrottleDeviceCmdStr(devList, rIO, fmt.Sprintf("%s/%s", cgroupPath, ReadIOFile))
	}

	if wIO != 0 {
		re += getThrottleDeviceCmdStr(devList, wIO, fmt.Sprintf("%s/%s", cgroupPath, WriteIOFile))
	}

	log.GetLogger().Debugf("blkio config: %s", re)
	return re[:len(re)-len(CmdSplit)]
}

func getThrottleDeviceCmdStr(devList []string, value int64, filename string) string {
	var re string
	for _, unitDec := range devList {
		re += fmt.Sprintf("echo %s %d > %s%s", unitDec, value, filename, CmdSplit)
	}

	return re
}

func CheckPidListCgroup(pidList []int) error {
	for _, unitP := range pidList {
		oldPath, err := getpidCurCgroup(unitP)
		if err != nil {
			return fmt.Errorf("get old cgroup path of process[%d] error: %s", unitP, err.Error())
		}

		if strings.Index(oldPath, BlkioCgroupName) >= 0 {
			return fmt.Errorf("%d is in experiment[%s]", unitP, oldPath)
		}
	}

	return nil
}

func GetPidListCurCgroup(pidList []int) (map[int]string, error) {
	var re = make(map[int]string)
	for _, unitP := range pidList {
		oldPath, err := getpidCurCgroup(unitP)
		if err != nil {
			return nil, fmt.Errorf("get old cgroup path of process[%d] error: %s", unitP, err.Error())
		}
		re[unitP] = oldPath
	}

	return re, nil
}

func getpidCurCgroup(pid int) (string, error) {
	reByte, err := RunBashCmdWithOutput(fmt.Sprintf("cat /proc/%d/cgroup | grep blkio", pid))
	if err != nil {
		return "", fmt.Errorf("run cmd error: %s", err.Error())
	}

	out := strings.TrimSpace(string(reByte))
	sArr := strings.Split(out, ":")
	if sArr[1] != "blkio" {
		return "", fmt.Errorf("out string is not valid: %s", out)
	}

	return sArr[2], nil
}

func MovePidListToCgroup(pidList []int, cgroupPath string) error {
	for _, unit := range pidList {
		if err := MoveToCgroup(unit, cgroupPath); err != nil {
			return fmt.Errorf("move pid[%d] to cgroup[%s] error: %s", unit, cgroupPath, err.Error())
		}
	}

	return nil
}

func MoveToCgroup(pid int, cgroupPath string) error {
	if err := RunBashCmdWithoutOutput(fmt.Sprintf("echo %d > %s/tasks", pid, cgroupPath)); err != nil {
		return err
	}

	return nil
}

func GetPidStrListByCgroup(cgroupPath string) ([]int, error) {
	reByte, err := RunBashCmdWithOutput(fmt.Sprintf("cat %s/tasks", cgroupPath))
	if err != nil {
		return nil, fmt.Errorf("run cmd error: %s", err.Error())
	}

	var pidList []int
	strList := strings.Split(string(reByte), "\n")
	for _, unit := range strList {
		if unit == "" {
			continue
		}

		pid, err := strconv.Atoi(unit)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid pid: %s", unit, err.Error())
		}

		pidList = append(pidList, pid)
	}

	return pidList, nil
}

func RemoveCgroup(cgroupPath string) error {
	if err := RunBashCmdWithoutOutput(fmt.Sprintf("rmdir %s", cgroupPath)); err != nil {
		return fmt.Errorf("cmd exec error: %s", err.Error())
	}

	return nil
}
