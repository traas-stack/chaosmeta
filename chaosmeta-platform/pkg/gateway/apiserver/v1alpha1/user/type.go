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
	namespace2 "chaosmeta-platform/pkg/models/namespace"
	"time"
)

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

type UserNamespace struct {
	User
	IsJoin     bool `json:"isJoin"`
	Permission int  `json:"permission"`
}

type UserListNamespaceResponse struct {
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int64            `json:"total"`
	Users    []*UserNamespace `json:"users"`
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginRequest UserCreateRequest

type UserLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UsersDeleteRequest struct {
	UserIDs []int `json:"user_ids"`
}

type UsersPasswordUpdateRequest struct {
	Password string `json:"password"`
}

type UserUpdateRoleRequest struct {
	UserIDs []int  `json:"user_ids"`
	Role    string `json:"role"`
}

type NameSpaceListResponse struct {
	Page       int                            `json:"page"`
	PageSize   int                            `json:"pageSize"`
	Total      int64                          `json:"total"`
	Namespaces []namespace2.UserNamespaceData `json:"namespaces"`
}
