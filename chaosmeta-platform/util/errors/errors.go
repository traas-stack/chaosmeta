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

package errors

import (
	"encoding/json"
	"fmt"
)

var _ Error = (*Err)(nil)

type Error interface {
	WithMessage(msg string) Error
	WithData(data interface{}) Error
	GetErrorCode() int
	GetErrorMessage() string
	ToString() string
	Error() string
	ToError() error
	CleanData() Error
	CleanMessage() Error
	IsOK() bool
}

type BaseErr struct {
	Code     int    `json:"code"`               // 业务编码
	Message  string `json:"message"`            // 错误描述
	TraceId  string `json:"trace_id"`           // traceID
	ShowType int    `json:"showType,omitempty"` // 错误信息展示方式： 0 静默处理; 1 警告; 2 错误; 4 通知; 9 跳转错误页
}

type Err struct {
	BaseErr
	Data interface{} `json:"data"` // 成功时返回的数据
	Path string      `json:"path,omitempty"`
}

func NewError(code int, msg string, ShowType int) Error {
	return &Err{
		BaseErr: BaseErr{
			Message:  msg,
			Code:     code,
			ShowType: ShowType,
		},
		Data: nil,
	}
}

func NewErrorWithPath(code int, msg string, ShowType int, path string) Error {
	return &Err{
		BaseErr: BaseErr{
			Message:  msg,
			Code:     code,
			ShowType: ShowType,
		},
		Data: nil,
		Path: path,
	}
}

func (e *Err) WithData(data interface{}) Error {
	e.Data = data
	return e
}

func (e *Err) WithMessage(msg string) Error {
	e.Message = msg
	return e
}

func (e *Err) GetErrorCode() int {
	return e.Code
}

func (e *Err) GetErrorMessage() string {
	return e.Message
}

func (e *Err) CleanData() Error {
	e.Data = nil
	return e
}

func (e *Err) CleanMessage() Error {
	e.Message = "nil"
	return e
}

// ToString 返回 JSON 格式的错误详情
func (e *Err) ToString() string {
	err := &struct {
		Code     int         `json:"code"`
		Message  string      `json:"message"`
		Data     interface{} `json:"data"`
		ShowType int         `json:"showType"`
	}{
		Code:     e.Code,
		Message:  e.Message,
		Data:     e.Data,
		ShowType: e.ShowType,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}

func (e *Err) ToError() error {
	return fmt.Errorf("code: %v, msg: %v", e.Code, e.Message)
}

func (e *Err) IsOK() bool {
	return e.Code == 200
}

func (e *Err) Error() string {
	return fmt.Sprintf("scode: %v, _msg: %v", e.Code, e.Message)
}
