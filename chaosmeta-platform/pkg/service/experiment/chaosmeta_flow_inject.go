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

import (
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	"errors"
	gyaml "github.com/ghodss/yaml"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type ChaosmetaFlowInjectInterface interface {
	Get(ctx context.Context, namespace, name string) (result *LoadTest, err error)
	List(ctx context.Context) (*LoadTestList, error)
	Create(ctx context.Context, ChaosmetaFlowInject *LoadTest) (*LoadTest, error)
	Update(ctx context.Context, ChaosmetaFlowInject *LoadTest) (*LoadTest, error)
	Delete(ctx context.Context, name string) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte) error
	DeleteExpiredList(ctx context.Context) error
}

type ChaosmetaFlowInjectService struct {
	Config *rest.Config
	Client dynamic.Interface
}

func NewChaosmetaFlowInjectService(config *rest.Config) ChaosmetaFlowInjectInterface {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil
	}

	return &ChaosmetaFlowInjectService{
		Config: config, Client: client,
	}
}

func (c *ChaosmetaFlowInjectService) Get(ctx context.Context, namespace, name string) (result *LoadTest, err error) {
	cb, err := c.Client.Resource(gvr).Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := cb.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var experiment LoadTest
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}

	return &experiment, nil
}

func (c *ChaosmetaFlowInjectService) List(ctx context.Context) (*LoadTestList, error) {
	list, err := c.Client.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var exList LoadTestList
	if err := json.Unmarshal(data, &exList); err != nil {
		return nil, err
	}

	return &exList, nil
}

func (c *ChaosmetaFlowInjectService) Create(ctx context.Context, ChaosmetaFlowInject *LoadTest) (*LoadTest, error) {
	d, err := json.Marshal(ChaosmetaFlowInject)
	if err != nil {
		return nil, err
	}

	y, err := gyaml.JSONToYAML(d)
	if err != nil {
		return nil, err
	}

	decoder := yamlutil.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode(y, &gvk, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvr).Create(ctx, obj, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var experiment LoadTest
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}
	if len(experiment.Status) > 0 {
		if experiment.Status == FailedStatus {
			return nil, errors.New("ChaosmetaFlowInject failed to execute")
		}
	}
	return &experiment, nil
}

func (c *ChaosmetaFlowInjectService) Update(ctx context.Context, ChaosmetaFlowInject *LoadTest) (*LoadTest, error) {
	d, err := json.Marshal(ChaosmetaFlowInject)
	if err != nil {
		return nil, err
	}

	y, err := gyaml.JSONToYAML(d)
	if err != nil {
		return nil, err
	}
	decoder := yamlutil.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode(y, &gvk, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvr).Get(ctx, obj.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = c.Client.Resource(gvr).Update(ctx, obj, v1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var experiment LoadTest
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}
	return &experiment, nil
}

func (c *ChaosmetaFlowInjectService) Delete(ctx context.Context, name string) error {
	return c.Client.Resource(gvr).Delete(ctx, name, v1.DeleteOptions{})
}

func (c *ChaosmetaFlowInjectService) Patch(ctx context.Context, name string, pt types.PatchType, data []byte) error {
	_, err := c.Client.Resource(gvr).Patch(ctx, name, pt, data, v1.PatchOptions{})
	return err
}

func (c *ChaosmetaFlowInjectService) DeleteExpiredList(ctx context.Context) error {
	ChaosmetaFlowInjectList, err := c.List(ctx)
	if err != nil {
		return err
	}
	for _, experiment := range ChaosmetaFlowInjectList.Items {
		if experiment.Status == SuccessStatus {
			err := c.Delete(ctx, experiment.Name)
			if err != nil {
				log.Errorf("failed to delete ChaosmetaFlowInject experiment %s: %v", experiment.Name, err.Error())
				return err
			} else {
				log.Errorf("ChaosmetaFlowInject experiment %s deleted", experiment.Name)
			}
		}
	}
	return nil
}
