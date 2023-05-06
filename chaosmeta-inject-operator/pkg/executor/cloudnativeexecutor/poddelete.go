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
	registerCloudExecutor(v1alpha1.PodCloudTarget, "delete", &PodDeleteExecutor{})
}

type PodDeleteExecutor struct{}

func (e *PodDeleteExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	ns, name, _, err := model.ParsePodInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected pod format: %s", err.Error())
	}

	return "", restclient.GetApiServerClientMap(v1alpha1.PodCloudTarget).
		Delete().Namespace(ns).Resource("pods").Name(name).Do(ctx).Error()
}

func (e *PodDeleteExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	return nil
}

func (e *PodDeleteExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
