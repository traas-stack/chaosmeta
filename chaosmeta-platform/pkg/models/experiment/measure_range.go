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

type MeasureRange struct {
	Id                       int    `json:"id,omitempty" orm:"pk;auto;column(id)"`
	WorkflowNodeInstanceUUID string `json:"workflow_node_instance_uuid,omitempty" orm:"index;column(workflow_node_instance_uuid);size(64)"`
	JudgeValue               string `json:"judgeValue" orm:"column(judge_value);size(255)"`
	JudgeType                string `json:"judgeType" orm:"column(judge_type);size(64)"`
	FailedCount              string `json:"failedCount" orm:"column(failed_Count);size(32)"`
	SuccessCount             string `json:"successCount" orm:"column(success_count);size(32)"`
	Interval                 string `json:"interval" orm:"column(interval);size(32)"`
	Duration                 string `json:"duration" orm:"column(duration);size(32)"`
	MeasureType              string `json:"measureType" orm:"column(measure_type);size(32)"`
	models.BaseTimeModel
}

func (er *MeasureRange) TableName() string {
	return TablePrefix + "measure_range"
}

func CreateMeasureRange(measureRange *MeasureRange) error {
	_, err := models.GetORM().Insert(measureRange)
	return err
}

func UpdateMeasureRange(measureRange *MeasureRange) error {
	_, err := models.GetORM().Update(measureRange)
	return err
}

func DeleteMeasureRangeByID(id int) error {
	measureRange := &MeasureRange{Id: id}
	_, err := models.GetORM().Delete(measureRange)
	return err
}

func GetMeasureRangeByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) (*MeasureRange, error) {
	var measureRange MeasureRange
	err := models.GetORM().QueryTable(new(MeasureRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).One(&measureRange)
	if err != nil {
		return nil, err
	}

	return &measureRange, nil
}

func ListMeasureRangeByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) ([]*MeasureRange, error) {
	measureRanges := []*MeasureRange{}
	_, err := models.GetORM().QueryTable(new(MeasureRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).OrderBy("id").All(&measureRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return measureRanges, nil
}

func ClearMeasureRangesByWorkflowNodeInstanceUUID(workflowNodeInstanceUUID string) error {
	_, err := models.GetORM().QueryTable(new(MeasureRange).TableName()).Filter("workflow_node_instance_uuid", workflowNodeInstanceUUID).Delete()
	return err
}

func BatchSearchMeasureRanges(searchCriteria map[string]interface{}) ([]*MeasureRange, error) {
	measureRanges := []*MeasureRange{}
	qs := models.GetORM().QueryTable(new(MeasureRange))
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&measureRanges)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return measureRanges, nil
}
