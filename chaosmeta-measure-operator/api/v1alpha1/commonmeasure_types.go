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
	"math"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	FinalizerName   = "chaosmeta/measure"
	TimeFormat      = "2006-01-02 15:04:05"
	JudgeValueSplit = ","
	KVListSplit     = ","
	KVSplit         = ":"

	IntervalMax = math.MaxFloat64
	IntervalMin = -math.MaxFloat64

	ConnectivityTrue  = "true"
	ConnectivityFalse = "false"
)

var (
	apiServer client.Client
)

func SetApiServer(c client.Client) {
	apiServer = c
}

func GetApiServer() client.Client {
	return apiServer
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CommonMeasureSpec defines the desired state of CommonMeasure
type CommonMeasureSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	MeasureType  MeasureType   `json:"measureType"`
	Duration     string        `json:"duration"`
	Interval     string        `json:"interval"`
	SuccessCount int           `json:"successCount,omitempty"`
	FailedCount  int           `json:"failedCount,omitempty"`
	Stopped      bool          `json:"stopped,omitempty"`
	Judgement    Judgement     `json:"judgement"`
	Args         []MeasureArgs `json:"args"`
}

type MeasureArgs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Judgement struct {
	JudgeType  JudgeType `json:"judgeType"`
	JudgeValue string    `json:"judgeValue"`
}

type MeasureType string

const (
	MonitorMeasureType MeasureType = "monitor"
	PodMeasureType     MeasureType = "pod"
	DeployMeasureType  MeasureType = "deploy"
	SvcMeasureType     MeasureType = "svc"
	HTTPMeasureType    MeasureType = "http"
	IPMeasureType      MeasureType = "ip"
	TCPMeasureType     MeasureType = "tcp"
	//UDPMeasureType     MeasureType = "udp"
)

type JudgeType string

const (
	AbsoluteValueJudgeType   JudgeType = "absolutevalue"
	RelativeValueJudgeType   JudgeType = "relativevalue"
	RelativePercentJudgeType JudgeType = "relativepercent"

	CountJudgeType JudgeType = "count"

	ConnectivityJudgeType JudgeType = "connectivity"
	CodeJudgeType         JudgeType = "code"
	BodyJudgeType         JudgeType = "body"
)

// CommonMeasureStatus defines the observed state of CommonMeasure
type CommonMeasureStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status         StatusType    `json:"status"`
	Message        string        `json:"message"`
	TotalMeasure   int           `json:"totalMeasure"`
	SuccessMeasure int           `json:"successMeasure"`
	FailedMeasure  int           `json:"failedMeasure"`
	InitialData    string        `json:"initialData"`
	CreateTime     string        `json:"createTime"`
	UpdateTime     string        `json:"updateTime"`
	NextTime       string        `json:"nextTime"`
	Measures       []MeasureTask `json:"measures,omitempty"`
}

type StatusType string

const (
	CreatedStatus StatusType = "created"
	RunningStatus StatusType = "running"
	SuccessStatus StatusType = "success"
	FailedStatus  StatusType = "failed"
)

type MeasureTask struct {
	Uid        string     `json:"uid"`
	CreateTime string     `json:"createTime"`
	UpdateTime string     `json:"updateTime"`
	Status     StatusType `json:"status"`
	Message    string     `json:"message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CommonMeasure is the Schema for the commonmeasures API
type CommonMeasure struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonMeasureSpec   `json:"spec,omitempty"`
	Status CommonMeasureStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CommonMeasureList contains a list of CommonMeasure
type CommonMeasureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonMeasure `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommonMeasure{}, &CommonMeasureList{})
}
