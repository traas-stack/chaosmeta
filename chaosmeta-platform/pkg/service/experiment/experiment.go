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
	"chaosmeta-platform/pkg/models/experiment"
	"chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/models/user"
	"chaosmeta-platform/util/log"
	"chaosmeta-platform/util/snowflake"
	"context"
	"errors"
	"fmt"
	"time"
)

func Init() {
	er := ExperimentRoutine{context: context.Background()}
	go er.Start()
}

const TimeLayout = "2006-01-02 15:04:05"

type ExperimentService struct{}

type ExperimentInfo struct {
	UUID         string    `json:"uuid,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ScheduleType string    `json:"schedule_type"`
	ScheduleRule string    `json:"schedule_rule"`
	NamespaceID  int       `json:"namespace_id"`
	Creator      int       `json:"creator,omitempty"`
	CreatorName  string    `json:"creator_name,omitempty"`
	Status       int       `json:"status"`
	CreateTime   time.Time `json:"create_time,omitempty"`
	UpdateTime   time.Time `json:"update_time,omitempty"`
}

type Label struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	NamespaceId int    `json:"namespaceId"`
}

type Experiment struct {
	ExperimentInfo
	Labels        []Label         `json:"labels,omitempty"`
	WorkflowNodes []*WorkflowNode `json:"workflow_nodes,omitempty"`
}

type ExperimentGet struct {
	UUID          string          `json:"uuid,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	ScheduleType  string          `json:"schedule_type"`
	ScheduleRule  string          `json:"schedule_rule"`
	NamespaceID   int             `json:"namespace_id"`
	Creator       string          `json:"creator,omitempty"`
	Status        int             `json:"status"`
	Labels        []Label         `json:"labels,omitempty"`
	WorkflowNodes []*WorkflowNode `json:"workflow_nodes,omitempty"`
}

type WorkflowNode struct {
	experiment.WorkflowNode
	ArgsValue  []*experiment.ArgsValue `json:"args_value,omitempty"`
	FaultRange *experiment.FaultRange  `json:"exec_range,omitempty"`
}

func (es *ExperimentService) createUUID(creator int, typeStr string) string {
	nodeSnow, err := snowflake.NewNode(1)
	if err != nil {
		log.Error(err)
		return ""
	}
	if typeStr != "" {
		return fmt.Sprintf("%d%d%s", nodeSnow.Generate(), creator, typeStr)
	}
	return fmt.Sprintf("%d%d", nodeSnow.Generate(), creator)
}

func getLabelIdsFromExperiment(labels []Label) []int {
	labelIds := make([]int, len(labels))
	for i, label := range labels {
		labelIds[i] = label.Id
	}
	return labelIds
}

func (es *ExperimentService) CreateExperiment(experimentParam *Experiment) (string, error) {
	if experimentParam == nil {
		return "", errors.New("experimentParam is nil")
	}
	experimentCreate := experiment.Experiment{
		UUID:         es.createUUID(experimentParam.Creator, ""),
		Name:         experimentParam.Name,
		NamespaceID:  experimentParam.NamespaceID,
		Description:  experimentParam.Description,
		ScheduleType: experimentParam.ScheduleType,
		ScheduleRule: experimentParam.ScheduleRule,
		Creator:      experimentParam.Creator,
	}

	// experiment
	if err := experiment.CreateExperiment(&experimentCreate); err != nil {
		return "", err
	}

	//label
	if len(experimentParam.Labels) > 0 {
		if err := experiment.AddLabelIDsToExperiment(experimentCreate.UUID, getLabelIdsFromExperiment(experimentParam.Labels)); err != nil {
			return experimentCreate.UUID, err
		}
	}

	//workflow_nodes
	for _, node := range experimentParam.WorkflowNodes {
		node.ExperimentUUID = experimentCreate.UUID
		workflowNodeCreate := experiment.WorkflowNode{
			UUID:           node.UUID,
			Name:           node.Name,
			ExperimentUUID: experimentCreate.UUID,
			Row:            node.Row,
			Column:         node.Column,
			Duration:       node.Duration,
			ScopeId:        node.ScopeId,
			TargetId:       node.TargetId,
			ExecType:       node.ExecType,
			ExecID:         node.ExecID,
		}
		if err := experiment.CreateWorkflowNode(&workflowNodeCreate); err != nil {
			return experimentCreate.UUID, err
		}

		//args_value
		if len(node.ArgsValue) > 0 {
			if err := experiment.BatchInsertArgsValues(node.UUID, node.ArgsValue); err != nil {
				return experimentCreate.UUID, err
			}
		}

		//exec_range
		if node.FaultRange != nil {
			node.FaultRange.WorkflowNodeInstanceUUID = node.UUID
			if err := experiment.CreateFaultRange(node.FaultRange); err != nil {
				return experimentCreate.UUID, err
			}
		}
	}
	return experimentCreate.UUID, nil
}

func (es *ExperimentService) UpdateExperiment(uuid string, experimentParam *Experiment) error {
	if experimentParam == nil {
		return errors.New("experimentParam is nil")
	}
	_, err := experiment.GetExperimentByUUID(uuid)
	if err != nil {
		return err
	}
	if err := es.DeleteExperimentByUUID(uuid); err != nil {
		return err
	}
	_, err = es.CreateExperiment(experimentParam)
	return err
}

func (es *ExperimentService) DeleteExperimentByUUID(uuid string) error {
	if err := experiment.ClearLabelIDsByExperimentUUID(uuid); err != nil {
		return err
	}

	workflowNodes, err := experiment.GetWorkflowNodesByExperimentUUID(uuid)
	if err != nil {
		return err
	}

	for _, workflowNode := range workflowNodes {
		if err := experiment.DeleteWorkflowNodeByUUID(workflowNode.UUID); err != nil {
			return err
		}
		// 删除args_value
		if err := experiment.ClearArgsValuesByWorkflowNodeUUID(workflowNode.UUID); err != nil {
			return err
		}
		// 删除fault_range
		if err := experiment.ClearFaultRangesByWorkflowNodeInstanceUUID(workflowNode.UUID); err != nil {
			return err
		}
	}
	return experiment.DeleteExperimentByUUID(uuid)
}

func (es *ExperimentService) GetExperimentByUUID(uuid string) (*Experiment, error) {
	experimentGet, err := experiment.GetExperimentByUUID(uuid)
	if err != nil || experimentGet == nil {
		return nil, fmt.Errorf("no experiment")
	}

	userGet := user.User{ID: experimentGet.Creator}
	if err := user.GetUserById(context.Background(), &userGet); err != nil {
		log.Error(err)
	}

	experimentReturn := Experiment{
		ExperimentInfo: ExperimentInfo{
			UUID:         experimentGet.UUID,
			Name:         experimentGet.Name,
			Description:  experimentGet.Description,
			ScheduleType: experimentGet.ScheduleType,
			ScheduleRule: experimentGet.ScheduleRule,
			NamespaceID:  experimentGet.NamespaceID,
			Creator:      experimentGet.Creator,
			CreatorName:  userGet.Email,
			Status:       int(experimentGet.Status),
			CreateTime:   experimentGet.CreateTime,
			UpdateTime:   experimentGet.UpdateTime,
		},
	}

	if err := es.GetLabelByExperiment(uuid, &experimentReturn); err != nil {
		return &experimentReturn, nil
	}

	return &experimentReturn, es.GetWorkflowNodesByExperiment(uuid, &experimentReturn)
}

func (es *ExperimentService) GetExperimentByModelExperiment(experimentGet *experiment.Experiment) (*Experiment, error) {
	if experimentGet == nil {
		return nil, errors.New("experimentGet is nil")
	}
	if experimentGet.UUID == "" {
		return nil, errors.New("experiment uuid is empty")
	}
	userGet := user.User{ID: experimentGet.Creator}
	if err := user.GetUserById(context.Background(), &userGet); err != nil {
		log.Error(err)
	}

	experimentReturn := Experiment{
		ExperimentInfo: ExperimentInfo{
			UUID:         experimentGet.UUID,
			Name:         experimentGet.Name,
			Description:  experimentGet.Description,
			ScheduleType: experimentGet.ScheduleType,
			ScheduleRule: experimentGet.ScheduleRule,
			NamespaceID:  experimentGet.NamespaceID,
			Creator:      experimentGet.Creator,
			CreatorName:  userGet.Email,
			Status:       int(experimentGet.Status),
			CreateTime:   experimentGet.CreateTime,
			UpdateTime:   experimentGet.UpdateTime,
		},
	}

	if err := es.GetLabelByExperiment(experimentGet.UUID, &experimentReturn); err != nil {
		return &experimentReturn, nil
	}

	return &experimentReturn, es.GetWorkflowNodesByExperiment(experimentGet.UUID, &experimentReturn)
}

func (es *ExperimentService) GetLabelByExperiment(uuid string, experimentReturn *Experiment) error {
	labelList, err := experiment.ListLabelIDsByExperimentUUID(uuid)
	if err != nil {
		return err
	}

	for _, labelId := range labelList {
		getLabel := namespace.Label{Id: labelId}
		if err := namespace.GetLabelById(context.Background(), &getLabel); err != nil {
			log.Error(err)
			continue
		}
		experimentReturn.Labels = append(experimentReturn.Labels, Label{
			Id:          labelId,
			Name:        getLabel.Name,
			Color:       getLabel.Color,
			NamespaceId: getLabel.NamespaceId,
		})
	}
	return nil
}

func (es *ExperimentService) GetWorkflowNodesByExperiment(uuid string, experimentReturn *Experiment) error {
	if experimentReturn == nil {
		return errors.New("experimentReturn is nil")
	}
	workflowNodesGet, err := experiment.GetWorkflowNodesByExperimentUUID(uuid)
	if err != nil {
		return err
	}

	var workflowNodes []*WorkflowNode

	for _, workflowNodeGet := range workflowNodesGet {
		nodeResult := WorkflowNode{
			WorkflowNode: experiment.WorkflowNode{
				UUID:     workflowNodeGet.UUID,
				Name:     workflowNodeGet.Name,
				Row:      workflowNodeGet.Row,
				Column:   workflowNodeGet.Column,
				Duration: workflowNodeGet.Duration,
				ScopeId:  workflowNodeGet.ScopeId,
				TargetId: workflowNodeGet.TargetId,
				ExecType: workflowNodeGet.ExecType,
				ExecID:   workflowNodeGet.ExecID,
			},
		}

		argsValue, err := experiment.GetArgsValuesByWorkflowNodeUUID(workflowNodeGet.UUID)
		if err != nil {
			return err
		}
		nodeResult.ArgsValue = append(nodeResult.ArgsValue, argsValue...)

		faultRange, err := experiment.GetFaultRangeByWorkflowNodeInstanceUUID(workflowNodeGet.UUID)
		if err != nil {
			return err
		}
		nodeResult.FaultRange = faultRange

		workflowNodes = append(workflowNodes, &nodeResult)

	}
	experimentReturn.WorkflowNodes = append(experimentReturn.WorkflowNodes, workflowNodes...)
	return nil
}

func (es *ExperimentService) SearchExperiments(lastInstance string, namespaceId int, creator int, name string, scheduleType string, timeType string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []Experiment, error) {
	var experimentList []Experiment
	total, experiments, err := experiment.SearchExperiments(lastInstance, namespaceId, creator, name, scheduleType, timeType, recentDays, startTime, endTime, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	for _, experiment := range experiments {
		experimentGet, err := es.GetExperimentByModelExperiment(experiment)
		if err != nil {
			return 0, nil, err
		}
		experimentList = append(experimentList, *experimentGet)
	}
	return total, experimentList, nil
}
