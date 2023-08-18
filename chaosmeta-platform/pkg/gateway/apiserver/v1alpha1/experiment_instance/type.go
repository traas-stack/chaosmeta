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

package experiment_instance

import (
	experimentInstanceModel "chaosmeta-platform/pkg/models/experiment_instance"
	"chaosmeta-platform/pkg/service/experiment_instance"
)

type GetExperimentInstancesResponse struct {
	WorkflowNodes []experiment_instance.WorkflowNodesInfo `json:"workflow_nodes"`
}

type GetExperimentInstanceResponse struct {
	WorkflowNode experiment_instance.WorkflowNodesDetail `json:"workflow_node"`
}

type DeleteExperimentInstanceRequest struct {
	ResultUUIDs []string `json:"result_uuids"`
}

type ExperimentInstanceListResponse struct {
	Page        int                                           `json:"page"`
	PageSize    int                                           `json:"pageSize"`
	Total       int64                                         `json:"total"`
	Experiments []*experimentInstanceModel.ExperimentInstance `json:"results"`
}

type GetFaultRangeInstanceResponse struct {
	FaultRangeInstance experimentInstanceModel.FaultRangeInstance `json:"subtask"`
}
