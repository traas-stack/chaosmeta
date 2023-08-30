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

package experiment

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/util/log"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type ArgsValue struct {
	Id               int    `json:"id" orm:"pk;auto;column(id)"`
	ArgsID           int    `json:"args_id" orm:"column(args_id);index"`
	WorkflowNodeUUID string `json:"workflow_node_uuid,omitempty" orm:"column(workflow_node_uuid);index"`
	Value            string `json:"value" orm:"column(value);size(1024)"`
	models.BaseTimeModel
}

func (av *ArgsValue) TableName() string {
	return TablePrefix + "args_value"
}

func (av *ArgsValue) TableUnique() [][]string {
	return [][]string{{"args_id", "workflow_node_uuid"}}
}

func InsertArgsValue(arg *ArgsValue) error {
	_, err := models.GetORM().Insert(arg)
	return err
}

func BatchInsertArgsValues(workflowNodeUUID string, argsValues []*ArgsValue) error {
	o := models.GetORM()
	if workflowNodeUUID == "" {
		return fmt.Errorf("workflowNodeUUID is empty")
	}
	if argsValues == nil {
		return fmt.Errorf("argsValues is nil")
	}

	oldValues := []*ArgsValue{}
	_, err := o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).All(&oldValues)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(oldValues) > 0 {
		if _, err = o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).Delete(); err != nil {
			log.Error(err)
			return err
		}
	}
	for _, argsValue := range argsValues {
		if argsValue == nil {
			return fmt.Errorf("argsValue is nil")
		}
		argsValue.WorkflowNodeUUID = workflowNodeUUID
		if err := InsertArgsValue(argsValue); err != nil {
			return err
		}
	}
	return nil
}

func ClearArgsValuesByWorkflowNodeUUID(workflowNodeUUID string) error {
	o := models.GetORM()
	_, err := o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).Delete()
	return err
}

func GetArgsValuesByWorkflowNodeUUID(workflowNodeUUID string) ([]*ArgsValue, error) {
	o := models.GetORM()
	var argsValues []*ArgsValue
	_, err := o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).OrderBy("-create_time").All(&argsValues)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return argsValues, nil
}

func BatchSearchArgsValues(searchCriteria map[string]interface{}) ([]*ArgsValue, error) {
	o := models.GetORM()
	argsValues := []*ArgsValue{}
	qs := o.QueryTable(new(ArgsValue).TableName())
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
