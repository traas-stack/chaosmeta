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

package disk

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	disk2 "github.com/ChaosMetaverse/chaosmetad/pkg/utils/disk"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	injector.Register(TargetDisk, FaultDiskFill, func() injector.IInjector { return &FillInjector{} })
}

type FillInjector struct {
	injector.BaseInjector
	Args    FillArgs
	Runtime FillRuntime
}

type FillArgs struct {
	Percent int    `json:"percent,omitempty"`
	Bytes   string `json:"bytes,omitempty"`
	Dir     string `json:"dir,omitempty"`
}

type FillRuntime struct {
}

func (i *FillInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *FillInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *FillInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Dir == "" {
		i.Args.Dir = DefaultDir
	}

	if i.Info.ContainerRuntime != "" {
		return
	}

	// TODO: SetDefault还是需要加上error返回参数才行，需要改参数的都放到SetDefault来，容器内的怎么获取
	var err error
	i.Args.Dir, err = filesys.GetAbsPath(i.Args.Dir)
	if err != nil {
		panic(any(fmt.Sprintf("\"dir\"[%s] get absolute path error: %s", i.Args.Dir, err.Error())))
	}
}

func (i *FillInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "disk fill target percent, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")
	cmd.Flags().StringVarP(&i.Args.Bytes, "bytes", "b", "", "disk fill bytes to add, support unit: KB/MB/GB/TB（default KB）")
	cmd.Flags().StringVarP(&i.Args.Dir, "dir", "d", "", "disk fill target dir")
}

func (i *FillInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Info.ContainerRuntime != "" {
		if !filepath.IsAbs(i.Args.Dir) {
			return fmt.Errorf("\"dir\" must provide absolute path")
		}

		//client, _ := crclient.GetClient(ctx, i.Info.ContainerRuntime)
		//if err := client.CpFile(ctx, i.Info.ContainerId, utils.GetToolPath(DiskFillKey), fmt.Sprintf("/tmp/%s", DiskFillKey)); err != nil {
		//	return fmt.Errorf("cp exec tool to container[%s] error: %s", i.Info.ContainerId, err.Error())
		//}

		if err := filesys.CpContainerFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, utils.GetToolPath(DiskFillKey),
			fmt.Sprintf("/tmp/%s", DiskFillKey)); err != nil {
			return fmt.Errorf("cp exec tool to container[%s] error: %s", i.Info.ContainerId, err.Error())
		}

		_, err := cmdexec.ExecContainer(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, []string{namespace.MNT},
			fmt.Sprintf("/tmp/%s %s %d %s %s", DiskFillKey, "validator", i.Args.Percent, i.Args.Bytes, i.Args.Dir), true)
		if err != nil {
			return fmt.Errorf("exec in container error: %s", err.Error())
		}

		return nil
	} else {
		return ValidatorDiskFill(ctx, i.Args.Percent, i.Args.Bytes, i.Args.Dir)
	}
}

func ValidatorDiskFill(ctx context.Context, percent int, bytes, dir string) error {
	if percent == 0 && bytes == "" {
		return fmt.Errorf("must provide \"percent\" or \"bytes\"")
	}

	if percent != 0 {
		if percent < 0 || percent > 100 {
			return fmt.Errorf("\"percent\"[%d] must be in (0,100]", percent)
		}
	}

	if dir == "" {
		return fmt.Errorf("\"dir\" is empty")
	}

	if err := filesys.CheckDir(dir); err != nil {
		return fmt.Errorf("\"dir\"[%s] check error: %s", dir, err.Error())
	}

	if _, err := getFillKBytes(dir, percent, bytes); err != nil {
		return fmt.Errorf("calculate fill bytes error: %s", err.Error())
	}

	if !cmdexec.SupportCmd("fallocate") && !cmdexec.SupportCmd("dd") {
		return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill disk")
	}

	return nil
}

func getFillFileName(uid string) string {
	return fmt.Sprintf("%s%s.dat", FillFileName, uid)
}

func (i *FillInjector) Inject(ctx context.Context) error {
	if i.Info.ContainerRuntime != "" {
		if err := filesys.CpContainerFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, utils.GetToolPath(DiskFillKey),
			fmt.Sprintf("/tmp/%s", DiskFillKey)); err != nil {
			return fmt.Errorf("cp exec tool to container[%s] error: %s", i.Info.ContainerId, err.Error())
		}

		_, err := cmdexec.ExecContainer(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, []string{namespace.MNT},
			fmt.Sprintf("/tmp/%s %s %d %s %s %s", DiskFillKey, "inject", i.Args.Percent, i.Args.Bytes, i.Args.Dir, i.Info.Uid), true)
		if err != nil {
			return fmt.Errorf("exec in container error: %s", err.Error())
		}

		return nil
	} else {
		return InjectDiskFill(ctx, i.Args.Percent, i.Args.Bytes, i.Args.Dir, i.Info.Uid)
	}
}

func InjectDiskFill(ctx context.Context, percent int, bytes, dir, uid string) error {
	logger := log.GetLogger(ctx)
	fillFile := fmt.Sprintf("%s/%s", dir, getFillFileName(uid))
	bytesKb, _ := getFillKBytes(dir, percent, bytes)

	if err := disk2.RunFillDisk(ctx, bytesKb, fillFile); err != nil {
		if err := os.Remove(fillFile); err != nil {
			logger.Warnf("run failed and delete fill file error: %s", err.Error())
		}
		return err
	}

	return nil
}

func (i *FillInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	if i.Info.ContainerRuntime != "" {
		if err := filesys.CpContainerFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, utils.GetToolPath(DiskFillKey),
			fmt.Sprintf("/tmp/%s", DiskFillKey)); err != nil {
			return fmt.Errorf("cp exec tool to container[%s] error: %s", i.Info.ContainerId, err.Error())
		}

		_, err := cmdexec.ExecContainer(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, []string{namespace.MNT},
			fmt.Sprintf("/tmp/%s %s %s %s", DiskFillKey, "recover", i.Args.Dir, i.Info.Uid), true)
		if err != nil {
			return fmt.Errorf("exec in container error: %s", err.Error())
		}

		return nil
	} else {
		return RecoverDiskFill(ctx, i.Args.Dir, i.Info.Uid)
	}
}

func RecoverDiskFill(ctx context.Context, dir, uid string) error {
	fillFile := fmt.Sprintf("%s/%s", dir, getFillFileName(uid))
	isExist, err := filesys.ExistPath(fillFile)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", fillFile, err.Error())
	}

	if isExist {
		return os.Remove(fillFile)
	}

	return nil
}

func getFillKBytes(dir string, percent int, bytes string) (int64, error) {
	var fillKBytes int64
	usage, err := disk.Usage(dir)
	if err != nil {
		return -1, fmt.Errorf("get disk info error: %s", err.Error())
	}

	if percent != 0 {
		if float64(percent) < usage.UsedPercent {
			return -1, fmt.Errorf("target path current disk usage is %.2f%%, no need to fill", usage.UsedPercent)
		}

		fillKBytes = int64(uint64((float64(percent) - usage.UsedPercent) / 100 * (float64(usage.Total) / 1024)))
	} else {
		fillKBytes, err = utils.GetKBytes(bytes)
		if err != nil {
			return -1, fmt.Errorf("\"bytes\" is invalid: %s", err.Error())
		}
	}

	freeKb := int64(usage.Free / 1024)
	if fillKBytes > freeKb {
		return -1, fmt.Errorf("space not enough, fill: %dKB, free: %dKB", fillKBytes, freeKb)
	}

	// fix bug: If it is the disk where the database file is located, the database cannot be read when it is full.
	// so temporarily free up 10k space to solve the problem
	if fillKBytes == freeKb {
		fillKBytes -= 10
	}

	// prevent overflow
	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}
