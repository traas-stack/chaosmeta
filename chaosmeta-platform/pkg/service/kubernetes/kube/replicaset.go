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
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ReplicaSetService interface {
	List(namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error)
	Get(namespace, name string) (*appsv1.ReplicaSet, error)
}

type replicaSetService struct {
	kubeClient kubernetes.Interface
}

func NewReplicaSetService(kubeClient kubernetes.Interface) ReplicaSetService {
	return &replicaSetService{
		kubeClient: kubeClient,
	}
}

func (r *replicaSetService) List(namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error) {
	return r.kubeClient.AppsV1().ReplicaSets(namespace).List(context.TODO(), opts)
}

func (r *replicaSetService) Get(namespace, name string) (*appsv1.ReplicaSet, error) {
	return r.kubeClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
