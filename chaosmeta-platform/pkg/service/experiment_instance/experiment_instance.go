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
	"chaosmeta-platform/pkg/models/experiment_instance"
	"chaosmeta-platform/pkg/models/namespace"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

type ExperimentInstanceService struct{}

type LabelInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
}

type ExperimentInstanceInfo struct {
	Uuid        string      `json:"uuid"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Creator     int         `json:"creator"`
	NamespaceId int         `json:"namespace_id"`
	CreateTime  string      `json:"create_time"`
	UpdateTime  string      `json:"update_time"`
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Labels      []LabelInfo `json:"labels"`
}

func (s *ExperimentInstanceService) GetExperimentInstanceByUUID(uuid string) (*ExperimentInstanceInfo, error) {
	exp, err := experiment_instance.GetExperimentInstanceByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if exp == nil {
		return nil, fmt.Errorf("no experiment instance found with uuid %s", uuid)
	}
	labels, err := experiment_instance.ListLabelsByExperimentInstanceUUID(uuid)
	if err != nil {
		return nil, err
	}

	expData := ExperimentInstanceInfo{
		Uuid:        exp.UUID,
		Name:        exp.Name,
		Description: exp.Description,
		Creator:     exp.Creator,
		NamespaceId: exp.NamespaceID,
		CreateTime:  exp.CreateTime.Format(time.RFC3339),
		UpdateTime:  exp.UpdateTime.Format(time.RFC3339),
		Status:      exp.Status,
		Message:     exp.Message,
	}

	for _, label := range labels {
		labelModel := namespace.Label{Id: label.LabelID, NamespaceId: expData.NamespaceId}
		if err := namespace.GetLabelByIdAndNamespaceId(context.Background(), &labelModel); err != nil {
			return nil, err
		}
		expData.Labels = append(expData.Labels, LabelInfo{
			Id:         label.LabelID,
			Name:       labelModel.Name,
			CreateTime: label.CreateTime.String(),
		})
	}
	return &expData, nil
}

type WorkflowNodesInfo struct {
	Uuid       string `json:"uuid"`
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Duration   string `json:"duration"`
	ExecType   string `json:"exec_type"`
	ExecId     int    `json:"exec_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (s *ExperimentInstanceService) GetWorkflowNodesInstanceByUUID(uuid string) (int, []WorkflowNodesInfo, error) {
	experiment, err := s.GetExperimentInstanceByUUID(uuid)
	if err != nil {
		return 0, nil, err
	}
	if experiment == nil {
		return 0, nil, errors.New("experiment not found")
	}
	nodes, err := experiment_instance.GetWorkflowNodeInstancesByExperimentUUID(uuid)
	if err != nil {
		return 0, nil, err
	}

	total := len(nodes)
	var workflowNodesInfos = []WorkflowNodesInfo{}

	for _, workflowNodeGet := range nodes {
		workflowNodesInfos = append(workflowNodesInfos, WorkflowNodesInfo{
			Uuid:       workflowNodeGet.UUID,
			Row:        workflowNodeGet.Row,
			Column:     workflowNodeGet.Column,
			Duration:   workflowNodeGet.Duration,
			ExecType:   workflowNodeGet.ExecType,
			ExecId:     workflowNodeGet.ExecID,
			Status:     workflowNodeGet.Status,
			Message:    workflowNodeGet.Message,
			CreateTime: workflowNodeGet.CreateTime.String(),
			UpdateTime: workflowNodeGet.UpdateTime.String()})
	}

	return total, workflowNodesInfos, nil
}

func (s *ExperimentInstanceService) GetWorkflowNodeInstanceByUUIDAndNodeId(uuid, nodeId string) (*WorkflowNodesInfo, error) {
	experiment, err := s.GetExperimentInstanceByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if experiment == nil {
		return nil, errors.New("experiment not found")
	}
	node, err := experiment_instance.GetWorkflowNodeInstanceByUUID(nodeId)
	if err != nil {
		return nil, err
	}
	return &WorkflowNodesInfo{
		Uuid:       node.UUID,
		Row:        node.Row,
		Column:     node.Column,
		Duration:   node.Duration,
		ExecType:   node.ExecType,
		ExecId:     node.ExecID,
		Status:     node.Status,
		Message:    node.Message,
		CreateTime: node.CreateTime.String(),
		UpdateTime: node.UpdateTime.String()}, nil
}

type ArgsValue struct {
	ArgsId int    `json:"args_id"`
	Value  string `json:"value"`
}

type WorkflowNodesDetail struct {
	WorkflowNodesInfo
	ArgsValue []ArgsValue                               `json:"args_value"`
	Subtasks  []*experiment_instance.FaultRangeInstance `json:"subtasks"`
}

func (s *ExperimentInstanceService) GetWorkflowNodeInstanceDetailByUUIDAndNodeId(uuid, nodeId string) (*WorkflowNodesDetail, error) {
	workflowNode, err := s.GetWorkflowNodeInstanceByUUIDAndNodeId(uuid, nodeId)
	if err != nil {
		return nil, err
	}
	if workflowNode == nil {
		return nil, errors.New("workflowNode not found")
	}
	workflowNodesDetail := WorkflowNodesDetail{
		WorkflowNodesInfo: *workflowNode,
	}
	argsValues, err := experiment_instance.GetArgsValueInstancesByWorkflowNodeUUID(nodeId)
	if err != nil {
		return &workflowNodesDetail, err
	}
	for _, argsValue := range argsValues {
		workflowNodesDetail.ArgsValue = append(workflowNodesDetail.ArgsValue, ArgsValue{ArgsId: argsValue.ArgsID, Value: argsValue.Value})
	}

	faultRanges, err := experiment_instance.ListFaultRangeInstancesByWorkflowNodeInstanceUUID(nodeId)
	if err != nil {
		return &workflowNodesDetail, err
	}
	workflowNodesDetail.Subtasks = faultRanges
	return &workflowNodesDetail, nil
}

func (s *ExperimentInstanceService) GetFaultRangeInstanceByWorkflowNodeInstanceUUID(uuid, nodeId, subtaskId string) (*experiment_instance.FaultRangeInstance, error) {
	if _, err := s.GetWorkflowNodeInstanceByUUIDAndNodeId(uuid, nodeId); err != nil {
		return nil, err
	}
	faultRangeInstance := experiment_instance.FaultRangeInstance{ID: cast.ToInt(subtaskId)}
	err := experiment_instance.GetFaultRangeInstanceById(&faultRangeInstance)
	return &faultRangeInstance, err
}

func (s *ExperimentInstanceService) DeleteExperimentInstanceByUUID(uuid string) error {
	if err := experiment_instance.ClearLabelIDsByExperimentInstanceUUID(uuid); err != nil {
		return err
	}

	workflowNodes, err := experiment_instance.GetWorkflowNodeInstancesByExperimentUUID(uuid)
	if err != nil {
		return err
	}

	for _, workflowNode := range workflowNodes {
		if err := experiment_instance.DeleteWorkflowNodeInstanceByUUID(workflowNode.UUID); err != nil {
			return err
		}
		if err := experiment_instance.ClearArgsValueInstancesByWorkflowNodeUUID(workflowNode.UUID); err != nil {
			return err
		}
		if err := experiment_instance.ClearFaultRangeInstancesByWorkflowNodeInstanceUUID(workflowNode.UUID); err != nil {
			return err
		}
	}
	return experiment_instance.DeleteExperimentInstanceByUUID(uuid)
}

func (s *ExperimentInstanceService) DeleteExperimentInstancesByUUID(uuids []string) error {
	for _, uuid := range uuids {
		if err := s.DeleteExperimentInstanceByUUID(uuid); err != nil {
			return err
		}
	}
	return nil
}

func (s *ExperimentInstanceService) SearchExperimentInstances(lastInstance string, namespaceId int, creator int, name string, scheduleType string, timeType string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []*experiment_instance.ExperimentInstance, error) {
	return experiment_instance.SearchExperimentInstances(lastInstance, namespaceId, creator, name, scheduleType, timeType, recentDays, startTime, endTime, orderBy, page, pageSize)
}
