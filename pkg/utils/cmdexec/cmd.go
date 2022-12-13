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

package cmdexec

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/crclient"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/namespace"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	InjectCheckInterval = time.Millisecond * 200
	cgroupWaitInterval  = time.Millisecond * 200
)

const (
	ExecWait   = "wait"
	ExecNormal = "normal"
)

//type CmdExecutor struct {
//	ContainerId      string
//	ContainerRuntime string
//	ContainerNs      []string
//	//Method           string
//}

func StartSleepRecover(ctx context.Context, sleepTime int64, uid string) error {
	return StartBashCmd(ctx, utils.GetSleepRecoverCmd(sleepTime, uid))
}

func waitProExec(ctx context.Context, stdout, stderr *bytes.Buffer) (err error) {
	var msg, timer = "", time.NewTimer(InjectCheckInterval)
	for {
		<-timer.C
		if stderr.String() != "" || stdout.String() != "" {
			msg = stdout.String() + stderr.String()
			break
		}
		timer.Reset(InjectCheckInterval)
	}

	log.GetLogger(ctx).Debugf(msg)

	if strings.Index(msg, "error") >= 0 {
		return fmt.Errorf("inject error: %s", msg)
	}

	if strings.Index(msg, "[success]") >= 0 {
		return nil
	}

	return fmt.Errorf("unexpected output: %s", msg)
}

func SupportCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}

	return true
}

func RunBashCmdWithOutput(ctx context.Context, cmd string) ([]byte, error) {
	log.GetLogger(ctx).Debugf("run cmd with output: %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
}

func RunBashCmdWithoutOutput(ctx context.Context, cmd string) error {
	log.GetLogger(ctx).Debugf("run cmd: %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).Run()
}

func StartBashCmd(ctx context.Context, cmd string) error {
	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).Start()
}

func StartBashCmdWithPid(ctx context.Context, cmd string) (int, error) {
	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)
	c := exec.Command("/bin/bash", "-c", cmd)
	if err := c.Start(); err != nil {
		return utils.NoPid, err
	}

	return c.Process.Pid, nil
}

func StartBashCmdAndWaitPid(ctx context.Context, cmd string) (int, error) {
	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)

	c := exec.Command("/bin/bash", "-c", cmd)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Start(); err != nil {
		return utils.NoPid, fmt.Errorf("cmd start error: %s", err.Error())
	}

	if err := waitProExec(ctx, &stdout, &stderr); err != nil {
		return c.Process.Pid, fmt.Errorf("wait process exec error: %s", err.Error())
	}

	return c.Process.Pid, nil
}

func StartBashCmdAndWaitByUser(ctx context.Context, cmd, user string) error {
	log.GetLogger(ctx).Debugf("user: %s, start cmd: %s", user, cmd)

	c := exec.Command("runuser", "-l", user, "-c", cmd)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Start(); err != nil {
		return fmt.Errorf("cmd start error: %s", err.Error())
	}

	if err := waitProExec(ctx, &stdout, &stderr); err != nil {
		return fmt.Errorf("wait process exec error: %s", err.Error())
	}

	return nil
}

// finish: false[wait success], true[finish and get all output]

func ExecContainer(ctx context.Context, cr, containerID string, namespaces []string, cmd string, finish bool) (string, error) {
	logger := log.GetLogger(ctx)

	// get container's init process
	client, err := crclient.GetClient(ctx, cr)
	if err != nil {
		return "", fmt.Errorf("get %s client error: %s", cr, err.Error())
	}

	targetPid, err := client.GetPidById(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("get pid of container[%s]'s init process error: %s", containerID, err.Error())
	}

	// exec ns
	c := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s -t %d %s -c \"%s\"",
		utils.GetToolPath(namespace.ExecnsKey), targetPid, namespace.GetNsOption(namespaces), cmd))

	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	logger.Debugf("container exec cmd: %s", c.Args)
	if err := c.Start(); err != nil {
		return "", fmt.Errorf("start process error: %s", err.Error())
	}

	// set cgroup for new process
	if err := addToProCgroup(c.Process.Pid, targetPid); err != nil {
		if err := c.Process.Kill(); err != nil {
			logger.Warnf("undo: kill container exec process[%d] error: %s", c.Process.Pid, err.Error())
		}

		return "", fmt.Errorf("add process[%d] to container[%d] cgroup error: %s", c.Process.Pid, targetPid, err.Error())
	}

	// signal continue
	time.Sleep(cgroupWaitInterval)
	if err := c.Process.Signal(syscall.SIGCONT); err != nil {
		return "", err
	}

	if finish {
		if err := c.Wait(); err != nil {
			return "", fmt.Errorf("wait process error: %s", err.Error())
		}

		if strings.Index(stdout.String()+stderr.String(), "error") >= 0 {
			return "", fmt.Errorf(stdout.String() + stderr.String())
		} else {
			return stdout.String() + stderr.String(), nil
		}

	} else {
		return "", waitProExec(ctx, &stdout, &stderr)
	}
}
