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

package log

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

func Debug(v ...interface{}) {
	logs.Debug(fmt.Sprint(v...))
}

func Debugf(f string, v ...interface{}) {
	logs.Debug(f, v...)
}

func Info(v ...interface{}) {
	logs.Info(fmt.Sprint(v...))
}

func Infof(f string, v ...interface{}) {
	logs.Info(f, v...)
}

func Warn(v ...interface{}) {
	logs.Warn(fmt.Sprint(v...))
}

func Warnf(f string, v ...interface{}) {
	logs.Warn(f, v...)
}

func Error(v ...interface{}) {
	logs.Error(fmt.Sprint(v...))
}

func Errorf(f string, v ...interface{}) {
	logs.Error(f, v...)
}
