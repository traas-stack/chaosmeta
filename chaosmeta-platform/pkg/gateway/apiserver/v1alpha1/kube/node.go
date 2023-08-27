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
	kubernetesService "chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/pkg/service/kubernetes/kube"
	"context"
)

func (c *KubeController) ListNodes() {
	id, _ := c.GetInt(":id")
	clusterService := cluster.ClusterService{}
	kubeClient, restConfig, err := clusterService.GetRestConfig(context.Background(), id)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	ns := kube.NewNodeService(&kubernetesService.KubernetesParam{KubernetesClient: kubeClient, RestConfig: restConfig})
	resp, err := ns.List(page.ParseDataSelectPathParameter(&c.Controller))
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, resp)
}
