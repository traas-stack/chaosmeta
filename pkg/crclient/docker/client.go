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

package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"sync"
	"time"
)

type Client struct {
	client *dockerClient.Client
}

var (
	clientInstance *Client
	mutex          sync.Mutex
)

func GetClient(ctx context.Context, host string) (d *Client, err error) {
	defer func() {
		if e := recover(); e != any(nil) {
			// catch exception from create client
			mutex.Unlock()
			err = fmt.Errorf("catch exception: %v", e)
		}
	}()

	if clientInstance == nil {
		mutex.Lock()
		if clientInstance == nil {
			cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithHost(host))
			if err != nil {
				return nil, fmt.Errorf("new docker client error: %s", err.Error())
			}

			clientInstance = &Client{
				client: cli,
			}
		}
		mutex.Unlock()
	}

	return clientInstance, nil
}

func (d *Client) GetPidById(ctx context.Context, containerID string) (int, error) {
	info, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return -1, fmt.Errorf("get meta data of container[%s] error: %s", containerID, err.Error())
	}

	return info.State.Pid, nil
}

func (d *Client) KillContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerKill(ctx, containerID, "SIGKILL")
}

func (d *Client) RmFContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

func (d *Client) RestartContainerById(ctx context.Context, containerID string, timeout *time.Duration) error {
	return d.client.ContainerRestart(ctx, containerID, timeout)
}

func (d *Client) ListId(ctx context.Context) ([]string, error) {
	containerList, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("get container list error: %s", err.Error())
	}

	var idList = make([]string, len(containerList))
	for i, c := range containerList {
		idList[i] = c.ID
	}

	return idList, nil
}
