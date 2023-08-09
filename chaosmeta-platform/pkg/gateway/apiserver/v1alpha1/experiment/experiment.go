package experiment

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/experiment"
	"encoding/json"
	"errors"
	beego "github.com/beego/beego/v2/server/web"
)

type ExperimentController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *ExperimentController) GetExperimentList() {
	lastInstanceStatus := c.GetString("last_instance_status")
	scheduleType := c.GetString("schedule_type")
	namespaceID, _ := c.GetInt("namespace_id")
	name := c.GetString("name")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	orderBy := c.GetString("sort")

	experimentService := experiment.ExperimentService{}

	total, experimentList, err := experimentService.SearchExperiments(lastInstanceStatus, namespaceID, name, scheduleType, orderBy, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	experimentListResponse := ExperimentListResponse{
		Page:        page,
		PageSize:    pageSize,
		Total:       total,
		Experiments: experimentList,
	}

	c.Success(&c.Controller, experimentListResponse)
}

func (c *ExperimentController) GetExperimentDetail() {
	uuid := c.Ctx.Input.Param(":uuid")
	experimentService := experiment.ExperimentService{}

	experimentGet, err := experimentService.GetExperimentByUUID(uuid)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, GetExperimentResponse{
		Experiment: *experimentGet,
	})
}

func (c *ExperimentController) CreateExperiment() {
	username := c.Ctx.Input.GetData("userName").(string)
	experimentService := experiment.ExperimentService{}
	creatorId, err := experimentService.GetCreator(username)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var createExperimentRequest experiment.Experiment
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &createExperimentRequest); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	uuid, err := experimentService.CreateExperiment(creatorId, &createExperimentRequest)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, CreateExperimentResponse{
		Uuid: uuid,
	})
}

func (c *ExperimentController) UpdateExperiment() {
	uuid := c.Ctx.Input.Param(":uuid")
	experimentService := experiment.ExperimentService{}

	var updateExperimentRequest experiment.Experiment
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &updateExperimentRequest); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	if err := experimentService.UpdateExperiment(uuid, &updateExperimentRequest); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *ExperimentController) DeleteExperiment() {
	uuid := c.Ctx.Input.Param(":uuid")
	if uuid == "" {
		c.Error(&c.Controller, errors.New("uuid is empty"))
		return
	}

	experimentService := experiment.ExperimentService{}
	if err := experimentService.DeleteExperimentByUUID(uuid); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}
