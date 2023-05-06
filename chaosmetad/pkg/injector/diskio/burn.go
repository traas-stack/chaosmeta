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

package diskio

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
)

//TODO: It needs to be stated in the document that the target directory has at least 1G disk space remaining

func init() {
	injector.Register(TargetDiskIO, FaultDiskIOBurn, func() injector.IInjector { return &BurnInjector{} })
}

type BurnInjector struct {
	injector.BaseInjector
	Args    BurnArgs
	Runtime BurnRuntime
}

type BurnArgs struct {
	Mode  string `json:"mode"`
	Block string `json:"block"`
	Dir   string `json:"dir"`
}

type BurnRuntime struct {
}

func (i *BurnInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *BurnInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *BurnInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Dir == "" {
		i.Args.Dir = DefaultDir
	}

	if i.Args.Mode == "" {
		i.Args.Mode = ModeRead
	}

	if i.Args.Block == "" {
		i.Args.Block = DefaultBlockSize
	}
}

func (i *BurnInjector) SetOption(cmd *cobra.Command) {
	//// i.BaseInjector.SetOption(cmd)
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("disk IO mode, support: %s、%s（default %s）", ModeRead, ModeWrite, ModeRead))
	cmd.Flags().StringVarP(&i.Args.Block, "block", "b", "", fmt.Sprintf("disk IO block size（default %s）, support unit: KB/MB（default KB）", DefaultBlockSize))
	cmd.Flags().StringVarP(&i.Args.Dir, "dir", "d", "", fmt.Sprintf("disk IO burn directory（default %s）", DefaultDir))
}

func (i *BurnInjector) getCmdExecutor(method, args string) *cmdexec.CmdExecutor {
	return &cmdexec.CmdExecutor{
		ContainerId:      i.Info.ContainerId,
		ContainerRuntime: i.Info.ContainerRuntime,
		ContainerNs:      []string{namespace.MNT, namespace.PID},
		ToolKey:          DiskIOExec,
		Method:           method,
		Fault:            FaultDiskIOBurn,
		Args:             args,
	}
}

func (i *BurnInjector) getFileName() string {
	return fmt.Sprintf("%s/%s_%s", i.Args.Dir, DiskIOBurnKey, i.Info.Uid)
}

func (i *BurnInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Dir == "" {
		return fmt.Errorf("\"dir\" is empty")
	}

	if i.Args.Mode != ModeRead && i.Args.Mode != ModeWrite {
		return fmt.Errorf("\"mode\" not support %s, only support: %s、%s", i.Args.Mode, ModeRead, ModeWrite)
	}

	kbytes, _, err := utils.GetBlockKbytes(i.Args.Block)
	if err != nil {
		return fmt.Errorf("\"block\"[%s] is invalid: %s", i.Args.Block, err.Error())
	}

	if kbytes <= 0 || kbytes > MaxBlockK {
		return fmt.Errorf("\"block\"[%s] value must be in (0, 1G]", i.Args.Block)
	}

	return i.getCmdExecutor(utils.MethodValidator, i.Args.Dir).ExecTool(ctx)
}

func (i *BurnInjector) Inject(ctx context.Context) error {
	logger := log.GetLogger(ctx)
	var timeout int64
	if i.Info.Timeout != "" {
		timeout, _ = utils.GetTimeSecond(i.Info.Timeout)
	}

	blockK, stdStr, _ := utils.GetBlockKbytes(i.Args.Block)
	count := MaxBlockK / blockK

	toolPath := utils.GetToolPath(DiskIOBurnKey)
	args := fmt.Sprintf("%s %s %s %s %d %s %d", i.Info.Uid, i.getFileName(), i.Args.Mode, stdStr, count, FlagDirect, timeout)

	//cmd := fmt.Sprintf("%s %s %s %s %s %d %s %d", utils.GetToolPath(DiskIOBurnKey), i.Info.Uid, i.getFileName(), i.Args.Mode, stdStr, count, FlagDirect, timeout)

	if i.Info.ContainerRuntime != "" {
		localPath := toolPath
		toolPath = utils.GetContainerPath(DiskIOBurnKey)
		if err := cmdexec.CpContainerFile(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, localPath, toolPath); err != nil {
			return fmt.Errorf("container cp from [%s] to [%s] error: %s", localPath, toolPath, err.Error())
		}
	}

	cmd := fmt.Sprintf("%s %s", toolPath, args)
	if err := i.getCmdExecutor("", "").StartCmdAndWait(ctx, cmd); err != nil {
		if err := i.Recover(ctx); err != nil {
			logger.Warnf("undo error: %s", err.Error())
		}

		return err
	}

	return nil
}

func (i *BurnInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return i.getCmdExecutor(utils.MethodRecover, fmt.Sprintf("%s %s", i.Info.Uid, i.Args.Dir)).ExecTool(ctx)
}
