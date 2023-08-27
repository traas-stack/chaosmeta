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
	"context"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CRDService interface {
	// todo lacks an interface to obtain resources under a certain crd
	Get(name string) (*v1.CustomResourceDefinition, error)
	List(dsQuery *page.DataSelectQuery) (*CrdListResponse, error)
}

type CrdCell v1.CustomResourceDefinition

func (n CrdCell) GetProperty(name page.PropertyName) page.ComparableValue {
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

func (c *crdService) toCells(std []v1.CustomResourceDefinition) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = CrdCell(std[i])
	}
	return cells
}

func (c *crdService) fromCells(cells []page.DataCell) []v1.CustomResourceDefinition {
	std := make([]v1.CustomResourceDefinition, len(cells))
	for i := range std {
		std[i] = v1.CustomResourceDefinition(cells[i].(CrdCell))
	}
	return std
}

type CrdListResponse struct {
	Total    int                           `json:"total"`
	Current  int                           `json:"current"`
	PageSize int                           `json:"pageSize"`
	List     []v1.CustomResourceDefinition `json:"list"`
}

type crdService struct {
	client apiextensionsclientset.Interface
}

func NewCRDService(
	client apiextensionsclientset.Interface,
) CRDService {
	return &crdService{
		client: client,
	}
}

func (c *crdService) Get(name string) (*v1.CustomResourceDefinition, error) {
	return c.client.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *crdService) List(dsQuery *page.DataSelectQuery) (*CrdListResponse, error) {
	var crdListResponse CrdListResponse
	crds, err := c.client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	crdCells, filteredTotal := page.GenericDataSelectWithFilter(c.toCells(crds.Items), dsQuery)
	dps := c.fromCells(crdCells)

	crdListResponse.List = dps
	crdListResponse.Current = dsQuery.PaginationQuery.Page + 1
	crdListResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	crdListResponse.Total = filteredTotal
	return &crdListResponse, nil
}
