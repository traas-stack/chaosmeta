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
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
	"time"
)

func init() {
	registerCloudExecutor(v1alpha1.DeploymentCloudTarget, "replicas", &DeploymentReplicasExecutor{})
}

type DeploymentReplicasExecutor struct{}

func (e *DeploymentReplicasExecutor) Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error) {
	ns, name, err := model.ParseDeploymentInfo(injectObject)
	if err != nil {
		return "", fmt.Errorf("unexpected deployment format: %s", err.Error())
	}

	replicasArgs, err := ParseReplicasArgs(args)
	if err != nil {
		return "", fmt.Errorf("args error: %s", err.Error())
	}

	c := restclient.GetApiServerClientMap(v1alpha1.DeploymentCloudTarget)
	deploy := &v1.Deployment{}
	if err := c.Get().Namespace(ns).Resource("deployments").Name(name).Do(ctx).Into(deploy); err != nil {
		return "", fmt.Errorf("get deployment error: %s", err.Error())
	}

	oldCount := int(*deploy.Spec.Replicas)
	var count = replicasArgs.getAbsoluteCount(oldCount)
	if oldCount == count {
		return strconv.Itoa(oldCount), nil
	}

	if err := c.Patch(types.MergePatchType).Namespace(ns).Resource("deployments").Name(name).
		Body([]byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, count))).SubResource("scale").Do(ctx).Error(); err != nil {
		return "", fmt.Errorf("patch deployment error: %s", err.Error())
	}

	return strconv.Itoa(oldCount), nil
}

func (e *DeploymentReplicasExecutor) Recover(ctx context.Context, injectObject, uid, backup string) error {
	ns, name, err := model.ParseDeploymentInfo(injectObject)
	if err != nil {
		return fmt.Errorf("unexpected deployment format: %s", err.Error())
	}

	oldCount, err := strconv.Atoi(backup)
	if err != nil {
		return fmt.Errorf("old replicas is not a num: %s", err.Error())
	}

	c := restclient.GetApiServerClientMap(v1alpha1.DeploymentCloudTarget)
	if err := c.Patch(types.MergePatchType).Namespace(ns).Resource("deployments").Name(name).
		Body([]byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, oldCount))).SubResource("scale").Do(ctx).Error(); err != nil {
		return fmt.Errorf("patch deployment error: %s", err.Error())
	}

	return nil
}

func (e *DeploymentReplicasExecutor) Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	return &model.SubExpInfo{
		UID:        uid,
		Status:     v1alpha1.SuccessStatusType,
		UpdateTime: time.Now().Format(model.TimeFormat),
	}, nil
}

type ReplicasArgs struct {
	Mode  ReplicasModeType `json:"mode,omitempty"`
	Value int              `json:"value"`
}

type ReplicasModeType string

const (
	AbsoluteCountMode   ReplicasModeType = "absolutecount"
	RelativeCountMode   ReplicasModeType = "relativecount"
	RelativePercentMode ReplicasModeType = "relativepercent"
)

func (a *ReplicasArgs) getAbsoluteCount(oldCount int) int {
	if a.Mode == RelativeCountMode {
		return a.Value + oldCount
	}

	if a.Mode == RelativePercentMode {
		return oldCount + a.Value*oldCount/100
	}

	return a.Value
}

func ParseReplicasArgs(args []v1alpha1.ArgsUnit) (*ReplicasArgs, error) {
	argsMap := make(map[string]string)
	for _, unitArgs := range args {
		argsMap[unitArgs.Key] = unitArgs.Value
	}

	value, err := strconv.Atoi(argsMap["value"])
	if err != nil {
		return nil, fmt.Errorf("value is not a num: %s", err.Error())
	}

	re := &ReplicasArgs{
		Mode:  ReplicasModeType(argsMap["mode"]),
		Value: value,
	}

	if re.Mode != AbsoluteCountMode && re.Mode != RelativeCountMode && re.Mode != RelativePercentMode {
		return nil, fmt.Errorf("not support mode: %s, only support: %s, %s, %s",
			re.Mode, AbsoluteCountMode, RelativeCountMode, RelativePercentMode)
	}

	return re, nil
}
