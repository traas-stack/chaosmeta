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

package selector

import (
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestGetTargetContainer(t *testing.T) {
	testStatus := []corev1.ContainerStatus{
		{Name: "chaosmeta", ContainerID: "docker://33124124"},
		{Name: "nginx", ContainerID: "pouch://esbersbsh"},
		{Name: "centos", ContainerID: ""},
	}

	type args struct {
		containerName string
		status        []corev1.ContainerStatus
	}
	tests := []struct {
		name     string
		args     args
		wantR    string
		wantId   string
		wantName string
		wantErr  bool
	}{
		{
			name: "first",
			args: args{
				containerName: v1alpha1.FirstContainer,
				status:        testStatus,
			},
			wantR:    "docker",
			wantId:   "33124124",
			wantName: "chaosmeta",
			wantErr:  false,
		},
		{
			name: "empty",
			args: args{
				containerName: v1alpha1.FirstContainer,
				status:        nil,
			},
			//wantR:   "docker",
			//wantId:  "33124124",
			wantErr: true,
		},
		{
			name: "target",
			args: args{
				containerName: "nginx",
				status:        testStatus,
			},
			wantR:    "pouch",
			wantId:   "esbersbsh",
			wantName: "nginx",
			wantErr:  false,
		},
		{
			name: "wrong format",
			args: args{
				containerName: "centos",
				status:        testStatus,
			},
			wantName: "centos",
			//wantR:   "pouch",
			//wantId:  "esbersbsh",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containers, err := GetTargetContainers(tt.args.containerName, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(containers) == 0 {
				return
			}
			gotR, gotId, gotName := containers[0].ContainerRuntime, containers[0].ContainerId, containers[0].ContainerName
			if gotR != tt.wantR {
				t.Errorf("GetTargetContainer() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotId != tt.wantId {
				t.Errorf("GetTargetContainer() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotName != tt.wantName {
				t.Errorf("GetTargetContainer() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
