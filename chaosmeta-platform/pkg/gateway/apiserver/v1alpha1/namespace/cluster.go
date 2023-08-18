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
	"chaosmeta-platform/pkg/service/namespace"
	"context"
	"encoding/json"
)

func (c *NamespaceController) SetAttackableCluster() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var reqBody SetAttackableClusterRequest
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	username := c.Ctx.Input.GetData("userName").(string)
	namespace := &namespace.NamespaceService{}
	if err := namespace.SetAttackableCluster(context.Background(), namespaceId, username, reqBody.ClusterID); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *NamespaceController) ListAttackableCluster() {
	namespaceId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	sort := c.GetString("sort")

	namespace := &namespace.NamespaceService{}
	total, clusterList, err := namespace.GetAttackableClustersByNamespaceID(context.Background(), namespaceId, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	getAttackableClusterResponse := GetAttackableClusterResponse{Total: total, Page: page, PageSize: pageSize}
	for _, cluster := range clusterList {
		getAttackableClusterResponse.Clusters = append(getAttackableClusterResponse.Clusters, ClusterNamespaceInfo{
			ID:         cluster.ID,
			Name:       cluster.Name,
			CreateTime: cluster.CreateTime,
			UpdateTime: cluster.UpdateTime,
		})
	}
	c.Success(&c.Controller, getAttackableClusterResponse)
}
