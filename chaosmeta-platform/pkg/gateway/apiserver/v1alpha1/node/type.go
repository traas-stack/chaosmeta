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
	"chaosmeta-platform/pkg/service/kubernetes/kube"
	corev1 "k8s.io/api/core/v1"
)

type Node struct {
	Env        string `json:"env" description:"env"`
	Cluster    string `json:"cluster" description:"cluster"`
	NodeName   string `json:"nodeName" description:"nodeName"`
	UnSchedule bool   `json:"unSchedule" description:"unSchedule"`
}

type PatchNodeTaint struct {
	Env      string         `json:"env" description:"env"`
	Cluster  string         `json:"cluster" description:"cluster"`
	NodeName string         `json:"nodeName" description:"nodeName"`
	Taints   []corev1.Taint `json:"taints" description:"taints"`
}

type NodeInfo struct {
	kube.NodeDetail `json:",inline"`
}

type ListNodeResponse struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
	List     []struct {
		corev1.Node        `json:",inline"`
		AllocatedResources kube.NodeAllocatedResources `json:"allocatedResources,omitempty"`
	} `json:"list"`
}

type GetNodePodsResponse kube.PodResponse

type GetNodeEventsResponse kube.EventResponse
