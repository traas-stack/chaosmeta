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

package models

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID         int64      `orm:"column(id);auto"`
	Name       string     `orm:"column(name)"`
	Password   string     `orm:"column(password)"`
	Role       string     `orm:"column(role)"`
	CreateTime *time.Time `orm:"column(create_time);auto_now_add;type(datetime)"`
	UpdateTime *time.Time `orm:"column(update_time);auto_now;type(datetime)"`
}

const (
	AdminRole  = "admin"
	NormalRole = "normal"
)

func InsertUser(ctx context.Context, u *User) (int64, error) {
	return GetORM().Insert(u)
}

func DeleteUser(ctx context.Context, id int64) error {
	_, err := GetORM().Delete(&User{ID: id})
	return err
}

func UpdateUserPasswd(ctx context.Context, id int64, passwd string) error {
	suc, err := GetORM().Update(&User{ID: id, Password: passwd}, "password")
	if suc == 0 {
		return fmt.Errorf("record[id: %d] not found", id)
	}

	return err
}

func UpdateUserRole(ctx context.Context, id int64, role string) error {
	suc, err := GetORM().Update(&User{ID: id, Role: role}, "role")
	if suc == 0 {
		return fmt.Errorf("record[id: %d] not found", id)
	}

	return err
}

func QueryUser(ctx context.Context) error {
	// TODO
	return nil
}
