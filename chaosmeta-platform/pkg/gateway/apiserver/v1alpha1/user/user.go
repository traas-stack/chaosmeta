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
	tocken, refreshToken, err := a.Login(context.Background(), UserLoginRequest.Name, UserLoginRequest.Password)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	c.Ctx.Output.Cookie("TOKEN", tocken)
	c.Ctx.Output.Cookie("REFRESH_TOKEN", refreshToken)
	c.Success(&c.Controller, UserLoginResponse{
		Token:        tocken,
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
