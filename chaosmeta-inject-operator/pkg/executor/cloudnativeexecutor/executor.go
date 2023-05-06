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

package cloudnativeexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
)

type CloudNativeExecutor interface {
	Inject(ctx context.Context, injectObject, uid, timeout string, args []v1alpha1.ArgsUnit) (string, error)
	Recover(ctx context.Context, injectObject, uid, backup string) error
	Query(ctx context.Context, injectObject, uid, backup string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error)
}

const (
	batchNameFormat = "%s-%d"
)

var (
	cloudNativeExecutorMap = make(map[string]CloudNativeExecutor)
	resourceCreateFuncMap  = make(map[string]func(ctx context.Context, namespace, name string) error)
)

func GetCloudNativeExecutor(target v1alpha1.CloudTargetType, fault string) CloudNativeExecutor {
	return cloudNativeExecutorMap[fmt.Sprintf("%s%s%s", target, model.ObjectNameSplit, fault)]
}

func registerCloudExecutor(target v1alpha1.CloudTargetType, fault string, e CloudNativeExecutor) {
	cloudNativeExecutorMap[fmt.Sprintf("%s%s%s", target, model.ObjectNameSplit, fault)] = e
}

func getResourceCreateFunc(target v1alpha1.CloudTargetType, fault string) func(ctx context.Context, namespace, name string) error {
	return resourceCreateFuncMap[fmt.Sprintf("%s%s%s", target, model.ObjectNameSplit, fault)]
}

func registerResourceCreateFunc(target v1alpha1.CloudTargetType, fault string, f func(ctx context.Context, namespace, name string) error) {
	resourceCreateFuncMap[fmt.Sprintf("%s%s%s", target, model.ObjectNameSplit, fault)] = f
}

func getNewFinalizers(ctx context.Context, oldFinalizers []string, args []v1alpha1.ArgsUnit) []string {
	reArgs := common.GetArgs(args, []string{"add", "delete"})
	var addStr, deleteStr = reArgs[0], reArgs[1]
	var addArr, deleteArr []string
	if addStr != "" {
		addArr = strings.Split(addStr, v1alpha1.ArgsListSplit)
	}

	if deleteStr != "" {
		deleteArr = strings.Split(deleteStr, v1alpha1.ArgsListSplit)
	}

	if len(addArr) == 0 && len(deleteArr) == 0 {
		return oldFinalizers
	}

	// delete first and add later
	deleteMap := make(map[string]bool)
	for _, unit := range deleteArr {
		deleteMap[unit] = true
	}

	var newFinalizers []string
	newExist := make(map[string]bool)
	for _, unit := range oldFinalizers {
		if deleteMap[unit] || newExist[unit] {
			continue
		}
		newFinalizers = append(newFinalizers, unit)
		newExist[unit] = true
	}

	for _, unit := range addArr {
		if newExist[unit] {
			continue
		}

		newFinalizers = append(newFinalizers, unit)
		newExist[unit] = true
	}

	return newFinalizers
}

func patchFinalizers(ctx context.Context, c rest.Interface, resource, ns, name string, finalizers []string) error {
	payload, err := json.Marshal(finalizers)
	if err != nil {
		return fmt.Errorf("get payload error: %s", err.Error())
	}

	if err := c.Patch(types.MergePatchType).Namespace(ns).Resource(resource).Name(name).
		Body([]byte(fmt.Sprintf(`{"metadata":{"finalizers":%s}}`, payload))).Do(ctx).Error(); err != nil {
		return fmt.Errorf("patch finalizers error: %s", err.Error())
	}

	return nil
}

func getBackupLabels(backup []byte, nowLabels map[string]string) ([]byte, error) {
	var backupMap map[string]interface{}

	if backup == nil || len(backup) == 0 {
		backupMap = make(map[string]interface{})
	} else {
		if err := json.Unmarshal(backup, &backupMap); err != nil {
			return nil, fmt.Errorf("backup labels is not a json: %s", err.Error())
		}
	}

	for k := range nowLabels {
		if _, ok := backupMap[k]; !ok {
			backupMap[k] = nil
		}
	}

	backupBytes, err := json.Marshal(backupMap)
	if err != nil {
		return nil, fmt.Errorf("backup to string error: %s", err.Error())
	}

	return backupBytes, nil
}

func getNewLabels(ctx context.Context, oldLabels map[string]string, args []v1alpha1.ArgsUnit) ([]byte, error) {
	reArgs := common.GetArgs(args, []string{"add", "delete"})
	var addStr, deleteStr = reArgs[0], reArgs[1]

	var addArr, deleteArr []string
	if addStr != "" {
		addArr = strings.Split(addStr, v1alpha1.ArgsListSplit)
	}
	if deleteStr != "" {
		deleteArr = strings.Split(deleteStr, v1alpha1.ArgsListSplit)
	}

	var reMap map[string]interface{}
	reMap = make(map[string]interface{})
	for _, unit := range addArr {
		if unit == "" {
			continue
		}

		tmpArr := strings.Split(unit, v1alpha1.LabelListSplit)
		if len(tmpArr) != 2 {
			return nil, fmt.Errorf("%s is error label format, true format is key=value", unit)
		}

		reMap[tmpArr[0]] = tmpArr[1]
	}

	deleteMap := make(map[string]bool)
	for _, unit := range deleteArr {
		deleteMap[unit] = true
	}

	for k, v := range oldLabels {
		if deleteMap[k] {
			reMap[k] = nil
			continue
		}

		if _, ok := reMap[k]; ok {
			continue
		}

		reMap[k] = v
	}

	reByte, err := json.Marshal(reMap)
	if err != nil {
		return nil, fmt.Errorf("labels to string error: %s", err.Error())
	}

	return reByte, nil
}

func getPatchLabels(labels []byte) []byte {
	if len(labels) == 0 {
		labels = []byte("{}")
	}

	return []byte(fmt.Sprintf(`{"metadata":{"labels":%s}}`, labels))
}

func patchLabels(ctx context.Context, c rest.Interface, resource, ns, name string, labels []byte) error {
	//patchBytes, err := json.Marshal(labels)
	//if err != nil {
	//	return fmt.Errorf("get patchBytes error: %s", err.Error())
	//}

	if err := c.Patch(types.MergePatchType).Namespace(ns).Resource(resource).Name(name).
		Body(getPatchLabels(labels)).Do(ctx).Error(); err != nil {
		return fmt.Errorf("patch error: %s", err.Error())
	}

	return nil
}

func createNs(ctx context.Context, name string) error {
	return restclient.GetApiServerClientMap(v1alpha1.NamespaceCloudTarget).Post().Resource("namespaces").
		Body(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}).Do(ctx).Error()
}

func deleteNs(ctx context.Context, name string) error {
	return restclient.GetApiServerClientMap(v1alpha1.NamespaceCloudTarget).Delete().Resource("namespaces").
		Name(name).Do(ctx).Error()
}

func batchResourceInject(ctx context.Context, args []v1alpha1.ArgsUnit, target v1alpha1.CloudTargetType, fault string) (string, error) {
	reArgs := common.GetArgs(args, []string{"count", "namespace", "name"})
	count, err := strconv.Atoi(reArgs[0])
	if err != nil {
		return "", fmt.Errorf("\"count\" is not a num: %s", err.Error())
	}

	namespace, name := reArgs[1], reArgs[2]
	if namespace == "" {
		return "", fmt.Errorf("namespace is empty")
	}
	if name == "" {
		return "", fmt.Errorf("name is empty")
	}

	// already exists will conflict
	if err := createNs(ctx, namespace); err != nil {
		return "", fmt.Errorf("create namespace error: %s", err.Error())
	}

	workers := runtime.NumCPU() / 2
	if workers == 0 {
		workers = 1
	}
	perCount, additional := count/workers, count%workers

	if ok := common.GetClusterCtrl().Run(int64(workers)); !ok {
		return "", fmt.Errorf("has other running cluster task, please retry later")
	}

	var startIndex int
	for i := 0; i < workers; i++ {
		endIndex := startIndex + perCount
		if i == workers-1 {
			endIndex += additional
		}

		go batchCreateResource(ctx, namespace, name, startIndex, endIndex, target, fault)
		startIndex = endIndex
	}

	return namespace, nil
}

func batchCreateResource(ctx context.Context, namespace, name string, start, end int, target v1alpha1.CloudTargetType, fault string) {
	defer common.GetClusterCtrl().FinishOne()

	ctrl, logger := common.GetClusterCtrl(), log.FromContext(ctx)
	for i := start; i < end; i++ {
		if ctrl.IsStopping() {
			return
		}
		if err := getResourceCreateFunc(target, fault)(ctx, namespace, fmt.Sprintf(batchNameFormat, name, i)); err != nil {
			logger.Error(err, "create resource error")
		}
	}
}
