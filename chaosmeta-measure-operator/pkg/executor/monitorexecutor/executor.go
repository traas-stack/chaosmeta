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

package monitorexecutor

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/config"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/monitorclient"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/utils"
	"strconv"
)

const (
	QueryArgsKey = "query"
)

func init() {
	e, err := NewMonitorExecutor(context.Background())
	if err != nil {
		fmt.Printf("new monitor executor error: %s\n", err.Error())
	} else {
		v1alpha1.SetMeasureExecutor(context.Background(), v1alpha1.MonitorMeasureType, e)
	}
}

type MonitorExecutor struct {
	client monitorclient.MonitorClient
}

func NewMonitorExecutor(ctx context.Context) (*MonitorExecutor, error) {
	client, err := monitorclient.NewMonitorClient(ctx,
		monitorclient.MonitorEngine(config.GetGlobalConfig().Monitor.Engine), config.GetGlobalConfig().Monitor.Url)
	if err != nil {
		return nil, fmt.Errorf("new %s monitor client error: %s", monitorclient.PrometheusMonitorEngine, err.Error())
	}

	return &MonitorExecutor{
		client: client,
	}, nil
}

func (e *MonitorExecutor) CheckConfig(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement) error {
	_, err := utils.GetArgsValueStr(args, QueryArgsKey)
	if err != nil {
		return fmt.Errorf("args error: %s", err.Error())
	}

	left, right, err := utils.GetIntervalValue(judgement.JudgeValue)
	if err != nil {
		return fmt.Errorf("get JudgeValue error: %s", err.Error())
	}

	if judgement.JudgeType == v1alpha1.RelativePercentJudgeType {
		if left < 0 || left > 100 {
			return fmt.Errorf("left percent value should meet [0, 100]")
		}
		if right < 0 || right > 100 {
			return fmt.Errorf("left percent value should meet [0, 100]")
		}
	} else if judgement.JudgeType != v1alpha1.RelativeValueJudgeType && judgement.JudgeType != v1alpha1.AbsoluteValueJudgeType {
		return fmt.Errorf("judge type of %s only support: %s,%s,%s", v1alpha1.MonitorMeasureType,
			v1alpha1.RelativeValueJudgeType, v1alpha1.RelativePercentJudgeType, v1alpha1.AbsoluteValueJudgeType)
	}

	return nil
}

func (e *MonitorExecutor) InitialData(ctx context.Context, args []v1alpha1.MeasureArgs) (string, error) {
	queryStr, _ := utils.GetArgsValueStr(args, QueryArgsKey)
	nowValue, err := e.client.GetNowValue(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("query monitor error: %s", err.Error())
	}

	return fmt.Sprintf("%f", nowValue), nil
}

func (e *MonitorExecutor) Measure(ctx context.Context, args []v1alpha1.MeasureArgs, judgement v1alpha1.Judgement, initialData string) error {
	queryStr, _ := utils.GetArgsValueStr(args, QueryArgsKey)
	nowValue, err := e.client.GetNowValue(ctx, queryStr)
	if err != nil {
		return fmt.Errorf("query monitor error: %s", err.Error())
	}

	left, right, err := utils.GetIntervalValue(judgement.JudgeValue)
	switch judgement.JudgeType {
	case v1alpha1.AbsoluteValueJudgeType:
		return utils.IfMeetInterval(nowValue, left, right)
	case v1alpha1.RelativeValueJudgeType:
		initialValue, _ := strconv.ParseFloat(initialData, 64)
		if left != v1alpha1.IntervalMin {
			left += initialValue
		}

		if right != v1alpha1.IntervalMax {
			right += initialValue
		}
		return utils.IfMeetInterval(nowValue, left, right)
	case v1alpha1.RelativePercentJudgeType:
		initialValue, _ := strconv.ParseFloat(initialData, 64)
		if left != v1alpha1.IntervalMin {
			left = initialValue * (1 + left/100)
		}

		if right != v1alpha1.IntervalMax {
			right = initialValue * (1 + right/100)
		}

		return utils.IfMeetInterval(nowValue, left, right)
	default:
		return fmt.Errorf("not support judge type: %s", judgement.JudgeType)
	}
}
