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

package clientset

import (
	cv1alpha1 "chaosmeta-platform/pkg/gateway/apis/chaosmetacluster/v1alpha1"
	"chaosmeta-platform/pkg/models/common/page"
	"errors"
	"fmt"
)

type ClusterCell cv1alpha1.ChaosmetaCluster

func (n ClusterCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(n.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(n.ObjectMeta.CreationTimestamp.Time)
	default:
		return nil
	}
}

func toCells(std []cv1alpha1.ChaosmetaCluster) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = ClusterCell(std[i])
	}
	return cells
}

func fromCells(cells []page.DataCell) []cv1alpha1.ChaosmetaCluster {
	std := make([]cv1alpha1.ChaosmetaCluster, len(cells))
	for i := range std {
		std[i] = cv1alpha1.ChaosmetaCluster(cells[i].(ClusterCell))
	}
	return std
}

// Get a list of all clusters without pagination
func (cs *clientset) List(env string) (*cv1alpha1.ChaosmetaClusterList, error) {
	clusterList, err := cs.ListClusters(env)
	if err != nil {
		return nil, err
	}

	return &cv1alpha1.ChaosmetaClusterList{Items: clusterList}, nil
}

func (cs *clientset) ListCluster(env string, dsQuery *page.DataSelectQuery) (*ClusterListResponse, error) {
	var (
		clusterList []cv1alpha1.ChaosmetaCluster
		csResponse  ClusterListResponse
		err         error
	)

	clusterList, err = cs.ListClusters(env)
	if err != nil {
		return nil, err
	}

	clusterCells, filteredTotal := page.GenericDataSelectWithFilter(toCells(clusterList), dsQuery)
	cls := fromCells(clusterCells)

	csResponse.List = cls
	csResponse.Current = dsQuery.PaginationQuery.Page + 1
	csResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	csResponse.Total = filteredTotal
	return &csResponse, nil
}

func (cs *clientset) GetCluster(env, clusterName string) (*cv1alpha1.ChaosmetaCluster, error) {
	clusterList, err := cs.ListClusters(env)
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusterList {
		if cluster.GetName() == clusterName {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("no clusters")
}

func (cs *clientset) GetClusterByClusterName(clusterName string) (*cv1alpha1.ChaosmetaCluster, error) {
	clusterList, err := cs.ListClusters("")
	if err != nil {
		return nil, err
	}
	for _, cluster := range clusterList {
		if cluster.GetName() == clusterName {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("no cluster")
}

func (cs *clientset) CreateCluster(cluster *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error) {
	return nil, fmt.Errorf("can not add cluster")
}

func (cs *clientset) DeleteCluster(env, cluster string) error {
	return fmt.Errorf("can not delete cluster")
}

func (cs *clientset) PatchCluster(originalObj, updatedObj *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error) {
	//updatedObj.TypeMeta = originalObj.TypeMeta
	//labels := updatedObj.ObjectMeta.Labels
	//updatedObj.ObjectMeta = originalObj.ObjectMeta
	//updatedObj.ObjectMeta.Labels = labels
	//
	//data, err := json.MargePatch(originalObj, updatedObj)
	//
	//if err != nil {
	//	return nil, err
	//}

	return nil, errors.New("cannot patch")
}

// Replace cluster meta information
func (cs *clientset) ReplaceCluster(originalObj, updatedObj *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error) {
	if originalObj == nil {
		return cs.CreateCluster(updatedObj)
	}

	return cs.PatchCluster(originalObj, updatedObj)
}

// Home cluster display dashboard
func (cs *clientset) ListClusterDashboardInfo(env string) (dashboard []ClusterDashboardInfo, err error) {
	return
}
