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

package tcpexecutor

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
	PortArgsKey    = "port"
	TimeoutArgsKey = "timeout"
)

func init() {
	e, err := NewTCPExecutor(context.Background())
	if err != nil {
		fmt.Printf("new tcp executor error: %s\n", err.Error())
	} else {
		v1alpha1.SetMeasureExecutor(context.Background(), v1alpha1.TCPMeasureType, e)
	}
}

type TCPExecutor struct {
}

func NewTCPExecutor(ctx context.Context) (*TCPExecutor, error) {
	return &TCPExecutor{}, nil
}

func (e *TCPExecutor) CheckConfig(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement) error {
	_, err := utils.GetArgsValueStr(args, IPArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	port, err := utils.GetArgsValueStr(args, PortArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port args is invalid")
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

func (e *TCPExecutor) InitialData(ctx context.Context, args []v1alpha1.MeasureArgs) (string, error) {
	return "", nil
}

func (e *TCPExecutor) Measure(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement, initialData string) error {
	ip, _ := utils.GetArgsValueStr(args, IPArgsKey)
	port, _ := utils.GetArgsValueStr(args, PortArgsKey)
	timeout, _ := utils.GetArgsValueStr(args, TimeoutArgsKey)
	t, _ := strconv.Atoi(timeout)
	duration := time.Duration(t) * time.Second

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), duration)
	if err == nil {
		conn.Close()
	}

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
