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

package network

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/net"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/process"
	"github.com/spf13/cobra"
)

func init() {
	injector.Register(TargetNetwork, FaultOccupy, func() injector.IInjector { return &OccupyInjector{} })
}

type OccupyInjector struct {
	injector.BaseInjector
	Args    OccupyArgs
	Runtime OccupyRuntime
}

type OccupyArgs struct {
	Port       int    `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Force      bool   `json:"force,omitempty"`
	RecoverCmd string `json:"recover_cmd,omitempty"`
}

type OccupyRuntime struct {
	//Pid int `json:"pid,omitempty"`
}

func (i *OccupyInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *OccupyInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *OccupyInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Protocol == "" {
		i.Args.Protocol = net.ProtocolTCP
	}
}

func (i *OccupyInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)
	cmd.Flags().IntVarP(&i.Args.Port, "port", "p", 0, "target port")
	cmd.Flags().StringVarP(&i.Args.Protocol, "protocol", "P", "",
		fmt.Sprintf("target protocol, support: %s、%s、%s、%s（default %s）",
			net.ProtocolTCP, net.ProtocolUDP, net.ProtocolTCP6, net.ProtocolUDP6, net.ProtocolTCP))
	cmd.Flags().BoolVarP(&i.Args.Force, "force", "f", false, "if kill the process which occupied target port")
	cmd.Flags().StringVarP(&i.Args.RecoverCmd, "recover-cmd", "r", "", "execute in recover stage")
}

func (i *OccupyInjector) Validator(ctx context.Context) error {
	if i.Args.Port <= 0 {
		return fmt.Errorf("\"port\" must larger than 0")
	}

	if i.Args.Protocol != net.ProtocolTCP && i.Args.Protocol != net.ProtocolUDP && i.Args.Protocol != net.ProtocolTCP6 && i.Args.Protocol != net.ProtocolUDP6 {
		return fmt.Errorf("\"protocol\" is not support %s", i.Args.Protocol)
	}

	return i.BaseInjector.Validator(ctx)
}

func (i *OccupyInjector) Inject(ctx context.Context) error {
	pid, err := net.GetPidByPort(ctx, i.Args.Port, i.Args.Protocol)
	if err != nil {
		return fmt.Errorf("get pid by port[%d] error: %s", i.Args.Port, err.Error())
	}

	if pid != utils.NoPid {
		if i.Args.Force {
			if err := process.KillPidWithSignal(ctx, pid, process.SIGKILL); err != nil {
				return fmt.Errorf("kill occupied process[%d] error: %s", pid, err.Error())
			}
		} else {
			return fmt.Errorf("port[%d] is occupied by process[%d], if want to force occupy, please add force args", i.Args.Port, pid)
		}
	}

	var timeout int64
	if i.Info.Timeout != "" {
		timeout, err = utils.GetTimeSecond(i.Info.Timeout)
	}

	_, err = cmdexec.StartBashCmdAndWaitPid(ctx, fmt.Sprintf("%s %s %d %s %d", utils.GetToolPath(OccupyKey), i.Info.Uid, i.Args.Port, i.Args.Protocol, timeout))
	if err != nil {
		return fmt.Errorf("start cmd error: %s", err.Error())
	}

	return nil
}

func (i *OccupyInjector) DelayRecover(ctx context.Context, timeout int64) error {
	return nil
}

func (i *OccupyInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	if err:= process.CheckExistAndKillByKey(ctx, fmt.Sprintf("%s %s", OccupyKey, i.Info.Uid));err != nil {
		return err
	}

	if i.Args.RecoverCmd != "" {
		return cmdexec.StartBashCmd(ctx, i.Args.RecoverCmd)
	}

	return nil
}
