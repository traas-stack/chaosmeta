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

package kube

import (
	"chaosmeta-platform/pkg/models/common/page"
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/util"
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"sync"
)

type NodeService interface {
	Get(name string) (*NodeDetail, error)
	List(dsQuery *page.DataSelectQuery) (*NodeResponse, error)
	GetEvents(name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
	GetPods(name string, dsQuery *page.DataSelectQuery) (*PodResponse, error)
	GetRawPods(name string) (*corev1.PodList, error)
	Patch(name string, param ReplaceNodeInfoParam) (*NodeDetail, error)
	CordonOrUnCordon(nodeName string, drain bool) (*NodeDetail, error)
	TaintOrUnTaint(nodeName, op string, taint *corev1.Taint) (*NodeDetail, error)
	PatchTaint(nodeName string, taints []corev1.Taint) (*NodeDetail, error)
}

// PatchStringValue  specifies a patch operation for a string.
type PatchStringValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type NodeDetail struct {
	corev1.Node `json:",inline"`
	// AllocatedResources node allocated resources
	AllocatedResources *NodeAllocatedResources `json:"allocatedResources,omitempty"`
}

type NodeResponse struct {
	Total    int          `json:"total"`
	Current  int          `json:"current"`
	PageSize int          `json:"pageSize"`
	List     []NodeDetail `json:"list"`
}

type ReplaceNodeInfoParam struct {
	Env      string `json:"env" description:"env"`
	Cluster  string `json:"cluster" description:"cluster"`
	NodeName string `json:"nodeName" description:"nodeName"`

	OperatorPath string            `json:"operator_path"` // Type：labels or annotations
	OperatorData map[string]string `json:"operator_data"` // content
}

type NodeCell NodeDetail

func (n NodeCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(n.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(n.ObjectMeta.CreationTimestamp.Time)
	case page.NamespaceProperty:
		return page.StdComparableString(n.ObjectMeta.Namespace)
	case page.LabelProperty:
		var labelList []string
		for k, _ := range n.Labels {
			labelList = append(labelList, k)
		}
		return page.StdComparableString(strings.Join(labelList, " "))

	default:
		return nil
	}
}

type NodeAllocatedResources struct {
	// CPURequests is number of allocated milicores.
	CPURequests int64 `json:"cpuRequests"`

	// CPURequestsFraction is a fraction of CPU, that is allocated.
	CPURequestsFraction float64 `json:"cpuRequestsFraction"`

	// CPULimits is defined CPU limit.
	CPULimits int64 `json:"cpuLimits"`

	// CPULimitsFraction is a fraction of defined CPU limit, can be over 100%, i.e.
	// overcommitted.
	CPULimitsFraction float64 `json:"cpuLimitsFraction"`

	// CPUCapacity is specified node CPU capacity in milicores.
	CPUCapacity int64 `json:"cpuCapacity"`

	// MemoryRequests is a fraction of memory, that is allocated.
	MemoryRequests int64 `json:"memoryRequests"`

	// MemoryRequestsFraction is a fraction of memory, that is allocated.
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`

	// MemoryLimits is defined memory limit.
	MemoryLimits int64 `json:"memoryLimits"`

	// MemoryLimitsFraction is a fraction of defined memory limit, can be over 100%, i.e.
	// overcommitted.
	MemoryLimitsFraction float64 `json:"memoryLimitsFraction"`

	// MemoryCapacity is specified node memory capacity in bytes.
	MemoryCapacity int64 `json:"memoryCapacity"`

	// AllocatedPods in number of currently allocated pods on the node.
	AllocatedPods int `json:"allocatedPods"`

	// PodCapacity is maximum number of pods, that can be allocated on the node.
	PodCapacity int64 `json:"podCapacity"`

	// PodFraction is a fraction of pods, that can be allocated on given node.
	PodFraction float64 `json:"podFraction"`
}

type nodeService struct {
	param  *kubernetes.KubernetesParam
	podCtl PodService
}

func NewNodeService(
	param *kubernetes.KubernetesParam,
) NodeService {
	return &nodeService{
		param:  param,
		podCtl: NewPodService(param),
	}
}

func toCells(std []NodeDetail) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = NodeCell(std[i])
	}
	return cells
}

func fromCells(cells []page.DataCell) []NodeDetail {
	std := make([]NodeDetail, len(cells))
	for i := range std {
		std[i] = NodeDetail(cells[i].(NodeCell))
	}
	return std
}

func getNodeAllocatedResources(node *corev1.Node, podList []corev1.Pod) (NodeAllocatedResources, error) {
	reqs, limits := map[corev1.ResourceName]resource.Quantity{}, map[corev1.ResourceName]resource.Quantity{}

	for _, pod := range podList {
		podReqs, podLimits, err := PodRequestsAndLimits(&pod)
		if err != nil {
			return NodeAllocatedResources{}, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}

	cpuRequests, cpuLimits, memoryRequests, memoryLimits := reqs[corev1.ResourceCPU],
		limits[corev1.ResourceCPU], reqs[corev1.ResourceMemory], limits[corev1.ResourceMemory]

	var cpuRequestsFraction, cpuLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Cpu().MilliValue()); capacity > 0 {
		cpuRequestsFraction = float64(cpuRequests.MilliValue()) / capacity * 100
		cpuLimitsFraction = float64(cpuLimits.MilliValue()) / capacity * 100
	}

	var memoryRequestsFraction, memoryLimitsFraction float64 = 0, 0
	if capacity := float64(node.Status.Capacity.Memory().MilliValue()); capacity > 0 {
		memoryRequestsFraction = float64(memoryRequests.MilliValue()) / capacity * 100
		memoryLimitsFraction = float64(memoryLimits.MilliValue()) / capacity * 100
	}

	var podFraction float64 = 0
	var podCapacity int64 = node.Status.Capacity.Pods().Value()
	if podCapacity > 0 {
		podFraction = float64(len(podList)) / float64(podCapacity) * 100
	}

	return NodeAllocatedResources{
		CPURequests:            cpuRequests.MilliValue(),
		CPURequestsFraction:    cpuRequestsFraction,
		CPULimits:              cpuLimits.MilliValue(),
		CPULimitsFraction:      cpuLimitsFraction,
		CPUCapacity:            node.Status.Capacity.Cpu().MilliValue(),
		MemoryRequests:         memoryRequests.Value(),
		MemoryRequestsFraction: memoryRequestsFraction,
		MemoryLimits:           memoryLimits.Value(),
		MemoryLimitsFraction:   memoryLimitsFraction,
		MemoryCapacity:         node.Status.Capacity.Memory().Value(),
		AllocatedPods:          len(podList),
		PodCapacity:            podCapacity,
		PodFraction:            podFraction,
	}, nil
}

func PodRequestsAndLimits(pod *corev1.Pod) (reqs, limits corev1.ResourceList, err error) {
	reqs, limits = corev1.ResourceList{}, corev1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	// Add overhead for running a pod to the sum of requests and to non-zero limits:
	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}

func addResourceList(list, new corev1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

func maxResourceList(list, new corev1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}

func (n *nodeService) Get(name string) (*NodeDetail, error) {
	var detail NodeDetail
	node, err := n.param.KubernetesClient.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	detail.Node = *node
	return &detail, nil
}

func (n *nodeService) List(dsQuery *page.DataSelectQuery) (*NodeResponse, error) {
	var nodeDetailList []NodeDetail
	var nodeResponse NodeResponse
	nodeList, err := n.param.KubernetesClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, node := range nodeList.Items {
		var nodeDetail NodeDetail
		nodeDetail.Node = node
		nodeDetailList = append(nodeDetailList, nodeDetail)
	}

	nodeCells, filteredTotal := page.GenericDataSelectWithFilter(toCells(nodeDetailList), dsQuery)
	nodes := fromCells(nodeCells)

	var nodeWithResources []NodeDetail

	defer func() {
		if err != nil {
			err = fmt.Errorf("error when list nodes |%v", err)
			return
		}
	}()

	var (
		wg  sync.WaitGroup
		gp  *ants.Pool
		lth = len(nodes)
		ech = make(chan error, lth)
		pch = make(chan NodeDetail, lth)
	)
	gp, err = ants.NewPool(20)
	if err != nil {
		err = fmt.Errorf("fail to new goroutine pool, caused by: %v", err)
		return &nodeResponse, err
	}
	defer gp.Release()

	for _, item := range nodes {
		wg.Add(1)
		item := item
		err = gp.Submit(func() {
			defer wg.Done()

			var detail NodeDetail
			detail.Node = item.Node
			pch <- detail
		})
		if err != nil {
			err = fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err)
			return &nodeResponse, err
		}
	}
	wg.Wait()
	close(pch)
	select {
	case err = <-ech:
		return &nodeResponse, err
	default:
		// do nothing
	}

	for p := range pch {
		nodeWithResources = append(nodeWithResources, p)
	}

	nodeResponse.List = nodeWithResources
	nodeResponse.Current = dsQuery.PaginationQuery.Page + 1
	nodeResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	nodeResponse.Total = filteredTotal
	return &nodeResponse, nil
}

func (n *nodeService) Create(obj *corev1.Node) error {
	_, err := n.param.KubernetesClient.CoreV1().Nodes().Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (n *nodeService) GetEvents(name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	node, err := n.Get(name)
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	groupVersion := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddKnownTypes(groupVersion, &corev1.Node{})
	eventList, err := n.param.KubernetesClient.CoreV1().Events(corev1.NamespaceAll).Search(scheme, &node.Node)
	if err != nil {
		return nil, err
	}

	var eventResponse EventResponse
	events := eventList.Items
	eventCells, filteredTotal := page.GenericDataSelectWithFilter(EventToCells(events), dsQuery)
	dps := EventFromCells(eventCells)

	var eventDetailList []EventDetail
	for _, tmp := range dps {
		eventDetailList = append(eventDetailList, toEvent(tmp))
	}

	eventResponse.List = eventDetailList
	eventResponse.Current = dsQuery.PaginationQuery.Page + 1
	eventResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	eventResponse.Total = filteredTotal
	return &eventResponse, nil
}

func (n *nodeService) GetPods(name string, dsQuery *page.DataSelectQuery) (*PodResponse, error) {
	var podResponse PodResponse

	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + name +
		",status.phase!=" + string(v1.PodSucceeded) +
		",status.phase!=" + string(v1.PodFailed))

	if err != nil {
		return nil, err
	}

	podList, err := n.param.KubernetesClient.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metaV1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
	if err != nil {
		return nil, err
	}

	podCells, filteredTotal := page.GenericDataSelectWithFilter(ToCells(podList.Items), dsQuery)
	ps := FromCells(podCells)

	var podDetailList []PodDetail

	for _, po := range ps {
		var detail PodDetail
		detail.Pod = po
		detail.PodPhase = getPodStatus(po)
		podDetailList = append(podDetailList, detail)
	}

	podResponse.List = podDetailList
	podResponse.Current = dsQuery.PaginationQuery.Page + 1
	podResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	podResponse.Total = filteredTotal

	return &podResponse, nil
}

func (n *nodeService) GetRawPods(name string) (*corev1.PodList, error) {
	return n.param.KubernetesClient.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", name)})
}

func (n *nodeService) Update(obj *corev1.Node) (*corev1.Node, error) {
	node, err := n.param.KubernetesClient.CoreV1().Nodes().Update(context.TODO(), obj, metav1.UpdateOptions{})
	return node, err
}

func (n *nodeService) Patch(name string, param ReplaceNodeInfoParam) (*NodeDetail, error) {
	node, err := n.param.KubernetesClient.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	dataType := param.OperatorPath
	operatorData := param.OperatorData
	operatorType := "replace"
	operatorPath := fmt.Sprintf("/metadata/%s", param.OperatorPath)

	_, err = util.IsContain(dataType, []string{"labels", "annotations"})
	if err != nil {
		return nil, fmt.Errorf("Unsurported path %s ", dataType)
	}

	var payloads []interface{}

	payload := PatchStringValue{
		Op:    operatorType,
		Path:  operatorPath,
		Value: operatorData,
	}

	payloads = append(payloads, payload)

	payloadBytes, _ := json.Marshal(payloads)

	_, err = n.param.KubernetesClient.CoreV1().Nodes().Patch(context.TODO(), name, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	var nodeDetail NodeDetail
	if err != nil {
		return nil, err
	}
	nodeDetail.Node = *node

	return &nodeDetail, err
}

func (n *nodeService) Delete(name string) error {
	return n.param.KubernetesClient.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (n *nodeService) CordonOrUnCordon(nodeName string, drain bool) (*NodeDetail, error) {
	var nodeDetail NodeDetail

	data := fmt.Sprintf(`{"spec":{"unschedulable":%t}}`, drain)
	node, err := n.param.KubernetesClient.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	nodeDetail.Node = *node
	return &nodeDetail, nil
}

func (n *nodeService) Drain(name string) error {
	return nil
}

func (n *nodeService) PatchTaint(nodeName string, taints []corev1.Taint) (*NodeDetail, error) {
	var nodeDetail NodeDetail
	node, err := n.param.KubernetesClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	newNode := node.DeepCopy()
	newNode.Spec.Taints = taints

	if _, err = n.param.KubernetesClient.CoreV1().Nodes().Update(context.TODO(), newNode, metav1.UpdateOptions{}); err != nil {
		log.Infof("Failed to update node object: %v", err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	nodeDetail.Node = *node
	return &nodeDetail, nil
}

func (n *nodeService) TaintOrUnTaint(nodeName, op string, taint *corev1.Taint) (*NodeDetail, error) {
	var (
		node       *corev1.Node
		err        error
		updated    bool
		nodeDetail NodeDetail
	)

	switch op {
	case "apply":
		node, err = n.param.KubernetesClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		node, updated = addOrUpdateTaint(node, taint)
		log.Infof("Node %q taints after removal; updated %v: %v", nodeName, updated, node.Spec.Taints)

		if updated {
			if _, err = n.param.KubernetesClient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{}); err != nil {
				log.Infof("Failed to update node object: %v", err)
				return nil, err
			}
		}
	case "remove":
		node, err = n.param.KubernetesClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Failed to remove taint: %v", err)
			return nil, err
		}
		var updated bool
		node, updated = removeTaint(node, taint)
		if updated {
			if _, err = n.param.KubernetesClient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{}); err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("Unsupported operations. ")
	}

	if err != nil {
		return nil, err
	}
	nodeDetail.Node = *node
	return &nodeDetail, nil
}

func addOrUpdateTaint(node *corev1.Node, taint *corev1.Taint) (*corev1.Node, bool) {
	newNode := node.DeepCopy()
	nodeTaints := newNode.Spec.Taints

	var newTaints []corev1.Taint
	updated := false
	for i := range nodeTaints {
		if taint.MatchTaint(&nodeTaints[i]) {
			if equality.Semantic.DeepEqual(*taint, nodeTaints[i]) {
				return newNode, false
			}
			newTaints = append(newTaints, *taint)
			updated = true
			continue
		}

		newTaints = append(newTaints, nodeTaints[i])
	}

	if !updated {
		newTaints = append(newTaints, *taint)
	}

	newNode.Spec.Taints = newTaints
	return newNode, true
}

func removeTaint(node *corev1.Node, taint *corev1.Taint) (*corev1.Node, bool) {
	newNode := node.DeepCopy()
	nodeTaints := newNode.Spec.Taints
	if len(nodeTaints) == 0 {
		return newNode, false
	}

	if !taintExists(nodeTaints, taint) {
		return newNode, false
	}

	newTaints, _ := deleteTaint(nodeTaints, taint)
	newNode.Spec.Taints = newTaints
	return newNode, true
}

func taintExists(taints []corev1.Taint, taintToFind *corev1.Taint) bool {
	for _, taint := range taints {
		if taint.MatchTaint(taintToFind) {
			return true
		}
	}
	return false
}

func deleteTaint(taints []corev1.Taint, taintToDelete *corev1.Taint) ([]corev1.Taint, bool) {
	var newTaints []corev1.Taint
	deleted := false
	for i := range taints {
		if taintToDelete.MatchTaint(&taints[i]) {
			deleted = true
			continue
		}
		newTaints = append(newTaints, taints[i])
	}
	return newTaints, deleted
}
