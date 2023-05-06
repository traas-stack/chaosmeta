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

import "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"

type SubExpInfo struct {
	UID        string              `json:"uid"`
	Status     v1alpha1.StatusType `json:"status"`
	CreateTime string              `json:"create_time,omitempty"`
	UpdateTime string              `json:"update_time,omitempty"`
	Message    string              `json:"error,omitempty"`
}

const (
	TimeFormat    = "2006-01-02 15:04:05"
	DockerRuntime = "docker"
)

//
//type SubExpStatus string
//
//const (
//	CreatedSubExpStatus   = "created"
//	SuccessSubExpStatus   = "success"
//	ErrorSubExpStatus     = "error"
//	DestroyedSubExpStatus = "destroyed"
//)
