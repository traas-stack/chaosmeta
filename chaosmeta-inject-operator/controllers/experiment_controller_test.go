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

package controllers

import (
	"context"
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	mockscopehandler "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/mock/scopehandler"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_solveRange(t *testing.T) {
	testCount := 5
	var testObjectList []model.AtomicObject
	for i := 0; i < testCount; i++ {
		testObjectList = append(testObjectList, &model.PodObject{
			Namespace: "ns2",
			PodName:   fmt.Sprintf("pod%d", i),
		})
		testObjectList = append(testObjectList, &model.PodObject{
			Namespace: "ns1",
			PodName:   fmt.Sprintf("pod%d", i),
		})
	}

	type args struct {
		initial   []model.AtomicObject
		rangeMode *v1alpha1.RangeMode
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "all success",
			args: args{
				initial: testObjectList,
				rangeMode: &v1alpha1.RangeMode{
					Type: v1alpha1.AllRangeType,
				},
			},
			want: testCount * 2,
		},
		{
			name: "count success",
			args: args{
				initial: testObjectList,
				rangeMode: &v1alpha1.RangeMode{
					Type:  v1alpha1.CountRangeType,
					Value: 3,
				},
			},
			want: 3,
		},
		{
			name: "count more then initial length",
			args: args{
				initial: testObjectList,
				rangeMode: &v1alpha1.RangeMode{
					Type:  v1alpha1.CountRangeType,
					Value: 15,
				},
			},
			want: 10,
		},
		{
			name: "percent success",
			args: args{
				initial: testObjectList,
				rangeMode: &v1alpha1.RangeMode{
					Type:  v1alpha1.PercentRangeType,
					Value: 65,
				},
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solveRange(tt.args.initial, tt.args.rangeMode)
			if len(got) != tt.want {
				t.Errorf("solveRange() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_initProcess(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		ctx  = context.Background()
		exp  = &v1alpha1.Experiment{
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
							Value: "nginx",
						},
					},
				},
				Selector: []v1alpha1.SelectorUnit{
					{
						Namespace: "chaosmeta",
					},
				},
				TargetPhase: v1alpha1.InjectPhaseType,
			},
			Status: v1alpha1.ExperimentStatus{},
		}
	)

	var reObject []model.AtomicObject
	reObject = append(reObject, &model.PodObject{
		Namespace: "chaosmeta",
		PodName:   "chaosmeta-0",
		PodUID:    "d32tg32",
		PodIP:     "1.2.3.4",
		NodeName:  "node-1",
		NodeIP:    "2.2.2.2",
		Containers: []model.ContainerInfo{
			{
				ContainerId:      "g3g3g",
				ContainerRuntime: "docker",
			},
		},
	})
	defer ctrl.Finish()
	scopeHandlerMock := mockscopehandler.NewMockScopeHandler(ctrl)
	scopeHandlerMock.EXPECT().ConvertSelector(ctx, &exp.Spec).Return(reObject, nil)

	gomonkey.ApplyFunc(scopehandler.GetScopeHandler, func(v1alpha1.ScopeType) scopehandler.ScopeHandler {
		return scopeHandlerMock
	})

	initProcess(ctx, exp)
	assert.Equal(t, "pod/chaosmeta/chaosmeta-0", exp.Status.Detail.Inject[0].InjectObjectName)
	assert.Equal(t, v1alpha1.CreatedStatusType, exp.Status.Detail.Inject[0].Status)
	assert.Equal(t, v1alpha1.CreatedStatusType, exp.Status.Status)
	assert.Equal(t, v1alpha1.InjectPhaseType, exp.Status.Phase)

	scopeHandlerMock.EXPECT().ConvertSelector(ctx, &exp.Spec).Return([]model.AtomicObject{}, nil)
	initProcess(ctx, exp)
	assert.Equal(t, v1alpha1.FailedStatusType, exp.Status.Status)
}

func Test_solveFinalizer(t *testing.T) {
	instance := &v1alpha1.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Finalizers: []string{"awbgrewga", v1alpha1.FinalizerName, "shbertbhersth"},
		},
		Status: v1alpha1.ExperimentStatus{
			Phase:  v1alpha1.RecoverPhaseType,
			Status: v1alpha1.SuccessStatusType,
		},
	}
	solveFinalizer(instance)
	assert.Equal(t, []string{"awbgrewga", "shbertbhersth"}, instance.ObjectMeta.Finalizers)

	instance.ObjectMeta.Finalizers = []string{v1alpha1.FinalizerName, "shbertbhersth"}
	solveFinalizer(instance)
	assert.Equal(t, []string{"shbertbhersth"}, instance.ObjectMeta.Finalizers)

	instance.ObjectMeta.Finalizers = []string{"awbgrewga", v1alpha1.FinalizerName}
	solveFinalizer(instance)
	assert.Equal(t, []string{"awbgrewga"}, instance.ObjectMeta.Finalizers)

	instance.ObjectMeta.Finalizers = []string{v1alpha1.FinalizerName}
	solveFinalizer(instance)
	assert.Equal(t, []string{}, instance.ObjectMeta.Finalizers)
}
