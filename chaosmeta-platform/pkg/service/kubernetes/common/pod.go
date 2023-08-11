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

package common

import (
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FilterDeploymentPodsByOwnerReference(deployment apps.Deployment, allRS []apps.ReplicaSet,
	allPods []corev1.Pod) []corev1.Pod {
	var matchingPods []corev1.Pod
	for _, rs := range allRS {
		if metav1.IsControlledBy(&rs, &deployment) {
			matchingPods = append(matchingPods, FilterPodsByControllerRef(&rs, allPods)...)
		}
	}

	return matchingPods
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

// GetPodInfo returns aggregate information about a group of pods.
func GetPodInfo(current int32, desired int32, pods []corev1.Pod) PodStatusInfo {
	result := PodStatusInfo{
		Current:  current,
		Desired:  desired,
		Warnings: make([]corev1.Event, 0),
	}

	for _, pod := range pods {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			result.Running++
		case corev1.PodPending:
			result.Pending++
		case corev1.PodFailed:
			result.Failed++
		case corev1.PodSucceeded:
			result.Succeeded++
		}
	}

	return result
}
