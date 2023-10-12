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
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/inject"
	beego "github.com/beego/beego/v2/server/web"
)

func injectInit() {
	beego.Router(NewWebServicePath("injects/scopes"), &inject.InjectController{}, "get:QueryScopes")
	beego.Router(NewWebServicePath("injects/scopes/:id/targets"), &inject.InjectController{}, "get:QueryTargets")
	beego.Router(NewWebServicePath("injects/scopes/:id/targets/:targets_id/faults"), &inject.InjectController{}, "get:QueryFaults")
	beego.Router(NewWebServicePath("injects/flows"), &inject.InjectController{}, "get:QueryFlows")
	beego.Router(NewWebServicePath("injects/measures"), &inject.InjectController{}, "get:QueryMeasures")
	beego.Router(NewWebServicePath("injects/faults/:id/args"), &inject.InjectController{}, "get:QueryFaultArgs")
	beego.Router(NewWebServicePath("injects/flows/:id/args"), &inject.InjectController{}, "get:QueryFlowArgs")
	beego.Router(NewWebServicePath("injects/measures/:id/args"), &inject.InjectController{}, "get:QueryMeasureArgs")
}
