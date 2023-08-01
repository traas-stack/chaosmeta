package cluster

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	"chaosmeta-platform/pkg/service/cluster"
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
)

type ClusterController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

func (c *ClusterController) Create() {
	var requestBody CreateClusterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	username := c.Ctx.Input.GetData("userName").(string)
	log.Error(username, "create:", requestBody.Name)
	clusterService := &cluster.ClusterService{}
	clusterId, err := clusterService.Create(context.Background(), requestBody.Name, requestBody.Kubeconfig)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, CreateClusterResponse{
		ID: clusterId,
	})
}

func (c *ClusterController) Get() {
	clusterId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	clusterService := &cluster.ClusterService{}
	cluster, err := clusterService.Get(context.Background(), clusterId)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, ClusterData{
		Id:         cluster.ID,
		Name:       cluster.Name,
		Kubeconfig: cluster.KubeConfig,
	})
}

func (c *ClusterController) GetList() {
	sort := c.GetString("sort")
	name := c.GetString("name")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	clusterService := &cluster.ClusterService{}
	total, clusterList, err := clusterService.GetList(context.Background(), name, sort, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	listClusterResponse := ListClusterResponse{Total: total, Page: page, PageSize: pageSize}

	for _, cluster := range clusterList {
		listClusterResponse.Clusters = append(listClusterResponse.Clusters, ClusterData{
			Id:   cluster.ID,
			Name: cluster.Name,
		})
	}
	c.Success(&c.Controller, listClusterResponse)
}

func (c *ClusterController) Update() {
	clusterId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}

	var requestBody UpdateClusterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	username := c.Ctx.Input.GetData("userName").(string)
	log.Error(username, "Update:", requestBody.Name)

	clusterService := &cluster.ClusterService{}
	if err := clusterService.Update(context.Background(), clusterId, requestBody.Name, requestBody.Kubeconfig); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}

func (c *ClusterController) Delete() {
	clusterId, err := c.GetInt(":id")
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	username := c.Ctx.Input.GetData("userName").(string)
	log.Error(username, "delete:", clusterId)
	clusterService := &cluster.ClusterService{}
	if err := clusterService.Delete(context.Background(), clusterId); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}
