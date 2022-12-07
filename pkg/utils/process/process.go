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

package process

import (
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/shirou/gopsutil/process"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// signal
const (
	SIGKILL = 9
	SIGTERM = 15
	SIGSTOP = 19
	SIGCONT = 18
)

func ExistPid(pid int) (bool, error) {
	return process.PidExists(int32(pid))
}

func ExistProcessByKey(key string) (bool, error) {
	reByte, err := cmdexec.RunBashCmdWithOutput(fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | wc -l", key, utils.RootName, utils.RootName))
	if err != nil {
		return false, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	if strings.TrimSpace(string(reByte)) == "0" {
		return false, nil
	}

	return true, nil
}

func KillPidWithSignal(pid int, signal int) error {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("find process [%d] error: %s", pid, err.Error())
	}

	switch signal {
	case SIGKILL:
		err = p.Kill()
	case SIGSTOP:
		err = p.Suspend()
	case SIGCONT:
		err = p.Resume()
	case SIGTERM:
		err = p.Terminate()
	default:
		err = p.SendSignal(syscall.Signal(signal))
	}

	if err != nil {
		return fmt.Errorf("kill process [%d] with signal [%d] error: %s", pid, signal, err.Error())
	}

	return nil
}

func KillProcessByKey(key string, signal int) error {
	return cmdexec.RunBashCmdWithoutOutput(fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | awk '{print $2}' | xargs kill -%d", key, utils.RootName, utils.RootName, signal))
}

func GetPidListByKey(key string) ([]int, error) {
	reByte, err := cmdexec.RunBashCmdWithOutput(fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | awk '{print $2}'", key, utils.RootName, utils.RootName))
	if err != nil {
		return nil, fmt.Errorf("get process list error: %s", err.Error())
	}

	var pidList []int
	pStrList := strings.Split(string(reByte), "\n")
	for _, unitPStr := range pStrList {
		if unitPStr == "" {
			continue
		}
		pid, err := strconv.Atoi(unitPStr)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid pid: %s", unitPStr, err.Error())
		}
		pidList = append(pidList, pid)
	}

	return pidList, nil
}

func GetPidListByStr(pidStr string) ([]int, error) {
	var pidList []int

	pidStrList := strings.Split(pidStr, ",")
	for _, unitPid := range pidStrList {
		unitPid = strings.TrimSpace(unitPid)
		if unitPid == "" {
			continue
		}

		pid, err := strconv.Atoi(unitPid)
		if err != nil {
			return nil, fmt.Errorf("pid[%s] is not a int: %s", unitPid, err.Error())
		}

		isExist, err := ExistPid(pid)
		if err != nil {
			return nil, fmt.Errorf("check pid[%d] exist error: %s", pid, err.Error())
		}

		if !isExist {
			return nil, fmt.Errorf("process[%d] is not exist", pid)
		}

		pidList = append(pidList, pid)
	}

	return pidList, nil
}

func GetPidListByListStrAndKey(pidListStr, key string) (pidList []int, err error) {
	if pidListStr != "" {
		pidList, err = GetPidListByStr(pidListStr)
		if err != nil {
			return nil, fmt.Errorf("get pid list by args[%s] error: %s", pidListStr, err.Error())
		}
	} else if key != "" {
		pidList, err = GetPidListByKey(key)
		if err != nil {
			return nil, fmt.Errorf("get pid list by grep key[%s] error: %s", key, err.Error())
		}
	} else {
		return nil, fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if len(pidList) == 0 {
		return nil, fmt.Errorf("no valid pid")
	}

	log.GetLogger().Debugf("pid list: %v", pidList)
	return
}

func GetPidByKeyWithoutRunUser(key string) (int, error) {
	reByte, err := cmdexec.RunBashCmdWithOutput(fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v runuser | awk '{print $2}'", key))
	if err != nil {
		return utils.NoPid, fmt.Errorf("grep process error: %s", err.Error())
	}

	pidStr := strings.TrimSpace(string(reByte))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return utils.NoPid, fmt.Errorf("\"%s\" change to int error: %s", pidStr, err.Error())
	}

	return pid, nil
}

func WaitDefunctProcess() {
	logger := log.GetLogger()

	reByte, err := cmdexec.RunBashCmdWithOutput(fmt.Sprintf("ps -ef | grep '%d' | grep '%s' | grep -v grep | awk '{print $2}'", os.Getpid(), "defunct"))
	if err != nil {
		logger.Warnf("get defunct process error: %s", err.Error())
		return
	}

	output := strings.TrimSpace(string(reByte))
	if output == "" {
		logger.Debugf("no defunct process")
		return
	}

	pidStrList := strings.Split(output, "\n")

	for _, pidStr := range pidStrList {
		if pidStr == "" {
			continue
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			logger.Warnf("\"%s\" is not a num: %s", pidStr, err.Error())
			continue
		}

		var wStatus syscall.WaitStatus
		_, err = syscall.Wait4(pid, &wStatus, 0, nil)
		if err != nil {
			logger.Warnf("wait child process[%d] error: %s", pid, err.Error())
		}
	}
}
