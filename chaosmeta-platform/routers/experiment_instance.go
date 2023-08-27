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
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/experiment_instance"
	beego "github.com/beego/beego/v2/server/web"
)

func experimentInstanceInit() {
	beego.Router(NewWebServicePath("experiments/results"), &experiment_instance.ExperimentInstanceController{}, "get:GetExperimentInstances")
	beego.Router(NewWebServicePath("experiments/results/:uuid"), &experiment_instance.ExperimentInstanceController{}, "get:GetExperimentInstanceDetail")
	beego.Router(NewWebServicePath("experiments/results/:uuid/nodes"), &experiment_instance.ExperimentInstanceController{}, "get:GetExperimentInstanceNodes")
	beego.Router(NewWebServicePath("experiments/results/:uuid/nodes/:node_id"), &experiment_instance.ExperimentInstanceController{}, "get:GetExperimentInstanceNode")
	beego.Router(NewWebServicePath("experiments/results/:uuid/nodes/:node_id/subtasks/:id"), &experiment_instance.ExperimentInstanceController{}, "get:GetExperimentInstanceNodeSubtask")
	beego.Router(NewWebServicePath("experiments/results/:uuid"), &experiment_instance.ExperimentInstanceController{}, "delete:DeleteExperimentInstance")
	beego.Router(NewWebServicePath("experiments/results"), &experiment_instance.ExperimentInstanceController{}, "delete:DeleteExperimentInstances")
}
