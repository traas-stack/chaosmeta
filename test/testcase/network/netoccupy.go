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

package network

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	net2 "github.com/traas-stack/chaosmetad/pkg/utils/net"
	"github.com/traas-stack/chaosmetad/pkg/utils/process"
	"github.com/traas-stack/chaosmetad/test/common"
	"os/exec"
	"strconv"
	"time"
)

var (
	netOccupySleepTime = 1 * time.Second
	netOccupyPort      = 20269
	netOccupyPid       = -1
)

func startTmpService() error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("python -m SimpleHTTPServer %d", netOccupyPort))
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("start server error: %s", err.Error())
	}

	netOccupyPid = cmd.Process.Pid
	time.Sleep(netOccupySleepTime)
	return err
	//listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", netOccupyPort))
	//if err != nil {
	//	panic(any(fmt.Sprintf("Listen() failed, err: %s", err)))
	//}
	//
	//for {
	//	conn, err := listen.Accept()
	//	if err != nil {
	//		continue
	//	}
	//
	//	conn.Close()
	//}
}

func GetNetOccupyTest() []common.TestCase {
	var tempCaseList = []common.TestCase{
		{
			Args:  "",
			Error: true,
		},
		{
			Args:  fmt.Sprintf("-p %d -P abc", netOccupyPort),
			Error: true,
		},
		{
			Args: fmt.Sprintf("-p %d", netOccupyPort),
			Check: func() error {
				return checkPortOccupy(utils.NoPid, false)
			},
			CheckRecover: func() error {
				return checkPortOccupy(utils.NoPid, true)
			},
		},
		{
			Args: fmt.Sprintf("-p %d", netOccupyPort),
			PreProcessor: func() error {
				return startTmpService()
			},
			Error: true,
			PostProcessor: func() error {
				return process.KillProcessByKey(context.Background(), "SimpleHTTPServer", process.SIGKILL)
			},
		},
		{
			Args: fmt.Sprintf("-p %d -f", netOccupyPort),
			PreProcessor: func() error {
				return startTmpService()
			},
			Check: func() error {
				return checkPortOccupy(netOccupyPid, false)
			},
			CheckRecover: func() error {
				return checkPortOccupy(utils.NoPid, true)
			},
		},
	}

	for i := range tempCaseList {
		tempCaseList[i].Target = "network"
		tempCaseList[i].Fault = "occupy"
		tempCaseList[i].Name = fmt.Sprintf("%s-%s-%s", tempCaseList[i].Target, tempCaseList[i].Fault, strconv.Itoa(i))
		if tempCaseList[i].CheckRecover == nil {
			tempCaseList[i].CheckRecover = func() error {
				return nil
			}
		}
	}

	return tempCaseList
}

func checkPortOccupy(targetPid int, eq bool) error {
	pid, err := net2.GetPidByPort(context.Background(), "", "", netOccupyPort, net2.ProtocolTCP)
	if err != nil {
		return err
	}

	if eq {
		if pid != targetPid {
			return fmt.Errorf("not eq: port's pid[%d], target pid[%d]", pid, targetPid)
		}
	} else {
		if pid == targetPid {
			return fmt.Errorf("eq: port's pid == not expected pid: %d", targetPid)
		}
	}

	return nil
}
