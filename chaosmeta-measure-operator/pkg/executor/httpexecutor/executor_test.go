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
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"reflect"
	"testing"
)

func TestHTTPExecutor_CheckConfig(t *testing.T) {
	testCases := []struct {
		name      string
		args      []v1alpha1.MeasureArgs
		judgement v1alpha1.Judgement
		wantErr   bool
	}{
		{
			name: "missing host",
			args: []v1alpha1.MeasureArgs{
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: v1alpha1.ConnectivityJudgeType, JudgeValue: v1alpha1.ConnectivityTrue},
			wantErr:   true,
		},
		{
			name: "invalid port",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "notanint"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: v1alpha1.ConnectivityJudgeType, JudgeValue: v1alpha1.ConnectivityTrue},
			wantErr:   true,
		},
		{
			name: "unsupported scheme",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: "notascheme"},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: v1alpha1.ConnectivityJudgeType, JudgeValue: v1alpha1.ConnectivityTrue},
			wantErr:   true,
		},
		{
			name: "unsupported method",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: "notamethod"},
				{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: v1alpha1.ConnectivityJudgeType, JudgeValue: v1alpha1.ConnectivityTrue},
			wantErr:   true,
		},
		{
			name: "invalid header",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Invalid-Header-Format"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: v1alpha1.ConnectivityJudgeType, JudgeValue: v1alpha1.ConnectivityTrue},
			wantErr:   true,
		},
		{
			name: "unsupported judge type",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
				{Key: BodyArgsKey, Value: "{\"name\": \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{JudgeType: "unsupported", JudgeValue: "unsupported"},
			wantErr:   true,
		},
		{
			name: "valid request",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type : application/json"},
				{Key: BodyArgsKey, Value: "{\"name\" : \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{
				JudgeType:  v1alpha1.CodeJudgeType,
				JudgeValue: "200",
			},
			wantErr: false,
		},
		{
			name: "valid request with relative percent judge type",
			args: []v1alpha1.MeasureArgs{
				{Key: HostArgsKey, Value: "example.com"},
				{Key: PortArgsKey, Value: "8080"},
				{Key: SchemeArgsKey, Value: SchemeHTTP},
				{Key: MethodArgsKey, Value: MethodGET},
				{Key: HeaderArgsKey, Value: "Content-Type : application/json"},
				{Key: BodyArgsKey, Value: "{\"name\" : \"test\"}"},
				{Key: TimeoutArgsKey, Value: "10"},
			},
			judgement: v1alpha1.Judgement{
				JudgeType:  v1alpha1.RelativePercentJudgeType,
				JudgeValue: "10",
			},
			wantErr: false,
		},
	}

	var httpexecutor = &HTTPExecutor{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := httpexecutor.CheckConfig(context.Background(), tc.args, tc.judgement)
			if tc.wantErr && err == nil {
				t.Errorf("CheckConfig(%v,%v) got nil error; want error", tc.args, tc.judgement)
			}
		})
	}
}

func Test_getJsonPaths(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
		"obj": map[string]interface{}{
			"hello": "world",
			"int":   456,
		},
	}
	expectedPaths := map[string]string{
		"foo":       "bar",
		"baz":       "123",
		"obj.hello": "world",
		"obj.int":   "456",
	}
	paths := getJsonPaths(data)

	if !reflect.DeepEqual(paths, expectedPaths) {
		t.Errorf("getJsonPaths() failed: got %v, expected %v", paths, expectedPaths)
	}
}

func TestHTTPExecutor_Measure(t *testing.T) {
	e := &HTTPExecutor{}
	t.Run("CodeJudgeType", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 200, "", nil
		})
		defer patch.Reset()
		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.CodeJudgeType,
			JudgeValue: "200",
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.NoError(t, err)
	})

	t.Run("ConnectivityJudgeType - true", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 0, "", context.DeadlineExceeded
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.ConnectivityJudgeType,
			JudgeValue: v1alpha1.ConnectivityTrue,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.Error(t, err)
	})

	t.Run("ConnectivityJudgeType - false", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 0, "", context.DeadlineExceeded
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "invalidhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.ConnectivityJudgeType,
			JudgeValue: v1alpha1.ConnectivityFalse,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.NoError(t, err)
	})

	t.Run("BodyJudgeType", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 200, `{"foo":"bar"}`, nil
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.BodyJudgeType,
			JudgeValue: `{"foo":"bar"}`,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.NoError(t, err)
	})

	t.Run("BodyJudgeType - invalid json", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 200, `{"foo":"bar"`, nil
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.BodyJudgeType,
			JudgeValue: `{"foo":"bar"}`,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.Error(t, err)
	})

	t.Run("BodyJudgeType - missing field", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 200, `{"foo":"bar"}`, nil
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.BodyJudgeType,
			JudgeValue: `{"foo":"bar","baz":"qux"}`,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.Error(t, err)
	})

	t.Run("BodyJudgeType - value not equal", func(t *testing.T) {
		patch := gomonkey.ApplyFunc(sendRequest, func(scheme, host, port, path, method, data, headers string, timeout int) (int, string, error) {
			return 200, `{"foo":"bar","baz":"qux"}`, nil
		})
		defer patch.Reset()

		args := []v1alpha1.MeasureArgs{
			{Key: HostArgsKey, Value: "localhost"},
			{Key: PortArgsKey, Value: "8080"},
			{Key: SchemeArgsKey, Value: "http"},
			{Key: MethodArgsKey, Value: "GET"},
			{Key: HeaderArgsKey, Value: "Content-Type: application/json"},
			{Key: PathArgsKey, Value: "/api/v1/users"},
			{Key: TimeoutArgsKey, Value: "5000"},
		}
		judgement := v1alpha1.Judgement{
			JudgeType:  v1alpha1.BodyJudgeType,
			JudgeValue: `{"foo":"bars"}`,
		}
		err := e.Measure(context.Background(), args, judgement, "")
		assert.Error(t, err)
	})
}
