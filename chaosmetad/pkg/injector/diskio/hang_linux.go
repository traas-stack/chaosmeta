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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/containercgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/disk"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/process"
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
	PidList string `json:"pid_list"`
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

func (i *HangInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	pidList, err := process.GetPidListByListStrAndKey(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.PidList, i.Args.Key)
	if err != nil {
		return fmt.Errorf("\"pid-list\" or \"key\" is invalid: %s", err.Error())
	}

	if err := cgroup.CheckPidListBlkioCgroup(ctx, pidList); err != nil {
		return fmt.Errorf("check cgroup of %v error: %s", pidList, err.Error())
	}

	i.Args.DevList = strings.TrimSpace(i.Args.DevList)
	if _, err := disk.GetDevList(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.DevList); err != nil {
		return fmt.Errorf("\"dev-list\"[%s] is invalid: %s", i.Args.DevList, err.Error())
	}

	if i.Args.Mode != ModeRead && i.Args.Mode != ModeWrite && i.Args.Mode != ModeAll {
		return fmt.Errorf("\"mode\" is not support: %s", i.Args.Mode)
	}

	return nil
}

func (i *HangInjector) Inject(ctx context.Context) error {
	logger := log.GetLogger(ctx)

	pidList, err := process.GetPidListByListStrAndKey(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.PidList, i.Args.Key)
	if err != nil {
		return err
	}

	i.Runtime.OldCgroupMap, err = cgroup.GetPidListCurCgroup(ctx, pidList, cgroup.BLKIO)
	if err != nil {
		return fmt.Errorf("get old path error: %s", err.Error())
	}
	logger.Debugf("old cgroup path: %v", i.Runtime.OldCgroupMap)

	devList, _ := disk.GetDevList(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.DevList)

	rByte, wByte := HangBytes, HangBytes
	if i.Args.Mode == ModeRead {
		wByte = ""
	} else if i.Args.Mode == ModeWrite {
		rByte = ""
	}

	var containerCgroup string
	if i.Info.ContainerRuntime != "" {
		containerCgroup, err = cgroup.GetContainerCgroup(ctx, i.Info.ContainerRuntime, i.Info.ContainerId)
		if err != nil {
			return fmt.Errorf("get cgroup path of container[%s] error: %s", i.Info.ContainerId, err.Error())
		}
	}

	blkioPath := cgroup.GetBlkioCPath(i.Info.Uid, containerCgroup)
	if err := cgroup.NewCgroup(ctx, blkioPath, cgroup.GetBlkioConfig(ctx, devList, rByte, wByte, 0, 0, blkioPath)); err != nil {
		if err := i.Recover(ctx); err != nil {
			logger.Warnf("undo error: %s", err.Error())
		}

		return fmt.Errorf("create cgroup[%s] error: %s", blkioPath, err.Error())
	}

	if err := cgroup.MovePidListToCgroup(ctx, pidList, blkioPath); err != nil {
		if err := i.Recover(ctx); err != nil {
			logger.Warnf("undo error: %s", err.Error())
		}

		return fmt.Errorf("move pid list to cgroup[%s] error: %s", blkioPath, err.Error())
	}

	return nil
}

func (i *HangInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	var (
		logger          = log.GetLogger(ctx)
		containerCgroup string
		err             error
		tmpPath         = TmpCgroup
	)

	if i.Info.ContainerRuntime != "" {
		containerCgroup, err = cgroup.GetContainerCgroup(ctx, i.Info.ContainerRuntime, i.Info.ContainerId)
		if err != nil {
			return fmt.Errorf("get cgroup path of container[%s] error: %s", i.Info.ContainerId, err.Error())
		}
		tmpPath = containerCgroup
	}

	cgroupPath := cgroup.GetBlkioCPath(i.Info.Uid, containerCgroup)
	isCgroupExist, err := filesys.ExistPath(cgroupPath)
	if err != nil {
		return fmt.Errorf("check cgroup[%s] exist error: %s", cgroupPath, err.Error())
	}

	if !isCgroupExist {
		return nil
	}

	pidList, err := cgroup.GetPidStrListByCgroup(ctx, cgroupPath)
	if err != nil {
		return fmt.Errorf("fail to get pid from cgroup[%s]: %s", cgroupPath, err.Error())
	}

	for _, pid := range pidList {
		oldPath, ok := i.Runtime.OldCgroupMap[pid]
		if !ok {
			logger.Warnf("fail to get pid[%d]'s old cgroup path, move to \"%s\" instead", pid, tmpPath)
			oldPath = tmpPath
		}

		if err := cgroup.MoveTaskToCgroup(ctx, pid, fmt.Sprintf("%s/%s%s", containercgroup.RootCgroupPath, cgroup.BLKIO, oldPath)); err != nil {
			return fmt.Errorf("recover pid[%d] error: %s", pid, err.Error())
		}
	}

	if err := cgroup.RemoveCgroup(ctx, cgroupPath); err != nil {
		return fmt.Errorf("remove cgroup[%s] error: %s", cgroupPath, err.Error())
	}

	return nil
}
