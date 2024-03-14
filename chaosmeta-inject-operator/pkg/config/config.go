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

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(path string) (*MainConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config file error: %s", err.Error())
	}
	var mainConfig MainConfig
	if err := json.Unmarshal(configBytes, &mainConfig); err != nil {
		return nil, fmt.Errorf("convert config file error: %s", err.Error())
	}

	return &mainConfig, nil
}

type MainConfig struct {
	Worker   WorkerConfig   `json:"worker"`
	Ticker   TickerConfig   `json:"ticker"`
	Executor ExecutorConfig `json:"executor"`
}

type WorkerConfig struct {
	PoolCount int `json:"poolCount"`
}

type TickerConfig struct {
	AutoCheckInterval int `json:"autoCheckInterval"`
}

type ExecutorConfig struct {
	Mode             string                  `json:"mode"`
	Executor         string                  `json:"executor"`
	Version          string                  `json:"version"`
	AgentConfig      AgentExecutorConfig     `json:"agentConfig"`
	DaemonsetConfig  DaemonsetExecutorConfig `json:"daemonsetConfig"`
	MiddlewareConfig MiddlewareConfig        `json:"middlewareConfig"`
}

type MiddlewareConfig struct {
	Url        string     `json:"url"`
	MistConfig MistConfig `json:"mistConfig"`
}

type MistConfig struct {
	AntVipUrl  string `json:"antVipUrl"`
	BkmiUrl    string `json:"bkmiUrl"`
	AppName    string `json:"appName"`
	Tenant     string `json:"tenant"`
	Mode       string `json:"mode"`
	SecretName string `json:"secretName"`
}

type AgentExecutorConfig struct {
	AgentPort int `json:"agentPort"`
}

type DaemonsetExecutorConfig struct {
	LocalExecPath string `json:"localExecPath"`

	DaemonNs          string            `json:"daemonNs"`
	DaemonLabel       map[string]string `json:"daemonLabel"`
	DaemonName        string            `json:"daemonName"`
	AutoLabelNode     bool              `json:"autoLabelNode"`
	NodeSelectorLabel map[string]string `json:"nodeSelectorLabel"`
}
