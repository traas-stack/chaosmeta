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

package experiment

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CommonMeasureSpec struct {
	MeasureType  string        `json:"measureType"`
	Duration     string        `json:"duration"`
	Interval     string        `json:"interval"`
	SuccessCount int           `json:"successCount,omitempty"`
	FailedCount  int           `json:"failedCount,omitempty"`
	Stopped      bool          `json:"stopped"`
	Judgement    Judgement     `json:"judgement"`
	Args         []MeasureArgs `json:"args"`
}

type MeasureArgs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Judgement struct {
	JudgeType  string `json:"judgeType"`
	JudgeValue string `json:"judgeValue"`
}

type CommonMeasureStatus struct {
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

type MeasureTask struct {
	Uid        string     `json:"uid"`
	CreateTime string     `json:"createTime"`
	UpdateTime string     `json:"updateTime"`
	Status     StatusType `json:"status"`
	Message    string     `json:"message"`
}

type CommonMeasureStruct struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonMeasureSpec   `json:"spec,omitempty"`
	Status CommonMeasureStatus `json:"status,omitempty"`
}

type CommonMeasureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonMeasureStruct `json:"items"`
}
