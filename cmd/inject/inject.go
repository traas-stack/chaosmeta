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

package inject

import (
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/container"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/cpu"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/disk"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/diskio"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/file"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/kernel"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/mem"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/network"
	_ "github.com/ChaosMetaverse/chaosmetad/pkg/injector/process"
	"github.com/spf13/cobra"
)

// NewInjectCommand injectCmd represents the inject command
func NewInjectCommand() *cobra.Command {
	var injectCmd = &cobra.Command{
		Use:   "inject",
		Short: "experiment create command",
	}

	targets := injector.GetTargets()

	var args = &injector.BaseInfo{}
	injectCmd.PersistentFlags().StringVarP(&args.Timeout, "timeout", "t", "", "experiment's duration（default 0, means need to stop manually）")
	injectCmd.PersistentFlags().StringVar(&args.Creator, "creator", "", "experiment's creator（default the cmd exec user）")

	injectCmd.PersistentFlags().StringVar(&args.ContainerRuntime, "container-runtime", "", "if attack a container of local host, need to provide the container runtime of target container")
	injectCmd.PersistentFlags().StringVar(&args.ContainerId, "container-id", "", "if attack a container of local host, need to provide the container id of target container")

	//var args = make([]string, 2)
	//injectCmd.PersistentFlags().StringVarP(&args[0], "timeout", "t", "", "experiment's duration（default 0, means need to stop manually）")
	//injectCmd.PersistentFlags().StringVar(&args[1], "creator", "", "experiment's creator（default the cmd exec user）")

	for _, target := range targets {
		targetCmd := injector.NewCmdByTarget(target, args)
		injectCmd.AddCommand(targetCmd)
	}

	return injectCmd
}
