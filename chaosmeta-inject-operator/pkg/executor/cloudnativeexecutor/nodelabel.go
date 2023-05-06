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
	corev1 "k8s.io/api/core/v1"
	"time"
)

const (
	faultNodeLabel = "label"
)

func init() {
	registerCloudExecutor(v1alpha1.NodeCloudTarget, faultNodeLabel, &NodeLabelExecutor{})
}

type NodeLabelExecutor struct{}

func (e *NodeLabelExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	name, _, err := model.ParseNodeInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected node format: %s", err.Error())
	}

	c, node := restclient.GetApiServerClientMap(v1alpha1.NodeCloudTarget), &corev1.Node{}
	if err := c.Get().Resource("nodes").Name(name).Do(ctx).Into(node); err != nil {
		return "", fmt.Errorf("get node error: %s", err.Error())
	}

	var backupBytes []byte
	if node.ObjectMeta.Labels != nil {
		backupBytes, err = json.Marshal(node.ObjectMeta.Labels)
		if err != nil {
			return "", fmt.Errorf("backup to string error: %s", err.Error())
		}
	}

	newLabels, err := getNewLabels(ctx, node.ObjectMeta.Labels, args)
	if err != nil {
		return "", fmt.Errorf("get new labels error: %s", err.Error())
	}

	return string(backupBytes), patchLabels(ctx, c, "nodes", "", name, newLabels)
}

func (e *NodeLabelExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	name, _, err := model.ParseNodeInfo(injectObject)
	if err != nil {
		return fmt.Errorf("unexpected node format: %s", err.Error())
	}

	c, node := restclient.GetApiServerClientMap(v1alpha1.NodeCloudTarget), &corev1.Node{}
	if err := c.Get().Resource("nodes").Name(name).Do(ctx).Into(node); err != nil {
		return fmt.Errorf("get node error: %s", err.Error())
	}

	backupBytes, err := getBackupLabels([]byte(backup), node.ObjectMeta.Labels)
	if err != nil {
		return fmt.Errorf("get backup labels error: %s", err.Error())
	}

	return patchLabels(ctx, c, "nodes", "", name, backupBytes)
}

func (e *NodeLabelExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
