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

package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type PrometheusClient struct {
	client api.Client
}

func (c *PrometheusClient) GetNowValue(ctx context.Context, query string) (float64, error) {
	v1api := v1.NewAPI(c.client)
	result, _, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}

	logger := log.FromContext(ctx)
	t := reflect.TypeOf(result).Name()
	logger.Info(fmt.Sprintf("prometheus data type: %s", t))
	switch t {
	case "Vector":
		v := result.(model.Vector)
		logger.Info(v.String())
		return float64(v[0].Value), nil
	case "Scalar":
		v := result.(*model.Scalar)
		logger.Info(v.String())
		return float64(v.Value), nil
	default:
		return 0, fmt.Errorf("not support prometheus data type: %s", t)
	}
}

func NewClient(ctx context.Context, url string) (*PrometheusClient, error) {
	config := api.Config{
		Address: url,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &PrometheusClient{
		client: client,
	}, nil
}
