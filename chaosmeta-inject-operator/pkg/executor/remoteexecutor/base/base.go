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

package base

import "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"

type RemoteExpStatus string

const (
	CreatedStatus   RemoteExpStatus = "created"
	SuccessStatus   RemoteExpStatus = "success"
	ErrorStatus     RemoteExpStatus = "error"
	DestroyedStatus RemoteExpStatus = "destroyed"
)

func ConvertStatus(status RemoteExpStatus, phase v1alpha1.PhaseType) v1alpha1.StatusType {
	switch phase {
	case v1alpha1.InjectPhaseType:
		switch status {
		case CreatedStatus:
			return v1alpha1.RunningStatusType
		case SuccessStatus:
			return v1alpha1.SuccessStatusType
		case ErrorStatus:
			return v1alpha1.FailedStatusType
		case DestroyedStatus:
			return v1alpha1.SuccessStatusType
		}
	case v1alpha1.RecoverPhaseType:
		switch status {
		case ErrorStatus:
			return v1alpha1.SuccessStatusType
		case DestroyedStatus:
			return v1alpha1.SuccessStatusType
		}
	}

	// unexpected status
	return v1alpha1.FailedStatusType
}

const (
	SucCode = 0
	//TaskNotFoundCode      = 1
	//ContainerNotFoundCode = 2
)

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"trace_id,omitempty"`
}

type ExperimentDataUnit struct {
	Uid              string          `json:"uid"`
	Target           string          `json:"target"`
	Fault            string          `json:"fault"`
	Args             string          `json:"args"`
	Runtime          string          `json:"runtime"`
	Status           RemoteExpStatus `json:"status"`
	Creator          string          `json:"creator"`
	Timeout          string          `json:"timeout,omitempty"`
	Error_           string          `json:"error,omitempty"`
	CreateTime       string          `json:"create_time,omitempty"`
	UpdateTime       string          `json:"update_time,omitempty"`
	ContainerId      string          `json:"container_id,omitempty"`
	ContainerRuntime string          `json:"container_runtime,omitempty"`
}

type InjectRequest struct {
	Target           string `json:"target"`
	Fault            string `json:"fault"`
	Timeout          string `json:"timeout"`
	Creator          string `json:"creator"`
	Args             string `json:"args"`
	ContainerId      string `json:"container_id"`
	ContainerRuntime string `json:"container_runtime"`
	TraceId          string `json:"trace_id"`
	Uid              string `json:"uid"`
}

type InjectResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    *InjectSuccessResponseData `json:"data,omitempty"`
	TraceId string                     `json:"trace_id,omitempty"`
}

type InjectSuccessResponseData struct {
	Experiment ExperimentDataUnit `json:"experiment,omitempty"`
}

type QueryRequest struct {
	Uid              string          `json:"uid,omitempty"`
	Status           RemoteExpStatus `json:"status,omitempty"`
	Target           string          `json:"target,omitempty"`
	Fault            string          `json:"fault,omitempty"`
	Creator          string          `json:"creator,omitempty"`
	ContainerId      string          `json:"container_id,omitempty"`
	ContainerRuntime string          `json:"container_runtime,omitempty"`
	Offset           int32           `json:"offset,omitempty"`
	Limit            int32           `json:"limit,omitempty"`
	TraceId          string          `json:"trace_id,omitempty"`
}

type QueryResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *QueryResponseData `json:"data,omitempty"`
	TraceId string             `json:"trace_id,omitempty"`
}

type QueryResponseData struct {
	Total       int64                `json:"total"`
	Experiments []ExperimentDataUnit `json:"experiments,omitempty"`
}

type RecoverRequest struct {
	Uid     string `json:"uid"`
	TraceId string `json:"trace_id"`
}

//type TaskResult struct {
//	Err     error
//	Status  v1alpha1.StatusType
//	Message string
//}

type VersionResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *VersionInfo `json:"data,omitempty"`
}

type VersionInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"build-date"`
}
