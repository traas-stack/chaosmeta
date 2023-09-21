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
	"chaosmeta-platform/pkg/models/experiment_instance"
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
	LastInstance string    `json:"last_instance,omitempty"`
}

type LabelGet struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	NamespaceId int    `json:"namespaceId"`
}

type ExperimentCreate struct {
	ExperimentInfo
	Labels        []int           `json:"labels,omitempty"`
	WorkflowNodes []*WorkflowNode `json:"workflow_nodes,omitempty"`
}

type ExperimentGet struct {
	UUID          string          `json:"uuid,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	ScheduleType  string          `json:"schedule_type"`
	ScheduleRule  string          `json:"schedule_rule"`
	NamespaceID   int             `json:"namespace_id"`
	Creator       int             `json:"creator,omitempty"`
	NextExec      time.Time       `json:"next_exec,omitempty"`
	CreatorName   string          `json:"creator_name,omitempty"`
	Status        int             `json:"status"`
	LastInstance  string          `json:"last_instance"`
	CreateTime    time.Time       `json:"create_time,omitempty"`
	UpdateTime    time.Time       `json:"update_time,omitempty"`
	Labels        []LabelGet      `json:"labels,omitempty"`
	WorkflowNodes []*WorkflowNode `json:"workflow_nodes,omitempty"`
	Number        int64           `json:"number,omitempty"`
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

func getLabelIdsFromLabelGet(labels []LabelGet) []int {
	labelIds := make([]int, len(labels))
	for i, label := range labels {
		labelIds[i] = label.Id
	}
	return labelIds
}

func (es *ExperimentService) CreateExperiment(experimentParam *ExperimentCreate) (string, error) {
	if experimentParam == nil {
		return "", errors.New("experimentParam is nil")
	}
	experimentUUid := es.createUUID(experimentParam.Creator, "")

	//label
	if len(experimentParam.Labels) > 0 {
		if err := experiment.AddLabelIDsToExperiment(experimentUUid, experimentParam.Labels); err != nil {
			return "", err
		}
	}

	//workflow_nodes
	for _, node := range experimentParam.WorkflowNodes {
		node.ExperimentUUID = experimentUUid
		workflowNodeCreate := experiment.WorkflowNode{
			UUID:           node.UUID,
			Name:           node.Name,
			ExperimentUUID: experimentUUid,
			Row:            node.Row,
			Column:         node.Column,
			Duration:       node.Duration,
			ScopeId:        node.ScopeId,
			TargetId:       node.TargetId,
			ExecType:       node.ExecType,
			ExecID:         node.ExecID,
		}
		if err := experiment.CreateWorkflowNode(&workflowNodeCreate); err != nil {
			return "", err
		}

		//args_value
		if len(node.ArgsValue) > 0 {
			if err := experiment.BatchInsertArgsValues(node.UUID, node.ArgsValue); err != nil {
				return "", err
			}
		}

		//exec_range
		if node.FaultRange != nil {
			node.FaultRange.WorkflowNodeInstanceUUID = node.UUID
			if err := experiment.CreateFaultRange(node.FaultRange); err != nil {
				return "", err
			}
		}
	}

	// experiment
	experimentCreate := experiment.Experiment{
		UUID:         experimentUUid,
		Name:         experimentParam.Name,
		NamespaceID:  experimentParam.NamespaceID,
		Description:  experimentParam.Description,
		ScheduleType: experimentParam.ScheduleType,
		ScheduleRule: experimentParam.ScheduleRule,
		Creator:      experimentParam.Creator,
	}
	if err := experiment.CreateExperiment(&experimentCreate); err != nil {
		return "", err
	}
	return experimentCreate.UUID, nil
}

func (es *ExperimentService) UpdateExperiment(uuid string, experimentParam *ExperimentCreate) error {
	if experimentParam == nil {
		return errors.New("experimentParam is nil")
	}
	getExperiment, err := experiment.GetExperimentByUUID(uuid)
	if err != nil {
		return fmt.Errorf("no this experiment")
	}

	experimentUUid := getExperiment.UUID
	log.Error(1)
	//label
	if len(experimentParam.Labels) > 0 {
		if err := experiment.ClearLabelIDsByExperimentUUID(uuid); err != nil {
			log.Error(err)
			return err
		}
		if err := experiment.AddLabelIDsToExperiment(uuid, experimentParam.Labels); err != nil {
			log.Error(err)
			return err
		}
	}
	log.Error(2)
	//workflow_nodes

	for _, node := range experimentParam.WorkflowNodes {
		node.ExperimentUUID = experimentUUid
		workflowNodeCreate := experiment.WorkflowNode{
			UUID:           node.UUID,
			Name:           node.Name,
			ExperimentUUID: experimentUUid,
			Row:            node.Row,
			Column:         node.Column,
			Duration:       node.Duration,
			ScopeId:        node.ScopeId,
			TargetId:       node.TargetId,
			ExecType:       node.ExecType,
			ExecID:         node.ExecID,
		}

		log.Error(3)
		if err := experiment.DeleteWorkflowNodeByUUID(node.UUID); err != nil {
			log.Error(err)
			return err
		}
		log.Error(4)
		if err := experiment.CreateWorkflowNode(&workflowNodeCreate); err != nil {
			log.Error(err)
			return err
		}
		log.Error(5)
		//args_value
		if len(node.ArgsValue) > 0 {
			if err := experiment.ClearArgsValuesByWorkflowNodeUUID(node.UUID); err != nil {
				log.Error(err)
				return err
			}
			if err := experiment.BatchInsertArgsValues(node.UUID, node.ArgsValue); err != nil {
				log.Error(err)
				return err
			}
		}

		log.Error(6)
		//exec_range
		if node.FaultRange != nil {
			node.FaultRange.WorkflowNodeInstanceUUID = node.UUID
			if err := experiment.ClearFaultRangesByWorkflowNodeInstanceUUID(node.UUID); err != nil {
				log.Error(err)
				return err
			}
			if err := experiment.CreateFaultRange(node.FaultRange); err != nil {
				return err
			}
		}
	}

	getExperiment.Name = experimentParam.Name
	getExperiment.Description = experimentParam.Description
	getExperiment.ScheduleType = experimentParam.ScheduleType
	getExperiment.ScheduleRule = experimentParam.ScheduleRule
	log.Error(7)
	return experiment.UpdateExperiment(getExperiment)
	//experimentParam.Creator = getExperiment.Creator
	//if err := es.DeleteExperimentByUUID(uuid); err != nil {
	//	return err
	//}
	//_, err = es.CreateExperiment(experimentParam)
	//return err
}

func (es *ExperimentService) UpdateExperimentStatusAndLastInstance(uuid string, status int, lastInstance string) error {
	experimentGet, err := experiment.GetExperimentByUUID(uuid)
	if err != nil || experimentGet == nil {
		return fmt.Errorf("no experiment")
	}
	if status >= 0 {
		experimentGet.Status = experiment.ExperimentStatus(status)
	}
	if lastInstance != "" {
		experimentGet.LastInstance = lastInstance
	}
	return experiment.UpdateExperiment(experimentGet)
}

func (es *ExperimentService) DeleteExperimentByUUID(uuid string) error {
	if err := experiment.ClearLabelIDsByExperimentUUID(uuid); err != nil {
		return err
	}

	workflowNodes, err := experiment.GetWorkflowNodesByExperimentUUID(uuid)
	if err != nil {
		log.Error(err)
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

func (es *ExperimentService) GetExperimentByUUID(uuid string) (*ExperimentGet, error) {
	experimentGet, err := experiment.GetExperimentByUUID(uuid)
	if err != nil || experimentGet == nil {
		return nil, fmt.Errorf("no experiment")
	}

	userGet := user.User{ID: experimentGet.Creator}
	if err := user.GetUserById(context.Background(), &userGet); err != nil {
		log.Errorf("can not find user, [user-id: %s]", err)
	}

	experimentReturn := ExperimentGet{
		UUID:         experimentGet.UUID,
		Name:         experimentGet.Name,
		Description:  experimentGet.Description,
		ScheduleType: experimentGet.ScheduleType,
		ScheduleRule: experimentGet.ScheduleRule,
		NamespaceID:  experimentGet.NamespaceID,
		CreatorName:  userGet.Email,
		Creator:      experimentGet.Creator,
		NextExec:     experimentGet.NextExec,
		Status:       int(experimentGet.Status),
		LastInstance: experimentGet.LastInstance,
		CreateTime:   experimentGet.CreateTime,
		UpdateTime:   experimentGet.UpdateTime,
	}

	experimentCount, _ := experiment_instance.CountExperimentInstances(0, experimentGet.UUID, "", 0)
	experimentReturn.Number = experimentCount
	if err := es.GetLabelByExperiment(uuid, &experimentReturn); err != nil {
		return &experimentReturn, nil
	}

	return &experimentReturn, es.GetWorkflowNodesByExperiment(uuid, &experimentReturn)
	//CountExperimentInstances()
}

func (es *ExperimentService) GetLabelByExperiment(uuid string, experimentReturn *ExperimentGet) error {
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
		experimentReturn.Labels = append(experimentReturn.Labels, LabelGet{
			Id:          labelId,
			Name:        getLabel.Name,
			Color:       getLabel.Color,
			NamespaceId: getLabel.NamespaceId,
		})
	}
	return nil
}

func (es *ExperimentService) GetWorkflowNodesByExperiment(uuid string, experimentReturn *ExperimentGet) error {
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
			log.Error(err)
		}
		nodeResult.FaultRange = faultRange
		workflowNodes = append(workflowNodes, &nodeResult)

	}
	experimentReturn.WorkflowNodes = append(experimentReturn.WorkflowNodes, workflowNodes...)
	return nil
}

func (es *ExperimentService) SearchExperiments(lastInstance string, namespaceId int, creatorName string, name string, scheduleType string, timeType string, timeSearchField string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []ExperimentGet, error) {
	log.Error(lastInstance, namespaceId, creatorName, name, scheduleType, timeType, timeSearchField, recentDays, startTime, endTime, orderBy, page, pageSize)
	var experimentList []ExperimentGet
	creator := 0
	if creatorName != "" {
		userGet := user.User{Email: creatorName}
		if err := user.GetUser(context.Background(), &userGet); err != nil {
			log.Error(err)
		} else {
			creator = userGet.ID
		}
	}

	total, experiments, err := experiment.SearchExperiments(lastInstance, namespaceId, creator, name, scheduleType, timeType, timeSearchField, recentDays, startTime, endTime, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	for _, experiment := range experiments {
		experimentGet, err := es.GetExperimentByUUID(experiment.UUID)
		if err != nil {
			return 0, nil, err
		}
		experimentList = append(experimentList, *experimentGet)
	}
	return total, experimentList, nil
}
