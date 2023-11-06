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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/containercgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/disk"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/process"
	"strings"
)

// TODO: It needs to be stated in the document that if the target process generates a child process, it will also be restricted, but it is impossible to determine which cgroup the child process should be put back when restoring, so put it back to "/user.slice"

func init() {
	injector.Register(TargetDiskIO, FaultDiskIOLimit, func() injector.IInjector { return &LimitInjector{} })
}

type LimitInjector struct {
	injector.BaseInjector
	Args    LimitArgs
	Runtime LimitRuntime
}

type LimitArgs struct {
	PidList    string `json:"pid_list"`
	Key        string `json:"key"`
	DevList    string `json:"dev_list"`
	ReadBytes  string `json:"read_bytes,omitempty"`
	WriteBytes string `json:"write_bytes,omitempty"`
	ReadIO     int64  `json:"read_io,omitempty"`
	WriteIO    int64  `json:"write_io,omitempty"`
}

type LimitRuntime struct {
	OldCgroupMap map[int]string
}

func (i *LimitInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *LimitInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *LimitInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().StringVarP(&i.Args.PidList, "pid-list", "p", "", "target process's pid, list split by \",\", eg: 9595,9696")
	cmd.Flags().StringVarP(&i.Args.Key, "key", "k", "", "the key used to grep to get target process, the effect is equivalent to \"ps -ef | grep [key]\". if \"pid-list\" provided, \"key\" will be ignored")
	cmd.Flags().StringVarP(&i.Args.DevList, "dev-list", "d", "", "target dev list, dev represent format: \"major-dev-num:minor-dev-num\",  use \"lsblk -a | grep disk\" to get dev num, eg:\"8:0,9:1\"\"")
	cmd.Flags().StringVar(&i.Args.ReadBytes, "read-bytes", "", "limit read bytes per second, must larger than 0, support unit: B/KB/MB/GB/TB（default B）")
	cmd.Flags().Int64Var(&i.Args.ReadIO, "read-io", 0, "limit read times per second, must larger than 0")
	cmd.Flags().StringVar(&i.Args.WriteBytes, "write-bytes", "", "limit write bytes per second, must larger than 0, support unit: B/KB/MB/GB/TB（default B）")
	cmd.Flags().Int64Var(&i.Args.WriteIO, "write-io", 0, "limit write times per second, must larger than 0")
}

func (i *LimitInjector) Validator(ctx context.Context) error {
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

	if i.Args.ReadBytes == "" && i.Args.WriteBytes == "" && i.Args.ReadIO <= 0 && i.Args.WriteIO <= 0 {
		return fmt.Errorf("must provide at least one valid args of: read-bytes、write-bytes、read-io、write-io")
	}

	if i.Args.ReadBytes != "" {
		if _, err := utils.GetBytes(i.Args.ReadBytes); err != nil {
			return fmt.Errorf("\"read-bytes\"[%s] is invalid: %s", i.Args.ReadBytes, err.Error())
		}
	}

	if i.Args.WriteBytes != "" {
		if _, err := utils.GetBytes(i.Args.WriteBytes); err != nil {
			return fmt.Errorf("\"write-bytes\"[%s] is invalid: %s", i.Args.WriteBytes, err.Error())
		}
	}

	return nil
}

func (i *LimitInjector) Inject(ctx context.Context) error {
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

	var containerCgroup string
	if i.Info.ContainerRuntime != "" {
		containerCgroup, err = cgroup.GetContainerCgroup(ctx, i.Info.ContainerRuntime, i.Info.ContainerId)
		if err != nil {
			return fmt.Errorf("get cgroup path of container[%s] error: %s", i.Info.ContainerId, err.Error())
		}
	}

	blkioPath := cgroup.GetBlkioCPath(i.Info.Uid, containerCgroup)
	if err := cgroup.NewCgroup(ctx, blkioPath, cgroup.GetBlkioConfig(ctx, devList, i.Args.ReadBytes, i.Args.WriteBytes, i.Args.ReadIO, i.Args.WriteIO, blkioPath)); err != nil {
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

func (i *LimitInjector) Recover(ctx context.Context) error {
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
	isCgroupExist, err := filesys.ExistPathLocal(cgroupPath)
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
			logger.Warnf("fail to get pid[%d]'s old cgroup path, move to \"%s\" instead", pid, TmpCgroup)
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
