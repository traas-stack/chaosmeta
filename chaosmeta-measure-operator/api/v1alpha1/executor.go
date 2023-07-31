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

package v1alpha1

import (
	"context"
)

// +kubebuilder:object:generate=false
type MeasureExecutor interface {
	CheckConfig(ctx context.Context, args []MeasureArgs, judgement Judgement) error
	InitialData(ctx context.Context, args []MeasureArgs) (string, error)
	Measure(ctx context.Context, args []MeasureArgs, judgement Judgement, initialData string) error
}

var executorMap = make(map[MeasureType]MeasureExecutor)

func GetMeasureExecutor(ctx context.Context, measureType MeasureType) MeasureExecutor {
	return executorMap[measureType]
}

func SetMeasureExecutor(ctx context.Context, measureType MeasureType, e MeasureExecutor) {
	executorMap[measureType] = e
}
