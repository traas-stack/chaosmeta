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

package kube

import (
	"chaosmeta-platform/pkg/models/common/page"
	kubernetesService "chaosmeta-platform/pkg/service/kubernetes"
	"encoding/json"
	"flag"
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"testing"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func initKubeClient() (*kubernetes.Clientset, *rest.Config) {
	var kubeconfig *string

	fmt.Println(flag.Lookup("kubeconfig"))

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("k8sconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("k8sconfig", "", "absolute path to the kubeconfig file")
	}
	fmt.Println("--------------", *kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset, config
}

func Test_podService_List(t *testing.T) {
	kubeClient, config := initKubeClient()
	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	stopper := make(chan struct{})
	defer close(stopper)

	podInformer := factory.Core().V1().Pods()
	informer := podInformer.Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) {},
		UpdateFunc: func(interface{}, interface{}) {},
		DeleteFunc: func(interface{}) {},
	})

	factory.Start(stopper)
	factory.WaitForCacheSync(stopper)

	var param kubernetesService.KubernetesParam
	param.KubernetesClient = kubeClient
	param.RestConfig = config
	//param.Factory = factory

	podCtrl := NewPodService(&param)
	podDetail, err := podCtrl.Get("chaosmeta", "chaosmeta-inject-controller-manager-c8d699557-bw6x7")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*podDetail)
	podResp, err := podCtrl.List("chaosmeta", page.DefaultDataSelect)
	if err != nil {
		t.Fatal(err)
	}
	podRespByte, _ := json.Marshal(podResp)
	fmt.Println(string(podRespByte))

}
