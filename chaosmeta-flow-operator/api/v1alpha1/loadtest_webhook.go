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
)

// log is for logging in this package.
var loadtestlog = logf.Log.WithName("loadtest-resource")

func (r *LoadTest) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-chaosmeta-io-v1alpha1-loadtest,mutating=true,failurePolicy=fail,sideEffects=None,groups=chaosmeta.io,resources=loadtests,verbs=create,versions=v1alpha1,name=mloadtest.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LoadTest{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LoadTest) Default() {
	loadtestlog.Info("default", "name", r.Name)
	if r.Spec.Stopped {
		loadtestlog.Info("update \"stopped\" to false")
		r.Spec.Stopped = false
	}

	var i int
	for i = 0; i < len(r.ObjectMeta.Finalizers); i++ {
		if r.ObjectMeta.Finalizers[i] == FinalizerName {
			break
		}
	}

	if i == len(r.ObjectMeta.Finalizers) {
		loadtestlog.Info(fmt.Sprintf("add \"%s\" finalizer", FinalizerName))
		r.ObjectMeta.Finalizers = append(r.ObjectMeta.Finalizers, FinalizerName)
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-chaosmeta-io-v1alpha1-loadtest,mutating=false,failurePolicy=fail,sideEffects=None,groups=chaosmeta.io,resources=loadtests,verbs=create;update,versions=v1alpha1,name=vloadtest.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &LoadTest{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateCreate() error {
	loadtestlog.Info("validate create", "name", r.Name)

	if r.Spec.Stopped {
		return fmt.Errorf("\"spec.stopped\" should be false")
	}

	if r.Spec.Parallelism < r.Spec.Source {
		return fmt.Errorf("must meet Spec.Parallelism >= Spec.Source")
	}

	if _, err := ConvertDuration(r.Spec.Duration); err != nil {
		return fmt.Errorf("spec.duration is not a valid duration: %s", err.Error())
	}

	if r.Spec.FlowType != HTTPFlowType {
		return fmt.Errorf("only support flowType: %s", HTTPFlowType)
	}

	argsMap := GetArgsMap(r.Spec.Args)
	header, ok := argsMap[HeaderArgsKey]
	if ok {
		_, err := GetHeaderMap(header)
		if err != nil {
			return fmt.Errorf("header args format is error: %s", err.Error())
		}
	}

	if argsMap[MethodArgsKey] != MethodGET && argsMap[MethodArgsKey] != MethodPOST {
		return fmt.Errorf("method only support: %s or %s", MethodPOST, MethodGET)
	}

	if _, ok := argsMap[HostArgsKey]; !ok {
		return fmt.Errorf("must provide %s args", HostArgsKey)
	}

	if _, ok := argsMap[PortArgsKey]; !ok {
		return fmt.Errorf("must provide %s args", PortArgsKey)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateUpdate(old runtime.Object) error {
	loadtestlog.Info("validate update", "name", r.Name)

	oldIns := old.(*LoadTest)
	if reflect.DeepEqual(r.Spec, oldIns.Spec) {
		return nil
	}

	if oldIns.Spec.Stopped != r.Spec.Stopped && r.Spec.Stopped == false {
		return fmt.Errorf("only support update spec.stopped from false to true")
	}

	if !reflect.DeepEqual(r.Spec.Args, oldIns.Spec.Args) ||
		!reflect.DeepEqual(r.Spec.FlowType, oldIns.Spec.FlowType) ||
		r.Spec.Duration != oldIns.Spec.Duration ||
		r.Spec.Parallelism != oldIns.Spec.Parallelism ||
		r.Spec.Source != oldIns.Spec.Source {
		return fmt.Errorf("only support update spec.stopped")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateDelete() error {
	loadtestlog.Info("validate delete", "name", r.Name)

	return nil
}

func ConvertDuration(d string) (int, error) {
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
		return v, nil
	case 'm':
		return v * 60, nil
	case 'h':
		return v * 60 * 60, nil
	default:
		return 0, fmt.Errorf("unknown time unit: %d", unit)
	}
}

func GetHeaderMap(headerStr string) (map[string]string, error) {
	headerStrList := strings.Split(headerStr, ",")
	headerMap := make(map[string]string)
	for _, unitStr := range headerStrList {
		kvList := strings.Split(unitStr, ":")
		if len(kvList) != 2 {
			return nil, fmt.Errorf("%s is format error, example: k1:v1,k2:v2", unitStr)
		}
		headerMap[kvList[0]] = kvList[1]
	}

	return headerMap, nil
}

func GetArgsMap(args []FlowArgs) map[string]string {
	argsMap := make(map[string]string)
	for _, unit := range args {
		argsMap[unit.Key] = unit.Value
	}

	return argsMap
}
