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

package cloudnativeexecutor

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getPatchLabels(t *testing.T) {
	type args struct {
		labels []byte
	}

	testMap1 := make(map[string]string)
	testMap1["a"] = "b"
	testBytes1, _ := json.Marshal(testMap1)

	testMap2 := make(map[string]string)
	testMap2["a"] = ""
	testBytes2, _ := json.Marshal(testMap2)

	testMap3 := make(map[string]string)
	testBytes3, _ := json.Marshal(testMap3)

	testMap4 := make(map[string]interface{})
	testMap4["a"] = nil
	testMap4["b"] = "gewgegwg"
	testBytes4, _ := json.Marshal(testMap4)

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{
				labels: testBytes1,
			},
			want: []byte("{\"metadata\":{\"labels\":{\"a\":\"b\"}}}"),
		},
		{
			name: "value_empty",
			args: args{
				labels: testBytes2,
			},
			want: []byte("{\"metadata\":{\"labels\":{\"a\":\"\"}}}"),
		},
		{
			name: "empty_map",
			args: args{
				labels: testBytes3,
			},
			want: []byte("{\"metadata\":{\"labels\":{}}}"),
		},
		{
			name: "empty",
			args: args{
				labels: nil,
			},
			want: []byte("{\"metadata\":{\"labels\":{}}}"),
		},
		{
			name: "null",
			args: args{
				labels: testBytes4,
			},
			want: []byte("{\"metadata\":{\"labels\":{\"a\":null,\"b\":\"gewgegwg\"}}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getPatchLabels(tt.args.labels), "getPatchLabels(%v)", tt.args.labels)
		})
	}
}

func Test_getBackupLabels(t *testing.T) {
	type args struct {
		backup    []byte
		nowLabels map[string]string
	}

	testMap1 := make(map[string]string)
	testMap1["a"] = "b"
	testMap1["a1"] = "b"

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{
				backup:    make([]byte, 0),
				nowLabels: testMap1,
			},
			want: []byte("{\"a\":null,\"a1\":null}"),
		},
		{
			name: "normal",
			args: args{
				backup:    nil,
				nowLabels: testMap1,
			},
			want: []byte("{\"a\":null,\"a1\":null}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := getBackupLabels(tt.args.backup, tt.args.nowLabels)
			assert.Equalf(t, tt.want, got, "getBackupLabels(%v, %v)", tt.args.backup, tt.args.nowLabels)
		})
	}
}
