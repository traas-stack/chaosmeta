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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/crclient/base"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	defaultSocket     = "unix:///var/run/docker.sock"
	dockerVersionKey  = "DOCKER_API_VERSION"
	defaultAPIVersion = "1.24"
)

type Client struct {
	client *dockerClient.Client
}

var (
	clientInstance *Client
	mutex          sync.Mutex
)

func GetClient(ctx context.Context) (d *Client, err error) {
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
			log.GetLogger(ctx).Debug("new docker client")
			if version := os.Getenv(dockerVersionKey); version == "" {
				if err := os.Setenv(dockerVersionKey, defaultAPIVersion); err != nil {
					return nil, fmt.Errorf("set DOCKER_API_VERSION error: %s", err.Error())
				}
			}

			cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithHost(defaultSocket))
			if err != nil {
				return nil, fmt.Errorf("new docker client error: %s", err.Error())
			}

			log.GetLogger(ctx).Debugf("docker client version: %s", cli.ClientVersion())
			clientInstance = &Client{
				client: cli,
			}
		}
		mutex.Unlock()
	}

	return clientInstance, nil
}

func (d *Client) GetAllPidList(ctx context.Context, containerID string) ([]base.SimpleProcess, error) {
	re, err := d.client.ContainerTop(ctx, containerID, nil)
	if err != nil {
		return nil, fmt.Errorf("get process info from client error: %s", err.Error())
	}

	var rePro = make([]base.SimpleProcess, len(re.Processes))
	for i := 0; i < len(re.Processes); i++ {
		for j := 0; j < len(re.Titles); j++ {
			if re.Titles[j] == "PID" {
				pid, err := strconv.Atoi(re.Processes[i][j])
				if err != nil {
					return nil, fmt.Errorf("PID[%s] is not a num: %s", re.Processes[i][j], err.Error())
				}
				rePro[i].Pid = pid
			} else if re.Titles[j] == "CMD" {
				rePro[i].Cmd = re.Processes[i][j]
			}
		}
	}

	return rePro, nil
}

func (d *Client) GetPidById(ctx context.Context, containerID string) (int, error) {
	info, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return -1, fmt.Errorf("get meta data of container[%s] error: %s", containerID, err.Error())
	}

	if info.HostConfig.Runtime != "runc" {
		return -1, fmt.Errorf("only support: runc, not support: %s", info.HostConfig.Runtime)
	}

	if info.State.Pid <= 0 {
		return -1, fmt.Errorf("no such container[%s]", containerID)
	}

	return info.State.Pid, nil
}

// Exec TODO: now output has extra space prefix, need to fix this bug
func (d *Client) Exec(ctx context.Context, containerID, cmd string) (string, error) {
	logger := log.GetLogger(ctx)
	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash", "-c", cmd},
	}

	resp, err := d.client.ContainerExecCreate(ctx, containerID, execOpts)
	if err != nil {
		return "", fmt.Errorf("container exec create error: %s", err.Error())
	}
	attach, err := d.client.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("container exec attach error: %s", err.Error())
	}

	defer attach.Close()
	dataBytes, err := ioutil.ReadAll(attach.Reader)
	if err != nil {
		return "", fmt.Errorf("read container exec data error: %s", err.Error())
	}

	data := string(dataBytes)
	logger.Debugf("container exec output: %s", data)
	execInspect, err := d.client.ContainerExecInspect(context.Background(), resp.ID)
	if err != nil {
		return "", fmt.Errorf("inspect container exec result error: %s", err.Error())
	}

	if execInspect.ExitCode != 0 {
		return "", fmt.Errorf("exit code: %d, output: %s", execInspect.ExitCode, string(data))
	}

	if len(data) >= 8 {
		data = data[8:]
	}

	return data, nil
}

// KillContainerById convert to static container
func (d *Client) KillContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerKill(ctx, containerID, "SIGKILL")
}

// RmFContainerById remove container
func (d *Client) RmFContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

func (d *Client) PauseContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerPause(ctx, containerID)
}

func (d *Client) UnPauseContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerUnpause(ctx, containerID)
}

func (d *Client) RestartContainerById(ctx context.Context, containerID string, timeout int64) error {
	var waitTime = time.Second * time.Duration(timeout)
	return d.client.ContainerRestart(ctx, containerID, &waitTime)
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

func (d *Client) CpFile(ctx context.Context, containerID, src, dst string) error {
	info, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get meta data of container[%s] error: %s", containerID, err.Error())
	}

	containerMergedDir, ok := info.GraphDriver.Data["MergedDir"]
	if !ok {
		return fmt.Errorf("get containerMergedDir error: not existed")
	}

	dst = filepath.Join(containerMergedDir, dst)
	log.GetLogger(ctx).Debugf("copy file from %s to container %s", dst)
	return base.CopyFile(src, dst)
}
