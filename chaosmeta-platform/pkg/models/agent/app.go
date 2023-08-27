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

package agent

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
)

type App struct {
	ID   int    `json:"id" orm:"pk;auto;column(id)"`
	Name string `json:"name" orm:"unique;index;column(name);size(255)"`
	models.BaseTimeModel
}

func (a *App) TableName() string {
	return "agent_app"
}

func InsertApp(ctx context.Context, a *App) (int64, error) {
	return models.GetORM().Insert(a)
}

func GetAppByName(ctx context.Context, a *App) error {
	return models.GetORM().Read(a, "name")
}

func GetAppById(ctx context.Context, a *App) error {
	return models.GetORM().Read(a)
}

//func DeleteUsersByIdList(ctx context.Context, userId []int) error {
//	user := User{}
//	querySeter := models.GetORM().QueryTable(user.TableName())
//	userQuery, err := models.NewDataSelectQuery(&querySeter)
//	if err != nil {
//		return err
//	}
//	userQuery.Filter("id", models.IN, false, userId)
//	_, err = userQuery.Update(orm.Params{
//		"is_deleted": true,
//	})
//	return err
//}
//
//func UpdateUser(ctx context.Context, u *User) error {
//	suc, err := models.GetORM().Update(u)
//	if suc == 0 {
//		return fmt.Errorf("record[email: %s] not found", u.Email)
//	}
//
//	return err
//}
