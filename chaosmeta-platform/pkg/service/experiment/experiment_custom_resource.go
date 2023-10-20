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
	"chaosmeta-platform/pkg/models/inject/basic"
	"chaosmeta-platform/pkg/service/experiment_instance"
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/util/log"
	"context"
	"errors"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/spf13/cast"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"sigs.k8s.io/yaml"
	"strings"
)

type WorkflowTemplateName string

var (
	RawSuspend           = WorkflowTemplateName("raw-suspend")
	ExperimentInjecFault = WorkflowTemplateName("experiment-fault")
	ExperimentInject     = WorkflowTemplateName("experiment-inject")
)

const (
	APIVersion       = "chaosmeta.io/v1alpha1"
	WorkflowMainStep = "main"
	ParametersName   = "Experiment"
)

var Manifest = `
- op: replace
  path: /spec/targetPhase
  value: recover
`

type ExecType string

const (
	MeasureExecType ExecType = "measure"
	FaultExecType   ExecType = "fault"
	FlowExecType    ExecType = "flow"
	WaitExecType    ExecType = "wait"
)

func getWorFlowName(experimentInstanceId string) string {
	return experimentInstanceId
}

func GetWorkflowStruct(experimentInstanceId string, nodes []*experiment_instance.WorkflowNodesDetail) *v1alpha1.Workflow {
	var newWorkflow = v1alpha1.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Workflow",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: config.DefaultRunOptIns.ArgoWorkflowNamespace,
		},
		Spec: v1alpha1.WorkflowSpec{
			ServiceAccountName: kubernetes.ServiceAccount,
			Entrypoint:         WorkflowMainStep,
			Templates: []v1alpha1.Template{
				{
					Name: string(ExperimentInjecFault),
					Inputs: v1alpha1.Inputs{
						Parameters: []v1alpha1.Parameter{
							{
								Name: ParametersName,
							},
						},
					},
					Resource: &v1alpha1.ResourceTemplate{
						Action:           "create",
						FailureCondition: "status.status == failed",
						SuccessCondition: "status.phase == recover,status.status == success",
						Manifest:         fmt.Sprintf("{{inputs.parameters.%s}}", ParametersName),
					},
				},
				{
					Name: string(ExperimentInject),
					Inputs: v1alpha1.Inputs{
						Parameters: []v1alpha1.Parameter{
							{
								Name: ParametersName,
							},
						},
					},
					Resource: &v1alpha1.ResourceTemplate{
						Action:           "create",
						FailureCondition: "status.status == failed",
						SuccessCondition: "status.status == success",
						Manifest:         fmt.Sprintf("{{inputs.parameters.%s}}", ParametersName),
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
			},
		},
	}
	newWorkflow.Name = getWorFlowName(experimentInstanceId)
	newWorkflow.Spec.Templates = append(newWorkflow.Spec.Templates, v1alpha1.Template{
		Name: WorkflowMainStep,
		DAG:  convertToSteps(experimentInstanceId, nodes),
	})

	return &newWorkflow
}

func getWaitStep(time string, experimentInstanceUUID string, nodeId string) *v1alpha1.DAGTask {
	waitStep := v1alpha1.DAGTask{
		Name:     getWaitStepName(experimentInstanceUUID, nodeId),
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

func getWaitStepName(experimentInstanceUUID string, nodeId string) string {
	return fmt.Sprintf("before-wait-%s-%s", experimentInstanceUUID, nodeId)
}

func getInjectStepName(scopeName, targetName, experimentInstanceUUID, nodeID string) string {
	return fmt.Sprintf("inject-fault-%s-%s-%s-%s", scopeName, targetName, "e", nodeID)
}

func getFlowStepName(flowtype, flowName, experimentInstanceUUID, nodeID string) string {
	return fmt.Sprintf("inject-flow-%s-%s-%s-%s", flowtype, flowName, "e", nodeID)
}

func getMeasureStepName(measureType, measureName, experimentInstanceUUID, nodeID string) string {
	return fmt.Sprintf("inject-measure-%s-%s-%s-%s", measureType, measureName, "e", nodeID)
}

func getFaultStep(experimentInstanceUUID string, node *experiment_instance.WorkflowNodesDetail, phaseType PhaseType) *v1alpha1.DAGTask {
	if node == nil {
		log.Error("node is nil")
		return nil
	}
	injectStep := v1alpha1.DAGTask{
		Template: string(ExperimentInjecFault),
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
	log.Error(scope.Name, target.Name, experimentInstanceUUID, node.UUID)
	injectStep.Name = getInjectStepName(scope.Name, target.Name, experimentInstanceUUID, node.UUID)
	experimentTemplate := ExperimentInjectStruct{
		TypeMeta: metav1.TypeMeta{
			APIVersion: APIVersion,
			Kind:       "Experiment",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      injectStep.Name,
			Namespace: config.DefaultRunOptIns.WorkflowNamespace,
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
				parts := strings.Split(pair, ":")
				if len(parts) >= 2 {
					labelMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
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
	log.Info("experimentTemplate", string(experimentTemplateBytes))
	injectStep.Arguments = v1alpha1.Arguments{
		Parameters: []v1alpha1.Parameter{
			{
				Name:  ParametersName,
				Value: v1alpha1.AnyStringPtr(string(experimentTemplateBytes)),
			},
		},
	}
	return &injectStep
}

func getFlowStep(experimentInstanceUUID string, node *experiment_instance.WorkflowNodesDetail) *v1alpha1.DAGTask {
	if node == nil {
		log.Error("node is nil")
		return nil
	}
	injectStep := v1alpha1.DAGTask{
		Template: string(ExperimentInject),
	}

	ctx := context.Background()
	flow, err := basic.GetFlowInjectByID(node.ExecId)
	if err != nil {
		log.Error(err)
		return nil
	}

	injectStep.Name = getFlowStepName(flow.FlowType, flow.Name, experimentInstanceUUID, node.UUID)

	flowSpec := LoadTest{
		TypeMeta: metav1.TypeMeta{
			APIVersion: APIVersion,
			Kind:       "LoadTest",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      injectStep.Name,
			Namespace: config.DefaultRunOptIns.WorkflowNamespace,
		},
	}

	if node.FlowSubtasks != nil {
		flowSpec.Spec.FlowType = node.FlowSubtasks.FlowType
		flowSpec.Spec.Duration = node.FlowSubtasks.Duration
		flowSpec.Spec.Parallelism = cast.ToInt(node.FlowSubtasks.Parallelism)
		flowSpec.Spec.Source = cast.ToInt(node.FlowSubtasks.Source)
	}

	for _, arg := range node.ArgsValues {
		argGet, err := basic.GetArgsById(ctx, arg.ArgsId)
		if err != nil {
			log.Error(err)
			return nil
		}
		flowSpec.Spec.Args = append(flowSpec.Spec.Args, FlowArgs{
			Key:   argGet.Key,
			Value: arg.Value,
		})
	}

	flowSpecBytes, err := yaml.Marshal(flowSpec)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("FlowSpecBytes:", string(flowSpecBytes))
	injectStep.Arguments = v1alpha1.Arguments{
		Parameters: []v1alpha1.Parameter{
			{
				Name:  ParametersName,
				Value: v1alpha1.AnyStringPtr(string(flowSpecBytes)),
			},
		},
	}
	return &injectStep
}

func getMeasureStep(experimentInstanceUUID string, node *experiment_instance.WorkflowNodesDetail) *v1alpha1.DAGTask {
	if node == nil {
		log.Error("node is nil")
		return nil
	}
	injectStep := v1alpha1.DAGTask{
		Template: string(ExperimentInject),
	}

	ctx := context.Background()
	measure, err := basic.GetMeasureInjectByID(node.ExecId)
	if err != nil {
		log.Error(err)
		return nil
	}

	injectStep.Name = getMeasureStepName(measure.MeasureType, measure.Name, experimentInstanceUUID, node.UUID)

	commonMeasureStruct := CommonMeasureStruct{
		TypeMeta: metav1.TypeMeta{
			APIVersion: APIVersion,
			Kind:       "CommonMeasure",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      injectStep.Name,
			Namespace: config.DefaultRunOptIns.WorkflowNamespace,
		},
	}

	if node.MeasureSubtasks != nil {
		commonMeasureStruct.Spec.MeasureType = node.MeasureSubtasks.MeasureType
		commonMeasureStruct.Spec.Duration = node.MeasureSubtasks.Duration
		commonMeasureStruct.Spec.Interval = node.MeasureSubtasks.Interval
		commonMeasureStruct.Spec.SuccessCount = cast.ToInt(node.MeasureSubtasks.SuccessCount)
		commonMeasureStruct.Spec.FailedCount = cast.ToInt(node.MeasureSubtasks.FailedCount)
		commonMeasureStruct.Spec.Judgement = Judgement{
			JudgeType:  node.MeasureSubtasks.JudgeType,
			JudgeValue: node.MeasureSubtasks.JudgeValue,
		}
	}

	for _, arg := range node.ArgsValues {
		argGet, err := basic.GetArgsById(ctx, arg.ArgsId)
		if err != nil {
			log.Error(err)
			return nil
		}
		commonMeasureStruct.Spec.Args = append(commonMeasureStruct.Spec.Args, MeasureArgs{
			Key:   argGet.Key,
			Value: arg.Value,
		})
	}

	commonMeasureStructBytes, err := yaml.Marshal(commonMeasureStruct)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("CommonMeasureStructBytes:", string(commonMeasureStructBytes))
	injectStep.Arguments = v1alpha1.Arguments{
		Parameters: []v1alpha1.Parameter{
			{
				Name:  ParametersName,
				Value: v1alpha1.AnyStringPtr(string(commonMeasureStructBytes)),
			},
		},
	}
	return &injectStep
}

//func getMaxRowAndColumn(nodes []*experiment_instance.WorkflowNodesDetail) (int, int) {
//	maxRow, maxColumn := 0, 0
//	for _, node := range nodes {
//		if node.Row > maxRow {
//			maxRow = node.Row
//		}
//		if node.Column > maxColumn {
//			maxColumn = node.Column
//		}
//	}
//	return maxRow, maxColumn
//}

func getStepArguments(experimentInstanceId string, node *experiment_instance.WorkflowNodesDetail) *v1alpha1.DAGTask {
	if node == nil {
		return &v1alpha1.DAGTask{}
	}
	switch node.ExecType {
	case string(WaitExecType):
		return getWaitStep(node.Duration, experimentInstanceId, node.UUID)
	case string(FaultExecType):
		return getFaultStep(experimentInstanceId, node, InjectPhaseType)
	case string(FlowExecType):
		return getFlowStep(experimentInstanceId, node)
	case string(MeasureExecType):
		return getMeasureStep(experimentInstanceId, node)
	default:
		return nil
	}
}

func convertToSteps(experimentInstanceId string, nodes []*experiment_instance.WorkflowNodesDetail) *v1alpha1.DAGTemplate {
	dAGTemplate := v1alpha1.DAGTemplate{}
	//_, maxColumn := getMaxRowAndColumn(nodes)
	var steps []v1alpha1.DAGTask
	beginTask := v1alpha1.DAGTask{
		Name:     "BeginWaitTask",
		Template: string(RawSuspend),
		Arguments: v1alpha1.Arguments{
			Parameters: []v1alpha1.Parameter{
				{
					Name:  "time",
					Value: v1alpha1.AnyStringPtr("0s"),
				},
			},
		},
	}
	endTask := beginTask
	endTask.Name = "EndWaitTask"

	steps = append(steps, beginTask)

	var prevNode *experiment_instance.WorkflowNodesDetail
	for _, node := range nodes {
		task := *getStepArguments(experimentInstanceId, node)
		if prevNode != nil && prevNode.Row != node.Row {
			endTask.Dependencies = append(endTask.Dependencies, getStepArguments(experimentInstanceId, prevNode).Name)
			//steps = append(steps, task)
			log.Debugf("End of row %d", prevNode.Row)
			task.Dependencies = []string{"BeginWaitTask"}
		}

		log.Debugf("%s(row:%d, column:%d) ", node.Name, node.Row, node.Column)
		if node.Column == 0 {
			task.Dependencies = []string{"BeginWaitTask"}
		}
		if prevNode != nil && prevNode.Row == node.Row {
			task.Dependencies = []string{getStepArguments(experimentInstanceId, prevNode).Name}
		}

		steps = append(steps, task)
		prevNode = node

	}
	endTask.Dependencies = append(endTask.Dependencies, getStepArguments(experimentInstanceId, prevNode).Name)
	steps = append(steps, endTask)
	dAGTemplate.Tasks = steps
	return &dAGTemplate
}

func getExperimentInstanceIdFromWorkflowName(workflowName string) (string, error) {
	parts := strings.Split(workflowName, "-")
	experimentID := ""
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasSuffix(parts[i], "experiment") {
			experimentID = parts[i]
			break
		}
	}
	return experimentID, nil
}

func getNodeIDFromStepName(name string) (string, error) {
	log.Info("ExperimentUUIDAndNodeIDFromStepName:", name)
	var reg *regexp.Regexp
	var match []string
	if isInjectStepName(name) {
		reg = regexp.MustCompile(`inject-\w+-\w+-\w+-(\w+)-(\w+)`)
		match = reg.FindStringSubmatch(name)
	} else if isWaitStepName(name) {
		reg = regexp.MustCompile(`before-wait-(\w+)-(\w+)`)
		match = reg.FindStringSubmatch(name)
	} else {
		return "", errors.New("invalid name")
	}

	if len(match) < 3 {
		return "", errors.New("failed to extract experimentInstanceUUID and nodeId from name")
	}
	return match[2], nil
}

func isInjectStepName(name string) bool {
	reg := regexp.MustCompile(`inject-\w+-\w+-\w+-\w+-\w+`)
	return reg.MatchString(name)
}

func getInjectSecondField(name string) (string, bool) {
	if !isInjectStepName(name) {
		return "", false
	}
	reg := regexp.MustCompile(`inject-(\w+)-\w+-\w+-\w+-\w+`)
	match := reg.FindStringSubmatch(name)
	if len(match) != 2 {
		return "", true
	}
	return match[1], true
}

func isWaitStepName(name string) bool {
	reg := regexp.MustCompile(`before-wait-\w+-\w+`)
	return reg.MatchString(name)
}
