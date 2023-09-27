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
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type WorkflowNode struct {
	Id             int    `json:"id,omitempty" orm:"pk;auto;column(id)"`
	Name           string `json:"name" orm:"index;column(name);size(255)"`
	UUID           string `json:"uuid,omitempty" orm:"column(uuid);index"`
	ExperimentUUID string `json:"experiment_uuid" orm:"index;column(experiment_uuid);size(64)"`
	Row            int    `json:"row" orm:"column(row)"`
	Column         int    `json:"column" orm:"column(column)"`
	Duration       string `json:"duration" orm:"column(duration);size(32)"`
	ScopeId        int    `json:"scope_id" orm:"column(scope_id); int(11)"`
	TargetId       int    `json:"target_id" orm:"column(target_id); int(11)"`
	ExecType       string `json:"exec_type" orm:"column(exec_type);size(32)"`
	ExecID         int    `json:"exec_id" orm:"column(exec_id); int(11)"`
	Version        int    `json:"-" orm:"column(version);default(0);version"`
	models.BaseTimeModel
}

func (wn *WorkflowNode) TableName() string {
	return TablePrefix + "workflow_node"
}

func (wn *WorkflowNode) TableUnique() [][]string {
	return [][]string{{"uuid", "experiment_uuid"}}
}

func GetWorkflowNodeByUUID(uuid string) (*WorkflowNode, error) {
	workflowNode := &WorkflowNode{UUID: uuid}
	err := models.GetORM().Read(workflowNode, "uuid")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNode, nil
}

func GetWorkflowNodesByExperimentUUID(experimentUUID string) ([]*WorkflowNode, error) {
	workflowNodes := []*WorkflowNode{}
	_, err := models.GetORM().QueryTable(new(WorkflowNode).TableName()).Filter("experiment_uuid", experimentUUID).OrderBy("row", "column").All(&workflowNodes)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNodes, nil
}

func CreateWorkflowNode(workflowNode *WorkflowNode) error {
	_, err := models.GetORM().Insert(workflowNode)
	return err
}

func UpdateWorkflowNode(workflowNode *WorkflowNode) error {
	o := models.GetORM()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	existing := WorkflowNode{UUID: workflowNode.UUID}
	err = tx.Read(&existing)
	if err != nil {
		tx.Rollback()
		return err
	}

	if workflowNode.Version != existing.Version {
		tx.Rollback()
		return errors.New("Concurrent modification detected")
	}

	workflowNode.Version = existing.Version + 1
	if _, err = tx.Update(workflowNode); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DeleteWorkflowNodeByUUID(uuid string) error {
	_, err := models.GetORM().QueryTable(new(WorkflowNode).TableName()).Filter("uuid", uuid).Delete()
	return err
}

func DeleteWorkflowNodeByExperimentUUID(uuid string) error {
	_, err := models.GetORM().QueryTable(new(WorkflowNode).TableName()).Filter("experiment_uuid", uuid).Delete()
	return err
}

// BatchSearchWorkflowNodes 批量搜索workflow_nodes
func BatchSearchWorkflowNodes(searchCriteria map[string]interface{}) ([]*WorkflowNode, error) {
	o := models.GetORM()
	workflowNodes := []*WorkflowNode{}
	qs := o.QueryTable(new(WorkflowNode).TableName())
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
