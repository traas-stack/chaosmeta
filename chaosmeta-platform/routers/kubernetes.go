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
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/cluster"
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/kube"
	beego "github.com/beego/beego/v2/server/web"
)

func kubernetesInit() {
	beego.Router(NewWebServicePath("kubernetes/cluster"), &cluster.ClusterController{}, "post:Create")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "get:Get")
	beego.Router(NewWebServicePath("kubernetes/cluster/list"), &cluster.ClusterController{}, "get:GetList")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "post:Update")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id"), &cluster.ClusterController{}, "delete:Delete")

	beego.Router(NewWebServicePath("kubernetes/cluster/:id/nodes"), &kube.KubeController{}, "get:ListNodes")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id/namespaces"), &kube.KubeController{}, "get:ListNamespaces")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id/namespace/:ns_name/pods"), &kube.KubeController{}, "get:ListPods")
	beego.Router(NewWebServicePath("kubernetes/cluster/:id/namespace/:ns_name/deployments"), &kube.KubeController{}, "get:ListDeployments")
}
