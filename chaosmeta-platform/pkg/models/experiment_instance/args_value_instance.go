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

package experiment_instance

import (
	models "chaosmeta-platform/pkg/models/common"
	"github.com/beego/beego/v2/client/orm"
)

type ArgsValueInstance struct {
	Id                       int    `json:"id" orm:"pk;auto;column(id)"`
	ArgsID                   int    `json:"args_id" orm:"column(args_id);index"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid,omitempty" orm:"column(workflow_node_instance_uuid);index"`
	Value                    string `json:"value" orm:"column(value);size(1024)"`
	models.BaseTimeModel
}

func (a *ArgsValueInstance) TableUnique() [][]string {
	return [][]string{{"args_id", "workflow_node_instance_uuid"}}
}

func (av *ArgsValueInstance) TableName() string {
	return TablePrefix + "args_value_instance"
}

func BatchInsertArgsValueInstances(workflowNodeInstanceUUID string, argsValues []*ArgsValueInstance) error {
	o := models.GetORM()
	oldValues := []*ArgsValueInstance{}
	_, err := o.QueryTable(new(ArgsValueInstance)).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).All(&oldValues)
	if err != nil {
		return err
	}
	if len(oldValues) > 0 {
		if _, err = o.QueryTable(new(ArgsValueInstance)).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete(); err != nil {
			return err
		}
	}
	for _, argsValue := range argsValues {
		argsValue.WorkflowNodeInstanceUUID = workflowNodeInstanceUUID
		if _, err := o.Insert(argsValue); err != nil {
			return err
		}
	}
	return nil
}

func ClearArgsValueInstancesByWorkflowNodeUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(ArgsValueInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func GetArgsValueInstancesByWorkflowNodeUUID(workflowNodeInstanceUUID string) ([]*ArgsValueInstance, error) {
	var argsValues []*ArgsValueInstance
	_, err := models.GetORM().QueryTable(new(ArgsValueInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("-create_time").All(&argsValues)
	if err != nil {
		return nil, err
	}

	return argsValues, nil
}

func BatchSearchArgsValueInstances(searchCriteria map[string]interface{}) ([]*ArgsValueInstance, error) {
	argsValues := []*ArgsValueInstance{}
	qs := models.GetORM().QueryTable(new(ArgsValueInstance).TableName())
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&argsValues)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return argsValues, nil
}
