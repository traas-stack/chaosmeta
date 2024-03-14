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

package selector

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	HostIPKey       = ".status.hostIP"
	PhaseKey        = ".status.phase"
	containerRegFmt = "^(%s)$"
)

var (
	globalAnalyzer IAnalyzer
)

func SetupAnalyzer(apiServer client.Client) {
	globalAnalyzer = &Analyzer{
		ApiServer: apiServer,
	}
}

func GetAnalyzer() IAnalyzer {
	return globalAnalyzer
}

type IAnalyzer interface {
	GetExperimentListByPhase(ctx context.Context, phase string) (*v1alpha1.ExperimentList, error)
	GetContainer(ctx context.Context, ns, podName, containerName string) (*model.ContainerObject, error)

	GetPod(ctx context.Context, ns, podName, containerName string) (*model.PodObject, error)
	GetPodListByLabelInNode(ctx context.Context, namespace string, label map[string]string, nodeIP string) ([]*model.PodObject, error)
	GetPodListByLabel(ctx context.Context, namespace string, label map[string]string, containerName string) ([]*model.PodObject, error)
	GetPodListByPodName(ctx context.Context, namespace string, podName []string, containerName string) ([]*model.PodObject, error)

	GetNodeListByLabel(ctx context.Context, label map[string]string, containerName string) ([]*model.NodeObject, error)
	GetNodeListByNodeName(ctx context.Context, nodeName []string, containerName string) ([]*model.NodeObject, error)
	GetNodeListByNodeIP(ctx context.Context, nodeIP []string, containerName string) ([]*model.NodeObject, error)

	GetDeploymentListByLabel(ctx context.Context, namespace string, label map[string]string) ([]*model.DeploymentObject, error)
	GetDeploymentListByName(ctx context.Context, namespace string, name []string) ([]*model.DeploymentObject, error)

	GetDaemonSetByName(ctx context.Context, namespace string, name string) (*model.DaemonSetObject, error)
}

type Analyzer struct {
	ApiServer client.Client
}

func (a *Analyzer) GetExperimentListByPhase(ctx context.Context, phase string) (*v1alpha1.ExperimentList, error) {
	opts := []client.ListOption{
		//client.MatchingFields{
		//	PhaseKey: phase,
		//},
		client.MatchingFields{
			PhaseKey: phase,
		},
	}

	expList := &v1alpha1.ExperimentList{}
	if err := a.ApiServer.List(ctx, expList, opts...); err != nil {
		return nil, fmt.Errorf("list experiment info by status error: %s", err.Error())
	}

	return expList, nil
}

func (a *Analyzer) GetPodListByLabelInNode(ctx context.Context, namespace string, label map[string]string, nodeIP string) ([]*model.PodObject, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(label),
		client.MatchingFields{
			HostIPKey: nodeIP,
		},
	}

	podList := &corev1.PodList{}
	if err := a.ApiServer.List(ctx, podList, opts...); err != nil {
		return nil, fmt.Errorf("list pod in node[%s] error: %s", nodeIP, err.Error())
	}

	var result = make([]*model.PodObject, len(podList.Items))
	for i, unitPod := range podList.Items {
		result[i] = &model.PodObject{
			PodName:   unitPod.Name,
			PodUID:    string(unitPod.UID),
			PodIP:     unitPod.Status.PodIP,
			Namespace: unitPod.Namespace,
			NodeName:  unitPod.Spec.NodeName,
			NodeIP:    unitPod.Status.HostIP,
		}
	}

	return result, nil
}

func (a *Analyzer) GetPodListByLabel(ctx context.Context, namespace string, label map[string]string, containerReg string) ([]*model.PodObject, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(label),
	}

	podList := &corev1.PodList{}
	if err := a.ApiServer.List(ctx, podList, opts...); err != nil {
		return nil, fmt.Errorf("list pod info by label error: %s", err.Error())
	}

	var result []*model.PodObject
	for _, unitPod := range podList.Items {
		podInfo := &model.PodObject{
			PodName:   unitPod.Name,
			PodUID:    string(unitPod.UID),
			PodIP:     unitPod.Status.PodIP,
			Namespace: unitPod.Namespace,
			NodeName:  unitPod.Spec.NodeName,
			NodeIP:    unitPod.Status.HostIP,
		}

		if containerReg != "" {
			var err error
			containers, err := GetTargetContainers(containerReg, unitPod.Status.ContainerStatuses)
			if err != nil {
				return nil, fmt.Errorf("get target container[%s] in pod[%s] error: %s", containerReg, unitPod.Name, err.Error())
			}
			podInfo.Containers = containers
		}

		result = append(result, podInfo)
	}

	return result, nil
}

func (a *Analyzer) GetPodListByPodName(ctx context.Context, namespace string, podName []string, containerName string) ([]*model.PodObject, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
	}

	podList := &corev1.PodList{}
	if err := a.ApiServer.List(ctx, podList, opts...); err != nil {
		return nil, fmt.Errorf("list pod info error: %s", err.Error())
	}

	podNameMap := make(map[string]bool)
	for _, unitP := range podName {
		podNameMap[unitP] = true
	}

	var result []*model.PodObject
	for _, unitPod := range podList.Items {
		if !podNameMap[unitPod.Name] {
			continue
		}

		podInfo := &model.PodObject{
			PodName:   unitPod.Name,
			PodUID:    string(unitPod.UID),
			PodIP:     unitPod.Status.PodIP,
			Namespace: unitPod.Namespace,
			NodeName:  unitPod.Spec.NodeName,
			NodeIP:    unitPod.Status.HostIP,
		}

		var err error
		containers, err := GetTargetContainers(containerName, unitPod.Status.ContainerStatuses)
		if err != nil {
			return nil, fmt.Errorf("get target container[%s] in pod[%s] error: %s", containerName, unitPod.Name, err.Error())
		}
		podInfo.Containers = containers

		result = append(result, podInfo)
	}

	return result, nil
}

func GetTargetContainer(containerName string, status []corev1.ContainerStatus) (container *model.ContainerInfo, err error) {
	if len(status) == 0 {
		err = fmt.Errorf("no container in pod")
		return
	}
	if containerName == v1alpha1.FirstContainer || containerName == "" {
		targetContainerInfo := status[0]
		info, err := getContainerInfo(targetContainerInfo)
		if err != nil {
			return nil, err
		}
		return info, nil
	} else {
		for _, containerStatus := range status {
			if containerName == containerStatus.Name {
				info, err := getContainerInfo(containerStatus)
				if err != nil {
					return nil, err
				}
				return info, nil
			}
		}
	}
	return nil, fmt.Errorf("fail to get target container")
}

func GetTargetContainers(containerReg string, status []corev1.ContainerStatus) (containers []model.ContainerInfo, err error) {
	if len(status) == 0 {
		err = fmt.Errorf("no container in pod")
		return
	}
	containerReg = strings.ReplaceAll(containerReg, ",", "|")
	containers = []model.ContainerInfo{}
	var targetContainerInfo corev1.ContainerStatus
	// no container setting equals to firstcontainer setting
	if containerReg == v1alpha1.FirstContainer || containerReg == "" {
		targetContainerInfo = status[0]
		info, err := getContainerInfo(targetContainerInfo)
		if err != nil {
			return nil, err
		}
		containers = append(containers, *info)
	} else {
		reg := regexp.MustCompile(fmt.Sprintf(containerRegFmt, containerReg))
		for _, containerStatus := range status {
			if reg.MatchString(containerStatus.Name) {
				info, err := getContainerInfo(containerStatus)
				if err != nil {
					return nil, err
				}
				containers = append(containers, *info)
			}
		}
	}
	return containers, nil
}

func getContainerInfo(containerStatus corev1.ContainerStatus) (*model.ContainerInfo, error) {
	r, id, err := model.ParseContainerID(containerStatus.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("parse container id[%s] error: %s", containerStatus.ContainerID, err.Error())
	}
	info := &model.ContainerInfo{
		ContainerId:      id,
		ContainerRuntime: r,
		ContainerName:    containerStatus.Name,
	}
	return info, nil
}

// GetNodeListByLabel return all node when label is empty map or nil
func (a *Analyzer) GetNodeListByLabel(ctx context.Context, label map[string]string, containerName string) ([]*model.NodeObject, error) {
	opts := []client.ListOption{
		client.MatchingLabels(label),
	}

	nodeList := &corev1.NodeList{}
	if err := a.ApiServer.List(ctx, nodeList, opts...); err != nil {
		return nil, fmt.Errorf("list node error: %s", err.Error())
	}

	var result = make([]*model.NodeObject, len(nodeList.Items))
	for i, unitNode := range nodeList.Items {
		result[i] = &model.NodeObject{
			NodeName: unitNode.Name,
		}

		for _, unitAddress := range unitNode.Status.Addresses {
			if unitAddress.Type == "InternalIP" {
				result[i].NodeInternalIP = unitAddress.Address
			} else if unitAddress.Type == "Hostname" {
				result[i].HostName = unitAddress.Address
			}
		}

		if containerName != "" {
			r, id, err := model.ParseContainerID(containerName)
			if err != nil {
				return nil, fmt.Errorf("parse container info error: %s", err.Error())
			}

			result[i].ContainerRuntime, result[i].ContainerID = r, id
		}
	}

	return result, nil
}

func (a *Analyzer) GetNodeListByNodeName(ctx context.Context, nodeName []string, containerName string) ([]*model.NodeObject, error) {
	nodeList := &corev1.NodeList{}

	if err := a.ApiServer.List(ctx, nodeList, []client.ListOption{}...); err != nil {
		return nil, fmt.Errorf("list node error: %s", err.Error())
	}

	nodeNameMap := make(map[string]bool)
	for _, unitNode := range nodeName {
		nodeNameMap[unitNode] = true
	}

	var result []*model.NodeObject
	for _, unitNode := range nodeList.Items {
		if !nodeNameMap[unitNode.Name] {
			continue
		}

		tmpNode := &model.NodeObject{
			NodeName: unitNode.Name,
		}

		for _, unitAddress := range unitNode.Status.Addresses {
			if unitAddress.Type == "InternalIP" {
				tmpNode.NodeInternalIP = unitAddress.Address
			} else if unitAddress.Type == "Hostname" {
				tmpNode.HostName = unitAddress.Address
			}
		}

		if containerName != "" {
			r, id, err := model.ParseContainerID(containerName)
			if err != nil {
				return nil, fmt.Errorf("parse container info error: %s", err.Error())
			}

			tmpNode.ContainerRuntime, tmpNode.ContainerID = r, id
		}

		result = append(result, tmpNode)
	}

	return result, nil
}

// GetNodeListByNodeIP In order to support non-k8s machines, so do not filter from apiServer
func (a *Analyzer) GetNodeListByNodeIP(ctx context.Context, nodeIP []string, containerName string) ([]*model.NodeObject, error) {
	nodeList := &corev1.NodeList{}

	if err := a.ApiServer.List(ctx, nodeList, []client.ListOption{}...); err != nil {
		return nil, fmt.Errorf("list node error: %s", err.Error())
	}

	nodeIPMap := make(map[string]bool)
	for _, unitIP := range nodeIP {
		nodeIPMap[unitIP] = true
	}

	var result []*model.NodeObject
	for _, unitNode := range nodeList.Items {
		var unitIP, unitHostName string
		for _, unitAddress := range unitNode.Status.Addresses {
			if unitAddress.Type == "InternalIP" {
				unitIP = unitAddress.Address
			} else if unitAddress.Type == "Hostname" {
				unitHostName = unitAddress.Address
			}
		}
		if unitIP == "" || !nodeIPMap[unitIP] {
			continue
		}

		tmpNode := &model.NodeObject{
			NodeName:       unitNode.Name,
			NodeInternalIP: unitIP,
			HostName:       unitHostName,
		}
		if containerName != "" {
			r, id, err := model.ParseContainerID(containerName)
			if err != nil {
				return nil, fmt.Errorf("parse container info error: %s", err.Error())
			}

			tmpNode.ContainerRuntime, tmpNode.ContainerID = r, id
		}

		result = append(result, tmpNode)
	}

	return result, nil
}

//	func (a *Analyzer) GetNodeListByNodeIP(ctx context.Context, nodeIP []string, containerName string) ([]*model.NodeObject, error) {
//		var nodeIPMap = make(map[string]bool)
//		var result []*model.NodeObject
//		for _, unit := range nodeIP {
//			if !nodeIPMap[unit] {
//				nodeIPMap[unit] = true
//				tmpNode := &model.NodeObject{
//					NodeInternalIP: unit,
//				}
//				if containerName != "" {
//					r, id, err := model.ParseContainerID(containerName)
//					if err != nil {
//						return nil, fmt.Errorf("parse container info error: %s", err.Error())
//					}
//
//					tmpNode.ContainerRuntime, tmpNode.ContainerID = r, id
//				}
//
//				result = append(result, tmpNode)
//			}
//		}
//
//		return result, nil
//	}
func (a *Analyzer) GetContainer(ctx context.Context, ns, podName, containerName string) (*model.ContainerObject, error) {
	pod := &corev1.Pod{}

	if err := a.ApiServer.Get(ctx, client.ObjectKey{
		Namespace: ns,
		Name:      podName,
	}, pod); err != nil {
		return nil, fmt.Errorf("get pod error: %s", err.Error())
	}
	container, err := GetTargetContainer(containerName, pod.Status.ContainerStatuses)
	if err != nil {
		return nil, fmt.Errorf("get target container[%s] in pod[%s] error: %s", containerName, pod.Name, err.Error())
	}
	containerObject := &model.ContainerObject{
		Namespace:        pod.Namespace,
		PodName:          pod.Name,
		PodUID:           string(pod.UID),
		PodIP:            pod.Status.PodIP,
		NodeName:         pod.Spec.NodeName,
		NodeIP:           pod.Status.HostIP,
		ContainerId:      container.ContainerId,
		ContainerRuntime: container.ContainerRuntime,
		ContainerName:    container.ContainerName,
	}
	return containerObject, nil
}

func (a *Analyzer) GetPod(ctx context.Context, ns, podName, containerName string) (*model.PodObject, error) {
	pod := &corev1.Pod{}

	if err := a.ApiServer.Get(ctx, client.ObjectKey{
		Namespace: ns,
		Name:      podName,
	}, pod); err != nil {
		return nil, fmt.Errorf("get pod error: %s", err.Error())
	}

	podInfo := &model.PodObject{
		Namespace: pod.Namespace,
		PodName:   pod.Name,
		PodUID:    string(pod.UID),
		PodIP:     pod.Status.PodIP,
		NodeName:  pod.Spec.NodeName,
		NodeIP:    pod.Status.HostIP,
	}

	var err error
	containers, err := GetTargetContainers(containerName, pod.Status.ContainerStatuses)
	if err != nil {
		return nil, fmt.Errorf("get target container[%s] in pod[%s] error: %s", containerName, pod.Name, err.Error())
	}
	podInfo.Containers = containers

	return podInfo, nil
}

func (a *Analyzer) GetDeploymentListByLabel(ctx context.Context, namespace string, label map[string]string) ([]*model.DeploymentObject, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(label),
	}

	deployList := &appsv1.DeploymentList{}
	if err := a.ApiServer.List(ctx, deployList, opts...); err != nil {
		return nil, fmt.Errorf("list deployment info error: %s", err.Error())
	}

	var result = make([]*model.DeploymentObject, len(deployList.Items))
	for i, unitDeploy := range deployList.Items {
		result[i] = &model.DeploymentObject{
			DeploymentName: unitDeploy.Name,
			Namespace:      unitDeploy.Namespace,
		}
	}

	return result, nil
}

func (a *Analyzer) GetDeploymentListByName(ctx context.Context, namespace string, name []string) ([]*model.DeploymentObject, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
	}

	deployList := &appsv1.DeploymentList{}
	if err := a.ApiServer.List(ctx, deployList, opts...); err != nil {
		return nil, fmt.Errorf("list deployment info error: %s", err.Error())
	}

	deployNameMap := make(map[string]bool)
	for _, unitP := range name {
		deployNameMap[unitP] = true
	}

	var result []*model.DeploymentObject
	for _, unitDeploy := range deployList.Items {
		if !deployNameMap[unitDeploy.Name] {
			continue
		}

		result = append(result, &model.DeploymentObject{
			DeploymentName: unitDeploy.Name,
			Namespace:      unitDeploy.Namespace,
		})
	}

	return result, nil
}

func (a *Analyzer) GetDaemonSetByName(ctx context.Context, namespace, name string) (*model.DaemonSetObject, error) {
	daemonSet := &appsv1.DaemonSet{}
	if err := a.ApiServer.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, daemonSet); err != nil {
		return nil, fmt.Errorf("failed to get daemonset by name: %s", name)
	}
	return &model.DaemonSetObject{
		DaemonSetName: daemonSet.Name,
		Namespace:     daemonSet.Namespace,
	}, nil
}
