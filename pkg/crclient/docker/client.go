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
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cgroup"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

const (
	defaultDockerSocket = "unix:///var/run/docker.sock"
	cgroupWaitInterval = time.Millisecond * 200
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
			cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithHost(defaultDockerSocket))
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

func (d *Client) GetCgroupPath(ctx context.Context, containerID, subSys string) (string, error) {
	client, err := GetClient(ctx)
	if err != nil {
		return "", fmt.Errorf("get client error: %s", err.Error())
	}

	pid, err := client.GetPidById(context.Background(), containerID)
	if err != nil {
		return "", fmt.Errorf("get pid of container[%s] error: %s", containerID, err.Error())
	}

	cPath, err := cgroup.GetpidCurCgroup(ctx, pid, subSys)
	if err != nil {
		return "", fmt.Errorf("get cgroup[%s] path of process[%d] error: %s", subSys, pid, err.Error())
	}

	return cPath, nil
}

func (d *Client) ExecContainer(ctx context.Context, containerID string, namespaces []string, cmd string) error {
	logger := log.GetLogger(ctx)
	// get container's init process
	client, err := GetClient(ctx)
	if err != nil {
		return fmt.Errorf("get client error: %s", err.Error())
	}

	targetPid, err := client.GetPidById(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get pid of container[%s]'s init process error: %s", containerID, err.Error())
	}

	// exec ns
	var nsOptionStr string
	for _, unitNs := range namespaces {
		switch unitNs {
		case namespace.MNT:
			nsOptionStr += " -m"
		case namespace.PID:
			nsOptionStr += " -p"
		case namespace.UTS:
			nsOptionStr += " -u"
		case namespace.NET:
			nsOptionStr += " -n"
		case namespace.IPC:
			nsOptionStr += " -i"
		}
	}

	c := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s -t %d %s -c \"%s\"", utils.GetToolPath(cmdexec.ExecnsKey), targetPid, nsOptionStr, cmd))
	logger.Debugf("container exec cmd: %s", c.Args)

	// TODO: need to catch execns error
	//var stdout, stderr bytes.Buffer
	//c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Start(); err != nil {
		return fmt.Errorf("start process error: %s", err.Error())
	}

	// set cgroup for new process
	if err := cgroup.AddToProCgroup(c.Process.Pid, targetPid); err != nil {
		if err := c.Process.Kill(); err != nil {
			logger.Warnf("undo: kill container exec process[%d] error: %s", c.Process.Pid, err.Error())
		}

		return fmt.Errorf("add process[%d] to container[%d] cgroup error: %s", c.Process.Pid, targetPid, err.Error())
	}

	time.Sleep(cgroupWaitInterval)
	if err := c.Process.Signal(syscall.SIGCONT); err != nil {
		return err
	}

	// signal continue
	return nil
}
