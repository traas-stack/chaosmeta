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

package containerd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/google/uuid"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/shirou/gopsutil/process"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/crclient/base"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	defaultSocket = "/run/containerd/containerd.sock"

	termExitCode = 143
	// TODO: need to be a command line args
	defaultNS = "k8s.io"
	//defaultNS = "moby"
)

type Client struct {
	client *containerd.Client
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
			log.GetLogger(ctx).Debugf("new containerd client, ns: %s, socket: %s", defaultNS, defaultSocket)
			cli, err := containerd.New(defaultSocket, containerd.WithDefaultNamespace(defaultNS))
			if err != nil {
				return nil, fmt.Errorf("new containerd client error: %s", err.Error())
			}

			clientInstance = &Client{
				client: cli,
			}
		}
		mutex.Unlock()
	}

	return clientInstance, nil
}

func (d *Client) getContainerTask(ctx context.Context, containerID string) (containerd.Task, error) {
	container, err := d.client.LoadContainer(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("load container error: %s", err.Error())
	}

	return container.Task(ctx, nil)
}

func (d *Client) GetAllPidList(ctx context.Context, containerID string) ([]base.SimpleProcess, error) {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("get task of container error: %s", err.Error())
	}

	procsList, err := task.Pids(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container's process error: %s", err.Error())
	}

	var reProList = make([]base.SimpleProcess, len(procsList))
	for i, proc := range procsList {
		reProList[i].Pid = int(proc.Pid)
		p, err := process.NewProcess(int32(reProList[i].Pid))
		if err != nil {
			return nil, fmt.Errorf("process[%d] is not exist, error: %s", reProList[i].Pid, err.Error())
		}

		reProList[i].Cmd, err = p.Cmdline()
		if err != nil {
			return nil, fmt.Errorf("get cmd of process[%d] error: %s", reProList[i].Pid, err.Error())
		}
	}

	return reProList, nil
}

func (d *Client) GetPidById(ctx context.Context, containerID string) (int, error) {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return 0, fmt.Errorf("get task of container error: %s", err.Error())
	}

	re := int(task.Pid())

	if re <= 0 {
		return -1, fmt.Errorf("no such container[%s]", containerID)
	}

	return re, nil
}

func (d *Client) Exec(ctx context.Context, containerID, cmd string) (string, error) {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("get task of container error: %s", err.Error())
	}

	pId := uuid.New().String()
	var stdout, stderr bytes.Buffer
	pro, err := task.Exec(ctx, pId, &specs.Process{
		Args: []string{"/bin/bash", "-c", cmd},
		Cwd:  "/",
	}, cio.NewCreator(append([]cio.Opt{cio.WithStreams(nil, &stdout, &stderr)})...))

	if err != nil {
		return "", fmt.Errorf("container exec error: %s", err.Error())
	}

	defer pro.Delete(ctx)

	eStatusC, err := pro.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("wait exec error: %s", err.Error())
	}

	if err := pro.Start(ctx); err != nil {
		return "", fmt.Errorf("task start error: %s", err.Error())
	}

	eStatus := <-eStatusC
	output := stdout.String() + stderr.String()
	if eStatus.ExitCode() != 0 {
		return output, fmt.Errorf("exit code: %d, exit error: %v, msg: %s", eStatus.ExitCode(), eStatus.Error(), output)
	}

	return output, nil
}

// KillContainerById convert to static container
func (d *Client) KillContainerById(ctx context.Context, containerID string) error {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	return task.Kill(ctx, syscall.SIGKILL)
}

// RmFContainerById remove container
func (d *Client) RmFContainerById(ctx context.Context, containerID string) error {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
		return fmt.Errorf("kill container error: %s", err.Error())
	}

	_, err = task.Delete(ctx)
	return err
}

func (d *Client) PauseContainerById(ctx context.Context, containerID string) error {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	return task.Pause(ctx)
}

func (d *Client) UnPauseContainerById(ctx context.Context, containerID string) error {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	return task.Resume(ctx)
}

func (d *Client) RestartContainerById(ctx context.Context, containerID string, timeout int64) error {
	container, err := d.client.LoadContainer(ctx, containerID)
	if err != nil {
		return fmt.Errorf("load container error: %s", err.Error())
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
		return fmt.Errorf("signal container error: %s", err.Error())
	}

	if _, err := task.Wait(ctx); err != nil {
		return fmt.Errorf("wait task error: %s", err.Error())
	}

	newTask, err := container.NewTask(ctx, cio.NullIO)
	if err != nil {
		return fmt.Errorf("create new task of container error: %s", err.Error())
	}

	if err := newTask.Start(ctx); err != nil {
		return fmt.Errorf("start task error: %s", err.Error())
	}

	return nil
}

func (d *Client) ListId(ctx context.Context) ([]string, error) {
	containerList, err := d.client.Containers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container list error: %s", err.Error())
	}

	var idList = make([]string, len(containerList))
	for i, c := range containerList {
		idList[i] = c.ID()
	}

	return idList, nil
}

func (d *Client) CpFile(ctx context.Context, containerID, src, dst string) error {
	rootfs := fmt.Sprintf("/run/containerd/io.containerd.runtime.v2.task/k8s.io/%s/rootfs", containerID)
	_, err := os.Stat(rootfs)
	if err != nil {
		return d.CpFileOld(ctx, containerID, src, dst)
	}

	dst = fmt.Sprintf("%s%s", rootfs, dst)
	log.GetLogger(ctx).Debugf("target merged file: %s", dst)
	return base.CopyFile(src, dst)
}

func (d *Client) CpFileOld(ctx context.Context, containerID, src, dst string) error {
	task, err := d.getContainerTask(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get task of container error: %s", err.Error())
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file error: %s", err.Error())
	}
	defer srcFile.Close()

	fileInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat file error: %s", err.Error())
	}

	perm := fmt.Sprintf("%o", fileInfo.Mode().Perm())
	pId := uuid.New().String()
	var stdout, stderr bytes.Buffer

	// TODO：To be optimized, should not rely on /bin/bash in the container
	pro, err := task.Exec(ctx, pId, &specs.Process{
		Args: []string{"/bin/bash", "-c", fmt.Sprintf("touch %s && chmod %s %s && cat > %s", dst, perm, dst, dst)},
		Cwd:  "/",
	}, cio.NewCreator(append([]cio.Opt{cio.WithStreams(srcFile, &stdout, &stderr)})...))

	if err != nil {
		return fmt.Errorf("container exec error: %s", err.Error())
	}

	defer pro.Delete(ctx)
	eStatusC, err := pro.Wait(ctx)
	if err != nil {
		return fmt.Errorf("wait exec error: %s", err.Error())
	}

	if err := pro.Start(ctx); err != nil {
		return fmt.Errorf("task start error: %s", err.Error())
	}

	time.Sleep(500 * time.Millisecond)
	if err := pro.Kill(ctx, syscall.SIGTERM); err != nil {
		return fmt.Errorf("terminal process error: %s", err.Error())
	}

	eStatus := <-eStatusC
	output := stdout.String() + stderr.String()
	if eStatus.ExitCode() != termExitCode && eStatus.ExitCode() != 0 {
		return fmt.Errorf("exit code: %d, exit error: %v, msg: %s", eStatus.ExitCode(), eStatus.Error(), output)
	}

	return nil
}
