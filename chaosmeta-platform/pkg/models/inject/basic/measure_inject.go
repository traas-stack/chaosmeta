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
	"github.com/beego/beego/v2/client/orm"
)

type MeasureInject struct {
	Id            int    `json:"id" orm:"pk;auto;column(id)"`
	MeasureType   string `json:"measure_type" orm:"column(measure_type);size(32)"`
	Name          string `json:"name" orm:"column(name);size(255)"`
	NameCn        string `json:"name_cn" orm:"column(name_cn);size(255)"`
	Description   string `json:"description" orm:"column(description);size(1024)"`
	DescriptionCn string `json:"description_cn" orm:"column(description_cn);size(1024)"`
}

func (m *MeasureInject) TableName() string {
	return TablePrefix + "measure_inject"
}

func GetMeasureInjectByID(id int) (*MeasureInject, error) {
	measureInject := &MeasureInject{Id: id}
	err := models.GetORM().Read(measureInject)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, nil
	} else {
		return measureInject, err
	}
}

// ListMeasureInjects retrieves all measure_injects
func ListMeasureInjects(orderBy string, page, pageSize int) (int64, []MeasureInject, error) {
	measureInject, measureInjects := MeasureInject{}, new([]MeasureInject)
	querySeter := models.GetORM().QueryTable(measureInject.TableName())
	measureInjectQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	var totalCount int64
	totalCount, err = measureInjectQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	measureInjectQuery.OrderBy(orderByList...)
	if err := measureInjectQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = measureInjectQuery.GetOamQuerySeter().All(measureInjects)
	return totalCount, *measureInjects, err
}

// InsertMeasureInject insert a new measure_inject
func InsertMeasureInject(measureInject *MeasureInject) error {
	_, err := models.GetORM().Insert(measureInject)
	return err
}

func UpdateMeasureInject(measureInject *MeasureInject) error {
	_, err := models.GetORM().Update(measureInject)
	return err
}

// DeleteMeasureInject deletes a measure_inject by its ID
func DeleteMeasureInject(id int) error {
	measureInject := MeasureInject{Id: id}
	_, err := models.GetORM().Delete(&measureInject)
	return err
}
