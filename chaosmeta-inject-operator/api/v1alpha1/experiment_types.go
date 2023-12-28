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

const (
	ArgsListSplit  = ","
	LabelListSplit = "="
	FinalizerName  = "chaosmeta/experiment"
	ContainerKey   = "containername"
	FirstContainer = "firstcontainer"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ScopeType string

const (
	PodScopeType        ScopeType = "pod"
	NodeScopeType       ScopeType = "node"
	KubernetesScopeType ScopeType = "kubernetes"
)

// ExperimentSpec defines the desired state of Experiment
type ExperimentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Scope Optional: node, pod. type of experiment object
	Scope      ScopeType         `json:"scope"`
	RangeMode  *RangeMode        `json:"rangeMode,omitempty"`
	Experiment *ExperimentCommon `json:"experiment"`
	// Selector The internal part of unit is "AND", and the external part is "OR" and de-duplication
	Selector []SelectorUnit `json:"selector,omitempty"`

	TargetPhase PhaseType `json:"targetPhase"`
	//SubObj      bool      `json:"subObj"`
}

type PhaseType string

const (
	InjectPhaseType  PhaseType = "inject"
	RecoverPhaseType PhaseType = "recover"
)

type StatusType string

const (
	CreatedStatusType     StatusType = "created"
	SuccessStatusType     StatusType = "success"
	FailedStatusType      StatusType = "failed"
	RunningStatusType     StatusType = "running"
	PartSuccessStatusType StatusType = "partSuccess"
)

// ExperimentStatus defines the observed state of Experiment
type ExperimentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase      PhaseType        `json:"phase"`
	Status     StatusType       `json:"status"`
	Message    string           `json:"message"`
	Detail     ExperimentDetail `json:"detail"`
	CreateTime string           `json:"createTime"`
	UpdateTime string           `json:"updateTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Experiment is the Schema for the experiments API
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExperimentSpec   `json:"spec,omitempty"`
	Status ExperimentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExperimentList contains a list of Experiment
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}

type RangeType string

const (
	AllRangeType     RangeType = "all"
	PercentRangeType RangeType = "percent"
	CountRangeType   RangeType = "count"
)

type RangeMode struct {
	// Type Optional: all、percent、count
	Type  RangeType `json:"type"`
	Value int       `json:"value,omitempty"`
}

type SelectorUnit struct {
	Namespace string            `json:"namespace,omitempty"`
	Name      []string          `json:"name,omitempty"`
	IP        []string          `json:"ip,omitempty"`
	Label     map[string]string `json:"label,omitempty"`
	SubName   string            `json:"subName,omitempty"`
}

//type TargetType string
//type FaultType string

type ExperimentCommon struct {
	// Duration support "h", "m", "s"
	Duration string     `json:"duration,omitempty"`
	Target   string     `json:"target"`
	Fault    string     `json:"fault"`
	Args     []ArgsUnit `json:"args,omitempty"`
}

type VType string

const (
	IntVType    VType = "int"
	StringVType VType = "string"
)

type ArgsUnit struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType VType  `json:"valueType,omitempty"`
}

type ExperimentDetail struct {
	Inject  []ExperimentDetailUnit `json:"inject,omitempty"`
	Recover []ExperimentDetailUnit `json:"recover,omitempty"`
}

type ExperimentDetailUnit struct {
	InjectObjectName string `json:"injectObjectName,omitempty"`
	//InjectObjectInfo string     `json:"injectObjectInfo,omitempty"`
	UID        string     `json:"uid,omitempty"`
	Message    string     `json:"message,omitempty"`
	Status     StatusType `json:"status,omitempty"`
	StartTime  string     `json:"startTime,omitempty"`
	UpdateTime string     `json:"updateTime,omitempty"`
	Backup     string     `json:"backup,omitempty"`
}

type CloudTargetType string

const (
	ClusterCloudTarget     CloudTargetType = "cluster"
	PodCloudTarget         CloudTargetType = "pod"
	NodeCloudTarget        CloudTargetType = "node"
	DeploymentCloudTarget  CloudTargetType = "deployment"
	StatefulsetCloudTarget CloudTargetType = "statefulset"
	DaemonsetCloudTarget   CloudTargetType = "daemonset"
	NamespaceCloudTarget   CloudTargetType = "namespace"
	JobCloudTarget         CloudTargetType = "job"
)
