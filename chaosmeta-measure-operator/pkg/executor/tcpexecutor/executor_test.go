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

package tcpexecutor

import (
	"context"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"testing"
)

func TestTCPExecutor_Measure(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		args        []v1alpha1.MeasureArgs
		judgement   v1alpha1.Judgement
		initialData string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				ctx: ctx,
				args: []v1alpha1.MeasureArgs{
					{Key: IPArgsKey, Value: "127.0.0.1"},
					{Key: PortArgsKey, Value: "8081"},
					{Key: TimeoutArgsKey, Value: "5"},
				},
				judgement: v1alpha1.Judgement{
					JudgeType:  v1alpha1.ConnectivityJudgeType,
					JudgeValue: v1alpha1.ConnectivityTrue,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TCPExecutor{}
			if err := e.Measure(tt.args.ctx, tt.args.args, tt.args.judgement, tt.args.initialData); (err != nil) != tt.wantErr {
				t.Errorf("Measure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
