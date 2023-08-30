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
	"chaosmeta-platform/pkg/models/experiment"
	"chaosmeta-platform/pkg/models/experiment_instance"
	"context"
)

func (s *NamespaceService) GetOverview(ctx context.Context, namespaceID int, recentDays int) (int64, int64, int64, error) {
	totalExperimentCount, err := experiment.CountExperiments(namespaceID, -1, recentDays)
	if err != nil {
		return totalExperimentCount, 0, 0, err
	}

	totalExperimentInstancesCount, err := experiment_instance.CountExperimentInstances(namespaceID, "", "", recentDays)
	if err != nil {
		return totalExperimentCount, totalExperimentInstancesCount, 0, err
	}
	failedExperimentInstancesCount, err := experiment_instance.CountExperimentInstances(namespaceID, "", "Failed", recentDays)
	return totalExperimentCount, totalExperimentInstancesCount, failedExperimentInstancesCount, err
}
