package routers

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/experiment"
	beego "github.com/beego/beego/v2/server/web"
)

func experimentInit() {
	beego.Router("/experiments", &experiment.ExperimentController{}, "get:GetExperimentList")
	beego.Router("/experiments/:uuid", &experiment.ExperimentController{}, "get:GetExperimentDetail")
	beego.Router("/experiments", &experiment.ExperimentController{}, "post:CreateExperiment")
	beego.Router("/experiments/:uuid", &experiment.ExperimentController{}, "post:UpdateExperiment")
	beego.Router("/experiments/:uuid", &experiment.ExperimentController{}, "delete:DeleteExperiment")

}
