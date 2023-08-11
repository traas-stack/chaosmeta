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

package node

import (
	"chaosmeta-platform/pkg/gateway/apiserver/v1alpha1"
	beego "github.com/beego/beego/v2/server/web"
)

type NodeController struct {
	v1alpha1.BeegoOutputController
	beego.Controller
}

/*
func (n *NodeController) ListNode() {
	cluster := n.GetString("cluster")
	dsQuery := page.ParseDataSelectPathParameter(gc)

	nodeController, err := kube.NewNodeService(cluster)
	if err != nil {
		n.Error(gc, err)
		return
	}

	nodes, err := nodeController.List(dsQuery)
	if err != nil {
		n.Error(gc, err)
		return
	}

	monCtrl, err := n.clientset.GetPrometheusClient(cluster)
	if err != nil {
		n.Error(gc, err)
		return
	}

	if monCtrl != nil {
		defer func() {
			if err != nil {
				n.ErrorWithMessage(gc, fmt.Sprintf("fail when get node metrics|%v", err))
				return
			}
		}()

		var (
			wg  sync.WaitGroup
			gp  *ants.Pool
			lth = len(nodes.List)
			ech = make(chan error, lth)
			pch = make(chan map[string][]monitor2.Metric, lth)
		)
		gp, err = ants.NewPool(20)
		if err != nil {
			n.ErrorWithMessage(gc, fmt.Sprintf("fail to new goroutine pool, caused by: %v", err))
			return
		}
		defer gp.Release()

		for _, node := range nodes.List {
			wg.Add(1)
			node := node
			err = gp.Submit(func() {
				defer wg.Done()

				metric := monCtrl.GetNamedMetrics([]string{
					"node_cpu_utilisation",
					"node_cpu_total",
					"node_cpu_usage",
					"node_memory_utilisation",
					"node_memory_available",
					"node_memory_total",
					"node_disk_size_capacity",
					"node_disk_size_available",
					"node_disk_size_usage",
					"node_disk_size_utilisation",
					"node_pod_count",
					"node_pod_quota",
					"node_pod_utilisation",
					"node_pod_running_count",
					"node_pod_succeeded_count",
					"node_pod_abnormal_count",
					"node_load1",
					"node_load5",
					"node_load15",
					"node_pod_abnormal_ratio",
				}, time.Now(), monitor2.NodeOption{
					NodeName: node.Name,
				})
				res := make(map[string][]monitor2.Metric)
				res[node.Name] = metric
				pch <- res
			})
			if err != nil {
				n.Error(gc, fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err))
				return
			}
		}
		wg.Wait()
		close(pch)
		select {
		case err = <-ech:
			return
		default:
			// do nothing
		}

		for p := range pch {
			for k, v := range p {
				for i, n := range nodes.List {
					if k == n.Name {
						nodes.List[i].Metrics = v
					}
				}
			}
		}
	}
	n.Success(gc, nodes)
}

func (n *NodeController) GetNode(gc *gin.Context) {
	var nodeInfo NodeInfo

	cluster := gc.Param("cluster")
	nodeName := gc.Param("node")

	nodeController, err := n.clientset.NewNodeController(cluster)
	if err != nil {
		n.Error(gc, err)
		return
	}

	node, err := nodeController.Get(nodeName)
	if err != nil {
		n.Error(gc, err)
		return
	}

	nodeInfo.Node = node.Node
	nodeInfo.AllocatedResources = node.AllocatedResources
	ackController, err := n.clientset.NewACKController(cluster)
	if err != nil {
		n.Success(gc, nodeInfo)
		return
	}
	cloudInfo, err := ackController.GetNodeDetail("", nodeName)
	if err != nil {
		n.Error(gc, err)
		return
	}

	nodeInfo.CloudInfo = cloudInfo
	n.Success(gc, nodeInfo)
	return
}

func (n *NodeController) ScheduleNode(gc *gin.Context) {
	jwtHeader := gc.GetHeader(constants.GwJWTHeader)
	user, err := sv1alpha1.GetUser(jwtHeader)
	if err != nil {
		n.ErrorWithMessage(gc, fmt.Sprintf("failed to parse params, error: %s", err))
		return
	}
	// 更改nodeschedule
	nodeInfo := Node{}
	err = gc.ShouldBindJSON(&nodeInfo)
	if err != nil {
		n.Error(gc, err)
		return
	}

	if nodeInfo.Cluster == "" || nodeInfo.NodeName == "" {
		n.ErrorWithMessage(gc, "missing cluster or node")
		return
	}

	node := v1alpha1.GalaxyNode{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("node-schedule-%v", time.Now().Unix()),
		},
		Spec: v1alpha1.GalaxyNodeSpec{
			User:    fmt.Sprintf("%s(%s)", user.RealmName, user.NickNameCn),
			Cluster: nodeInfo.Cluster,
			Env:     nodeInfo.Env,
			Operate: v1alpha1.GalaxyNodeOperate{
				Type: v1alpha1.ScheduleGalaxyNode,
				Node: nodeInfo.NodeName,
				OperateParams: v1alpha1.OperateParams{
					Schedule: nodeInfo.UnSchedule,
				},
			},
		},
	}
	err = n.galaxyNodeController.Create(&node)
	if err != nil {
		n.Error(gc, err)
		return
	}
	n.SuccessNoData(gc)
}

func (n *NodeController) PatchNode(gc *gin.Context) {
	jwtHeader := gc.GetHeader(constants.GwJWTHeader)
	user, err := sv1alpha1.GetUser(jwtHeader)
	if err != nil {
		n.ErrorWithMessage(gc, fmt.Sprintf("failed to parse params, error: %s", err))
		return
	}
	// 更新node的label或者annotation
	nodePatch := kube.ReplaceNodeInfoParam{}

	err = gc.ShouldBindJSON(&nodePatch)
	if err != nil {
		n.Error(gc, err)
		return
	}

	node := v1alpha1.GalaxyNode{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("node-patch-%v", time.Now().Unix()),
		},
		Spec: v1alpha1.GalaxyNodeSpec{
			User:    fmt.Sprintf("%s(%s)", user.RealmName, user.NickNameCn),
			Env:     nodePatch.Env,
			Cluster: nodePatch.Cluster,
			Operate: v1alpha1.GalaxyNodeOperate{
				Type: v1alpha1.PatchGalaxyNode,
				Node: nodePatch.NodeName,
				OperateParams: v1alpha1.OperateParams{
					LabelOrAnnotation: v1alpha1.ParamMeta{
						Path: nodePatch.OperatorPath,
						Data: nodePatch.OperatorData,
					},
				},
			},
		},
	}
	err = n.galaxyNodeController.Create(&node)
	if err != nil {
		n.Error(gc, err)
		return
	}
	n.SuccessNoData(gc)
}

func (n *nodeService) TaintNode(gc *gin.Context) {
	jwtHeader := gc.GetHeader(constants.GwJWTHeader)
	user, err := sv1alpha1.GetUser(jwtHeader)
	if err != nil {
		n.ErrorWithMessage(gc, fmt.Sprintf("failed to parse params, error: %s", err))
		return

	}
	nodeTaint := PatchNodeTaint{}

	err = gc.ShouldBindJSON(&nodeTaint)
	if err != nil {
		n.Error(gc, err)
		return
	}

	node := v1alpha1.GalaxyNode{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("node-taint-%v", time.Now().Unix()),
		},
		Spec: v1alpha1.GalaxyNodeSpec{
			User:    fmt.Sprintf("%s(%s)", user.RealmName, user.NickNameCn),
			Cluster: nodeTaint.Cluster,
			Env:     nodeTaint.Env,
			Operate: v1alpha1.GalaxyNodeOperate{
				Type: v1alpha1.TaintGalaxyNode,
				Node: nodeTaint.NodeName,
				OperateParams: v1alpha1.OperateParams{
					Taints: nodeTaint.Taints,
				},
			},
		},
	}
	err = n.galaxyNodeController.Create(&node)
	if err != nil {
		n.Error(gc, err)
		return
	}
	n.SuccessNoData(gc)
}

func (n *nodeService) GetNodePods(gc *gin.Context) {
	cluster := gc.Param("cluster")
	nodeName := gc.Param("node")
	dsQuery := n.ParseDataSelectPathParameter(gc)

	nodeController, err := n.clientset.NewNodeController(cluster)
	if err != nil {
		n.Error(gc, err)
		return
	}

	pods, err := nodeController.GetPods(nodeName, dsQuery)
	if err != nil {
		n.Error(gc, err)
		return
	}
	n.Success(gc, pods)
}

func (n *NodeController) GetNodeEvents(gc *gin.Context) {
	cluster := gc.Param("cluster")
	node := gc.Param("node")
	dsQuery := page.ParseDataSelectPathParameter(gc)

	nodeController, err := n.clientset.NewNodeController(cluster)
	if err != nil {
		n.Error(gc, err)
		return
	}

	eventList, err := nodeController.GetEvents(node, dsQuery)
	if err != nil {
		n.Error(gc, err)
		return
	}
	n.Success(gc, eventList)
}

*/
