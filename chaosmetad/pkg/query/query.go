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
	"encoding/json"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/storage"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/web/handler"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/web/model"
)

type OptionExpQuery struct {
	Uid              string `json:"uid,omitempty"`
	Status           string `json:"status,omitempty"`
	Target           string `json:"target,omitempty"`
	Fault            string `json:"fault,omitempty"`
	Creator          string `json:"creator,omitempty"`
	ContainerId      string `json:"container_id,omitempty"`
	ContainerRuntime string `json:"container_runtime,omitempty"`
	Offset           uint   `json:"offset"`
	Limit            uint   `json:"limit"`
}

const (
	TableFormat = "table"
	JsonFormat  = "json"
)

func PrintExpByOption(ctx context.Context, o *OptionExpQuery, ifAll bool, format string) {
	if format != TableFormat && format != JsonFormat {
		errutil.SolveErr(ctx, errutil.BadArgsErr, fmt.Sprintf("not support format: %s", format))
	}

	if o == nil {
		errutil.SolveErr(ctx, errutil.BadArgsErr, fmt.Sprintf("option is empty"))
	}

	temp, err := json.Marshal(o)
	if err != nil {
		errutil.SolveErr(ctx, errutil.BadArgsErr, err.Error())
	}

	db, dbErr := storage.GetExperimentStore()
	if dbErr != nil {
		errutil.SolveErr(ctx, errutil.DBErr, dbErr.Error())
	}
	exps, total, queryErr := db.QueryByOption(o.Uid, o.Status, o.Target, o.Fault, o.Creator, o.ContainerRuntime, o.ContainerId, o.Offset, o.Limit)
	if queryErr != nil {
		errutil.SolveErr(ctx, errutil.DBErr, queryErr.Error())
	}

	if format == JsonFormat {
		printJson(ctx, exps, total)
	} else {
		log.GetLogger(ctx).Infof("query args: %s", string(temp))
		printTable(ctx, exps, total, ifAll)
	}
}

func printJson(ctx context.Context, exps []*storage.Experiment, total int64) {
	logger := log.GetLogger(ctx)
	reList := make([]model.ExperimentDataUnit, len(exps))
	for i, exp := range exps {
		reList[i] = handler.ExpToExperimentDataUnit(exp)
	}

	res := &model.QueryResponseData{
		Experiments: reList,
		Total:       total,
	}

	reBytes, err := json.Marshal(res)
	if err != nil {
		errutil.SolveErr(ctx, errutil.InternalErr, fmt.Sprintf("query response change to string error: %s", err.Error()))
	}

	if log.Path != "" {
		logger.Info(string(reBytes))
	} else {
		fmt.Println(string(reBytes))
	}
}

func printTable(ctx context.Context, exps []*storage.Experiment, total int64, ifAll bool) {
	logger := log.GetLogger(ctx)
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

	logger.Infof("total count of experiments: %d\n%s\n", total, formatData)
}
