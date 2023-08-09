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
