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
	"sync/atomic"
	"time"
)

type GoroutinePool struct {
	n    int
	pool chan struct{}
}

var (
	workerPool  *GoroutinePool
	clusterCtrl = &ClusterCtrl{}
)

func SetGoroutinePool(n int) {
	workerPool = &GoroutinePool{
		n:    n,
		pool: make(chan struct{}, n),
	}
}

func GetGoroutinePool() *GoroutinePool {
	return workerPool
}

func (g *GoroutinePool) GetSize() int {
	return g.n
}

func (g *GoroutinePool) GetLen() int {
	return len(g.pool)
}

func (g *GoroutinePool) GetGoroutine() {
	g.pool <- struct{}{}
}

func (g *GoroutinePool) ReleaseGoroutine() {
	<-g.pool
}

type ClusterCtrl struct {
	runningWorker int64
	stopping      bool
}

func GetClusterCtrl() *ClusterCtrl {
	return clusterCtrl
}

func (c *ClusterCtrl) Run(n int64) bool {
	return !c.stopping && atomic.CompareAndSwapInt64(&c.runningWorker, 0, n)
}

func (c *ClusterCtrl) FinishOne() {
	atomic.AddInt64(&c.runningWorker, -1)
}

func (c *ClusterCtrl) GetRunningWorker() int64 {
	return c.runningWorker
}

func (c *ClusterCtrl) IsRunning() bool {
	return c.runningWorker > 0
}

func (c *ClusterCtrl) IsStopping() bool {
	return c.stopping
}

func (c *ClusterCtrl) Stop() {
	c.stopping = true
	for {
		if !c.IsRunning() {
			break
		}

		time.Sleep(time.Second)
	}
	c.stopping = false
}
