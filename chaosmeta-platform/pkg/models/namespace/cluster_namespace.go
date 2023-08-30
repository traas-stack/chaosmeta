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
	if err == orm.ErrNoRows {
		return nil, nil
	}
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
