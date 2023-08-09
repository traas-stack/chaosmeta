package kube

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	beego "github.com/beego/beego/v2/server/web"
)

type KubeController struct {
	beego.Controller
	v1alpha1.BeegoOutputController
}
