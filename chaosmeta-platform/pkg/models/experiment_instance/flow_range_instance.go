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

type FlowRangeInstance struct {
	Id                       int    `json:"id" orm:"pk;auto;column(id)"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid" orm:"index;column(workflow_node_instance_uuid);size(64)"`
	Source                   string `json:"source" orm:"column(source);size(32)"`
	Parallelism              string `json:"parallelism" orm:"column(parallelism);size(32)"`
	Duration                 string `json:"duration" orm:"column(duration);size(32)"`
	FlowType                 string `json:"flowType" orm:"column(flow_type);size(32)"`
	ExecLog                  string `json:"exec_log" orm:"column(exec_log);size(2048)"`
	Status                   string `json:"status" orm:"column(status);size(32);index"`
	Message                  string `json:"message" orm:"column(message);size(1024)"`
	models.BaseTimeModel
}

func (er *FlowRangeInstance) TableName() string {
	return TablePrefix + "flow_range_instance"
}

func CreateFlowRangeInstance(f *FlowRangeInstance) error {
	_, err := models.GetORM().Insert(f)
	return err
}

func UpdateFlowRangeInstance(f *FlowRangeInstance) error {
	_, err := models.GetORM().Update(f)
	return err
}

func DeleteFlowRangeInstanceByID(id int) error {
	flowRangeInstance := &FlowRangeInstance{Id: id}
	_, err := models.GetORM().Delete(flowRangeInstance)
	return err
}

func GetFlowRangeInstanceById(f *FlowRangeInstance) error {
	return models.GetORM().Read(f)
}

func GetFlowRangeInstancesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) (*FlowRangeInstance, error) {
	var flowRangeInstance FlowRangeInstance
	err := models.GetORM().QueryTable(new(FlowRangeInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("id").One(&flowRangeInstance)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &flowRangeInstance, nil
}

func ClearFlowRangeInstancesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(FlowRangeInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func BatchSearchFlowRangeInstances(searchCriteria map[string]interface{}) ([]*FlowRangeInstance, error) {
	flowRangeInstances := []*FlowRangeInstance{}
	qs := models.GetORM().QueryTable(new(FlowRangeInstance))
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&flowRangeInstances)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flowRangeInstances, nil
}
