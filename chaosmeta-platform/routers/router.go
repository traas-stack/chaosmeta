package routers

import (
	"chaosmeta-platform/pkg/service"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &service.MainController{})
}
