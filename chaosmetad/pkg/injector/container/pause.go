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
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/crclient"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
)

func init() {
	injector.Register(TargetContainer, FaultContainerPause, func() injector.IInjector { return &PauseInjector{} })
}

type PauseInjector struct {
	injector.BaseInjector
	Args    PauseArgs
	Runtime PauseRuntime
}

type PauseArgs struct {
}

type PauseRuntime struct {
}

func (i *PauseInjector) GetArgs() interface{} {
	return &i.Args
}

func (i *PauseInjector) GetRuntime() interface{} {
	return &i.Runtime
}

func (i *PauseInjector) SetOption(cmd *cobra.Command) {
}

func (i *PauseInjector) Validator(ctx context.Context) error {
	if i.Info.ContainerRuntime == "" || i.Info.ContainerId == "" {
		return fmt.Errorf("please provide container runtime and id")
	}

	return i.BaseInjector.Validator(ctx)
}

func (i *PauseInjector) Inject(ctx context.Context) error {
	client, err := crclient.GetClient(ctx, i.Info.ContainerRuntime)
	if err != nil {
		return fmt.Errorf("get %s client error: %s", i.Info.ContainerRuntime, err.Error())
	}

	return client.PauseContainerById(ctx, i.Info.ContainerId)
}

func (i *PauseInjector) Recover(ctx context.Context) error {
	if i.BaseInjector.Recover(ctx) == nil {
		return nil
	}

	client, err := crclient.GetClient(ctx, i.Info.ContainerRuntime)
	if err != nil {
		return fmt.Errorf("get %s client error: %s", i.Info.ContainerRuntime, err.Error())
	}

	return client.UnPauseContainerById(ctx, i.Info.ContainerId)
}

//func (i *PauseInjector) DelayRecover(ctx context.Context, timeout int64) error {
//	return nil
//}
