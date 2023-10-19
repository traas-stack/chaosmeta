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

package cgroup

const (
	BLKIO  = "blkio"
	CPUSET = "cpuset"
	MEMORY = "memory"
)

const (
	MemUnLimit             = 9223372036854771712
	MemoryLimitInBytesFile = "memory.limit_in_bytes"
	MemoryStatFile         = "memory.stat"
	MemoryUsageInBytesFile = "memory.usage_in_bytes"
	CpusetCoreFile         = "cpuset.cpus"
	WriteBytesFile         = "blkio.throttle.write_bps_device"
	ReadBytesFile          = "blkio.throttle.read_bps_device"
	WriteIOFile            = "blkio.throttle.write_iops_device"
	ReadIOFile             = "blkio.throttle.read_iops_device"
	BlkioCgroupName        = "chaosmeta_blkio"
)
