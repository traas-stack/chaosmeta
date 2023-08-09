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
