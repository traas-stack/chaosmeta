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

package node

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
)

type NodeScopeHandler struct {
}

var globalNodeHandler = &NodeScopeHandler{}

func GetGlobalNodeHandler() *NodeScopeHandler {
	return globalNodeHandler
}

func (h *NodeScopeHandler) ConvertSelector(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error) {
	var (
		result  []model.AtomicObject
		isExist = make(map[string]bool)
	)

	argsList := common.GetArgs(spec.Experiment.Args, []string{v1alpha1.ContainerKey})

	for _, unitSelector := range spec.Selector {
		resultUnitSelector, err := getNodeObjectList(ctx, unitSelector, argsList[0])
		if err != nil {
			return nil, err
		}

		for _, unitObj := range resultUnitSelector {
			// Node Deduplication
			if isExist[unitObj.GetObjectName()] {
				continue
			}
			isExist[unitObj.GetObjectName()] = true
			result = append(result, unitObj)
		}
	}

	return result, nil
}

func (h *NodeScopeHandler) CheckAlive(ctx context.Context, injectObject model.AtomicObject) error {
	node, ok := injectObject.(*model.NodeObject)
	if !ok {
		return fmt.Errorf("inject object change to node error")
	}

	return remoteexecutor.GetRemoteExecutor().CheckAlive(ctx, node.NodeInternalIP)
}

func (h *NodeScopeHandler) QueryExperiment(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	node, ok := injectObject.(*model.NodeObject)
	if !ok {
		return nil, fmt.Errorf("inject object change to node error")
	}

	return remoteexecutor.GetRemoteExecutor().Query(ctx, node.NodeInternalIP, UID, phase)
}

func (h *NodeScopeHandler) ExecuteInject(ctx context.Context, injectObject model.AtomicObject, UID string, expArgs *v1alpha1.ExperimentCommon) (string, error) {
	node, ok := injectObject.(*model.NodeObject)
	if !ok {
		return "", fmt.Errorf("inject object change to node error")
	}

	if node.ContainerID != "" {
		return "", remoteexecutor.GetRemoteExecutor().Inject(ctx, node.NodeInternalIP, expArgs.Target, expArgs.Fault, UID, expArgs.Duration, node.ContainerID, node.ContainerRuntime, expArgs.Args)
	}

	return "", remoteexecutor.GetRemoteExecutor().Inject(ctx, node.NodeInternalIP, expArgs.Target, expArgs.Fault, UID, expArgs.Duration, "", "", expArgs.Args)
}

func (h *NodeScopeHandler) ExecuteRecover(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon) error {
	node, ok := injectObject.(*model.NodeObject)
	if !ok {
		return fmt.Errorf("inject object change to node error")
	}

	return remoteexecutor.GetRemoteExecutor().Recover(ctx, node.NodeInternalIP, UID)
}

func (h *NodeScopeHandler) GetInjectObject(ctx context.Context, exp *v1alpha1.ExperimentCommon, objectName string) (model.AtomicObject, error) {
	nodeName, nodeIP, err := model.ParseNodeInfo(objectName)
	if err != nil {
		return nil, fmt.Errorf("unexpected node format: %s", err.Error())
	}

	nodeInfo := &model.NodeObject{
		NodeName:       nodeName,
		NodeInternalIP: nodeIP,
	}

	containerName := common.GetArgs(exp.Args, []string{v1alpha1.ContainerKey})[0]
	if containerName != "" {
		r, id, err := model.ParseContainerID(containerName)
		if err != nil {
			return nil, fmt.Errorf("parse container info error: %s", err.Error())
		}

		nodeInfo.ContainerRuntime, nodeInfo.ContainerID = r, id
	}

	return nodeInfo, nil
}

// getInjectObjectList IP > nodeName > label
func getNodeObjectList(ctx context.Context, selectorUnit v1alpha1.SelectorUnit, containerName string) ([]model.AtomicObject, error) {
	var err error
	analyzer := selector.GetAnalyzer()
	var nodeList []*model.NodeObject
	if len(selectorUnit.IP) > 0 {
		nodeList, err = analyzer.GetNodeListByNodeIP(ctx, selectorUnit.IP, containerName)
	} else if len(selectorUnit.Name) > 0 {
		nodeList, err = analyzer.GetNodeListByNodeName(ctx, selectorUnit.Name, containerName)
	} else if len(selectorUnit.Label) > 0 {
		nodeList, err = analyzer.GetNodeListByLabel(ctx, selectorUnit.Label, containerName)
	} // other skip

	if err != nil {
		return nil, fmt.Errorf("get node list error: %s", err.Error())
	}

	var result = make([]model.AtomicObject, len(nodeList))
	for i := range nodeList {
		result[i] = nodeList[i]
	}

	return result, err
}
