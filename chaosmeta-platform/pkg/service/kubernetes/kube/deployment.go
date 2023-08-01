package kube

import (
	"chaosmeta-platform/pkg/models/common/page"
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/pkg/service/kubernetes/common"
	"chaosmeta-platform/util/json"
	"context"
	"fmt"
	"github.com/panjf2000/ants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

// DeploymentService defines the interface contains deployment manages methods.
type DeploymentService interface {
	List(namespace string, opts metav1.ListOptions, dsQuery *page.DataSelectQuery) (*DeploymentResponse, error)
	Get(namespace, name string) (*DeploymentDetail, error)
	Create(deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(namespace, name string) error
	Patch(originalObj, updatedObj *appsv1.Deployment) (*appsv1.Deployment, error)
	Replace(originalObj, updatedObj *appsv1.Deployment) (*appsv1.Deployment, error)
	GetRawPods(namespace, name string) ([]corev1.Pod, error)
	GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error)
	GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
}

type deploymentService struct {
	param   *kubernetes.KubernetesParam
	podCtrl PodService
}

// NewDeploymentService returns an instance of deployment Service.
func NewDeploymentService(
	param *kubernetes.KubernetesParam,
) DeploymentService {
	return &deploymentService{
		param:   param,
		podCtrl: NewPodService(param),
	}
}

type DeploymentResponse struct {
	Total    int                `json:"total"`
	Current  int                `json:"current"`
	PageSize int                `json:"pageSize"`
	List     []DeploymentDetail `json:"list"`
}

type DeploymentDetail struct {
	appsv1.Deployment `json:",inline"`
	PodStatusInfo     common.PodStatusInfo     `json:"podStatusInfo"`
	ReplicaStatusInfo common.ReplicaStatusInfo `json:"replicaStatusInfo"`
}

type DeploymentCell appsv1.Deployment

func (n DeploymentCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(n.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(n.ObjectMeta.CreationTimestamp.Time)
	case page.NamespaceProperty:
		return page.StdComparableString(n.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func (dp *deploymentService) toCells(std []appsv1.Deployment) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = DeploymentCell(std[i])
	}
	return cells
}

func (dp *deploymentService) fromCells(cells []page.DataCell) []appsv1.Deployment {
	std := make([]appsv1.Deployment, len(cells))
	for i := range std {
		std[i] = appsv1.Deployment(cells[i].(DeploymentCell))
	}
	return std
}

func (dp *deploymentService) getDeploymentPodStatus(deployment *appsv1.Deployment) DeploymentDetail {
	var dpDetail DeploymentDetail
	dpDetail.Deployment = *deployment
	dpDetail.ReplicaStatusInfo.Desired = deployment.Status.Replicas
	dpDetail.ReplicaStatusInfo.Available = deployment.Status.AvailableReplicas

	return dpDetail
}

func (dp *deploymentService) List(namespace string, opts metav1.ListOptions, dsQuery *page.DataSelectQuery) (*DeploymentResponse, error) {
	var deploymentResponse DeploymentResponse
	deployments, err := dp.param.KubernetesClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var deploymentList []appsv1.Deployment
	for _, dd := range deployments.Items {
		deploymentList = append(deploymentList, dd)
	}
	deploymentsCells, filteredTotal := page.GenericDataSelectWithFilter(dp.toCells(deploymentList), dsQuery)
	dps := dp.fromCells(deploymentsCells)

	var deploymentDetailList []DeploymentDetail

	defer func() {
		if err != nil {
			err = fmt.Errorf("error when list pods by names|%v", err)
			return
		}
	}()

	var (
		wg  sync.WaitGroup
		gp  *ants.Pool
		lth = len(dps)
		ech = make(chan error, lth)
		pch = make(chan DeploymentDetail, lth)
	)
	gp, err = ants.NewPool(20)
	if err != nil {
		err = fmt.Errorf("fail to new goroutine pool, caused by: %v", err)
		return nil, err
	}
	defer gp.Release()

	for _, d := range dps {
		wg.Add(1)
		d := d
		err = gp.Submit(func() {
			defer wg.Done()
			p := dp.getDeploymentPodStatus(&d)
			pch <- p
		})
		if err != nil {
			err = fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err)
			return nil, err
		}
	}
	wg.Wait()
	close(pch)
	select {
	case err = <-ech:
		return nil, err
	default:
		// do nothing
	}

	for p := range pch {
		deploymentDetailList = append(deploymentDetailList, p)
	}

	deploymentResponse.List = deploymentDetailList
	deploymentResponse.Current = dsQuery.PaginationQuery.Page + 1
	deploymentResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	deploymentResponse.Total = filteredTotal
	return &deploymentResponse, nil
}

func (dp *deploymentService) Get(namespace, name string) (*DeploymentDetail, error) {
	dep, err := dp.param.KubernetesClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	detail := dp.getDeploymentPodStatus(dep)
	return &detail, nil
}

func (dp *deploymentService) Create(dm *appsv1.Deployment) (*appsv1.Deployment, error) {
	return dp.param.KubernetesClient.AppsV1().Deployments(dm.Namespace).Create(context.TODO(), dm, metav1.CreateOptions{})
}

func (dp *deploymentService) Update(dm *appsv1.Deployment) (*appsv1.Deployment, error) {
	return dp.param.KubernetesClient.AppsV1().Deployments(dm.Namespace).Update(context.TODO(), dm, metav1.UpdateOptions{})

}

func (dp *deploymentService) Delete(namespace, name string) error {
	return dp.param.KubernetesClient.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (dp *deploymentService) Patch(originalObj, updatedObj *appsv1.Deployment) (*appsv1.Deployment, error) {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return nil, err
	}

	info, err := dp.param.KubernetesClient.AppsV1().Deployments(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return info, nil
}

func (dp *deploymentService) Replace(originalObj, updatedObj *appsv1.Deployment) (*appsv1.Deployment, error) {
	if originalObj == nil {
		return dp.Create(updatedObj)
	}

	return dp.Patch(originalObj, updatedObj)
}

func (dp *deploymentService) GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error) {
	var podResponse PodResponse

	podList, err := dp.GetRawPods(namespace, name)
	if err != nil {
		return nil, err
	}

	podCells, filteredTotal := page.GenericDataSelectWithFilter(ToCells(podList), dsQuery)
	ps := FromCells(podCells)

	var podDetailList []PodDetail

	for _, po := range ps {
		var detail PodDetail
		detail.Pod = po
		detail.PodPhase = getPodStatus(po)
		podDetailList = append(podDetailList, detail)
	}

	podResponse.List = podDetailList
	podResponse.Current = dsQuery.PaginationQuery.Page + 1
	podResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	podResponse.Total = filteredTotal

	return &podResponse, nil
}

func (dp *deploymentService) GetRawPods(namespace, name string) ([]corev1.Pod, error) {
	deployment, err := dp.param.KubernetesClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	rawPods, err := dp.podCtrl.ListWithOptions(namespace, deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}

	replicaList, err := dp.param.KubernetesClient.AppsV1().ReplicaSets(namespace).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	matchingPods := common.FilterDeploymentPodsByOwnerReference(*deployment, replicaList.Items, rawPods)

	return matchingPods, nil
}

func (dp *deploymentService) GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	_, err := dp.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	eventCtrl := NewEventService(dp.param.KubernetesClient)
	eventResponse, err := eventCtrl.GetResourceEvents(namespace, name, dsQuery)
	if err != nil {
		return nil, err
	}
	return eventResponse, nil
}
