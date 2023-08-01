package kube

import (
	"chaosmeta-platform/util/json"
	"context"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// JobServicedefines the interface contains job manages methods.
type JobService interface {
	List(namespace string, opts metav1.ListOptions) (*batchv1.JobList, error)
	Get(namespace, name string) (*batchv1.Job, error)
	Create(job *batchv1.Job) error
	Update(job *batchv1.Job) error
	Delete(namespace, name string) error
	Patch(originalObj, updatedObj *batchv1.Job) error
	Replace(originalObj, updatedObj *batchv1.Job) error
}

type jobService struct {
	kubeClient kubernetes.Interface
}

func NewJobService(kubeClient kubernetes.Interface) JobService {
	return &jobService{
		kubeClient: kubeClient,
	}
}

func (jc *jobService) List(namespace string, opts metav1.ListOptions) (*batchv1.JobList, error) {
	return jc.kubeClient.BatchV1().Jobs(namespace).List(context.TODO(), opts)
}

func (jc *jobService) Get(namespace, name string) (*batchv1.Job, error) {
	return jc.kubeClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (jc *jobService) Create(job *batchv1.Job) error {
	_, err := jc.kubeClient.BatchV1().Jobs(job.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	return err
}

func (jc *jobService) Delete(namespace, name string) error {
	return jc.kubeClient.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (jc *jobService) Update(job *batchv1.Job) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := jc.kubeClient.BatchV1().Jobs(job.Namespace).Update(context.TODO(), job, metav1.UpdateOptions{})
		return err
	})
}

func (jc *jobService) Patch(originalObj, updatedObj *batchv1.Job) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = jc.kubeClient.BatchV1().Jobs(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return err
}

func (jc *jobService) Replace(originalObj, updatedObj *batchv1.Job) error {
	if originalObj == nil {
		return jc.Create(updatedObj)
	}

	return jc.Patch(originalObj, updatedObj)
}
