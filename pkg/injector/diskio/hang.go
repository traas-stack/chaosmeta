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

package diskio

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/spf13/cobra"
	"strings"
)

// TODO: It needs to be stated in the document that if the target process generates a child process, it will also be restricted, but it is impossible to determine which cgroup the child process should be put back when restoring, so put it back to "/user.slice"

func init() {
	injector.Register(TargetDiskIO, FaultDiskIOHang, func() injector.IInjector { return &HangInjector{} })
}

type HangInjector struct {
	injector.BaseInjector
	Args    HangArgs
	Runtime HangRuntime
}

type HangArgs struct {
	PidList string `json:"pid_list"` // 需要校验是否存在
	Key     string `json:"key"`
	DevList string `json:"dev_list"`
	Mode    string `json:"mode"`
}

type HangRuntime struct {
	OldCgroupMap map[int]string
}

func (i *HangInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *HangInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *HangInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeAll
	}
}

func (i *HangInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().StringVarP(&i.Args.PidList, "pid-list", "p", "", "target process's pid, list split by \",\", eg: 9595,9696")
	cmd.Flags().StringVarP(&i.Args.Key, "key", "k", "", "the key used to grep to get target process, the effect is equivalent to \"ps -ef | grep [key]\". if \"pid-list\" provided, \"key\" will be ignored")
	cmd.Flags().StringVarP(&i.Args.DevList, "dev-list", "d", "", "target dev list, dev represent format: \"major-dev-num:minor-dev-num\",  use \"lsblk -a | grep disk\" to get dev num, eg:\"8:0,9:1\"\"")
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("target IO mode to hang, support: %s、%s、%s（default %s）", ModeAll, ModeRead, ModeWrite, ModeAll))
}

func (i *HangInjector) Validator() error {
	pidList, err := utils.GetPidListByListStrAndKey(i.Args.PidList, i.Args.Key)
	if err != nil {
		return fmt.Errorf("\"pid-list\" or \"key\" is invalid: %s", err.Error())
	}

	if err := utils.CheckPidListCgroup(pidList); err != nil {
		return fmt.Errorf("check cgroup of %v error: %s", pidList, err.Error())
	}

	i.Args.DevList = strings.TrimSpace(i.Args.DevList)
	if _, err := utils.GetDevList(i.Args.DevList); err != nil {
		return fmt.Errorf("\"dev-list\"[%s] is invalid: %s", i.Args.DevList, err.Error())
	}

	if i.Args.Mode != ModeRead && i.Args.Mode != ModeWrite && i.Args.Mode != ModeAll {
		return fmt.Errorf("\"mode\" is not support: %s", i.Args.Mode)
	}

	return i.BaseInjector.Validator()
}

func (i *HangInjector) Inject() error {
	pidList, err := utils.GetPidListByListStrAndKey(i.Args.PidList, i.Args.Key)
	if err != nil {
		return err
	}

	i.Runtime.OldCgroupMap, err = utils.GetPidListCurCgroup(pidList)
	if err != nil {
		return fmt.Errorf("get old path error: %s", err.Error())
	}
	log.WithUid(i.Info.Uid).Debugf("old cgroup path: %v", i.Runtime.OldCgroupMap)

	devList, _ := utils.GetDevList(i.Args.DevList)

	rByte, wByte := HangBytes, HangBytes
	if i.Args.Mode == ModeRead {
		wByte = ""
	} else if i.Args.Mode == ModeWrite {
		rByte = ""
	}

	// 先new cgroup
	blkioPath := utils.GetBlkioCPath(i.Info.Uid)
	if err := utils.NewCgroup(blkioPath, utils.GetBlkioConfig(devList, rByte, wByte, 0, 0, blkioPath)); err != nil {
		return fmt.Errorf("create cgroup[%s] error: %s", blkioPath, err.Error())
	}

	// 然后加进程
	if err := utils.MovePidListToCgroup(pidList, blkioPath); err != nil {
		// need to undo, use recover?
		if err := i.Recover(); err != nil {
			log.WithUid(i.Info.Uid).Warnf("undo error: %s", err.Error())
		}

		return fmt.Errorf("move pid list to cgroup[%s] error: %s", blkioPath, err.Error())
	}

	return nil
}

func (i *HangInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	cgroupPath := utils.GetBlkioCPath(i.Info.Uid)
	isCgroupExist, err := utils.ExistPath(cgroupPath)
	if err != nil {
		return fmt.Errorf("check cgroup[%s] exist error: %s", cgroupPath, err.Error())
	}

	if !isCgroupExist {
		return nil
	}

	pidList, err := utils.GetPidStrListByCgroup(cgroupPath)
	if err != nil {
		return fmt.Errorf("fail to get pid from cgroup[%s]: %s", cgroupPath, err.Error())
	}

	for _, pid := range pidList {
		oldPath, ok := i.Runtime.OldCgroupMap[pid]
		// 目标进程产生的子进程可能会遇到这种情况
		if !ok {
			log.WithUid(i.Info.Uid).Warnf("fail to get pid[%d]'s old cgroup path, move to \"%s\" instead", pid, TmpCgroup)
			oldPath = TmpCgroup
		}

		if err := utils.MoveToCgroup(pid, fmt.Sprintf("%s%s", utils.BlkioPath, oldPath)); err != nil {
			return fmt.Errorf("recover pid[%d] error: %s", pid, err.Error())
		}
	}

	if err := utils.RemoveCgroup(cgroupPath); err != nil {
		return fmt.Errorf("remove cgroup[%s] error: %s", cgroupPath, err.Error())
	}

	return nil
}
