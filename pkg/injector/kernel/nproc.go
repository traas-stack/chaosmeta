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

package kernel

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/user"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

// View the total number of processes (including threads)： ps h -Led -o user | sort | uniq -c | grep temp | awk '{print $1}'
// TODO：It needs to be stated in the document: the root user is not open for the time being, because root's nproc is too large, and it is probably the first to trigger oom

func init() {
	injector.Register(TargetKernel, FaultKernelNproc, func() injector.IInjector { return &NprocInjector{} })
}

type NprocInjector struct {
	injector.BaseInjector
	Args    NprocArgs
	Runtime NprocRuntime
}

type NprocArgs struct {
	User  string `json:"user"`
	Count int    `json:"count"`
}

type NprocRuntime struct {
	//Pid int `json:"pid,omitempty"`
}

func (i *NprocInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *NprocInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *NprocInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.User, "user", "u", "", "affected user")
	cmd.Flags().IntVarP(&i.Args.Count, "count", "c", 0, "count of proc to add（default 0, means add to nproc）, you can check nproc by \"ulimit -u\"")
}

func (i *NprocInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Count < 0 {
		return fmt.Errorf("\"count\" must larger than 0")
	}

	if _, err := user.LookupUser(i.Args.User); err != nil {
		return fmt.Errorf("\"user\" is invalid: %s", err.Error())
	}

	isExist, err := process.ExistProcessByKey(ctx, fmt.Sprintf("%s %s", FdFullKey, i.Args.User))
	if err != nil {
		return fmt.Errorf("check if running error: %s", err.Error())
	}

	if isExist {
		return fmt.Errorf("nproc experiment of user[%s] is running, please recover first", i.Args.User)
	}

	return nil
}

func (i *NprocInjector) Inject(ctx context.Context) error {
	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	return cmdexec.StartBashCmdAndWaitByUser(ctx, fmt.Sprintf("%s %s %s %d %d",
		utils.GetToolPath(NprocKey), i.Args.User, i.Args.User, i.Args.Count, timeout), i.Args.User)
}

//func (i *NprocInjector) DelayRecover(ctx context.Context, timeout int64) error {
//	return nil
//}

func (i *NprocInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	grepKey := fmt.Sprintf("%s %s", NprocKey, i.Args.User)
	pid, err := process.GetPidByKeyWithoutRunUser(ctx, grepKey)
	if err != nil {
		return fmt.Errorf("get pid from key[%s] error: %s", grepKey, err.Error())
	}

	time.Sleep(1 * time.Second)

	processKey := strconv.Itoa(pid)
	exist, err := process.ExistProcessByKey(ctx, processKey)
	if err != nil {
		return fmt.Errorf("check process exist by key[%d] error: %s", pid, err.Error())
	}

	if exist {
		return process.KillProcessByKey(ctx, processKey, process.SIGKILL)
	}

	return nil
}
