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

package version

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/version"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "description of version",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintVersion(utils.GetCtxWithTraceId(context.Background(), utils.TraceId))
		},
	}
}
