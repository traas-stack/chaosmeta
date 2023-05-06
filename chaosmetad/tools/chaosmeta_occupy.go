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
	"github.com/traas-stack/chaosmeta/chaosmetad/tools/common"
	"net"
	"os"
	"strconv"
)

// [uid] [port] [protocol] [timeout]
func main() {
	args := os.Args
	if len(args) < 5 {
		common.ExitWithErr("must provide 4 args: uid、port、protocol、timeout")
	}

	p, proto, t := args[2], args[3], args[4]
	port, err := strconv.Atoi(p)
	if err != nil || port <= 0 {
		common.ExitWithErr("port is invalid")
	}

	if proto != "tcp" && proto != "udp" && proto != "tcp6" && proto != "udp6" {
		common.ExitWithErr("proto only support: udp、tcp、udp6、tcp6")
	}

	if proto == "tcp" {
		proto = "tcp4"
	} else if proto == "udp" {
		proto = "udp4"
	}

	var timeout int
	timeout, err = strconv.Atoi(t)
	if err != nil {
		common.ExitWithErr(fmt.Sprintf("timeout value is not a valid int, error: %s", err.Error()))
	}

	if proto == "tcp4" || proto == "tcp6" {
		if _, err := net.Listen(proto, fmt.Sprintf(":%d", port)); err != nil {
			common.ExitWithErr(fmt.Sprintf("%s listen on %d error: %s", proto, port, err.Error()))
		}
	} else {
		if _, err := net.ListenUDP(proto, &net.UDPAddr{
			Port: port,
		}); err != nil {
			common.ExitWithErr(fmt.Sprintf("%s listen on %d error: %s", proto, port, err.Error()))
		}
	}

	fmt.Println("[success]inject success")

	common.SleepWait(timeout)
}
