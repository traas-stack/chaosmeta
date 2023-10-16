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
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlutil "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"time"
)

const (
	FlowKind     = "LoadTest"
	FlowResource = "loadtests"
)

type ChaosmetaFlowInterface interface {
	Get(ctx context.Context, namespace, name string) (result *LoadTest, err error)
	List(ctx context.Context, namespace string) (*LoadTestList, error)
	Create(ctx context.Context, chaosmeta *LoadTest) (*LoadTest, error)
	Update(ctx context.Context, chaosmeta *LoadTest) (*LoadTest, error)
	Delete(ctx context.Context, namespace, name string) error
	Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error
	DeleteExpiredList(ctx context.Context, namespace string) error
	Recover(namespace, name string) error
}

type ChaosmetaFlowService struct {
	Config *rest.Config
	Client dynamic.Interface
}

var gvrFlow = schema.GroupVersionResource{
	Group:    Group,
	Version:  Version,
	Resource: FlowResource,
}

var gvkFlow = schema.GroupVersionKind{
	Group:   Group,
	Version: Version,
	Kind:    FlowKind,
}

func NewChaosmetaFlowService(config *rest.Config) ChaosmetaFlowInterface {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil
	}

	return &ChaosmetaFlowService{
		Config: config, Client: client,
	}
}

func (c *ChaosmetaFlowService) Get(ctx context.Context, namespace, name string) (result *LoadTest, err error) {
	cb, err := c.Client.Resource(gvrFlow).Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := cb.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var flow LoadTest
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}

	return &flow, nil
}

func (c *ChaosmetaFlowService) List(ctx context.Context, namespace string) (*LoadTestList, error) {
	list, err := c.Client.Resource(gvrFlow).Namespace(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var ltList LoadTestList
	if err := json.Unmarshal(data, &ltList); err != nil {
		return nil, err
	}

	return &ltList, nil
}

func (c *ChaosmetaFlowService) Create(ctx context.Context, chaosmeta *LoadTest) (*LoadTest, error) {
	d, err := json.Marshal(chaosmeta)
	if err != nil {
		return nil, err
	}

	y, err := gyaml.JSONToYAML(d)
	if err != nil {
		return nil, err
	}

	decoder := yamlutil.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode(y, &gvkFlow, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvrFlow).Create(ctx, obj, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var flow LoadTest
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}
	if len(flow.Status.Status) > 0 {
		if flow.Status.Status == FailedStatusType {
			return nil, errors.New("chaosmeta failed to execute")
		}
	}
	return &flow, nil
}

func (c *ChaosmetaFlowService) Update(ctx context.Context, chaosmeta *LoadTest) (*LoadTest, error) {
	d, err := json.Marshal(chaosmeta)
	if err != nil {
		return nil, err
	}

	y, err := gyaml.JSONToYAML(d)
	if err != nil {
		return nil, err
	}
	decoder := yamlutil.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode(y, &gvkFlow, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvrFlow).Namespace(obj.GetNamespace()).Get(ctx, obj.GetName(), v1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = c.Client.Resource(gvrFlow).Namespace(obj.GetNamespace()).Update(ctx, obj, v1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var flow LoadTest
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}
	return &flow, nil
}

func (c *ChaosmetaFlowService) Delete(ctx context.Context, namespace, name string) error {
	return c.Client.Resource(gvrFlow).Namespace(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

func (c *ChaosmetaFlowService) Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error {
	_, err := c.Client.Resource(gvrFlow).Namespace(namespace).Patch(ctx, name, pt, data, v1.PatchOptions{})
	return err
}

func (c *ChaosmetaFlowService) DeleteExpiredList(ctx context.Context, namespace string) error {
	chaosmetaList, err := c.List(ctx, namespace)
	if err != nil {
		return err
	}
	for _, measure := range chaosmetaList.Items {
		expirationTime := time.Now().AddDate(0, 0, -1)
		experimentCreateTime, err := time.Parse(TimeLayout, measure.Status.CreateTime)
		if err != nil {
			return err
		}

		if measure.Status.Status == SuccessStatusType && measure.Status.CreateTime != "" && experimentCreateTime.Before(expirationTime) {
			err := c.Delete(ctx, measure.Namespace, measure.Name)
			if err != nil {
				log.Infof("failed to delete chaosmeta measure %s: %v", measure.Name, err.Error())
				return err
			} else {
				log.Infof("chaosmeta measure %s deleted", measure.Name)
			}
		}
	}
	return nil
}

func (c *ChaosmetaFlowService) Recover(namespace, name string) error {
	chaosmetaCR, err := c.Get(context.Background(), namespace, name)
	if err != nil {
		log.Error(err)
		return err
	}
	chaosmetaCR.Spec.Stopped = true
	if _, err := c.Update(context.Background(), chaosmetaCR); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
