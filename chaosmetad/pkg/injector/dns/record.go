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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/namespace"
)

func init() {
	injector.Register(TargetDNS, FaultDNSRecord, func() injector.IInjector { return &RecordInjector{} })
}

type RecordInjector struct {
	injector.BaseInjector
	Args    RecordArgs
	Runtime RecordRuntime
}

type RecordArgs struct {
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
	Mode   string `json:"mode"`
}

type RecordRuntime struct {
}

func (i *RecordInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *RecordInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *RecordInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.Mode == "" {
		i.Args.Mode = ModeAdd
	}
}

func (i *RecordInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&i.Args.Domain, "domain", "d", "", "dns record's domain")
	cmd.Flags().StringVarP(&i.Args.Ip, "ip", "i", "", "dns record's ip")
	cmd.Flags().StringVarP(&i.Args.Mode, "mode", "m", "", fmt.Sprintf("inject mode, support: %s, %s", ModeAdd, ModeDelete))
}

// Validator delete: cannot delete records that have been deleted
func (i *RecordInjector) Validator(ctx context.Context) error {
	if err := i.BaseInjector.Validator(ctx); err != nil {
		return err
	}

	if i.Args.Mode == ModeAdd {
		if i.Args.Domain == "" || i.Args.Ip == "" {
			return fmt.Errorf("must provide args \"domain\" and \"ip\" in mode \"%s\"", ModeAdd)
		}
	} else if i.Args.Mode == ModeDelete {
		if i.Args.Domain == "" {
			return fmt.Errorf("must provide args \"domain\" in mode \"%s\"", ModeDelete)
		}
	} else {
		return fmt.Errorf("args \"mode\" only support: %s, %s", ModeAdd, ModeDelete)
	}

	return nil
}

// Inject add: insert in first line
func (i *RecordInjector) Inject(ctx context.Context) error {
	var cmd string
	if i.Args.Mode == ModeAdd {
		cmd = getRecordAddInjectCmd(i.Info.Uid, i.Args.Domain, i.Args.Ip)
	} else {
		cmd = getRecordDeleteInjectCmd(i.Info.Uid, i.Args.Domain)
	}

	_, err := cmdexec.ExecCommonWithNS(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, cmd, []string{namespace.MNT})
	return err
}

func (i *RecordInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	var cmd string
	if i.Args.Mode == ModeAdd {
		cmd = getRecordAddRecoverCmd(i.Info.Uid)
	} else {
		cmd = getRecordDeleteRecoverCmd(i.Info.Uid)
	}

	_, err := cmdexec.ExecCommonWithNS(ctx, i.Info.ContainerRuntime, i.Info.ContainerId, cmd, []string{namespace.MNT})
	return err
}

func getRecordAddInjectCmd(uid, domain, ip string) string {
	return fmt.Sprintf("sed '1s/^/%s %s %s\\n/' %s > %s && cat %s > %s && rm -rf %s", ip, domain, getFlag(uid, ModeAdd), ConfRecord, ConfRecordBak, ConfRecordBak, ConfRecord, ConfRecordBak)
}

func getRecordDeleteInjectCmd(uid, domain string) string {
	return fmt.Sprintf("sed '/%s/s/^/%s/' %s > %s && cat %s > %s && rm -rf %s", domain, getFlag(uid, ModeDelete), ConfRecord, ConfRecordBak, ConfRecordBak, ConfRecord, ConfRecordBak)
}

func getRecordAddRecoverCmd(uid string) string {
	return fmt.Sprintf("sed '/%s/d' %s > %s && cat %s > %s && rm -rf %s", getFlag(uid, ModeAdd), ConfRecord, ConfRecordBak, ConfRecordBak, ConfRecord, ConfRecordBak)
}

func getRecordDeleteRecoverCmd(uid string) string {
	return fmt.Sprintf("sed 's/%s//' %s > %s && cat %s > %s && rm -rf %s", getFlag(uid, ModeDelete), ConfRecord, ConfRecordBak, ConfRecordBak, ConfRecord, ConfRecordBak)
}

func getFlag(uid, mode string) string {
	return fmt.Sprintf("# ChaosMeta-%s-%s ", mode, uid)
}
