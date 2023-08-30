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

package namespace

import (
	"chaosmeta-platform/pkg/service/namespace"
	"context"
)

func (c *NamespaceController) GetOverview() {
	recentDays, err := c.GetInt("recent_day", 7)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	namespaceId, _ := c.GetInt(":id", 0)

	namespace := &namespace.NamespaceService{}
	totalExperimentCount, totalExperimentInstancesCount, failedExperimentInstancesCount, err := namespace.GetOverview(context.Background(), namespaceId, recentDays)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	response := GetOverviewResponse{TotalExperiments: totalExperimentCount, TotalExperimentInstances: totalExperimentInstancesCount, FailedExperimentInstances: failedExperimentInstancesCount}
	c.Success(&c.Controller, response)
}
