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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	disk2 "github.com/ChaosMetaverse/chaosmetad/pkg/utils/disk"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/filesys"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
	"os"
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
}

func (i *FillInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "disk fill target percent, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")
	cmd.Flags().StringVarP(&i.Args.Bytes, "bytes", "b", "", "disk fill bytes to add, support unit: KB/MB/GB/TB（default KB）")
	cmd.Flags().StringVarP(&i.Args.Dir, "dir", "d", "", "disk fill target dir")
}

func (i *FillInjector) Validator() error {
	if i.Args.Dir == "" {
		return fmt.Errorf("\"dir\" is empty")
	}

	if err := filesys.CheckDir(i.Args.Dir); err != nil {
		return fmt.Errorf("\"dir\"[%s] check error: %s", i.Args.Dir, err.Error())
	}

	path, err := filesys.GetAbsPath(i.Args.Dir)
	if err != nil {
		return fmt.Errorf("\"dir\"[%s] get absolute path error: %s", i.Args.Dir, err.Error())
	}

	i.Args.Dir = path

	if i.Args.Percent == 0 && i.Args.Bytes == "" {
		return fmt.Errorf("must provide \"percent\" or \"bytes\"")
	}

	if i.Args.Percent != 0 {
		if i.Args.Percent < 0 || i.Args.Percent > 100 {
			return fmt.Errorf("\"percent\"[%d] must be in (0,100]", i.Args.Percent)
		}
	}

	_, err = getFillKBytes(i.Args.Dir, i.Args.Percent, i.Args.Bytes)
	if err != nil {
		return fmt.Errorf("calculate fill bytes error: %s", err.Error())
	}

	if !cmdexec.SupportCmd("fallocate") && !cmdexec.SupportCmd("dd") {
		return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill disk")
	}

	return i.BaseInjector.Validator()
}

func getFillFileName(uid string) string {
	return fmt.Sprintf("%s%s.dat", FillFileName, uid)
}

func (i *FillInjector) Inject() error {
	fillFile := fmt.Sprintf("%s/%s", i.Args.Dir, getFillFileName(i.Info.Uid))
	bytesKb, _ := getFillKBytes(i.Args.Dir, i.Args.Percent, i.Args.Bytes)

	if err := disk2.RunFillDisk(bytesKb, fillFile); err != nil {
		if err := os.Remove(fillFile); err != nil {
			log.WithUid(i.Info.Uid).Warnf("run failed and delete fill file error: %s", err.Error())
		}
		return err
	}

	return nil
}

func (i *FillInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	fillFile := fmt.Sprintf("%s/%s", i.Args.Dir, getFillFileName(i.Info.Uid))
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

	// fix bug: 如果是数据库文件所在磁盘，填充满后会导致无法数据库，暂时腾出1k空间解决
	if fillKBytes == freeKb {
		fillKBytes -= 10
	}

	// 防止溢出情况
	if fillKBytes <= 0 {
		return -1, fmt.Errorf("fill bytes[%dKB]must larget than 0", fillKBytes)
	}

	return fillKBytes, nil
}
