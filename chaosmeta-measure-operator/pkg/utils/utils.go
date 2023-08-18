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
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"strconv"
	"strings"
	"time"
)

func ParseKV(str string) (map[string]string, error) {
	kvList := strings.Split(str, v1alpha1.KVListSplit)
	var re = make(map[string]string)
	for _, unit := range kvList {
		kvs := strings.Split(strings.TrimSpace(unit), v1alpha1.KVSplit)
		if len(kvs) != 2 {
			return nil, fmt.Errorf("%s is invalid format, expected format: k1:v1,k2:v2", unit)
		}
		re[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
	}

	return re, nil
}

func GetArgsValueStr(args []v1alpha1.MeasureArgs, key string) (string, error) {
	for i := range args {
		if args[i].Key == key {
			return args[i].Value, nil
		}
	}

	return "", fmt.Errorf("value of key[%s] not found", key)
}

func GetIntervalValue(str string) (float64, float64, error) {
	var (
		left  float64
		right float64
		err   error
	)

	valueList := strings.Split(str, v1alpha1.JudgeValueSplit)
	if len(valueList) == 1 {
		left, err = strconv.ParseFloat(valueList[0], 64)
		if err != nil {
			return left, right, fmt.Errorf("%s if not a float: %s", valueList[0], err.Error())
		}
		right = left
	} else if len(valueList) == 2 {
		if valueList[0] == "" && valueList[1] == "" {
			return left, right, fmt.Errorf("must provide at least one value in a interval")
		}

		if valueList[0] == "" {
			left = v1alpha1.IntervalMin
		} else {
			left, err = strconv.ParseFloat(valueList[0], 64)
			if err != nil {
				return left, right, fmt.Errorf("%s if not a float: %s", valueList[0], err.Error())
			}
		}

		if valueList[1] == "" {
			right = v1alpha1.IntervalMax
		} else {
			right, err = strconv.ParseFloat(valueList[1], 64)
			if err != nil {
				return left, right, fmt.Errorf("%s if not a float: %s", valueList[1], err.Error())
			}
		}
	} else {
		return left, right, fmt.Errorf("too many value, expect 1 or 2, but get %d", len(valueList))
	}

	if left > right {
		return left, right, fmt.Errorf("left value[%f] should not larger than right value[%f] in a interval", left, right)
	}

	return left, right, nil
}

func IfMeetInterval(nowValue, left, right float64) error {
	if nowValue >= left && nowValue <= right {
		return nil
	} else {
		intervalStr := fmt.Sprintf("[%f, %f]", left, right)
		if left == v1alpha1.IntervalMin {
			intervalStr = fmt.Sprintf("<= %f", right)
		} else if right == v1alpha1.IntervalMax {
			intervalStr = fmt.Sprintf(">= %f", left)
		}
		return fmt.Errorf("now value is %f, not in expected interval: %s", nowValue, intervalStr)
	}
}

func CheckSum(msg []byte) uint16 {
	sum := uint32(0)

	// Calculate cumulative sum
	for i := 0; i < len(msg)-1; i += 2 {
		sum += uint32(msg[i])<<8 | uint32(msg[i+1])
	}

	// If the message length is odd, add the last byte to the cumulative sum
	if len(msg)%2 == 1 {
		sum += uint32(msg[len(msg)-1])
	}

	// Add the upper 16 bits of the sum to the lower 16 bits
	sum = (sum >> 16) + (sum & 0xffff)

	// Add the carry to the lower 16 bits and invert to get the checksum
	checksum := uint16(^sum)

	return checksum
}

func NewUid() string {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	return fmt.Sprintf("%s%04d", timeStr, t.Nanosecond()/1000%100000%10000)
}

func IsTimeout(createTimeStr string, durationStr string) (bool, error) {
	// duration is empty means no timeout
	if durationStr == "" {
		return false, nil
	}

	duration, err := v1alpha1.ConvertDuration(durationStr)
	if err != nil {
		return false, fmt.Errorf("get duration error: %s", err.Error())
	}

	createTime, err := time.ParseInLocation(v1alpha1.TimeFormat, createTimeStr, time.Local)
	if err != nil {
		return false, fmt.Errorf("get createTime error: %s", err.Error())
	}

	return createTime.Add(duration).Before(time.Now()), nil
}
