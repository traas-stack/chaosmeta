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
	"time"
)

type Namespace struct {
	ID          int        `orm:"column(id);auto"`
	Name        string     `orm:"column(name)"`
	CreatorId   int        `orm:"column(creator_id)"`
	Description string     `orm:"column(description)"`
	CreateTime  *time.Time `orm:"column(create_time);auto_now_add;type(datetime)"`
	UpdateTime  *time.Time `orm:"column(update_time);auto_now;type(datetime)"`
}

func InsertNamespace(ctx context.Context, n *Namespace) error {
	_, err := GetORM().Insert(n)
	return err
}
