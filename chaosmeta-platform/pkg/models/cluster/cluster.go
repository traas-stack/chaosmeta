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

package cluster

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
)

type Cluster struct {
	ID         int    `json:"id" orm:"pk;auto;column(id)"`
	Name       string `json:"name" orm:"unique;index;column(name);size(255)"`
	KubeConfig string `json:"kubeConfig" orm:"column(kube_config);type(text)"`
	//AppKey     string `json:"appKey" orm:"column(app_key);size(255)"`
	Version string `json:"version" orm:"column(version);size(32);index"`
	models.BaseTimeModel
}

func (c *Cluster) TableName() string {
	return "cluster"
}

func InsertCluster(ctx context.Context, cluster *Cluster) (int64, error) {
	return models.GetORM().Insert(cluster)
}

func GetClusterById(ctx context.Context, cluster *Cluster) error {
	if cluster == nil {
		return errors.New("cluster is nil")
	}
	return models.GetORM().Read(cluster)
}

func GetClustersByIdList(ctx context.Context, ids []int, orderBy string, page, pageSize int) (int64, []*Cluster, error) {
	if len(ids) == 0 {
		return 0, nil, errors.New("empty ids")
	}
	var clusters []*Cluster
	cluster := Cluster{}
	query := models.GetORM().QueryTable(cluster.TableName()).Filter("id__in", ids)

	var totalCount int64
	totalCount, err := query.Count()
	if err != nil {
		return totalCount, nil, err
	}
	query = query.Limit(pageSize, (page-1)*pageSize)
	if len(orderBy) > 0 {
		query = query.OrderBy(orderBy)
	}
	_, err = query.All(clusters)
	return totalCount, clusters, err
}

func GetClusterByName(ctx context.Context, cluster *Cluster) error {
	if cluster == nil {
		return errors.New("cluster is nil")
	}
	return models.GetORM().Read(cluster, "name")
}

func UpdateCluster(ctx context.Context, cluster *Cluster) (int64, error) {
	if cluster == nil {
		return 0, errors.New("cluster is nil")
	}
	num, err := models.GetORM().Update(cluster)
	return num, err
}

func QueryCluster(ctx context.Context, name, version, orderBy string, page, pageSize int) (int64, []Cluster, error) {
	c, clusters := Cluster{}, new([]Cluster)
	querySeter := models.GetORM().QueryTable(c.TableName())
	clusterQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}
	if len(name) > 0 {
		clusterQuery.Filter("name", models.CONTAINS, true, name)
	}
	if len(version) > 0 {
		clusterQuery.Filter("version", models.NEGLECT, false, version)
	}
	var totalCount int64
	totalCount, err = clusterQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	clusterQuery.OrderBy(orderByList...)

	if err := clusterQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = clusterQuery.GetOamQuerySeter().All(clusters)
	return totalCount, *clusters, err
}

func ListCluster() ([]Cluster, error) {
	c, clusters := Cluster{}, new([]Cluster)
	querySeter := models.GetORM().QueryTable(c.TableName())
	clusterQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return nil, err
	}
	_, err = clusterQuery.GetOamQuerySeter().All(clusters)
	return *clusters, err
}

func DeleteClustersByIdList(ctx context.Context, ids []int) error {
	cluster := Cluster{}
	querySeter := models.GetORM().QueryTable(cluster.TableName())
	clusterQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return err
	}
	clusterQuery.Filter("id", models.IN, false, ids)
	_, err = clusterQuery.Delete()
	return err
}

func GetDefaultCluster(ctx context.Context, cluster *Cluster) error {
	if cluster == nil {
		return errors.New("cluster is nil")
	}
	return models.GetORM().QueryTable(cluster.TableName()).Filter("name", "noKubernetes").One(cluster)
}
