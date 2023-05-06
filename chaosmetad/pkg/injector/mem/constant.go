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

package mem

const (
	TargetMem = "mem"

	FaultMemFill = "fill"
	FaultMemOOM  = "oom"
	PercentOOM   = 200

	ModeRam   = "ram"
	ModeCache = "cache"

	FillDir   = "/tmp/chaosmeta_mem_tmpfs"
	OOMDir    = "/tmp/chaosmeta_oom_tmpfs"
	TmpFsFile = "chaosmeta_mem_tmpfs"

	MemFillKey = "chaosmeta_memfill"

	MemExec = "chaosmeta_mem"
)
