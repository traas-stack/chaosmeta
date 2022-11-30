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
	"strconv"
	"strings"
)

func CheckSpeedValue(sp string) error {
	_, unit, err := getValueAndUnit(sp)
	if err != nil {
		return err
	}

	if unit != "" && unit != "bit" && unit != "kbit" && unit != "mbit" && unit != "gbit" && unit != "tbit" {
		return fmt.Errorf("unit %s is not support", unit)
	}

	return nil
}

func CheckTimeValue(timeStr string) error {
	_, unit, err := getValueAndUnit(timeStr)
	if err != nil {
		return err
	}

	if unit == "us" || unit == "" || unit == "ms" || unit == "s" {
		return nil
	}

	return fmt.Errorf("unit %s is not support", unit)
}

func GetTimeSecond(timeStr string) (int64, error) {
	value, unit, err := getValueAndUnit(timeStr)
	if err != nil {
		return -1, err
	}

	if unit == "s" || unit == "" {
		return value, nil
	}

	if unit == "m" {
		return value * 60, nil
	}

	if unit == "h" {
		return value * 3600, nil
	}

	return -1, fmt.Errorf("unit %s is not support", unit)
}

func GetKBytes(byteStr string) (int64, error) {
	value, unit, err := getValueAndUnit(byteStr)
	if err != nil {
		return -1, err
	}

	if unit == "kb" || unit == "" {
		return value, nil
	}

	if unit == "mb" {
		return value * 1024, nil
	}

	if unit == "gb" {
		return value * 1024 * 1024, nil
	}

	if unit == "tb" {
		return value * 1024 * 1024 * 1024, nil
	}

	if unit == "pb" {
		return value * 1024 * 1024 * 1024 * 1024, nil
	}

	return -1, fmt.Errorf("unit %s is not support", unit)
}

func GetBytes(byteStr string) (int64, error) {
	value, unit, err := getValueAndUnit(byteStr)
	if err != nil {
		return -1, err
	}

	if unit == "b" || unit == "" {
		return value, nil
	}

	if unit == "kb" {
		return value * 1024, nil
	}

	if unit == "mb" {
		return value * 1024 * 1024, nil
	}

	if unit == "gb" {
		return value * 1024 * 1024 * 1024, nil
	}

	if unit == "tb" {
		return value * 1024 * 1024 * 1024 * 1024, nil
	}

	return -1, fmt.Errorf("unit %s is not support", unit)
}

/*GetBlockKbytes
  @Description:
  @param valueStr: input value string
  @return int64: value in unit kb
  @return string: standardized value string
  @return error
*/
func GetBlockKbytes(valueStr string) (int64, string, error) {
	value, unit, err := getValueAndUnit(valueStr)
	if err != nil {
		return -1, "", err
	}

	if unit == "kb" || unit == "" {
		unit = "k"
	}

	if unit == "mb" {
		unit = "m"
	}

	if unit != "k" && unit != "m" {
		return -1, "", fmt.Errorf("not support unit: %s", unit)
	}

	reStr := fmt.Sprintf("%d%s", value, strings.ToUpper(unit))

	if unit == "m" {
		value *= 1024
	}

	return value, reStr, nil
}

func getValueAndUnit(strWithUnit string) (int64, string, error) {
	var index = -1
	for i, unitC := range strWithUnit {
		if unitC < '0' || unitC > '9' {
			index = i
			break
		}
	}

	if index == 0 {
		return -1, "", fmt.Errorf("value is empty")
	}

	if index == -1 {
		index = len(strWithUnit)
	}

	unit := strings.ToLower(strWithUnit[index:])
	valueStr := strWithUnit[:index]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return -1, "", fmt.Errorf("value is not a valid num: %s", valueStr)
	}

	if value < 0 {
		return -1, "", fmt.Errorf("value is less than 0")
	}

	return value, unit, nil
}
