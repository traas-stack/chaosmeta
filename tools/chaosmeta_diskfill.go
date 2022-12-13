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
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector/disk"
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"os"
	"strconv"
)

// [func] [args]
func main() {
	fName, args := os.Args[1], os.Args[2:]
	var err error
	if fName == "validator" {
		err = validatorExp(args)
	} else if fName == "inject" {
		err = injectExp(args)
	} else if fName == "recover" {
		err = recoverExp(args)
	} else {
		common.ExitWithErr(fmt.Sprintf("not support method: %s", fName))
	}

	if err != nil {
		common.ExitWithErr(err.Error())
	}
}

func validatorExp(args []string) error {
	percentStr, bytes, dir := args[0], args[1], args[2]
	percent, _ := strconv.Atoi(percentStr)
	return disk.ValidatorDiskFill(context.Background(), percent, bytes, dir)
}

func injectExp(args []string) error {
	percentStr, bytes, dir, uid := args[0], args[1], args[2], args[3]
	percent, _ := strconv.Atoi(percentStr)
	return disk.InjectDiskFill(context.Background(), percent, bytes, dir, uid)
}

func recoverExp(args []string) error {
	return disk.RecoverDiskFill(context.Background(), args[0], args[1])
}
