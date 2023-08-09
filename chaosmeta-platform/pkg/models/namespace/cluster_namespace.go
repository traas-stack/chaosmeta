package namespace

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/util/log"
	"github.com/beego/beego/v2/client/orm"
	"github.com/spf13/cast"
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

func GetClusterIDsByNamespaceID(namespaceID int) ([]int, error) {
	clusterNamespace := ClusterNamespace{}
	var clusterIDs orm.ParamsList
	_, err := models.GetORM().QueryTable(clusterNamespace.TableName()).Filter("namespace_id", namespaceID).ValuesFlat(&clusterIDs, "cluster_id")
	if err != nil {
		return nil, err
	}
	return cast.ToIntSlice(clusterIDs), nil
}

func SetClusterIDsForNamespace(namespaceID int, clusterIDs []int) error {
	if err := ClearClusterIDsForNamespace(namespaceID); err != nil {
		return err
	}
	for _, clusterID := range clusterIDs {
		if _, err := models.GetORM().Insert(&ClusterNamespace{ClusterID: clusterID, NamespaceID: namespaceID}); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func ClearClusterIDsForNamespace(namespaceID int) error {
	clusterNamespace := ClusterNamespace{}
	_, err := models.GetORM().QueryTable(clusterNamespace.TableName()).Filter("namespace_id", namespaceID).Delete()
	return err
}
