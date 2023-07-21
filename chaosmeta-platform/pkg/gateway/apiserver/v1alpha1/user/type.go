package user

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role" `
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type UserData struct {
	Total int64   `json:"total"`
	Users []*User `json:"users"`
}

type UserListResponse struct {
	v1alpha1.ResponseData
	Data UserData `json:"data"`
}

type SingleUserResponse struct {
	v1alpha1.ResponseData
	Data struct {
		User *User `json:"user"`
	} `json:"data"`
}

type UserCreateResponse struct {
	v1alpha1.ResponseData
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

type UserLoginResponse struct {
	v1alpha1.ResponseData
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type UserDeleteRequest struct {
	UserIDs []int `json:"user_ids"`
}

type UserDeleteResponse struct {
	v1alpha1.ResponseData
}

type UserChangePasswordResponse struct {
	v1alpha1.ResponseData
}

type UserChangeRoleResponse struct {
	v1alpha1.ResponseData
}
