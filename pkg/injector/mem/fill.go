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

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/memory"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	injector.Register(TargetMem, FaultMemFill, func() injector.IInjector { return &FillInjector{} })
}

type FillInjector struct {
	injector.BaseInjector
	Args    FillArgs
	Runtime FillRuntime
}

type FillArgs struct {
	Percent int    `json:"percent,omitempty"`
	Bytes   string `json:"bytes,omitempty"`
	Mode    string `json:"mode"`
}

type FillRuntime struct {
	//Pid int `json:"pid,omitempty"`
}

func (i *FillInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *FillInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *FillInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeCache
	}
}

func (i *FillInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "mem fill target percent, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")
	cmd.Flags().StringVarP(&i.Args.Bytes, "bytes", "b", "", "mem fill bytes to add, support unit: KB/MB/GB/TB（default KB）")
	i.Args.Mode = ModeRam
	if runtime.GOOS == utils.LINUX {
		cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("mem fill mode, support: %s、%s（default %s）", ModeRam, ModeCache, ModeCache))
	}
}

// Validator percent > bytes
func (i *FillInjector) Validator(ctx context.Context) error {
	if i.Args.Percent == 0 && i.Args.Bytes == "" {
		return fmt.Errorf("must provide \"percent\" or \"bytes\"")
	}

	if i.Args.Percent != 0 {
		if i.Args.Percent < 0 || i.Args.Percent > 100 {
			return fmt.Errorf("\"percent\" must be in (0,100]")
		}
	} else {
		if _, err := utils.GetKBytes(i.Args.Bytes); err != nil {
			return fmt.Errorf("\"bytes\" is invalid: %s", err.Error())
		}
	}

	if i.Args.Mode != ModeCache && i.Args.Mode != ModeRam {
		return fmt.Errorf("\"mode\" is not support: %s, only support: %s、%s", i.Args.Mode, ModeCache, ModeRam)
	}

	if i.Args.Mode == ModeCache {
		if !cmdexec.SupportCmd("fallocate") && !cmdexec.SupportCmd("dd") {
			return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill cache")
		}

		if !cmdexec.SupportCmd("mount") {
			return fmt.Errorf("not support cmd \"mount\", can not fill cache")
		}

		//fillDir := getFillDir(i.Info.Uid)
		//isExist, err := filesys.ExistPath(fillDir)
		//if err != nil {
		//	return fmt.Errorf("check tmpfs[%s] exist error: %s", fillDir, err.Error())
		//}
		//
		//if isExist {
		//	return fmt.Errorf("tmpfs[%s] exist, if another cache_fill experiment is running, please recover first", fillDir)
		//}
	}

	return i.BaseInjector.Validator(ctx)
}

func getFillDir(uid string) string {
	return fmt.Sprintf("%s%s", FillDir, uid)
}

func (i *FillInjector) Inject(ctx context.Context) error {
	logger := log.GetLogger(ctx)
	fillKBytes, err := memory.CalculateFillKBytes(ctx, i.Args.Percent, i.Args.Bytes)
	if err != nil {
		return err
	}

	logger.Debugf("need to fill mem: %dKB", fillKBytes)

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	// 获取命令并执行
	if i.Args.Mode == ModeRam {
		_, err := cmdexec.StartBashCmdAndWaitPid(ctx, fmt.Sprintf("%s %s %d %dkb %d", utils.GetToolPath(MemFillKey), i.Info.Uid, -999, fillKBytes, timeout))
		if err != nil {
			return fmt.Errorf("fill ram exec error: %s", err.Error())
		}
	} else {
		if err := memory.FillCache(ctx, fillKBytes, getFillDir(i.Info.Uid), TmpFsFile); err != nil {
			return fmt.Errorf("fill cache error: %s", err.Error())
		}
	}

	return nil
}

func (i *FillInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	if i.Args.Mode == ModeCache {
		fillDir := getFillDir(i.Info.Uid)

		isDirExist, err := filesys.ExistPath(fillDir)
		if err != nil {
			return fmt.Errorf("check tmpfs[%s] exist error: %s", fillDir, err.Error())
		}

		if isDirExist {
			return memory.UndoTmpfs(ctx, fillDir)
		}

		return nil
	} else {
		return process.CheckExistAndKillByKey(ctx, fmt.Sprintf("%s %s", MemFillKey, i.Info.Uid))
	}
}
