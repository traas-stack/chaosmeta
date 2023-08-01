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
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

const (
	AdminRole  = "admin"
	NormalRole = "normal"
)

type User struct {
	ID            int       `json:"id" orm:"pk;auto;column(id)"`
	Email         string    `json:"email" orm:"unique;index;column(email);size(255)"`
	Password      string    `json:"password" orm:"column(password);size(255)"`
	Role          string    `json:"role" orm:"index; column(role);size(32)"`
	Token         string    `json:"token" orm:"column(token);size(255)"`
	Disabled      bool      `json:"disabled" orm:"column(disabled)"`
	IsDeleted     bool      `json:"isDeleted" orm:"column(is_deleted);default(0)"`
	LastLoginTime time.Time `json:"lastLoginTime" orm:"column(last_login_time);auto_now;type(datetime)"`
	//Namespace     []*Namespace `json:"namespaces" orm:"rel(m2m);rel_through(chaosmeta-platform/pkg/models.UserNamespace);on_delete(cascade)"`
	models.BaseTimeModel
}

func (u *User) TableName() string {
	return "user"
}

func InsertUser(ctx context.Context, u *User) (int64, error) {
	return models.GetORM().Insert(u)
}

func DeleteUsersByIdList(ctx context.Context, userId []int) error {
	user := User{}
	querySeter := models.GetORM().QueryTable(user.TableName())
	userQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return err
	}
	userQuery.Filter("id", models.IN, false, userId)
	_, err = userQuery.Update(orm.Params{
		"is_deleted": true,
	})
	return err
}

func UpdateUser(ctx context.Context, u *User) error {
	suc, err := models.GetORM().Update(u)
	if suc == 0 {
		return fmt.Errorf("record[email: %s] not found", u.Email)
	}

	return err
}

func UpdateUsersRole(ctx context.Context, userId []int, role string) error {
	user := User{}
	querySeter := models.GetORM().QueryTable(user.TableName())
	userQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return err
	}
	userQuery.Filter("id", models.IN, false, userId)
	_, err = userQuery.Update(orm.Params{
		"role": role,
	})
	return err
}

func GetUser(ctx context.Context, u *User) error {
	return models.GetORM().Read(u, "email")
}

func QueryUser(ctx context.Context, name, role, orderBy string, page, pageSize int) (int64, []User, error) {
	u, users := User{}, new([]User)
	querySeter := models.GetORM().QueryTable(u.TableName())
	userQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	userQuery.Filter("is_deleted", models.NEGLECT, false, false)
	if len(name) > 0 {
		userQuery.Filter("email", models.CONTAINS, true, name)
	}

	if len(role) > 0 {
		userQuery.Filter("role", models.NEGLECT, false, role)
	}

	var totalCount int64
	totalCount, err = userQuery.GetOamQuerySeter().Count()
	if err := userQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}
	if len(orderBy) > 0 {
		userQuery.OrderBy(orderBy)
	}

	_, err = userQuery.GetOamQuerySeter().All(users)
	return totalCount, *users, err
}

func DeleteUsersByNameList(ctx context.Context, names []string) error {
	user := User{}
	querySeter := models.GetORM().QueryTable(user.TableName())
	userQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return err
	}
	userQuery.Filter("email", models.IN, false, names)
	_, err = userQuery.Delete()
	return err
}
