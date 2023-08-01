package kube

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ReplicaSetService interface {
	List(namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error)
	Get(namespace, name string) (*appsv1.ReplicaSet, error)
}

type replicaSetService struct {
	kubeClient kubernetes.Interface
}

func NewReplicaSetService(kubeClient kubernetes.Interface) ReplicaSetService {
	return &replicaSetService{
		kubeClient: kubeClient,
	}
}

func (r *replicaSetService) List(namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error) {
	return r.kubeClient.AppsV1().ReplicaSets(namespace).List(context.TODO(), opts)
}

func (r *replicaSetService) Get(namespace, name string) (*appsv1.ReplicaSet, error) {
	return r.kubeClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
