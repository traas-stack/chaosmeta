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

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmetad/pkg/storage"
	"github.com/traas-stack/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmetad/pkg/web/model"
	"net/http"
)

func ExperimentQueryPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var (
		ctx      = context.Background()
		queryReq = &model.QueryRequest{}
		queryRes *model.QueryResponse
	)

	if err := json.NewDecoder(r.Body).Decode(queryReq); err != nil {
		queryRes = getExperimentQueryPostResponse(ctx, errutil.BadArgsErr, fmt.Sprintf("req body format error: %s", err.Error()), nil, 0)
	} else {
		ctx = utils.GetCtxWithTraceId(ctx, queryReq.TraceId)
		db, dbErr := storage.GetExperimentStore()
		if dbErr != nil {
			queryRes = getExperimentQueryPostResponse(ctx, errutil.DBErr, fmt.Sprintf("get db error: %s", dbErr.Error()), nil, 0)
		} else {
			exps, total, qErr := db.QueryByOption(queryReq.Uid, queryReq.Status, queryReq.Target, queryReq.Fault,
				queryReq.Creator, queryReq.ContainerRuntime, queryReq.ContainerId, uint(queryReq.Offset), uint(queryReq.Limit))
			if qErr != nil {
				queryRes = getExperimentQueryPostResponse(ctx, errutil.DBErr, fmt.Sprintf("db query error: %s", qErr.Error()), nil, 0)
			}
			queryRes = getExperimentQueryPostResponse(ctx, errutil.NoErr, "success", exps, total)
		}
	}

	WriteResponse(ctx, w, queryRes)
}

func getExperimentQueryPostResponse(ctx context.Context, code int, msg string, exps []*storage.Experiment, total int64) *model.QueryResponse {
	var re = &model.QueryResponse{
		Code:    code,
		Message: msg,
		TraceId: utils.GetTraceId(ctx),
	}
	if exps != nil {
		reList := make([]model.ExperimentDataUnit, len(exps))
		for i, exp := range exps {
			reList[i] = ExpToExperimentDataUnit(exp)
		}

		re.Data = &model.QueryResponseData{
			Experiments: reList,
			Total:       total,
		}
	}

	return re
}

func ExpToExperimentDataUnit(exp *storage.Experiment) model.ExperimentDataUnit {
	return model.ExperimentDataUnit{
		Uid:              exp.Uid,
		Target:           exp.Target,
		Fault:            exp.Fault,
		Args:             exp.Args,
		Runtime:          exp.Runtime,
		Timeout:          exp.Timeout,
		Status:           exp.Status,
		Creator:          exp.Creator,
		Error_:           exp.Error,
		CreateTime:       exp.CreateTime,
		UpdateTime:       exp.UpdateTime,
		ContainerId:      exp.ContainerId,
		ContainerRuntime: exp.ContainerRuntime,
	}
}
