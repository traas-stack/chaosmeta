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

package namespace

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	namespace2 "chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/models/user"
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
	namespaceId, err := namespace.Create(context.Background(), requestBody.Name, requestBody.Description, username)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, CreateNamespaceResponse{
		ID: namespaceId,
	})
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

func (c *NamespaceController) QueryList() {
	nameSpaceName := c.GetString("name")
	userNameQuery := c.GetString("userName")
	username := c.Ctx.Input.GetData("userName").(string)
	namespaceClass := c.GetString("namespaceClass")

	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	sort := c.GetString("sort")

	userGet := user.User{Email: username}

	if err := user.GetUser(context.Background(), &userGet); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	namespaceService := &namespace.NamespaceService{}
	var (
		total         int64
		namespaceList []namespace.NamespaceData
		err           error
	)

	switch namespaceClass {
	case "read":
		total, namespaceList, err = namespaceService.GroupNamespacesByUsername(context.Background(), userGet.ID, nameSpaceName, userNameQuery, 0, sort, page, pageSize)
		if err != nil {
			c.Error(&c.Controller, err)
			return
		}
	case "write":
		total, namespaceList, err = namespaceService.GroupNamespacesByUsername(context.Background(), userGet.ID, nameSpaceName, userNameQuery, 1, sort, page, pageSize)
		if err != nil {
			c.Error(&c.Controller, err)
			return
		}
	case "all":
		total, namespaceList, err = namespaceService.GroupAllNamespaces(context.Background(), nameSpaceName, userNameQuery, sort, page, pageSize)
		if err != nil {
			c.Error(&c.Controller, err)
			return
		}
	case "not":
		total, namespaceList, err = namespaceService.GroupNamespacesUserNotIn(context.Background(), userGet.ID, nameSpaceName, userNameQuery, sort, page, pageSize)
		if err != nil {
			c.Error(&c.Controller, err)
			return
		}
	}

	queryNamespaceResponse := QueryNamespaceResponse{Total: total, Page: page, PageSize: pageSize, NameSpaces: namespaceList}
	c.Success(&c.Controller, queryNamespaceResponse)
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
	total, usrList, err := namespace.GroupedUserInNamespaces(context.Background(), namespaceId, "", usernameQuery, permission, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	userListResponse := UserListResponse{Total: total, Page: page, PageSize: pageSize}
	for _, user := range usrList {
		userListResponse.Users = append(userListResponse.Users, &User{
			ID:         user.Id,
			Name:       user.Name,
			Permission: user.Permission,
			CreateTime: user.CreateTime,
		})
	}
	c.Success(&c.Controller, userListResponse)
}

//func (c *NamespaceController) GetUserList() {
//	namespaceId, err := c.GetInt(":id")
//	if err != nil {
//		c.Error(&c.Controller, err)
//		return
//	}
//
//	page, _ := c.GetInt("page", 1)
//	pageSize, _ := c.GetInt("page_size", 10)
//	sort := c.GetString("sort")
//	usernameQuery := c.GetString("username")
//	permission, err := c.GetInt("permission")
//	if err != nil {
//		permission = -1
//	}
//
//	namespace := &namespace.NamespaceService{}
//	total, usrList, err := namespace.GroupedUserNamespaces(context.Background(), namespaceId, usernameQuery, permission, sort, page, pageSize)
//	if err != nil {
//		c.Error(&c.Controller, err)
//		return
//	}
//
//	userListResponse := UserListResponse{Total: total, Page: page, PageSize: pageSize}
//	for _, user := range usrList {
//		userListResponse.Users = append(userListResponse.Users, &User{
//			ID:         user.Id,
//			Name:       user.Name,
//			Permission: user.Permission,
//			CreateTime: user.CreateTime,
//		})
//	}
//	c.Success(&c.Controller, userListResponse)
//}

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
