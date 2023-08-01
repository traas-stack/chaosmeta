package kube

import (
	"chaosmeta-platform/util/json"
	"context"

	"k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type HorizontalPodAutoscalerService interface {
	Get(namespace, name string) (*v2beta2.HorizontalPodAutoscaler, error)
	Create(hpa *v2beta2.HorizontalPodAutoscaler) error
	Update(hpa *v2beta2.HorizontalPodAutoscaler) error
	Patch(originalObj, updatedObj *v2beta2.HorizontalPodAutoscaler) error
	Replace(originalObj, updatedObj *v2beta2.HorizontalPodAutoscaler) error
}

type horizontalPodAutoscalerService struct {
	kubeClient kubernetes.Interface
}

func NewHorizontalPodAutoscalerService(kubeClient kubernetes.Interface) HorizontalPodAutoscalerService {
	return &horizontalPodAutoscalerService{kubeClient}
}

func (hpac *horizontalPodAutoscalerService) Get(namespace, name string) (*v2beta2.HorizontalPodAutoscaler, error) {
	return hpac.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (hpac *horizontalPodAutoscalerService) Create(hpa *v2beta2.HorizontalPodAutoscaler) error {
	_, err := hpac.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(hpa.GetNamespace()).Create(context.TODO(), hpa, metav1.CreateOptions{})
	return err
}

func (hpac *horizontalPodAutoscalerService) Update(hpa *v2beta2.HorizontalPodAutoscaler) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := hpac.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(hpa.GetNamespace()).Update(context.TODO(), hpa, metav1.UpdateOptions{})
		return err
	})
}

func (hpac *horizontalPodAutoscalerService) Patch(originalObj, updatedObj *v2beta2.HorizontalPodAutoscaler) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = hpac.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return err
}

func (hpac *horizontalPodAutoscalerService) Replace(originalObj, updatedObj *v2beta2.HorizontalPodAutoscaler) error {
	if originalObj == nil {
		return hpac.Create(updatedObj)
	}

	return hpac.Patch(originalObj, updatedObj)
}
