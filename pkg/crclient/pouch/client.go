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

package pouch

import (
	"context"
	"fmt"
	"github.com/alibaba/pouch/apis/types"
	pouchClient "github.com/alibaba/pouch/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/traas-stack/chaosmetad/pkg/crclient/base"
	"github.com/traas-stack/chaosmetad/pkg/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	defaultSocket = "unix:///var/run/pouchd.sock"
)

type Client struct {
	client pouchClient.CommonAPIClient
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
			log.GetLogger(ctx).Debug("new pouch client")
			cli, err := pouchClient.NewAPIClient(defaultSocket, pouchClient.TLSConfig{})
			if err != nil {
				return nil, fmt.Errorf("new pouch client error: %s", err.Error())
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
	info, err := d.client.ContainerGet(ctx, containerID)
	if err != nil {
		return -1, fmt.Errorf("get meta data of container[%s] error: %s", containerID, err.Error())
	}

	if info.State.Pid <= 0 {
		return -1, fmt.Errorf("no such container[%s]", containerID)
	}

	return int(info.State.Pid), nil
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

// Exec TODO: now output has extra space prefix, need to fix this bug
func (d *Client) Exec(ctx context.Context, containerID, cmd string) (string, error) {
	logger := log.GetLogger(ctx)
	logger.Debugf("container exec cmd: %s", cmd)
	execOpts := &types.ExecCreateConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash", "-c", cmd},
	}

	resp, err := d.client.ContainerCreateExec(ctx, containerID, execOpts)
	if err != nil {
		return "", fmt.Errorf("container exec create error: %s", err.Error())
	}

	conn, r, err := d.client.ContainerStartExec(ctx, resp.ID, &types.ExecStartConfig{})
	if err != nil {
		return "", fmt.Errorf("container exec start error: %s", err.Error())
	}

	defer conn.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("read container exec date error: %s", err.Error())
	}

	logger.Debugf("container exec output: %s", string(data))
	return string(data), nil
}

// KillContainerById convert to static container
func (d *Client) KillContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerKill(ctx, containerID, "SIGKILL")
}

func (d *Client) PauseContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerPause(ctx, containerID)
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

func (d *Client) UnPauseContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerUnpause(ctx, containerID)
}

// RmFContainerById remove container
func (d *Client) RmFContainerById(ctx context.Context, containerID string) error {
	return d.client.ContainerRemove(ctx, containerID, &types.ContainerRemoveOptions{Force: true})
}

func (d *Client) RestartContainerById(ctx context.Context, containerID string, timeout int64) error {
	return d.client.ContainerRestart(ctx, containerID, strconv.Itoa(int(timeout)))
}

func (d *Client) CpFile(ctx context.Context, containerID, src, dst string) error {
	dstInfo := archive.CopyInfo{Path: dst}
	dstStat, err := d.client.ContainerStatPath(ctx, containerID, dst)

	if err == nil && os.FileMode(dstStat.Mode)&os.ModeSymlink != 0 {
		linkTarget := dstStat.Path
		if !system.IsAbs(linkTarget) {
			dstParent, _ := archive.SplitPathDirEntry(dst)
			linkTarget = filepath.Join(dstParent, linkTarget)
		}

		dstInfo.Path = linkTarget
		dstStat, err = d.client.ContainerStatPath(ctx, containerID, linkTarget)
	}

	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, os.FileMode(dstStat.Mode).IsDir()
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

	return d.client.CopyToContainer(ctx, containerID, resolvedDstPath, content)
}
