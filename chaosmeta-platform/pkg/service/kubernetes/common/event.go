package common

import (
	"strings"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"
)

// FailedReasonPartials  is an array of partial strings to correctly filter warning events.
// Have to be lower case for correct case insensitive comparison.
// Based on k8s official events reason file:
// https://github.com/kubernetes/blob/886e04f1fffbb04faf8a9f9ee141143b2684ae68/pkg/kubelet/events/event.go
// Partial strings that are not in k8s_event.go file are added in order to support
// older versions of k8s which contained additional event reason messages.
var FailedReasonPartials = []string{"failed", "err", "exceeded", "invalid", "unhealthy",
	"mismatch", "insufficient", "conflict", "outof", "nil", "backoff"}

// GetPodsEventWarnings returns warning pod events by filtering out events targeting only given pods
func GetPodsEventWarnings(events []corev1.Event, pods []corev1.Pod) []corev1.Event {
	result := make([]corev1.Event, 0)

	// Filter out only warning events
	events = getWarningEvents(events)
	failedPods := make([]corev1.Pod, 0)

	// Filter out ready and successful pods
	for _, pod := range pods {
		if !isReadyOrSucceeded(pod) {
			failedPods = append(failedPods, pod)
		}
	}

	// Filter events by failed pods UID
	events = filterEventsByPodsUID(events, failedPods)
	events = removeDuplicates(events)

	for _, event := range events {
		result = append(result, corev1.Event{
			Message: event.Message,
			Reason:  event.Reason,
			Type:    event.Type,
		})
	}

	return result
}

// Returns filtered list of event objects. Events list is filtered to get only events targeting
// pods on the list.
func filterEventsByPodsUID(events []corev1.Event, pods []corev1.Pod) []corev1.Event {
	result := make([]corev1.Event, 0)
	podEventMap := make(map[types.UID]bool, 0)

	if len(pods) == 0 || len(events) == 0 {
		return result
	}

	for _, pod := range pods {
		podEventMap[pod.UID] = true
	}

	for _, event := range events {
		if _, exists := podEventMap[event.InvolvedObject.UID]; exists {
			result = append(result, event)
		}
	}

	return result
}

func FillEventsType(events []corev1.Event) []corev1.Event {
	for i := range events {
		// Fill in only events with empty type.
		if len(events[i].Type) == 0 {
			if isFailedReason(events[i].Reason, FailedReasonPartials...) {
				events[i].Type = corev1.EventTypeWarning
			} else {
				events[i].Type = corev1.EventTypeNormal
			}
		}
	}

	return events
}

// Returns filtered list of event objects.
// Event list object is filtered to get only warning events.
func getWarningEvents(events []corev1.Event) []corev1.Event {
	return filterEventsByType(FillEventsType(events), corev1.EventTypeWarning)
}

// Filters kubernetes corev1 event objects based on event type.
// Empty string will return all events.
func filterEventsByType(events []corev1.Event, eventType string) []corev1.Event {
	if len(eventType) == 0 || len(events) == 0 {
		return events
	}

	result := make([]corev1.Event, 0)
	for _, event := range events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}

	return result
}

// Returns true if reason string contains any partial string indicating that this may be a
// warning, false otherwise
func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}

	return false
}

// Removes duplicate strings from the slice
func removeDuplicates(slice []corev1.Event) []corev1.Event {
	visited := make(map[string]bool, 0)
	result := make([]corev1.Event, 0)

	for _, elem := range slice {
		if !visited[elem.Reason] {
			visited[elem.Reason] = true
			result = append(result, elem)
		}
	}

	return result
}

// Returns true if given pod is in state ready or succeeded, false otherwise
func isReadyOrSucceeded(pod corev1.Pod) bool {
	if pod.Status.Phase == corev1.PodSucceeded {
		return true
	}
	if pod.Status.Phase == corev1.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == corev1.PodReady {
				if c.Status == corev1.ConditionFalse {
					return false
				}
			}
		}

		return true
	}

	return false
}
