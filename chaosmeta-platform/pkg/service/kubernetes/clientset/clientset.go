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
	"chaosmeta-platform/config"
	cv1alpha1 "chaosmeta-platform/pkg/gateway/apis/chaosmetacluster/v1alpha1"
	"chaosmeta-platform/pkg/models/cluster"
	"chaosmeta-platform/pkg/models/common/page"
	"chaosmeta-platform/pkg/service/kubernetes/clients"
	kube3 "chaosmeta-platform/pkg/service/kubernetes/kube"
	"chaosmeta-platform/util/enc_dec"
	"chaosmeta-platform/util/log"
	"context"
	goerr "errors"
	"fmt"
	"github.com/panjf2000/ants"
	"github.com/robfig/cron"
	"k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"strings"
	"sync"
)

const (
	Node        = "Node"
	Pod         = "Pod"
	Deployment  = "Deployment"
	DaemonSet   = "DaemonSet"
	StatefulSet = "StatefulSet"
)

type Interface interface {
	Run(ctx context.Context, checkHealthy bool)
	RunWorker()
	GetRestConfiguration(cluster string) (*rest.Config, error)
	GetKubernetesClient(cluster string) (kubernetes.Interface, error)
	GetCRDClient(cluster string) (apiextensionsclientset.Interface, error)
	GetKubernetesFactory(cluster string) (informers.SharedInformerFactory, error)
	NewNodeService(cluster string) (kube3.NodeService, error)
	NewPodService(cluster string) (kube3.PodService, error)
	NewEventService(cluster string) (kube3.EventService, error)
	NewReplicaSetService(cluster string) (kube3.ReplicaSetService, error)
	NewDeploymentService(cluster string) (kube3.DeploymentService, error)
	NewJobService(cluster string) (kube3.JobService, error)
	NewServiceService(cluster string) (kube3.ServiceService, error)
	NewIngressService(cluster string) (kube3.IngressService, error)
	NewConfigMapService(cluster string) (kube3.ConfigMapService, error)
	NewHPAService(cluster string) (kube3.HorizontalPodAutoscalerService, error)
	NewStatefulSetService(cluster string) (kube3.StatefulsetService, error)
	NewNamespaceService(cluster string) (kube3.NamespaceService, error)
	NewContainerService(cluster string) (kube3.ContainerService, error)
	NewDaemonSetService(cluster string) (kube3.DaemonsetService, error)
	NewCrdService(cluster string) (kube3.CRDService, error)

	ListCluster(env string, dsQuery *page.DataSelectQuery) (*ClusterListResponse, error)
	List(env string) (*cv1alpha1.ChaosmetaClusterList, error)
	GetCluster(env, cluster string) (*cv1alpha1.ChaosmetaCluster, error)
	CreateCluster(cluster *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error)
	DeleteCluster(env, cluster string) error
	PatchCluster(originalObj, updatedObj *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error)
	ReplaceCluster(originalObj, updatedObj *cv1alpha1.ChaosmetaCluster) (*cv1alpha1.ChaosmetaCluster, error)
	ListClusterDashboardInfo(env string) (dashboard []ClusterDashboardInfo, err error)
}

var NotSupportedError = goerr.New("Cluster Not Supported")
var NotRegistedError = goerr.New("Cluster Not Register")
var DefaultClientSet *clientset

type clientset struct {
	sync.RWMutex
	ServiceMutex sync.RWMutex

	ctx context.Context

	clusterSynced   cache.InformerSynced
	kubeLoadFile    string
	opClusterClient versioned.Interface
	//recorder            map[string]record.EventRecorder
	restConfig       sync.Map
	opClient         sync.Map
	crdClient        sync.Map
	kubernetesClient sync.Map
	prometheusClient sync.Map
	//prometheusClientSet sync.Map
	factory     sync.Map
	factoryStop sync.Map

	// Services
	clusterServices         sync.Map
	kubeServices            sync.Map
	nodeServices            sync.Map
	eventServices           sync.Map
	podServices             sync.Map
	replicaSetServices      sync.Map
	deploymentServices      sync.Map
	jobServices             sync.Map
	serviceServices         sync.Map
	ingressServices         sync.Map
	configMapServices       sync.Map
	statefulSetServices     sync.Map
	namespaceServices       sync.Map
	hpaServices             sync.Map
	containerServices       sync.Map
	validateResourceService sync.Map
	topologyService         sync.Map
	daemonSetService        sync.Map
	kubectlService          sync.Map
	injectionService        sync.Map
	crdService              sync.Map
	eventResource           sync.Map

	localCron *cron.Cron
	//goroutine pool
	goroutinePool *ants.Pool
	queue         workqueue.RateLimitingInterface
}

// NewClientset returns the clientset interface.
func NewClientset(opClusterClient versioned.Interface,
) Interface {

	cs := &clientset{
		opClusterClient: opClusterClient,
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cluster"),
	}
	cs.goroutinePool, _ = ants.NewPool(30)
	defer cs.goroutinePool.Release()

	DefaultClientSet = cs

	return cs
}

// worker runs a worker goroutine that invokes processNextWorkItem until the
// Service's queue is closed.
func (cs *clientset) RunWorker() {
	if err := cs.syncHandler(); err != nil {
		log.Error(err)
	}
}

func (cs *clientset) syncCluster() error {
	clusters, err := cluster.ListCluster()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, cluster := range clusters {
		kubeConf, err := enc_dec.Decrypt([]byte(cluster.KubeConfig), []byte(config.DefaultRunOptIns.SecretKey))
		if err != nil {
			log.Error(err)
			continue
		}
		clusterMeta := &ChaosmetaCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:              strings.ToLower(cluster.Name),
				CreationTimestamp: metav1.NewTime(cluster.CreateTime),
			},
			Spec: ChaosmetaClusterSpec{
				KubernetesOption: &KubernetesOption{
					LoadMode: "LoadFromAdmin",
					KubeConf: string(kubeConf),
				},
				Description: cluster.Version,
			},
		}
		clusterMeta.Status = ChaosmetaClusterStatus{
			Phase:                "Ready",
			LastUpdatedTimestamp: metav1.NewTime(cluster.UpdateTime),
		}

		if err := cs.addCluster(clusterMeta); err != nil {
			return err
		}
	}
	return nil
}

func (cs *clientset) syncHandler() error {
	if err := cs.syncCluster(); err != nil {
		log.Error(err)
		return err
	}
	localCron := cron.New()
	if err := localCron.AddFunc("@every 10m", func() {
		if err := cs.syncCluster(); err != nil {
			log.Error(err)
		}
	}); err != nil {
		return err
	}

	localCron.Start()
	cs.localCron = localCron
	return nil
}

func (cs *clientset) RunInformer(factory informers.SharedInformerFactory, resource string) {
	var indexedInformer cache.SharedIndexInformer
	switch resource {
	case Pod:
		indexedInformer = factory.Core().V1().Pods().Informer()
	case Node:
		indexedInformer = factory.Core().V1().Nodes().Informer()
	case Deployment:
		indexedInformer = factory.Apps().V1().Deployments().Informer()
	case StatefulSet:
		indexedInformer = factory.Apps().V1().StatefulSets().Informer()
	case DaemonSet:
		indexedInformer = factory.Apps().V1().DaemonSets().Informer()
	}

	indexedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) {},
		UpdateFunc: func(interface{}, interface{}) {},
		DeleteFunc: func(interface{}) {},
	})
}

func (cs *clientset) InitFactory(factory informers.SharedInformerFactory, name string) {
	stopper := make(chan struct{})
	defer close(stopper)

	cs.factoryStop.Store(name, stopper)
	defer runtime.HandleCrash()

	resourceList := []string{Pod, StatefulSet, Node, Deployment, DaemonSet}

	for _, tmp := range resourceList {
		cs.RunInformer(factory, tmp)
	}

	factory.Start(stopper)
	factory.WaitForCacheSync(stopper)

	<-stopper
}

func (cs *clientset) addCluster(cluster *ChaosmetaCluster) error {
	// parsing kubernetes configuration
	if err := func() error {
		if cluster.Spec.KubernetesOption == nil || (cluster.Spec.KubernetesOption.LoadMode == "" && cluster.Spec.KubernetesOption.KubeConf == "") {
			log.Errorf("cluster spec: kubernetes option or options.kubeconf is empty")
			cs.restConfig.Store(cluster.GetName(), nil)
			cs.kubernetesClient.Store(cluster.GetName(), nil)
			cs.opClient.Store(cluster.GetName(), nil)
			cs.crdClient.Store(cluster.GetName(), nil)
			cs.factory.Store(cluster.GetName(), nil)
			cs.clusterServices.Store(cluster.GetName(), nil)
			cs.kubeServices.Store(cluster.GetName(), nil)
			return nil
		}

		kubeOpt := cluster.Spec.KubernetesOption
		restClient, err := GetKubeRestConf(KubeLoadMode(kubeOpt.LoadMode), kubeOpt.KubeConf)
		if err != nil {
			return err
		}

		kubeClient, err := kubernetes.NewForConfig(restClient)
		if err != nil {
			return err
		}

		//factory := informers.NewSharedInformerFactory(kubeClient, 0)

		opClient, err := clients.NewForConfig(restClient)
		if err != nil {
			return err
		}

		crdClient, err := apiextensionsclientset.NewForConfig(restClient)
		if err != nil {
			return err
		}

		cs.restConfig.Store(cluster.GetName(), restClient)
		cs.kubernetesClient.Store(cluster.GetName(), kubeClient)
		cs.opClient.Store(cluster.GetName(), opClient)
		cs.crdClient.Store(cluster.GetName(), crdClient)
		cs.clusterServices.Store(cluster.GetName(), cluster)

		kubeCtl, err := cs.NewKubeService(cluster.GetName())
		if err != nil {
			return err
		}
		cs.kubeServices.Store(cluster.GetName(), kubeCtl)
		return nil
	}(); err != nil {
		return fmt.Errorf("KubernetesOption Error: %s", err)
	}

	return nil
}

func (cs *clientset) ListClusters(env string) ([]cv1alpha1.ChaosmetaCluster, error) {
	var clusterList []cv1alpha1.ChaosmetaCluster
	cs.clusterServices.Range(func(key, value interface{}) bool {
		cluster := value.(*cv1alpha1.ChaosmetaCluster)
		if cluster != nil {
			if env != "" {
				if strings.Contains(cluster.GetName(), strings.ToLower(env)) || strings.Contains(cluster.GetName(), strings.ToUpper(env)) {
					clusterList = append(clusterList, *cluster)
				}
			} else {
				clusterList = append(clusterList, *cluster)
			}
			log.Info("ListClusters name:", cluster.GetName(), " ", "cluster:", cluster)
		}
		return true
	})
	return clusterList, nil
}

func (cs *clientset) onClusterDelete(obj interface{}) {
	// todo 禁止删除默认集群 否则创建crd会出现问题 待修复
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorf("failed to delete key from obj %#v, error: %w", obj, err)
		return
	}

	cs.queue.AddRateLimited(key)

	cluster, ok := obj.(*ChaosmetaCluster)
	if !ok {
		return
	}
	if factoryStopInterface, _ := cs.factoryStop.Load(cluster.GetName()); factoryStopInterface != nil {
		factoryStopInterface.(chan struct{}) <- struct{}{}
		cs.factoryStop.Store(cluster.GetName(), factoryStopInterface.(chan struct{}))
	}

	//cs.recorder.Delete(cluster.GetName())
	cs.restConfig.Delete(cluster.GetName())
	cs.kubernetesClient.Delete(cluster.GetName())
	cs.opClient.Delete(cluster.GetName())
	cs.crdClient.Delete(cluster.GetName())
	cs.prometheusClient.Delete(cluster.GetName())
	cs.nodeServices.Delete(cluster.GetName())
	cs.kubeServices.Delete(cluster.GetName())
	cs.clusterServices.Delete(cluster.GetName())
	cs.podServices.Delete(cluster.GetName())
	cs.replicaSetServices.Delete(cluster.GetName())
	cs.deploymentServices.Delete(cluster.GetName())
	cs.jobServices.Delete(cluster.GetName())
	cs.configMapServices.Delete(cluster.GetName())
	cs.hpaServices.Delete(cluster.GetName())
	cs.ingressServices.Delete(cluster.GetName())
	cs.serviceServices.Delete(cluster.GetName())
	cs.statefulSetServices.Delete(cluster.GetName())
	cs.namespaceServices.Delete(cluster.GetName())
	cs.containerServices.Delete(cluster.GetName())
	cs.topologyService.Delete(cluster.GetName())
	cs.validateResourceService.Delete(cluster.GetName())
	cs.daemonSetService.Delete(cluster.GetName())
	cs.kubectlService.Delete(cluster.GetName())
	cs.crdService.Delete(cluster.GetName())
}

func (cs *clientset) checkCluster(cluster *cv1alpha1.ChaosmetaCluster) {
	var err error
	updatedObj := cluster.DeepCopy()

	// check kubernetes healthy
	if err = func() error {
		kubeClient, err := cs.GetKubernetesClient(cluster.GetName())
		if err != nil {
			if err == NotSupportedError {
				return nil
			}

			return err
		}

		if _, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{Limit: 1}); err != nil {
			return fmt.Errorf("failed to list node in cluster, error: %s", err)
		}

		return nil
	}(); err != nil {
		goto Error
	}

	if updatedObj.Status.Phase != cv1alpha1.ReadyChaosmetaClusterPhase && !strings.Contains(updatedObj.Status.Reason, "EdasOption Error") {
		updatedObj.Status.Phase = cv1alpha1.ReadyChaosmetaClusterPhase
		updatedObj.Status.LastUpdatedTimestamp = metav1.Now()
		updatedObj.Status.Reason = ""

		if _, err := cs.PatchCluster(cluster, updatedObj); err != nil {
			log.Errorf("failed to update cluster %s status, error: %s", cluster.GetName(), err)
		}
	}

	return

Error:
	updatedObj.Status.Phase = cv1alpha1.FailedChaosmetaClusterPhase
	updatedObj.Status.LastUpdatedTimestamp = metav1.Now()
	updatedObj.Status.Reason = err.Error()

	_, err = cs.PatchCluster(cluster, updatedObj)
	if err != nil {
		log.Errorf("failed to update cluster %s status, error: %s", cluster.GetName(), err)
	}
}

func (cs *clientset) Run(ctx context.Context, checkHealthy bool) {
}

func (cs *clientset) GetRestConfiguration(cluster string) (*rest.Config, error) {
	cfgInterface, ok := cs.restConfig.Load(cluster)
	if !ok {
		log.Errorf("无法获取集群%s的restconfig", cluster)
		return nil, NotRegistedError
	}
	if cfgInterface == nil {
		return nil, NotSupportedError
	}

	cfg := cfgInterface.(*rest.Config)
	if cfg == nil {
		return nil, NotSupportedError
	}

	return cfg, nil
}

func (cs *clientset) ListRestConfiguration() ([]*rest.Config, error) {
	var restConfigList []*rest.Config
	cs.restConfig.Range(func(key, value interface{}) bool {
		restConfig := value.(*rest.Config)
		if restConfig != nil {
			restConfigList = append(restConfigList, restConfig)
			log.Info("ListRestConfiguration cluster:", key, " ", "restConfig:", value)
		}
		return true
	})
	return restConfigList, nil
}

func (cs *clientset) GetKubernetesClient(cluster string) (kubernetes.Interface, error) {
	clientInterface, ok := cs.kubernetesClient.Load(cluster)
	if !ok {
		log.Errorf("无法获取集群%s的kubeClient", cluster)
		return nil, NotRegistedError
	}
	if clientInterface == nil {
		return nil, NotSupportedError
	}

	client := clientInterface.(kubernetes.Interface)
	if client == nil {
		return nil, NotSupportedError
	}

	return client, nil
}

func (cs *clientset) GetCRDClient(cluster string) (apiextensionsclientset.Interface, error) {
	clientInterface, ok := cs.crdClient.Load(cluster)
	if !ok {
		log.Errorf("无法获取集群%s的crdClient", cluster)
		return nil, NotRegistedError
	}
	if clientInterface == nil {
		return nil, NotSupportedError
	}

	client := clientInterface.(apiextensionsclientset.Interface)
	if client == nil {
		return nil, NotSupportedError
	}

	return client, nil
}

func (cs *clientset) GetKubernetesFactory(cluster string) (informers.SharedInformerFactory, error) {
	clientInterface, ok := cs.factory.Load(cluster)
	if !ok {
		log.Errorf("无法获取集群%s的KubernetesFactory", cluster)
		return nil, NotRegistedError
	}
	if clientInterface == nil {
		return nil, NotSupportedError
	}

	client := clientInterface.(informers.SharedInformerFactory)
	if client == nil {
		return nil, NotSupportedError
	}

	return client, nil
}
