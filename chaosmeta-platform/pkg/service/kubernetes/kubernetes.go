package kubernetes

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesParam struct {
	Cluster          string
	RestConfig       *rest.Config
	KubernetesClient kubernetes.Interface
	Factory          informers.SharedInformerFactory
}
