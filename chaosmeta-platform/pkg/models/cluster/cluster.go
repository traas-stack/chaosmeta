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
	Version string `json:"version" orm:"column(version);size(32)"`
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
	if err := clusterQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}
	if len(orderBy) > 0 {
		clusterQuery.OrderBy(orderBy)
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
