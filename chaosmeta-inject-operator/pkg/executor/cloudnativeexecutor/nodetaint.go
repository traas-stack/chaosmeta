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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.NodeCloudTarget, "taint", &NodeTaintExecutor{})
}

type NodeTaintExecutor struct{}

func (e *NodeTaintExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	name, _, err := model.ParseNodeInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected node format: %s", err.Error())
	}

	c, node := restclient.GetApiServerClientMap(v1alpha1.NodeCloudTarget), &corev1.Node{}
	if err := c.Get().Resource("nodes").Name(name).Do(ctx).Into(node); err != nil {
		return "", fmt.Errorf("get node error: %s", err.Error())
	}

	backupBytes, err := json.Marshal(node.Spec.Taints)
	if err != nil {
		return "", fmt.Errorf("backup to string error: %s", err.Error())
	}

	return string(backupBytes), patchTaints(ctx, c, name, getNewTaints(ctx, node.Spec.Taints, args))
}

func patchTaints(ctx context.Context, c rest.Interface, name string, taints []corev1.Taint) error {
	payload, err := json.Marshal(taints)
	if err != nil {
		return fmt.Errorf("get payload error: %s", err.Error())
	}

	if err := c.Patch(types.MergePatchType).Resource("nodes").Name(name).
		Body([]byte(fmt.Sprintf(`{"spec":{"taints":%s}}`, payload))).Do(ctx).Error(); err != nil {
		return fmt.Errorf("patch taints error: %s", err.Error())
	}

	return nil
}

func getNewTaints(ctx context.Context, oldTaints []corev1.Taint, args []v1alpha1.ArgsUnit) []corev1.Taint {
	logger := log.FromContext(ctx)
	var oldTainsArr []string
	for _, taint := range oldTaints {
		oldTainsArr = append(oldTainsArr, fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect))
	}

	newTainsArr := getNewFinalizers(ctx, oldTainsArr, args)
	var newTains []corev1.Taint
	for _, unit := range newTainsArr {
		kvArr := strings.Split(unit, "=")
		if len(kvArr) != 2 {
			logger.Error(fmt.Errorf("not found \"=\""), "parse taints error")
			continue
		}
		valueArr := strings.Split(kvArr[1], ":")
		if len(valueArr) != 2 {
			logger.Error(fmt.Errorf("not found \":\""), "parse taints error")
			continue
		}
		newTains = append(newTains, corev1.Taint{
			Key:    kvArr[0],
			Value:  valueArr[0],
			Effect: corev1.TaintEffect(valueArr[1]),
		})
	}

	return newTains
}

func (e *NodeTaintExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	name, _, err := model.ParseNodeInfo(injectObject)
	if err != nil {
		return fmt.Errorf("unexpected node format: %s", err.Error())
	}

	var oldTaints []corev1.Taint
	if backup != "" {
		if err := json.Unmarshal([]byte(backup), &oldTaints); err != nil {
			return fmt.Errorf("get old taints error: %s", err.Error())
		}
	}

	c := restclient.GetApiServerClientMap(v1alpha1.NodeCloudTarget)
	return patchTaints(ctx, c, name, oldTaints)
}

func (e *NodeTaintExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}
