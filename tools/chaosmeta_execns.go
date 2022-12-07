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

package main

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cgroup"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"os"
	"strconv"
	"strings"
)

// [containerPid] [namespaces] [cmd]
func main() {
	args := os.Args

	if len(args) < 4 {
		common.ExitWithErr("must provide args: [containerPid] [namespaces] [cmd]")
	}
	pidStr, nsStr, cmdStr := args[1], args[2], args[3:]
	cPid, err := strconv.Atoi(pidStr)
	if err != nil {
		common.ExitWithErr("pid is not an integer")
	}

	nsList := strings.Split(nsStr, ",")

	mPid := os.Getpid()
	if err := cgroup.AddToProCgroup(mPid, cPid); err != nil {
		common.ExitWithErr(err.Error())
	}

	for _, ns := range nsList {
		if err := namespace.JoinProcNs(cPid, ns); err != nil {
			common.ExitWithErr(fmt.Sprintf("join ns[%s] of process[%d] error: %s", ns, cPid, err.Error()))
		}
	}

	if _, err := cmdexec.StartBashCmdAndWaitPid(strings.Join(cmdStr, " ")); err != nil {
		common.ExitWithErr(fmt.Sprintf("start process error: %s", err.Error()))
	}

	fmt.Println("[success]inject success")
}
