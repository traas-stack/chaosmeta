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
	"chaosmeta-platform/pkg/service/cluster"
	"chaosmeta-platform/pkg/service/kubernetes/kube"
	"context"
	"encoding/json"
)

func (c *KubeController) ListContainers() {
	id, _ := c.GetInt(":id", 0)
	nsName := c.GetString(":ns_name")
	clusterService := cluster.ClusterService{}
	kubeClient, config, err := clusterService.GetRestConfig(context.Background(), id)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	var requestBody QueryContainerRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	containerService := kube.NewContainerService(kubeClient, config)
	resp, err := containerService.ListContainers(nsName, requestBody.TargetPods, requestBody.TargetLabel, page.ParseDataSelectPathParameter(&c.Controller))
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, resp)
}
