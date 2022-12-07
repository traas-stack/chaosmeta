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

package main

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"github.com/shirou/gopsutil/cpu"
	"os"
	"strconv"
	"time"
)

var nowTargetPercent, worktime, sleeptime int

// uid core percent timeout
func main() {
	args := os.Args
	if len(args) < 5 {
		common.ExitWithErr("must provide 4 args: uid、core、percent、timeout")
	}

	coreStr, percentStr, timeoutStr := args[2], args[3], args[4]
	core, err := strconv.Atoi(coreStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("core[%s] is not a num: %s", coreStr, err.Error()))
	}

	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("percent[%s] is not a num: %s", percentStr, err.Error()))
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("timeout[%s] is not a num: %s", timeoutStr, err.Error()))
	}

	go adjustPercent(core, percent)
	go burnCpu()

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}

func burnCpu() {
	var starttime, endtime int64
	for {
		worktime, sleeptime = nowTargetPercent, 100-nowTargetPercent
		starttime = time.Now().UnixMicro()
		endtime = starttime

		for (endtime - starttime) < int64(worktime)*1000 {
			endtime = time.Now().UnixMicro()
		}

		time.Sleep(time.Microsecond * time.Duration(sleeptime*1000))
	}
}

func adjustPercent(core, maxPercent int) {
	for {
		// TODO: Need to change the implementation of "cpu.Percent" to only get the cpu usage rate of this container. Consider using "docker.CgroupCPU" to achieve
		p, err := cpu.Percent(2*time.Second, true)
		if err != nil {
			common.ExitWithErr(fmt.Sprintf("check cpu usage error: %s", err.Error()))
		}

		needAdd := maxPercent - int(p[core])
		if needAdd+nowTargetPercent < 0 {
			nowTargetPercent = 0
		} else {
			nowTargetPercent += needAdd
		}
	}
}
