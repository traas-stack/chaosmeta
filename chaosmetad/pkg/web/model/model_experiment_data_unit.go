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

package model

type ExperimentDataUnit struct {
	Uid              string `json:"uid"`
	Target           string `json:"target"`
	Fault            string `json:"fault"`
	Args             string `json:"args"`
	Runtime          string `json:"runtime"`
	Status           string `json:"status"`
	Creator          string `json:"creator"`
	Timeout          string `json:"timeout,omitempty"`
	Error_           string `json:"error,omitempty"`
	CreateTime       string `json:"create_time,omitempty"`
	UpdateTime       string `json:"update_time,omitempty"`
	ContainerId      string `json:"container_id,omitempty"`
	ContainerRuntime string `json:"container_runtime,omitempty"`
}
