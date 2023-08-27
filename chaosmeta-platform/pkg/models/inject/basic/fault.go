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

type Fault struct {
	ID            int    `json:"id" orm:"pk;auto;column(id)"`
	TargetId      int    `json:"targetId" orm:"column(target_id)"`
	Name          string `json:"name" orm:"size(255);column(name)"`
	NameCn        string `json:"nameCn" orm:"size(255);column(name_cn)"`
	Description   string `json:"description" orm:"size(1024);column(description)"`
	DescriptionCn string `json:"descriptionCn" orm:"size(1024);column(description_cn)"`
	models.BaseTimeModel
}

func (f *Fault) TableName() string {
	return TablePrefix + "fault"
}

func InsertFault(ctx context.Context, fault *Fault) error {
	_, err := models.GetORM().Insert(fault)
	return err
}

func InsertFaultsMulti(ctx context.Context, faultList []*Fault) error {
	_, err := models.GetORM().InsertMulti(len(faultList), faultList)
	return err
}

func DeleteFault(ctx context.Context, id int) error {
	fault := &Fault{ID: id}
	_, err := models.GetORM().Delete(fault)
	return err
}

func UpdateFault(ctx context.Context, fault *Fault) error {
	if models.GetORM().Read(fault) == nil {
		_, err := models.GetORM().Update(fault)
		return err
	}
	return errors.New("fault not found")
}

func GetFaultById(ctx context.Context, id int) (*Fault, error) {
	o := models.GetORM()

	fault := &Fault{ID: id}
	err := o.Read(fault)

	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, nil
	} else {
		return fault, err
	}
}

func ListFaults(ctx context.Context, targetId int, orderBy string, page, pageSize int) (int64, []Fault, error) {
	fault, faults := Fault{}, new([]Fault)

	querySeter := models.GetORM().QueryTable(fault.TableName())
	faultQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	if targetId > 0 {
		faultQuery.Filter("target_id", models.NEGLECT, false, targetId)
	}

	var totalCount int64
	totalCount, err = faultQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	faultQuery.OrderBy(orderByList...)
	if err := faultQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = faultQuery.GetOamQuerySeter().All(faults)
	return totalCount, *faults, err
}
