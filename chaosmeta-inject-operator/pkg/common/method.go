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

package common

import (
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"strings"
	"time"
)

func IsTimeout(createTimeStr string, durationStr string) (bool, error) {
	// duration is empty means no timeout
	if durationStr == "" {
		return false, nil
	}

	duration, err := v1alpha1.ConvertDuration(durationStr)
	if err != nil {
		return false, fmt.Errorf("get duration error: %s", err.Error())
	}

	createTime, err := time.ParseInLocation(model.TimeFormat, createTimeStr, time.Local)
	if err != nil {
		return false, fmt.Errorf("get createTime error: %s", err.Error())
	}

	return createTime.Add(duration).Before(time.Now()), nil
}

func GetArgs(args []v1alpha1.ArgsUnit, keys []string) []string {
	reList := make([]string, len(keys))
	for i, k := range keys {
		for _, unit := range args {
			if unit.Key == k {
				reList[i] = unit.Value
				break
			}
		}
	}

	return reList
}

func IsKeyUniqueErr(err error) bool {
	return strings.Index(err.Error(), "UNIQUE") >= 0 && strings.Index(err.Error(), "uid") >= 0
}

func IsNetErr(err error) bool {
	return false
}

func IsNotFoundErr(err error) bool {
	return strings.Index(err.Error(), "not found") >= 0
}
