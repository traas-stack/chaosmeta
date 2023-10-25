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

type MeasureRangeInstance struct {
	Id                       int    `json:"id" orm:"pk;auto;column(id)"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid" orm:"index;column(workflow_node_instance_uuid);size(64)"`
	JudgeValue               string `json:"judgeValue" orm:"column(judge_value);size(255)"`
	JudgeType                string `json:"judgeType" orm:"column(judge_type);size(64)"`
	FailedCount              string `json:"failedCount" orm:"column(failed_Count);size(32)"`
	SuccessCount             string `json:"successCount" orm:"column(success_count);size(32)"`
	Interval                 string `json:"interval" orm:"column(interval);size(32)"`
	Duration                 string `json:"duration" orm:"column(duration);size(32)"`
	MeasureType              string `json:"measureType" orm:"column(measure_type);size(32)"`
	ExecLog                  string `json:"exec_log" orm:"column(exec_log);size(2048)"`
	Status                   string `json:"status" orm:"column(status);size(32);index"`
	Message                  string `json:"message" orm:"column(message);size(1024)"`
	models.BaseTimeModel
}

func (er *MeasureRangeInstance) TableName() string {
	return TablePrefix + "measure_range_instance"
}

func CreateMeasureRangeInstance(m *MeasureRangeInstance) error {
	_, err := models.GetORM().Insert(m)
	return err
}

func UpdateMeasureRangeInstance(m *MeasureRangeInstance) error {
	_, err := models.GetORM().Update(m)
	return err
}

func DeleteMeasureRangeInstanceByID(id int) error {
	measureRangeInstance := &MeasureRangeInstance{Id: id}
	_, err := models.GetORM().Delete(measureRangeInstance)
	return err
}

func GetMeasureRangeInstanceById(m *MeasureRangeInstance) error {
	return models.GetORM().Read(m)
}

func GetMeasureRangeInstancesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) (*MeasureRangeInstance, error) {
	var measureRangeInstance MeasureRangeInstance
	err := models.GetORM().QueryTable(new(MeasureRangeInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("id").One(&measureRangeInstance)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &measureRangeInstance, nil
}

func ClearMeasureRangeInstancesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(MeasureRangeInstance).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func BatchSearchMeasureRangeInstances(searchCriteria map[string]interface{}) ([]*MeasureRangeInstance, error) {
	measureRangeInstances := []*MeasureRangeInstance{}
	qs := models.GetORM().QueryTable(new(MeasureRangeInstance))
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&measureRangeInstances)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return measureRangeInstances, nil
}
