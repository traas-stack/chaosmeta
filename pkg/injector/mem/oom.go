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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	injector.Register(TargetMem, FaultMemOOM, func() injector.IInjector { return &OOMInjector{} })
}

type OOMInjector struct {
	injector.BaseInjector
	Args    OOMArgs
	Runtime OOMRuntime
}

type OOMArgs struct {
	Mode string `json:"mode,omitempty"`
}

type OOMRuntime struct {
	Pid int `json:"pid,omitempty"`
}

func (i *OOMInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *OOMInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *OOMInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeCache
	}
}

func (i *OOMInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)
	i.Args.Mode = ModeRam
	if runtime.GOOS == utils.LINUX {
		cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("mem fill mode, support: %s、%s（default %s）", ModeRam, ModeCache, ModeCache))
	}
}

func (i *OOMInjector) Validator() error {
	if i.Args.Mode != ModeCache && i.Args.Mode != ModeRam {
		return fmt.Errorf("\"mode\" is not support: %s, only support: %s、%s", i.Args.Mode, ModeCache, ModeRam)
	}

	if i.Args.Mode == ModeCache {
		if !utils.SupportCmd("fallocate") && !utils.SupportCmd("dd") {
			return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill cache")
		}

		if !utils.SupportCmd("mount") {
			return fmt.Errorf("not support cmd \"mount\", can not fill cache")
		}

		fillDir := getOOMDir(i.Info.Uid)
		isExist, err := utils.ExistPath(fillDir)
		if err != nil {
			return fmt.Errorf("check tmpfs[%s] exist error: %s", fillDir, err.Error())
		}

		if isExist {
			return fmt.Errorf("tmpfs[%s] exist, if another cache_fill experiment is running, please recover first", fillDir)
		}
	}

	return i.BaseInjector.Validator()
}

func getOOMDir(uid string) string {
	return fmt.Sprintf("%s%s", OOMDir, uid)
}

func (i *OOMInjector) Inject() error {

	fillKBytes, err := utils.CalculateFillKBytes(PercentOOM, "")
	if err != nil {
		return err
	}

	log.WithUid(i.Info.Uid).Debugf("need to fill mem: %dKB", fillKBytes)

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	// 获取命令并执行
	if i.Args.Mode == ModeRam {
		pid, err := utils.FillRam(MemFillKey, fillKBytes, i.Info.Uid, timeout)
		if err != nil {
			return fmt.Errorf("fill ram error: %s", err.Error())
		}

		i.Runtime.Pid = pid
	} else {
		if err := utils.FillCache(fillKBytes, getOOMDir(i.Info.Uid), TmpFsFile); err != nil {
			return fmt.Errorf("fill cache error: %s", err.Error())
		}
	}

	return nil
}

func (i *OOMInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	if i.Args.Mode == ModeCache {
		fillDir := getOOMDir(i.Info.Uid)

		isDirExist, err := utils.ExistPath(fillDir)
		if err != nil {
			return fmt.Errorf("check tmpfs[%s] exist error: %s", fillDir, err.Error())
		}

		if isDirExist {
			return utils.UndoTmpfs(fillDir)
		}

		return nil
	} else {
		if i.Runtime.Pid == 0 {
			return nil
		}

		isPidExist, err := utils.ExistPid(i.Runtime.Pid)
		if err != nil {
			return fmt.Errorf("check pid[%d] exist error: %s", i.Runtime.Pid, err.Error())
		}

		if isPidExist {
			return utils.KillPidWithSignal(i.Runtime.Pid, utils.SIGKILL)
		}

		return nil
	}
}
