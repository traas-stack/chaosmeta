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
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"reflect"
	"testing"
	"time"
)

func TestParseKV(t *testing.T) {
	testCases := []struct {
		input       string
		expected    map[string]string
		expectedErr error
	}{
		{
			input: "k1:v1,k2:v2,k3:v3",
			expected: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
			},
			expectedErr: nil,
		},
		{
			input:       "k1:v1,k2",
			expected:    nil,
			expectedErr: fmt.Errorf("k2 is invalid format, expected format: k1:v1,k2:v2"),
		},
	}

	for _, tc := range testCases {
		actual, err := ParseKV(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) || !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("ParseKV(%s) = (%v, %v), expected (%v, %v)", tc.input, actual, err, tc.expected, tc.expectedErr)
		}
	}
}

func TestGetIntervalValue(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantL   float64
		wantR   float64
		wantErr bool
	}{
		{
			name:  "normal case",
			args:  args{str: "1.0,2.0"},
			wantL: 1.0,
			wantR: 2.0,
		},
		{
			name:  "single value",
			args:  args{str: "1.0"},
			wantL: 1.0,
			wantR: 1.0,
		},
		{
			name:    "empty string",
			args:    args{str: ""},
			wantErr: true,
		},
		{
			name:    "invalid input",
			args:    args{str: "1.0,a"},
			wantErr: true,
		},
		{
			name:    "too many values",
			args:    args{str: "1.0,2.0,3.0"},
			wantErr: true,
		},
		{
			name:  "left empty",
			args:  args{str: ",2"},
			wantL: v1alpha1.IntervalMin,
			wantR: 2.0,
		},
		{
			name:  "right empty",
			args:  args{str: "1.0,"},
			wantL: 1.0,
			wantR: v1alpha1.IntervalMax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL, gotR, err := GetIntervalValue(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntervalValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (gotL != tt.wantL || gotR != tt.wantR) {
				t.Errorf("GetIntervalValue() = (%v, %v), want (%v, %v)", gotL, gotR, tt.wantL, tt.wantR)
			}
		})
	}
}

func TestIfMeetInterval(t *testing.T) {
	type args struct {
		nowValue float64
		left     float64
		right    float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				nowValue: 3.5,
				left:     2,
				right:    4,
			},
			wantErr: false,
		},
		{
			name: "test case 2",
			args: args{
				nowValue: 1,
				left:     2,
				right:    4,
			},
			wantErr: true,
		},
		{
			name: "test case 3",
			args: args{
				nowValue: 5,
				left:     2,
				right:    4,
			},
			wantErr: true,
		},
		{
			name: "test case 4",
			args: args{
				nowValue: 2,
				left:     2,
				right:    4,
			},
			wantErr: false,
		},
		{
			name: "test case 5",
			args: args{
				nowValue: 4,
				left:     2,
				right:    4,
			},
			wantErr: false,
		},
		{
			name: "larger",
			args: args{
				nowValue: 10000,
				left:     2,
				right:    v1alpha1.IntervalMax,
			},
			wantErr: false,
		},
		{
			name: "less",
			args: args{
				nowValue: -9999,
				left:     v1alpha1.IntervalMin,
				right:    4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IfMeetInterval(tt.args.nowValue, tt.args.left, tt.args.right); (err != nil) != tt.wantErr {
				t.Errorf("IfMeetInterval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckSum(t *testing.T) {
	testCases := []struct {
		name string
		msg  []byte
		want uint16
	}{
		{
			name: "empty message",
			msg:  []byte{},
			want: 65535,
		},
		{
			name: "even length message",
			msg:  []byte{0x01, 0x02, 0x03, 0x04},
			want: 64505,
		},
		{
			name: "odd length message",
			msg:  []byte{0x01, 0x02, 0x03},
			want: 65274,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckSum(tc.msg)
			if got != tc.want {
				t.Errorf("CheckSum(%v) = %v, want %v", tc.msg, got, tc.want)
			}
		})
	}
}

func TestGetArgsValueStr(t *testing.T) {
	testCases := []struct {
		name      string
		args      []v1alpha1.MeasureArgs
		key       string
		wantValue string
		wantErr   bool
	}{
		{
			name:    "empty args",
			args:    []v1alpha1.MeasureArgs{},
			key:     "foo",
			wantErr: true,
		},
		{
			name: "key not found",
			args: []v1alpha1.MeasureArgs{
				{Key: "bar", Value: "value1"},
				{Key: "baz", Value: "value2"},
			},
			key:     "foo",
			wantErr: true,
		},
		{
			name: "key found",
			args: []v1alpha1.MeasureArgs{
				{Key: "foo", Value: "value1"},
				{Key: "bar", Value: "value2"},
			},
			key:       "foo",
			wantValue: "value1",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetArgsValueStr(tc.args, tc.key)
			if tc.wantErr != (err != nil) {
				t.Errorf("GetArgsValueStr(%v, %s) error = %v, wantErr %v", tc.args, tc.key, err, tc.wantErr)
				return
			}
			if !tc.wantErr && got != tc.wantValue {
				t.Errorf("GetArgsValueStr(%v, %s) = %s, want %s", tc.args, tc.key, got, tc.wantValue)
			}
		})
	}
}

func TestIsTimeout(t *testing.T) {
	tests := []struct {
		name          string
		createTimeStr string
		durationStr   string
		wantIsTimeout bool
		wantErr       bool
	}{
		{
			name:          "duration empty",
			createTimeStr: "2022-01-01 00:00:00",
			durationStr:   "",
			wantIsTimeout: false,
			wantErr:       false,
		},
		{
			name:          "invalid duration",
			createTimeStr: "2022-01-01 00:00:00",
			durationStr:   "abc",
			wantIsTimeout: false,
			wantErr:       true,
		},
		{
			name:          "invalid createTime",
			createTimeStr: "invalidTime",
			durationStr:   "10s",
			wantIsTimeout: false,
			wantErr:       true,
		},
		{
			name:          "not timeout",
			createTimeStr: time.Now().Format(v1alpha1.TimeFormat),
			durationStr:   "10s",
			wantIsTimeout: false,
			wantErr:       false,
		},
		{
			name:          "timeout",
			createTimeStr: "2022-01-01 00:00:00",
			durationStr:   "10s",
			wantIsTimeout: true,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsTimeout, err := IsTimeout(tt.createTimeStr, tt.durationStr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantIsTimeout, gotIsTimeout)
		})
	}
}
