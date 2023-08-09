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

func (c *InjectController) QueryArgs() {
	faultId, _ := c.GetInt(":id")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 100)
	injectService := inject.InjectService{}

	total, args, err := injectService.ListArg(context.Background(), faultId, "", page, pageSize)
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
