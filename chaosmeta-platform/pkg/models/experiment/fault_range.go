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
	"github.com/beego/beego/v2/client/orm"
)

type FaultRange struct {
	Id                       int    `json:"id,omitempty" orm:"pk;auto;column(id)"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid,omitempty" orm:"index;column(workflow_node_instance_uuid);size(64)"`
	TargetName               string `json:"target_name" orm:"column(target_name);size(255)"`
	TargetIP                 string `json:"target_ip" orm:"column(target_ip);size(32)"`
	TargetHostname           string `json:"target_hostname" orm:"column(target_hostname);size(255)"`
	TargetLabel              string `json:"target_label" orm:"column(target_label);size(1024)"`
	TargetApp                string `json:"target_app" orm:"column(target_app);size(255)"`
	TargetNamespace          string `json:"target_namespace" orm:"column(target_namespace);size(255)"`
	RangeType                string `json:"range_type" orm:"column(range_type);size(32)"`
	models.BaseTimeModel
}

func (er *FaultRange) TableName() string {
	return TablePrefix + "fault_range"
}

func CreateFaultRange(faultRange *FaultRange) error {
	_, err := models.GetORM().Insert(faultRange)
	return err
}

func UpdateFaultRange(faultRange *FaultRange) error {
	_, err := models.GetORM().Update(faultRange)
	return err
}

func DeleteFaultRangeByID(id int) error {
	faultRange := &FaultRange{Id: id}
	_, err := models.GetORM().Delete(faultRange)
	return err
}

func GetFaultRangeByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) (*FaultRange, error) {
	var faultRange FaultRange
	err := models.GetORM().QueryTable(new(FaultRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).One(&faultRange)
	if err != nil {
		return nil, err
	}

	return &faultRange, nil
}

func ListFaultRangesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) ([]*FaultRange, error) {
	faultRanges := []*FaultRange{}
	_, err := models.GetORM().QueryTable(new(FaultRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("id").All(&faultRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return faultRanges, nil
}

func ClearFaultRangesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(FaultRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func BatchSearchFaultRanges(searchCriteria map[string]interface{}) ([]*FaultRange, error) {
	faultRanges := []*FaultRange{}
	qs := models.GetORM().QueryTable(new(FaultRange))
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&faultRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return faultRanges, nil
}
