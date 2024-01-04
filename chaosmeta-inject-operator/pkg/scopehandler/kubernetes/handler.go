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

package kubernetes

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/cloudnativeexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler/node"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler/pod"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
)

type KubernetesScopeHandler struct {
}

var globalKubernetesHandler = &KubernetesScopeHandler{}

func GetGlobalKubernetesHandler() *KubernetesScopeHandler {
	return globalKubernetesHandler
}

func (k KubernetesScopeHandler) ConvertSelector(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error) {
	switch v1alpha1.CloudTargetType(spec.Experiment.Target) {
	case v1alpha1.PodCloudTarget:
		return pod.GetGlobalPodHandler().ConvertSelector(ctx, spec)
	case v1alpha1.DeploymentCloudTarget:
		return convertDeploy(ctx, spec)
	case v1alpha1.NodeCloudTarget:
		return node.GetGlobalNodeHandler().ConvertSelector(ctx, spec)
	case v1alpha1.ClusterCloudTarget:
		return convertCluster(ctx, spec)
	default:
		return nil, fmt.Errorf("ConvertSelector not support target: %s", spec.Experiment.Target)
	}
}

func (k KubernetesScopeHandler) GetInjectObject(ctx context.Context, exp *v1alpha1.ExperimentCommon, objectName string) (model.AtomicObject, error) {
	switch v1alpha1.CloudTargetType(exp.Target) {
	case v1alpha1.PodCloudTarget:
		return pod.GetGlobalPodHandler().GetInjectObject(ctx, exp, objectName)
	case v1alpha1.DeploymentCloudTarget:
		ns, name, err := model.ParseDeploymentInfo(objectName)
		if err != nil {
			return nil, fmt.Errorf("unexpected deployment object name: %s", objectName)
		}

		return &model.DeploymentObject{
			Namespace:      ns,
			DeploymentName: name,
		}, nil
	case v1alpha1.NodeCloudTarget:
		return node.GetGlobalNodeHandler().GetInjectObject(ctx, exp, objectName)
	case v1alpha1.ClusterCloudTarget:
		return &model.NamespaceObject{
			Namespace: objectName,
		}, nil
	default:
		return nil, fmt.Errorf("GetInjectObject not support target: %s", exp.Target)
	}
}

func (k KubernetesScopeHandler) QueryExperiment(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return cloudnativeexecutor.GetCloudNativeExecutor(v1alpha1.CloudTargetType(expArgs.Target), expArgs.Fault).
		Query(ctx, injectObject.GetObjectName(), UID, backup, phase)
}

func (k KubernetesScopeHandler) ExecuteInject(ctx context.Context, injectObject model.AtomicObject, UID string, expArgs *v1alpha1.ExperimentCommon) (string, error) {
	p, ok := injectObject.(*model.PodObject)
	if !ok {
		return cloudnativeexecutor.GetCloudNativeExecutor(v1alpha1.CloudTargetType(expArgs.Target), expArgs.Fault).
			Inject(ctx, injectObject.GetObjectName(), UID, expArgs.Duration, expArgs.Args)
	}
	// pod object convert to container object
	subObjects := p.GetSubObjects()
	if len(subObjects) <= 0 {
		return "", fmt.Errorf("not found inject object")
	}
	var err error
	var resStr string
	for _, subObject := range subObjects {
		resStr, err = cloudnativeexecutor.GetCloudNativeExecutor(v1alpha1.CloudTargetType(expArgs.Target), expArgs.Fault).
			Inject(ctx, subObject.GetObjectName(), UID, expArgs.Duration, expArgs.Args)
		if err != nil {
			return resStr, err
		}
	}
	return resStr, err
}

func (k KubernetesScopeHandler) ExecuteRecover(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon) error {
	p, ok := injectObject.(*model.PodObject)
	if !ok {
		return cloudnativeexecutor.GetCloudNativeExecutor(v1alpha1.CloudTargetType(expArgs.Target), expArgs.Fault).
			Recover(ctx, injectObject.GetObjectName(), UID, backup)
	}
	subObjects := p.GetSubObjects()
	if len(subObjects) <= 0 {
		return fmt.Errorf("not found inject object")
	}
	for _, subObject := range subObjects {
		err := cloudnativeexecutor.GetCloudNativeExecutor(v1alpha1.CloudTargetType(expArgs.Target), expArgs.Fault).Recover(ctx, subObject.GetObjectName(), UID, backup)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k KubernetesScopeHandler) CheckAlive(ctx context.Context, injectObject model.AtomicObject) error {
	return nil
}

func convertCluster(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error) {
	args := common.GetArgs(spec.Experiment.Args, []string{"namespace"})
	if args[0] == "" {
		return nil, fmt.Errorf("namespace is empty")
	}

	var result = make([]model.AtomicObject, 1)
	result[0] = &model.NamespaceObject{
		Namespace: args[0],
	}
	return result, nil
}

func convertDeploy(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error) {
	var (
		result  []model.AtomicObject
		isExist = make(map[string]bool)
	)

	for _, unitSelector := range spec.Selector {
		if unitSelector.Namespace == "" {
			return nil, fmt.Errorf("selector of scope deployment must provide namespace")
		}

		resultUnitSelector, err := getDeployObjectFromSelector(ctx, unitSelector)
		if err != nil {
			return nil, err
		}

		for _, unitObj := range resultUnitSelector {
			// Deduplication
			if isExist[unitObj.GetObjectName()] {
				continue
			}
			isExist[unitObj.GetObjectName()] = true
			result = append(result, unitObj)
		}
	}

	return result, nil
}

func getDeployObjectFromSelector(ctx context.Context, selectorUnit v1alpha1.SelectorUnit) ([]model.AtomicObject, error) {
	var err error
	analyzer := selector.GetAnalyzer()
	var reList []*model.DeploymentObject
	if len(selectorUnit.Name) != 0 {
		reList, err = analyzer.GetDeploymentListByName(ctx, selectorUnit.Namespace, selectorUnit.Name)
		if err != nil {
			return nil, fmt.Errorf("get pod info by podname list error: %s", err.Error())
		}
	} else {
		reList, err = analyzer.GetDeploymentListByLabel(ctx, selectorUnit.Namespace, selectorUnit.Label)
		if err != nil {
			return nil, fmt.Errorf("get pod info by podname list error: %s", err.Error())
		}
	}

	var result = make([]model.AtomicObject, len(reList))
	for i := range reList {
		result[i] = reList[i]
	}

	return result, err
}
