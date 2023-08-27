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

package basic

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type Target struct {
	ID            int    `json:"id" orm:"pk;auto;column(id)"`
	ScopeId       int    `json:"scopeId" orm:"column(scope_id)"`
	Name          string `json:"name" orm:"size(255);column(name)"`
	NameCn        string `json:"nameCn" orm:"size(255);column(name_cn)"`
	Description   string `json:"description" orm:"size(1024);column(description)"`
	DescriptionCn string `json:"descriptionCn" orm:"size(1024);column(description_cn)"`
	models.BaseTimeModel
}

func (c *Target) TableName() string {
	return TablePrefix + "target"
}

func (c *Target) TableUnique() [][]string {
	return [][]string{{"scope_id", "name"}}
}

func InsertTarget(ctx context.Context, target *Target) error {
	_, err := models.GetORM().Insert(target)
	return err
}

func DeleteTarget(ctx context.Context, id int) error {
	target := &Target{ID: id}
	_, err := models.GetORM().Delete(target)
	return err
}

func UpdateTarget(ctx context.Context, target *Target) error {
	o := models.GetORM()
	if o.Read(target) == nil {
		_, err := o.Update(target)
		return err
	}
	return errors.New("target Not Found")
}

func GetTargetById(ctx context.Context, id int) (*Target, error) {
	o := models.GetORM()
	target := &Target{ID: id}
	err := o.Read(target)

	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, nil
	} else {
		return target, err
	}
}

func ListTargets(ctx context.Context, scopeId int, orderBy string, page, pageSize int) (int64, []Target, error) {
	target, targets := Target{}, new([]Target)

	querySeter := models.GetORM().QueryTable(target.TableName())
	scopeQuery, err := models.NewDataSelectQuery(&querySeter)

	if scopeId > 0 {
		scopeQuery.Filter("scope_id", models.NEGLECT, false, scopeId)
	}

	var totalCount int64
	totalCount, err = scopeQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	scopeQuery.OrderBy(orderByList...)
	if err := scopeQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = scopeQuery.GetOamQuerySeter().All(targets)
	return totalCount, *targets, err
}
