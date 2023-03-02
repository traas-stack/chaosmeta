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

package query

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/query"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/spf13/cobra"
)

// NewQueryCommand queryCmd represents the query command
func NewQueryCommand() *cobra.Command {
	var (
		optionQuery = &query.OptionExpQuery{}
		ifAll       bool
		format      string
	)

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "experiment query command",
		Run: func(cmd *cobra.Command, args []string) {
			query.PrintExpByOption(utils.GetCtxWithTraceId(context.Background(), utils.TraceId), optionQuery, ifAll, format)
		},
	}

	queryCmd.Flags().StringVarP(&optionQuery.Uid, "uid", "u", "", "query experiment by uid, eg: chaosmetad query -u [uid]")
	queryCmd.Flags().StringVarP(&optionQuery.Status, "status", "s", "", "query experiment by status, eg: chaosmetad query -s success")
	queryCmd.Flags().StringVarP(&optionQuery.Target, "target", "t", "", "query experiment by target, eg: chaosmetad query -t cpu")
	queryCmd.Flags().StringVarP(&optionQuery.Fault, "fault", "f", "", "query experiment by target and fault, eg: chaosmetad query -t cpu -f burn")
	queryCmd.Flags().StringVarP(&optionQuery.Creator, "creator", "c", "", "query experiment by creator, eg: chaosmetad query -c root")
	queryCmd.Flags().UintVarP(&optionQuery.Offset, "offset", "o", 0, "query experiment records with offset, eg: chaosmetad query -o 5")
	queryCmd.Flags().UintVarP(&optionQuery.Limit, "limit", "l", 10, "query experiment records with limit, eg: chaosmetad query -o 5 -l 5")
	queryCmd.Flags().BoolVarP(&ifAll, "all", "a", false, "if show all")
	queryCmd.Flags().StringVar(&format, "format", query.TableFormat, fmt.Sprintf("data show format, support: %s(default), %s", query.TableFormat, query.JsonFormat))

	return queryCmd
}
