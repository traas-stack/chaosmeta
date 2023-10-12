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
	MeasureKind     = "CommonMeasure"
	MeasureResource = "commonmeasures"
)

type ChaosmetaMeasureInterface interface {
	Get(ctx context.Context, namespace, name string) (result *CommonMeasureStruct, err error)
	List(ctx context.Context, namespace string) (*CommonMeasureList, error)
	Create(ctx context.Context, chaosmeta *CommonMeasureStruct) (*CommonMeasureStruct, error)
	Update(ctx context.Context, chaosmeta *CommonMeasureStruct) (*CommonMeasureStruct, error)
	Delete(ctx context.Context, namespace, name string) error
	Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error
	DeleteExpiredList(ctx context.Context, namespace string) error
}

type ChaosmetaMeasureService struct {
	Config *rest.Config
	Client dynamic.Interface
}

var gvrMeasure = schema.GroupVersionResource{
	Group:    Group,
	Version:  Version,
	Resource: MeasureResource,
}

var gvkMeasure = schema.GroupVersionKind{
	Group:   Group,
	Version: Version,
	Kind:    MeasureKind,
}

func NewChaosmetaMeasureService(config *rest.Config) ChaosmetaMeasureInterface {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil
	}

	return &ChaosmetaMeasureService{
		Config: config, Client: client,
	}
}

func (c *ChaosmetaMeasureService) Get(ctx context.Context, namespace, name string) (result *CommonMeasureStruct, err error) {
	cb, err := c.Client.Resource(gvrMeasure).Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := cb.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var measure CommonMeasureStruct
	if err := json.Unmarshal(data, &measure); err != nil {
		return nil, err
	}

	return &measure, nil
}

func (c *ChaosmetaMeasureService) List(ctx context.Context, namespace string) (*CommonMeasureList, error) {
	list, err := c.Client.Resource(gvrMeasure).Namespace(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var cmList CommonMeasureList
	if err := json.Unmarshal(data, &cmList); err != nil {
		return nil, err
	}

	return &cmList, nil
}

func (c *ChaosmetaMeasureService) Create(ctx context.Context, chaosmeta *CommonMeasureStruct) (*CommonMeasureStruct, error) {
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
	if _, _, err := decoder.Decode(y, &gvkMeasure, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvrMeasure).Create(ctx, obj, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var measure CommonMeasureStruct
	if err := json.Unmarshal(data, &measure); err != nil {
		return nil, err
	}
	if len(measure.Status.Status) > 0 {
		if measure.Status.Status == FailedStatusType {
			return nil, errors.New("chaosmeta failed to execute")
		}
	}
	return &measure, nil
}

func (c *ChaosmetaMeasureService) Update(ctx context.Context, chaosmeta *CommonMeasureStruct) (*CommonMeasureStruct, error) {
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
	if _, _, err := decoder.Decode(y, &gvkMeasure, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvrMeasure).Namespace(obj.GetNamespace()).Get(ctx, obj.GetName(), v1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = c.Client.Resource(gvrMeasure).Namespace(obj.GetNamespace()).Update(ctx, obj, v1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var measure CommonMeasureStruct
	if err := json.Unmarshal(data, &measure); err != nil {
		return nil, err
	}
	return &measure, nil
}

func (c *ChaosmetaMeasureService) Delete(ctx context.Context, namespace, name string) error {
	return c.Client.Resource(gvrMeasure).Namespace(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

func (c *ChaosmetaMeasureService) Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error {
	_, err := c.Client.Resource(gvrMeasure).Namespace(namespace).Patch(ctx, name, pt, data, v1.PatchOptions{})
	return err
}

func (c *ChaosmetaMeasureService) DeleteExpiredList(ctx context.Context, namespace string) error {
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
