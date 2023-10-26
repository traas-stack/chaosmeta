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

package inject

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/inject"
	"context"
	beego "github.com/beego/beego/v2/server/web"
)

type InjectController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *InjectController) QueryScopes() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}
	total, scopes, err := injectService.ListScopes(context.Background(), "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	scopesListResponse := ScopesListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Scopes:   scopes,
	}
	c.Success(&c.Controller, scopesListResponse)
}

func (c *InjectController) QueryTargets() {
	scopeId, _ := c.GetInt(":id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}
	total, targets, err := injectService.ListTargets(context.Background(), scopeId, "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	targetsListResponse := TargetsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Targets:  targets,
	}
	c.Success(&c.Controller, targetsListResponse)
}

func (c *InjectController) GetTarget() {
	id, _ := c.GetInt(":id")
	injectService := inject.InjectService{}
	target, err := injectService.GetTarget(context.Background(), id)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, target)
}

func (c *InjectController) QueryFaults() {
	targetId, _ := c.GetInt(":targets_id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}
	total, faults, err := injectService.ListFault(context.Background(), targetId, "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	faultsListResponse := FaultsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Faults:   faults,
	}
	c.Success(&c.Controller, faultsListResponse)
}

func (c *InjectController) QueryFlows() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}
	total, faults, err := injectService.ListFlows(context.Background(), "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	flowsListResponse := FlowsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Flows:    faults,
	}
	c.Success(&c.Controller, flowsListResponse)
}

func (c *InjectController) QueryMeasures() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}
	total, measures, err := injectService.ListMeasures(context.Background(), "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	flowsListResponse := MeasuresListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Measures: measures,
	}
	c.Success(&c.Controller, flowsListResponse)
}

func (c *InjectController) QueryMeasureArgs() {
	faultId, _ := c.GetInt(":id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}

	total, args, err := injectService.ListArg(context.Background(), []string{inject.ExecMeasureCommon, inject.ExecMeasure}, faultId, "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	argsListResponse := ArgsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Args:     args,
	}
	c.Success(&c.Controller, argsListResponse)
}

func (c *InjectController) QueryFaultArgs() {
	faultId, _ := c.GetInt(":id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}

	total, args, err := injectService.ListArg(context.Background(), []string{inject.ExecInject}, faultId, "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	argsListResponse := ArgsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Args:     args,
	}
	c.Success(&c.Controller, argsListResponse)
}

func (c *InjectController) QueryFlowArgs() {
	faultId, _ := c.GetInt(":id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}

	total, args, err := injectService.ListArg(context.Background(), []string{inject.ExecFlowCommon, inject.ExecFlow}, faultId, "", page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	argsListResponse := ArgsListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Args:     args,
	}
	c.Success(&c.Controller, argsListResponse)
}
