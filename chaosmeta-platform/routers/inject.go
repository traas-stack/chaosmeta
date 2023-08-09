package routers

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/inject"
	beego "github.com/beego/beego/v2/server/web"
)

func injectInit() {
	beego.Router(NewWebServicePath("injects/scopes"), &inject.InjectController{}, "get:QueryScopes")
	beego.Router(NewWebServicePath("injects/scopes/:id/targets"), &inject.InjectController{}, "get:QueryTargets")
	beego.Router(NewWebServicePath("injects/scopes/:id/targets/:targets_id/faults"), &inject.InjectController{}, "get:QueryFaults")
	beego.Router(NewWebServicePath("injects/faults/:id/args"), &inject.InjectController{}, "get:QueryArgs")
}
