package routers

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1/namespace"
	beego "github.com/beego/beego/v2/server/web"
)

func nameSpaceInit() {
	beego.Router(NewWebServicePath("namespaces"), &namespace.NamespaceController{}, "post:Create")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "get:Get")
	beego.Router(NewWebServicePath("namespaces/list"), &namespace.NamespaceController{}, "get:GetList")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "post:Update")
	beego.Router(NewWebServicePath("namespaces/:id"), &namespace.NamespaceController{}, "delete:Delete")
	beego.Router(NewWebServicePath("namespaces/:id/users"), &namespace.NamespaceController{}, "get:GetUserList")
	beego.Router(NewWebServicePath("namespaces/:id/users/batch"), &namespace.NamespaceController{}, "post:AddUsers")
	beego.Router(NewWebServicePath("namespaces/:id/users/:user_id"), &namespace.NamespaceController{}, "delete:RemoveUser")
	beego.Router(NewWebServicePath("namespaces/:id/users"), &namespace.NamespaceController{}, "delete:RemoveUsers")
	beego.Router(NewWebServicePath("namespaces/:id/users/permission"), &namespace.NamespaceController{}, "post:ChangePermissions")

	beego.Router(NewWebServicePath("/namespaces/:id/labels"), &namespace.NamespaceController{}, "get:ListLabel")
	beego.Router(NewWebServicePath("/namespaces/:id/labels"), &namespace.NamespaceController{}, "post:LabelCreate")
	beego.Router(NewWebServicePath("/namespaces/:ns_id/labels/:id"), &namespace.NamespaceController{}, "delete:LabelDelete")
}
