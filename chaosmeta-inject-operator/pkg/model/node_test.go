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

func TestParseNodeInfo(t *testing.T) {
	type args struct {
		nodeStr string
	}
	tests := []struct {
		name         string
		args         args
		wantNodeName string
		wantNodeIP   string
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				nodeStr: "node//1.2.3.4",
			},
			wantNodeName: "",
			wantNodeIP:   "1.2.3.4",
			wantErr:      false,
		},
		{
			name: "success",
			args: args{
				nodeStr: "node/node1/1.2.3.4",
			},
			wantNodeName: "node1",
			wantNodeIP:   "1.2.3.4",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNodeName, gotNodeIP, err := ParseNodeInfo(tt.args.nodeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNodeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNodeName != tt.wantNodeName {
				t.Errorf("ParseNodeInfo() gotNodeName = %v, want %v", gotNodeName, tt.wantNodeName)
			}
			if gotNodeIP != tt.wantNodeIP {
				t.Errorf("ParseNodeInfo() gotNodeIP = %v, want %v", gotNodeIP, tt.wantNodeIP)
			}
		})
	}
}
