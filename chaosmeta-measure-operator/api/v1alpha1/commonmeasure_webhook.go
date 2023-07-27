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
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strconv"
	"time"
)

// log is for logging in this package.
var commonmeasurelog = logf.Log.WithName("commonmeasure-resource")

func (r *CommonMeasure) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-chaosmeta-io-v1alpha1-commonmeasure,mutating=true,failurePolicy=fail,sideEffects=None,groups=chaosmeta.io,resources=commonmeasures,verbs=create;update,versions=v1alpha1,name=mcommonmeasure.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CommonMeasure{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CommonMeasure) Default() {
	commonmeasurelog.Info("default", "name", r.Name)
	r.Status.Status = CreatedStatus
	r.Spec.Stopped = false

	var i int
	for i = 0; i < len(r.ObjectMeta.Finalizers); i++ {
		if r.ObjectMeta.Finalizers[i] == FinalizerName {
			break
		}
	}

	if i == len(r.ObjectMeta.Finalizers) {
		r.ObjectMeta.Finalizers = append(r.ObjectMeta.Finalizers, FinalizerName)
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-chaosmeta-io-v1alpha1-commonmeasure,mutating=false,failurePolicy=fail,sideEffects=None,groups=chaosmeta.io,resources=commonmeasures,verbs=create;update,versions=v1alpha1,name=vcommonmeasure.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CommonMeasure{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CommonMeasure) ValidateCreate() error {
	commonmeasurelog.Info("validate create", "name", r.Name)

	if r.Spec.SuccessCount == 0 && r.Spec.FailedCount == 0 {
		return fmt.Errorf("one of spec.successCount and spec.failedCount must provide")
	}

	ctx := context.Background()
	e := GetMeasureExecutor(ctx, r.Spec.MeasureType)
	if e == nil {
		return fmt.Errorf("%s measure executor is not set, please check operator config and reload", r.Spec.MeasureType)
	}

	if err := e.CheckConfig(ctx, r.Spec.Args, r.Spec.Judgement); err != nil {
		return fmt.Errorf("measure config is error: %s", err.Error())
	}

	if _, err := ConvertDuration(r.Spec.Duration); err != nil {
		return fmt.Errorf("spec.duration is not a valid duration: %s", err.Error())
	}

	if _, err := ConvertDuration(r.Spec.Interval); err != nil {
		return fmt.Errorf("spec.interval is not a valid duration: %s", err.Error())
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CommonMeasure) ValidateUpdate(old runtime.Object) error {
	commonmeasurelog.Info("validate update", "name", r.Name)
	oldIns := old.(*CommonMeasure)
	if reflect.DeepEqual(r.Spec, oldIns.Spec) {
		return nil
	}

	if oldIns.Spec.Stopped == true || r.Spec.Stopped == false {
		return fmt.Errorf("only support update spec.stopped from false to true")
	}

	if !reflect.DeepEqual(r.Spec.Args, oldIns.Spec.Args) ||
		!reflect.DeepEqual(r.Spec.Judgement, oldIns.Spec.Judgement) ||
		r.Spec.MeasureType != oldIns.Spec.MeasureType ||
		r.Spec.SuccessCount != oldIns.Spec.SuccessCount ||
		r.Spec.FailedCount != oldIns.Spec.FailedCount ||
		r.Spec.Duration != oldIns.Spec.Duration ||
		r.Spec.Interval != oldIns.Spec.Interval {
		return fmt.Errorf("only support update spec.stopped")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CommonMeasure) ValidateDelete() error {
	commonmeasurelog.Info("validate delete", "name", r.Name)
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
