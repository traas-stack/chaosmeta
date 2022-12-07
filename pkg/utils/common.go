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
	"os"
	"path/filepath"
	"time"
)

// common value
const (
	NoPid        = -1
	BuilderSplit = "-"
	CmdSplit     = " && "
	PortSplit    = "-"

	RootName   = "chaosmetad"
	TimeFormat = "2006-01-02 15:04:05"
	RecoverLog = "/tmp/chaosmetad_recover.log" //TODO: Need to add log cleanup strategy
)

// os
const (
	DARWIN = "darwin"
	LINUX  = "linux"
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

func GetRunPath() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

func GetToolPath(tool string) string {
	return fmt.Sprintf("%s/tools/%s", GetRunPath(), tool)
}

func GetSleepRecoverCmd(sleepTime int64, uid string) string {
	return fmt.Sprintf("sleep %ds; %s/%s recover %s >> %s 2>&1", sleepTime, GetRunPath(), RootName, uid, RecoverLog)
}
