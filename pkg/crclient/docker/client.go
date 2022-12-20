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
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultDockerSocket = "unix:///var/run/docker.sock"
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

func (d *Client) Exec(ctx context.Context, containerID, cmd string) (string, error) {
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
	data, err := ioutil.ReadAll(attach.Reader)
	if err != nil {
		return "", fmt.Errorf("read container exec date error: %s", err.Error())
	}

	return string(data), nil
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

func (d *Client) CpFile(ctx context.Context, containerID, src, dst string) error {
	dstInfo := archive.CopyInfo{Path: dst}
	dstStat, err := d.client.ContainerStatPath(ctx, containerID, dst)

	if err == nil && dstStat.Mode&os.ModeSymlink != 0 {
		linkTarget := dstStat.LinkTarget
		if !system.IsAbs(linkTarget) {
			dstParent, _ := archive.SplitPathDirEntry(dst)
			linkTarget = filepath.Join(dstParent, linkTarget)
		}

		dstInfo.Path = linkTarget
		dstStat, err = d.client.ContainerStatPath(ctx, containerID, linkTarget)
	}

	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
	}

	var (
		content         io.Reader
		resolvedDstPath string
	)

	srcInfo, err := archive.CopyInfoSourcePath(src, true)
	if err != nil {
		return err
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}

	defer srcArchive.Close()

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()

	resolvedDstPath = dstDir
	content = preparedArchive

	return d.client.CopyToContainer(ctx, containerID, resolvedDstPath, content, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}
