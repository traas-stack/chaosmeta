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

package agentexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/base"
	httpclient "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/http"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"strconv"
	"strings"
)

type AgentRemoteExecutor struct {
	Client      *httpclient.HTTPClient
	ServicePort int
	Version     string
}

func (r *AgentRemoteExecutor) CheckExecutorWay(ctx context.Context) error {

	return nil
}

func (r *AgentRemoteExecutor) CheckAlive(ctx context.Context, injectObject string) error {
	resBytes, err := r.Client.Get(ctx, fmt.Sprintf("http://%s:%d/v1/version", injectObject, r.ServicePort))
	if err != nil {
		return fmt.Errorf("get response error: %s", err.Error())
	}

	var resp base.VersionResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return fmt.Errorf("resp[%s] format error: %s", string(resBytes), err.Error())
	}

	if resp.Data == nil || resp.Code != 0 {
		return fmt.Errorf("query version error: %s", resp.Message)
	}

	if resp.Data.Version != r.Version {
		return fmt.Errorf("expected version %s, but get %s", r.Version, resp.Data.Version)
	}

	return nil
}

// Init install agent
func (r *AgentRemoteExecutor) Init(ctx context.Context, target string) error {
	return nil
}

func (r *AgentRemoteExecutor) Inject(ctx context.Context, injectObject string, target, fault, uid, timeout, cID, cRuntime string, args []v1alpha1.ArgsUnit) error {
	if err := r.CheckAlive(ctx, injectObject); err != nil {
		return fmt.Errorf("check target's status error: %s", err.Error())
	}

	var argsMap = make(map[string]interface{})
	for _, unitArgs := range args {
		if unitArgs.Key == v1alpha1.ContainerKey {
			continue
		}

		unitArgs.Key = strings.ReplaceAll(unitArgs.Key, "-", "_")
		if unitArgs.ValueType == v1alpha1.IntVType {
			argsInt, err := strconv.Atoi(unitArgs.Value)
			if err != nil {
				return fmt.Errorf("args[%s]'s value[%s] require int type", unitArgs.Key, unitArgs.Value)
			}

			argsMap[unitArgs.Key] = argsInt
		} else if unitArgs.ValueType == v1alpha1.StringVType {
			argsMap[unitArgs.Key] = unitArgs.Value
		} else {
			return fmt.Errorf("args[%s] not support value type: %s", unitArgs.Key, unitArgs.ValueType)
		}
	}

	argsBytes, err := json.Marshal(argsMap)
	if err != nil {
		return fmt.Errorf("args to json string error: %s", err.Error())
	}

	bytesData, err := json.Marshal(base.InjectRequest{
		Target:           target,
		Fault:            fault,
		Timeout:          timeout,
		ContainerId:      cID,
		ContainerRuntime: cRuntime,
		Uid:              uid,
		Args:             string(argsBytes),
	})
	if err != nil {
		return fmt.Errorf("request to string error: %s", err.Error())
	}

	resBytes, err := r.Client.Post(ctx, fmt.Sprintf("http://%s:%d/v1/experiment/inject", injectObject, r.ServicePort), bytesData)
	if err != nil {
		return fmt.Errorf("get response error: %s", err.Error())
	}

	var resp base.InjectResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return fmt.Errorf("resp[%s] format error: %s", string(resBytes), err.Error())
	}

	if resp.Code == base.SucCode {
		return nil
	} else {
		return fmt.Errorf("err code: {%d}, err msg: %s", resp.Code, resp.Message)
	}
}

func (r *AgentRemoteExecutor) Recover(ctx context.Context, injectObject string, uid string) error {
	bytesData, err := json.Marshal(base.RecoverRequest{
		Uid: uid,
	})

	if err != nil {
		return fmt.Errorf("request to string error: %s", err.Error())
	}

	resBytes, err := r.Client.Post(ctx, fmt.Sprintf("http://%s:%d/v1/experiment/recover", injectObject, r.ServicePort), bytesData)
	if err != nil {
		return fmt.Errorf("get response error: %s", err.Error())
	}

	var resp base.CommonResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return fmt.Errorf("resp[%s] format error: %s", string(resBytes), err.Error())
	}

	if resp.Code == base.SucCode {
		return nil
	} else {
		return fmt.Errorf("err code: {%d}, err msg: %s", resp.Code, resp.Message)
	}
}

func (r *AgentRemoteExecutor) Query(ctx context.Context, injectObject string, uid string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	bytesData, err := json.Marshal(base.QueryRequest{
		Uid: uid,
	})

	if err != nil {
		return nil, fmt.Errorf("request to string error: %s", err.Error())
	}

	resBytes, err := r.Client.Post(ctx, fmt.Sprintf("http://%s:%d/v1/experiment/query", injectObject, r.ServicePort), bytesData)
	if err != nil {
		return nil, fmt.Errorf("get response error: %s", err.Error())
	}

	var resp base.QueryResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return nil, fmt.Errorf("resp[%s] format error: %s", string(resBytes), err.Error())
	}

	if resp.Code == base.SucCode {
		if resp.Data == nil || resp.Data.Total == 0 {
			return nil, fmt.Errorf("task not found")
		}
		task := resp.Data.Experiments[0]

		return &model.SubExpInfo{
			UID:        uid,
			CreateTime: task.CreateTime,
			UpdateTime: task.UpdateTime,
			Message:    task.Error_,
			Status:     base.ConvertStatus(task.Status, phase),
		}, nil
	} else {
		return nil, fmt.Errorf("err code: {%d}, err msg: %s", resp.Code, resp.Message)
	}
}
