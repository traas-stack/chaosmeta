package routers

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/user"
	beego "github.com/beego/beego/v2/server/web"
)

func userInit() {
	beego.Router("/users/token/create", &user.UserController{}, "post:Create")
	beego.Router("/users/token/login", &user.UserController{}, "post:Login")
	beego.Router("/users/token/refresh", &user.UserController{}, "post:RefreshToken")
	beego.Router(NewWebServicePath("users/list"), &user.UserController{}, "get:GetList")
	beego.Router(NewWebServicePath("users/namespace/list"), &user.UserController{}, "get:GetNamespaceList")
	beego.Router(NewWebServicePath("users/namespace/:id/user_list"), &user.UserController{}, "get:GetListWithNamespaceInfo")
	beego.Router(NewWebServicePath("users/:name"), &user.UserController{}, "get:Get")
	beego.Router(NewWebServicePath("users/:id"), &user.UserController{}, "delete:Delete")
	beego.Router(NewWebServicePath("users"), &user.UserController{}, "delete:DeleteList")
	beego.Router(NewWebServicePath("users/password"), &user.UserController{}, "post:UpdateUserPassword")
	beego.Router(NewWebServicePath("users/role"), &user.UserController{}, "post:UpdateListRole")
}
