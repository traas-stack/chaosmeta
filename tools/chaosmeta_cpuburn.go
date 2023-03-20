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
	"github.com/traas-stack/chaosmetad/pkg/utils/containercgroup"
	"github.com/traas-stack/chaosmetad/tools/common"
	"os"
	"strconv"
	"time"
)

var nowTargetPercent, worktime, sleeptime int

// uid core percent pid timeout
func main() {
	args := os.Args
	if len(args) < 6 {
		common.ExitWithErr("must provide 5 args: uid、core、percent、target pid、timeout")
	}

	coreStr, percentStr, targetPidStr, timeoutStr := args[2], args[3], args[4], args[5]
	core, err := strconv.Atoi(coreStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("core[%s] is not a num: %s", coreStr, err.Error()))
	}

	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("percent[%s] is not a num: %s", percentStr, err.Error()))
	}

	targetPid, err := strconv.Atoi(targetPidStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("pid[%s] is not a num: %s", targetPidStr, err.Error()))
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("timeout[%s] is not a num: %s", timeoutStr, err.Error()))
	}

	if percent < 100 {
		go adjustPercent(targetPid, core, percent)
	} else {
		nowTargetPercent = 100
	}

	go burnCpu()

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}

func burnCpu() {
	if nowTargetPercent == 100 {
		for {
		}
	}

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

func adjustPercent(targetPid, core, maxPercent int) {
	for {
		p, err := containercgroup.CalculateNowPercent(targetPid)
		//p, err := cpu.Percent(2*time.Second, true)
		if err != nil {
			common.ExitWithErr(fmt.Sprintf("get cpu usage error: %s", err.Error()))
		}

		needAdd := maxPercent - int(p[core])
		if needAdd+nowTargetPercent < 0 {
			nowTargetPercent = 0
		} else {
			nowTargetPercent += needAdd
		}
	}
}
