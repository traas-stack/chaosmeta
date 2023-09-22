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

package file

import (
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
)

const (
	TargetFile = "file"

	FaultFileAdd = "add"

	FaultFileAppend = "append"

	FaultFileMv = "mv"

	FaultFileDelete = "del"

	FaultFileChmod = "chmod"
	//FileExec       = "chaosmeta_file"

	BackUpDir = "/tmp/chaosmeta_backup_file"
)

func getAppendFlag(uid string) string {
	return fmt.Sprintf(" %s-%s", utils.RootName, uid)
}

func getBackupDir(uid string) string {
	return fmt.Sprintf("%s%s", BackUpDir, uid)
}
