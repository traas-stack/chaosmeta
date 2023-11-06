/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	TimeFormat    = "2006-01-02 15:04:05"
	FinalizerName = "chaosmeta/flow"

	HostArgsKey   = "host"
	PortArgsKey   = "port"
	MethodArgsKey = "method"
	HeaderArgsKey = "header"
	PathArgsKey   = "path"
	BodyArgsKey   = "body"

	MethodPOST = "POST"
	MethodGET  = "GET"
)

var (
	JmeterConfigStr string
	JobYamlStr      string
)

// LoadTestSpec defines the desired state of LoadTest
type LoadTestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	FlowType    FlowType   `json:"flowType"`
	Duration    string     `json:"duration"`
	Parallelism int        `json:"parallelism"`
	Source      int        `json:"source"`
	Stopped     bool       `json:"stopped"`
	Args        []FlowArgs `json:"args"`
}

type FlowType string

const (
	HTTPFlowType FlowType = "http"
)

type FlowArgs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LoadTestStatus defines the observed state of LoadTest
type LoadTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status       StatusType `json:"status"`
	Message      string     `json:"message"`
	SuccessCount int        `json:"successCount"`
	TotalCount   int        `json:"totalCount"`
	AvgRPS       int        `json:"avgRPS"`
	CreateTime   string     `json:"createTime"`
	UpdateTime   string     `json:"updateTime"`
}

type StatusType string

const (
	CreatedStatus StatusType = "created"
	RunningStatus StatusType = "running"
	SuccessStatus StatusType = "success"
	FailedStatus  StatusType = "failed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LoadTest is the Schema for the loadtests API
type LoadTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadTestSpec   `json:"spec,omitempty"`
	Status LoadTestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LoadTestList contains a list of LoadTest
type LoadTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadTest{}, &LoadTestList{})
}
