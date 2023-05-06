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

package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoroutinePool(t *testing.T) {
	SetGoroutinePool(3)

	pool := GetGoroutinePool()
	pool.GetGoroutine()
	pool.GetGoroutine()
	pool.GetGoroutine()
	assert.Equal(t, 3, pool.GetLen())

	var flag bool
	go func(t *bool) {
		pool.GetGoroutine()
		*t = true
	}(&flag)
	time.Sleep(time.Second)
	assert.Equal(t, false, flag)
	pool.ReleaseGoroutine()
	time.Sleep(time.Second)
	assert.Equal(t, true, flag)

	pool.ReleaseGoroutine()
	pool.ReleaseGoroutine()
	assert.Equal(t, 1, pool.GetLen())
	pool.ReleaseGoroutine()
	assert.Equal(t, 0, pool.GetLen())
	assert.Equal(t, 3, pool.GetSize())
}

func TestClusterCtrl(t *testing.T) {
	assert.Equal(t, false, GetClusterCtrl().IsStopping())
	assert.Equal(t, false, GetClusterCtrl().IsRunning())
	suc := GetClusterCtrl().Run(3)
	assert.Equal(t, true, GetClusterCtrl().IsRunning())
	assert.Equal(t, false, GetClusterCtrl().IsStopping())
	assert.Equal(t, true, suc)
	for i := 0; i < 3; i++ {
		go func() {
			defer GetClusterCtrl().FinishOne()

			for {
				if GetClusterCtrl().IsStopping() {
					break
				}
				time.Sleep(time.Second)
			}
		}()
	}
	suc = GetClusterCtrl().Run(1)
	assert.Equal(t, false, suc)
	assert.Equal(t, int64(3), GetClusterCtrl().GetRunningWorker())

	GetClusterCtrl().Stop()
	time.Sleep(time.Second)
	assert.Equal(t, false, GetClusterCtrl().IsStopping())
	assert.Equal(t, false, GetClusterCtrl().IsRunning())
}
