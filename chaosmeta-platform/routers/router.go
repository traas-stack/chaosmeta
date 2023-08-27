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

package routers

import (
	"chaosmeta-platform/pkg/service"
	userService "chaosmeta-platform/pkg/service/user"
	"chaosmeta-platform/util/errors"
	"chaosmeta-platform/util/log"
	"context"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"strings"
)

const RootAPI = "/chaosmeta/api/v1/%s"

func routerInit() {
	userInit()
	nameSpaceInit()
	kubernetesInit()
	injectInit()
	experimentInit()
	experimentInstanceInit()
}

func Init() {
	beego.InsertFilter("/chaosmeta/api/*", beego.BeforeRouter, CheckTokenMiddleware)
	routerInit()
	beego.Router("/", &service.MainController{})
}

func CheckTokenMiddleware(ctx *beecontext.Context) {
	token := ctx.Input.Header("Authorization")
	if token == "" {
		log.Error("token is empty")
		ctx.Output.JSON(errors.ErrUnauthorized().WithMessage("no token"), false, false)
		return
	}

	a := &userService.UserService{}
	userName, err := a.CheckToken(context.Background(), strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		log.Error(err)
		ctx.Output.JSON(errors.ErrUnauthorized().WithMessage(err.Error()), false, false)
		return
	}
	ctx.Input.SetData("userName", userName)
}

func NewWebServicePath(prefix string) string {
	if prefix != "" {
		return fmt.Sprintf(RootAPI, prefix)
	}
	return ""
}
