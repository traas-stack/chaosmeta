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
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

// k8S-Target
func InitK8STarget(ctx context.Context, scope basic.Scope) error {
	var (
		K8SPodTarget        = basic.Target{Name: "pod", NameCn: "Pod", Description: "Fault injection capabilities related to cloud-native resource pod instances", DescriptionCn: "云原生资源pod实例相关的故障注入能力"}
		K8SDeploymentTarget = basic.Target{Name: "deployment", NameCn: "Deployment", Description: "Fault injection capabilities related to cloud-native resource deployment instances", DescriptionCn: "云原生资源deployment实例相关的故障注入能力"}
		K8SNodeTarget       = basic.Target{Name: "node", NameCn: "Node", Description: "Fault injection capabilities related to cloud-native resource node instances", DescriptionCn: "云原生资源node实例相关的故障注入能力"}
		K8SClusterTarget    = basic.Target{Name: "cluster", NameCn: "Cluster", Description: "Fault injection capabilities related to kubernetes macro cluster risks", DescriptionCn: "kubernetes宏观的集群性风险相关的故障注入能力"}
	)
	K8SPodTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &K8SPodTarget); err != nil {
		return err
	}
	if err := InitPodFault(ctx, K8SPodTarget); err != nil {
		return err
	}
	K8SDeploymentTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &K8SDeploymentTarget); err != nil {
		return err
	}
	if err := InitDeploymentFault(ctx, K8SDeploymentTarget); err != nil {
		return err
	}
	K8SNodeTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &K8SNodeTarget); err != nil {
		return err
	}
	if err := InitNodeFault(ctx, K8SNodeTarget); err != nil {
		return err
	}
	K8SClusterTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &K8SClusterTarget); err != nil {
		return err
	}
	return InitClusterFault(ctx, K8SClusterTarget)
}

// pod
func InitPodFault(ctx context.Context, podTarget basic.Target) error {
	var (
		podFaultDelete         = basic.Fault{TargetId: podTarget.ID, Name: "delete", NameCn: "删除Pod", Description: "Delete the target Pod instance", DescriptionCn: "删除目标Pod实例"}
		podFaultLabel          = basic.Fault{TargetId: podTarget.ID, Name: "label", NameCn: "增删Pod标签", Description: "Add or delete the label of the target Pod instance", DescriptionCn: "增删目标Pod实例的标签"}
		podFaultFinalizer      = basic.Fault{TargetId: podTarget.ID, Name: "finalizer", NameCn: "Pod增加finalizer", Description: "Add a finalizer to the target Pod instance", DescriptionCn: "为目标Pod实例增加finalizer"}
		podFaultContainerKill  = basic.Fault{TargetId: podTarget.ID, Name: "containerkill", NameCn: "杀掉Pod中的容器", Description: "Kill the specified container in the target Pod instance", DescriptionCn: "杀掉目标Pod实例中指定的容器"}
		podFaultContainerPause = basic.Fault{TargetId: podTarget.ID, Name: "containerpause", NameCn: "暂停Pod中的容器", Description: "Pauses the specified container in the target Pod instance", DescriptionCn: "暂停目标Pod实例中指定的容器"}
		podFaultContainerImage = basic.Fault{TargetId: podTarget.ID, Name: "containerimage", NameCn: "修改Pod容器镜像", Description: "Modify the image of the specified container in the target Pod instance", DescriptionCn: "修改目标Pod实例中指定容器的镜像"}
	)
	if err := basic.InsertFault(ctx, &podFaultDelete); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &podFaultLabel); err != nil {
		return err
	}
	if err := InitPodTargetArgsLabel(ctx, podFaultLabel); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &podFaultFinalizer); err != nil {
		return err
	}
	if err := InitPodTargetArgsFinalizer(ctx, podFaultFinalizer); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &podFaultContainerKill); err != nil {
		return err
	}
	if err := InitPodTargetArgsContainerKillAndPause(ctx, podFaultContainerKill); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &podFaultContainerPause); err != nil {
		return err
	}
	if err := InitPodTargetArgsContainerKillAndPause(ctx, podFaultContainerPause); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &podFaultContainerImage); err != nil {
		return err
	}
	return InitPodTargetArgsContainerImage(ctx, podFaultContainerImage)
}

func InitPodTargetArgsLabel(ctx context.Context, podFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的标签", ValueType: "string", Description: "Added labels;comma-separated key-value pair list", DescriptionCn: "增加的标签;逗号分隔的键值对列表"}
	argsDelete := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的标签key", ValueType: "string", Description: "Deleted label key; comma-separated key-value pair list", DescriptionCn: "删除的标签"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

func InitPodTargetArgsFinalizer(ctx context.Context, podFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的finalizer", ValueType: "string", Description: "Added finalizers; comma-separated key-value pair list", DescriptionCn: "增加的finalizer;逗号分隔的字符串列表"}
	argsDelete := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的标签finalizer", ValueType: "string", Description: "Deleted finalizer key; comma-separated key-value pair list", DescriptionCn: "删除的finalizer;逗号分隔的字符串列表"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

func InitPodTargetArgsContainerKillAndPause(ctx context.Context, podFault basic.Fault) error {
	argsContainerName := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "containername", KeyCn: "目标容器名称", ValueType: "string", DefaultValue: "", Description: "Target container name; specific container name, or 'firstcontainer' which represents the first container in the pod", DescriptionCn: "目标容器名称;具体的容器名称,或者“firstcontainer”,表示pod中第一个容器"}
	return basic.InsertArgs(ctx, &argsContainerName)
}

func InitPodTargetArgsContainerImage(ctx context.Context, podFault basic.Fault) error {
	argsContainerName := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "containername", KeyCn: "目标容器名称", ValueType: "string", Description: "Target container name; specific container name, or 'firstcontainer' which represents the first container in the pod", DescriptionCn: "目标容器名称;具体的容器名称,或者“firstcontainer”,表示pod中第一个容器"}
	argsImage := basic.Args{InjectId: podFault.ID, ExecType: ExecInject, Key: "image", KeyCn: "目标镜像名称", UnitCn: "", ValueType: "string", DefaultValue: "", Description: "Target image name", DescriptionCn: "目标镜像名称"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsContainerName, &argsImage})
}

// deployment
func InitDeploymentFault(ctx context.Context, deploymentTarget basic.Target) error {
	var (
		deploymentFaultDelete    = basic.Fault{TargetId: deploymentTarget.ID, Name: "delete", NameCn: "删除deployment", Description: "Delete target deployment", DescriptionCn: "删除目标deployment"}
		deploymentFaultLabel     = basic.Fault{TargetId: deploymentTarget.ID, Name: "label", NameCn: "修改deployment标签", Description: "Modify the label of the target deployment", DescriptionCn: "修改目标deployment的标签"}
		deploymentFaultFinalizer = basic.Fault{TargetId: deploymentTarget.ID, Name: "finalizer", NameCn: "增加deployment finalizer", Description: "Add a specified finalizer to the deployment instance so that its deletion is blocked or processed by the corresponding recycler.", DescriptionCn: "给deployment实例增加指定的finalizer，使之删除被阻塞或被对应回收器处理"}
		deploymentFaultReplicas  = basic.Fault{TargetId: deploymentTarget.ID, Name: "replicas", NameCn: "篡改deployment副本数", Description: "Tampering with the number of copies of the target deployment", DescriptionCn: "篡改目标deployment的副本数"}
	)
	if err := basic.InsertFault(ctx, &deploymentFaultDelete); err != nil {
		return err
	}
	if err := InitDeploymentDeleteArgs(ctx, deploymentFaultDelete); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &deploymentFaultLabel); err != nil {
		return err
	}
	if err := InitDeploymentLabelArgs(ctx, deploymentFaultLabel); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &deploymentFaultFinalizer); err != nil {
		return err
	}
	if err := InitDeploymentFinalizerArgs(ctx, deploymentFaultFinalizer); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &deploymentFaultReplicas); err != nil {
		return err
	}
	return InitDeploymentReplicasArgs(ctx, deploymentFaultReplicas)
}

func InitDeploymentDeleteArgs(ctx context.Context, deploymentFault basic.Fault) error {
	return nil
}

func InitDeploymentLabelArgs(ctx context.Context, deploymentFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的标签", ValueType: "string", Description: "Added labels; a comma-separated list of key-value pairs in the format: k1=v1,k2=v2", DescriptionCn: "增加的标签;逗号分隔的键值对列表,比如:k1=v1,k2=v2"}
	argsDelete := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的标签key", ValueType: "string", Description: "Deleted label; a comma-separated list of strings in the format: k1,k2", DescriptionCn: "删除的标签;逗号分隔的字符串列表,比如:k1,k2"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

func InitDeploymentFinalizerArgs(ctx context.Context, deploymentFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的finalizer", ValueType: "string", Description: "Added finalizer; a comma-separated list of finalizer names in the format: c/1,c/2", DescriptionCn: "增加的finalizer;逗号分隔的finalizer名称列表,格式为:c/1,c/2"}
	argsDelete := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的finalizer", ValueType: "string", DefaultValue: "", Description: "Removed finalizers; a comma-separated list of strings in the format: c/1,c/2", DescriptionCn: "删除的finalizer;逗号分隔的finalizer名称列表,格式为:c/1,c/2"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

func InitDeploymentReplicasArgs(ctx context.Context, deploymentFault basic.Fault) error {
	argsMode := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "扩缩模式", ValueType: "string", DefaultValue: "", Description: "Scaling mode", DescriptionCn: "扩缩容模式", ValueRule: "absolutecount,relativecount,relativepercent"}
	argsValue := basic.Args{InjectId: deploymentFault.ID, ExecType: ExecInject, Key: "value", KeyCn: "扩缩容大小", Description: "Numerical values, with different meanings in the three modes, absolutecount: the final target number of copies, relativecount: how much has been increased or decreased relative to the number of old copies, relativepercent: the percentage increase or decrease relative to the number of old copies", DescriptionCn: "扩缩容数值,在三种模式下表示不同含义.absolutecount:最终目标副本数;relativecount:相对旧副本数增加或减少了多少个,relativepercent:相对旧副本数增加或减少了百分之多少", ValueType: "string", DefaultValue: "", ValueRule: ""}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsMode, &argsValue})
}

// node
func InitNodeFault(ctx context.Context, nodeTarget basic.Target) error {
	var (
		nodeFaultLabel = basic.Fault{TargetId: nodeTarget.ID, Name: "label", NameCn: "修改node标签", Description: "The label of the node instance is dynamically modified", DescriptionCn: "node实例的label被动态修改"}
		nodeFaultTaint = basic.Fault{TargetId: nodeTarget.ID, Name: "taint", NameCn: "为node增加taint", Description: "Add specified stains to node instances to affect pod scheduling logic", DescriptionCn: "给node实例增加指定的污点,影响pod调度逻辑"}
	)
	if err := basic.InsertFault(ctx, &nodeFaultLabel); err != nil {
		return err
	}
	if err := InitNodeLabelArgs(ctx, nodeFaultLabel); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &nodeFaultTaint); err != nil {
		return err
	}
	return InitNodeTaintArgs(ctx, nodeFaultTaint)
}

func InitNodeLabelArgs(ctx context.Context, nodeFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: nodeFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的标签", ValueType: "string", Description: "Increased label;a comma-separated list of labels in the format: k1=v1,k2=v2", DescriptionCn: "增加的label,逗号分隔的label列表,格式为:k1=v1,k2=v2"}
	argsDelete := basic.Args{InjectId: nodeFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的标签", ValueType: "string", Description: "Removed label; a comma-separated list of labels in the format: k1,k2", DescriptionCn: "删除的label,逗号分隔的label列表,格式为:k1,k2"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

func InitNodeTaintArgs(ctx context.Context, nodeFault basic.Fault) error {
	argsAdd := basic.Args{InjectId: nodeFault.ID, ExecType: ExecInject, Key: "add", KeyCn: "增加的taint", ValueType: "string", DefaultValue: "", Description: "Increased taint； a comma-separated list of taints in the format: k1=v1:NoSchedule,k2=v2:NoSchedule", DescriptionCn: "增加的taint;逗号分隔的taint列表,格式为:k1=v1:NoSchedule,k2=v2:NoSchedule"}
	argsDelete := basic.Args{InjectId: nodeFault.ID, ExecType: ExecInject, Key: "delete", KeyCn: "删除的taint", ValueType: "string", DefaultValue: "", Description: "Removed taint; a comma-separated list of taints in the format: k1=v1:NoSchedule,k2=v2:NoSchedule", DescriptionCn: "删除的taint;逗号分隔的taint列表,格式为:k1=v1:NoSchedule,k2=v2:NoSchedule"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsAdd, &argsDelete})
}

//func GetSelectorArgs(selectorTypeName string, fault basic.Fault) []*basic.Args {
//	argsNamespace := basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "namespace", KeyCn: "命名空间", ValueType: "string", DefaultValue: "", Description: "removed taint; a comma-separated list of taints in the format: k1=v1:NoSchedule,k2=v2:NoSchedule", DescriptionCn: "删除的taint;逗号分隔的taint列表，格式为：k1=v1:NoSchedule,k2=v2:NoSchedule"}
//
//	return []*basic.Args{&NetworkArgsInterface, &NetworkArgsDstIP, &NetworkArgsSrcIP, &NetworkArgsDstPort, &NetworkArgsSrcPort, &NetworkArgsMode, &NetworkArgsForce}
//}

// cluster
func InitClusterFault(ctx context.Context, clusterTarget basic.Target) error {
	var (
		clusterFaultPendingPod   = basic.Fault{TargetId: clusterTarget.ID, Name: "pendingpod", NameCn: "堆积pending状态的pod", Description: "Accumulate a large number of pods in the pending state for the cluster in batches", DescriptionCn: "给集群批量堆积大量pending状态的pod"}
		clusterFaultCompletedJob = basic.Fault{TargetId: clusterTarget.ID, Name: "completedjob", NameCn: "堆积completed状态的job", Description: "Accumulate a large number of jobs in the completed state in batches for the cluster", DescriptionCn: "给集群批量堆积大量completed状态的job"}
	)
	if err := basic.InsertFault(ctx, &clusterFaultPendingPod); err != nil {
		return err
	}
	if err := InitClusterPendingPodArgs(ctx, clusterFaultPendingPod); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &clusterFaultCompletedJob); err != nil {
		return err
	}
	return InitClusterCompletedJobArgs(ctx, clusterFaultCompletedJob)
}

func InitClusterPendingPodArgs(ctx context.Context, clusterFault basic.Fault) error {
	return InitClusterCompletedJobArgs(ctx, clusterFault)
}

func InitClusterCompletedJobArgs(ctx context.Context, clusterFault basic.Fault) error {
	argsCount := basic.Args{InjectId: clusterFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "数量", ValueType: "int", DefaultValue: "", Required: true, Description: "Count", DescriptionCn: "数量", ValueRule: ">0"}
	argsNamespace := basic.Args{InjectId: clusterFault.ID, ExecType: ExecInject, Key: "namespace", KeyCn: "命名空间", ValueType: "string", DefaultValue: "", Required: true, Description: "A namespace that does not exist in the current cluster, such as: \"pendingattack\"", DescriptionCn: "当前集群不存在的namespace,比如:\"pendingattack\""}
	argsName := basic.Args{InjectId: clusterFault.ID, ExecType: ExecInject, Key: "name", KeyCn: "pod名称前缀", ValueType: "string", DefaultValue: "", Required: true, Description: "Pod name prefix", DescriptionCn: "pod名称前缀"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsCount, &argsNamespace, &argsName})
}
