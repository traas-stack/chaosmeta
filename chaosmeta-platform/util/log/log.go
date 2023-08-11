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

// from https://github.com/ngaut/log/blob/master/log.go
package log

import (
	"fmt"
	"go.uber.org/zap"
)

const CallerSkipDepth = 1

func SetOutput(outPutType string) {
	var err error
	option := DefaultLogger.getLogOption()
	option.OutPutType = outPutType

	DefaultLogger, err = IntiLoggerSizeRotate(option)
	if err != nil {
		panic("init default logger failed")
	}
}

func SetLevel(level string) {
	var err error
	option := DefaultLogger.getLogOption()
	option.Level = level

	DefaultLogger, err = IntiLoggerSizeRotate(option)
	if err != nil {
		panic("init default logger failed")
	}
}

func SetDefaultLogOption(logOption LogOption) {
	var err error
	DefaultLogger, err = IntiLoggerSizeRotate(logOption)
	if err != nil {
		panic("init default logger failed")
	}
}

func Info(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Info(fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Debug(fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Warn(fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Warn(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Error(fmt.Sprintf(format, v...))
}

func Panic(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Panic(fmt.Sprint(v...))
}

func Panicf(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Panic(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Fatal(fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Fatal(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Fatal(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Fatal(fmt.Sprint(v...))
}

func (l *LoggerStruct) Fatalf(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Fatal(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Panic(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Panic(fmt.Sprint(v...))
}

func (l *LoggerStruct) Panicf(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Panic(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Error(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Error(fmt.Sprint(v...))
}

func (l *LoggerStruct) Errorf(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Error(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Warning(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Warn(fmt.Sprint(v...))
}

func (l *LoggerStruct) Warningf(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Warn(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Info(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Info(fmt.Sprint(v...))
}

func (l *LoggerStruct) Infof(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Info(fmt.Sprintf(format, v...))
}

func (l *LoggerStruct) Debug(v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Debug(fmt.Sprint(v...))
}

func (l *LoggerStruct) Debugf(format string, v ...interface{}) {
	l.globalLogger.WithOptions(zap.AddCallerSkip(CallerSkipDepth)).Debug(fmt.Sprintf(format, v...))
}
