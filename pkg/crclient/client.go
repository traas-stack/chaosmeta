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

package crclient

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/crclient/docker"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"time"
)

const (
	CrLocal      = "local"
	CrDocker     = "docker"
	CrContainerd = "containerd"

	defaultDockerSocket = "unix:///var/run/docker.sock"
)

type Client interface {
	GetPidById(ctx context.Context, containerID string) (int, error)
	ListId(ctx context.Context) ([]string, error)
	KillContainerById(ctx context.Context, containerID string) error
	RmFContainerById(ctx context.Context, containerID string) error
	RestartContainerById(ctx context.Context, containerID string, timeout *time.Duration) error
}

func GetClient(cr string) (Client, error) {
	log.GetLogger().Debugf("create %s client", cr)

	switch cr {
	case CrDocker:
		return docker.GetClient(defaultDockerSocket)
	case CrContainerd:
		return nil, fmt.Errorf("to be supported")
	default:
		return nil, fmt.Errorf("not support container runtime: %s", cr)
	}
}
