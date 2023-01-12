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

package container

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/crclient"
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	injector.Register(TargetCpu, FaultContainerRestart, func() injector.IInjector { return &RestartInjector{} })
}

type RestartInjector struct {
	injector.BaseInjector
	Args    RestartArgs
	Runtime RestartRuntime
}

type RestartArgs struct {
	WaitTime int64 `json:"wait_time"`
}

type RestartRuntime struct {
}

func (i *RestartInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *RestartInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *RestartInjector) SetDefault() {
	i.BaseInjector.SetDefault()

	if i.Args.WaitTime <= 0 {
		i.Args.WaitTime = DefaultWaitTime
	}
}

func (i *RestartInjector) SetOption(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&i.Args.WaitTime, "wait-time", "w", 0, "tolerable time-consuming(seconds) to restart the container")
}

func (i *RestartInjector) Validator(ctx context.Context) error {
	if i.Info.ContainerRuntime == "" || i.Info.ContainerId == "" {
		return fmt.Errorf("please provide container runtime and id")
	}

	return i.BaseInjector.Validator(ctx)
}

func (i *RestartInjector) Inject(ctx context.Context) error {
	client, err := crclient.GetClient(ctx, i.Info.ContainerRuntime)
	if err != nil {
		return fmt.Errorf("get %s client error: %s", i.Info.ContainerRuntime, err.Error())
	}

	var waitTime = time.Second * time.Duration(i.Args.WaitTime)
	return client.RestartContainerById(ctx, i.Info.ContainerId, &waitTime)
}

func (i *RestartInjector) Recover(ctx context.Context) error {
	return nil
}

//func (i *RestartInjector) DelayRecover(ctx context.Context, timeout int64) error {
//	return nil
//}
