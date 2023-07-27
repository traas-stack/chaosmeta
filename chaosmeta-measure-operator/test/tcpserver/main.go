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
	"net"
)

func handleTCP(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	reqData := string(buf)
	respData := "Hello, " + reqData

	_, err = conn.Write([]byte(respData))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}
}

func main() {
	tcpListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("failed to start TCP server:", err)
		return
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("failed to accept TCP connection:", err)
			continue
		}
		go handleTCP(conn)
	}
}
