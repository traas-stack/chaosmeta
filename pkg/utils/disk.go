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
	"strings"
)

func GetDevList(devStr string) ([]string, error) {
	if devStr == "" {
		return nil, fmt.Errorf("args dev-list is empty")
	}

	devStrList := strings.Split(devStr, ",")
	for _, unit := range devStrList {
		isExist, err := existDev(unit)
		if err != nil {
			return nil, fmt.Errorf("check dev[%s] exist error: %s", unit, err.Error())
		}

		if !isExist {
			return nil, fmt.Errorf("dev[%s] is not exist", unit)
		}
	}

	return devStrList, nil
}

func existDev(devNum string) (bool, error) {
	reByte, err := RunBashCmdWithOutput(fmt.Sprintf("lsblk -a | grep disk | awk '{print $2}' | grep \"%s\" | wc -l", devNum))
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(string(reByte)) == "1" {
		return true, nil
	}

	return false, nil
}

func RunFillDisk(size int64, file string) error {
	unit := "K"
	// 100M以上使用M单位
	if size/1024 >= 100 {
		unit = "M"
		size /= 1024
	}

	if SupportCmd("fallocate") {
		return RunBashCmdWithoutOutput(fmt.Sprintf("fallocate -l %d%s %s", size, unit, file))
	}

	if SupportCmd("dd") {
		return RunBashCmdWithoutOutput(fmt.Sprintf("dd if=/dev/zero of=%s bs=1%s count=%d iflag=fullblock", file, unit, size))
	}

	return fmt.Errorf("not support \"fallocate\" and \"dd\"")
}
