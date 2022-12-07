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

package cpu

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cgroup"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	injector.Register(TargetCpu, FaultCpuBurn, func() injector.IInjector { return &BurnInjector{} })
}

type BurnInjector struct {
	injector.BaseInjector
	Args    BurnArgs
	Runtime BurnRuntime
}

type BurnArgs struct {
	Percent int    `json:"percent"`
	Count   int    `json:"count,omitempty"`
	List    string `json:"list,omitempty"`
}

type BurnRuntime struct {
}

func (i *BurnInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *BurnInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *BurnInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "cpu burn usage percent to add, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")
	cmd.Flags().StringVarP(&i.Args.List, "list", "l", "", "cpu burn core number list, start from 0, eg: \"0-2,6\" means \"0,1,2,6\" core")
	cmd.Flags().IntVarP(&i.Args.Count, "count", "c", 0, "cpu burn core count（default 0, means all core）. if provide args \"list\", \"count\" will be ignored.")
}

// Validator list 优先级大于 count
func (i *BurnInjector) Validator() error {
	if i.Args.Percent <= 0 || i.Args.Percent > 100 {
		return fmt.Errorf("\"percent\"[%d] must be in (0,100]", i.Args.Percent)
	}

	cpuList, err := getAllCpuList(i.Info.ContainerRuntime, i.Info.ContainerId)
	if err != nil {
		return fmt.Errorf("get all available cpu list error: %s", err.Error())
	}

	if i.Args.List != "" {
		targetList, err := utils.GetNumArrByList(i.Args.List)
		if err != nil {
			return fmt.Errorf("\"list\"[%s] is not valid: %s", i.Args.List, err.Error())
		}

		for _, core := range targetList {
			var exist bool
			for _, availCore := range cpuList {
				if availCore == core {
					exist = true
					break
				}
			}

			if !exist {
				return fmt.Errorf("core[%d] is not available", core)
			}
		}
	} else {
		if i.Args.Count == 0 || i.Args.Count > len(cpuList) {
			i.Args.Count = len(cpuList)
		}

		if i.Args.Count < 0 {
			return fmt.Errorf("\"count\"[%d] can not less than 0", i.Args.Count)
		}
	}

	if !cmdexec.SupportCmd("taskset") {
		return fmt.Errorf("not support cmd \"taskset\"")
	}

	return i.BaseInjector.Validator()
}

func (i *BurnInjector) Inject() error {
	var coreList []int
	if i.Args.List != "" {
		coreList, _ = utils.GetNumArrByList(i.Args.List)
	} else {
		cpuList, _ := getAllCpuList(i.Info.ContainerRuntime, i.Info.ContainerId)
		coreList = utils.GetNumArrByCount(i.Args.Count, cpuList)
	}

	log.WithUid(i.Info.Uid).Debugf("burn core list: %v", coreList)

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	for c := 0; c < len(coreList); c++ {
		var err error
		cmd := fmt.Sprintf("taskset -c %d %s %s %d %d %d", coreList[c], utils.GetToolPath(CpuBurnKey), i.Info.Uid, coreList[c], i.Args.Percent, timeout)

		if i.Info.ContainerRuntime != "" {
			_, err = cmdexec.ExecContainer(cmd, i.Info.ContainerRuntime, i.Info.ContainerId, namespace.PID)
		} else {
			_, err = cmdexec.StartBashCmdAndWaitPid(cmd)
		}

		if err != nil {
			if err := i.Recover(); err != nil {
				log.WithUid(i.Info.Uid).Warnf("undo error: %s", err.Error())
			}
			return fmt.Errorf("burn cpu of core[%d] error: %s", coreList[c], err.Error())
		}
	}

	return nil
}

func (i *BurnInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	processKey := fmt.Sprintf("%s %s", CpuBurnKey, i.Info.Uid)
	isProExist, err := process.ExistProcessByKey(processKey)
	if err != nil {
		return fmt.Errorf("check process exist by key[%s] error: %s", processKey, err.Error())
	}

	if isProExist {
		if err := process.KillProcessByKey(processKey, process.SIGKILL); err != nil {
			return fmt.Errorf("kill process by key[%s] error: %s", processKey, err.Error())
		}
	}

	return nil
}

func (i *BurnInjector) DelayRecover(timeout int64) error {
	return nil
}

func getAllCpuList(cr, cId string) (cpuList []int, err error) {
	var cpusetPath = "/"
	if cr != "" {
		cpusetPath, err = cgroup.GetContainerCgroupPath(cr, cId, cgroup.CPUSET)
		if err != nil {
			return nil, fmt.Errorf("get cgroup[%s] path of container[%s] error: %s", cgroup.CPUSET, cId, err.Error())
		}
	}

	return getCpuList(cpusetPath)
}

func getCpuList(path string) ([]int, error) {
	cpusetFile := fmt.Sprintf("%s/%s%s/%s", cgroup.RootPath, cgroup.CPUSET, path, cgroup.CpusetCoreFile)
	reByte, err := ioutil.ReadFile(cpusetFile)
	if err != nil {
		return nil, fmt.Errorf("read cpu list info from file[%s] error: %s", cpusetFile, err.Error())
	}

	cpuListStr := string(reByte)
	cpuList, err := utils.GetNumArrByList(cpuListStr)
	if err != nil {
		return nil, fmt.Errorf("format cpu list string error: %s", err.Error())
	}

	return cpuList, nil
}
