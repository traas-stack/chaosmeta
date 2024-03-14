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

package remoteexecutor

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/config"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/agentexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/daemonsetexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/middlewareexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/middlewareexecutor/tse"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/middlewareexecutor/tse/auth"
	httpclient "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/http"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"net/http"
)

type RemoteModeType string

const (
	AgentRemoteMode      RemoteModeType = "agent"
	DaemonsetRemoteMode  RemoteModeType = "daemonset"
	MiddlewareRemoteMode RemoteModeType = "middleware"
)

type RemoteExecutor interface {
	CheckExecutorWay(ctx context.Context) error
	// CheckAlive check target object service alive
	CheckAlive(ctx context.Context, injectObject string) error
	// Init install agent
	Init(ctx context.Context, target string) error
	Inject(ctx context.Context, injectObject string, target, fault, uid, timeout, cID, cRuntime string, args []v1alpha1.ArgsUnit) error
	Recover(ctx context.Context, injectObject string, uid string) error
	Query(ctx context.Context, injectObject string, uid string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error)
	//SyncStatus(ctx context.Context, exp *v1alpha1.ExperimentStatus)
}

var globalRemoteExecutor RemoteExecutor

func AutoSelectRemoteExecutor(config *config.ExecutorConfig, restConfig *rest.Config, schema *runtime.Scheme) error {
	/**
	 * 启动时通道自动选择
	 */
	ctx := context.Background()
	daemonSetRemoteExecutor := &daemonsetexecutor.DaemonsetRemoteExecutor{
		RESTConfig: restConfig,
		Schema:     schema,

		LocalExecPath: config.DaemonsetConfig.LocalExecPath,
		Executor:      config.Executor,
		Version:       config.Version,

		DaemonsetNs:    config.DaemonsetConfig.DaemonNs,
		DaemonsetLabel: config.DaemonsetConfig.DaemonLabel,
		DaemonsetName:  config.DaemonsetConfig.DaemonName,
	}
	if err := daemonSetRemoteExecutor.CheckExecutorWay(ctx); err == nil {
		globalRemoteExecutor = daemonSetRemoteExecutor
		return nil
	}
	agentExecutor := &agentexecutor.AgentRemoteExecutor{
		Client: &httpclient.HTTPClient{
			Client: &http.Client{},
		},
		Version:     config.Version,
		ServicePort: config.AgentConfig.AgentPort,
	}
	if err := agentExecutor.CheckExecutorWay(ctx); err == nil {
		globalRemoteExecutor = agentExecutor
	}
	middlewareRemoteExecutor := &middlewareexecutor.MiddleWareExecutor{
		Middleware: &tse.TseMiddleware{
			Config: config.MiddlewareConfig,
			MistClient: auth.MistClient{
				Config: config.MiddlewareConfig.MistConfig,
			},
		},
	}
	if err := middlewareRemoteExecutor.CheckExecutorWay(ctx); err == nil {
		globalRemoteExecutor = middlewareRemoteExecutor
	}
	return nil
}

func SetGlobalRemoteExecutor(config *config.ExecutorConfig, restConfig *rest.Config, schema *runtime.Scheme) error {
	switch RemoteModeType(config.Mode) {
	case AgentRemoteMode:
		globalRemoteExecutor = &agentexecutor.AgentRemoteExecutor{
			Client: &httpclient.HTTPClient{
				Client: &http.Client{},
			},
			Version:     config.Version,
			ServicePort: config.AgentConfig.AgentPort,
		}
	case DaemonsetRemoteMode:
		globalRemoteExecutor = &daemonsetexecutor.DaemonsetRemoteExecutor{
			//ApiServer:  apiServer,
			RESTConfig: restConfig,
			Schema:     schema,

			LocalExecPath: config.DaemonsetConfig.LocalExecPath,
			Executor:      config.Executor,
			Version:       config.Version,

			DaemonsetNs:    config.DaemonsetConfig.DaemonNs,
			DaemonsetLabel: config.DaemonsetConfig.DaemonLabel,

			//AutoLabelNode:     config.DaemonsetConfig.AutoLabelNode,
			//NodeSelectorLabel: config.DaemonsetConfig.NodeSelectorLabel,
		}
	case MiddlewareRemoteMode:
		globalRemoteExecutor = &middlewareexecutor.MiddleWareExecutor{
			Middleware: &tse.TseMiddleware{
				Config: config.MiddlewareConfig,
				MistClient: auth.MistClient{
					Config: config.MiddlewareConfig.MistConfig,
				},
			},
		}
	default:
		return fmt.Errorf("not support remote executor: %s", config.Mode)
	}

	return nil
}

func GetRemoteExecutor() RemoteExecutor {
	return globalRemoteExecutor
}
