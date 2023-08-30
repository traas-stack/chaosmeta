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
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/namespace"
	beego "github.com/beego/beego/v2/server/web"
)

func nameSpaceInit() {
	beego.Router(NewWebServicePath("namespaces"), &namespace.NamespaceController{}, "post:Create")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "get:Get")
	beego.Router(NewWebServicePath("namespaces/:id/permission"), &namespace.NamespaceController{}, "get:GetPermission")
	beego.Router(NewWebServicePath("namespaces/:id/overview"), &namespace.NamespaceController{}, "get:GetOverview")
	beego.Router(NewWebServicePath("namespaces/list"), &namespace.NamespaceController{}, "get:GetList")
	beego.Router(NewWebServicePath("namespaces/query"), &namespace.NamespaceController{}, "get:QueryList")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "post:Update")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "delete:Delete")
	beego.Router(NewWebServicePath("namespaces/:id/users"), &namespace.NamespaceController{}, "get:GetUserList")
	beego.Router(NewWebServicePath("namespaces/:id/users/batch"), &namespace.NamespaceController{}, "post:AddUsers")
	beego.Router(NewWebServicePath("namespaces/:id/users/:user_id"), &namespace.NamespaceController{}, "delete:RemoveUser")
	beego.Router(NewWebServicePath("namespaces/:id/users"), &namespace.NamespaceController{}, "delete:RemoveUsers")
	beego.Router(NewWebServicePath("namespaces/:id/users/permission"), &namespace.NamespaceController{}, "post:ChangePermissions")

	beego.Router(NewWebServicePath("namespaces/:id/labels"), &namespace.NamespaceController{}, "get:ListLabel")
	beego.Router(NewWebServicePath("namespaces/:id/labels"), &namespace.NamespaceController{}, "post:LabelCreate")
	beego.Router(NewWebServicePath("namespaces/:ns_id/labels/:id"), &namespace.NamespaceController{}, "delete:LabelDelete")
	beego.Router(NewWebServicePath("namespaces/:ns_id/labels/:name"), &namespace.NamespaceController{}, "get:LabelGet")

	beego.Router(NewWebServicePath("namespaces/:id/cluster"), &namespace.NamespaceController{}, "post:SetAttackableCluster")
	beego.Router(NewWebServicePath("namespaces/:id/cluster"), &namespace.NamespaceController{}, "get:ListAttackableCluster")
}
