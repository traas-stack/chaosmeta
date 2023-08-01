package namespace

import (
	models "chaosmeta-platform/pkg/models/common"
)

type ClusterNamespace struct {
	ID          int `json:"id" orm:"pk;auto;column(id)"`
	ClusterID   int `json:"clusterID" orm:"column(cluster_id);index"`
	NamespaceID int `json:"namespaceId" orm:"column(namespace_id);index"`
	models.BaseTimeModel
}

func (c *ClusterNamespace) TableName() string {
	return "cluster_namespace"
}

func (c *ClusterNamespace) TableUnique() [][]string {
	return [][]string{{"cluster_id", "namespace_id"}}
}
