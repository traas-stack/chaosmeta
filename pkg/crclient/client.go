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
	"github.com/traas-stack/chaosmetad/pkg/crclient/base"
	"github.com/traas-stack/chaosmetad/pkg/crclient/docker"
	"github.com/traas-stack/chaosmetad/pkg/crclient/pouch"
	"github.com/traas-stack/chaosmetad/pkg/log"
)

const (
	CrDocker     = "docker"
	CrContainerd = "containerd"
	CrPouch      = "pouch"

	// /var/run/pouchd.sock
	// /var/run/containerd.sock
)

type Client interface {
	GetPidById(ctx context.Context, containerID string) (int, error)
	ListId(ctx context.Context) ([]string, error)
	KillContainerById(ctx context.Context, containerID string) error
	PauseContainerById(ctx context.Context, containerID string) error
	UnPauseContainerById(ctx context.Context, containerID string) error
	RmFContainerById(ctx context.Context, containerID string) error
	RestartContainerById(ctx context.Context, containerID string, timeout int64) error
	CpFile(ctx context.Context, containerID, src, dst string) error
	Exec(ctx context.Context, containerID, cmd string) (string, error)
	GetAllPidList(ctx context.Context, containerID string) ([]base.SimpleProcess, error)
}

func GetClient(ctx context.Context, cr string) (Client, error) {
	log.GetLogger(ctx).Debugf("get %s client", cr)

	switch cr {
	case CrDocker:
		return docker.GetClient(ctx)
	case CrContainerd:
		return nil, fmt.Errorf("to be supported")
	case CrPouch:
		return pouch.GetClient(ctx)
	default:
		return nil, fmt.Errorf("not support container runtime: %s", cr)
	}
}
