/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strconv"
	"strings"
	"time"
)

// log is for logging in this package.
var (
	experimentlog = logf.Log.WithName("experiment-resource")
)

func (r *Experiment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-inject-chaosmeta-io-v1alpha1-experiment,mutating=true,failurePolicy=fail,sideEffects=None,groups=inject.chaosmeta.io,resources=experiments,verbs=create,versions=v1alpha1,name=mexperiment.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Experiment{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Experiment) Default() {
	experimentlog.Info("mutate", "name", r.Name)

	if r.Status.Phase != "" {
		return
	}

	var i int
	for i = 0; i < len(r.ObjectMeta.Finalizers); i++ {
		if r.ObjectMeta.Finalizers[i] == FinalizerName {
			break
		}
	}

	if i == len(r.ObjectMeta.Finalizers) {
		r.ObjectMeta.Finalizers = append(r.ObjectMeta.Finalizers, FinalizerName)
	}

	if r.Spec.Scope == PodScopeType || (r.Spec.Scope == KubernetesScopeType && strings.Index(r.Spec.Experiment.Target, "container") >= 0) {
		var i int
		for i = 0; i < len(r.Spec.Experiment.Args); i++ {
			if r.Spec.Experiment.Args[i].Key == ContainerKey {
				break
			}
		}

		if i == len(r.Spec.Experiment.Args) {
			r.Spec.Experiment.Args = append(r.Spec.Experiment.Args, ArgsUnit{
				Key:       ContainerKey,
				Value:     FirstContainer,
				ValueType: StringVType,
			})
		}
	}
}

//+kubebuilder:webhook:path=/validate-inject-chaosmeta-io-v1alpha1-experiment,mutating=false,failurePolicy=fail,sideEffects=None,groups=inject.chaosmeta.io,resources=experiments,verbs=create;update;delete,versions=v1alpha1,name=vexperiment.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Experiment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Experiment) ValidateCreate() error {
	experimentlog.Info("validate create", "name", r.Name)
	if r.Spec.Experiment.Duration == "" {
		return fmt.Errorf("experiment's duration is empty")
	}

	_, err := ConvertDuration(r.Spec.Experiment.Duration)
	if err != nil {
		return fmt.Errorf("experiment's duration is invalid: %s", err.Error())
	}

	if r.Spec.Scope != PodScopeType && r.Spec.Scope != NodeScopeType && r.Spec.Scope != KubernetesScopeType {
		return fmt.Errorf("\"scope\" not support: %s, only support: %s, %s, %s", r.Spec.Scope, PodScopeType, NodeScopeType, KubernetesScopeType)
	}

	if r.Spec.TargetPhase != InjectPhaseType {
		return fmt.Errorf("initial \"targetPhase\" only support: %s", InjectPhaseType)
	}

	if r.Spec.RangeMode != nil {
		if r.Spec.RangeMode.Type != AllRangeType && r.Spec.RangeMode.Type != PercentRangeType && r.Spec.RangeMode.Type != CountRangeType {
			return fmt.Errorf("\"rangeMode.type\" not support: %s, only support: %s, %s, %s", r.Spec.RangeMode.Type, AllRangeType, PercentRangeType, CountRangeType)
		}

		if r.Spec.RangeMode.Type == PercentRangeType {
			if r.Spec.RangeMode.Value <= 0 || r.Spec.RangeMode.Value > 100 {
				return fmt.Errorf("\"rangeMode.value\" should be in (0,100]")
			}
		}

		if r.Spec.RangeMode.Type == CountRangeType {
			if r.Spec.RangeMode.Value <= 0 {
				return fmt.Errorf("\"rangeMode.value\" should larger than 0")
			}
		}
	}

	if len(r.Spec.Selector) == 0 && r.Spec.Scope != KubernetesScopeType {
		return fmt.Errorf("length of \"selector\" must not be 0")
	}

	if r.Spec.Scope == PodScopeType {
		for _, unitSelector := range r.Spec.Selector {
			if unitSelector.Namespace == "" {
				return fmt.Errorf("namespace in selector must not empty")
			}
		}
	} else if r.Spec.Scope == NodeScopeType {
		for _, unitSelector := range r.Spec.Selector {
			if len(unitSelector.Name) == 0 && len(unitSelector.Label) == 0 && len(unitSelector.IP) == 0 {
				return fmt.Errorf("must provide one of \"name\"、\"label\"、\"ip\" in selector")
			}

			var emptyCount int
			if len(unitSelector.Name) == 0 {
				emptyCount++
			}
			if len(unitSelector.Label) == 0 {
				emptyCount++
			}
			if len(unitSelector.IP) == 0 {
				emptyCount++
			}

			if emptyCount == 0 {
				return fmt.Errorf("must provide one type of selector in \"nodeName\"、\"label\"、\"ip\"")
			}

			if emptyCount != 1 {
				return fmt.Errorf("can only provide one type of selector in \"nodeName\"、\"label\"、\"ip\"")
			}
		}
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Experiment) ValidateUpdate(old runtime.Object) error {
	experimentlog.Info("validate update", "name", r.Name)
	oldExp := old.(*Experiment)
	if reflect.DeepEqual(r.Spec, oldExp.Spec) {
		return nil
	}

	if !reflect.DeepEqual(r.Spec.Experiment, oldExp.Spec.Experiment) ||
		!reflect.DeepEqual(r.Spec.Selector, oldExp.Spec.Selector) ||
		!reflect.DeepEqual(r.Spec.RangeMode, oldExp.Spec.RangeMode) ||
		r.Spec.Scope != oldExp.Spec.Scope {
		return fmt.Errorf("spec only support update \"targetPhase\"")
	}

	if !(oldExp.Status.Phase == InjectPhaseType && (oldExp.Status.Status == SuccessStatusType || oldExp.Status.Status == FailedStatusType || oldExp.Status.Status == PartSuccessStatusType)) {
		return fmt.Errorf("only support update when \"status.phase == inject and status.status == success/failed/partSuccess\"")
	}

	if r.Spec.TargetPhase != RecoverPhaseType {
		return fmt.Errorf("only can update \"targetPhase\" to \"recover\"")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Experiment) ValidateDelete() error {
	experimentlog.Info("validate delete", "name", r.Name)
	return nil
}

func ConvertDuration(d string) (time.Duration, error) {
	unit := d[len(d)-1]
	var value string
	if unit != 'h' && unit != 'm' && unit != 's' {
		value, unit = d, 's'
	} else {
		value = d[:len(d)-1]
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 's':
		return time.Duration(v) * time.Second, nil
	case 'm':
		return time.Duration(v) * time.Minute, nil
	case 'h':
		return time.Duration(v) * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown time unit: %d", unit)
	}
}
