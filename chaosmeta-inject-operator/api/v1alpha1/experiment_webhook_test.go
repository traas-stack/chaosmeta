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

package v1alpha1

import (
	"testing"
	"time"
)

func TestConvertDuration(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "value error",
			args: args{
				d: "fe3",
			},
			wantErr: true,
		},
		{
			name: "unit error",
			args: args{
				d: "5p",
			},
			wantErr: true,
		},
		{
			name: "s, true",
			args: args{
				d: "5s",
			},
			want:    time.Second * 5,
			wantErr: false,
		},
		{
			name: "m, true",
			args: args{
				d: "5m",
			},
			want:    time.Minute * 5,
			wantErr: false,
		},
		{
			name: "h, true",
			args: args{
				d: "10h",
			},
			want:    time.Hour * 10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDuration(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
