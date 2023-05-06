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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/disk"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"os"
	"strconv"
)

const (
	FillFileName = "chaosmeta_fill"
)

// [func] [fault] [level] [args]
func main() {
	var (
		err                   error
		fName, _, level, args = os.Args[1], os.Args[2], os.Args[3], os.Args[4:]
		ctx                   = context.Background()
	)
	log.Level = level

	switch fName {
	case utils.MethodValidator:
		err = execValidator(ctx, args)
	case utils.MethodInject:
		err = execInject(ctx, args)
	case utils.MethodRecover:
		err = execRecover(ctx, args)
	default:
		errutil.ExitExpectedErr(fmt.Sprintf("not support method: %s", fName))
	}

	if err != nil {
		errutil.ExitExpectedErr(err.Error())
	}
}

func execValidator(ctx context.Context, args []string) error {
	percentStr, bytes, dir := args[0], args[1], args[2]
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		return fmt.Errorf("percent is not a num")
	}

	return validatorDiskFill(ctx, percent, bytes, dir)
}

func execInject(ctx context.Context, args []string) error {
	percentStr, bytes, dir, uid := args[0], args[1], args[2], args[3]
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		return fmt.Errorf("pecent is not a num")
	}

	return injectDiskFill(ctx, percent, bytes, dir, uid)
}

func execRecover(ctx context.Context, args []string) error {
	return recoverDiskFill(ctx, args[0], args[1])
}

func validatorDiskFill(ctx context.Context, percent int, bytes, dir string) error {
	if percent == 0 && bytes == "" {
		return fmt.Errorf("must provide \"percent\" or \"bytes\"")
	}

	if percent != 0 {
		if percent < 0 || percent > 100 {
			return fmt.Errorf("\"percent\"[%d] must be in (0,100]", percent)
		}
	}

	if dir == "" {
		return fmt.Errorf("\"dir\" is empty")
	}

	if err := filesys.CheckDir(dir); err != nil {
		return fmt.Errorf("\"dir\"[%s] check error: %s", dir, err.Error())
	}

	if _, err := disk.GetFillKBytes(dir, percent, bytes); err != nil {
		return fmt.Errorf("calculate fill bytes error: %s", err.Error())
	}

	if !cmdexec.SupportCmd("fallocate") && !cmdexec.SupportCmd("dd") {
		return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill disk")
	}

	return nil
}

func getFillFileName(uid string) string {
	return fmt.Sprintf("%s%s.dat", FillFileName, uid)
}

func injectDiskFill(ctx context.Context, percent int, bytes, dir, uid string) error {
	logger := log.GetLogger(ctx)
	fillFile := fmt.Sprintf("%s/%s", dir, getFillFileName(uid))
	bytesKb, _ := disk.GetFillKBytes(dir, percent, bytes)

	if err := disk.RunFillDisk(ctx, bytesKb, fillFile); err != nil {
		if err := os.Remove(fillFile); err != nil {
			logger.Warnf("run failed and delete fill file error: %s", err.Error())
		}
		return err
	}

	return nil
}

func recoverDiskFill(ctx context.Context, dir, uid string) error {
	fillFile := fmt.Sprintf("%s/%s", dir, getFillFileName(uid))
	isExist, err := filesys.ExistPath(fillFile)
	if err != nil {
		return fmt.Errorf("check file[%s] exist error: %s", fillFile, err.Error())
	}

	if isExist {
		return os.Remove(fillFile)
	}

	return nil
}
