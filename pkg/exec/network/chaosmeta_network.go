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

package main

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmetad/pkg/utils/net"
	"os"
	"strconv"
)

const (
	DirectionOut = "out"

	FaultLimit     = "limit"
	FaultDelay     = "delay"
	FaultLoss      = "loss"
	FaultDuplicate = "duplicate"
	FaultCorrupt   = "corrupt"
	FaultReorder   = "reorder"
)

// [func] [fault] [level] [args]
func main() {
	var (
		err                       error
		fName, fault, level, args = os.Args[1], os.Args[2], os.Args[3], os.Args[4:]
		ctx                       = context.Background()
	)
	log.Level = level

	switch fName {
	case utils.MethodValidator:
		err = execValidator(ctx, fault, args)
	case utils.MethodInject:
		err = execInject(ctx, fault, args)
	case utils.MethodRecover:
		err = execRecover(ctx, args[0])
	default:
		errutil.ExitExpectedErr(fmt.Sprintf("not support method: %s", fName))
	}

	if err != nil {
		errutil.ExitExpectedErr(err.Error())
	}
}

func execValidator(ctx context.Context, fault string, args []string) error {
	force, err := strconv.ParseBool(args[7])
	if err != nil {
		return fmt.Errorf("\"force\" is not a bool: %s", args[1])
	}

	switch fault {
	case FaultDelay:
		return validatorDelay(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, args[8], args[9])
	case FaultLoss:
		percent, err := strconv.Atoi(args[8])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return validatorLoss(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, percent)
	case FaultCorrupt:
		percent, err := strconv.Atoi(args[8])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return validatorCorrupt(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, percent)

	case FaultDuplicate:
		percent, err := strconv.Atoi(args[8])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return validatorDuplicate(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, percent)

	case FaultReorder:
		gap, err := strconv.Atoi(args[9])
		if err != nil {
			return fmt.Errorf("\"gap\" is not a num: %s", args[1])
		}

		return validatorReorder(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, args[8], gap)
	case FaultLimit:
		return validatorLimit(ctx, args[0], args[1], args[2], args[3], args[4], args[5], args[6], force, args[8])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execInject(ctx context.Context, fault string, args []string) error {
	force, err := strconv.ParseBool(args[6])
	if err != nil {
		return fmt.Errorf("\"force\" is not a bool: %s", args[1])
	}

	switch fault {
	case FaultDelay:
		return injectDelay(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, args[7], args[8])
	case FaultLoss:
		percent, err := strconv.Atoi(args[7])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return injectLoss(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, percent)
	case FaultCorrupt:
		percent, err := strconv.Atoi(args[7])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return injectCorrupt(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, percent)
	case FaultDuplicate:
		percent, err := strconv.Atoi(args[7])
		if err != nil {
			return fmt.Errorf("\"percent\" is not a num: %s", args[1])
		}

		return injectDuplicate(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, percent)

	case FaultReorder:
		gap, err := strconv.Atoi(args[8])
		if err != nil {
			return fmt.Errorf("\"gap\" is not a num: %s", args[1])
		}

		return injectReorder(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, args[7], gap)
	case FaultLimit:
		return injectLimit(ctx, args[0], args[1], args[2], args[3], args[4], args[5], force, args[7])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execRecover(ctx context.Context, netInterface string) error {
	isTcExist, err := net.ExistTCRootQdisc(ctx, netInterface)
	if err != nil {
		return fmt.Errorf("check tc rule exist error: %s", err.Error())
	}

	if isTcExist {
		return net.ClearTcRule(ctx, netInterface)
	}

	return nil
}

func validatorDelay(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, latency, jitter string) error {
	if latency == "" {
		return fmt.Errorf("\"latency\" must provide")
	}

	if err := utils.CheckTimeValue(latency); err != nil {
		return fmt.Errorf("\"latency\" is invalid: %s", err.Error())
	}

	if jitter != "" {
		if err := utils.CheckTimeValue(jitter); err != nil {
			return fmt.Errorf("\"jitter\" is invalid: %s", err.Error())
		}
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func validatorLoss(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	if percent <= 0 {
		return fmt.Errorf("\"percent\" must larger than 0")
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func validatorDuplicate(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	if percent <= 0 {
		return fmt.Errorf("\"percent\" must larger than 0")
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func validatorCorrupt(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	if percent <= 0 {
		return fmt.Errorf("\"percent\" must larger than 0")
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func validatorReorder(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, latency string, gap int) error {
	if latency == "" {
		return fmt.Errorf("args latency must provide")
	}

	if err := utils.CheckTimeValue(latency); err != nil {
		return fmt.Errorf("args latency is invalid: %s", err.Error())
	}

	if gap <= 0 {
		return fmt.Errorf("args [gap] must larger than 0")
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func validatorLimit(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool, rate string) error {
	if rate == "" {
		return fmt.Errorf("\"rate\" must provide")
	}

	if err := utils.CheckSpeedValue(rate); err != nil {
		return fmt.Errorf("\"rate\" is invalid: %s", err.Error())
	}

	return validatorTcCommon(ctx, netInterface, direction, mode, sIp, dIp, sPort, dPort, force)
}

func injectDelay(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, latency, jitter string) error {
	return injectTcCommon(ctx, netInterface, mode, sIp, dIp, sPort, dPort, force, FaultDelay, fmt.Sprintf("%s %s", latency, jitter))
}

func injectLoss(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	return injectTcCommon(ctx, netInterface, mode, sIp, dIp, sPort, dPort, force, FaultLoss, fmt.Sprintf("%d", percent))
}

func injectDuplicate(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	return injectTcCommon(ctx, netInterface, mode, sIp, dIp, sPort, dPort, force, FaultDuplicate, fmt.Sprintf("%d", percent))
}

func injectCorrupt(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, percent int) error {
	return injectTcCommon(ctx, netInterface, mode, sIp, dIp, sPort, dPort, force, FaultCorrupt, fmt.Sprintf("%d", percent))
}

func injectReorder(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, latency string, gap int) error {
	return injectTcCommon(ctx, netInterface, mode, sIp, dIp, sPort, dPort, force, FaultReorder, fmt.Sprintf("100 gap %d delay %s", gap, latency))
}

func undoTcWithErr(ctx context.Context, netInterface, msg string) error {
	if err := execRecover(ctx, netInterface); err != nil {
		log.GetLogger(ctx).Warnf("undo tc rule error: %s", err.Error())
	}

	return fmt.Errorf(msg)
}

func injectLimit(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, rate string) error {
	if force {
		exist, _ := net.ExistTCRootQdisc(ctx, netInterface)
		if exist {
			if err := net.ClearTcRule(ctx, netInterface); err != nil {
				return fmt.Errorf("reset tc rule for %s error: %s", netInterface, err.Error())
			}
		}
	}

	if err := net.AddHTBQdisc(ctx, netInterface); err != nil {
		return fmt.Errorf("add htb qdisc for %s error: %s", netInterface, err.Error())
	}

	if err := net.AddLimitClass(ctx, netInterface, rate, mode); err != nil {
		return undoTcWithErr(ctx, netInterface, fmt.Sprintf("add limit class for %s error: %s", netInterface, err.Error()))
	}

	if sIp != "" || dIp != "" || sPort != "" || dPort != "" {
		if err := net.AddFilter(ctx, netInterface, "1:2", sIp, dIp, sPort, dPort); err != nil {
			return undoTcWithErr(ctx, netInterface, fmt.Sprintf("add filter for %s error: %s", netInterface, err.Error()))
		}
	}

	return nil
}

func validatorTcCommon(ctx context.Context, netInterface, direction, mode, sIp, dIp, sPort, dPort string, force bool) error {
	if !cmdexec.SupportCmd("tc") {
		return fmt.Errorf("not support command \"tc\"")
	}

	if netInterface == "" {
		return fmt.Errorf("\"interface\" is empty")
	}

	if !net.ExistInterface(netInterface) {
		return fmt.Errorf("\"interface\"[%s] is not exist", netInterface)
	}

	if direction != DirectionOut {
		return fmt.Errorf("\"direction\" only support: %s", DirectionOut)
	}

	if mode != net.ModeNormal && mode != net.ModeExclude {
		return fmt.Errorf("\"mode\" is not support: %s, only support: %s, %s", mode, net.ModeNormal, net.ModeExclude)
	}

	if sIp != "" {
		if _, err := net.GetValidIPList(sIp, true); err != nil {
			return fmt.Errorf("\"src-ip\"[%s] is invalid: %s", sIp, err.Error())
		}
	}

	if dIp != "" {
		if _, err := net.GetValidIPList(dIp, true); err != nil {
			return fmt.Errorf("\"dst-ip\"[%s] is invalid: %s", dIp, err.Error())
		}
	}

	if sPort != "" {
		if _, err := net.GetValidPortList(sPort); err != nil {
			return fmt.Errorf("\"src-port\"[%s] is invalid: %s", sPort, err.Error())
		}
	}

	if dPort != "" {
		if _, err := net.GetValidPortList(dPort); err != nil {
			return fmt.Errorf("\"dst-port\"[%s] is invalid: %s", dPort, err.Error())
		}
	}

	exist, err := net.ExistTCRootQdisc(ctx, netInterface)
	if err != nil {
		return fmt.Errorf("check tc rule error: %s", err.Error())
	}

	if exist && !force {
		return fmt.Errorf("has other tc root rule, if want to force to execute, please provide [-f] or [--force] args")
	}

	return nil
}

func injectTcCommon(ctx context.Context, netInterface, mode, sIp, dIp, sPort, dPort string, force bool, fault, faultArgs string) error {
	if force {
		exist, _ := net.ExistTCRootQdisc(ctx, netInterface)
		if exist {
			if err := net.ClearTcRule(ctx, netInterface); err != nil {
				return fmt.Errorf("reset tc rule for %s error: %s", netInterface, err.Error())
			}
		}
	}

	if sIp == "" && dIp == "" && sPort == "" && dPort == "" {
		return net.AddNetemQdisc(ctx, netInterface, "", fault, faultArgs)
	}

	if err := net.AddPrioQdisc(ctx, netInterface, "", "1:"); err != nil {
		return fmt.Errorf("add root prio qdisc for %s error: %s", netInterface, err.Error())
	}

	if mode == net.ModeNormal {
		parent := "1:4"
		if err := net.AddNetemQdisc(ctx, netInterface, parent, fault, faultArgs); err != nil {
			return undoTcWithErr(ctx, netInterface, fmt.Sprintf("add parent %s netem qdisc for %s error: %s", parent, netInterface, err.Error()))
		}
	} else {
		for subIndex := 1; subIndex < 4; subIndex++ {
			parent := fmt.Sprintf("1:%d", subIndex)
			if err := net.AddNetemQdisc(ctx, netInterface, parent, fault, faultArgs); err != nil {
				return undoTcWithErr(ctx, netInterface, fmt.Sprintf("add parent %s netem qdisc for %s error: %s", parent, netInterface, err.Error()))
			}
		}
	}

	if err := net.AddFilter(ctx, netInterface, "1:4", sIp, dIp, sPort, dPort); err != nil {
		return undoTcWithErr(ctx, netInterface, fmt.Sprintf("add filter for %s error: %s", netInterface, err.Error()))
	}

	return nil
}
