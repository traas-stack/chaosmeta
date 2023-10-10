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
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/net"
)

func init() {
	injector.Register(TargetNetwork, FaultCorrupt, func() injector.IInjector { return &CorruptInjector{} })
}

type CorruptInjector struct {
	injector.BaseInjector
	Args    CorruptArgs
	Runtime CorruptRuntime
}

type CorruptArgs struct {
	Interface string `json:"interface"`
	Percent   int    `json:"percent"`
	Direction string `json:"direction"`
	Mode      string `json:"mode"`
	SrcIp     string `json:"src_ip,omitempty"`
	DstIp     string `json:"dst_ip,omitempty"`
	SrcPort   string `json:"src_port,omitempty"`
	DstPort   string `json:"dst_port,omitempty"`
	Force     bool   `json:"force,omitempty"`
}

type CorruptRuntime struct{}

func (i *CorruptInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *CorruptInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *CorruptInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Direction == "" {
		i.Args.Direction = DirectionOut
	}

	if i.Args.Mode == "" {
		i.Args.Mode = net.ModeNormal
	}
}

func (i *CorruptInjector) SetOption(cmd *cobra.Command) {
	// i.BaseInjector.SetOption(cmd)

	cmd.Flags().IntVarP(&i.Args.Percent, "percent", "p", 0, "packets corrupt percent, an integer in (0,100] without \"%\", eg: \"30\" means \"30%\"")

	cmd.Flags().StringVarP(&i.Args.Direction, "direction", "d", "", fmt.Sprintf("flow direction to inject, support: %s（default %s）", DirectionOut, DirectionOut))
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("inject mode, support: %s（default）、%s(means white list mode)", net.ModeNormal, net.ModeExclude))
	cmd.Flags().BoolVarP(&i.Args.Force, "force", "f", false, "force will overwrite the network rule if old rule exist")

	cmd.Flags().StringVarP(&i.Args.Interface, "interface", "i", "", "filter condition: network interface. eg: lo")
	cmd.Flags().StringVar(&i.Args.SrcIp, "src-ip", "", "filter condition: source ip. eg: 10.10.0.0/16,192.168.2.5,192.168.1.0/24")
	cmd.Flags().StringVar(&i.Args.DstIp, "dst-ip", "", "filter condition: destination ip. eg: 10.10.0.0/16,192.168.2.5,192.168.1.0/24")
	cmd.Flags().StringVar(&i.Args.SrcPort, "src-port", "", "filter condition: source port. eg: 8080,9090,12000/8")
	cmd.Flags().StringVar(&i.Args.DstPort, "dst-port", "", "filter condition: destination port. eg: 8080,9090,12000/8")
}

// Validator Only one tc network failure can be executed at the same time
func (i *CorruptInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}
	if i.Args.Percent <= 0 {
		return fmt.Errorf("\"percent\" must larger than 0")
	}

	if !cmdexec.SupportCmd("tc") {
		return fmt.Errorf("not support command \"tc\"")
	}

	if i.Args.Interface == "" {
		return fmt.Errorf("\"interface\" is empty")
	}

	if i.Args.Direction != DirectionOut {
		return fmt.Errorf("\"direction\" only support: %s", DirectionOut)
	}

	if i.Args.Mode != net.ModeNormal && i.Args.Mode != net.ModeExclude {
		return fmt.Errorf("\"mode\" is not support: %s, only support: %s, %s", i.Args.Mode, net.ModeNormal, net.ModeExclude)
	}

	if i.Args.SrcIp != "" {
		if _, err := net.GetValidIPList(i.Args.SrcIp, true); err != nil {
			return fmt.Errorf("\"src-ip\"[%s] is invalid: %s", i.Args.SrcIp, err.Error())
		}
	}

	if i.Args.DstIp != "" {
		if _, err := net.GetValidIPList(i.Args.DstIp, true); err != nil {
			return fmt.Errorf("\"dst-ip\"[%s] is invalid: %s", i.Args.DstIp, err.Error())
		}
	}

	if i.Args.SrcPort != "" {
		if _, err := net.GetValidPortList(i.Args.SrcPort); err != nil {
			return fmt.Errorf("\"src-port\"[%s] is invalid: %s", i.Args.SrcPort, err.Error())
		}
	}

	if i.Args.DstPort != "" {
		if _, err := net.GetValidPortList(i.Args.DstPort); err != nil {
			return fmt.Errorf("\"dst-port\"[%s] is invalid: %s", i.Args.DstPort, err.Error())
		}
	}

	exist, err := net.ExistTCRootQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface)
	if err != nil {
		return fmt.Errorf("check tc rule error: %s", err.Error())
	}

	if exist && !i.Args.Force {
		return fmt.Errorf("has other tc root rule, if want to force to execute, please provide [-f] or [--force] args")
	}

	return nil
}

func (i *CorruptInjector) Inject(ctx context.Context) error {
	if i.Args.Force {
		exist, _ := net.ExistTCRootQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface)
		if exist {
			if err := net.ClearTcRule(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface); err != nil {
				return fmt.Errorf("reset tc rule for %s error: %s", i.Args.Interface, err.Error())
			}
		}
	}

	if i.Args.SrcIp == "" && i.Args.DstIp == "" && i.Args.SrcPort == "" && i.Args.DstPort == "" {
		return net.AddNetemQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, "", FaultCorrupt, fmt.Sprintf("%d", i.Args.Percent))
	}

	if err := net.AddPrioQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, "", "1:"); err != nil {
		return fmt.Errorf("add root prio qdisc for %s error: %s", i.Args.Interface, err.Error())
	}

	if i.Args.Mode == net.ModeNormal {
		parent := "1:4"
		if err := net.AddNetemQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, parent, FaultCorrupt, fmt.Sprintf("%d", i.Args.Percent)); err != nil {
			return undoTcWithErr(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, fmt.Sprintf("add parent %s netem qdisc for %s error: %s", parent, i.Args.Interface, err.Error()))
		}
	} else {
		for subIndex := 1; subIndex < 4; subIndex++ {
			parent := fmt.Sprintf("1:%d", subIndex)
			if err := net.AddNetemQdisc(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, parent, FaultCorrupt, fmt.Sprintf("%d", i.Args.Percent)); err != nil {
				return undoTcWithErr(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, fmt.Sprintf("add parent %s netem qdisc for %s error: %s", parent, i.Args.Interface, err.Error()))
			}
		}
	}

	if err := net.AddFilter(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, "1:4", i.Args.SrcIp, i.Args.DstIp, i.Args.SrcPort, i.Args.DstPort); err != nil {
		return undoTcWithErr(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface, fmt.Sprintf("add filter for %s error: %s", i.Args.Interface, err.Error()))
	}

	return nil
}

func (i *CorruptInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	return execRecover(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, i.Args.Interface)
}
