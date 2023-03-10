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
	"github.com/traas-stack/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmetad/pkg/utils/process"
	"os"
)

const (
	DiskIOBurnKey   = "chaosmeta_diskburn"
	FaultDiskIOBurn = "burn"
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
	case FaultDiskIOBurn:
		return validatorBurn(ctx, args[0])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func validatorBurn(ctx context.Context, dir string) error {
	if err := filesys.CheckDir(dir); err != nil {
		return fmt.Errorf("\"dir\"[%s] check error: %s", dir, err.Error())
	}

	if !cmdexec.SupportCmd("dd") {
		return fmt.Errorf("not support cmd \"dd\"")
	}

	return nil
}

func execInject(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultDiskIOBurn:
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func getBurnFileName(uid, dir string) string {
	return fmt.Sprintf("%s/%s_%s", dir, DiskIOBurnKey, uid)
}

func execRecover(ctx context.Context, fault string, args []string) error {
	switch fault {
	case FaultDiskIOBurn:
		return recoverBurn(ctx, args[0], args[1])
	default:
		return fmt.Errorf("not support fault: %s", fault)
	}
}

func recoverBurn(ctx context.Context, uid, dir string) error {
	if err := process.CheckExistAndKillByKey(ctx, fmt.Sprintf("%s %s", DiskIOBurnKey, uid)); err != nil {
		return err
	}

	file := getBurnFileName(uid, dir)
	isFileExist, err := filesys.ExistPath(file)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", file, err.Error())
	}

	if isFileExist {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("remove file[%s] error: %s", file, err.Error())
		}
	}

	return nil
}
