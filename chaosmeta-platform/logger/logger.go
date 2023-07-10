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

package logger

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

var (
	FileName = "logs/chaosmeta-platform.log"
	MaxDays  = 7
	Daily    = true
)

const (
	TraceIdKey = "TraceId"
)

func Setup() {
	if err := logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","daily":%t,"maxdays":%d}`, FileName, Daily, MaxDays)); err != nil {
		panic(any(fmt.Sprintf("set logger error: %s", err.Error())))
	}
}

func Debug(ctx context.Context, msg string) {
	traceId := getTraceId(ctx)
	logs.Debug(fmt.Sprintf("[trace: %s] %s", traceId, msg))
}

func Info(ctx context.Context, msg string) {
	traceId := getTraceId(ctx)
	logs.Info(fmt.Sprintf("[trace: %s] %s", traceId, msg))
}

func Warn(ctx context.Context, msg string) {
	traceId := getTraceId(ctx)
	logs.Warn(fmt.Sprintf("[trace: %s] %s", traceId, msg))
}

func Error(ctx context.Context, msg string) {
	traceId := getTraceId(ctx)
	logs.Error(fmt.Sprintf("[trace: %s] %s", traceId, msg))
}

func getTraceId(ctx context.Context) string {
	if ctx.Value(TraceIdKey) == nil {
		return "system"
	}

	return ctx.Value(TraceIdKey).(string)
}

//logs.Debug("my book is bought in the year of ", 2016)
//logs.Info("this %s cat is %v years old", "yellow", 3)
//logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
//logs.Error(1024, "is a very", "good game")
//logs.Critical("oh,crash")
