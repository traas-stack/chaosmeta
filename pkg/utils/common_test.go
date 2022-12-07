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
	"reflect"
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
			got, err := GetNumArrByList(tt.args.listStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNumArrByList() error = %v, wantErr %v, args: %v", err, tt.wantErr, tt.args)
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
				count: 4,
			},
			length: 4,
		},
		{
			args: args{
				count: 1,
			},
			length: 1,
		},
	}

	list := []int{
		0, 1, 2, 3,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNumArrByCount(tt.args.count, list); len(got) != tt.length {
				t.Errorf("getNumArrByCount() = %v, want length %v", got, tt.length)
			}
		})
	}
}
