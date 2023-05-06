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
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	v1 "k8s.io/api/apps/v1"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.DeploymentCloudTarget, "finalizer", &DeploymentFinalizerExecutor{})
}

type DeploymentFinalizerExecutor struct{}

func (e *DeploymentFinalizerExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	ns, name, err := model.ParseDeploymentInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected deployment format: %s", err.Error())
	}

	c, deploy := restclient.GetApiServerClientMap(v1alpha1.DeploymentCloudTarget), &v1.Deployment{}
	if err := c.Get().Namespace(ns).Resource("deployments").Name(name).Do(ctx).Into(deploy); err != nil {
		return "", fmt.Errorf("get deployment error: %s", err.Error())
	}

	var backupBytes []byte
	if deploy.ObjectMeta.Finalizers != nil {
		backupBytes, err = json.Marshal(deploy.ObjectMeta.Finalizers)
		if err != nil {
			return "", fmt.Errorf("backup to string error: %s", err.Error())
		}
	}

	return string(backupBytes), patchFinalizers(ctx, c, "deployments", ns, name, getNewFinalizers(ctx, deploy.ObjectMeta.Finalizers, args))
}

func (e *DeploymentFinalizerExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	ns, name, err := model.ParseDeploymentInfo(injectObject)
	if err != nil {
		return fmt.Errorf("unexpected deployment format: %s", err.Error())
	}

	var oldFinalizers []string
	if backup != "" {
		if err := json.Unmarshal([]byte(backup), &oldFinalizers); err != nil {
			return fmt.Errorf("get old finalizers error: %s", err.Error())
		}
	}

	c := restclient.GetApiServerClientMap(v1alpha1.DeploymentCloudTarget)
	return patchFinalizers(ctx, c, "deployments", ns, name, oldFinalizers)
}

func (e *DeploymentFinalizerExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
