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

import models "chaosmeta-platform/pkg/models/common"

type FlowInject struct {
	Id            int    `json:"id" orm:"column(id);pk"`
	Name          string `json:"name" orm:"column(name);size(255)"`
	NameCn        string `json:"name_cn" orm:"column(name_cn);size(255)"`
	Description   string `json:"description" orm:"column(description);size(1024)"`
	DescriptionCn string `json:"description_cn" orm:"column(description_cn);size(1024)"`
}

func (m *FlowInject) TableName() string {
	return TablePrefix + "flow_inject"
}

func GetFlowInjectByID(id int) (FlowInject, error) {
	FlowInject := FlowInject{Id: id}
	err := models.GetORM().Read(&FlowInject)
	return FlowInject, err
}

func ListFlowInjects(orderBy string, page, pageSize int) (int64, []FlowInject, error) {
	flowInject, flowInjects := FlowInject{}, new([]FlowInject)
	querySeter := models.GetORM().QueryTable(flowInject.TableName())
	FlowInjectQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	var totalCount int64
	totalCount, err = FlowInjectQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	FlowInjectQuery.OrderBy(orderByList...)
	if err := FlowInjectQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = FlowInjectQuery.GetOamQuerySeter().All(flowInjects)
	return totalCount, *flowInjects, err
}

func InsertFlowInject(FlowInject *FlowInject) error {
	_, err := models.GetORM().Insert(FlowInject)
	return err
}

func UpdateFlowInject(FlowInject *FlowInject) error {
	_, err := models.GetORM().Update(FlowInject)
	return err
}

// DeleteFlowInject deletes a Flow_inject by its ID
func DeleteFlowInject(id int) error {
	FlowInject := FlowInject{Id: id}
	_, err := models.GetORM().Delete(&FlowInject)
	return err
}
