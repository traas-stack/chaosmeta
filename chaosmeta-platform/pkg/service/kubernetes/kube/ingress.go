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
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// IngressService defines the interface contains ingress manages methods.
type IngressService interface {
	List(namespace string, opts metav1.ListOptions) (*v1beta1.IngressList, error)
	Get(namespace, name string) (*v1beta1.Ingress, error)
	Create(ingress *v1beta1.Ingress) error
	Update(ingress *v1beta1.Ingress) error
	Delete(ingress *v1beta1.Ingress) error
	Patch(originalObj, updatedObj *v1beta1.Ingress) error
	Replace(originalObj, updatedObj *v1beta1.Ingress) error
}

type ingressService struct {
	kubeClient kubernetes.Interface
}

// NewIngressService returns an instance of ingress Service.
func NewIngressService(kubeClient kubernetes.Interface) IngressService {
	return &ingressService{kubeClient}
}

func (ic *ingressService) List(namespace string, opts metav1.ListOptions) (*v1beta1.IngressList, error) {
	return ic.kubeClient.ExtensionsV1beta1().Ingresses(namespace).List(context.TODO(), opts)
}

func (ic *ingressService) Get(namespace, name string) (*v1beta1.Ingress, error) {
	return ic.kubeClient.ExtensionsV1beta1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (ic *ingressService) Create(ingress *v1beta1.Ingress) error {
	_, err := ic.kubeClient.ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Create(context.TODO(), ingress, metav1.CreateOptions{})
	return err
}

func (ic *ingressService) Update(ingress *v1beta1.Ingress) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := ic.kubeClient.ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Update(context.TODO(), ingress, metav1.UpdateOptions{})
		return err
	})
}

func (ic *ingressService) Delete(ingress *v1beta1.Ingress) error {
	return ic.kubeClient.ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Delete(context.TODO(), ingress.GetName(), metav1.DeleteOptions{})
}

func (ic *ingressService) Patch(originalObj, updatedObj *v1beta1.Ingress) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = ic.kubeClient.ExtensionsV1beta1().Ingresses(originalObj.GetNamespace()).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return err
}

func (ic *ingressService) Replace(originalObj, updatedObj *v1beta1.Ingress) error {
	if originalObj == nil {
		return ic.Create(updatedObj)
	}

	return ic.Patch(originalObj, updatedObj)
}
