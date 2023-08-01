package namespace

import (
	"chaosmeta-platform/pkg/models/namespace"
	"time"
)

type CreateNamespaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateNamespaceResponse struct {
	ID int64 `json:"id"`
}

type NameSpace struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role" `
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type ListNamespaceResponse struct {
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	Total      int64                 `json:"total"`
	NameSpaces []namespace.Namespace `json:"namespaces"`
}

type UpdateNamespaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetNamespaceResponse struct {
	NameSpace namespace.Namespace `json:"namespace"`
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role" `
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type UserListResponse struct {
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	Total    int64   `json:"total"`
	Users    []*User `json:"users"`
}

type AddUsersRequest struct {
	Users []struct {
		Id         int `json:"id"`
		Permission int `json:"permission"`
	} `json:"users"`
}

type RemoveUsersRequest struct {
	UserIds []int `json:"user_ids"`
}

type ChangeUsersPermissionRequest struct {
	UserIds    []int `json:"user_ids"`
	Permission int   `json:"permission"`
}

type LabelCreateRequest struct {
	Name string `json:"name"`
}

type LabelCreateResponse struct {
	Id interface{} `json:"id"`
}
