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
	Group               = "chaosmeta.io"
	Version             = "v1alpha1"
	ExperimentKind      = "Experiment"
	ExperimentsResource = "experiments"
	RecoverTargetPhase  = "recover"
)

type ChaosmetaInterface interface {
	Get(ctx context.Context, namespace, name string) (result *ExperimentInjectStruct, err error)
	List(ctx context.Context, namespace string) (*ExperimentInjectStructList, error)
	Create(ctx context.Context, chaosmeta *ExperimentInjectStruct) (*ExperimentInjectStruct, error)
	Update(ctx context.Context, chaosmeta *ExperimentInjectStruct) (*ExperimentInjectStruct, error)
	Delete(ctx context.Context, namespace, name string) error
	Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error
	DeleteExpiredList(ctx context.Context, namespace string) error
	Recover(namespace, name string) error
}

type ChaosmetaService struct {
	Config *rest.Config
	Client dynamic.Interface
}

var gvr = schema.GroupVersionResource{
	Group:    Group,
	Version:  Version,
	Resource: ExperimentsResource,
}

var gvk = schema.GroupVersionKind{
	Group:   Group,
	Version: Version,
	Kind:    ExperimentKind,
}

func NewChaosmetaService(config *rest.Config) ChaosmetaInterface {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil
	}

	return &ChaosmetaService{
		Config: config, Client: client,
	}
}

func (c *ChaosmetaService) Get(ctx context.Context, namespace, name string) (result *ExperimentInjectStruct, err error) {
	cb, err := c.Client.Resource(gvr).Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := cb.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var experiment ExperimentInjectStruct
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}

	return &experiment, nil
}

func (c *ChaosmetaService) List(ctx context.Context, namespace string) (*ExperimentInjectStructList, error) {
	list, err := c.Client.Resource(gvr).Namespace(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var exList ExperimentInjectStructList
	if err := json.Unmarshal(data, &exList); err != nil {
		return nil, err
	}

	return &exList, nil
}

func (c *ChaosmetaService) Create(ctx context.Context, chaosmeta *ExperimentInjectStruct) (*ExperimentInjectStruct, error) {
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
	var experiment ExperimentInjectStruct
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}
	if len(experiment.Status.Status) > 0 {
		if experiment.Status.Status == FailedStatusType {
			return nil, errors.New("chaosmeta failed to execute")
		}
	}
	return &experiment, nil
}

func (c *ChaosmetaService) Update(ctx context.Context, chaosmeta *ExperimentInjectStruct) (*ExperimentInjectStruct, error) {
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
	if _, _, err := decoder.Decode(y, &gvk, obj); err != nil {
		return nil, err
	}

	utd, err := c.Client.Resource(gvr).Namespace(obj.GetNamespace()).Get(ctx, obj.GetName(), v1.GetOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = c.Client.Resource(gvr).Namespace(obj.GetNamespace()).Update(ctx, obj, v1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var experiment ExperimentInjectStruct
	if err := json.Unmarshal(data, &experiment); err != nil {
		return nil, err
	}
	return &experiment, nil
}

func (c *ChaosmetaService) Delete(ctx context.Context, namespace, name string) error {
	return c.Client.Resource(gvr).Namespace(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

func (c *ChaosmetaService) Patch(ctx context.Context, namespace, name string, pt types.PatchType, data []byte) error {
	_, err := c.Client.Resource(gvr).Namespace(namespace).Patch(ctx, name, pt, data, v1.PatchOptions{})
	return err
}

func (c *ChaosmetaService) DeleteExpiredList(ctx context.Context, namespace string) error {
	chaosmetaList, err := c.List(ctx, namespace)
	if err != nil {
		return err
	}
	for _, experiment := range chaosmetaList.Items {
		expirationTime := time.Now().AddDate(0, 0, -1)
		experimentCreateTime, err := time.Parse(TimeLayout, experiment.Status.CreateTime)
		if err != nil {
			return err
		}

		if experiment.Status.Status == SuccessStatusType && experiment.Status.CreateTime != "" && experimentCreateTime.Before(expirationTime) {
			err := c.Delete(ctx, experiment.Namespace, experiment.Name)
			if err != nil {
				log.Infof("failed to delete chaosmeta experiment %s: %v", experiment.Name, err.Error())
				return err
			} else {
				log.Infof("chaosmeta experiment %s deleted", experiment.Name)
			}
		}
	}
	return nil
}

func (c *ChaosmetaService) Recover(namespace, name string) error {
	chaosmetaCR, err := c.Get(context.Background(), namespace, name)
	if err != nil {
		log.Error(err)
		return err
	}
	chaosmetaCR.Spec.TargetPhase = RecoverTargetPhase
	if _, err := c.Update(context.Background(), chaosmetaCR); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
