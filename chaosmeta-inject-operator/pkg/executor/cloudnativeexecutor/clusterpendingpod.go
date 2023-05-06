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

package cloudnativeexecutor

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.ClusterCloudTarget, faultClusterPendingPod, &ClusterPendingPodExecutor{})
	registerResourceCreateFunc(v1alpha1.ClusterCloudTarget, faultClusterPendingPod, createPendingPod)
}

type ClusterPendingPodExecutor struct{}

const (
	pendingImage           = "chaosmetapending:v1"
	pendingContainer       = "chaosmeta"
	resourceReq            = "1000000000000000000"
	faultClusterPendingPod = "pendingpod"
)

func createPendingPod(ctx context.Context, namespace, name string) error {
	quota, _ := resource.ParseQuantity(resourceReq)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  pendingContainer,
					Image: pendingImage,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    quota,
							corev1.ResourceMemory: quota,
						},
					},
				},
			},
		},
	}

	return restclient.GetApiServerClientMap(v1alpha1.PodCloudTarget).Post().Resource("pods").Namespace(pod.Namespace).Body(pod).Do(ctx).Error()
}

func (e *ClusterPendingPodExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	return batchResourceInject(ctx, args, v1alpha1.ClusterCloudTarget, faultClusterPendingPod)
}

func (e *ClusterPendingPodExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	common.GetClusterCtrl().Stop()

	if err := deleteNs(ctx, injectObject); err != nil {
		return fmt.Errorf("delete namespace error: %s", err.Error())
	}

	return nil
}

func (e *ClusterPendingPodExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	status := v1alpha1.SuccessStatusType
	if phase == v1alpha1.InjectPhaseType {
		if common.GetClusterCtrl().IsRunning() {
			status = v1alpha1.RunningStatusType
		}
	}

	return &model.SubExpInfo{
		UID:        uid,
		Status:     status,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
