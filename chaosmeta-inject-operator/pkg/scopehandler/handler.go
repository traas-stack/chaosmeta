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

package scopehandler

import (
	"context"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler/kubernetes"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler/node"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler/pod"
)

type ScopeHandler interface {
	ConvertSelector(ctx context.Context, spec *v1alpha1.ExperimentSpec) ([]model.AtomicObject, error)
	QueryExperiment(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon, phase v1alpha1.PhaseType) (*model.SubExpInfo, error)
	ExecuteInject(ctx context.Context, injectObject model.AtomicObject, UID string, expArgs *v1alpha1.ExperimentCommon) (string, error)
	ExecuteRecover(ctx context.Context, injectObject model.AtomicObject, UID, backup string, expArgs *v1alpha1.ExperimentCommon) error
	GetInjectObject(ctx context.Context, exp *v1alpha1.ExperimentCommon, objectName string) (model.AtomicObject, error)
	CheckAlive(ctx context.Context, injectObject model.AtomicObject) error
}

func GetScopeHandler(scope v1alpha1.ScopeType) ScopeHandler {
	switch scope {
	case v1alpha1.PodScopeType:
		return pod.GetGlobalPodHandler()
	case v1alpha1.NodeScopeType:
		return node.GetGlobalNodeHandler()
	case v1alpha1.KubernetesScopeType:
		return kubernetes.GetGlobalKubernetesHandler()
	default:
		return nil
	}
}
