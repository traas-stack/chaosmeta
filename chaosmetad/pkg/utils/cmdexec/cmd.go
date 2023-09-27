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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/crclient"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/containercgroup"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	InjectCheckInterval = time.Millisecond * 200
	cgroupWaitInterval  = time.Millisecond * 200

	ExecWait  = "wait"
	ExecStart = "start"
	ExecRun   = "run"
)

type CmdExecutor struct {
	ContainerId      string
	ContainerRuntime string
	ContainerNs      []string
	ToolKey          string
	Method           string
	Fault            string
	Args             string
}

func (e *CmdExecutor) GetTargetPid(ctx context.Context) (int, error) {
	if e.ContainerRuntime != "" {
		client, err := crclient.GetClient(ctx, e.ContainerRuntime)
		if err != nil {
			return -1, fmt.Errorf("get %s client error: %s", e.ContainerRuntime, err.Error())
		}

		return client.GetPidById(ctx, e.ContainerId)
	} else {
		return -1, nil
	}
}

func (e *CmdExecutor) StartCmdAndWait(ctx context.Context, cmd string) error {
	var err error
	if e.ContainerRuntime != "" {
		_, err = ExecContainer(ctx, e.ContainerRuntime, e.ContainerId, e.ContainerNs, cmd, ExecWait)
	} else {
		_, err = StartBashCmdAndWaitPid(ctx, cmd, 0)
	}
	return err
}

func (e *CmdExecutor) StartCmd(ctx context.Context, cmd string) error {
	var err error
	if e.ContainerRuntime != "" {
		_, err = ExecContainer(ctx, e.ContainerRuntime, e.ContainerId, e.ContainerNs, cmd, ExecStart)
	} else {
		err = StartBashCmd(ctx, cmd)
	}
	return err
}

// Exec cmd should not block
func (e *CmdExecutor) Exec(ctx context.Context, cmd string) (string, error) {
	if e.ContainerRuntime != "" {
		return ExecContainerRaw(ctx, e.ContainerRuntime, e.ContainerId, cmd)
	} else {
		return RunBashCmdWithOutput(ctx, cmd)
	}
}

func (e *CmdExecutor) ExecTool(ctx context.Context) error {
	logger, commonArgs := log.GetLogger(ctx), fmt.Sprintf("%s %s %s %s", e.Method, e.Fault, log.Level, e.Args)

	if e.ContainerRuntime != "" {
		var execTool = utils.GetToolPath(e.ToolKey)
		if utils.StrListContain(e.ContainerNs, namespace.MNT) {
			srcTool := execTool
			execTool = utils.GetContainerPath(e.ToolKey)
			if err := CpContainerFile(ctx, e.ContainerRuntime, e.ContainerId, srcTool, execTool); err != nil {
				return fmt.Errorf("cp exec tool to container[%s] error: %s", e.ContainerId, err.Error())
			}
		}

		re, err := ExecContainer(ctx, e.ContainerRuntime, e.ContainerId, e.ContainerNs, fmt.Sprintf("%s %s", execTool, commonArgs), ExecRun)
		logger.Debugf(re)
		if err != nil {
			return fmt.Errorf("exec in container error: %s", err.Error())
		}
	} else {
		re, err := RunBashCmdWithOutput(ctx, fmt.Sprintf("%s %s", utils.GetToolPath(e.ToolKey), commonArgs))
		logger.Debugf(re)
		if err != nil {
			return err
		}
	}

	return nil
}

func CpContainerFile(ctx context.Context, cr, containerID, src, dst string) error {
	log.GetLogger(ctx).Debugf("cp from %s to %s in %s", src, dst, containerID)
	client, err := crclient.GetClient(ctx, cr)
	if err != nil {
		return fmt.Errorf("get %s client error: %s", cr, err.Error())
	}

	return client.CpFile(ctx, containerID, src, dst)
}

func StartSleepRecover(ctx context.Context, sleepTime int64, uid string) error {
	return StartBashCmd(ctx, utils.GetSleepRecoverCmd(sleepTime, uid))
}

func waitProExec(ctx context.Context, stdout, stderr *bytes.Buffer, timeoutSec int) (err error) {
	var msg, timer = "", time.NewTimer(InjectCheckInterval)
	var startTime = time.Now()

	for {
		<-timer.C
		if stderr.String() != "" || stdout.String() != "" {
			msg = stdout.String() + stderr.String()
			break
		}

		if timeoutSec > 0 && time.Now().After(startTime.Add(time.Second*time.Duration(timeoutSec))) {
			break
		}

		timer.Reset(InjectCheckInterval)
	}

	log.GetLogger(ctx).Debugf(msg)

	if strings.Index(msg, "error") >= 0 || strings.Index(msg, "Error") >= 0 {
		return fmt.Errorf("inject error: %s", msg)
	}

	if timeoutSec <= 0 {
		if strings.Index(msg, "[success]") >= 0 {
			return nil
		}
	} else {
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

func RunBashCmdWithOutput(ctx context.Context, cmd string) (string, error) {
	log.GetLogger(ctx).Debugf("run cmd with output: %s", cmd)
	c := exec.Command("/bin/bash", "-c", cmd)

	reByte, err := c.CombinedOutput()
	re := string(reByte)
	log.GetLogger(ctx).Debugf("output: %s, err: %v", re, err)
	if err != nil {
		if c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() == errutil.ExpectedErr {
			return "", fmt.Errorf("output: %s, error: %s", re, err.Error())
		}

		if c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus() == errutil.TestFileErr {
			return "", fmt.Errorf("exit code: %d, output: %s, error: %s", errutil.TestFileErr, re, err.Error())
		}

		return "", fmt.Errorf("error: %s, output: %s", err.Error(), re)
	}

	return re, nil
}

func RunBashCmdWithoutOutput(ctx context.Context, cmd string) error {
	log.GetLogger(ctx).Debugf("run cmd: %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).Run()
}

func StartBashCmd(ctx context.Context, cmd string) error {
	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).Start()
}

//func StartBashCmdWithPid(ctx context.Context, cmd string) (int, error) {
//	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)
//	c := exec.Command("/bin/bash", "-c", cmd)
//	if err := c.Start(); err != nil {
//		return utils.NoPid, err
//	}
//
//	return c.Process.Pid, nil
//}

func StartBashCmdAndWaitPid(ctx context.Context, cmd string, timeoutSec int) (int, error) {
	log.GetLogger(ctx).Debugf("start cmd: %s", cmd)

	c := exec.Command("/bin/bash", "-c", cmd)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Start(); err != nil {
		return utils.NoPid, fmt.Errorf("cmd start error: %s", err.Error())
	}

	if err := waitProExec(ctx, &stdout, &stderr, timeoutSec); err != nil {
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

	if err := waitProExec(ctx, &stdout, &stderr, 0); err != nil {
		return fmt.Errorf("wait process exec error: %s", err.Error())
	}

	return nil
}

// finish: false[wait success], true[finish and get all output]

func ExecContainer(ctx context.Context, cr, containerID string, namespaces []string, cmd string, method string) (string, error) {
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
	if err := containercgroup.AddToProCgroup(c.Process.Pid, targetPid); err != nil {
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

	// solve return
	switch method {
	case ExecWait:
		return "", waitProExec(ctx, &stdout, &stderr, 0)
	case ExecRun:
		if err := c.Wait(); err != nil {
			return "", fmt.Errorf("wait process error: %s", err.Error())
		}

		combinedOutput := stdout.String() + stderr.String()
		// TODO: not use exit code judge task if success?
		if strings.Index(combinedOutput, "[error]") >= 0 {
			return "", fmt.Errorf(combinedOutput)
		} else {
			return combinedOutput, nil
		}
	case ExecStart:
		return "", nil
	default:
		return "", fmt.Errorf("unknown exec method")
	}
}

func ExecContainerRaw(ctx context.Context, cr, cId, cmd string) (string, error) {
	client, err := crclient.GetClient(ctx, cr)
	if err != nil {
		return "", fmt.Errorf("get %s client error: %s", cr, err.Error())
	}

	log.GetLogger(ctx).Debugf("container: %s, exec cmd: %s", cId, cmd)
	re, err := client.Exec(ctx, cId, cmd)
	log.GetLogger(ctx).Debugf("container: %s, output: %s, err: %v", cId, re, err)
	return re, err
}

func ExecCommon(ctx context.Context, cr, cId, cmd string) (string, error) {
	if cr == "" {
		return RunBashCmdWithOutput(ctx, cmd)
	} else {
		// TODO: need to transfer to ExecContainer?
		return ExecContainerRaw(ctx, cr, cId, cmd)
	}
}

func ExecBackGroundCommon(ctx context.Context, cr, cId, cmd string) error {
	if cr == "" {
		return StartBashCmd(ctx, cmd)
	} else {
		_, err := ExecContainer(ctx, cr, cId, []string{namespace.MNT, namespace.IPC, namespace.NET, namespace.PID, namespace.UTS}, cmd, ExecStart)
		return err
	}
}
