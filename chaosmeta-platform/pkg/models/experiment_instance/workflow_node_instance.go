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
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type WorkflowNodeInstance struct {
	UUID                   string `json:"uuid,omitempty" orm:"column(uuid);pk"`
	ExperimentInstanceUUID string `json:"experiment_instance_uuid" orm:"index;column(experiment_instance_uuid);size(64)"`
	Row                    int    `json:"row" orm:"column(row)"`
	Column                 int    `json:"column" orm:"column(column)"`
	Duration               string `json:"duration" orm:"column(duration);size(32)"`
	ScopeId                int    `json:"scope_id" orm:"column(scope_id); int(11)"`
	TargetId               int    `json:"target_id" orm:"column(target_id); int(11)"`
	ExecType               string `json:"exec_type" orm:"column(exec_type);size(32)"`
	ExecID                 int    `json:"exec_id" orm:"column(exec_id)"`
	Status                 string `json:"status" orm:"column(status);size(32);default(to_be_executed);index"`
	Message                string `json:"message" orm:"column(message);size(1024)"`
	Version                int    `json:"-" orm:"column(version);default(0);column(version)"`
	models.BaseTimeModel
}

func (wn *WorkflowNodeInstance) TableName() string {
	return TablePrefix + "workflow_node_instance"
}

func GetWorkflowNodeInstanceByUUID(uuid string) (*WorkflowNodeInstance, error) {
	workflowNode := &WorkflowNodeInstance{UUID: uuid}
	err := models.GetORM().Read(workflowNode, "uuid")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNode, nil
}

func UpdateWorkflowNodeInstance(workflowNodeInstance *WorkflowNodeInstance) error {
	o := models.GetORM()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	existing := WorkflowNodeInstance{UUID: workflowNodeInstance.UUID}
	err = tx.Read(&existing)
	if err != nil {
		tx.Rollback()
		return err
	}

	if workflowNodeInstance.Version != existing.Version {
		tx.Rollback()
		return errors.New("Concurrent modification detected")
	}

	workflowNodeInstance.Version = existing.Version + 1
	if _, err = tx.Update(workflowNodeInstance); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func UpdateWorkflowNodeInstanceStatus(uuid string, status, message string) error {
	nodeInstance, err := GetWorkflowNodeInstanceByUUID(uuid)
	if err != nil {
		return err
	}
	nodeInstance.Status = status
	if message != "" {
		nodeInstance.Message = message
	}
	return UpdateWorkflowNodeInstance(nodeInstance)
}

func GetWorkflowNodeInstancesByExperimentUUID(experimentUUID string) ([]*WorkflowNodeInstance, error) {
	workflowNodes := []*WorkflowNodeInstance{}
	_, err := models.GetORM().QueryTable(new(WorkflowNodeInstance).TableName()).Filter("experiment_instance_uuid", experimentUUID).OrderBy("row", "column").All(&workflowNodes)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNodes, nil
}

func CreateWorkflowNodeInstance(workflowNode *WorkflowNodeInstance) error {
	_, err := models.GetORM().Insert(workflowNode)
	return err
}

func DeleteWorkflowNodeInstanceByUUID(uuid string) error {
	workflowNode := &WorkflowNodeInstance{UUID: uuid}
	_, err := models.GetORM().Delete(workflowNode)
	return err
}

func BatchSearchWorkflowNodeInstances(searchCriteria map[string]interface{}) ([]*WorkflowNodeInstance, error) {
	o := models.GetORM()
	workflowNodes := []*WorkflowNodeInstance{}
	qs := o.QueryTable(new(WorkflowNodeInstance).TableName())
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&workflowNodes)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNodes, nil
}
