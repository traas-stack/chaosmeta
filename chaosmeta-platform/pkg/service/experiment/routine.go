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
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models/experiment"
	experimentInstanceModel "chaosmeta-platform/pkg/models/experiment_instance"
	"chaosmeta-platform/pkg/service/cluster"
	"chaosmeta-platform/pkg/service/experiment_instance"
	"chaosmeta-platform/pkg/service/user"
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/robfig/cron"
	"time"
)

const (
	DefaultFormat = "2006-01-02 15:04:05"
)

type ExperimentRoutine struct {
	context   context.Context
	localCron *cron.Cron
}

func convertToExperimentInstance(experiment *ExperimentGet, status string) *experiment_instance.ExperimentInstance {
	experimentInstance := &experiment_instance.ExperimentInstance{
		ExperimentInstanceInfo: experiment_instance.ExperimentInstanceInfo{
			UUID:        experiment.UUID,
			Name:        experiment.Name,
			Description: experiment.Description,
			Creator:     experiment.Creator,
			NamespaceId: experiment.NamespaceID,
			Status:      status,
		},
		Labels: getLabelIdsFromLabelGet(experiment.Labels),
	}

	for _, node := range experiment.WorkflowNodes {
		workflowNodeDetail := &experiment_instance.WorkflowNodesDetail{
			WorkflowNodesInfo: experiment_instance.WorkflowNodesInfo{
				UUID:     node.UUID,
				Name:     node.Name,
				Row:      node.Row,
				Column:   node.Column,
				Duration: node.Duration,
				ScopeId:  node.ScopeId,
				TargetId: node.TargetId,
				ExecType: node.ExecType,
				ExecId:   node.ExecID,
			},
			Subtasks: &experimentInstanceModel.FaultRangeInstance{
				WorkflowNodeInstanceUUID: node.UUID,
			},
		}

		if node.FaultRange != nil {
			workflowNodeDetail.Subtasks.TargetName = node.FaultRange.TargetName
			workflowNodeDetail.Subtasks.TargetIP = node.FaultRange.TargetIP
			workflowNodeDetail.Subtasks.TargetHostname = node.FaultRange.TargetHostname
			workflowNodeDetail.Subtasks.TargetLabel = node.FaultRange.TargetLabel
			workflowNodeDetail.Subtasks.TargetApp = node.FaultRange.TargetApp
			workflowNodeDetail.Subtasks.TargetNamespace = node.FaultRange.TargetNamespace
			workflowNodeDetail.Subtasks.RangeType = node.FaultRange.RangeType
		}
		for _, arg := range node.ArgsValue {
			workflowNodeDetail.ArgsValues = append(workflowNodeDetail.ArgsValues, experiment_instance.ArgsValue{ArgsId: arg.ArgsID, Value: arg.Value})
		}
		experimentInstance.WorkflowNodes = append(experimentInstance.WorkflowNodes, workflowNodeDetail)
	}

	experimentInstanceBytes, _ := json.Marshal(experimentInstance)
	log.Error("convertToExperimentInstance:", string(experimentInstanceBytes))
	return experimentInstance
}

func StartExperiment(experimentID string, creatorName string) error {
	experimentService := ExperimentService{}
	experimentGet, err := experimentService.GetExperimentByUUID(experimentID)
	if err != nil || experimentGet == nil {
		return fmt.Errorf("error %v", err)
	}

	experimentInstance := convertToExperimentInstance(experimentGet, string(experimentInstanceModel.Running))
	if creatorName != "" {
		creatorId, err := user.GetIdByName(creatorName)
		if err != nil {
			log.Error(err)
			return err
		}
		experimentInstance.Creator = creatorId
	}
	experimentInstanceService := experiment_instance.ExperimentInstanceService{}
	experimentInstanceId, err := experimentInstanceService.CreateExperimentInstance(experimentInstance, "Pending")
	if err != nil {
		return err
	}

	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), config.DefaultRunOptIns.RunMode.Int())
	if err != nil {
		return err
	}

	argoWorkFlowCtl, err := NewArgoWorkFlowService(restConfig, config.DefaultRunOptIns.ArgoWorkflowNamespace)
	if err != nil {
		return err
	}

	nodes, err := experimentInstanceService.GetWorkflowNodeInstanceDetailList(experimentInstanceId)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = argoWorkFlowCtl.Create(*GetWorkflowStruct(experimentInstanceId, nodes))
	return err
}

func StopExperiment(experimentInstanceID string) error {
	experimentInstanceInfo, err := experimentInstanceModel.GetExperimentInstanceByUUID(experimentInstanceID)
	if err != nil || experimentInstanceInfo == nil {
		return fmt.Errorf("can not find experimentInstance")
	}
	if experimentInstanceInfo.Status == "Succeeded" {
		return errors.New("walkthrough is over")
	}

	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), config.DefaultRunOptIns.RunMode.Int())
	if err != nil {
		return err
	}

	argoWorkFlowCtl, err := NewArgoWorkFlowService(restConfig, config.DefaultRunOptIns.WorkflowNamespace)
	if err != nil {
		log.Error(err)
		return err
	}

	workFlowGet, status, err := argoWorkFlowCtl.Get(getWorFlowName(experimentInstanceID))
	if err != nil {
		log.Error(err)
		return nil
	}
	if status == "Succeeded" {
		return errors.New("experiment has ended")
	}

	chaosmetaService := NewChaosmetaService(restConfig)
	for _, node := range workFlowGet.Status.Nodes {
		if isInjectStepName(node.DisplayName) {
			chaosmetaCR, err := chaosmetaService.Get(context.Background(), config.DefaultRunOptIns.WorkflowNamespace, node.DisplayName)
			if err != nil {
				log.Error(err)
				return err
			}
			chaosmetaCR.Spec.TargetPhase = "recover"
			if _, err := chaosmetaService.Update(context.Background(), chaosmetaCR); err != nil {
				log.Error(err)
				return err
			}
			_, nodeId, err := getExperimentUUIDAndNodeIDFromStepName(node.DisplayName)
			if err != nil {
				log.Error(err)
				continue
			}

			if err := experimentInstanceModel.UpdateWorkflowNodeInstanceStatus(nodeId, "Succeeded", ""); err != nil {
				log.Error(err)
				continue
			}
		}

	}

	if err := argoWorkFlowCtl.Delete(getWorFlowName(experimentInstanceID)); err != nil {
		log.Error(err)
		return err
	}

	experimentInstanceInfo.Status = "Succeeded"
	for _, node := range workFlowGet.Status.Nodes {
		if node.TemplateName == string(ExperimentInject) || node.TemplateName == string(RawSuspend) {
			if node.Phase == "Failed" || node.Phase == "Error" {
				experimentInstanceInfo.Status = string(node.Phase)
				return experimentInstanceModel.UpdateExperimentInstance(experimentInstanceInfo)
			}
		}
	}

	return experimentInstanceModel.UpdateExperimentInstance(experimentInstanceInfo)
}

func (e *ExperimentRoutine) DealOnceExperiment() {
	_, experiments, err := experiment.ListExperimentsByScheduleTypeAndStatus(experiment.OnceMode, experiment.ToBeExecuted)
	if err != nil {
		log.Error(err)
		return
	}

	for _, experimentGet := range experiments {
		nextExec, _ := time.Parse(DefaultFormat, experimentGet.ScheduleRule)
		if time.Now().After(nextExec) {
			log.Error("create an experiment")
			if err := StartExperiment(experimentGet.UUID, ""); err != nil {
				log.Error(err)
				continue
			}
			experimentGet.Status = experiment.Executed
			experimentGet.LastInstance = time.Now().Format(TimeLayout)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}
		} else {
			continue
		}
	}

}

func (e *ExperimentRoutine) DealCronExperiment() {
	_, experiments, err := experiment.ListExperimentsByScheduleTypeAndStatus(experiment.CronMode, experiment.ToBeExecuted)
	if err != nil {
		log.Error(err)
		return
	}
	for _, experimentGet := range experiments {
		cronExpr, err := cron.Parse(experimentGet.ScheduleRule)
		if err != nil {
			continue
		}
		now := time.Now().Add(time.Minute)
		if experimentGet.NextExec.IsZero() {
			experimentGet.NextExec = cronExpr.Next(now)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
			}
			continue
		}

		if time.Now().After(experimentGet.NextExec) {
			experimentGet.Status = experiment.Executed
			experimentGet.NextExec = cronExpr.Next(now)
			log.Error(experimentGet.UUID, " next exec time", experimentGet.NextExec)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}

			if err := StartExperiment(experimentGet.UUID, ""); err != nil {
				log.Error(err)
			}

			experimentGet.Status = experiment.ToBeExecuted
			experimentGet.LastInstance = time.Now().Format(TimeLayout)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}
		}

	}

}

func (e *ExperimentRoutine) syncExperimentStatus(workflow v1alpha1.Workflow) error {
	log.Info("syncExperimentStatus.Name:", workflow.Name)
	experimentInstanceId, err := getExperimentInstanceIdFromWorkflowName(workflow.Name)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := experimentInstanceModel.UpdateExperimentInstanceStatus(experimentInstanceId, string(workflow.Status.Phase), workflow.Status.Message); err != nil {
		log.Error("UpdateExperimentInstanceStatus err:", err)
		return err
	}

	for _, node := range workflow.Status.Nodes {
		if node.TemplateName == string(ExperimentInject) || node.TemplateName == string(RawSuspend) {
			_, nodeId, err := getExperimentUUIDAndNodeIDFromStepName(node.DisplayName)
			if err != nil {
				log.Error("getExperimentUUIDAndNodeIDFromStepName", err)
				continue
			}

			if err := experimentInstanceModel.UpdateWorkflowNodeInstanceStatus(nodeId, string(node.Phase), node.Message); err != nil {
				log.Error("UpdateWorkflowNodeInstanceStatus", err)
				continue
			}
		}
	}
	return nil
}

func (e *ExperimentRoutine) SyncExperimentsStatus() {
	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), config.DefaultRunOptIns.RunMode.Int())
	if err != nil {
		log.Error(err)
		return
	}

	argoWorkFlowCtl, err := NewArgoWorkFlowService(restConfig, config.DefaultRunOptIns.ArgoWorkflowNamespace)
	pendingArgos, finishArgos, err := argoWorkFlowCtl.ListPendingAndFinishWorkflows()
	if err != nil {
		log.Error(err)
		return
	}

	errCh, doneCh := make(chan error), make(chan struct{})
	go func() {
		for _, pendingArgo := range pendingArgos {
			go func(argo v1alpha1.Workflow) {
				if err := e.syncExperimentStatus(argo); err != nil {
					errCh <- err
				}
			}(*pendingArgo)
		}
	}()

	go func() {
		for _, finishArgo := range finishArgos {
			go func(argo v1alpha1.Workflow) {
				if err := e.syncExperimentStatus(argo); err != nil {
					errCh <- err
				}
				if argo.Status.Phase == "Succeeded" {
					if err := argoWorkFlowCtl.Delete(argo.Name); err != nil {
						errCh <- err
					}
				}
			}(*finishArgo)
		}
	}()

	go func() {
		for range pendingArgos {
			<-doneCh
		}
		for range finishArgos {
			<-doneCh
		}
		close(errCh)
	}()

	for err := range errCh {
		log.Error(err)
	}

	close(doneCh)
}

func (e *ExperimentRoutine) DeleteExecutedInstanceCR() {
	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), config.DefaultRunOptIns.RunMode.Int())
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("expired Workflows have been deleted successfully.")

	ctx := context.Background()
	chaosmetaService := NewChaosmetaService(restConfig)
	if err := chaosmetaService.DeleteExpiredList(ctx); err != nil {
		log.Error(err)
	}
	log.Info("expired chaosmeta experiment  have been deleted successfully.")
	chaosmetaFlowInjectService := NewChaosmetaFlowInjectService(restConfig)
	if err := chaosmetaFlowInjectService.DeleteExpiredList(ctx); err != nil {
		log.Error(err)
	}
	log.Info("expired chaosmeta flow experiment  have been deleted successfully.")

}

func (e *ExperimentRoutine) Start() {
	localCron := cron.New()
	spec := "@every 3s"

	if err := localCron.AddFunc(spec, e.DealOnceExperiment); err != nil {
		log.Error(err)
		return
	}
	if err := localCron.AddFunc(spec, e.DealCronExperiment); err != nil {
		log.Error(err)
		return
	}

	if err := localCron.AddFunc(spec, e.SyncExperimentsStatus); err != nil {
		log.Error(err)
		return
	}

	if err := localCron.AddFunc("@every 6h", e.DeleteExecutedInstanceCR); err != nil {
		log.Error(err)
		return
	}

	localCron.Start()
	e.localCron = localCron

	select {
	case <-e.context.Done():
		log.Info("Receive stop signal")
	}
}
