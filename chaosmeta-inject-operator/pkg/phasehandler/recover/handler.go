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

package recoverhandler

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

type RecoverPhaseHandler struct {
}

func (h *RecoverPhaseHandler) SolveCreated(ctx context.Context, exp *v1alpha1.Experiment) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveCreated start", exp.Namespace, exp.Name))

	isTimeout, err := common.IsTimeout(exp.Status.CreateTime, exp.Spec.Experiment.Duration)
	if err != nil {
		// Unexpected case are treated as timeout
		isTimeout = true
		logger.Error(err, "check if experiment timeout error")
	}

	var (
		targetSubExp = exp.Status.Detail.Recover
		wg           = sync.WaitGroup{}
	)

	for i := range targetSubExp {
		if targetSubExp[i].Status != v1alpha1.CreatedStatusType {
			continue
		}

		common.GetGoroutinePool().GetGoroutine()
		wg.Add(1)
		go solveCreated(ctx, &wg, exp, i, isTimeout)
	}

	wg.Wait()

	var failCount, createdCount int
	for i := range targetSubExp {
		if targetSubExp[i].Status == v1alpha1.FailedStatusType {
			failCount++
		} else if targetSubExp[i].Status == v1alpha1.CreatedStatusType {
			createdCount++
		}
	}

	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveCreated: totalCount[%d], failCount[%d], createdCount[%d]", exp.Namespace, exp.Name, len(targetSubExp), failCount, createdCount))

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

func (h *RecoverPhaseHandler) SolveRunning(ctx context.Context, exp *v1alpha1.Experiment) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("experiment: %s/%s, SolveRunning start", exp.Namespace, exp.Name))

	isTimeout, err := common.IsTimeout(exp.Status.CreateTime, exp.Spec.Experiment.Duration)
	if err != nil {
		// Unexpected case are treated as timeout
		isTimeout = true
		logger.Error(err, "check if experiment timeout error")
	}

	var (
		targetSubExp = exp.Status.Detail.Recover
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

	exp.Status.UpdateTime = time.Now().Format(model.TimeFormat)
}

func (h *RecoverPhaseHandler) SolveSuccess(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolveSuccess start", exp.Namespace, exp.Name))
}

func (h *RecoverPhaseHandler) SolvePartSuccess(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolvePartSuccess start", exp.Namespace, exp.Name))
}

func (h *RecoverPhaseHandler) SolveFailed(ctx context.Context, exp *v1alpha1.Experiment) {
	log.FromContext(ctx).Info(fmt.Sprintf("experiment: %s/%s, SolveFailed start", exp.Namespace, exp.Name))
}

func solveCreated(ctx context.Context, wg *sync.WaitGroup, exp *v1alpha1.Experiment, i int, isTimeout bool) {
	var (
		logger       = log.FromContext(ctx)
		scopeHandler = scopehandler.GetScopeHandler(exp.Spec.Scope)
		targetSubExp = exp.Status.Detail.Recover
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
		} else if common.IsNotFoundErr(err) {
			// not found as success in recover stage
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.SuccessStatusType, err.Error()
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.SuccessStatusType, fmt.Sprintf("GetInjectObject error: %s", err.Error())
		}

		return
	}

	if err := scopeHandler.ExecuteRecover(ctx, commonObject, targetSubExp[i].UID, targetSubExp[i].Backup, exp.Spec.Experiment); err != nil {
		if common.IsNotFoundErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.SuccessStatusType, err.Error()
		} else if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.CreatedStatusType, "experiment recover network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "experiment recover network error, timeout"
			}
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("experiment recover error: %s", err.Error())
		}
	} else {
		targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.RunningStatusType, "experiment recover start success"
	}
}

func solveRunning(ctx context.Context, wg *sync.WaitGroup, exp *v1alpha1.Experiment, i int, isTimeout bool) {
	var (
		scopeHandler = scopehandler.GetScopeHandler(exp.Spec.Scope)
		targetSubExp = exp.Status.Detail.Recover
		commonObject model.AtomicObject
		err          error
		logger       = log.FromContext(ctx)
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

	expInfo, err := scopeHandler.QueryExperiment(ctx, commonObject, targetSubExp[i].UID, targetSubExp[i].Backup, exp.Spec.Experiment, v1alpha1.RecoverPhaseType)
	if err != nil {
		if common.IsNetErr(err) {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.RunningStatusType, "experiment query network error, need to retry"
			if isTimeout {
				targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, "experiment query network error, timeout"
			}
		} else {
			targetSubExp[i].Status, targetSubExp[i].Message = v1alpha1.FailedStatusType, fmt.Sprintf("experiment query error: %s", err.Error())
		}

		return
	} else {
		if expInfo.Status == v1alpha1.SuccessStatusType || expInfo.Status == v1alpha1.FailedStatusType || expInfo.Status == v1alpha1.RunningStatusType {
			targetSubExp[i].Status, targetSubExp[i].Message = expInfo.Status, expInfo.Message
			targetSubExp[i].StartTime, targetSubExp[i].UpdateTime = expInfo.CreateTime, expInfo.UpdateTime
		} else {
			logger.Error(fmt.Errorf("unexpected status"), fmt.Sprintf("expInfo.Status is %s", expInfo.Status))
			return
		}
	}
}
