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
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/experiment"
	beego "github.com/beego/beego/v2/server/web"
)

func experimentInit() {
	beego.Router(NewWebServicePath("experiments"), &experiment.ExperimentController{}, "get:GetExperimentList")
	beego.Router(NewWebServicePath("experiments/:uuid"), &experiment.ExperimentController{}, "get:GetExperimentDetail")
	beego.Router(NewWebServicePath("experiments"), &experiment.ExperimentController{}, "post:CreateExperiment")
	beego.Router(NewWebServicePath("experiments/:uuid"), &experiment.ExperimentController{}, "post:UpdateExperiment")
	beego.Router(NewWebServicePath("experiments/:uuid"), &experiment.ExperimentController{}, "delete:DeleteExperiment")

	beego.Router(NewWebServicePath("experiments/:uuid/start"), &experiment.ExperimentController{}, "post:StartExperiment")
	beego.Router(NewWebServicePath("experiments/:uuid/stop"), &experiment.ExperimentController{}, "post:StopExperiment")
}
