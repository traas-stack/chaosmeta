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

type FlowRange struct {
	Id                       int    `json:"id,omitempty" orm:"pk;auto;column(id)"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid,omitempty" orm:"index;column(workflow_node_instance_uuid);size(64)"`
	Source                   string `json:"source" orm:"column(source);size(32)"`
	Parallelism              string `json:"parallelism" orm:"column(parallelism);size(32)"`
	Duration                 string `json:"duration" orm:"column(duration);size(32)"`
	FlowType                 string `json:"flowType" orm:"column(flow_type);size(32)"`
	models.BaseTimeModel
}

func (er *FlowRange) TableName() string {
	return TablePrefix + "flow_range"
}

func CreateFlowRange(flowRange *FlowRange) error {
	_, err := models.GetORM().Insert(flowRange)
	return err
}

func UpdateFlowRange(flowRange *FlowRange) error {
	_, err := models.GetORM().Update(flowRange)
	return err
}

func DeleteFlowRangeByID(id int) error {
	flowRange := &FlowRange{Id: id}
	_, err := models.GetORM().Delete(flowRange)
	return err
}

func GetFlowRangeByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) (*FlowRange, error) {
	var flowRange FlowRange
	err := models.GetORM().QueryTable(new(FlowRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).One(&flowRange)
	if err != nil {
		return nil, err
	}

	return &flowRange, nil
}

func ListFlowRangesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) ([]*FlowRange, error) {
	flowRanges := []*FlowRange{}
	_, err := models.GetORM().QueryTable(new(FlowRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("id").All(&flowRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flowRanges, nil
}

func ClearFlowRangesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(FlowRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func BatchSearchFlowRanges(searchCriteria map[string]interface{}) ([]*FlowRange, error) {
	flowRanges := []*FlowRange{}
	qs := models.GetORM().QueryTable(new(FlowRange))
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&flowRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flowRanges, nil
}
