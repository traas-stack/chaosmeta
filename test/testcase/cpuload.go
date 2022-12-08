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

package testcase

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector/cpu"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/test/common"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var cpuLoadSleepTime = 2 * time.Second

func GetCpuLoadTest() []common.TestCase {
	ctx := context.Background()
	var tempCaseList = []common.TestCase{
		{
			Args:  "greogh3wg",
			Error: true,
		},
		{
			Args:  "-c -1",
			Error: true,
		},
		{
			Args: "-c 0",
			Check: func() error {
				return checkProcessCountByKey(ctx, cpu.CpuLoadKey, runtime.NumCPU()*4+1)
			},
		},
		{
			Args: "",
			Check: func() error {
				return checkProcessCountByKey(ctx, cpu.CpuLoadKey, runtime.NumCPU()*4+1)
			},
		},
		{
			Args: "-c 5",
			Check: func() error {
				return checkProcessCountByKey(ctx, cpu.CpuLoadKey, 6)
			},
		},
		{
			Args: "-c 100",
			Check: func() error {
				return checkProcessCountByKey(ctx, cpu.CpuLoadKey, 101)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "cpu"
		tempCaseList[i].Fault = "load"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		tempCaseList[i].CheckRecover = func() error {
			return checkProcessCountByKey(ctx, cpu.CpuLoadKey, 0)
		}
	}

	return tempCaseList
}

func checkProcessCountByKey(ctx context.Context, key string, count int) error {
	time.Sleep(cpuLoadSleepTime)
	re, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("ps -ef | grep %s | grep -v grep | wc -l", key))
	if err != nil {
		return fmt.Errorf("cmd run error: %s", err.Error())
	}

	nowCount := strings.TrimSpace(string(re))
	if nowCount != strconv.Itoa(count) {
		return fmt.Errorf("expected count: %d, now count: %s", count, nowCount)
	}

	return nil
}
