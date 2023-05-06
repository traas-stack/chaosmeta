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

package recover

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
)

func NewRecoverCommand() *cobra.Command {
	recoverCmd := &cobra.Command{
		Use:   "recover",
		Short: "experiment recover command",
		Long:  "experiment recover command, usage: recover [uid]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := utils.GetCtxWithTraceId(context.Background(), utils.TraceId)
			if len(args) != 1 {
				errutil.SolveErr(ctx, errutil.BadArgsErr, fmt.Sprintf("please add target experiment's uid, eg: recover [uid]"))
			}

			code, msg := injector.ProcessRecover(ctx, args[0])
			errutil.SolveErr(ctx, code, msg)
		},
	}

	return recoverCmd
}
