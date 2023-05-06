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
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"reflect"
	"testing"
	"time"
)

func TestIsTimeout(t *testing.T) {
	type args struct {
		createTimeStr string
		durationStr   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "no auto recover",
			args: args{
				"",
				"",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "no auto recover",
			args: args{
				"",
				"5p",
			},
			wantErr: true,
		},
		{
			name: "time format error",
			args: args{
				"2023-03-0319:01:50",
				"5m",
			},
			wantErr: true,
		},
		{
			name: "timeout",
			args: args{
				func() string {
					return time.Now().Add(-2 * time.Minute).Format(model.TimeFormat)
				}(),
				"5m",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "timeout",
			args: args{
				func() string {
					return time.Now().Add(-7 * time.Minute).Format(model.TimeFormat)
				}(),
				"5m",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsTimeout(tt.args.createTimeStr, tt.args.durationStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsTimeout() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetArgs(t *testing.T) {
	type args struct {
		args []v1alpha1.ArgsUnit
		keys []string
	}
	var testArgs = []v1alpha1.ArgsUnit{
		{Key: "count", Value: "50"},
		{Key: "percent", Value: "90"},
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal",
			args: args{
				args: testArgs,
				keys: []string{"count", "percent"},
			},
			want: []string{"50", "90"},
		},
		{
			name: "more than",
			args: args{
				args: testArgs,
				keys: []string{"count", "util", "percent"},
			},
			want: []string{"50", "", "90"},
		},
		{
			name: "oneEmpty",
			args: args{
				args: testArgs,
				keys: []string{"a", "percent"},
			},
			want: []string{"", "90"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetArgs(tt.args.args, tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
