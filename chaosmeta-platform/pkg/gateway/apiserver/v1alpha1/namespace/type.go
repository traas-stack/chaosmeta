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
	"chaosmeta-platform/pkg/models/namespace"
	namespaceService "chaosmeta-platform/pkg/service/namespace"
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

type QueryNamespaceResponse struct {
	Page       int                              `json:"page"`
	PageSize   int                              `json:"pageSize"`
	Total      int64                            `json:"total"`
	NameSpaces []namespaceService.NamespaceData `json:"namespaces"`
}

type UpdateNamespaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetNamespaceResponse struct {
	NameSpace namespace.Namespace `json:"namespace"`
}

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	CreateTime string `json:"create_time"`
}

type UserNamespace struct {
	User
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
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LabelCreateResponse struct {
	Id interface{} `json:"id"`
}

type LabelListResponse struct {
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Total    int64             `json:"total"`
	Labels   []namespace.Label `json:"labels"`
}

type SetAttackableClusterRequest struct {
	ClusterID int `json:"cluster_id"`
}

type ClusterNamespaceInfo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type GetAttackableClusterResponse struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Total    int64                  `json:"total"`
	Clusters []ClusterNamespaceInfo `json:"clusters,omitempty"`
}

type GetNamespaceListResponse struct {
	Page       int                           `json:"page"`
	PageSize   int                           `json:"pageSize"`
	Total      int64                         `json:"total"`
	Namespaces []namespace.UserNamespaceData `json:"namespaces,omitempty"`
}
