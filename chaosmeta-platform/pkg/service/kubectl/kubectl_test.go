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
	_, kubeCtrl := NewkubectlService(config)
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
	_, kubeCtrl := NewkubectlService(config)
	content := ""
	err, getDeploymentList := kubeCtrl.GetByContent(nil, []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(getDeploymentList)
}
