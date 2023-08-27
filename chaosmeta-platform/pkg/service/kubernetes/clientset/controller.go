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

package clientset

import (
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/pkg/service/kubernetes/kube"
	"chaosmeta-platform/pkg/service/kubernetes/kubectl"
)

func (cs *clientset) initKubernetesParam(cluster string) (*kubernetes.KubernetesParam, error) {
	var param kubernetes.KubernetesParam
	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	config, err := cs.GetRestConfiguration(cluster)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, NotSupportedError
	}

	param.KubernetesClient = kubeClient
	param.RestConfig = config
	param.Cluster = cluster

	return &param, nil
}

func (cs *clientset) NewKubeService(cluster string) (kubectl.KubectlServiceInterface, error) {
	cfg, err := cs.GetRestConfiguration(cluster)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, NotSupportedError
	}

	kubeService, err := kubectl.NewkubectlService(cfg)
	if err != nil {
		return nil, err
	}
	cs.kubeServices.Store(cluster, kubeService)
	return kubeService, nil
}

func (cs *clientset) NewNodeService(cluster string) (kube.NodeService, error) {
	if nodeServiceInterface, ok := cs.nodeServices.Load(cluster); ok {
		return nodeServiceInterface.(kube.NodeService), nil
	}

	param, err := cs.initKubernetesParam(cluster)
	if err != nil {
		return nil, err
	}

	nodeService := kube.NewNodeService(param)
	cs.nodeServices.Store(cluster, nodeService)

	return nodeService, nil
}

func (cs *clientset) NewEventService(cluster string) (kube.EventService, error) {
	if eventServiceInterface, ok := cs.eventServices.Load(cluster); ok {
		return eventServiceInterface.(kube.EventService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	eventService := kube.NewEventService(kubeClient)
	cs.eventServices.Store(cluster, eventService)

	return eventService, nil
}

func (cs *clientset) NewPodService(cluster string) (kube.PodService, error) {
	param, err := cs.initKubernetesParam(cluster)
	if err != nil {
		return nil, err
	}

	if podServiceInterface, ok := cs.podServices.Load(cluster); ok {
		return podServiceInterface.(kube.PodService), nil
	}

	podService := kube.NewPodService(param)
	cs.podServices.Store(cluster, podService)

	return podService, nil
}

func (cs *clientset) NewReplicaSetService(cluster string) (kube.ReplicaSetService, error) {
	if replicaSetServiceInterface, ok := cs.replicaSetServices.Load(cluster); ok {
		return replicaSetServiceInterface.(kube.ReplicaSetService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	replicaSetService := kube.NewReplicaSetService(kubeClient)
	cs.replicaSetServices.Store(cluster, replicaSetService)

	return replicaSetService, nil
}

func (cs *clientset) NewDeploymentService(cluster string) (kube.DeploymentService, error) {
	if deploymentServiceInterface, ok := cs.deploymentServices.Load(cluster); ok {
		return deploymentServiceInterface.(kube.DeploymentService), nil
	}

	param, err := cs.initKubernetesParam(cluster)
	if err != nil {
		return nil, err
	}

	deploymentService := kube.NewDeploymentService(param)
	cs.deploymentServices.Store(cluster, deploymentService)

	return deploymentService, nil
}

func (cs *clientset) NewJobService(cluster string) (kube.JobService, error) {
	if jobServiceInterface, ok := cs.jobServices.Load(cluster); ok {
		return jobServiceInterface.(kube.JobService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	jobService := kube.NewJobService(kubeClient)
	cs.jobServices.Store(cluster, jobService)

	return jobService, nil
}

func (cs *clientset) NewServiceService(cluster string) (kube.ServiceService, error) {
	if serviceServiceInterface, ok := cs.serviceServices.Load(cluster); ok {
		return serviceServiceInterface.(kube.ServiceService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	serviceService := kube.NewServiceService(kubeClient)
	cs.serviceServices.Store(cluster, serviceService)

	return serviceService, nil
}

func (cs *clientset) NewIngressService(cluster string) (kube.IngressService, error) {
	if ingressServiceInterface, ok := cs.ingressServices.Load(cluster); ok {
		return ingressServiceInterface.(kube.IngressService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	ingressService := kube.NewIngressService(kubeClient)
	cs.ingressServices.Store(cluster, ingressService)

	return ingressService, nil
}

func (cs *clientset) NewNamespaceService(cluster string) (kube.NamespaceService, error) {
	if namespaceServiceInterface, ok := cs.namespaceServices.Load(cluster); ok {
		return namespaceServiceInterface.(kube.NamespaceService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	namespaceService := kube.NewNamespaceService(kubeClient)
	cs.namespaceServices.Store(cluster, namespaceService)

	return namespaceService, nil
}

func (cs *clientset) NewConfigMapService(cluster string) (kube.ConfigMapService, error) {
	if configmapServiceInterface, ok := cs.configMapServices.Load(cluster); ok {
		return configmapServiceInterface.(kube.ConfigMapService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	configmapService := kube.NewConfigMapService(kubeClient)
	cs.configMapServices.Store(cluster, configmapService)

	return configmapService, nil
}

func (cs *clientset) NewStatefulSetService(cluster string) (kube.StatefulsetService, error) {
	if statefulSetServiceInterface, ok := cs.statefulSetServices.Load(cluster); ok {
		return statefulSetServiceInterface.(kube.StatefulsetService), nil
	}

	param, err := cs.initKubernetesParam(cluster)
	if err != nil {
		return nil, err
	}

	statefulSetService := kube.NewStatefulsetService(param)
	cs.statefulSetServices.Store(cluster, statefulSetService)

	return statefulSetService, nil
}

func (cs *clientset) NewHPAService(cluster string) (kube.HorizontalPodAutoscalerService, error) {
	if hpaServicesInterface, ok := cs.hpaServices.Load(cluster); ok {
		return hpaServicesInterface.(kube.HorizontalPodAutoscalerService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	hpaService := kube.NewHorizontalPodAutoscalerService(kubeClient)
	cs.hpaServices.Store(cluster, hpaService)

	return hpaService, nil
}

func (cs *clientset) NewContainerService(cluster string) (kube.ContainerService, error) {
	if containerServicesInterface, ok := cs.containerServices.Load(cluster); ok {
		return containerServicesInterface.(kube.ContainerService), nil
	}

	kubeClient, err := cs.GetKubernetesClient(cluster)
	if err != nil {
		return nil, err
	}

	if kubeClient == nil {
		return nil, NotSupportedError
	}

	cfg, err := cs.GetRestConfiguration(cluster)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, NotSupportedError
	}

	containerService := kube.NewContainerService(kubeClient, cfg)
	cs.containerServices.Store(cluster, containerService)

	return containerService, nil
}

func (cs *clientset) NewDaemonSetService(cluster string) (kube.DaemonsetService, error) {
	if daemonSetCtrInterface, ok := cs.daemonSetService.Load(cluster); ok {
		return daemonSetCtrInterface.(kube.DaemonsetService), nil
	}

	param, err := cs.initKubernetesParam(cluster)
	if err != nil {
		return nil, err
	}

	daemonSetService := kube.NewDaemonSetService(param)
	cs.daemonSetService.Store(cluster, daemonSetService)

	return daemonSetService, nil
}

func (cs *clientset) NewCrdService(cluster string) (kube.CRDService, error) {
	if crdInterface, ok := cs.crdService.Load(cluster); ok {
		return crdInterface.(kube.CRDService), nil
	}

	crdClient, err := cs.GetCRDClient(cluster)
	if err != nil {
		return nil, err
	}

	return kube.NewCRDService(crdClient), nil
}
