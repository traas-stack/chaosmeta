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
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"testing"
	"time"
)

func TestClusterPendingPodExecutor_Inject(t *testing.T) {
	ctx := context.Background()
	e := GetCloudNativeExecutor(v1alpha1.ClusterCloudTarget, faultClusterPendingPod)
	count, ns, podname := "16", "chaosmeta-pending", "pending-test"
	var args = []v1alpha1.ArgsUnit{
		{Key: "count", Value: count},
		{Key: "namespace", Value: ns},
		{Key: "name", Value: podname},
	}

	gomonkey.ApplyFunc(createNs, func(ctx context.Context, name string) error {
		return nil
	})
	gomonkey.ApplyFunc(createPendingPod, func(ctx context.Context, namespace, name string) error {
		return nil
	})

	backup, err := e.Inject(ctx, "", "", "", args)
	time.Sleep(time.Second)
	assert.Equal(t, backup, ns)
	assert.Equal(t, err, nil)
	assert.Equal(t, common.GetClusterCtrl().IsRunning(), false)

	gomonkey.ApplyFunc(createPendingPod, func(ctx context.Context, namespace, name string) error {
		time.Sleep(10 * time.Second)
		return nil
	})

	gomonkey.ApplyFunc(deleteNs, func(ctx context.Context, name string) error {
		return nil
	})

	backup, err = e.Inject(ctx, "", "", "", args)
	time.Sleep(time.Second)
	assert.Equal(t, backup, ns)
	assert.Equal(t, err, nil)
	assert.Equal(t, common.GetClusterCtrl().IsRunning(), true)
	// running status
	re, err := e.Query(ctx, "", "", ns, v1alpha1.InjectPhaseType)
	assert.Equal(t, err, nil)
	assert.Equal(t, v1alpha1.RunningStatusType, re.Status)
	// has running, retry later
	backup, err = e.Inject(ctx, "", "", "", args)
	assert.NotEqual(t, err, nil)
	// execute recover
	err = e.Recover(ctx, "", "", ns)
	assert.Equal(t, err, nil)
	time.Sleep(3 * time.Second)
	re, err = e.Query(ctx, "", "", ns, v1alpha1.InjectPhaseType)
	assert.Equal(t, err, nil)
	assert.Equal(t, v1alpha1.SuccessStatusType, re.Status)
	assert.Equal(t, common.GetClusterCtrl().IsRunning(), false)
}
