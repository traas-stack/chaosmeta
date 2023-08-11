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
	"chaosmeta-platform/pkg/service/kubernetes/common"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type EventService interface {
	List(namespace string, opt metav1.ListOptions) (*v1.EventList, error)
	Get(namespace, name string) (*v1.Event, error)
	GetResourceEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
}

var FailedReasonPartials = []string{"failed", "err", "exceeded", "invalid", "unhealthy",
	"mismatch", "insufficient", "conflict", "outof", "nil", "backoff"}

type eventService struct {
	kubeClient kubernetes.Interface
}

type EventResponse struct {
	Total    int           `json:"total"`
	Current  int           `json:"current"`
	PageSize int           `json:"pageSize"`
	List     []EventDetail `json:"list"`
}

type EventDetail struct {
	v1.Event `json:",inline"`
	Age      string `json:"age"`
	Source   string `json:"source"`
}

type EventCell v1.Event

func (e EventCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(e.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(e.ObjectMeta.CreationTimestamp.Time)
	case page.NamespaceProperty:
		return page.StdComparableString(e.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func EventToCells(std []v1.Event) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = EventCell(std[i])
	}
	return cells
}

func EventFromCells(cells []page.DataCell) []v1.Event {
	std := make([]v1.Event, len(cells))
	for i := range std {
		std[i] = v1.Event(cells[i].(EventCell))
	}
	return std
}

func NewEventService(kubeClient kubernetes.Interface) EventService {
	return &eventService{kubeClient}
}

func (ec *eventService) List(namespace string, opt metav1.ListOptions) (*v1.EventList, error) {
	return ec.kubeClient.CoreV1().Events(namespace).List(context.TODO(), opt)
}

func (ec *eventService) Get(namespace, name string) (*v1.Event, error) {
	return ec.kubeClient.CoreV1().Events(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func toEvent(event v1.Event) EventDetail {
	var detail EventDetail
	detail.Event = event
	// 持续时间
	delta := time.Now().Sub(detail.FirstTimestamp.Time).Hours()
	// 最近一次时间
	currentDelta := time.Now().Sub(detail.LastTimestamp.Time).Truncate(time.Second).String()
	detail.Age = fmt.Sprintf("%v (*%v over %vh)", currentDelta, detail.Count, fmt.Sprintf("%.f", delta))

	detail.Source = fmt.Sprintf("%s", event.Source.Component)
	return detail
}

func (ec *eventService) GetResourceEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	resourceEvents, err := GetEvents(ec.kubeClient, namespace, name)
	if err != nil {
		return nil, err
	}

	var eventResponse EventResponse
	eventList := resourceEvents
	eventCells, filteredTotal := page.GenericDataSelectWithFilter(EventToCells(eventList), dsQuery)
	dps := EventFromCells(eventCells)

	var eventDetailList []EventDetail
	for _, tmp := range dps {
		eventDetailList = append(eventDetailList, toEvent(tmp))
	}

	eventResponse.List = eventDetailList
	eventResponse.Current = dsQuery.PaginationQuery.Page + 1
	eventResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	eventResponse.Total = filteredTotal
	return &eventResponse, nil
}

func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}

	return false
}

func FillEventsType(events []v1.Event) []v1.Event {
	for i := range events {
		// Fill in only events with empty type.
		if len(events[i].Type) == 0 {
			if isFailedReason(events[i].Reason, FailedReasonPartials...) {
				events[i].Type = v1.EventTypeWarning
			} else {
				events[i].Type = v1.EventTypeNormal
			}
		}
	}

	return events
}

func GetEvents(client kubernetes.Interface, namespace, resourceName string) ([]v1.Event, error) {
	fieldSelector, err := fields.ParseSelector("involvedObject.name" + "=" + resourceName)

	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		EventList: common.GetEventListChannelWithOptions(
			client,
			common.NewSameNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}

	eventList := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	return FillEventsType(eventList.Items), nil
}
