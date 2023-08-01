package common

import (
	corev1 "k8s.io/api/core/v1"
)

type PodStatusInfo struct {
	// Number of pods that are created.
	Current int32 `json:"current"`

	// Number of pods that are desired.
	Desired int32 `json:"desired,omitempty"`

	// Number of pods that are currently running.
	Running int32 `json:"running"`

	// Number of pods that are currently waiting.
	Pending int32 `json:"pending"`

	// Number of pods that are failed.
	Failed int32 `json:"failed"`

	// Number of pods that are succeeded.
	Succeeded int32 `json:"succeeded"`

	// Unique warning messages related to pods in this resource.
	Warnings []corev1.Event `json:"warnings"`
}

type ReplicaStatusInfo struct {
	Desired   int32 `json:"desired,omitempty"`
	Available int32 `json:"available"`
}
