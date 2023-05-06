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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.ClusterCloudTarget, faultClusterCompletedJob, &ClusterCompletedJobExecutor{})
	registerResourceCreateFunc(v1alpha1.ClusterCloudTarget, faultClusterCompletedJob, createJob)
}

type ClusterCompletedJobExecutor struct{}

const (
	jobImage                 = "centos:centos7"
	jobContainer             = "chaosjob"
	faultClusterCompletedJob = "completedjob"
)

func createJob(ctx context.Context, namespace, name string) error {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    jobContainer,
							Image:   jobImage,
							Command: []string{"echo", "ok"},
						},
					},
				},
			},
		},
	}

	return restclient.GetApiServerClientMap(v1alpha1.JobCloudTarget).Post().Resource("jobs").Namespace(job.Namespace).Body(job).Do(ctx).Error()
}

func (e *ClusterCompletedJobExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	return batchResourceInject(ctx, args, v1alpha1.ClusterCloudTarget, faultClusterCompletedJob)
}

func (e *ClusterCompletedJobExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	common.GetClusterCtrl().Stop()
	if err := restclient.GetApiServerClientMap(v1alpha1.NamespaceCloudTarget).Delete().Resource("namespaces").
		Name(injectObject).Do(ctx).Error(); err != nil {
		return fmt.Errorf("create namespace error: %s", err.Error())
	}

	return nil
}

func (e *ClusterCompletedJobExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
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
