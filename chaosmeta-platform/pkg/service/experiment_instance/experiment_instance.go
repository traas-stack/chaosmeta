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
	"chaosmeta-platform/pkg/models/user"
	"chaosmeta-platform/util/log"
	"chaosmeta-platform/util/snowflake"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

type ExperimentInstanceService struct{}

type ExperimentInstance struct {
	ExperimentInstanceInfo
	Labels        []int                                   `json:"labels,omitempty"`
	WorkflowNodes []*WorkflowNodesDetail                  `json:"workflow_nodes,omitempty"`
	FaultRange    *experiment_instance.FaultRangeInstance `json:"exec_range,omitempty"`
}

func (s *ExperimentInstanceService) createUUID(creator int, typeStr string) string {
	nodeSnow, err := snowflake.NewNode(1)
	if err != nil {
		log.Error(err)
		return ""
	}
	return fmt.Sprintf("%d%d%d%s", nodeSnow.Generate(), creator, time.Now().Unix(), typeStr)
}

func (s *ExperimentInstanceService) CreateExperimentInstance(experimentParam *ExperimentInstance, status string) (string, error) {
	experimentCreate := experiment_instance.ExperimentInstance{
		UUID:           s.createUUID(experimentParam.Creator, "experiment"),
		Name:           experimentParam.Name,
		NamespaceID:    experimentParam.NamespaceId,
		Description:    experimentParam.Description,
		ExperimentUUID: experimentParam.UUID,
		Creator:        experimentParam.Creator,
		Message:        experimentParam.Message,
		Status:         status,
	}

	// experiment
	if err := experiment_instance.CreateExperimentInstance(&experimentCreate); err != nil {
		return "", err
	}

	//label
	if err := experiment_instance.AddLabelIDsToExperiment(experimentCreate.UUID, experimentParam.Labels); err != nil {
		return experimentCreate.UUID, err
	}
	//workflow_nodes
	for _, node := range experimentParam.WorkflowNodes {
		workflowNodeCreate := experiment_instance.WorkflowNodeInstance{
			UUID:                   s.createUUID(experimentParam.Creator, "node"),
			Name:                   node.Name,
			ExperimentInstanceUUID: experimentCreate.UUID,
			Row:                    node.Row,
			Column:                 node.Column,
			Duration:               node.Duration,
			ScopeId:                node.ScopeId,
			TargetId:               node.TargetId,
			ExecType:               node.ExecType,
			ExecID:                 node.ExecId,
			Message:                node.Message,
		}
		if err := experiment_instance.CreateWorkflowNodeInstance(&workflowNodeCreate); err != nil {
			return experimentCreate.UUID, err
		}

		var argsValues []*experiment_instance.ArgsValueInstance
		for _, argsValue := range node.ArgsValues {
			argsValues = append(argsValues, &experiment_instance.ArgsValueInstance{
				ArgsID: argsValue.ArgsId,
				Value:  argsValue.Value,
			})
		}

		//args_value
		if len(argsValues) > 0 {
			if err := experiment_instance.BatchInsertArgsValueInstances(workflowNodeCreate.UUID, argsValues); err != nil {
				return experimentCreate.UUID, err
			}
		}

		//exec_range
		if node.Subtasks != nil {
			node.Subtasks.WorkflowNodeInstanceUUID = workflowNodeCreate.UUID
			if err := experiment_instance.CreateFaultRangeInstance(node.Subtasks); err != nil {
				return experimentCreate.UUID, err
			}
		}
	}
	return experimentCreate.UUID, nil
}

type LabelInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
}

type ExperimentInstanceInfo struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     int    `json:"creator"`
	CreatorName string `json:"creator_name,omitempty"`
	NamespaceId int    `json:"namespace_id"`

	CreateTime string      `json:"create_time"`
	UpdateTime string      `json:"update_time"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Labels     []LabelInfo `json:"labels"`
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

	userGet := user.User{ID: exp.Creator}
	if err := user.GetUserById(context.Background(), &userGet); err != nil {
		log.Error(err)
	}

	expData := ExperimentInstanceInfo{
		UUID:        exp.UUID,
		Name:        exp.Name,
		Description: exp.Description,
		Creator:     exp.Creator,
		CreatorName: userGet.Email,
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
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Duration   string `json:"duration"`
	ScopeId    int    `json:"scope_id"`
	TargetId   int    `json:"target_id"`
	ExecType   string `json:"exec_type"`
	ExecId     int    `json:"exec_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (s *ExperimentInstanceService) GetWorkflowNodesInstanceInfoByUUID(experimentUUID string) (int, []WorkflowNodesInfo, error) {
	experiment, err := s.GetExperimentInstanceByUUID(experimentUUID)
	if err != nil {
		return 0, nil, err
	}
	if experiment == nil {
		return 0, nil, errors.New("experiment not found")
	}
	nodes, err := experiment_instance.GetWorkflowNodeInstancesByExperimentUUID(experimentUUID)
	if err != nil {
		return 0, nil, err
	}

	total := len(nodes)
	var workflowNodesInfos = []WorkflowNodesInfo{}

	for _, workflowNodeGet := range nodes {
		workflowNodesInfos = append(workflowNodesInfos, WorkflowNodesInfo{
			UUID:       workflowNodeGet.UUID,
			Name:       workflowNodeGet.Name,
			Row:        workflowNodeGet.Row,
			Column:     workflowNodeGet.Column,
			Duration:   workflowNodeGet.Duration,
			ScopeId:    workflowNodeGet.ScopeId,
			TargetId:   workflowNodeGet.TargetId,
			ExecType:   workflowNodeGet.ExecType,
			ExecId:     workflowNodeGet.ExecID,
			Status:     workflowNodeGet.Status,
			Message:    workflowNodeGet.Message,
			CreateTime: workflowNodeGet.CreateTime.String(),
			UpdateTime: workflowNodeGet.UpdateTime.String()})
	}

	return total, workflowNodesInfos, nil
}

func (s *ExperimentInstanceService) GetWorkflowNodeInstanceByUUIDAndNodeId(experimentUUID, nodeId string) (*WorkflowNodesInfo, error) {
	experiment, err := s.GetExperimentInstanceByUUID(experimentUUID)
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
		UUID:       node.UUID,
		Name:       node.Name,
		Row:        node.Row,
		Column:     node.Column,
		Duration:   node.Duration,
		ScopeId:    node.ScopeId,
		TargetId:   node.TargetId,
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
	ArgsValues      []ArgsValue                               `json:"args_value"`
	Subtasks        *experiment_instance.FaultRangeInstance   `json:"subtasks"`
	FlowSubtasks    *experiment_instance.FlowRangeInstance    `json:"flow_subtasks"`
	MeasureSubtasks *experiment_instance.MeasureRangeInstance `json:"measure_subtasks"`
}

func (s *ExperimentInstanceService) GetWorkflowNodeInstanceDetailByUUIDAndNodeId(experimentUUID, nodeId string) (*WorkflowNodesDetail, error) {
	workflowNode, err := s.GetWorkflowNodeInstanceByUUIDAndNodeId(experimentUUID, nodeId)
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
		workflowNodesDetail.ArgsValues = append(workflowNodesDetail.ArgsValues, ArgsValue{ArgsId: argsValue.ArgsID, Value: argsValue.Value})
	}

	faultRange, err := experiment_instance.GetFaultRangeInstancesByWorkflowNodeInstanceUUID(nodeId)
	if err != nil {
		log.Error(err)
		return &workflowNodesDetail, err
	}
	workflowNodesDetail.Subtasks = faultRange

	flowRange, err := experiment_instance.GetFlowRangeInstancesByWorkflowNodeInstanceUUID(nodeId)
	if err != nil {
		log.Error(err)
		return &workflowNodesDetail, err
	}
	workflowNodesDetail.FlowSubtasks = flowRange

	measureRange, err := experiment_instance.GetMeasureRangeInstancesByWorkflowNodeInstanceUUID(nodeId)
	if err != nil {
		log.Error(err)
		return &workflowNodesDetail, err
	}
	workflowNodesDetail.MeasureSubtasks = measureRange
	return &workflowNodesDetail, nil
}

func (s *ExperimentInstanceService) GetWorkflowNodeInstanceDetailList(experimentUUID string) ([]*WorkflowNodesDetail, error) {
	experiment, err := s.GetExperimentInstanceByUUID(experimentUUID)
	if err != nil {
		return nil, err
	}
	if experiment == nil {
		return nil, errors.New("experiment not found")
	}
	nodes, err := experiment_instance.GetWorkflowNodeInstancesByExperimentUUID(experimentUUID)
	if err != nil {
		return nil, err
	}

	var workflowNodesDetails = []*WorkflowNodesDetail{}
	for _, workflowNodeGet := range nodes {
		workflowNodesDetail, err := s.GetWorkflowNodeInstanceDetailByUUIDAndNodeId(experimentUUID, workflowNodeGet.UUID)
		if err != nil {
			return nil, err
		}
		workflowNodesDetails = append(workflowNodesDetails, workflowNodesDetail)

	}
	return workflowNodesDetails, nil
}

func (s *ExperimentInstanceService) GetFaultRangeInstanceByWorkflowNodeInstanceUUID(uuid, nodeId, subtaskId string) (*experiment_instance.FaultRangeInstance, error) {
	if _, err := s.GetWorkflowNodeInstanceByUUIDAndNodeId(uuid, nodeId); err != nil {
		return nil, err
	}
	faultRangeInstance := experiment_instance.FaultRangeInstance{Id: cast.ToInt(subtaskId)}
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

func (s *ExperimentInstanceService) SearchExperimentInstances(lastInstance string, experimentUUID string, namespaceId int, creatorName string, name string, timeType string, timeSearchField string, status string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []*ExperimentInstanceInfo, error) {
	var experimentInstanceInfoList []*ExperimentInstanceInfo
	creator := 0
	userGet := user.User{Email: creatorName}
	if err := user.GetUser(context.Background(), &userGet); err != nil {
		log.Error(err)
	} else {
		creator = userGet.ID
	}
	total, experiments, err := experiment_instance.SearchExperimentInstances(lastInstance, experimentUUID, namespaceId, creator, name, timeType, timeSearchField, status, recentDays, startTime, endTime, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	for _, experiment := range experiments {
		labels, err := experiment_instance.ListLabelsByExperimentInstanceUUID(experiment.UUID)
		if err != nil {
			log.Error(err)
		}
		userGet := user.User{ID: experiment.Creator}
		if err := user.GetUserById(context.Background(), &userGet); err != nil {
			log.Error(err)
		}

		expData := ExperimentInstanceInfo{
			UUID:        experiment.UUID,
			Name:        experiment.Name,
			Description: experiment.Description,
			Creator:     experiment.Creator,
			CreatorName: userGet.Email,
			NamespaceId: experiment.NamespaceID,
			CreateTime:  experiment.CreateTime.Format(time.RFC3339),
			UpdateTime:  experiment.UpdateTime.Format(time.RFC3339),
			Status:      experiment.Status,
			Message:     experiment.Message,
		}

		for _, label := range labels {
			labelModel := namespace.Label{Id: label.LabelID, NamespaceId: expData.NamespaceId}
			if err := namespace.GetLabelByIdAndNamespaceId(context.Background(), &labelModel); err != nil {
				log.Error(err)
			}
			expData.Labels = append(expData.Labels, LabelInfo{
				Id:         label.LabelID,
				Name:       labelModel.Name,
				CreateTime: label.CreateTime.String(),
			})
		}
		experimentInstanceInfoList = append(experimentInstanceInfoList, &expData)
	}
	return total, experimentInstanceInfoList, nil
}
