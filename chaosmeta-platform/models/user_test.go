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

package models

import (
	"chaosmeta-platform/logger"
	"context"
	"testing"
)

func TestUpdateUserRole(t *testing.T) {
	Setup()
	ctx := context.WithValue(context.Background(), logger.TraceIdKey, "erg3g42g432g")

	type args struct {
		id   int
		role string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				id:   1,
				role: AdminRole,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserRole(ctx, tt.args.id, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
