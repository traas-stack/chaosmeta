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
	"k8s.io/client-go/rest"
	"time"
)

const (
	DefaultFormat = "2006-01-02 15:04:05"

	WorkflowPending   = "Pending" // pending some set-up - rarely used
	WorkflowRunning   = "Running" // any node has started; pods might not be running yet, the workflow maybe suspended too
	WorkflowSucceeded = "Succeeded"
	WorkflowFailed    = "Failed" // it maybe that the workflow was terminated
	WorkflowError     = "Error"
)

type ExperimentRoutine struct {
	context   context.Context
	localCron *cron.Cron
}

func convertToWorkflowNodesDetail(node *WorkflowNode, workflowNodesDetail *experiment_instance.WorkflowNodesDetail) {
	if node == nil || workflowNodesDetail == nil {
		return
	}
	if node.FaultRange != nil {
		workflowNodesDetail.Subtasks = &experimentInstanceModel.FaultRangeInstance{
			TargetName:      node.FaultRange.TargetName,
			TargetIP:        node.FaultRange.TargetIP,
			TargetHostname:  node.FaultRange.TargetHostname,
			TargetLabel:     node.FaultRange.TargetLabel,
			TargetApp:       node.FaultRange.TargetApp,
			TargetNamespace: node.FaultRange.TargetNamespace,
			RangeType:       node.FaultRange.RangeType,
		}
	}
	if node.FlowRange != nil {
		workflowNodesDetail.FlowSubtasks = &experimentInstanceModel.FlowRangeInstance{
			Source:      node.FlowRange.Source,
			Parallelism: node.FlowRange.Parallelism,
			Duration:    node.FlowRange.Duration,
			FlowType:    node.FlowRange.FlowType,
		}
	}
	if node.MeasureRange != nil {
		workflowNodesDetail.MeasureSubtasks = &experimentInstanceModel.MeasureRangeInstance{
			JudgeValue:   node.MeasureRange.JudgeValue,
			JudgeType:    node.MeasureRange.JudgeType,
			FailedCount:  node.MeasureRange.FailedCount,
			SuccessCount: node.MeasureRange.SuccessCount,
			Interval:     node.MeasureRange.Interval,
			Duration:     node.MeasureRange.Duration,
			MeasureType:  node.MeasureRange.MeasureType,
		}
	}
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
				ExecName: node.ExecName,
				ExecId:   node.ExecID,
			},
			Subtasks: &experimentInstanceModel.FaultRangeInstance{
				WorkflowNodeInstanceUUID: node.UUID,
			},
		}
		convertToWorkflowNodesDetail(node, workflowNodeDetail)
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
	experimentInstanceId, err := experimentInstanceService.CreateExperimentInstance(experimentInstance, WorkflowPending)
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

func injectRecoverByArgo(node v1alpha1.NodeStatus, experimentStatus *string, restConfig *rest.Config) error {
	injectType, isInject := getInjectSecondField(node.DisplayName)
	if isInject {
		if node.Phase == v1alpha1.NodeFailed || node.Phase == v1alpha1.NodeError {
			*experimentStatus = string(v1alpha1.WorkflowFailed)
			nodeId, err := getNodeIDFromStepName(node.DisplayName)
			if err != nil {
				log.Error(err)
				return err
			}
			if err := experimentInstanceModel.UpdateWorkflowNodeInstanceStatus(nodeId, string(node.Phase), ""); err != nil {
				log.Error(err)
			}
			return err
		}
		switch injectType {
		case string(FaultExecType):
			chaosmetaService := NewChaosmetaService(restConfig)
			if err := chaosmetaService.Recover(config.DefaultRunOptIns.WorkflowNamespace, node.DisplayName); err != nil {
				log.Error("fault CR recover failed, err:", err)
				return err
			}
		case string(FlowExecType):
			chaosmetaService := NewChaosmetaFlowService(restConfig)
			if err := chaosmetaService.Recover(config.DefaultRunOptIns.WorkflowNamespace, node.DisplayName); err != nil {
				log.Error("flow CR recover failed, err:", err)
				return err
			}
		case string(MeasureExecType):
			chaosmetaService := NewChaosmetaMeasureService(restConfig)
			if err := chaosmetaService.Recover(config.DefaultRunOptIns.WorkflowNamespace, node.DisplayName); err != nil {
				log.Error("measure CR recover failed, err:", err)
				return err
			}
		}
		nodeId, err := getNodeIDFromStepName(node.DisplayName)
		if err != nil {
			log.Error(err)
			return err
		}

		if err := experimentInstanceModel.UpdateWorkflowNodeInstanceStatus(nodeId, WorkflowSucceeded, ""); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func stopExperiment(experimentInstanceID string, experimentStatus *string, tolerateFailure bool) error {
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

	if status == WorkflowSucceeded {
		return errors.New("experiment has ended")
	}

	for _, node := range workFlowGet.Status.Nodes {
		if err := injectRecoverByArgo(node, experimentStatus, restConfig); err != nil {
			if !tolerateFailure {
				log.Error(err)
				return err
			}
		}
	}

	workFlowGet.Spec.Shutdown = v1alpha1.ShutdownStrategyStop
	if _, err := argoWorkFlowCtl.Update(*workFlowGet); err != nil {
		log.Error(err)
		return err
	}
	return argoWorkFlowCtl.Delete(getWorFlowName(experimentInstanceID))
}

func StopExperiment(experimentInstanceID string, tolerateFailure bool) error {
	experimentInstanceInfo, err := experimentInstanceModel.GetExperimentInstanceByUUID(experimentInstanceID)
	if err != nil || experimentInstanceInfo == nil {
		return fmt.Errorf("can not find experimentInstance")
	}
	var experimentStatus = WorkflowSucceeded
	if err := stopExperiment(experimentInstanceID, &experimentStatus, tolerateFailure); err != nil {
		log.Error("stopExperiment error:", err)
	}
	experimentInstanceInfo.Status = experimentStatus
	return experimentInstanceModel.UpdateExperimentInstance(experimentInstanceInfo)
}

func UserStopExperiment(experimentInstanceID string) error {
	experimentInstanceInfo, err := experimentInstanceModel.GetExperimentInstanceByUUID(experimentInstanceID)
	if err != nil || experimentInstanceInfo == nil {
		return fmt.Errorf("can not find experimentInstance")
	}
	if experimentInstanceInfo.Status == WorkflowSucceeded || experimentInstanceInfo.Status == WorkflowFailed || experimentInstanceInfo.Status == WorkflowError {
		return errors.New("experiment is over")
	}
	return StopExperiment(experimentInstanceID, false)
}

func (e *ExperimentRoutine) DealOnceExperiment() {
	_, experiments, err := experiment.ListExperimentsByScheduleTypeAndStatus(experiment.OnceMode, experiment.ToBeExecuted)
	if err != nil {
		log.Error(err)
		return
	}

	for _, experimentGet := range experiments {
		nextExec, _ := time.Parse(DefaultFormat, experimentGet.ScheduleRule)
		timeNow, _ := time.Parse(DefaultFormat, time.Now().Format(DefaultFormat))
		if timeNow.After(nextExec) {
			experimentGet.LastInstance = timeNow.Format(TimeLayout)
			log.Info(experimentGet.UUID, "next exec time", experimentGet.NextExec)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}

			log.Error("create an experiment")
			if err := StartExperiment(experimentGet.UUID, ""); err != nil {
				log.Error(err)
				continue
			}
			experimentGet.Status = experiment.Executed
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
			experimentGet.LastInstance = time.Now().Format(TimeLayout)
			log.Info(experimentGet.UUID, "next exec time", experimentGet.NextExec)
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}

			if err := StartExperiment(experimentGet.UUID, ""); err != nil {
				log.Error(err)
			}

			experimentGet.Status = experiment.ToBeExecuted
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func (e *ExperimentRoutine) syncExperimentStatus(workflow v1alpha1.Workflow) error {
	log.Info("syncExperimentStatus.Name:", workflow.Name, "workflow.Status", workflow.Status)
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
			nodeId, err := getNodeIDFromStepName(node.DisplayName)
			if err != nil {
				log.Error("getExperimentUUIDAndNodeIDFromStepName", err)
				continue
			}
			if node.Phase == v1alpha1.NodeFailed || node.Phase == v1alpha1.NodeError {
				return StopExperiment(experimentInstanceId, true)
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
	log.Info("expired workflows have been deleted successfully.")

	ctx := context.Background()
	chaosmetaService := NewChaosmetaService(restConfig)
	if err := chaosmetaService.DeleteExpiredList(ctx, config.DefaultRunOptIns.WorkflowNamespace); err != nil {
		log.Error(err)
	}
	log.Info("expired chaosmeta fault experiment have been deleted successfully.")
	chaosmetaFlowInjectService := NewChaosmetaFlowService(restConfig)
	if err := chaosmetaFlowInjectService.DeleteExpiredList(ctx, config.DefaultRunOptIns.WorkflowNamespace); err != nil {
		log.Error(err)
	}
	log.Info("expired chaosmeta flow experiment have been deleted successfully.")
	chaosmetaMeasureService := NewChaosmetaMeasureService(restConfig)
	if err := chaosmetaMeasureService.DeleteExpiredList(ctx, config.DefaultRunOptIns.WorkflowNamespace); err != nil {
		log.Error(err)
	}
	log.Info("expired chaosmeta measure experiment have been deleted successfully.")
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
