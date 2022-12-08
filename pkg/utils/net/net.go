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

package net

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"net"
	"strconv"
	"strings"
)

const (
	ModeNormal   = "normal"
	ModeExclude  = "exclude"
	PortBit      = 16
	MaxRuleCount = 625 // 已测过864不行

	ProtocolTCP  = "tcp"
	ProtocolTCP6 = "tcp6"
	ProtocolUDP  = "udp"
	ProtocolUDP6 = "udp6"
)

func ExistInterface(name string) bool {
	if _, err := net.InterfaceByName(name); err != nil {
		return false
	}

	return true
}

// GetValidIPList only support ipv4
func GetValidIPList(ipStr string, ifSubNet bool) ([]string, error) {
	ipStrList := strings.Split(ipStr, ",")
	var re = make([]string, len(ipStrList))
	for i, unit := range ipStrList {
		unit = strings.TrimSpace(unit)
		if address := net.ParseIP(unit); address != nil {
			re[i] = unit
			continue
		}

		if !ifSubNet {
			return nil, fmt.Errorf("%s is not a valid ip", unit)
		}

		if _, _, err := net.ParseCIDR(unit); err == nil {
			re[i] = unit
			continue
		}

		return nil, fmt.Errorf("%s is not a valid ip or subNet", unit)
	}

	return re, nil
}

func GetValidPortList(portStr string) ([]string, error) {
	portStrList := strings.Split(portStr, ",")
	var re = make([]string, len(portStrList))
	for i, unit := range portStrList {
		unit = strings.TrimSpace(unit)
		portArr := strings.Split(unit, "/")
		if len(portArr) < 2 {
			portArr = append(portArr, fmt.Sprintf("%d", PortBit))
		}

		port, err := strconv.Atoi(portArr[0])
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid port", portArr[0])
		}

		if port <= 0 || port >= 65535 {
			return nil, fmt.Errorf("%d is invalid, should in (0, 65536)", port)
		}

		mask, err := strconv.Atoi(portArr[1])
		if err != nil {
			return nil, fmt.Errorf("%s port mask is not a num", portArr[1])
		}

		if mask < 0 || mask > PortBit {
			return nil, fmt.Errorf("%d port mask should in [0,%d]", mask, PortBit)
		}

		re[i] = fmt.Sprintf("%s%s%s", portArr[0], utils.PortSplit, getPortMask(mask))
	}

	return re, nil
}

func getPortMask(mask int) string {
	needZero := PortBit - mask
	var maskValue int
	for i := 0; i < mask; i++ {
		maskValue = maskValue<<1 + 1
	}

	for i := 0; i < needZero; i++ {
		maskValue <<= 1
	}

	return fmt.Sprintf("0x%x", maskValue)
}

func ClearTcRule(ctx context.Context, netInterface string) error {
	return cmdexec.RunBashCmdWithoutOutput(ctx, GetClearTcRuleCmd(netInterface))
}

func GetClearTcRuleCmd(netInterface string) string {
	return fmt.Sprintf("tc qdisc del dev %s root", netInterface)
}

func AddNetemQdisc(ctx context.Context, netInterface, parent, fault string, args string) error {
	if parent == "" {
		parent = "root handle 1:"
	} else {
		parent = fmt.Sprintf("parent %s", parent)
	}

	cmd := fmt.Sprintf("tc qdisc add dev %s %s netem %s %s", netInterface, parent, fault, args)
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, cmd); err != nil {
		return err
	}

	return nil
}

func AddPrioQdisc(ctx context.Context, netInterface, parent, name string) error {
	if parent == "" {
		parent = "root"
	} else {
		parent = fmt.Sprintf("parent %s", parent)
	}

	cmd := fmt.Sprintf("tc qdisc add dev %s %s handle %s prio bands 4", netInterface, parent, name)
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, cmd); err != nil {
		return err
	}

	return nil
}

// AddHTBQdisc default 1:1
func AddHTBQdisc(ctx context.Context, netInterface string) error {
	cmdStr := fmt.Sprintf("tc qdisc add dev %s root handle 1: htb default 1", netInterface)
	log.GetLogger(ctx).Debugf("add htb qdisc cmd: %s", cmdStr)

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, cmdStr); err != nil {
		return err
	}

	return nil
}

func AddLimitClass(ctx context.Context, netInterface string, rate string, mode string) error {
	subNum := 1
	if mode == ModeNormal {
		subNum = 2
	}

	cmdStr := fmt.Sprintf("tc class add dev %s parent 1: classid 1:%d htb rate %s", netInterface, subNum, rate)
	log.GetLogger(ctx).Debugf("add class cmd: %s", cmdStr)

	if err := cmdexec.RunBashCmdWithoutOutput(ctx, cmdStr); err != nil {
		return err
	}

	return nil
}

func AddFilter(ctx context.Context, netInterface, target, srcIpListStr, dstIpListStr, srcPortListStr, dstPortListStr string) error {
	cmd, err := getAddFilterCmd(ctx, netInterface, target, srcIpListStr, dstIpListStr, srcPortListStr, dstPortListStr)
	if err != nil {
		return fmt.Errorf("get filter cmd error: %s", err.Error())
	}

	log.GetLogger(ctx).Debugf("add filter cmd: %s", cmd)
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, cmd); err != nil {
		return fmt.Errorf("run cmd error: %s", err.Error())
	}

	return nil
}

func getAddFilterCmd(ctx context.Context, netInterface, target, srcIpListStr, dstIpListStr, srcPortListStr, dstPortListStr string) (tcFilterStr string, err error) {
	srcIpList, dstIpList, srcPortList, dstPortList, err := getStrList(srcIpListStr, dstIpListStr, srcPortListStr, dstPortListStr)
	if err != nil {
		return
	}

	var si, di, sp, dp int
	var ruleArr []string
	var siLen, diLen, spLen, dpLen = len(srcIpList), len(dstIpList), len(srcPortList), len(dstPortList)
	for {
		if si >= siLen && di >= diLen && sp >= spLen && dp >= dpLen {
			break
		}

		var ifCtn = true
		if (siLen > 0 && si == siLen) || (diLen > 0 && di == diLen) || (spLen > 0 && sp == spLen) || (dpLen > 0 && dp == dpLen) {
			ifCtn = false
		}

		if ifCtn {
			var args string
			if siLen > 0 {
				args += fmt.Sprintf("match ip src %s ", srcIpList[si])
			}

			if diLen > 0 {
				args += fmt.Sprintf("match ip dst %s ", dstIpList[di])
			}

			if spLen > 0 {
				portArr := strings.Split(srcPortList[sp], utils.PortSplit)
				args += fmt.Sprintf("match ip sport %s %s ", portArr[0], portArr[1])
			}

			if dpLen > 0 {
				portArr := strings.Split(dstPortList[dp], utils.PortSplit)
				args += fmt.Sprintf("match ip dport %s %s ", portArr[0], portArr[1])
			}

			if args != "" {
				ruleArr = append(ruleArr, fmt.Sprintf("tc filter add dev %s parent 1: prio 1 protocol ip u32 %sflowid %s", netInterface, args, target))
				if len(ruleArr) > MaxRuleCount {
					err = fmt.Errorf("filter rule count is larget than %d", MaxRuleCount)
					return
				}
			}
		}

		if si < len(srcIpList) {
			si++
		} else if di < len(dstIpList) {
			di++
			si = 0
		} else if sp < len(srcPortList) {
			sp++
			si, di = 0, 0
		} else if dp < len(dstPortList) {
			dp++
			si, di, sp = 0, 0, 0
		}
	}

	log.GetLogger(ctx).Debugf("filter rule count: %d", len(ruleArr))
	tcFilterStr = strings.Join(ruleArr, utils.CmdSplit)
	return
}

func getStrList(srcIpListStr, dstIpListStr, srcPortListStr, dstPortListStr string) (srcIpList, dstIpList, srcPortList, dstPortList []string, err error) {
	if srcIpListStr != "" {
		srcIpList, err = GetValidIPList(srcIpListStr, true)
		if err != nil {
			err = fmt.Errorf("get valid src ip list from [%s] error: %s", srcIpListStr, err.Error())
			return
		}
	}

	if dstIpListStr != "" {
		dstIpList, err = GetValidIPList(dstIpListStr, true)
		if err != nil {
			err = fmt.Errorf("get valid dst ip list from [%s] error: %s", dstIpListStr, err.Error())
			return
		}
	}

	if srcPortListStr != "" {
		srcPortList, err = GetValidPortList(srcPortListStr)
		if err != nil {
			err = fmt.Errorf("get valid src port list from [%s] error: %s", srcPortListStr, err.Error())
			return
		}
	}

	if dstPortListStr != "" {
		dstPortList, err = GetValidPortList(dstPortListStr)
		if err != nil {
			err = fmt.Errorf("get valid dst port list from [%s] error: %s", dstPortListStr, err.Error())
			return
		}
	}

	return
}

func ExistTCRootQdisc(ctx context.Context, netInterface string) (bool, error) {
	out, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("tc qdisc ls dev %s | grep -w \"1: root\" | grep -v grep | wc -l", netInterface))
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(string(out)) == "0" {
		return false, nil
	}

	return true, nil
}

func GetPidByPort(ctx context.Context, port int, proto string) (int, error) {
	var cmd string
	if proto == ProtocolTCP || proto == ProtocolTCP6 {
		cmd = fmt.Sprintf("netstat -anpt | grep -w %s | awk '{print $4,$7}' | grep -w %d | grep :%d | awk '{print $2}' | awk -F'/' '{print $1}'", proto, port, port)
	} else if proto == ProtocolUDP || proto == ProtocolUDP6 {
		cmd = fmt.Sprintf("netstat -anpu | grep -w %s | awk '{print $4,$6}' | grep -w %d | grep :%d | awk '{print $2}' | awk -F'/' '{print $1}'", proto, port, port)
	} else {
		return utils.NoPid, fmt.Errorf("protocol not support: %s、%s、%s、%s", ProtocolTCP, ProtocolUDP, ProtocolTCP6, ProtocolUDP6)
	}

	log.GetLogger(ctx).Debugf("get pid by port cmd: %s", cmd)
	out, err := cmdexec.RunBashCmdWithOutput(ctx, cmd)
	if err != nil {
		return utils.NoPid, fmt.Errorf("cmd exec error: %s", err.Error())
	}

	pidStr := strings.TrimSpace(string(out))
	if pidStr == "" {
		return utils.NoPid, nil
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return utils.NoPid, fmt.Errorf("pid[%s] to int error: %s", pidStr, err.Error())
	}

	return pid, nil
}
