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

package cgroup

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
)

func GetBlkioConfig(devList []string, rBytes, wBytes string, rIO, wIO int64, cgroupPath string) string {
	var re = ""
	if rBytes != "" {
		b, _ := utils.GetBytes(rBytes)
		re += getThrottleDeviceCmdStr(devList, b, fmt.Sprintf("%s/%s", cgroupPath, ReadBytesFile))
	}

	if wBytes != "" {
		b, _ := utils.GetBytes(wBytes)
		re += getThrottleDeviceCmdStr(devList, b, fmt.Sprintf("%s/%s", cgroupPath, WriteBytesFile))
	}

	if rIO != 0 {
		re += getThrottleDeviceCmdStr(devList, rIO, fmt.Sprintf("%s/%s", cgroupPath, ReadIOFile))
	}

	if wIO != 0 {
		re += getThrottleDeviceCmdStr(devList, wIO, fmt.Sprintf("%s/%s", cgroupPath, WriteIOFile))
	}

	log.GetLogger().Debugf("blkio config: %s", re)
	return re[:len(re)-len(utils.CmdSplit)]
}

func getThrottleDeviceCmdStr(devList []string, value int64, filename string) string {
	var re string
	for _, unitDec := range devList {
		re += fmt.Sprintf("echo %s %d > %s%s", unitDec, value, filename, utils.CmdSplit)
	}

	return re
}
