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
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	registerCloudExecutor(v1alpha1.PodCloudTarget, "containerpause", &PodContainerPauseExecutor{})
}

type PodContainerPauseExecutor struct{}

func (e *PodContainerPauseExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	// parse experiment object info
	ns, name, containerName, err := model.ParseContainerInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected pod format: %s", err.Error())
	}

	if containerName == "" {
		return "", fmt.Errorf("container name not provide")
	}

	// get container id and host ip
	c := restclient.GetApiServerClientMap(v1alpha1.PodCloudTarget)
	pod := &corev1.Pod{}
	if err := c.Get().Namespace(ns).Resource("pods").Name(name).Do(ctx).Into(pod); err != nil {
		return "", fmt.Errorf("get pod error: %s", err.Error())
	}

	hostIP := pod.Status.HostIP
	containers, err := selector.GetTargetContainers(containerName, pod.Status.ContainerStatuses)
	if err != nil || len(containers) == 0 {
		return "", fmt.Errorf("get target container[%s] in pod[%s] error: %s", containerName, pod.Name, err.Error())
	}

	return hostIP, remoteexecutor.GetRemoteExecutor().Inject(ctx, hostIP, "container", "pause", uid, timeout, containers[0].ContainerId, containers[0].ContainerRuntime, nil)
}

func (e *PodContainerPauseExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	return remoteexecutor.GetRemoteExecutor().Recover(ctx, backup, uid)
}
func (e *PodContainerPauseExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return remoteexecutor.GetRemoteExecutor().Query(ctx, backup, uid, phase)
}
