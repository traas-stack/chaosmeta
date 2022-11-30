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

package cpu

import (
	"reflect"
	"runtime"
	"testing"
)

func Test_getNumArrByList(t *testing.T) {
	type args struct {
		listStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{name: "normal",
			args: args{
				listStr: "1-3,5",
			},
			want:    []int{1, 2, 3, 5},
			wantErr: false,
		},
		{name: "normal",
			args: args{
				listStr: "    1-3,5     ",
			},
			want:    []int{1, 2, 3, 5},
			wantErr: false,
		},
		{name: "normal",
			args: args{
				listStr: "1    -3,     5",
			},
			want:    []int{1, 2, 3, 5},
			wantErr: false,
		},
		{name: "normal",
			args: args{
				listStr: "1-3,2",
			},
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{name: "normal",
			args: args{
				listStr: "2-2,3",
			},
			want:    []int{2, 3},
			wantErr: false,
		},
		{name: "err",
			args: args{
				listStr: "-3,5",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "2-,5",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "3-2,5",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "2-3-6,5",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "2-5,99999999",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "2-5,s",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "s-5,3",
			},
			want:    nil,
			wantErr: true,
		},
		{name: "err",
			args: args{
				listStr: "2-s,3",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNumArrByList(tt.args.listStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNumArrByList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNumArrByList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNumArrByCount(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name   string
		args   args
		length int
	}{
		{
			args: args{
				count: runtime.NumCPU(),
			},
			length: runtime.NumCPU(),
		},
		{
			args: args{
				count: 1,
			},
			length: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumArrByCount(tt.args.count); len(got) != tt.length {
				t.Errorf("getNumArrByCount() = %v, want length %v", got, tt.length)
			}
		})
	}
}
