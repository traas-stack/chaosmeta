package routers

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/cluster"
	beego "github.com/beego/beego/v2/server/web"
)

func kubernetesInit() {
	beego.Router(NewWebServicePath("kubernetes/cluster"), &cluster.ClusterController{}, "post:Create")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "get:Get")
	beego.Router(NewWebServicePath("kubernetes/cluster/list"), &cluster.ClusterController{}, "get:GetList")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "post:Update")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "delete:Delete")
}
