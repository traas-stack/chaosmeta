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
	"bytes"
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.PodCloudTarget, "containerimage", &PodContainerImageExecutor{})
}

type PodContainerImageExecutor struct{}

func (e *PodContainerImageExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	ns, name, containerName, err := model.ParseContainerInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected pod format: %s", err.Error())
	}

	reArgs := common.GetArgs(args, []string{"image"})
	var newImage = reArgs[0]

	if containerName == "" {
		return "", fmt.Errorf("container name not provide")
	}
	if newImage == "" {
		return "", fmt.Errorf("image not provide")
	}

	// get target container info
	return patchImage(ctx, ns, name, containerName, newImage)
}

func (e *PodContainerImageExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	ns, name, containerName, err := model.ParseContainerInfo(injectObject)
	if err != nil {
		return fmt.Errorf("unexpected pod format: %s", err.Error())
	}

	_, err = patchImage(ctx, ns, name, containerName, backup)
	return err
}

func (e *PodContainerImageExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}

func patchImage(ctx context.Context, ns, name, containerName, newImage string) (string, error) {
	var c, pod = restclient.GetApiServerClientMap(v1alpha1.PodCloudTarget), &corev1.Pod{}
	var oldImage string
	var index int
	// get container info
	if err := c.Get().Namespace(ns).Resource("pods").Name(name).Do(ctx).Into(pod); err != nil {
		return "", fmt.Errorf("get pod error: %s", err.Error())
	}

	for i, unitC := range pod.Spec.Containers {
		if unitC.Name == containerName {
			oldImage, index = unitC.Image, i
			break
		}
	}

	if oldImage == "" {
		return "", fmt.Errorf("not found container[%s] in pod[%s/%s]", containerName, ns, name)
	}

	// execute
	patchReader := bytes.NewReader([]byte(fmt.Sprintf(`{"spec":{"containers":[{"name":"%s","image":"%s"}]}}`, pod.Spec.Containers[index].Name, newImage)))
	if err := c.Patch(types.StrategicMergePatchType).Namespace(ns).Resource("pods").Name(name).Body(patchReader).Do(ctx).Error(); err != nil {
		return "", fmt.Errorf("patch error: %s", err.Error())
	}

	return oldImage, nil
}
