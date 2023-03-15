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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// common value
const (
	NoPid        = -1
	BuilderSplit = "-"
	CmdSplit     = " && "
	PortSplit    = "-"

	RootName       = "chaosmetad"
	TimeFormat     = "2006-01-02 15:04:05"
	RecoverLog     = "/tmp/chaosmetad_recover.log"
	RootCgroupPath = "/sys/fs/cgroup"
)

// TraceId for command line
var TraceId string

// os
const (
	DARWIN = "darwin"
	LINUX  = "linux"
)

const (
	CtxTraceId = "TraceId"
)

const (
	MethodValidator = "validator"
	MethodInject    = "inject"
	MethodRecover   = "recover"
)

// task status
const (
	StatusCreated   = "created"
	StatusSuccess   = "success"
	StatusError     = "error"
	StatusDestroyed = "destroyed"
)

func NewUid() string {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	return fmt.Sprintf("%s%04d", timeStr, t.Nanosecond()/1000%100000%10000)
}

func IsValidUid(uid string) error {
	if len(uid) > 36 || len(uid) < 5 {
		return fmt.Errorf("length should be in [5, 36]")
	}

	for _, letter := range uid {
		if ('a' <= letter && letter <= 'z') || ('A' <= letter && letter <= 'Z') || ('0' <= letter && letter <= '9') || letter == '-' || letter == '_' {
			continue
		} else {
			return fmt.Errorf("must consist of numbers、characters、'-' and '_'")
		}
	}

	return nil
}

func StrListContain(arr []string, target string) bool {
	for _, unit := range arr {
		if target == unit {
			return true
		}
	}

	return false
}

func GetRunPath() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

func GetToolPath(tool string) string {
	return fmt.Sprintf("%s/tools/%s", GetRunPath(), tool)
}

func GetContainerPath(tool string) string {
	return fmt.Sprintf("/tmp/%s", tool)
}

func GetSleepRecoverCmd(sleepTime int64, uid string) string {
	return fmt.Sprintf("sleep %ds; %s/%s recover %s >> %s 2>&1", sleepTime, GetRunPath(), RootName, uid, RecoverLog)
}

func GetTraceId(ctx context.Context) string {
	if ctx.Value(CtxTraceId) == nil {
		return ""
	}

	return ctx.Value(CtxTraceId).(string)
}

func GetCtxWithTraceId(ctx context.Context, traceId string) context.Context {
	//if traceId == "" {
	//	traceId = NewUuid()
	//}

	return context.WithValue(ctx, CtxTraceId, traceId)
}

func GetNumArrByList(listStr string) ([]int, error) {
	var listArr []int
	var ifExist = make(map[int]bool)
	strArr := strings.Split(listStr, ",")
	for _, unitStr := range strArr {
		unitStr = strings.TrimSpace(unitStr)
		if strings.Index(unitStr, "-") >= 0 {
			rangeArr := strings.Split(unitStr, "-")
			if len(rangeArr) != 2 {
				return nil, fmt.Errorf("core range format is error. true format: 1-3")
			}

			rangeArr[0], rangeArr[1] = strings.TrimSpace(rangeArr[0]), strings.TrimSpace(rangeArr[1])
			sCore, err := strconv.Atoi(rangeArr[0])
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", rangeArr[0], err.Error())
			}

			eCore, err := strconv.Atoi(rangeArr[1])
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", rangeArr[1], err.Error())
			}

			if sCore > eCore {
				return nil, fmt.Errorf("core range must: startIndex <= endIndex")
			}

			for i := sCore; i <= eCore; i++ {
				if i < 0 {
					return nil, fmt.Errorf("core[%d] is less than 0", i)
				}

				if !ifExist[i] {
					ifExist[i] = true
					listArr = append(listArr, i)
				}
			}
		} else {
			unitCore, err := strconv.Atoi(unitStr)
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", unitStr, err.Error())
			}

			if unitCore < 0 {
				return nil, fmt.Errorf("core[%d] is less than 0", unitCore)
			}

			if !ifExist[unitCore] {
				ifExist[unitCore] = true
				listArr = append(listArr, unitCore)
			}
		}
	}

	return listArr, nil
}

func GetNumArrByCount(count int, listArr []int) []int {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//r.Shuffle(len(listArr), func(i, j int) {
	//	listArr[i], listArr[j] = listArr[j], listArr[i]
	//})

	return listArr[:count]
}
