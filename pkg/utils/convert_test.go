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

import "testing"

func Test_getValueAndUnit(t *testing.T) {
	type args struct {
		strWithUnit string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   string
		wantErr bool
	}{
		{
			args:    args{strWithUnit: ""},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{strWithUnit: "-2"},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{strWithUnit: "250mnb"},
			want:    250,
			want1:   "mnb",
			wantErr: false,
		},
		{
			args:    args{strWithUnit: "250"},
			want:    250,
			want1:   "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getValueAndUnit(tt.args.strWithUnit)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueAndUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getValueAndUnit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getValueAndUnit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetBlockKbytes(t *testing.T) {
	type args struct {
		valueStr string
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		want1   string
		wantErr bool
	}{
		{
			args:    args{valueStr: ""},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{valueStr: "250gb"},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{valueStr: "250a"},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{valueStr: "kb"},
			want:    -1,
			want1:   "",
			wantErr: true,
		},
		{
			args:    args{valueStr: "256"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256k"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256K"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256kB"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256Kb"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256KB"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256kb"},
			want:    256,
			want1:   "256k",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256m"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256M"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256mB"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256Mb"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256MB"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
		{
			args:    args{valueStr: "256mb"},
			want:    262144,
			want1:   "256m",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetBlockKbytes(tt.args.valueStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockKbytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBlockKbytes() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetBlockKbytes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
