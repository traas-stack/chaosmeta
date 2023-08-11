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
	"chaosmeta-platform/util/json"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type ConfigMapService interface {
	Get(namespace, name string) (*corev1.ConfigMap, error)
	Create(cm *corev1.ConfigMap) error
	Update(cm *corev1.ConfigMap) error
	Patch(originalObj, updatedObj *corev1.ConfigMap) error
	Replace(originalObj, updatedObj *corev1.ConfigMap) error
}

type configMapService struct {
	kubeClient kubernetes.Interface
}

func NewConfigMapService(
	kubeClient kubernetes.Interface,
) ConfigMapService {
	return &configMapService{
		kubeClient: kubeClient,
	}
}

func (cmc *configMapService) Get(namespace, name string) (*corev1.ConfigMap, error) {
	return cmc.kubeClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (cmc *configMapService) Create(cm *corev1.ConfigMap) error {
	_, err := cmc.kubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

func (cmc *configMapService) Update(cm *corev1.ConfigMap) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := cmc.kubeClient.CoreV1().ConfigMaps(cm.GetNamespace()).Update(context.TODO(), cm, metav1.UpdateOptions{})
		return err
	})
}

func (cmc *configMapService) Patch(originalObj, updatedObj *corev1.ConfigMap) error {
	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = cmc.kubeClient.CoreV1().ConfigMaps(originalObj.GetNamespace()).Patch(
		context.TODO(),
		originalObj.GetName(),
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)
	return err
}

func (cmc *configMapService) Replace(originalObj, updatedObj *corev1.ConfigMap) error {
	if originalObj != nil {
		return cmc.Create(updatedObj)
	}

	return cmc.Patch(originalObj, updatedObj)
}
