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
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/tools/common"
	"os"
	"strconv"
)

// uid dir prefix start end timeout
func main() {
	args := os.Args
	if len(args) < 6 {
		common.ExitWithErr("must provide 6 args: uid, dir prefix start end timeout")
	}

	dirname, prefix, startStr, endStr, timeoutStr := args[2], args[3], args[4], args[5], args[6]

	var timeout int
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("timeout value is not a valid int, error: %s\n", err.Error()))
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("start value is not a valid int, error: %s\n", err.Error()))
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("end value is not a valid int, error: %s\n", err.Error()))
	}

	for i := start; i < end; i++ {
		_, err := os.Open(fmt.Sprintf("%s/%s%d", dirname, prefix, i))
		if err != nil {
			fmt.Printf("[warn]open file %d error: %s\n", i, err.Error())
		}
	}

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}

// It takes 8M to enable the coroutine, and only 4M to disable it
//func occupyFd(file string) {
//	time.Sleep(5 * time.Second)
//	for {
//		_, err := os.Open(file)
//		if err == nil {
//			break
//		}
//		// 获取错误，不管，继续获取
//		//fmt.Println(err.Error())
//		time.Sleep(1 * time.Second)
//	}
//
//	select {}
//}
