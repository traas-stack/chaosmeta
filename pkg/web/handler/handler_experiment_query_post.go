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
	"encoding/json"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/storage"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/web/model"
	"net/http"
)

func ExperimentQueryPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var queryReq = &model.ExperimentQueryRequest{}
	var queryRes *model.ExperimentQueryResponse

	if err := json.NewDecoder(r.Body).Decode(queryReq); err != nil {
		queryRes = getExperimentQueryPostResponse(utils.BadArgsErr, fmt.Sprintf("req body format error: %s", err.Error()), nil, 0)
	} else {
		db, dbErr := storage.GetExperimentStore()
		if dbErr != nil {
			queryRes = getExperimentQueryPostResponse(utils.DBErr, fmt.Sprintf("get db error: %s", dbErr.Error()), nil, 0)
		} else {
			exps, total, qErr := db.QueryByOption(queryReq.Uid, queryReq.Status, queryReq.Target, queryReq.Fault, queryReq.Creator, uint(queryReq.Offset), uint(queryReq.Limit))
			if qErr != nil {
				queryRes = getExperimentQueryPostResponse(utils.DBErr, fmt.Sprintf("db query error: %s", qErr.Error()), nil, 0)
			}
			queryRes = getExperimentQueryPostResponse(utils.NoErr, "success", exps, total)
		}
	}

	WriteResponse(w, queryRes)
}

func getExperimentQueryPostResponse(code int, msg string, exps []*storage.Experiment, total int64) *model.ExperimentQueryResponse {
	var re = &model.ExperimentQueryResponse{
		Code:    code,
		Message: msg,
	}
	if exps != nil {
		reList := make([]model.ExperimentDataUnit, len(exps))
		for i, exp := range exps {
			reList[i] = expToExperimentDataUnit(exp)
		}

		re.Data = &model.ExperimentQueryResponseData{
			Experiments: reList,
			Total:       total,
		}
	}

	return re
}

func expToExperimentDataUnit(exp *storage.Experiment) model.ExperimentDataUnit {
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
