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

import (
	"github.com/stretchr/testify/assert"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"testing"
)

func TestConvertStatus(t *testing.T) {
	type args struct {
		status RemoteExpStatus
		phase  v1alpha1.PhaseType
	}
	tests := []struct {
		name string
		args args
		want v1alpha1.StatusType
	}{
		{
			name: "inject_created",
			args: args{
				status: CreatedStatus,
				phase:  v1alpha1.InjectPhaseType,
			},
			want: v1alpha1.RunningStatusType,
		},
		{
			name: "inject_success",
			args: args{
				status: SuccessStatus,
				phase:  v1alpha1.InjectPhaseType,
			},
			want: v1alpha1.SuccessStatusType,
		},
		{
			name: "inject_error",
			args: args{
				status: ErrorStatus,
				phase:  v1alpha1.InjectPhaseType,
			},
			want: v1alpha1.FailedStatusType,
		},
		{
			name: "inject_destroyed",
			args: args{
				status: DestroyedStatus,
				phase:  v1alpha1.InjectPhaseType,
			},
			want: v1alpha1.SuccessStatusType,
		},
		{
			name: "recover_created",
			args: args{
				status: CreatedStatus,
				phase:  v1alpha1.RecoverPhaseType,
			},
			want: v1alpha1.FailedStatusType,
		},
		{
			name: "recover_success",
			args: args{
				status: SuccessStatus,
				phase:  v1alpha1.RecoverPhaseType,
			},
			want: v1alpha1.FailedStatusType,
		},
		{
			name: "recover_error",
			args: args{
				status: ErrorStatus,
				phase:  v1alpha1.RecoverPhaseType,
			},
			want: v1alpha1.SuccessStatusType,
		},
		{
			name: "recover_destroyed",
			args: args{
				status: DestroyedStatus,
				phase:  v1alpha1.RecoverPhaseType,
			},
			want: v1alpha1.SuccessStatusType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ConvertStatus(tt.args.status, tt.args.phase), "ConvertStatus(%v, %v)", tt.args.status, tt.args.phase)
		})
	}
}
