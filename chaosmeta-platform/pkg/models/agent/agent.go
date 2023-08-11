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

package agent

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
)

type Agent struct {
	ID               int    `json:"id" orm:"pk;auto;column(id)"`
	AgentType        string `json:"agentType" orm:"column(agent_type);size(32);index"`
	ClusterID        int    `json:"clusterId" orm:"column(cluster_id);index"`
	ContainerRuntime string `json:"containerRuntime" orm:"column(container_runtime);size(32);default(docker)"`
	Version          string `json:"version" orm:"column(version);size(32);index"`
	NodeName         string `json:"nodeName" orm:"column(node_name);size(255);"`
	Hostname         string `json:"hostname" orm:"column(host_name);size(255);index"`
	IP               string `json:"ip" orm:"column(ip);size(32);index"`
	Status           string `json:"status"  orm:"column(status);size(32);index"`
	AppID            int    `json:"appId" orm:"column(app_id);index"`
	SelectLabel      string `json:"selectLabel" orm:"column(select_label);size(1024);"`
	SelectNamespace  string `json:"selectNamespace" orm:"column(select_namespace);size(255);"`
	models.BaseTimeModel
}

func (a *Agent) TableName() string {
	return "agent"
}

func (a *Agent) TableUnique() [][]string {
	return [][]string{{"agent_type", "cluster_id", "version", "host_name", "ip", "status"}}
}

func InsertAgent(ctx context.Context, agent *Agent) (int64, error) {
	if agent == nil {
		return 0, errors.New("agent is nil")
	}
	id, err := models.GetORM().Insert(agent)
	return id, err
}

func InsertOrUpdateAgent(ctx context.Context, agent *Agent) (int64, error) {
	if err := GetAgentByHostname(ctx, agent); err != nil {
		return InsertAgent(ctx, agent)
	}
	return UpdateAgent(ctx, agent)
}

func UpdateAgent(ctx context.Context, agent *Agent) (int64, error) {
	if agent == nil {
		return 0, errors.New("agent is nil")
	}
	num, err := models.GetORM().Update(agent)
	return num, err
}

func DeleteAgent(ctx context.Context, id int) (int64, error) {
	num, err := models.GetORM().Delete(&Agent{ID: id})
	return num, err
}

func GetAgentById(ctx context.Context, agent *Agent) error {
	if agent == nil {
		return errors.New("agent is nil")
	}
	return models.GetORM().Read(agent)
}

func GetAgentByHostname(ctx context.Context, agent *Agent) error {
	return models.GetORM().Read(agent, "host_name")
}

func QueryAgents(ctx context.Context, hostname, ip, selectLabel, version, status, orderBy string, page, pageSize int) (int64, []Agent, error) {
	a, agents := Agent{}, new([]Agent)
	querySeter := models.GetORM().QueryTable(a.TableName())
	agentQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}
	if len(hostname) > 0 {
		agentQuery.Filter("host_name", models.CONTAINS, true, hostname)
	}
	if len(ip) > 0 {
		agentQuery.Filter("ip", models.CONTAINS, false, ip)
	}

	if len(selectLabel) > 0 {
		agentQuery.Filter("select_label", models.NEGLECT, false, selectLabel)
	}

	if len(version) > 0 {
		agentQuery.Filter("version", models.NEGLECT, false, version)
	}

	if len(status) > 0 {
		agentQuery.Filter("status", models.NEGLECT, false, status)
	}

	var totalCount int64
	totalCount, err = agentQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	agentQuery.OrderBy(orderByList...)

	if err := agentQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = agentQuery.GetOamQuerySeter().All(agents)
	return totalCount, *agents, err
}
