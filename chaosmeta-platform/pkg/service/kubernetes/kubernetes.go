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

package kubernetes

import (
	"chaosmeta-platform/pkg/models/cluster"
	"chaosmeta-platform/util/log"
	"context"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Init() {
	ctx := context.Background()
	defaultCluster := cluster.Cluster{}
	if err := cluster.GetDefaultCluster(ctx, &defaultCluster); err == nil {
		return
	}

	noKubernetes := &cluster.Cluster{
		Name: "noKubernetes",
	}
	_, err := cluster.InsertCluster(ctx, noKubernetes)
	if err != nil {
		log.Panic(err)
	}
}

type KubernetesParam struct {
	Cluster          string
	RestConfig       *rest.Config
	KubernetesClient kubernetes.Interface
	Factory          informers.SharedInformerFactory
}
