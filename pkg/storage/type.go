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

package storage

type Experiment struct {
	Uid              string `gorm:"primary_key" json:"uid"`
	Target           string `gorm:"index:target" json:"target"`
	Fault            string `gorm:"index:fault" json:"fault"`
	Args             string `json:"args"`
	Runtime          string `json:"runtime"`
	Timeout          string `json:"timeout"`
	Status           string `gorm:"index:status" json:"status"`
	Creator          string `gorm:"index:creator" json:"creator"`
	Error            string `json:"error"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
	ContainerId      string `json:"container_id"`
	ContainerRuntime string `json:"container_runtime"`
}
