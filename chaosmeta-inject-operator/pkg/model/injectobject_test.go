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

package model

import "testing"

func TestParseContainerID(t *testing.T) {
	type args struct {
		cID string
	}
	tests := []struct {
		name    string
		args    args
		wantR   string
		wantId  string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				cID: "pouch://123456",
			},
			wantR:   "pouch",
			wantId:  "123456",
			wantErr: false,
		},
		{
			name: "success_default",
			args: args{
				cID: "123456",
			},
			wantR:   "docker",
			wantId:  "123456",
			wantErr: false,
		},
		{
			name: "fault_empty",
			args: args{
				cID: "",
			},
			wantErr: true,
		},
		{
			name: "fault_err_format",
			args: args{
				cID: "temp://egwag://asfqwe",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotId, err := ParseContainerID(tt.args.cID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContainerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("ParseContainerID() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotId != tt.wantId {
				t.Errorf("ParseContainerID() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
