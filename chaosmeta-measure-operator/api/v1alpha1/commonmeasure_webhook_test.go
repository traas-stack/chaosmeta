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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertDuration(t *testing.T) {
	tests := []struct {
		name    string
		d       string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "seconds",
			d:       "10s",
			want:    10 * time.Second,
			wantErr: false,
		},
		{
			name:    "minutes",
			d:       "5m",
			want:    5 * time.Minute,
			wantErr: false,
		},
		{
			name:    "hours",
			d:       "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:    "invalid unit",
			d:       "10t",
			want:    0,
			wantErr: true,
		},
		{
			name:    "missing unit",
			d:       "10",
			want:    10 * time.Second,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDuration(tt.d)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
