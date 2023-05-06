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

package injecthandler

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/scopehandler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sync"
	"time"
)

type InjectPhaseHandler struct {
}

func (h *InjectPhaseHandler) SolveCreated(ctx context.Context, exp *v1alpha1.Experiment) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveCreated start", exp.Namespace, exp.Name))

	isTimeout, err := common.IsTimeout(exp.Status.CreateTime, exp.Spec.Experiment.Duration)
	if err != nil {
		// Unexpected case are treated as timeout
		isTimeout = true
		logger.Error(err, "check if experiment timeout error")
	}

	var (
		targetSubExp = exp.Status.Detail.Inject
		wg           = sync.WaitGroup{}
	)

	for i := range exp.Status.Detail.Inject {
		if targetSubExp[i].Status != v1alpha1.CreatedStatusType {
			continue
		}

		common.GetGoroutinePool().GetGoroutine()
		wg.Add(1)
		go solveCreated(ctx, &wg, exp, i, isTimeout)
	}

	wg.Wait()
	// Summarize subtask execution results
	var failCount, createdCount int
	for i := range targetSubExp {
		if targetSubExp[i].Status == v1alpha1.FailedStatusType {
			failCount++
		} else if targetSubExp[i].Status == v1alpha1.CreatedStatusType {
			createdCount++
		}
	}

	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveCreated: totalCount[%d], failCount[%d], createdCount[%d]", exp.Namespace, exp.Name, len(targetSubExp), failCount, createdCount))
	// Update the overall task status
	if createdCount > 0 {
		exp.Status.Status, exp.Status.Message = v1alpha1.CreatedStatusType, "created count is more than 0, need to retry"
	} else {
		if failCount == len(targetSubExp) {
			exp.Status.Status, exp.Status.Message = v1alpha1.FailedStatusType, "create failed"
		} else {
			exp.Status.Status, exp.Status.Message = v1alpha1.RunningStatusType, "create finish, start to solve running status"
		}
	}

	exp.Status.UpdateTime = time.Now().Format(model.TimeFormat)
}

func solveCreated(ctx context.Context, wg *sync.WaitGroup, exp *v1alpha1.Experiment, i int, isTimeout bool) {
	var (
		logger       = log.FromContext(ctx)
		targetSubExp = exp.Status.Detail.Inject
		scopeHandler = scopehandler.GetScopeHandler(exp.Spec.Scope)
		commonObject model.AtomicObject
		err          error
	)

	logger.Info(fmt.Sprintf("experiment: %s/%s/%s, solveCreated start, now Goroutine: %d", exp.Namespace, exp.Name, targetSubExp[i].InjectObjectName, common.GetGoroutinePool().GetLen()))

	defer func() {
		common.GetGoroutinePool().ReleaseGoroutine()
		wg.Done()
		logger.Info(fmt.Sprintf("experiment: %s/%s/%s, solveCreated finish, status: %s, now Goroutine: %d", exp.Namespace, exp.Name, targetSubExp[i].InjectObjectName, targetSubExp[i].Status, common.GetGoroutinePool().GetLen()))
	}()

	commonObject, err = scopeHandler.GetInjectObject(ctx, exp.Spec.Experiment, targetSubExp[i].InjectObjectName)
	if err != nil {
		if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.CreatedStatusType, "GetInjectObject network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "GetInjectObject network error, timeout"
			}
		} else {
			// include not found
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("GetInjectObject error: %s", err.Error())
		}

		return
	}

	backup, err := scopeHandler.ExecuteInject(ctx, commonObject, targetSubExp[i].UID, exp.Spec.Experiment)
	if err != nil {
		if common.IsKeyUniqueErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.RunningStatusType, "experiment start success"
		} else if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.CreatedStatusType, "experiment inject network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "experiment inject network error, timeout"
			}
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("experiment inject error: %s", err.Error())
		}
	} else {
		targetSubExp[i].Backup, targetSubExp[i].Status, targetSubExp[i].Message = backup, v1alpha1.RunningStatusType, "experiment inject start success"
	}
}

func (h *InjectPhaseHandler) SolveRunning(ctx context.Context, exp *v1alpha1.Experiment) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveRunning start", exp.Namespace, exp.Name))

	isTimeout, err := common.IsTimeout(exp.Status.CreateTime, exp.Spec.Experiment.Duration)
	if err != nil {
		// Unexpected case are treated as timeout
		isTimeout = true
		logger.Error(err, "check if experiment timeout error")
	}

	var (
		targetSubExp = exp.Status.Detail.Inject
		wg           = sync.WaitGroup{}
	)

	for i := range targetSubExp {
		if targetSubExp[i].Status != v1alpha1.RunningStatusType {
			continue
		}

		common.GetGoroutinePool().GetGoroutine()
		wg.Add(1)
		go solveRunning(ctx, &wg, exp, i, isTimeout)
	}

	wg.Wait()

	var runCount, failCount int
	for i := range targetSubExp {
		if targetSubExp[i].Status == v1alpha1.RunningStatusType {
			runCount++
		} else if targetSubExp[i].Status == v1alpha1.FailedStatusType {
			failCount++
		}
	}

	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveRunning: totalCount[%d], failCount[%d], runCount[%d]", exp.Namespace, exp.Name, len(targetSubExp), failCount, runCount))

	if runCount > 0 {
		exp.Status.Status, exp.Status.Message = v1alpha1.RunningStatusType, "run count is more than 0, need to retry"
	} else {
		if failCount == 0 {
			exp.Status.Status, exp.Status.Message = v1alpha1.SuccessStatusType, "run success"
		} else if failCount == len(targetSubExp) {
			exp.Status.Status, exp.Status.Message = v1alpha1.FailedStatusType, "run all failed"
		} else {
			exp.Status.Status, exp.Status.Message = v1alpha1.PartSuccessStatusType, "run part success"
		}
	}

	time.Sleep(time.Second)
	exp.Status.UpdateTime = time.Now().Format(model.TimeFormat)
}

func solveRunning(ctx context.Context, wg *sync.WaitGroup, exp *v1alpha1.Experiment, i int, isTimeout bool) {
	var (
		logger       = log.FromContext(ctx)
		targetSubExp = exp.Status.Detail.Inject
		scopeHandler = scopehandler.GetScopeHandler(exp.Spec.Scope)
		commonObject model.AtomicObject
		err          error
	)

	logger.Info(fmt.Sprintf("experiment: %s/%s/%s, solveRunning start, now Goroutine: %d", exp.Namespace, exp.Name, targetSubExp[i].InjectObjectName, common.GetGoroutinePool().GetLen()))

	defer func() {
		common.GetGoroutinePool().ReleaseGoroutine()
		wg.Done()
		logger.Info(fmt.Sprintf("experiment: %s/%s/%s, solveRunning finish, status: %s, now Goroutine: %d", exp.Namespace, exp.Name, targetSubExp[i].InjectObjectName, targetSubExp[i].Status, common.GetGoroutinePool().GetLen()))
	}()

	commonObject, err = scopeHandler.GetInjectObject(ctx, exp.Spec.Experiment, targetSubExp[i].InjectObjectName)
	if err != nil {
		if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.RunningStatusType, "GetInjectObject network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "GetInjectObject network error, timeout"
			}
		} else if common.IsNotFoundErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.SuccessStatusType, err.Error()
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("GetInjectObject error: %s", err.Error())
		}

		return
	}

	expInfo, err := scopeHandler.QueryExperiment(ctx, commonObject, targetSubExp[i].UID, targetSubExp[i].Backup, exp.Spec.Experiment, v1alpha1.InjectPhaseType)
	if err != nil {
		if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.RunningStatusType, "experiment query network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "experiment query network error, timeout"
			}
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("experiment query error: %s", err.Error())
			targetSubExp[i].UpdateTime = time.Now().Format(model.TimeFormat)
		}

		return
	} else {
		if expInfo.Status == v1alpha1.SuccessStatusType || expInfo.Status == v1alpha1.FailedStatusType || expInfo.Status == v1alpha1.RunningStatusType {
			targetSubExp[i].StartTime, targetSubExp[i].UpdateTime = expInfo.CreateTime, expInfo.UpdateTime
			targetSubExp[i].Status, targetSubExp[i].Message = expInfo.Status, expInfo.Message
		} else {
			logger.Error(fmt.Errorf("unexpected status"), fmt.Sprintf("expInfo.Status is %s", expInfo.Status))
			return
		}
	}
}

func (h *InjectPhaseHandler) SolveSuccess(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolveSuccess start", exp.Namespace, exp.Name))
	solveFinalStatus(ctx, exp)
}

func (h *InjectPhaseHandler) SolvePartSuccess(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolvePartSuccess start", exp.Namespace, exp.Name))
	solveFinalStatus(ctx, exp)
}

func (h *InjectPhaseHandler) SolveFailed(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolveFailed start", exp.Namespace, exp.Name))
	solveFinalStatus(ctx, exp)
}

func solveFinalStatus(ctx context.Context, exp *v1alpha1.Experiment) {
	if exp.Spec.TargetPhase == exp.Status.Phase || exp.Spec.TargetPhase != v1alpha1.RecoverPhaseType {
		return
	}

	injectDetail := exp.Status.Detail.Inject
	recoverDetail := make([]v1alpha1.ExperimentDetailUnit, len(injectDetail))
	nowTime := time.Now().Format(model.TimeFormat)
	for i := range injectDetail {
		recoverDetail[i] = v1alpha1.ExperimentDetailUnit{
			InjectObjectName: injectDetail[i].InjectObjectName,
			UID:              injectDetail[i].UID,
			Status:           v1alpha1.CreatedStatusType,
			Message:          "start to recover",
			StartTime:        nowTime,
			Backup:           injectDetail[i].Backup,
		}
	}

	exp.Status.Phase, exp.Status.Status, exp.Status.Detail.Recover = v1alpha1.RecoverPhaseType, v1alpha1.CreatedStatusType, recoverDetail
	exp.Status.UpdateTime = nowTime
}
