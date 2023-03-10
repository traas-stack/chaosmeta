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
	"github.com/traas-stack/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmetad/pkg/utils/process"
	"os"
	"strconv"
)

const (
	FaultProcessKill = "kill"
	FaultProcessStop = "stop"
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
		err = execRecover(ctx, fault, args)
	default:
		errutil.ExitExpectedErr(fmt.Sprintf("not support method: %s", fName))
	}

	if err != nil {
		errutil.ExitExpectedErr(err.Error())
	}
}

func execValidator(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultProcessKill:
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("\"pid\" is not a num: %s", args[1])
		}
		signal, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("\"signal\" is not a num: %s", args[1])
		}

		return validatorKill(ctx, pid, args[1], signal)
	case FaultProcessStop:
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("\"pid\" is not a num: %s", args[1])
		}

		return validatorStop(ctx, pid, args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execInject(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultProcessKill:
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("\"pid\" is not a num: %s", args[1])
		}
		signal, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("\"signal\" is not a num: %s", args[1])
		}

		return injectKill(ctx, pid, args[1], signal)
	case FaultProcessStop:
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("\"pid\" is not a num: %s", args[1])
		}

		return injectStop(ctx, pid, args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func execRecover(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultProcessKill:
		return nil
	case FaultProcessStop:
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("\"pid\" is not a num: %s", args[1])
		}

		return recoverStop(ctx, pid, args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func validatorKill(ctx context.Context, pid int, key string, signal int) error {
	if signal <= 0 {
		return fmt.Errorf("signal[%d] is invalid, must larget than 0", signal)
	}

	return validatorCommon(ctx, pid, key)
}

func validatorStop(ctx context.Context, pid int, key string) error {
	return validatorCommon(ctx, pid, key)
}

func injectKill(ctx context.Context, pid int, key string, signal int) error {
	return injectCommon(ctx, pid, key, signal)
}

func injectStop(ctx context.Context, pid int, key string) error {
	return injectCommon(ctx, pid, key, process.SIGSTOP)
}

func recoverStop(ctx context.Context, pid int, key string) error {
	return injectCommon(ctx, pid, key, process.SIGCONT)
}

func validatorCommon(ctx context.Context, pid int, key string) error {
	if pid < 0 {
		return fmt.Errorf("\"pid\" can not less than 0")
	}

	if pid == 0 && key == "" {
		return fmt.Errorf("must provide \"pid\" or \"key\"")
	}

	if pid > 0 {
		exist, err := process.ExistPid(ctx, pid)
		if err != nil {
			return fmt.Errorf("check pid[%d] exist error: %s", pid, err.Error())
		}

		if !exist {
			return fmt.Errorf("pid[%d] not exist", pid)
		}
	} else {
		exist, err := process.ExistProcessByKey(ctx, key)
		if err != nil {
			return fmt.Errorf("check pid by key[%s] error: %s", key, err.Error())
		}

		if !exist {
			return fmt.Errorf("no process grep by key[%s]", key)
		}
	}

	return nil
}

func injectCommon(ctx context.Context, pid int, key string, signal int) error {
	if pid > 0 {
		if err := process.KillPidWithSignal(ctx, pid, signal); err != nil {
			return err
		}
	} else {
		if err := process.KillProcessByKey(ctx, key, signal); err != nil {
			return err
		}
	}

	return nil
}
