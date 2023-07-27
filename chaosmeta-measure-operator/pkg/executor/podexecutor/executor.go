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

package podexecutor

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	NamespaceArgsKey  = "namespace"
	LabelArgsKey      = "label"
	NameprefixArgsKey = "nameprefix"
)

func init() {
	e, err := NewPodExecutor(context.Background())
	if err != nil {
		fmt.Printf("new pod executor error: %s\n", err.Error())
	} else {
		v1alpha1.SetMeasureExecutor(context.Background(), v1alpha1.PodMeasureType, e)
	}
}

type PodExecutor struct {
}

func NewPodExecutor(ctx context.Context) (*PodExecutor, error) {
	return &PodExecutor{}, nil
}

func (e *PodExecutor) CheckConfig(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement) error {
	_, err := utils.GetArgsValueStr(args, NamespaceArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	labels, err := utils.GetArgsValueStr(args, LabelArgsKey)
	if err == nil {
		if _, err := utils.ParseKV(labels); err != nil {
			return fmt.Errorf("label args error: %s", err.Error())
		}
	}

	if judgement.JudgeType != v1alpha1.CountJudgeType {
		return fmt.Errorf("judge type of pod measure only support: %s", v1alpha1.CountJudgeType)
	}

	_, _, err = utils.GetIntervalValue(judgement.JudgeValue)
	if err != nil {
		return fmt.Errorf("get JudgeValue error: %s", err.Error())
	}

	return nil
}

func (e *PodExecutor) InitialData(ctx context.Context, args []v1alpha1.MeasureArgs) (string, error) {
	return "", nil
}

func (e *PodExecutor) Measure(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement, initialData string) error {
	actualCount, err := e.getPodCount(ctx, args)
	if err != nil {
		return fmt.Errorf("get pod count error: %s", err.Error())
	}

	if judgement.JudgeType != v1alpha1.CountJudgeType {
		return fmt.Errorf("judge type of pod measure only support: %s", v1alpha1.CountJudgeType)
	}

	left, right, _ := utils.GetIntervalValue(judgement.JudgeValue)
	return utils.IfMeetInterval(float64(actualCount), left, right)
}

func (e *PodExecutor) getPodCount(ctx context.Context, args []v1alpha1.MeasureArgs) (count int, err error) {
	ns, _ := utils.GetArgsValueStr(args, NamespaceArgsKey)
	opts := []client.ListOption{
		client.InNamespace(ns),
	}

	labelStr, _ := utils.GetArgsValueStr(args, LabelArgsKey)
	if labelStr != "" {
		labels, _ := utils.ParseKV(labelStr)
		opts = append(opts, client.MatchingLabels(labels))
	}

	podList := &corev1.PodList{}
	if err = v1alpha1.GetApiServer().List(ctx, podList, opts...); err != nil {
		return 0, fmt.Errorf("list pod info by label error: %s", err.Error())
	}

	prefix, _ := utils.GetArgsValueStr(args, NameprefixArgsKey)
	if prefix == "" {
		return len(podList.Items), nil
	}

	for _, unitPod := range podList.Items {
		if !strings.HasPrefix(unitPod.Name, prefix) {
			continue
		}
		count++
	}

	return
}
