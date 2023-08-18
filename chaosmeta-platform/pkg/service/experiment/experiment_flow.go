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
	"chaosmeta-platform/pkg/models/inject/basic"
	"chaosmeta-platform/pkg/service/experiment_instance"
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/util/log"
	"context"
	"errors"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

type WorkflowTemplateName string

var (
	RawSuspend       = WorkflowTemplateName("raw-suspend")
	ExperimentInject = WorkflowTemplateName("experiment-inject")
)

const (
	WorkflowNamespace = "chaosmeta"
	WorkflowMainStep  = "main"
)

var workflow = v1alpha1.Workflow{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "argoproj.io/v1alpha1",
		Kind:       "Workflow",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "",
		Namespace: "chaosmeta",
	},
	Spec: v1alpha1.WorkflowSpec{
		ServiceAccountName: kubernetes.ServiceAccount,
		Entrypoint:         WorkflowMainStep,
		Templates: []v1alpha1.Template{
			{
				Name: string(ExperimentInject),
				Inputs: v1alpha1.Inputs{
					Parameters: []v1alpha1.Parameter{
						{
							Name: "experiment",
						},
					},
				},
				Resource: &v1alpha1.ResourceTemplate{
					Action:           "create",
					FailureCondition: "status.status == failed",
					SuccessCondition: "status.phase == recover,status.status == success",
					Manifest:         "{{inputs.parameters.experiment}}",
				},
			},
			{
				Name: string(RawSuspend),
				Inputs: v1alpha1.Inputs{
					Parameters: []v1alpha1.Parameter{
						{
							Name: "time",
						},
					},
				},
				Suspend: &v1alpha1.SuspendTemplate{
					Duration: "{{inputs.parameters.time}}",
				},
			},
			{
				Name: "experiment-recover",
				Inputs: v1alpha1.Inputs{
					Parameters: []v1alpha1.Parameter{
						{
							Name: "experiment",
						},
					},
				},
				Resource: &v1alpha1.ResourceTemplate{
					Action:           "patch",
					FailureCondition: "status.status == failed",
					SuccessCondition: "status.phase == recover,status.status == success",
					MergeStrategy:    "json",
					Flags: []string{
						"experiments.inject.chaosmeta.io",
						"{{inputs.parameters.experiment}}",
					},
					Manifest: `- op: replace
								path: /spec/targetPhase
  								value: recover`,
				},
			},
		},
	},
}

type ExecType string

const (
	MeasureExecType ExecType = "measure"
	FaultExecType   ExecType = "fault"
	FlowExecType    ExecType = "flow"
	WaitExecType    ExecType = "wait"
)

func GetWorkWorkflow(experimentInstanceId string, nodes []*experiment_instance.WorkflowNodesDetail) *v1alpha1.Workflow {
	newWorkflow := workflow
	newWorkflow.Name = fmt.Sprintf("chaosmeta-inject-%s", experimentInstanceId)
	newWorkflow.Spec.Templates = append(newWorkflow.Spec.Templates, v1alpha1.Template{
		Name:  WorkflowMainStep,
		Steps: convertToSteps(experimentInstanceId, nodes),
	})
	return &newWorkflow
}

func getWaitStep(time string, experimentInstanceUUID string, nodeId string) *v1alpha1.WorkflowStep {
	waitStep := v1alpha1.WorkflowStep{
		Name:     fmt.Sprintf("before-wait-%s-%s", experimentInstanceUUID, nodeId),
		Template: string(RawSuspend),
		Arguments: v1alpha1.Arguments{
			Parameters: []v1alpha1.Parameter{
				{
					Name:  "time",
					Value: v1alpha1.AnyStringPtr(time),
				},
			},
		},
	}
	return &waitStep
}

func getInjectStep(experimentInstanceUUID string, node *experiment_instance.WorkflowNodesDetail, phaseType PhaseType) *v1alpha1.WorkflowStep {
	if node == nil {
		log.Error("node is nil")
		return nil
	}
	injectStep := v1alpha1.WorkflowStep{
		Template: string(ExperimentInject),
	}

	ctx := context.Background()
	scope, err := basic.GetScopeById(ctx, node.ScopeId)
	if err != nil {
		log.Error(err)
		return nil
	}
	target, err := basic.GetTargetById(ctx, node.TargetId)
	if err != nil {
		log.Error(err)
		return nil
	}

	fault, err := basic.GetFaultById(ctx, node.ExecId)
	if err != nil {
		log.Error(err)
		return nil
	}

	injectStep.Name = fmt.Sprintf("%s-%s-experiment-%s", scope.Name, target.Name, node.ExecType)

	experimentTemplate := ExperimentInjectStruct{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "inject.chaosmeta.io/v1alpha1",
			Kind:       "Experiment",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("inject-%s-%s-%s-%s", scope.Name, target.Name, experimentInstanceUUID, node.UUID),
			Namespace: WorkflowNamespace,
		},

		Spec: ExperimentSpec{
			Scope:       ScopeType(scope.Name),
			TargetPhase: phaseType,
			Experiment: &ExperimentCommon{
				Target:   target.Name,
				Fault:    fault.Name,
				Duration: node.Duration,
			},
		},
	}
	if node.Subtasks != nil {
		var selector SelectorUnit
		if node.Subtasks.TargetNamespace != "" {
			selector.Namespace = node.Subtasks.TargetNamespace

		}
		if node.Subtasks.TargetName != "" {
			selector.Name = strings.Split(node.Subtasks.TargetName, ",")
		}
		if node.Subtasks.TargetIP != "" {
			selector.IP = strings.Split(node.Subtasks.TargetIP, ",")
		}

		if node.Subtasks.TargetLabel != "" {
			labelMap := make(map[string]string)
			for _, pair := range strings.Split(node.Subtasks.TargetLabel, ",") {
				parts := strings.Split(pair, "=")
				labelMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
			selector.Label = labelMap
		}
		experimentTemplate.Spec.Selector = append(experimentTemplate.Spec.Selector, selector)
	}

	for _, arg := range node.ArgsValues {
		argGet, err := basic.GetArgsById(ctx, arg.ArgsId)
		if err != nil {
			log.Error(err)
			return nil
		}
		experimentTemplate.Spec.Experiment.Args = append(experimentTemplate.Spec.Experiment.Args, ArgsUnit{
			Key:       argGet.Key,
			Value:     arg.Value,
			ValueType: VType(argGet.ValueType),
		})
	}

	experimentTemplateBytes, err := yaml.Marshal(experimentTemplate)
	if err != nil {
		log.Error(err)
		return nil
	}

	injectStep.Arguments = v1alpha1.Arguments{
		Parameters: []v1alpha1.Parameter{
			{
				Name:  "experiment",
				Value: v1alpha1.AnyStringPtr(experimentTemplateBytes),
			},
		},
	}

	return &injectStep
}

func convertToSteps(experimentInstanceId string, nodes []*experiment_instance.WorkflowNodesDetail) []v1alpha1.ParallelSteps {
	maxRow, maxColumn := getMaxRowAndColumn(nodes)
	steps := make([]v1alpha1.ParallelSteps, maxRow+1)
	for i := range steps {
		steps[i].Steps = make([]v1alpha1.WorkflowStep, maxColumn+1)
	}
	step := &v1alpha1.WorkflowStep{}
	for _, node := range nodes {
		switch node.ExecType {
		case string(WaitExecType):
			step = getWaitStep(node.Duration, experimentInstanceId, node.UUID)
		case string(FaultExecType):
			step = getInjectStep(experimentInstanceId, node, InjectPhaseType)
		default:
			continue
		}
		if step != nil {
			steps[node.Row].Steps[node.Column] = *step
		}
	}
	return steps
}

func getMaxRowAndColumn(nodes []*experiment_instance.WorkflowNodesDetail) (int, int) {
	maxRow, maxColumn := 0, 0
	for _, node := range nodes {
		if node.Row > maxRow {
			maxRow = node.Row
		}
		if node.Column > maxColumn {
			maxColumn = node.Column
		}
	}
	return maxRow, maxColumn
}

func getNodeStatus(node v1alpha1.NodeStatus) string {
	switch node.Phase {
	case v1alpha1.NodePending:
		return "Pending"
	case v1alpha1.NodeRunning:
		return "Running"
	case v1alpha1.NodeSucceeded:
		return "Succeeded"
	case v1alpha1.NodeSkipped:
		return "Skipped"
	case v1alpha1.NodeFailed:
		return "Failed"
	case v1alpha1.NodeError:
		return "Error"
	case v1alpha1.NodeOmitted:
		return "Omitted"
	default:
		return ""
	}
}

//func syncWorkflowNodeStatus(workflow *v1alpha1.Workflow, nodes []*WorkflowNode) []*WorkflowNode {
//	updatedNodes := make([]*WorkflowNode, len(nodes))
//	nodeMap := make(map[string]*WorkflowNode)
//	for _, node := range nodes {
//		nodeMap[node.UUID] = node
//	}
//	for _, node := range workflow.Status.Nodes {
//		if node.Type == v1alpha1.NodeTypePod {
//			// 根据Pod节点的名称和结束时间，确定对应的WorkflowNode的执行状态
//			nodeName := getNodeName(node.Name)
//			if wNode, ok := nodeMap[nodeName]; ok {
//				wNode.Status = getNodeStatus(node)
//				updatedNodes = append(updatedNodes, wNode)
//			}
//		}
//	}
//	// 将更新后的WorkflowNode保存到数据库中
//	for _, node := range updatedNodes {
//		_, err := orm.NewOrm().Update(node)
//		if err != nil {
//			log.Errorf("Failed to update workflow node: %v", err)
//		}
//	}
//	return updatedNodes
//}
//
//func convertToWorkflowNode(workflowNode *v1alpha1.Nodes) *experiment.WorkflowNode {
//
//	return nil
//}

func getExperimentInstanceIdFromWorkflowName(workflowName string) (string, error) {
	reg := regexp.MustCompile(`chaosmeta-inject-(\w+)`)
	match := reg.FindStringSubmatch(workflowName)
	if len(match) < 2 {
		return "", fmt.Errorf("Failed to extract experimentInstanceId from workflowName")
	}
	return match[1], nil
}

func getExperimentUUIDAndNodeIDFromStepName(name string) (string, string, error) {
	var reg *regexp.Regexp
	var match []string

	if isInjectStepName(name) {
		reg = regexp.MustCompile(`inject-\w+-\w+-(\w+)-(\w+)`)
		match = reg.FindStringSubmatch(name)
	} else if isWaitStepName(name) {
		reg = regexp.MustCompile(`before-wait-(\w+)-(\w+)`)
		match = reg.FindStringSubmatch(name)
	} else {
		return "", "", errors.New("invalid name")
	}

	if len(match) < 3 {
		return "", "", errors.New("failed to extract experimentInstanceUUID and nodeId from name")
	}

	experimentInstanceUUID := match[1]
	nodeId := match[2]
	return experimentInstanceUUID, nodeId, nil
}

func isInjectStepName(name string) bool {
	reg := regexp.MustCompile(`inject-\w+-\w+-\w+-\w+`)
	return reg.MatchString(name)
}

func isWaitStepName(name string) bool {
	reg := regexp.MustCompile(`before-wait-\w+-\w+`)
	return reg.MatchString(name)
}
