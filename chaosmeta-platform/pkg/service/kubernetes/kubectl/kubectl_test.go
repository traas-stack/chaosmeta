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

package kubectl

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func TestApplyController_ApplyByContent(t *testing.T) {
	_, config := initKubeClient()
	kubeCtrl, _ := NewkubectlService(config)
	content := ""
	err, ctls := kubeCtrl.ApplyByContent(nil, []byte(content), true)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(ctls)
}

func TestGetDeployment(t *testing.T) {
	_, config := initKubeClient()
	kubeCtrl, _ := NewkubectlService(config)
	content := ""
	err, getDeploymentList := kubeCtrl.GetByContent(nil, []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(getDeploymentList)
}
