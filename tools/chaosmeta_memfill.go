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
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/memory"
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func parseByteValue(value int64) (int, string, error) {
	//splitIndex := len(byteStr) - 2
	//if splitIndex <= 0 {
	//	return -1, "", fmt.Errorf("byte value must contain unit: KB/MB/GB/TB")
	//}
	//
	//unit := strings.ToLower(byteStr[splitIndex:])
	//valueStr := byteStr[:splitIndex]
	//
	//value, err := strconv.ParseInt(valueStr, 10, 64)
	//if err != nil {
	//	return -1, "", fmt.Errorf("value is not a num: %s", valueStr)
	//}
	//
	//if value <= 0 {
	//	return -1, "", fmt.Errorf("value must larger than 0")
	//}
	//
	//if unit != "kb" && unit != "mb" && unit != "gb" && unit != "tb" {
	//	return -1, "", fmt.Errorf("unit only support: KB/MB/GB/TB")
	//}

	unit := "kb"

	//  Prevent the integer from being too large, as long as it is less than the maximum value of int
	for value > math.MaxInt {
		if unit == "tb" {
			return -1, "", fmt.Errorf("fill bytes is larger than support")
		} else if unit == "gb" {
			unit = "tb"
			value /= 1024
		} else if unit == "mb" {
			unit = "gb"
			value /= 1024
		} else if unit == "kb" {
			unit = "mb"
			value /= 1024
		}
	}

	return int(value), unit, nil
}

func getStrUnit(unit string) string {
	var tempStr = "a"
	var fillUnit string
	fillUnit = strings.Repeat(tempStr, 1024)
	if unit == "kb" {
		return fillUnit
	}

	fillUnit = strings.Repeat(fillUnit, 1024)
	if unit == "mb" {
		return fillUnit
	}

	fillUnit = strings.Repeat(fillUnit, 1024)
	if unit == "gb" {
		return fillUnit
	}

	fillUnit = strings.Repeat(fillUnit, 1024)
	if unit == "tb" {
		return fillUnit
	}

	return ""
}

func writeScore(scoreStr string) error {
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return fmt.Errorf("[error]score is not a valid int: %s", err.Error())
	}

	f, err := os.OpenFile("/proc/self/oom_score_adj", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open score file fail: %s", err.Error())
	}

	if _, err := f.Write([]byte(strconv.Itoa(score))); err != nil {
		return fmt.Errorf("write score error: %s", err.Error())
	}

	return nil
}

// [uid] [score] [fill bytes: KB/MB/GB/TB] [timeout second]
// [uid] [score] [percent] [bytes] [timeout second]
func main() {
	args := os.Args
	if len(args) < 5 {
		common.ExitWithErr("args must at lease 4. format: [uid] [score] [percent] [bytes] [timeout second]")
	}

	score, percentStr, bytes := args[2], args[3], args[4]

	if err := writeScore(score); err != nil {
		common.ExitWithErr(fmt.Sprintf("set score error: %s", err.Error()))
	}

	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("percent is not a num: %s", err.Error()))
	}

	fillKBytes, err := memory.CalculateFillKBytes(context.Background(), percent, bytes)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("get fill KBytes error: %s", err.Error()))
	}

	value, unit, err := parseByteValue(fillKBytes)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("parse byte value error: %s", err.Error()))
	}

	unitStr := getStrUnit(unit)
	runtime.GC()
	if value > 1 {
		_ = strings.Repeat(unitStr, value)
	}

	var timeout int
	if len(args) > 4 {
		timeout, err = strconv.Atoi(args[5])
		if err != nil {
			common.ExitWithErr(fmt.Sprintf("timeout value is not a valid int, error: %s", err.Error()))
		}
	}

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}
