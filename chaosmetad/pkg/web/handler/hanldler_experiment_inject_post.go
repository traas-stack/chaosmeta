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
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/injector"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/storage"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/utils/errutil"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/web/model"
	"net/http"
)

func ExperimentInjectPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var (
		ctx       = context.Background()
		injectReq = &model.InjectRequest{}
		injectRes *model.InjectResponse
	)

	if err := json.NewDecoder(r.Body).Decode(injectReq); err != nil {
		injectRes = getExperimentInjectPostResponse(ctx, errutil.BadArgsErr, fmt.Sprintf("req body format error: %s", err.Error()), nil)
	} else {
		ctx = utils.GetCtxWithTraceId(ctx, injectReq.TraceId)
		i, err := injector.NewInjector(injectReq.Target, injectReq.Fault)
		if err != nil {
			injectRes = getExperimentInjectPostResponse(ctx, errutil.BadArgsErr, fmt.Sprintf("get injector error: %s", err.Error()), nil)
		} else {
			creator := injectReq.Creator
			if creator == "" {
				creator = r.Host
			}

			if err := i.LoadInjector(&storage.Experiment{
				Uid:              injectReq.Uid,
				Target:           injectReq.Target,
				Fault:            injectReq.Fault,
				Args:             injectReq.Args,
				Timeout:          injectReq.Timeout,
				ContainerRuntime: injectReq.ContainerRuntime,
				ContainerId:      injectReq.ContainerId,
				Creator:          creator,
				Runtime:          "{}",
			}, i.GetArgs(), i.GetRuntime()); err != nil {
				injectRes = getExperimentInjectPostResponse(ctx, errutil.BadArgsErr, fmt.Sprintf("args load error: %s", err.Error()), nil)
			} else {
				code, msg := injector.ProcessInject(ctx, i)
				if code == errutil.NoErr {
					exp, err := i.OptionToExp(i.GetArgs(), i.GetRuntime())
					if err != nil {
						injectRes = getExperimentInjectPostResponse(ctx, errutil.NoErr, fmt.Sprintf("inject success but get exp info error: %s", err.Error()), nil)
					} else {
						injectRes = getExperimentInjectPostResponse(ctx, errutil.NoErr, "success", exp)
					}
				} else {
					injectRes = getExperimentInjectPostResponse(ctx, errutil.InjectErr, fmt.Sprintf("injector error: %s", msg), nil)
				}
			}
		}
	}

	WriteResponse(ctx, w, injectRes)
}

func getExperimentInjectPostResponse(ctx context.Context, code int, msg string, exp *storage.Experiment) *model.InjectResponse {
	var re = &model.InjectResponse{
		Code:    code,
		Message: msg,
		TraceId: utils.GetTraceId(ctx),
	}

	if exp != nil {
		re.Data = &model.InjectSuccessResponseData{
			Experiment: ExpToExperimentDataUnit(exp),
		}
	}

	return re
}
