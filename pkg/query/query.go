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
	"encoding/json"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/log"
	"github.com/ChaosMetaverse/chaosmetad/pkg/storage"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/bndr/gotabulate"
)

type OptionExpQuery struct {
	Uid     string `json:"uid,omitempty"`
	Status  string `json:"status,omitempty"`
	Target  string `json:"target,omitempty"`
	Fault   string `json:"fault,omitempty"`
	Creator string `json:"creator,omitempty"`
	Offset  uint   `json:"offset"`
	Limit   uint   `json:"limit"`
}

func GetExpByOption(o *OptionExpQuery, ifAll bool) {
	if o == nil {
		utils.SolveErr(utils.BadArgsErr, fmt.Sprintf("option is empty"))
	}

	temp, err := json.Marshal(o)
	if err != nil {
		utils.SolveErr(utils.BadArgsErr, err.Error())
	}

	log.GetLogger().Infof("query args: %s", string(temp))

	db, dbErr := storage.GetExperimentStore()
	if dbErr != nil {
		utils.SolveErr(utils.DBErr, dbErr.Error())
	}
	exps, total, queryErr := db.QueryByOption(o.Uid, o.Status, o.Target, o.Fault, o.Creator, o.Offset, o.Limit)
	if queryErr != nil {
		utils.SolveErr(utils.DBErr, queryErr.Error())
	}

	printExp(exps, total, ifAll)
}

func printExp(exps []*storage.Experiment, total int64, ifAll bool) {
	logger := log.GetLogger()
	var formatData string
	if len(exps) != 0 {
		var data [][]interface{}
		for _, exp := range exps {
			var aData []interface{}
			if ifAll {
				aData = []interface{}{exp.Uid, exp.Status, exp.Target, exp.Fault, exp.Args, exp.Creator, exp.Runtime,
					exp.ContainerId, exp.ContainerRuntime, exp.Timeout, exp.Error, exp.CreateTime, exp.UpdateTime}
			} else {
				aData = []interface{}{exp.Uid, exp.Status, exp.Target, exp.Fault, exp.Args}
			}

			data = append(data, aData)
		}

		t := gotabulate.Create(data)
		if ifAll {
			t.SetHeaders([]string{"UID", "STATUS", "TARGET", "FAULT", "ARGS", "CREATOR", "RUNTIME",
				"CONTAINER_ID", "CONTAINER_RUNTIME", "TIMEOUT", "ERROR", "CREATE_TIME", "UPDATE_TIME"})
		} else {
			t.SetHeaders([]string{"UID", "STATUS", "TARGET", "FAULT", "ARGS"})
		}
		t.SetEmptyString("None")
		t.SetAlign("left")
		t.SetWrapStrings(true)
		formatData = t.Render("grid")
	}

	logger.Infof("total count of experiments: %d\n%s", total, formatData)
}
