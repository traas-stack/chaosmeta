package namespace

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	namespace2 "chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/service/namespace"
	"context"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
)

type NamespaceController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *NamespaceController) Create() {
	var requestBody CreateNamespaceRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.Create(context.Background(), requestBody.Name, requestBody.Description, username); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) Get() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	namespace := &namespace.NamespaceService{}
	ns, err := namespace.Get(context.Background(), namespaceId)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	getNamespaceResponse := GetNamespaceResponse{}
	getNamespaceResponse.NameSpace = *ns
	c.Success(&c.Controller, getNamespaceResponse)
}

func (c *NamespaceController) GetList() {
	sort := c.GetString("sort")
	name := c.GetString("name")
	creator := c.GetString("creator")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	namespace := &namespace.NamespaceService{}
	total, namespaceList, err := namespace.GetList(context.Background(), name, creator, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	listNamespaceResponse := ListNamespaceResponse{Total: total, Page: page, PageSize: pageSize}
	listNamespaceResponse.NameSpaces = append(listNamespaceResponse.NameSpaces, namespaceList...)
	c.Success(&c.Controller, listNamespaceResponse)
}

func (c *NamespaceController) Update() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var requestBody UpdateNamespaceRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.Update(context.Background(), username, namespaceId, requestBody.Name, requestBody.Description); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) Delete() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.Delete(context.Background(), username, namespaceId); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) GetUserList() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	sort := c.GetString("sort")
	usernameQuery := c.GetString("username")
	permission, err := c.GetInt("permission")
	if err != nil {
		permission = -1
	}

	namespace := &namespace.NamespaceService{}
	usrList, total, err := namespace.GetUsers(context.Background(), namespaceId, usernameQuery, permission, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	userListResponse := UserListResponse{Total: total, Page: page, PageSize: pageSize}
	for _, user := range usrList {
		userListResponse.Users = append(userListResponse.Users, &User{
			ID:         user.ID,
			Name:       user.Email,
			Role:       user.Role,
			CreateTime: user.CreateTime,
			UpdateTime: user.UpdateTime,
		})
	}
	c.Success(&c.Controller, userListResponse)
}

func (c *NamespaceController) AddUsers() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var reqBody AddUsersRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	namespace := &namespace.NamespaceService{}
	username := c.Ctx.Input.GetData("userName").(string)
	var userParam namespace2.AddUsersParam
	for _, users := range reqBody.Users {
		userParam.Users = append(userParam.Users, namespace2.UserData{
			Id:         users.Id,
			Permission: users.Permission,
		})
	}

	if err := namespace.AddUsers(context.Background(), username, namespaceId, userParam); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) RemoveUser() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	userId, err := c.GetInt(":user_id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	name := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.RemoveUsers(context.Background(), name, []int{userId}, namespaceId); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) RemoveUsers() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var reqBody RemoveUsersRequest
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.RemoveUsers(context.Background(), username, reqBody.UserIds, namespaceId); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) ChangePermissions() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var reqBody ChangeUsersPermissionRequest
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.ChangeUsersPermission(context.Background(), username, reqBody.UserIds, namespaceId, namespace2.Permission(reqBody.Permission)); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}
