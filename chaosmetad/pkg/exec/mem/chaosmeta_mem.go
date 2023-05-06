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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/filesys"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/memory"
	"os"
	"strconv"
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
	if !cmdexec.SupportCmd("fallocate") && !cmdexec.SupportCmd("dd") {
		return fmt.Errorf("not support cmd \"fallocate\" and \"dd\", can not fill cache")
	}

	if !cmdexec.SupportCmd("mount") {
		return fmt.Errorf("not support cmd \"mount\", can not fill cache")
	}

	return nil
}

func execInject(ctx context.Context, fault string, args []string) error {
	percent, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("percent is not a num: %s", err.Error())
	}
	if err := memory.FillCache(ctx, percent, args[1], args[2], args[3]); err != nil {
		return fmt.Errorf("fill cache error: %s", err.Error())
	}

	return nil
}

func execRecover(ctx context.Context, fault string, args []string) error {
	fillDir := args[0]

	isDirExist, err := filesys.ExistPath(fillDir)
	if err != nil {
		return fmt.Errorf("check tmpfs[%s] exist error: %s", fillDir, err.Error())
	}

	if isDirExist {
		return memory.UndoTmpfs(ctx, fillDir)
	}

	return nil
}
