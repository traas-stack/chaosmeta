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

func TestGetBlkioConfig(t *testing.T) {
	type args struct {
		devList    []string
		rBytes     string
		wBytes     string
		rIO        int64
		wIO        int64
		cgroupPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				devList: []string{"8:0", "8:1"},
				rBytes: "200kb",
				wBytes: "500KB",
				rIO: 5,
				wIO: 6,
				cgroupPath: "/etc/sys/fs/cgroup/cpu/chaosmetad_1241q52",
			},
			want: "echo 8:0 204800 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.read_bps_device &&" +
				" echo 8:1 204800 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.read_bps_device &&" +
				" echo 8:0 512000 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.write_bps_device &&" +
				" echo 8:1 512000 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.write_bps_device &&" +
				" echo 8:0 5 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.read_iops_device &&" +
				" echo 8:1 5 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.read_iops_device &&" +
				" echo 8:0 6 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.write_iops_device &&" +
				" echo 8:1 6 > /etc/sys/fs/cgroup/cpu/chaosmetad_1241q52/blkio.throttle.write_iops_device",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBlkioConfig(tt.args.devList, tt.args.rBytes, tt.args.wBytes, tt.args.rIO, tt.args.wIO, tt.args.cgroupPath); got != tt.want {
				t.Errorf("GetBlkioConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
