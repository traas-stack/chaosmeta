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

package httpexecutor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HostArgsKey    = "host"
	PortArgsKey    = "port"
	SchemeArgsKey  = "scheme"
	MethodArgsKey  = "method"
	PathArgsKey    = "path"
	HeaderArgsKey  = "header"
	BodyArgsKey    = "body"
	TimeoutArgsKey = "timeout"

	SchemeHTTP  = "HTTP"
	SchemeHTTPS = "HTTPS"
	MethodGET   = "GET"
	MethodPOST  = "POST"
)

func init() {
	h, err := NewHTTPExecutor(context.Background())
	if err != nil {
		fmt.Printf("new http executor error: %s\n", err.Error())
	} else {
		v1alpha1.SetMeasureExecutor(context.Background(), v1alpha1.HTTPMeasureType, h)
	}
}

type HTTPExecutor struct {
}

func NewHTTPExecutor(ctx context.Context) (*HTTPExecutor, error) {
	return &HTTPExecutor{}, nil
}

func (e *HTTPExecutor) CheckConfig(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement) error {
	_, err := utils.GetArgsValueStr(args, HostArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	portStr, err := utils.GetArgsValueStr(args, PortArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	_, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("port args is not a int: %s", err.Error())
	}

	timeout, err := utils.GetArgsValueStr(args, TimeoutArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	_, err = strconv.Atoi(timeout)
	if err != nil {
		return fmt.Errorf("timeout args is not a int: %s", err.Error())
	}

	scheme, err := utils.GetArgsValueStr(args, SchemeArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	if scheme != SchemeHTTP && scheme != SchemeHTTPS {
		return fmt.Errorf("only support scheme: %s, %s", SchemeHTTPS, SchemeHTTP)
	}

	method, err := utils.GetArgsValueStr(args, MethodArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	if method != MethodPOST && method != MethodGET {
		return fmt.Errorf("only support method: %s, %s", MethodGET, MethodPOST)
	}

	headers, err := utils.GetArgsValueStr(args, HeaderArgsKey)
	if err == nil {
		if _, err := utils.ParseKV(headers); err != nil {
			return fmt.Errorf("header args error: %s", err.Error())
		}
	}

	if judgement.JudgeType == v1alpha1.ConnectivityJudgeType {
		if judgement.JudgeValue != v1alpha1.ConnectivityTrue && judgement.JudgeValue != v1alpha1.ConnectivityFalse {
			return fmt.Errorf("value of %s judge type only support: %s, %s", v1alpha1.ConnectivityJudgeType, v1alpha1.ConnectivityTrue, v1alpha1.ConnectivityFalse)
		}
	} else if judgement.JudgeType == v1alpha1.CodeJudgeType {
		_, err := strconv.Atoi(judgement.JudgeValue)
		if err != nil {
			return fmt.Errorf("value of %s judge type only support code int: %s", v1alpha1.CodeJudgeType, err.Error())
		}
	} else if judgement.JudgeType == v1alpha1.BodyJudgeType {
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(judgement.JudgeValue), &body); err != nil {
			return fmt.Errorf("value of %s judge type only support json format: %s", v1alpha1.BodyJudgeType, err.Error())
		}
	} else {
		return fmt.Errorf("http measure only support judge type: %s, %s, %s", v1alpha1.ConnectivityJudgeType, v1alpha1.CodeJudgeType, v1alpha1.BodyJudgeType)
	}

	return nil
}

func (e *HTTPExecutor) InitialData(ctx context.Context, args []v1alpha1.MeasureArgs) (string, error) {
	return "", nil
}

func (e *HTTPExecutor) Measure(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement, initialData string) error {
	host, _ := utils.GetArgsValueStr(args, HostArgsKey)
	port, _ := utils.GetArgsValueStr(args, PortArgsKey)
	scheme, _ := utils.GetArgsValueStr(args, SchemeArgsKey)
	method, _ := utils.GetArgsValueStr(args, MethodArgsKey)
	headers, _ := utils.GetArgsValueStr(args, HeaderArgsKey)
	body, _ := utils.GetArgsValueStr(args, BodyArgsKey)
	path, _ := utils.GetArgsValueStr(args, PathArgsKey)
	timeoutStr, _ := utils.GetArgsValueStr(args, TimeoutArgsKey)

	timeout, _ := strconv.Atoi(timeoutStr)
	code, res, err := sendRequest(scheme, host, port, path, method, body, headers, timeout)
	switch judgement.JudgeType {
	case v1alpha1.CodeJudgeType:
		if err != nil {
			return fmt.Errorf("send request error: %s", err.Error())
		}
		expectedCode, _ := strconv.Atoi(judgement.JudgeValue)
		if expectedCode != code {
			return fmt.Errorf("expect code %d, but get %d", expectedCode, code)
		}
	case v1alpha1.ConnectivityJudgeType:
		if judgement.JudgeValue == v1alpha1.ConnectivityTrue {
			return err
		} else {
			if err == nil {
				return fmt.Errorf("expect connectivity false, but get true")
			}
		}
	case v1alpha1.BodyJudgeType:
		if err != nil {
			return fmt.Errorf("send request error: %s", err.Error())
		}

		var actualJson map[string]interface{}
		if err := json.Unmarshal([]byte(res), &actualJson); err != nil {
			return fmt.Errorf("%s judge type only support json format response", v1alpha1.BodyJudgeType)
		}

		actualPathsMap := getJsonPaths(actualJson)
		var expectedJson map[string]interface{}
		_ = json.Unmarshal([]byte(judgement.JudgeValue), &expectedJson)

		expectedPathMap := getJsonPaths(expectedJson)
		for expectedKey, expectedValue := range expectedPathMap {
			actualValue, isExist := actualPathsMap[expectedKey]
			if !isExist {
				return fmt.Errorf("not exist json path: %s", expectedKey)
			}

			if expectedValue != actualValue {
				return fmt.Errorf("value of path[%s] expect %s, but get %s", expectedKey, expectedValue, actualValue)
			}
		}
	}

	return nil
}

func getJsonPaths(data map[string]interface{}) map[string]string {
	paths := make(map[string]string)
	for key, value := range data {
		if subData, ok := value.(map[string]interface{}); ok {
			subPaths := getJsonPaths(subData)
			for subKey, subValue := range subPaths {
				paths[key+"."+subKey] = subValue
			}
		} else {
			paths[key] = fmt.Sprintf("%v", value)
		}
	}
	return paths
}

func sendRequest(scheme, host, port, path, method, data, headers string, timeout int) (code int, res string, err error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	if scheme == SchemeHTTPS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	var payload = &strings.Reader{}
	if data != "" {
		payload = strings.NewReader(data)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s:%s%s", strings.ToLower(scheme), host, port, path), payload)
	if err != nil {
		return
	}

	h, _ := utils.ParseKV(headers)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	res = string(body)
	code = resp.StatusCode
	return
}
