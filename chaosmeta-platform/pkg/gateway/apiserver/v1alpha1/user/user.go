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

package user

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/user"
	"chaosmeta-platform/util/errors"
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"strings"
)

type UserController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *UserController) GetList() {
	sort := c.GetString("sort")
	name := c.GetString("name")
	role := c.GetString("role")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)

	a := &user.UserService{}
	total, usrList, err := a.GetList(context.Background(), name, role, sort, page, pageSize)
	if err != nil {
		c.Data["json"] = errors.ErrServer().WithMessage(err.Error())
		c.ServeJSON()
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

func (c *UserController) GetListWithNamespaceInfo() {
	nsId, _ := c.GetInt(":id")
	sort := c.GetString("sort")
	name := c.GetString("name")
	role := c.GetString("role")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)

	a := &user.UserService{}
	total, userNamespaceList, err := a.GetListWithNamespaceInfo(context.Background(), nsId, name, role, sort, page, pageSize)
	if err != nil {
		c.Data["json"] = errors.ErrServer().WithMessage(err.Error())
		c.ServeJSON()
		return
	}

	userListNamespaceResponse := UserListNamespaceResponse{Total: total, Page: page, PageSize: pageSize}
	for _, userNamespace := range userNamespaceList {
		userListNamespaceResponse.Users = append(userListNamespaceResponse.Users, &UserNamespace{
			User: User{
				ID:         userNamespace.User.ID,
				Name:       userNamespace.User.Email,
				Role:       userNamespace.User.Role,
				CreateTime: userNamespace.User.CreateTime,
				UpdateTime: userNamespace.User.UpdateTime,
			},
			IsJoin:     userNamespace.IsJoin,
			Permission: userNamespace.Permission,
		})
	}
	c.Success(&c.Controller, userListNamespaceResponse)
}

func (c *UserController) Get() {
	name := c.Ctx.Input.Param(":name")

	a := &user.UserService{}
	user, err := a.Get(context.Background(), name)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, User{
		ID:         user.ID,
		Name:       user.Email,
		Role:       user.Role,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	})
}

func (c *UserController) Create() {
	var UserCreateRequest UserCreateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &UserCreateRequest); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	a := &user.UserService{}
	if _, err := a.Create(context.Background(), UserCreateRequest.Name, UserCreateRequest.Password, string(user.NormalRole)); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *UserController) Login() {
	var UserLoginRequest UserLoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &UserLoginRequest); err != nil {
		c.Data["json"] = errors.ErrServer().WithMessage(err.Error())
		c.ServeJSON()
		log.Error(err)
		return
	}
	a := &user.UserService{}
	token, refreshToken, err := a.Login(context.Background(), UserLoginRequest.Name, UserLoginRequest.Password)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	c.Ctx.Output.Cookie("TOKEN", token)
	c.Ctx.Output.Cookie("REFRESH_TOKEN", refreshToken)
	c.Success(&c.Controller, UserLoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (c *UserController) RefreshToken() {
	a := &user.UserService{}

	token := strings.TrimPrefix(c.Ctx.Input.Header("Authorization"), "Bearer ")
	if token == "" {
		c.ErrorWithMessage(&c.Controller, "no token")
		return
	}

	refreshToken, err := a.RefreshToken(context.Background(), token)
	if err != nil {
		c.ErrUnauthorized(&c.Controller, err)
		return
	}

	c.Ctx.Output.Cookie("TOKEN", refreshToken)
	c.Success(&c.Controller, UserLoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (c *UserController) Delete() {
	userName := c.Ctx.Input.GetData("userName").(string)

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	a := &user.UserService{}
	if err := a.DeleteList(context.Background(), userName, []int{id}); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	c.Success(&c.Controller, "ok")
}

func (c *UserController) DeleteList() {
	var usersDeleteRequest UsersDeleteRequest
	userName := c.Ctx.Input.GetData("userName").(string)

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &usersDeleteRequest)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	a := &user.UserService{}
	if err := a.DeleteList(context.Background(), userName, usersDeleteRequest.UserIDs); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *UserController) UpdateUserPassword() {
	userName := c.Ctx.Input.GetData("userName").(string)

	var usersPasswordUpdateRequest UsersPasswordUpdateRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &usersPasswordUpdateRequest)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	a := &user.UserService{}
	if err := a.UpdatePassword(context.Background(), userName, usersPasswordUpdateRequest.Password); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *UserController) UpdateListRole() {
	userName := c.Ctx.Input.GetData("userName").(string)

	var userUpdateRoleRequest UserUpdateRoleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userUpdateRoleRequest)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	a := &user.UserService{}
	if err := a.UpdateListRole(context.Background(), userName, userUpdateRoleRequest.UserIDs, userUpdateRoleRequest.Role); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

// 获取用户下的命名空间列表
func (c *UserController) GetNamespaceList() {
	userName := c.Ctx.Input.GetData("userName").(string)

	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	sort := c.GetString("sort")
	permission, err := c.GetInt("permission", -1)

	a := &user.UserService{}
	total, namespaceList, err := a.GetNamespaceList(context.Background(), userName, permission, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	nameSpaceListResponse := NameSpaceListResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		Namespaces: namespaceList,
	}
	c.Success(&c.Controller, nameSpaceListResponse)
}
