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

package experiment_instance

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/experiment"
	"chaosmeta-platform/pkg/service/experiment_instance"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type ExperimentInstanceController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *ExperimentInstanceController) GetExperimentInstances() {
	lastInstance := c.GetString("last_instance")
	scheduleType := c.GetString("schedule_type")
	namespaceId, _ := c.GetInt("namespace_id")
	name := c.GetString("name")
	creator, _ := c.GetInt("creator", 0)
	timeType := c.GetString("time_type")
	recentDays, _ := c.GetInt("recent_days", 0)
	startTime, _ := time.Parse(experiment.TimeLayout, c.GetString("start_time"))
	endTime, _ := time.Parse(experiment.TimeLayout, c.GetString("end_time"))
	orderBy := c.GetString("sort")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	es := experiment_instance.ExperimentInstanceService{}
	total, experiments, err := es.SearchExperimentInstances(lastInstance, namespaceId, creator, name, scheduleType, timeType, recentDays, startTime, endTime, orderBy, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, ExperimentInstanceListResponse{
		Page:        page,
		PageSize:    pageSize,
		Total:       total,
		Experiments: experiments,
	})
}

func (c *ExperimentInstanceController) GetExperimentInstanceDetail() {
	uuid := c.GetString(":uuid")
	es := experiment_instance.ExperimentInstanceService{}
	experiment, err := es.GetExperimentInstanceByUUID(uuid)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, experiment)
}

func (c *ExperimentInstanceController) GetExperimentInstanceNodes() {
	uuid := c.GetString(":uuid")
	es := experiment_instance.ExperimentInstanceService{}
	total, nodes, err := es.GetWorkflowNodesInstanceInfoByUUID(uuid)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, GetExperimentInstancesResponse{Total: total, WorkflowNodes: nodes})
}

func (c *ExperimentInstanceController) GetExperimentInstanceNode() {
	uuid := c.GetString(":uuid")
	nodeId := c.GetString(":node_id")
	es := experiment_instance.ExperimentInstanceService{}
	nodeDetail, err := es.GetWorkflowNodeInstanceDetailByUUIDAndNodeId(uuid, nodeId)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, GetExperimentInstanceResponse{WorkflowNode: *nodeDetail})
}

func (c *ExperimentInstanceController) GetExperimentInstanceNodeSubtask() {
	uuid := c.GetString(":uuid")
	nodeId := c.GetString(":node_id")
	subtaskId := c.GetString(":subtask_id")
	es := experiment_instance.ExperimentInstanceService{}
	rangeInstance, err := es.GetFaultRangeInstanceByWorkflowNodeInstanceUUID(uuid, nodeId, subtaskId)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, GetFaultRangeInstanceResponse{FaultRangeInstance: *rangeInstance})
}

func (c *ExperimentInstanceController) DeleteExperimentInstance() {
	uuid := c.GetString(":uuid")
	es := experiment_instance.ExperimentInstanceService{}
	if err := es.DeleteExperimentInstanceByUUID(uuid); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *ExperimentInstanceController) DeleteExperimentInstances() {
	var reqBody DeleteExperimentInstanceRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	es := experiment_instance.ExperimentInstanceService{}
	if err := es.DeleteExperimentInstancesByUUID(reqBody.ResultUUIDs); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}
