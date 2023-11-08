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
	"chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

const TimeLayout = "2006-01-02 15:04:05"

type Label struct {
	Id          int    `json:"id" orm:"pk;auto;column(id)"`
	Name        string `json:"name" orm:"column(name);size(255);index"`
	Color       string `json:"color" orm:"column(color);size(255);index"`
	NamespaceId int    `json:"namespaceId" orm:"column(namespace_id);index"`
	Creator     string `json:"creator" orm:"column(creator); size(255); index"`
	models.BaseTimeModel
}

func (l *Label) TableName() string {
	return "namespace_label"
}

func (l *Label) TableUnique() [][]string {
	return [][]string{
		{"name", "namespace_id"},
	}
}

func InsertLabel(ctx context.Context, label *Label) (int64, error) {
	if label == nil {
		return 0, errors.New("label is nil")
	}
	id, err := models.GetORM().Insert(label)
	return id, err
}

func UpdateLabel(ctx context.Context, label *Label) (int64, error) {
	if label == nil {
		return 0, errors.New("label is nil")
	}
	num, err := models.GetORM().Update(label)
	return num, err
}

func GetLabelById(ctx context.Context, label *Label) error {
	if label == nil {
		return errors.New("label is nil")
	}
	return models.GetORM().Read(label)
}

func GetLabelByIdAndNamespaceId(ctx context.Context, label *Label) error {
	return models.GetORM().Read(label, "id", "namespace_id")
}

func GetLabelByName(ctx context.Context, label *Label) error {
	return models.GetORM().Read(label, "name", "namespace_id")
}

func DeleteLabel(ctx context.Context, id int) (int64, error) {
	num, err := models.GetORM().Delete(&Label{Id: id})
	return num, err
}

func QueryLabels(ctx context.Context, nameSpaceId int, name, creator string, orderBy string, page, pageSize int) (int64, []Label, error) {
	label, labelList := Label{}, new([]Label)
	querySeter := models.GetORM().QueryTable(label.TableName())
	labelQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	labelQuery.Filter("namespace_id", models.NEGLECT, false, nameSpaceId)
	if len(name) > 0 {
		labelQuery.Filter("name", models.CONTAINS, true, name)
	}

	if len(creator) > 0 {
		labelQuery.Filter("creator", models.CONTAINS, true, creator)
	}

	totalCount, err := labelQuery.GetOamQuerySeter().Count()
	if err != nil {
		return 0, nil, err
	}

	orderByList := []string{}
	if orderBy != "" {
		orderByList = append(orderByList, orderBy)
	} else {
		orderByList = append(orderByList, "id")
	}
	labelQuery.OrderBy(orderByList...)
	if err := labelQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = labelQuery.GetOamQuerySeter().All(labelList)
	if err == orm.ErrNoRows {
		return 0, nil, nil
	}
	return totalCount, *labelList, err
}
