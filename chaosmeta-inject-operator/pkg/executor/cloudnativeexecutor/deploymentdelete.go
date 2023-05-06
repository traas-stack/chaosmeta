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
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.DeploymentCloudTarget, "delete", &DeploymentDeleteExecutor{})
}

type DeploymentDeleteExecutor struct{}

func (e *DeploymentDeleteExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	ns, name, err := model.ParseDeploymentInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected deployment format: %s", err.Error())
	}

	return "", restclient.GetApiServerClientMap(v1alpha1.DeploymentCloudTarget).Delete().Namespace(ns).
		Resource("deployments").Name(name).Do(ctx).Error()
}

func (e *DeploymentDeleteExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	return nil
}

func (e *DeploymentDeleteExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
