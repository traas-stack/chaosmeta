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

package ipexecutor

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/utils"
	"net"
	"strconv"
	"time"
)

const (
	IPArgsKey      = "ip"
	TimeoutArgsKey = "timeout"
)

func init() {
	e, err := NewIPExecutor(context.Background())
	if err != nil {
		fmt.Printf("new ip executor error: %s\n", err.Error())
	} else {
		v1alpha1.SetMeasureExecutor(context.Background(), v1alpha1.IPMeasureType, e)
	}
}

type IPExecutor struct {
}

func NewIPExecutor(ctx context.Context) (*IPExecutor, error) {
	return &IPExecutor{}, nil
}

func (e *IPExecutor) CheckConfig(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement) error {
	_, err := utils.GetArgsValueStr(args, IPArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	timeout, err := utils.GetArgsValueStr(args, TimeoutArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	_, err = strconv.Atoi(timeout)
	if err != nil {
		return fmt.Errorf("timeout args is not a int: %s", err.Error())
	}

	if judgement.JudgeType != v1alpha1.ConnectivityJudgeType {
		return fmt.Errorf("tcp measure only support judge type: %s", v1alpha1.ConnectivityJudgeType)
	}

	if judgement.JudgeValue != v1alpha1.ConnectivityTrue && judgement.JudgeValue != v1alpha1.ConnectivityFalse {
		return fmt.Errorf("value of %s judge type only support: %s, %s", v1alpha1.ConnectivityJudgeType, v1alpha1.ConnectivityTrue, v1alpha1.ConnectivityFalse)
	}

	return nil
}

func (e *IPExecutor) InitialData(ctx context.Context, args []v1alpha1.MeasureArgs) (string, error) {
	return "", nil
}

func ping(ip string, timeout time.Duration) error {
	conn, err := net.DialTimeout("ip4:icmp", ip, timeout)
	if err != nil {
		return fmt.Errorf("create connect error: %s", err.Error())
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("set timeout for connection error: %s", err.Error())
	}

	msg := make([]byte, 1024)
	msg[0] = 8  // echo
	msg[1] = 0  // code
	msg[2] = 0  // checksum
	msg[3] = 0  // checksum
	msg[4] = 0  // identifier[0]
	msg[5] = 13 // identifier[1]
	msg[6] = 0  // sequence[0]
	msg[7] = 37 // sequence[1]

	checksum := utils.CheckSum(msg)

	msg[2] = byte(checksum >> 8)
	msg[3] = byte(checksum & 0xff)

	_, err = conn.Write(msg[:8])
	if err != nil {
		return fmt.Errorf("send icmp message error: %s", err.Error())
	}

	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	if err != nil {
		return fmt.Errorf("recv icmp message error: %s", err.Error())
	}

	return err
}

func (e *IPExecutor) Measure(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement, initialData string) error {
	ipStr, _ := utils.GetArgsValueStr(args, IPArgsKey)
	timeout, _ := utils.GetArgsValueStr(args, TimeoutArgsKey)
	t, _ := strconv.Atoi(timeout)
	duration := time.Duration(t) * time.Second
	err := ping(ipStr, duration)

	if judgement.JudgeType == v1alpha1.ConnectivityJudgeType {
		if judgement.JudgeValue == v1alpha1.ConnectivityTrue {
			return err
		} else {
			if err == nil {
				return fmt.Errorf("expect connectivity is false, but get true")
			} else {
				return nil
			}
		}
	} else {
		return fmt.Errorf("not support judge type: %s", judgement.JudgeType)
	}
}
