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

// ServiceService defines the interface contains service manages methods.
type ServiceService interface {
	List(namespace string, opts metav1.ListOptions) (*corev1.ServiceList, error)
	Get(namespace, name string) (*corev1.Service, error)
	Create(svc *corev1.Service) error
	Update(svc *corev1.Service) error
	Delete(svc *corev1.Service) error
	Patch(originalObj, updatedObj *corev1.Service) error
	Replace(originalObj, updatedObj *corev1.Service) error
}

type serviceService struct {
	kubeClient kubernetes.Interface
}

// NewServiceService returns an instance of service Service.
func NewServiceService(kubeClient kubernetes.Interface) ServiceService {
	return &serviceService{kubeClient}
}

func (s *serviceService) List(namespace string, opts metav1.ListOptions) (*corev1.ServiceList, error) {
	return s.kubeClient.CoreV1().Services(namespace).List(context.TODO(), opts)
}

func (s *serviceService) Get(namespace, name string) (*corev1.Service, error) {
	return s.kubeClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (s *serviceService) Create(svc *corev1.Service) error {
	_, err := s.kubeClient.CoreV1().Services(svc.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
	return err
}

func (s *serviceService) Update(svc *corev1.Service) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := s.kubeClient.CoreV1().Services(svc.Namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
		return err
	})
}

func (s *serviceService) Delete(svc *corev1.Service) error {
	return s.kubeClient.CoreV1().Services(svc.GetNamespace()).Delete(context.TODO(), svc.GetName(), metav1.DeleteOptions{})
}

func (s *serviceService) Patch(originalObj, updatedObj *corev1.Service) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = s.kubeClient.CoreV1().Services(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return err
}

func (s *serviceService) Replace(originalObj, updatedObj *corev1.Service) error {
	if originalObj == nil {
		return s.Create(updatedObj)
	}

	return s.Patch(originalObj, updatedObj)
}
