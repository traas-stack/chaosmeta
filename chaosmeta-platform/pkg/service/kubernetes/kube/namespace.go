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
	"chaosmeta-platform/util/json"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type NamespaceService interface {
	List(opts metav1.ListOptions, dsQuery *page.DataSelectQuery) (*NamespaceResponse, error)
	Get(name string) (*corev1.Namespace, error)
	Create(namespace *corev1.Namespace) (*corev1.Namespace, error)
	Update(namespace *corev1.Namespace) error
	Delete(name string) error
	Patch(originalObj, updatedObj *corev1.Namespace) (*corev1.Namespace, error)
}

type namespaceService struct {
	kubeClient kubernetes.Interface
}

// NewNamespaceService  returns an instance of namespace Service.
func NewNamespaceService(kubeClient kubernetes.Interface) NamespaceService {
	return &namespaceService{
		kubeClient: kubeClient,
	}
}

type NamespaceResponse struct {
	Total    int                `json:"total"`
	Current  int                `json:"current"`
	PageSize int                `json:"pageSize"`
	List     []corev1.Namespace `json:"list"`
}

type NamespaceCell corev1.Namespace

func (n NamespaceCell) GetProperty(name page.PropertyName) page.ComparableValue {
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

func (ns *namespaceService) toCells(std []corev1.Namespace) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = NamespaceCell(std[i])
	}
	return cells
}

func (ns *namespaceService) fromCells(cells []page.DataCell) []corev1.Namespace {
	std := make([]corev1.Namespace, len(cells))
	for i := range std {
		std[i] = corev1.Namespace(cells[i].(NamespaceCell))
	}
	return std
}

func (ns *namespaceService) List(opts metav1.ListOptions, dsQuery *page.DataSelectQuery) (*NamespaceResponse, error) {
	var namespaceResponse NamespaceResponse
	namespaces, err := ns.kubeClient.CoreV1().Namespaces().List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	namespaceList := namespaces.Items
	namespaceCells, filteredTotal := page.GenericDataSelectWithFilter(ns.toCells(namespaceList), dsQuery)
	nss := ns.fromCells(namespaceCells)

	namespaceResponse.List = nss
	namespaceResponse.Current = dsQuery.PaginationQuery.Page + 1
	namespaceResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	namespaceResponse.Total = filteredTotal
	return &namespaceResponse, nil
}

func (ns *namespaceService) Get(name string) (*corev1.Namespace, error) {
	return ns.kubeClient.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

func (ns *namespaceService) Create(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	namespace, err := ns.kubeClient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	return namespace, err
}

func (ns *namespaceService) Update(namespace *corev1.Namespace) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := ns.kubeClient.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
		return err
	})
}

func (ns *namespaceService) Delete(name string) error {
	return ns.kubeClient.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (ns *namespaceService) Patch(originalObj, updatedObj *corev1.Namespace) (*corev1.Namespace, error) {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return nil, err
	}

	namespace, err := ns.kubeClient.CoreV1().Namespaces().Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return namespace, err
}
