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
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
	"os"
)

// TODO：It needs to be explained in the document: 1. When the maxfd of "fill" mode is too large, oom may occur first instead of fd full; 2. It can only affect the fd acquisition of non-root processes, and does not affect root user
func init() {
	injector.Register(TargetKernel, FaultKernelFdfull, func() injector.IInjector { return &FdfullInjector{} })
}

type FdfullInjector struct {
	injector.BaseInjector
	Args    FdfullArgs
	Runtime FdfullRuntime
}

type FdfullArgs struct {
	Count int    `json:"count"`
	Mode  string `json:"mode"`
}

type FdfullRuntime struct {
	FileMax int `json:"file_max,omitempty"`
}

func (i *FdfullInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *FdfullInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *FdfullInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeFileMax
	}
}

func (i *FdfullInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("mode to make full of fd. \"%s\"(default): change config of max fd, \"%s\": add fd to max of os", ModeFileMax, ModeFdFill))
	cmd.Flags().IntVarP(&i.Args.Count, "count", "c", 0, fmt.Sprintf("count of fd to fill, args of \"%s\" mode（default 0, means add to max）, you can check by \"cat %s\"", ModeFdFill, FileNrPath))
}

func (i *FdfullInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Mode != ModeFdFill && i.Args.Mode != ModeFileMax {
		return fmt.Errorf(fmt.Sprintf("\"mode\" not support: %s, only support: %s, %s", i.Args.Mode, ModeFdFill, ModeFileMax))
	}

	if i.Args.Mode == ModeFdFill && i.Args.Count < 0 {
		return fmt.Errorf("\"count\" must larger than 0")
	}

	nowFd, maxFd, err := filesys.GetKernelFdStatus(ctx)
	if err != nil {
		return fmt.Errorf("get kernel max fd count error: %s", err.Error())
	}

	if nowFd >= maxFd {
		return fmt.Errorf("now fd[%d] is larger than max fd[%d], no need to inject", nowFd, maxFd)
	}

	return nil
}

func (i *FdfullInjector) getFdFullDir() string {
	return fmt.Sprintf("%s%s", FdFullDir, i.Info.Uid)
}

func (i *FdfullInjector) fdfill(ctx context.Context, maxFd, nowFd int) error {
	proFd, err := filesys.GetProMaxFd(ctx)
	if err != nil {
		return fmt.Errorf("get process max fd count error: %s", err.Error())
	}

	if i.Args.Count == 0 || i.Args.Count > maxFd {
		i.Args.Count = maxFd - nowFd + proFd
	}

	step := proFd - 10

	fdFullDir := i.getFdFullDir()
	if err := filesys.CreateFdFile(ctx, fdFullDir, FdFullFile, step); err != nil {
		return fmt.Errorf("create tmp file[%s] error: %s", fdFullDir, err.Error())
	}

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	proCount := i.Args.Count/step + 1

	for proCount > 0 {
		if err := cmdexec.StartBashCmd(ctx, fmt.Sprintf("%s %s %s %s %d %d %d",
			utils.GetToolPath(FdFullKey), i.Info.Uid, fdFullDir, FdFullFile, 0, step, timeout)); err != nil {
			return fmt.Errorf("start fd full error: %s", err.Error())
		}

		proCount--
	}

	return nil
}

func changeFileMax(ctx context.Context, fileMax int) error {
	return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo %d > %s", fileMax, FileMaxPath))
}

func (i *FdfullInjector) Inject(ctx context.Context) error {
	nowFd, maxFd, _ := filesys.GetKernelFdStatus(ctx)
	if i.Args.Mode == ModeFdFill {
		if err := i.fdfill(ctx, maxFd, nowFd); err != nil {
			return i.getErrWithUndo(ctx, err.Error())
		}
	} else {
		targetFileMax := nowFd - 2000
		if targetFileMax < 3 {
			targetFileMax = 3
		}

		if err := changeFileMax(ctx, targetFileMax); err != nil {
			return i.getErrWithUndo(ctx, err.Error())
		}

		i.Runtime.FileMax = maxFd
	}

	return nil
}

func (i *FdfullInjector) getErrWithUndo(ctx context.Context, msg string) error {
	if err := i.Recover(ctx); err != nil {
		log.GetLogger(ctx).Warnf("undo error: %s", err.Error())
	}

	return fmt.Errorf(msg)
}

func getFdfullKey(uid string) string {
	return fmt.Sprintf("%s %s", FdFullKey, uid)
}

func (i *FdfullInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	logger := log.GetLogger(ctx)
	if i.Args.Mode == ModeFdFill {
		processKey, fdFullDir := getFdfullKey(i.Info.Uid), i.getFdFullDir()
		isProExist, err := process.ExistProcessByKey(ctx, processKey)
		if err != nil {
			return fmt.Errorf("check process exist by key[%s] error: %s", processKey, err.Error())
		}

		if isProExist {
			if err := process.KillProcessByKey(ctx, processKey, process.SIGKILL); err != nil {
				logger.Warnf("kill process by key[%s] error: %s", processKey, err.Error())
			}
		}

		isDirExist, err := filesys.ExistPath(fdFullDir)
		if err != nil {
			return fmt.Errorf("check dir[%s] exist error: %s", fdFullDir, err.Error())
		}

		if isDirExist {
			return os.RemoveAll(fdFullDir)
		}

		return nil
	} else {
		_, maxFd, err := filesys.GetKernelFdStatus(ctx)
		if err != nil {
			return fmt.Errorf("get kernel max fd count error: %s", err.Error())
		}

		if i.Runtime.FileMax != maxFd {
			return changeFileMax(ctx, i.Runtime.FileMax)

		}
	}

	return nil
}
