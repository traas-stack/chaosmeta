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
	"context"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/crclient"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
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

func getProcessPidCmd(pid int) string {
	return fmt.Sprintf("ps -eo pid | grep -w %d", pid)
}

func getProcessKeyCmd(key string) string {
	return fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | grep -v 'chaosmeta_process ' | awk '{print $2}'", key, utils.RootName, utils.RootName)
}

func getProcessSignalPidCmd(pid int, signal int) string {
	return fmt.Sprintf("kill -%d %d", signal, pid)
}

func getProcessSignalKeyCmd(key string, signal int) string {
	return fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | grep -v 'chaosmeta_process ' | awk '{print $2}' | xargs kill -%d", key, utils.RootName, utils.RootName, signal)
}

// GetProcessByPid in container's pid namespace
func GetProcessByPid(ctx context.Context, cr, cId string, pid int) (int, error) {
	if pid < 0 {
		return -1, fmt.Errorf("\"pid\" can not less than 0")
	}

	var (
		reStr string
		cmd   = getProcessPidCmd(pid)
		err   error
	)

	reStr, err = cmdexec.ExecCommon(ctx, cr, cId, cmd)
	if err != nil {
		return -1, fmt.Errorf("exec cmd error: %s", err.Error())
	}

	reStr = strings.TrimSpace(reStr)
	if reStr != strconv.Itoa(pid) {
		return -1, fmt.Errorf("pid[%d] is not exist, get: %s", pid, reStr)
	}

	return pid, nil
}

// GetProcessByKey in container's pid namespace
func GetProcessByKey(ctx context.Context, cr, cId string, key string) ([]int, error) {
	if key == "" {
		return nil, fmt.Errorf("\"key\" can not be empty")
	}

	var (
		reStr   string
		cmd     = getProcessKeyCmd(key)
		err     error
		pidList []int
	)

	reStr, err = cmdexec.ExecCommon(ctx, cr, cId, cmd)
	if err != nil {
		return nil, fmt.Errorf("exec cmd error: %s", err.Error())
	}

	reStr = strings.TrimSpace(reStr)
	pStrList := strings.Split(reStr, "\n")
	for _, unitPStr := range pStrList {
		unitPStr = strings.TrimSpace(unitPStr)
		if unitPStr == "" {
			continue
		}
		pid, err := strconv.Atoi(unitPStr)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid pid: %s", unitPStr, err.Error())
		}
		pidList = append(pidList, pid)
	}

	if len(pidList) == 0 {
		return nil, fmt.Errorf("no pid grep by key: %s", key)
	}

	return pidList, nil
}

// SignalProcessByPid in container's pid namespace
func SignalProcessByPid(ctx context.Context, cr, cId string, pid, signal int) error {
	if pid < 0 {
		return fmt.Errorf("\"pid\" can not less than 0")
	}

	var (
		cmd = getProcessSignalPidCmd(pid, signal)
		err error
	)

	_, err = cmdexec.ExecCommon(ctx, cr, cId, cmd)
	if err != nil {
		return fmt.Errorf("exec cmd error: %s", err.Error())
	}

	return nil
}

// SignalProcessByKey in container's pid namespace
func SignalProcessByKey(ctx context.Context, cr, cId string, key string, signal int) error {
	if key == "" {
		return fmt.Errorf("\"key\" can not be empty")
	}

	var (
		cmd = getProcessSignalKeyCmd(key, signal)
		err error
	)

	_, err = cmdexec.ExecCommon(ctx, cr, cId, cmd)
	if err != nil {
		return fmt.Errorf("exec cmd error: %s", err.Error())
	}

	return nil
}

func ExistPid(ctx context.Context, pid int) (bool, error) {
	return process.PidExists(int32(pid))
}

func CheckExistAndKillByKey(ctx context.Context, processKey string) error {
	isProExist, err := ExistProcessByKey(ctx, processKey)
	if err != nil {
		return fmt.Errorf("check process exist by key[%s] error: %s", processKey, err.Error())
	}

	if isProExist {
		if err := KillProcessByKey(ctx, processKey, SIGKILL); err != nil {
			return fmt.Errorf("kill process by key[%s] error: %s", processKey, err.Error())
		}
	}

	return nil
}

func ExistProcessByKey(ctx context.Context, key string) (bool, error) {
	re, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | grep -v 'chaosmeta_process ' | wc -l", key, utils.RootName, utils.RootName))
	if err != nil {
		return false, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	if strings.TrimSpace(re) == "0" {
		return false, nil
	}

	return true, nil
}

func KillPidWithSignal(ctx context.Context, pid int, signal int) error {
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

func KillProcessByKey(ctx context.Context, key string, signal int) error {
	return cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v '%s inject' | grep -v '%s recover' | grep -v 'chaosmeta_process ' | awk '{print $2}' | xargs kill -%d", key, utils.RootName, utils.RootName, signal))
}

// GetPidListByPidOrKeyInContainer return pidList in host's or container's pid ns
func GetPidListByPidOrKeyInContainer(ctx context.Context, cr, cId string, pid int, key string) ([]int, error) {
	//logger := log.GetLogger(ctx)
	var pidList []int
	if pid > 0 {
		if cr != "" {
			cmd := fmt.Sprintf("ps -eo pid | grep -w %d", pid)
			reStr, err := cmdexec.ExecContainerRaw(ctx, cr, cId, cmd)
			if err != nil {
				return nil, fmt.Errorf("exec cmd[%s] in container[%s] error: %s", cmd, cId, err.Error())
			}
			reStr = strings.TrimSpace(reStr)
			if reStr != strconv.Itoa(pid) {
				return nil, fmt.Errorf("pid[%d] is not exist", pid)
			}
		} else {
			ifExist, err := ExistPid(ctx, pid)
			if err != nil {
				return nil, fmt.Errorf("check pid exist error: %s", err.Error())
			}
			if !ifExist {
				return nil, fmt.Errorf("pid[%d] is not exist", pid)
			}
		}
		pidList = append(pidList, pid)
	} else if key != "" {
		cmd := getProcessKeyCmd(key)
		var (
			reStr string
			err   error
		)

		if cr != "" {
			reStr, err = cmdexec.ExecContainerRaw(ctx, cr, cId, cmd)
			if err != nil {
				return nil, fmt.Errorf("exec cmd[%s] in container[%s] error: %s", cmd, cId, err.Error())
			}
		} else {
			reStr, err = cmdexec.RunBashCmdWithOutput(ctx, getProcessKeyCmd(key))
			if err != nil {
				return nil, fmt.Errorf("get process list error: %s", err.Error())
			}
		}

		pStrList := strings.Split(reStr, "\n")
		for _, unitPStr := range pStrList {
			unitPStr = strings.TrimSpace(unitPStr)
			if unitPStr == "" {
				continue
			}
			unitPid, err := strconv.Atoi(unitPStr)
			if err != nil {
				return nil, fmt.Errorf("%s is not a valid pid: %s", unitPStr, err.Error())
			}
			pidList = append(pidList, unitPid)
		}

	} else {
		return nil, fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if len(pidList) == 0 {
		return nil, fmt.Errorf("target pid is empty")
	}

	return pidList, nil
}

// GetPidListByKey return pidList in host's pid ns, not in container's pid ns
func GetPidListByKey(ctx context.Context, cr, cId string, key string) ([]int, error) {
	var pidList []int
	if cr == "" {
		re, err := cmdexec.RunBashCmdWithOutput(ctx, getProcessKeyCmd(key))
		if err != nil {
			return nil, fmt.Errorf("get process list error: %s", err.Error())
		}

		pStrList := strings.Split(re, "\n")
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
	} else {
		client, err := crclient.GetClient(ctx, cr)
		if err != nil {
			return nil, fmt.Errorf("get %s client error: %s", cr, err.Error())
		}

		existPro, err := client.GetAllPidList(ctx, cId)
		if err != nil {
			return nil, fmt.Errorf("get pid of %s error: %s", cId, err.Error())
		}

		for _, unit := range existPro {
			if strings.Index(unit.Cmd, key) >= 0 {
				pidList = append(pidList, unit.Pid)
			}
		}
	}

	return pidList, nil
}

func GetPidListByStr(ctx context.Context, pidStr string) ([]int, error) {
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

		isExist, err := ExistPid(ctx, pid)
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

func GetPidListByListStrAndKey(ctx context.Context, cr, cId string, pidListStr, key string) (pidList []int, err error) {
	if pidListStr != "" {
		if cr != "" {
			return nil, fmt.Errorf("not support \"pid-list\" args in container")
		}

		pidList, err = GetPidListByStr(ctx, pidListStr)
		if err != nil {
			return nil, fmt.Errorf("get pid list by args[%s] error: %s", pidListStr, err.Error())
		}
	} else if key != "" {
		pidList, err = GetPidListByKey(ctx, cr, cId, key)
		if err != nil {
			return nil, fmt.Errorf("get pid list by grep key[%s] error: %s", key, err.Error())
		}
	} else {
		return nil, fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if len(pidList) == 0 {
		return nil, fmt.Errorf("no valid pid")
	}

	log.GetLogger(ctx).Debugf("pid list: %v", pidList)
	return
}

func GetPidByKeyWithoutRunUser(ctx context.Context, key string) (int, error) {
	re, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("ps -ef | grep '%s' | grep -v grep | grep -v runuser | awk '{print $2}'", key))
	if err != nil {
		return utils.NoPid, fmt.Errorf("grep process error: %s", err.Error())
	}

	pidStr := strings.TrimSpace(re)
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return utils.NoPid, fmt.Errorf("\"%s\" change to int error: %s", pidStr, err.Error())
	}

	return pid, nil
}

func WaitDefunctProcess(ctx context.Context) {
	logger := log.GetLogger(ctx)

	re, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("ps -ef | grep '%d' | grep '%s' | grep -v grep | awk '{print $2}'", os.Getpid(), "defunct"))
	if err != nil {
		logger.Warnf("get defunct process error: %s", err.Error())
		return
	}

	output := strings.TrimSpace(re)
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
