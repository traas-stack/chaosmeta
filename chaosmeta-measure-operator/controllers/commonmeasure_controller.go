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
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/config"
	"github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	measurev1alpha1 "github.com/traas-stack/chaosmeta/chaosmeta-measure-operator/api/v1alpha1"
)

// CommonMeasureReconciler reconciles a CommonMeasure object
type CommonMeasureReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=chaosmeta.io,resources=commonmeasures,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=chaosmeta.io,resources=commonmeasures/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=chaosmeta.io,resources=commonmeasures/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CommonMeasure object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *CommonMeasureReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance, logger := &measurev1alpha1.CommonMeasure{}, log.FromContext(ctx)
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get instance error: %s", err.Error())
	}

	defer func() {
		if e := recover(); e != any(nil) {
			// catch exception from solve experiment
			logger.Error(fmt.Errorf("catch exception: %v", e), fmt.Sprintf("when processing measure: %s/%s", instance.Namespace, instance.Name))
		}
	}()

	// TODO: need to test
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if instance.Status.Status != measurev1alpha1.SuccessStatus && instance.Status.Status != measurev1alpha1.FailedStatus {
			if !instance.Spec.Stopped {
				instance.Spec.Stopped = true
				logger.Info(fmt.Sprintf("update spec.stopped of %s/%s from false to true", instance.Namespace, instance.Name))
				return ctrl.Result{}, r.Update(ctx, instance)
			}
		} else {
			solveFinalizer(instance)
			logger.Info(fmt.Sprintf("update Finalizer of %s/%s to: %s", instance.Namespace, instance.Name, instance.ObjectMeta.Finalizers))
			return ctrl.Result{}, r.Update(ctx, instance)
		}
	}

	if instance.Status.Status == "" {
		instance.Status.Status = measurev1alpha1.CreatedStatus
	}

	logger.Info(fmt.Sprintf("process instance %s/%s, status: %s", instance.Namespace, instance.Name, instance.Status.Status))
	switch instance.Status.Status {
	case measurev1alpha1.CreatedStatus:
		initialData(ctx, instance)
	case measurev1alpha1.RunningStatus:
		processTask(ctx, instance)
	default:
		return ctrl.Result{}, nil
	}

	status, _ := json.Marshal(instance.Status)
	logger.Info(fmt.Sprintf("measure: %s/%s, start to update status: %s", instance.Namespace, instance.Name, string(status)))
	if err := r.Client.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("update instance error: %s", err.Error())
	}

	return ctrl.Result{}, nil
}

func initialData(ctx context.Context, ins *measurev1alpha1.CommonMeasure) {
	var e = measurev1alpha1.GetMeasureExecutor(ctx, ins.Spec.MeasureType)
	data, err := e.InitialData(ctx, ins.Spec.Args)
	if err != nil {
		ins.Status.Status = measurev1alpha1.FailedStatus
		ins.Status.Message = fmt.Sprintf("initial data error: %s", err.Error())
	} else {
		ins.Status.Status = measurev1alpha1.RunningStatus
		ins.Status.Message = "initial data success"
	}

	nowTime := time.Now()
	interval, _ := measurev1alpha1.ConvertDuration(ins.Spec.Interval)
	ins.Status.CreateTime, ins.Status.UpdateTime = nowTime.Format(measurev1alpha1.TimeFormat), nowTime.Format(measurev1alpha1.TimeFormat)
	ins.Status.NextTime = nowTime.Add(interval).Format(measurev1alpha1.TimeFormat)
	ins.Status.InitialData = data
}

func processTask(ctx context.Context, ins *measurev1alpha1.CommonMeasure) {
	// check judgement if meet
	if judge(ctx, ins) {
		return
	}

	// if meet interval do: measure
	if meetInterval, _ := utils.IsTimeout(ins.Status.NextTime, "0s"); meetInterval {
		measure(ctx, ins)
	}
}

func execJudge(ctx context.Context, ins *measurev1alpha1.CommonMeasure) bool {
	if ins.Spec.SuccessCount > 0 {
		if ins.Status.SuccessMeasure >= ins.Spec.SuccessCount {
			return true
		} else {
			return false
		}
	}

	if ins.Spec.FailedCount > 0 {
		if ins.Status.FailedMeasure >= ins.Spec.FailedCount {
			return false
		} else {
			return true
		}
	}

	return false
}

func judge(ctx context.Context, ins *measurev1alpha1.CommonMeasure) bool {
	time.Sleep(time.Second)
	result := execJudge(ctx, ins)
	timeout, _ := utils.IsTimeout(ins.Status.CreateTime, ins.Spec.Duration)
	nowTime := time.Now().Format(measurev1alpha1.TimeFormat)
	ins.Status.UpdateTime = nowTime

	if timeout || ins.Spec.Stopped {
		if result {
			ins.Status.Status = measurev1alpha1.SuccessStatus
			ins.Status.Message = "measure success"
		} else {
			ins.Status.Status = measurev1alpha1.FailedStatus
			ins.Status.Message = "measure failed"
		}

		return true
	}
	return false
}

func measure(ctx context.Context, ins *measurev1alpha1.CommonMeasure) {
	logger := log.FromContext(ctx)
	var e = measurev1alpha1.GetMeasureExecutor(ctx, ins.Spec.MeasureType)
	err := e.Measure(ctx, ins.Spec.Args, ins.Spec.Judgement, ins.Status.InitialData)
	nowTime := time.Now()
	status := measurev1alpha1.SuccessStatus
	msg := "measure success"

	ins.Status.TotalMeasure++
	if err == nil {
		ins.Status.SuccessMeasure++
	} else {
		ins.Status.FailedMeasure++
		status = measurev1alpha1.FailedStatus
		msg = err.Error()
	}

	interval, _ := measurev1alpha1.ConvertDuration(ins.Spec.Interval)
	ins.Status.NextTime = nowTime.Add(interval).Format(measurev1alpha1.TimeFormat)

	ins.Status.Measures = append([]measurev1alpha1.MeasureTask{
		{
			Uid:        utils.NewUid(),
			CreateTime: nowTime.Format(measurev1alpha1.TimeFormat),
			UpdateTime: nowTime.Format(measurev1alpha1.TimeFormat),
			Message:    msg,
			Status:     status,
		},
	}, ins.Status.Measures...)

	if config.GetGlobalConfig().TaskLimit > 0 && len(ins.Status.Measures) > config.GetGlobalConfig().TaskLimit {
		logger.Info(fmt.Sprintf("now length is %d, need to truncate to %d", len(ins.Status.Measures), config.GetGlobalConfig().TaskLimit))
		ins.Status.Measures = ins.Status.Measures[:config.GetGlobalConfig().TaskLimit]
	}

	ins.Status.Message = fmt.Sprintf("total measures: %d, success: %d, failed: %d", ins.Status.TotalMeasure, ins.Status.SuccessMeasure, ins.Status.FailedMeasure)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommonMeasureReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&measurev1alpha1.CommonMeasure{}).
		Complete(r)
}

func solveFinalizer(instance *measurev1alpha1.CommonMeasure) {
	for index := 0; index < len(instance.ObjectMeta.Finalizers); index++ {
		if instance.ObjectMeta.Finalizers[index] == measurev1alpha1.FinalizerName {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers[:index], instance.ObjectMeta.Finalizers[index+1:]...)
			return
		}
	}
}
