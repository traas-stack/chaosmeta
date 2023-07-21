package user

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/user"
	"chaosmeta-platform/util/errors"
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
)

type UserController struct {
	beego.Controller
}

func init() {
	beego.Router("/users", &UserController{}, "get:GetUserList")
}

//func (c *UserController) Prepare() {
//	// Verify that the token is valid
//	if !c.IsLogin() {
//		c.Ctx.Output.Status = http.StatusUnauthorized
//		c.Data["json"] = &UserListResponse{
//			ResponseData: v1alpha1.NewResponseData(http.StatusUnauthorized, "Invalid token", c.GetTraceID()),
//		}
//		c.ServeJSON()
//		return
//	}
//}

// @Title Get user list information
// @Description 获取用户列表信息
// @Tags User
// @Param sort query string false "排序方式，正序或倒序，例如：create_time 或 -create_time"
// @Param name query string false "筛选名称"
// @Param role query string false "筛选角色"
// @Param offset query int false "偏移量"
// @Param limit query int false "每页数量"
// @Success 200 {object} UserListResponse
// @Failure 401 {object} ResponseData
// @Failure 500 {object} ResponseData
// @router / [get]
func (c *UserController) GetUserList() {
	sort := c.GetString("sort")
	name := c.GetString("name")
	role := c.GetString("role")
	offset, _ := c.GetInt("offset", 0)
	limit, _ := c.GetInt("limit", 10)
	a := &user.User{}

	total, usrList, err := a.GetList(context.Background(), name, role, sort, offset, limit)
	if err != nil {
		c.Data["json"] = errors.ErrServer().WithMessage(err.Error())
		c.ServeJSON()
		return
	}

	userListResponse := UserListResponse{
		ResponseData: v1alpha1.NewResponseData(http.StatusInternalServerError, err.Error(), ""),
		Data: UserData{
			Total: total,
		},
	}

	for _, user := range usrList {
		userListResponse.Data.Users = append(userListResponse.Data.Users, &User{
			ID:         user.ID,
			Name:       user.Email,
			Role:       user.Role,
			CreateTime: user.CreateTime,
			UpdateTime: user.UpdateTime,
		})
	}

	c.Data["json"] = errors.OK().WithData(userListResponse)
	c.ServeJSON()
}
