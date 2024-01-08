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

package pod

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
)

type PodScopeHandler struct {
}

var globalPodHandler = &PodScopeHandler{}

func GetGlobalPodHandler() *PodScopeHandler {
	return globalPodHandler
}

func (h *PodScopeHandler) ConvertSelector(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error) {
	var (
		result  []model.AtomicObject
		isExist = make(map[string]bool)
	)

	//argsList := common.GetArgs(spec.Experiment.Args, []string{v1alpha1.ContainerKey})
	//if argsList[0] == "" {
	//	return nil, fmt.Errorf("container is not provide")
	//}

	for _, unitSelector := range spec.Selector {
		if unitSelector.Namespace == "" {
			return nil, fmt.Errorf("selector of scope pod must provide namespace")
		}

		resultUnitSelector, err := getPodObjectList(ctx, unitSelector)
		if err != nil {
			return nil, err
		}

		for _, unitObj := range resultUnitSelector {
			// Pod Deduplication
			if isExist[unitObj.GetObjectName()] {
				continue
			}
			isExist[unitObj.GetObjectName()] = true
			result = append(result, unitObj)
		}
	}

	return result, nil
}

func (h *PodScopeHandler) GetInjectObject(ctx context.Context, exp *v1alpha1.ExperimentCommon, objectName string) (model.AtomicObject, error) {
	analyzer := selector.GetAnalyzer()
	ns, podName, containerName, err := model.ParseContainerInfo(objectName)
	if err != nil {
		return nil, fmt.Errorf("unexpected pod object name: %s", objectName)
	}
	return analyzer.GetContainer(ctx, ns, podName, containerName)
}

func (h *PodScopeHandler) CheckAlive(ctx context.Context, injectObject model.AtomicObject) error {
	pod, ok := injectObject.(*model.ContainerObject)
	if !ok {
		return fmt.Errorf("inject object change to pod error")
	}

	return remoteexecutor.GetRemoteExecutor().CheckAlive(ctx, pod.NodeIP)
}

func (h *PodScopeHandler) QueryExperiment(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	container, ok := injectObject.(*model.ContainerObject)
	if !ok {
		return nil, fmt.Errorf("inject object change to container error")
	}

	return remoteexecutor.GetRemoteExecutor().Query(ctx, container.NodeIP, UID, phase)

}

func (h *PodScopeHandler) ExecuteInject(ctx context.Context, injectObject model.AtomicObject, UID string, expArgs *v1alpha1.ExperimentCommon) (string, error) {
	c, ok := injectObject.(*model.ContainerObject)
	if !ok {
		return "", fmt.Errorf("inject object change to pod error")
	}

	err := remoteexecutor.GetRemoteExecutor().Inject(ctx, c.NodeIP, expArgs.Target, expArgs.Fault, UID, expArgs.Duration, c.ContainerId, c.ContainerRuntime, expArgs.Args)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (h *PodScopeHandler) ExecuteRecover(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon) error {
	container, ok := injectObject.(*model.ContainerObject)
	if !ok {
		return fmt.Errorf("inject object change to pod error")
	}

	return remoteexecutor.GetRemoteExecutor().Recover(ctx, container.NodeIP, UID)
}

func getPodObjectList(ctx context.Context, selectorUnit v1alpha1.SelectorUnit) ([]model.AtomicObject, error) {
	var err error
	analyzer := selector.GetAnalyzer()
	var podList []*model.PodObject
	if len(selectorUnit.Name) != 0 {
		podList, err = analyzer.GetPodListByPodName(ctx, selectorUnit.Namespace, selectorUnit.Name, selectorUnit.SubName)
		if err != nil {
			return nil, fmt.Errorf("get pod info by podname list error: %s", err.Error())
		}
	} else {
		podList, err = analyzer.GetPodListByLabel(ctx, selectorUnit.Namespace, selectorUnit.Label, selectorUnit.SubName)
		if err != nil {
			return nil, fmt.Errorf("get pod info by podname list error: %s", err.Error())
		}
	}

	var result = make([]model.AtomicObject, 0)
	for _, pod := range podList {
		subObjects := pod.GetSubObjects()
		for i, _ := range subObjects {
			result = append(result, &subObjects[i])
		}
	}

	return result, err
}
