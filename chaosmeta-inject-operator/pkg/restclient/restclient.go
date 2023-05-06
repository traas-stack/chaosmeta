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

package restclient

import (
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var (
	apiServerClientMap = make(map[v1alpha1.CloudTargetType]rest.Interface)
)

func GetApiServerClientMap(targetType v1alpha1.CloudTargetType) rest.Interface {
	return apiServerClientMap[targetType]
}

func SetApiServerClientMap(c *rest.Config, s *runtime.Scheme, t []v1alpha1.CloudTargetType) error {
	for _, unitTarget := range t {
		e, err := newClient(unitTarget, c, s)
		if err != nil {
			return fmt.Errorf("create apiserver client for %s error: %s", unitTarget, err.Error())
		}

		apiServerClientMap[unitTarget] = e
	}

	return nil
}

func newClient(target v1alpha1.CloudTargetType, c *rest.Config, s *runtime.Scheme) (e rest.Interface, err error) {
	switch target {
	case v1alpha1.PodCloudTarget:
		e, err = newRESTClientForGVK("", "v1", "Pod", c, s)
	case v1alpha1.DeploymentCloudTarget:
		e, err = newRESTClientForGVK("apps", "v1", "Deployment", c, s)
	case v1alpha1.NodeCloudTarget:
		e, err = newRESTClientForGVK("", "v1", "Node", c, s)
	case v1alpha1.NamespaceCloudTarget:
		e, err = newRESTClientForGVK("", "v1", "Namespace", c, s)
	case v1alpha1.JobCloudTarget:
		e, err = newRESTClientForGVK("batch", "v1", "Job", c, s)
	default:
		err = fmt.Errorf("not support target: %s", target)
	}
	return
}

func newRESTClientForGVK(group, version, kind string, c *rest.Config, s *runtime.Scheme) (rest.Interface, error) {
	return apiutil.RESTClientForGVK(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, false, c, serializer.NewCodecFactory(s))
}
