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
	experimentInstanceModel "chaosmeta-platform/pkg/models/experiment_instance"
	"chaosmeta-platform/pkg/service/cluster"
	"chaosmeta-platform/pkg/service/experiment_instance"
	"chaosmeta-platform/util/log"
	"context"
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

func convertToExperimentInstance(experiment *Experiment, status string) *experiment_instance.ExperimentInstance {
	experimentInstance := &experiment_instance.ExperimentInstance{
		ExperimentInstanceInfo: experiment_instance.ExperimentInstanceInfo{
			UUID:        experiment.UUID,
			Name:        experiment.Name,
			Description: experiment.Description,
			Creator:     experiment.Creator,
			NamespaceId: experiment.NamespaceID,
			Status:      status,
		},
		Labels: experiment.Labels,
	}

	for _, node := range experiment.WorkflowNodes {
		workflowNodeDetail := &experiment_instance.WorkflowNodesDetail{
			WorkflowNodesInfo: experiment_instance.WorkflowNodesInfo{
				UUID:     node.UUID,
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
				TargetName:               node.FaultRange.TargetName,
				TargetIP:                 node.FaultRange.TargetIP,
				TargetHostname:           node.FaultRange.TargetHostname,
				TargetLabel:              node.FaultRange.TargetLabel,
				TargetApp:                node.FaultRange.TargetApp,
				TargetNamespace:          node.FaultRange.TargetNamespace,
				RangeType:                node.FaultRange.RangeType,
			},
		}
		for _, arg := range node.ArgsValue {
			workflowNodeDetail.ArgsValues = append(workflowNodeDetail.ArgsValues, experiment_instance.ArgsValue{ArgsId: arg.ArgsID, Value: arg.Value})
		}
		experimentInstance.WorkflowNodes = append(experimentInstance.WorkflowNodes, workflowNodeDetail)
	}

	return experimentInstance
}

func (e *ExperimentRoutine) DoExperiment(experimentID string) error {
	experimentService := ExperimentService{}
	experimentGet, err := experimentService.GetExperimentByUUID(experimentID)
	if err != nil {
		return err
	}

	experimentInstance := convertToExperimentInstance(experimentGet, string(experimentInstanceModel.Running))

	experimentInstanceService := experiment_instance.ExperimentInstanceService{}
	experimentInstanceId, err := experimentInstanceService.CreateExperimentInstance(experimentInstance)
	if err != nil {
		return err
	}

	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), -1)
	if err != nil {
		return err
	}

	argoWorkFlowCtl, err := NewArgoWorkFlowService(restConfig, WorkflowNamespace)
	if err != nil {
		return err
	}

	nodes, err := experimentInstanceService.GetWorkflowNodeInstanceDetailList(experimentInstanceId)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = argoWorkFlowCtl.Create(*GetWorkWorkflow(experimentInstanceId, nodes))
	return err
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
			if err := e.DoExperiment(experimentGet.UUID); err != nil {
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
			log.Error("parsing cron expression failed parsing cron table parsing cron expression failed:", err)
			continue
		}
		if experimentGet.NextExec.IsZero() {
			experimentGet.NextExec = cronExpr.Next(time.Now())
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
			}
			continue
		}

		nextExec := cronExpr.Next(experimentGet.NextExec)

		if time.Now().After(nextExec) {
			experimentGet.Status = experiment.Executed
			experimentGet.NextExec = cronExpr.Next(time.Now())
			if err := experiment.UpdateExperiment(experimentGet); err != nil {
				log.Error(err)
				continue
			}

			if err := e.DoExperiment(experimentGet.UUID); err != nil {
				log.Error(err)
				continue
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
	experimentInstanceId, err := getExperimentInstanceIdFromWorkflowName(workflow.Name)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := experimentInstanceModel.UpdateExperimentInstanceStatus(experimentInstanceId, string(workflow.Status.Phase)); err != nil {
		return err
	}

	for _, node := range workflow.Status.Nodes {
		_, nodeId, err := getExperimentUUIDAndNodeIDFromStepName(node.DisplayName)
		if err != nil {
			log.Error(err)
			continue
		}

		if err := experimentInstanceModel.UpdateWorkflowNodeInstanceStatus(nodeId, string(node.Phase)); err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func (e *ExperimentRoutine) SyncExperimentsStatus() {
	clusterService := cluster.ClusterService{}
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), -1)
	if err != nil {
		log.Error(err)
		return
	}

	argoWorkFlowCtl, err := NewArgoWorkFlowService(restConfig, WorkflowNamespace)
	pendingArgos, finishArgos, err := argoWorkFlowCtl.ListPendingAndRecentWorkflows()
	if err != nil {
		log.Error(err)
		return
	}

	errCh, doneCh := make(chan error), make(chan struct{})
	go func() {
		for _, pendingArgo := range pendingArgos.Items {
			go func(argo v1alpha1.Workflow) {
				if err := e.syncExperimentStatus(argo); err != nil {
					errCh <- err
				}
			}(pendingArgo)
		}
	}()

	go func() {
		for _, finishArgo := range finishArgos {
			go func(argo v1alpha1.Workflow) {
				if err := e.syncExperimentStatus(argo); err != nil {
					errCh <- err
				}
				if err := argoWorkFlowCtl.Delete(argo.Name); err != nil {
					errCh <- err
				}
			}(*finishArgo)
		}
	}()

	go func() {
		for range pendingArgos.Items {
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
	_, restConfig, err := clusterService.GetRestConfig(context.Background(), -1)
	if err != nil {
		log.Error(err)
		return
	}

	//argoWorkFlowCtl, err := workflow_channel.NewArgoWorkFlowService(restConfig, WorkflowNamespace)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//if err := argoWorkFlowCtl.DeleteExpiredList(); err != nil {
	//	log.Error(err)
	//	return
	//}
	log.Info("expired Workflows have been deleted successfully.")

	chaosmetaService := NewChaosmetaService(restConfig)
	if err := chaosmetaService.DeleteExpiredList(context.Background()); err != nil {
		log.Error(err)
		return
	}
	log.Info("expired chaosmeta experiment  have been deleted successfully.")

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
