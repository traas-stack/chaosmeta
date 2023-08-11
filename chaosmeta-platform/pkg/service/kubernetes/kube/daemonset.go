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
	"chaosmeta-platform/pkg/service/kubernetes/common"
	"chaosmeta-platform/util/json"
	"context"
	"fmt"
	"github.com/panjf2000/ants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

// DaemonsetService defines the interface contains daemonest manages methods.
type DaemonsetService interface {
	Get(namespace, name string) (*DaemonSetDetail, error)
	List(namespace string, dsQuery *page.DataSelectQuery) (*DaemonSetResponse, error)
	Create(ss *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Update(ss *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Patch(originalObj, updatedObj *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Replace(originalObj, updatedObj *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	GetRawPods(namespace, name string) ([]corev1.Pod, error)
	GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error)
	GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
}

type daemonsetService struct {
	param   *kubernetes.KubernetesParam
	podCtrl PodService
}

type DaemonsetCell appsv1.DaemonSet

func (n DaemonsetCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(n.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(n.ObjectMeta.CreationTimestamp.Time)
	case page.NamespaceProperty:
		return page.StdComparableString(n.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func (d *daemonsetService) toCells(std []appsv1.DaemonSet) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = DaemonsetCell(std[i])
	}
	return cells
}

func (d *daemonsetService) fromCells(cells []page.DataCell) []appsv1.DaemonSet {
	std := make([]appsv1.DaemonSet, len(cells))
	for i := range std {
		std[i] = appsv1.DaemonSet(cells[i].(DaemonsetCell))
	}
	return std
}

type DaemonSetResponse struct {
	Total    int               `json:"total"`
	Current  int               `json:"current"`
	PageSize int               `json:"pageSize"`
	List     []DaemonSetDetail `json:"list"`
}

type DaemonSetDetail struct {
	appsv1.DaemonSet  `json:",inline"`
	PodStatusInfo     common.PodStatusInfo     `json:"podStatusInfo"`
	ReplicaStatusInfo common.ReplicaStatusInfo `json:"replicaStatusInfo"`
}

// NewDaemonSetService returns an instance of DaemonsetService.
func NewDaemonSetService(param *kubernetes.KubernetesParam) DaemonsetService {
	return &daemonsetService{param: param, podCtrl: NewPodService(param)}
}

func (d *daemonsetService) getDaemonSetPodStatus(daemonSet *appsv1.DaemonSet) DaemonSetDetail {
	var daemonSetDetail DaemonSetDetail
	daemonSetDetail.DaemonSet = *daemonSet
	//daemonSetDetail.PodStatusInfo = podInfo
	daemonSetDetail.ReplicaStatusInfo.Desired = daemonSet.Status.DesiredNumberScheduled
	daemonSetDetail.ReplicaStatusInfo.Available = daemonSet.Status.NumberAvailable
	return daemonSetDetail
}

func (d *daemonsetService) List(namespace string, dsQuery *page.DataSelectQuery) (*DaemonSetResponse, error) {
	var daemonSetResponse DaemonSetResponse
	daemonSets, err := d.param.KubernetesClient.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var daemonSetList []appsv1.DaemonSet
	for _, dd := range daemonSets.Items {
		daemonSetList = append(daemonSetList, dd)
	}

	daemonsetsCells, filteredTotal := page.GenericDataSelectWithFilter(d.toCells(daemonSetList), dsQuery)
	dps := d.fromCells(daemonsetsCells)

	var daemonSetDetailList []DaemonSetDetail

	defer func() {
		if err != nil {
			err = fmt.Errorf("error when list pods by names|%v", err)
			return
		}
	}()

	var (
		wg  sync.WaitGroup
		gp  *ants.Pool
		lth = len(dps)
		ech = make(chan error, lth)
		pch = make(chan DaemonSetDetail, lth)
	)
	gp, err = ants.NewPool(20)
	if err != nil {
		err = fmt.Errorf("fail to new goroutine pool, caused by: %v", err)
		return nil, err
	}
	defer gp.Release()

	for _, dm := range dps {
		wg.Add(1)
		dm := dm
		err = gp.Submit(func() {
			defer wg.Done()
			p := d.getDaemonSetPodStatus(&dm)
			pch <- p
		})
		if err != nil {
			err = fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err)
			return nil, err
		}
	}
	wg.Wait()
	close(pch)
	select {
	case err = <-ech:
		return nil, err
	default:
		// do nothing
	}

	for p := range pch {
		daemonSetDetailList = append(daemonSetDetailList, p)
	}

	daemonSetResponse.List = daemonSetDetailList
	daemonSetResponse.Current = dsQuery.PaginationQuery.Page + 1
	daemonSetResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	daemonSetResponse.Total = filteredTotal
	return &daemonSetResponse, nil
}

func (d *daemonsetService) GetRawPods(namespace, name string) ([]corev1.Pod, error) {
	daemonset, err := d.param.KubernetesClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	rawPods, err := d.podCtrl.ListWithOptions(namespace, daemonset.Spec.Selector)
	if err != nil {
		return nil, err
	}

	return FilterPodsByControllerRef(daemonset, rawPods), nil
}

func (d *daemonsetService) GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error) {
	var podResponse PodResponse

	podList, err := d.GetRawPods(namespace, name)
	if err != nil {
		return nil, err
	}

	podCells, filteredTotal := page.GenericDataSelectWithFilter(ToCells(podList), dsQuery)
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

func (d *daemonsetService) Get(namespace, name string) (*DaemonSetDetail, error) {
	daemonset, err := d.param.KubernetesClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	detail := d.getDaemonSetPodStatus(daemonset)
	return &detail, nil
}

func (d *daemonsetService) Create(ss *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return d.param.KubernetesClient.AppsV1().DaemonSets(ss.Namespace).Create(context.TODO(), ss, metav1.CreateOptions{})
}

func (d *daemonsetService) Update(ss *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return d.param.KubernetesClient.AppsV1().DaemonSets(ss.Namespace).Update(context.TODO(), ss, metav1.UpdateOptions{})
}

func (d *daemonsetService) Patch(originalObj, updatedObj *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return nil, err
	}

	info, err := d.param.KubernetesClient.AppsV1().DaemonSets(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return info, nil
}

func (d *daemonsetService) Replace(originalObj, updatedObj *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	if originalObj == nil {
		return d.Create(updatedObj)
	}

	return d.Patch(originalObj, updatedObj)
}

func (d *daemonsetService) GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	eventCtrl := NewEventService(d.param.KubernetesClient)
	eventResponse, err := eventCtrl.GetResourceEvents(namespace, name, dsQuery)
	if err != nil {
		return nil, err
	}
	return eventResponse, nil
}
