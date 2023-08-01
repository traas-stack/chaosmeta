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
	"k8s.io/client-go/util/retry"
	"sync"
)

// StatefulsetService defines the interface contains statefulset manages methods.
type StatefulsetService interface {
	Get(namespace, name string) (*StatefulSetDetail, error)
	Create(ss *appsv1.StatefulSet) error
	Update(ss *appsv1.StatefulSet) error
	Patch(originalObj, updatedObj *appsv1.StatefulSet) error
	Replace(originalObj, updatedObj *appsv1.StatefulSet) error
	List(namespace string, dsQuery *page.DataSelectQuery) (*StatefulSetResponse, error)
	GetRawPods(namespace, name string) ([]corev1.Pod, error)
	GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error)
	GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
}

type statefulsetService struct {
	param   *kubernetes.KubernetesParam
	podCtrl PodService
}

type StatefulSetCell appsv1.StatefulSet

func (n StatefulSetCell) GetProperty(name page.PropertyName) page.ComparableValue {
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

func (s *statefulsetService) toCells(std []appsv1.StatefulSet) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = StatefulSetCell(std[i])
	}
	return cells
}

func (s *statefulsetService) fromCells(cells []page.DataCell) []appsv1.StatefulSet {
	std := make([]appsv1.StatefulSet, len(cells))
	for i := range std {
		std[i] = appsv1.StatefulSet(cells[i].(StatefulSetCell))
	}
	return std
}

type StatefulSetResponse struct {
	Total    int                 `json:"total"`
	Current  int                 `json:"current"`
	PageSize int                 `json:"pageSize"`
	List     []StatefulSetDetail `json:"list"`
}

type StatefulSetDetail struct {
	appsv1.StatefulSet `json:",inline"`
	PodStatusInfo      common.PodStatusInfo     `json:"podStatusInfo"`
	ReplicaStatusInfo  common.ReplicaStatusInfo `json:"replicaStatusInfo"`
}

// NewStatefulsetService returns an instance of StatefulsetService.
func NewStatefulsetService(param *kubernetes.KubernetesParam) StatefulsetService {
	return &statefulsetService{param: param,
		podCtrl: NewPodService(param)}
}

func (s *statefulsetService) Get(namespace, name string) (*StatefulSetDetail, error) {
	st, err := s.param.KubernetesClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	detail := s.getStatefulSetPodStatus(st)
	return &detail, nil
}

func (s *statefulsetService) getStatefulSetPodStatus(statefulSet *appsv1.StatefulSet) StatefulSetDetail {
	var statefulSetDetailInfo StatefulSetDetail
	//matchingPods, err := s.GetRawPods(statefulSet.Namespace, statefulSet.Name)
	//if err != nil {
	//	return statefulSetDetailInfo
	//}
	////matchingEvents, err := s.GetEvents(statefulSet.Namespace, statefulSet.Name, page.NewDataSelectQuery(page.NewPaginationQuery(int(100), 0),
	////	page.NewSortQuery(strings.Split("", ",")),
	////	page.NewFilterQuery(strings.Split("", ","))))
	////if err != nil {
	////	return podStatusInfo
	////}
	//podInfo := common.GetPodInfo(statefulSet.Status.Replicas, *statefulSet.Spec.Replicas, matchingPods)
	//podInfo.Warnings = common.GetPodsEventWarnings(matchingEvents.List, matchingPods)
	statefulSetDetailInfo.StatefulSet = *statefulSet
	//statefulSetDetailInfo.PodStatusInfo = podInfo

	statefulSetDetailInfo.ReplicaStatusInfo.Desired = statefulSet.Status.Replicas
	statefulSetDetailInfo.ReplicaStatusInfo.Available = statefulSet.Status.ReadyReplicas

	return statefulSetDetailInfo
}

func (s *statefulsetService) List(namespace string, dsQuery *page.DataSelectQuery) (*StatefulSetResponse, error) {
	var statefulSetResponse StatefulSetResponse
	statefulSets, err := s.param.KubernetesClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var stsList []appsv1.StatefulSet
	for _, dd := range statefulSets.Items {
		stsList = append(stsList, dd)
	}

	statefulSetCells, filteredTotal := page.GenericDataSelectWithFilter(s.toCells(stsList), dsQuery)
	dps := s.fromCells(statefulSetCells)

	var statefulSetDetailList []StatefulSetDetail

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
		pch = make(chan StatefulSetDetail, lth)
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
			p := s.getStatefulSetPodStatus(&d)
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
		statefulSetDetailList = append(statefulSetDetailList, p)
	}
	statefulSetResponse.List = statefulSetDetailList
	statefulSetResponse.Current = dsQuery.PaginationQuery.Page + 1
	statefulSetResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	statefulSetResponse.Total = filteredTotal
	return &statefulSetResponse, nil
}

func FilterPodsByControllerRef(owner metav1.Object, allPods []corev1.Pod) []corev1.Pod {
	var matchingPods []corev1.Pod
	for _, pod := range allPods {
		if metav1.IsControlledBy(&pod, owner) {
			matchingPods = append(matchingPods, pod)
		}
	}
	return matchingPods
}

func (s *statefulsetService) GetRawPods(namespace, name string) ([]corev1.Pod, error) {
	statefulSet, err := s.param.KubernetesClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	rawPods, err := s.podCtrl.ListWithOptions(namespace, statefulSet.Spec.Selector)
	return FilterPodsByControllerRef(statefulSet, rawPods), nil
}

func (s *statefulsetService) GetPods(namespace, name string, dsQuery *page.DataSelectQuery) (*PodResponse, error) {
	var podResponse PodResponse

	podList, err := s.GetRawPods(namespace, name)
	if err != nil {
		return nil, err
	}

	podCells, filteredTotal := page.GenericDataSelectWithFilter(ToCells(podList), dsQuery)

	var podDetailList []PodDetail

	for _, po := range FromCells(podCells) {
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

func (s *statefulsetService) Create(ss *appsv1.StatefulSet) error {
	_, err := s.param.KubernetesClient.AppsV1().StatefulSets(ss.Namespace).Create(context.TODO(), ss, metav1.CreateOptions{})
	return err
}

func (s *statefulsetService) Update(ss *appsv1.StatefulSet) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := s.param.KubernetesClient.AppsV1().StatefulSets(ss.Namespace).Update(context.TODO(), ss, metav1.UpdateOptions{})
		return err
	})
}

func (s *statefulsetService) Patch(originalObj, updatedObj *appsv1.StatefulSet) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = s.param.KubernetesClient.AppsV1().StatefulSets(originalObj.Namespace).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)

	return err
}

func (s *statefulsetService) Replace(originalObj, updatedObj *appsv1.StatefulSet) error {
	if originalObj == nil {
		return s.Create(updatedObj)
	}

	return s.Patch(originalObj, updatedObj)
}

func (s *statefulsetService) GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	eventCtrl := NewEventService(s.param.KubernetesClient)
	eventResponse, err := eventCtrl.GetResourceEvents(namespace, name, dsQuery)
	return eventResponse, err
}
