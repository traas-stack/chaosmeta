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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/phasehandler"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
	"time"
)

const defaultConcurrentReconciles = 10

// ExperimentReconciler reconciles a Experiment object
type ExperimentReconciler struct {
	client.Client
	//RESTClient rest.Interface
	//RESTConfig *rest.Config
	//Scheme     *runtime.Scheme
}

//+kubebuilder:rbac:groups=chaosmeta.io,resources=experiments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=chaosmeta.io,resources=experiments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=chaosmeta.io,resources=experiments/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods;pods/exec;services;namespaces;nodes,verbs=*
//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=*
//+kubebuilder:rbac:groups=batchs,resources=jobs,verbs=*

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Experiment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *ExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance, logger := &v1alpha1.Experiment{}, log.FromContext(ctx)
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, fmt.Errorf("get instance error: %s", err.Error())
	}

	defer func() {
		if e := recover(); e != any(nil) {
			// catch exception from solve experiment
			logger.Error(fmt.Errorf("catch exception: %v", e), fmt.Sprintf("when processing experiment: %s/%s", instance.Namespace, instance.Name))
		}
	}()

	status, _ := json.Marshal(instance.Status)
	logger.Info(fmt.Sprintf("experiment: %s/%s, get status: %s", instance.Namespace, instance.Name, string(status)))

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if instance.Status.Status == v1alpha1.SuccessStatusType || instance.Status.Status == v1alpha1.FailedStatusType || instance.Status.Status == v1alpha1.PartSuccessStatusType {
			if instance.Spec.TargetPhase == v1alpha1.InjectPhaseType && instance.Status.Phase == v1alpha1.InjectPhaseType {
				instance.Spec.TargetPhase = v1alpha1.RecoverPhaseType
				logger.Info(fmt.Sprintf("update TargetPhase of %s/%s to: %s", instance.Namespace, instance.Name, instance.Spec.TargetPhase))
				return ctrl.Result{}, r.Update(ctx, instance)
			} else if instance.Status.Phase == v1alpha1.RecoverPhaseType {
				solveFinalizer(instance)
				logger.Info(fmt.Sprintf("update Finalizer of %s/%s to: %s", instance.Namespace, instance.Name, instance.ObjectMeta.Finalizers))
				return ctrl.Result{}, r.Update(ctx, instance)
			}
		}
	} else {
		if instance.Status.Phase == v1alpha1.RecoverPhaseType && (instance.Status.Status == v1alpha1.SuccessStatusType ||
			instance.Status.Status == v1alpha1.FailedStatusType || instance.Status.Status == v1alpha1.PartSuccessStatusType) {
			solveFinalizer(instance)
			logger.Info(fmt.Sprintf("update Finalizer of %s/%s to: %s", instance.Namespace, instance.Name, instance.ObjectMeta.Finalizers))
			return ctrl.Result{}, r.Update(ctx, instance)
		}
	}

	if instance.Status.Phase == "" {
		initProcess(ctx, instance)
	} else {
		statusProcess(ctx, instance)
	}

	status, _ = json.Marshal(instance.Status)
	logger.Info(fmt.Sprintf("experiment: %s/%s, start to update status: %s", instance.Namespace, instance.Name, string(status)))
	if err := r.Client.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("update instance error: %s", err.Error())
	}

	return ctrl.Result{}, nil
}

func initProcess(ctx context.Context, instance *v1alpha1.Experiment) {
	// var init
	logger, nowTime := log.FromContext(ctx), time.Now().Format(model.TimeFormat)
	instance.Status.Phase, instance.Status.CreateTime, instance.Status.UpdateTime = v1alpha1.InjectPhaseType, nowTime, nowTime

	spec, _ := json.Marshal(instance.Status)
	logger.Info(fmt.Sprintf("experiment: %s/%s, spec info: %s", instance.Namespace, instance.Name, string(spec)))
	// search experiment object
	injectObjects, err := scopehandler.GetScopeHandler(instance.Spec.Scope).ConvertSelector(ctx, &instance.Spec)
	if err != nil {
		instance.Status.Status, instance.Status.Message = v1alpha1.FailedStatusType, fmt.Sprintf("convert selector to inject object error: %s", err.Error())
		return
	}
	if len(injectObjects) == 0 {
		instance.Status.Status, instance.Status.Message = v1alpha1.FailedStatusType, "no matching target"
		return
	}
	// process with range args
	injectObjects = solveRange(injectObjects, instance.Spec.RangeMode)
	details := make([]v1alpha1.ExperimentDetailUnit, len(injectObjects))
	for i, unitInjectObj := range injectObjects {
		details[i] = v1alpha1.ExperimentDetailUnit{
			InjectObjectName: unitInjectObj.GetObjectName(),
			//InjectObjectInfo: string(objBytes),
			UID:       newUid(),
			Status:    v1alpha1.CreatedStatusType,
			Message:   "Initial experiment created",
			StartTime: nowTime,
		}
	}

	instance.Status.Message = "Initial experiment created"
	instance.Status.Status, instance.Status.Detail.Inject = v1alpha1.CreatedStatusType, details
}

func statusProcess(ctx context.Context, instance *v1alpha1.Experiment) {
	handler := phasehandler.GetHandler(instance.Status.Phase)

	switch instance.Status.Status {
	case v1alpha1.CreatedStatusType:
		handler.SolveCreated(ctx, instance)
	case v1alpha1.RunningStatusType:
		handler.SolveRunning(ctx, instance)
	case v1alpha1.SuccessStatusType:
		handler.SolveSuccess(ctx, instance)
	case v1alpha1.PartSuccessStatusType:
		handler.SolvePartSuccess(ctx, instance)
	case v1alpha1.FailedStatusType:
		handler.SolveFailed(ctx, instance)
	}
}

func newUid() string {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	return fmt.Sprintf("%s%04d", timeStr, t.Nanosecond()/1000%100000%10000)
}

func solveRange(initial []model.AtomicObject, rangeMode *v1alpha1.RangeMode) []model.AtomicObject {
	if rangeMode == nil || rangeMode.Type == v1alpha1.AllRangeType {
		return initial
	}

	var count int
	if rangeMode.Type == v1alpha1.CountRangeType {
		count = rangeMode.Value
	}

	if rangeMode.Type == v1alpha1.PercentRangeType {
		count = rangeMode.Value * len(initial) / 100
	}

	if count >= len(initial) {
		return initial
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(initial), func(i int, j int) {
		initial[i], initial[j] = initial[j], initial[i]
	})

	res := initial[:count]
	sort.Slice(res, func(i, j int) bool {
		return res[i].GetObjectName() < res[j].GetObjectName()
	})

	return res
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Pod{}, selector.HostIPKey, func(rawObj client.Object) []string {
		pod := rawObj.(*corev1.Pod)
		return []string{pod.Status.HostIP}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.Experiment{}, selector.PhaseKey, func(rawObj client.Object) []string {
		exp := rawObj.(*v1alpha1.Experiment)
		return []string{string(exp.Status.Phase)}
	}); err != nil {
		return err
	}

	//if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.Experiment{}, selector.StatusKey, func(rawObj client.Object) []string {
	//	exp := rawObj.(*v1alpha1.Experiment)
	//	return []string{string(exp.Status.Status)}
	//}); err != nil {
	//	return err
	//}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Experiment{}).WithOptions(controller.Options{MaxConcurrentReconciles: defaultConcurrentReconciles}).
		Complete(r)
}

func solveFinalizer(instance *v1alpha1.Experiment) {
	for index := 0; index < len(instance.ObjectMeta.Finalizers); index++ {
		if instance.ObjectMeta.Finalizers[index] == v1alpha1.FinalizerName {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers[:index], instance.ObjectMeta.Finalizers[index+1:]...)
			return
		}
	}
}
