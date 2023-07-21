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
	"chaosmeta-platform/pkg/models/common"
	"context"
)

type Namespace struct {
	Id          int    `orm:"auto; column(id)"`
	Name        string `json:"name" orm:"column(name); size(255);index"`
	Description string `json:"description" orm:"column(description); size(1024)"`
	Creator     int    `json:"creator" orm:"column(creator); index"`
	//User        []*User `json:"users" orm:"reverse(many)"`
	//Members     []*UserNamespace `json:"members" orm:"reverse(many)"`
	models.BaseTimeModel
}

func (u *Namespace) TableName() string {
	return "namespace"
}

func InsertNamespace(ctx context.Context, namespace *Namespace) (int64, error) {
	id, err := models.GetORM().Insert(namespace)
	return id, err
}

func UpdateNamespace(ctx context.Context, namespace *Namespace) (int64, error) {
	num, err := models.GetORM().Update(namespace)
	return num, err
}

func DeleteNamespace(ctx context.Context, id int) (int64, error) {
	num, err := models.GetORM().Delete(&Namespace{Id: id})
	return num, err
}

func GetNamespace(ctx context.Context, namespace *Namespace) error {
	return models.GetORM().Read(namespace)
}

func GetAllNamespaces() ([]*Namespace, error) {
	var namespaces []*Namespace
	_, err := models.GetORM().QueryTable("namespace").All(&namespaces)
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}
