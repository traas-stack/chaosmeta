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
	"github.com/spf13/cobra"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
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

func (i *BurnInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.List == "" && (i.Args.Count == 0 || i.Args.Count > runtime.NumCPU()) {
		i.Args.Count = runtime.NumCPU()
	}
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

	if i.Args.List != "" {
		if _, err := getNumArrByList(i.Args.List); err != nil {
			return fmt.Errorf("\"list\"[%s] is not valid: %s", i.Args.List, err.Error())
		}
	} else {
		if i.Args.Count < 0 {
			return fmt.Errorf("\"count\"[%d] can not less than 0", i.Args.Count)
		}
	}

	if !utils.SupportCmd("taskset") {
		return fmt.Errorf("not support cmd \"taskset\"")
	}

	return i.BaseInjector.Validator()
}

func (i *BurnInjector) Inject() error {
	var coreList []int
	if i.Args.List != "" {
		coreList, _ = getNumArrByList(i.Args.List)
	} else {
		coreList = getNumArrByCount(i.Args.Count)
	}

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	log.WithUid(i.Info.Uid).Debugf("burn core list: %v", coreList)
	if err := burnCpu(i.Info.Uid, i.Args.Percent, coreList, timeout); err != nil {
		if err := i.Recover(); err != nil {
			log.WithUid(i.Info.Uid).Warnf("undo error: %s", err.Error())
		}
		return fmt.Errorf("burn cpu error: %s", err.Error())
	}

	return nil
}

func (i *BurnInjector) Recover() error {
	if i.BaseInjector.Recover() == nil {
		return nil
	}

	processKey := fmt.Sprintf("%s %s", CpuBurnKey, i.Info.Uid)
	isProExist, err := utils.ExistProcessByKey(processKey)
	if err != nil {
		return fmt.Errorf("check process exist by key[%s] error: %s", processKey, err.Error())
	}

	if isProExist {
		if err := utils.KillProcessByKey(processKey, utils.SIGKILL); err != nil {
			return fmt.Errorf("kill process by key[%s] error: %s", processKey, err.Error())
		}
	}

	return nil
}

func (i *BurnInjector) DelayRecover(timeout int64) error {
	return nil
}

func getNumArrByList(listStr string) ([]int, error) {
	maxIndex := runtime.NumCPU() - 1
	var listArr []int
	var ifExist = make(map[int]bool)
	strArr := strings.Split(listStr, ",")
	for _, unitStr := range strArr {
		unitStr = strings.TrimSpace(unitStr)

		if strings.Index(unitStr, "-") >= 0 {
			rangeArr := strings.Split(unitStr, "-")
			if len(rangeArr) != 2 {
				return nil, fmt.Errorf("core range format is error. true format: 1-3")
			}

			rangeArr[0], rangeArr[1] = strings.TrimSpace(rangeArr[0]), strings.TrimSpace(rangeArr[1])
			sCore, err := strconv.Atoi(rangeArr[0])
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", rangeArr[0], err.Error())
			}

			eCore, err := strconv.Atoi(rangeArr[1])
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", rangeArr[1], err.Error())
			}

			if sCore > eCore {
				return nil, fmt.Errorf("core range must: startIndex <= endIndex")
			}

			for i := sCore; i <= eCore; i++ {
				if i < 0 || i > maxIndex {
					return nil, fmt.Errorf("core[%d] is out of core num range: [%d,%d]", i, 0, maxIndex)
				}

				if !ifExist[i] {
					ifExist[i] = true
					listArr = append(listArr, i)
				}
			}
		} else {
			unitCore, err := strconv.Atoi(unitStr)
			if err != nil {
				return nil, fmt.Errorf("core[%s] is not a num: %s", unitStr, err.Error())
			}

			if unitCore < 0 || unitCore > maxIndex {
				return nil, fmt.Errorf("core[%d] is out of core num range: [%d,%d]", unitCore, 0, maxIndex)
			}

			if !ifExist[unitCore] {
				ifExist[unitCore] = true
				listArr = append(listArr, unitCore)
			}
		}
	}

	return listArr, nil
}

func getNumArrByCount(count int) []int {
	total := runtime.NumCPU()
	var listArr = make([]int, total)
	for i := 0; i < total; i++ {
		listArr[i] = i
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(listArr), func(i, j int) {
		listArr[i], listArr[j] = listArr[j], listArr[i]
	})

	return listArr[:count]
}

func burnCpu(uid string, percent int, corelist []int, timeout int64) error {
	for i := 0; i < len(corelist); i++ {
		if _, err := utils.StartBashCmdAndWaitPid(fmt.Sprintf("taskset -c %d %s %s %d %d %d",
			corelist[i], utils.GetToolPath(CpuBurnKey), uid, corelist[i], percent, timeout)); err != nil {
			return err
		}
	}

	return nil
}
