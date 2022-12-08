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
	"github.com/ChaosMetaverse/chaosmetad/pkg/injector"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/errutil"
	"github.com/ChaosMetaverse/chaosmetad/pkg/web/model"
	"github.com/gorilla/mux"
	"net/http"
)

// TODO: Consider whether to add "TraceId" args

func ExperimentUidRecoverPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var (
		recoverRes *model.CommonResponse
		ctx        = utils.GetCtxWithTraceId(context.Background(), "")
	)

	vars := mux.Vars(r)
	uid := vars["uid"]
	if uid == "" {
		recoverRes = getCommonResponse(ctx, errutil.BadArgsErr, "uid is empty")
	} else {
		code, msg := injector.ProcessRecover(ctx, uid)
		recoverRes = getCommonResponse(ctx, code, msg)
	}

	WriteResponse(ctx, w, recoverRes)
}

func getCommonResponse(ctx context.Context, code int, msg string) *model.CommonResponse {
	return &model.CommonResponse{
		Code:    code,
		Message: msg,
		TraceId: utils.GetTraceId(ctx),
	}
}
