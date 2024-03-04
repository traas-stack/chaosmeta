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

package pod

import (
	"context"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	mockselector "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/mock/selector"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
	"testing"
)

func TestPodScopeHandler_ConvertSelector(t *testing.T) {
	var (
		namespace     = "ns"
		containerName = "nginx"
		exp           = &v1alpha1.Experiment{
			Spec: v1alpha1.ExperimentSpec{
				Scope: v1alpha1.PodScopeType,
				RangeMode: &v1alpha1.RangeMode{
					Type:  v1alpha1.CountRangeType,
					Value: 3,
				},
				Experiment: &v1alpha1.ExperimentCommon{
					Duration: "2m",
					Target:   "cpu",
					Fault:    "burn",
					Args: []v1alpha1.ArgsUnit{
						{
							Key:       "percent",
							Value:     "90",
							ValueType: v1alpha1.IntVType,
						},
						{
							Key:   v1alpha1.ContainerKey,
							Value: containerName,
						},
					},
				},
				Selector: []v1alpha1.SelectorUnit{
					{
						Namespace: namespace,
						Label: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
						SubName: containerName,
					},
				},
				TargetPhase: v1alpha1.InjectPhaseType,
			},
		}
		podList = []*model.PodObject{
			{
				Namespace: namespace,
				PodName:   "pod1",
				NodeName:  "node1",
				NodeIP:    "1.1.1.1",
				Containers: []model.ContainerInfo{{
					ContainerId:      "ef2g24g21",
					ContainerRuntime: "docker",
				}},
			},
			{
				Namespace: namespace,
				PodName:   "pod2",
				NodeName:  "node2",
				NodeIP:    "1.1.1.2",
				Containers: []model.ContainerInfo{{
					ContainerId:      "ef2g24g22",
					ContainerRuntime: "docker",
				}},
			},
		}
	)

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	analyzerMock := mockselector.NewMockIAnalyzer(ctrl)
	analyzerMock.EXPECT().GetPodListByLabel(ctx, namespace, exp.Spec.Selector[0].Label, containerName).Return(podList, nil)
	gomonkey.ApplyFunc(selector.GetAnalyzer, func() selector.IAnalyzer {
		return analyzerMock
	})

	reList, err := GetGlobalPodHandler().ConvertSelector(ctx, &exp.Spec)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(reList))
}
