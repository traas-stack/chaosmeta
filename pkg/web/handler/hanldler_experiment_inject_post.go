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
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/storage"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/web/model"
	"net/http"
)

func ExperimentInjectPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var injectReq = &model.InjectRequest{}
	var injectRes *model.InjectResponse

	if err := json.NewDecoder(r.Body).Decode(injectReq); err != nil {
		injectRes = getExperimentInjectPostResponse(utils.BadArgsErr, fmt.Sprintf("req body format error: %s", err.Error()), nil)
	} else {
		i, err := injector.NewInjector(injectReq.Target, injectReq.Fault)
		if err != nil {
			injectRes = getExperimentInjectPostResponse(utils.BadArgsErr, fmt.Sprintf("get injector error: %s", err.Error()), nil)
		} else {
			// load args
			creator := injectReq.Creator
			if creator == "" {
				creator = r.Host
			}
			if err := i.LoadInjector(&storage.Experiment{
				// r.Host
				Target:  injectReq.Target,
				Fault:   injectReq.Fault,
				Args:    injectReq.Args,
				Timeout: injectReq.Timeout,
				Creator: creator,
				Runtime: "{}",
			}, i.GetArgs(), i.GetRuntime()); err != nil {
				injectRes = getExperimentInjectPostResponse(utils.BadArgsErr, fmt.Sprintf("args load error: %s", err.Error()), nil)
			} else {
				code, msg := injector.ProcessInject(i)
				if code == utils.NoErr {
					exp, err := i.OptionToExp(i.GetArgs(), i.GetRuntime())
					if err != nil {
						injectRes = getExperimentInjectPostResponse(utils.NoErr, fmt.Sprintf("inject success but get exp info error: %s", err.Error()), nil)
					} else {
						injectRes = getExperimentInjectPostResponse(utils.NoErr, "success", exp)
					}
				} else {
					injectRes = getExperimentInjectPostResponse(utils.InjectErr, fmt.Sprintf("injector error: %s", msg), nil)
				}
			}
		}
	}

	WriteResponse(w, injectRes)
}

func getExperimentInjectPostResponse(code int, msg string, exp *storage.Experiment) *model.InjectResponse {
	var re = &model.InjectResponse{
		Code:    code,
		Message: msg,
	}

	if exp != nil {
		re.Data = &model.InjectSuccessResponseData{
			Experiment: model.ExperimentDataUnit{
				Uid:     exp.Uid,
				Target:  exp.Target,
				Fault:   exp.Fault,
				Status:  exp.Status,
				Error_:  exp.Error,
				Creator: exp.Creator,
				Timeout: exp.Timeout,
				Args:    exp.Args,
			},
		}
	}

	return re
}
