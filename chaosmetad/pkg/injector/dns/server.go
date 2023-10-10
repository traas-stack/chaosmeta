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

package dns

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/cmdexec"
)

func init() {
	injector.Register(TargetDNS, FaultDNSServer, func() injector.IInjector { return &ServerInjector{} })
}

type ServerInjector struct {
	injector.BaseInjector
	Args    ServerArgs
	Runtime ServerRuntime
}

type ServerArgs struct {
	Ip   string `json:"ip"`
	Mode string `json:"mode"`
}

type ServerRuntime struct {
}

func (i *ServerInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *ServerInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *ServerInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeAdd
	}
}

func (i *ServerInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Ip, "ip", "i", "", "dns server's ip")
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("inject mode, support: %s, %s", ModeAdd, ModeDelete))
}

// Validator delete: cannot delete records that have been deleted
func (i *ServerInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Ip == "" {
		return fmt.Errorf("must provide args \"ip\"")
	}

	if i.Args.Mode != ModeAdd && i.Args.Mode != ModeDelete {
		return fmt.Errorf("args \"mode\" only support: %s, %s", ModeAdd, ModeDelete)
	}

	return nil
}

// Inject add: insert in first line
func (i *ServerInjector) Inject(ctx context.Context) error {
	var cmd string
	if i.Args.Mode == ModeAdd {
		cmd = getServerAddInjectCmd(i.Info.Uid, i.Args.Ip)
	} else {
		cmd = getServerDeleteInjectCmd(i.Info.Uid, i.Args.Ip)
	}

	_, err := cmdexec.ExecCommon(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, cmd)
	return err
}

func (i *ServerInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	var cmd string
	if i.Args.Mode == ModeAdd {
		cmd = getServerAddRecoverCmd(i.Info.Uid)
	} else {
		cmd = getServerDeleteRecoverCmd(i.Info.Uid)
	}

	_, err := cmdexec.ExecCommon(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, cmd)
	return err
}

func getServerAddInjectCmd(uid, ip string) string {
	return fmt.Sprintf("sed '1s/^/nameserver %s %s\\n/' %s > %s && cat %s > %s && rm -rf %s", ip, getFlag(uid, ModeAdd), ConfServer, ConfServerBak, ConfServerBak, ConfServer, ConfServerBak)
}

func getServerDeleteInjectCmd(uid, ip string) string {
	return fmt.Sprintf("sed '/%s/s/^/%s/' %s > %s && cat %s > %s && rm -rf %s", ip, getFlag(uid, ModeDelete), ConfServer, ConfServerBak, ConfServerBak, ConfServer, ConfServerBak)
}

func getServerAddRecoverCmd(uid string) string {
	return fmt.Sprintf("sed '/%s/d' %s > %s && cat %s > %s && rm -rf %s", getFlag(uid, ModeAdd), ConfServer, ConfServerBak, ConfServerBak, ConfServer, ConfServerBak)
}

func getServerDeleteRecoverCmd(uid string) string {
	return fmt.Sprintf("sed 's/%s//' %s > %s && cat %s > %s && rm -rf %s", getFlag(uid, ModeDelete), ConfServer, ConfServerBak, ConfServerBak, ConfServer, ConfServerBak)
}
