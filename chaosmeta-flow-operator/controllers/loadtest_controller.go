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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-flow-operator/pkg/config"
	"io"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/traas-stack/chaosmeta/chaosmeta-flow-operator/api/v1alpha1"
)

// LoadTestReconciler reconciles a LoadTest object
type LoadTestReconciler struct {
	Client    client.Client
	Scheme    *runtime.Scheme
	ClientSet *kubernetes.Clientset
	//RestfulClient rest.Interface
}

//+kubebuilder:rbac:groups=chaosmeta.io,resources=loadtests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=chaosmeta.io,resources=loadtests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=chaosmeta.io,resources=loadtests/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods;pods/log,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LoadTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *LoadTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance, logger := &v1alpha1.LoadTest{}, log.FromContext(ctx)
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get instance error: %s", err.Error())
	}
	if instance.Status.Status == "" {
		instance.Status.Status = v1alpha1.CreatedStatus
	}

	defer func() {
		if e := recover(); e != any(nil) {
			// catch exception from solve experiment
			logger.Error(fmt.Errorf("catch exception: %v", e), fmt.Sprintf("when processing measure: %s/%s", instance.Namespace, instance.Name))
		}
	}()

	if instance.Status.Status == v1alpha1.RunningStatus || instance.Status.Status == v1alpha1.CreatedStatus {
		if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
			if !instance.Spec.Stopped {
				logger.Info(fmt.Sprintf("update spec.stopped of %s/%s from false to true", instance.Namespace, instance.Name))
				instance.Spec.Stopped = true
				return ctrl.Result{}, r.Client.Update(ctx, instance)
			}
		}
	} else {
		solveFinalizer(instance)
		logger.Info(fmt.Sprintf("update Finalizer of %s/%s to: %s", instance.Namespace, instance.Name, instance.ObjectMeta.Finalizers))
		return ctrl.Result{}, r.Client.Update(ctx, instance)
	}

	logger.Info(fmt.Sprintf("process instance %s/%s, status: %s", instance.Namespace, instance.Name, instance.Status.Status))
	var err error
	switch instance.Status.Status {
	case v1alpha1.CreatedStatus:
		err = r.createJob(ctx, instance)
	case v1alpha1.RunningStatus:
		err = r.syncStatus(ctx, instance)
	default:
		return ctrl.Result{}, nil
	}

	if err != nil {
		instance.Status.Status = v1alpha1.FailedStatus
		instance.Status.Message = fmt.Sprintf("solve instance error: %s", err.Error())
	}

	instance.Status.UpdateTime = time.Now().Format(v1alpha1.TimeFormat)

	status, _ := json.Marshal(instance.Status)
	logger.Info(fmt.Sprintf("measure: %s/%s, start to update status: %s", instance.Namespace, instance.Name, string(status)))
	if err := r.Client.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("update instance error: %s", err.Error())
	}

	return ctrl.Result{}, nil
}

func (r *LoadTestReconciler) createJob(ctx context.Context, ins *v1alpha1.LoadTest) error {
	configFileStr := loadConfig(ctx, ins)
	job, err := loadJob(ctx, ins, configFileStr)
	if err != nil {
		return fmt.Errorf("load job error: %v", err)
	}

	if err := r.Client.Create(ctx, job); err != nil {
		return fmt.Errorf("create job[%s] error: %s", job.Name, err.Error())
	}

	ins.Status.Status = v1alpha1.RunningStatus
	ins.Status.CreateTime = time.Now().Format(v1alpha1.TimeFormat)
	return nil
}

func (r *LoadTestReconciler) syncStatus(ctx context.Context, ins *v1alpha1.LoadTest) error {
	time.Sleep(time.Second * 2)
	job := &batchv1.Job{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: ins.Namespace, Name: ins.Name}, job); err != nil {
		if errors.IsNotFound(err) {
			ins.Status.Status = v1alpha1.SuccessStatus
			ins.Status.Message = "job not found"
			return nil
		}
		return fmt.Errorf("failed to get job: %s", err.Error())
	}

	podList := &corev1.PodList{}
	if err := r.Client.List(context.Background(), podList, []client.ListOption{
		client.InNamespace(ins.Namespace),
		client.MatchingLabels(job.Spec.Selector.MatchLabels),
	}...); err != nil {
		return fmt.Errorf("failed to list pods: %v", err.Error())
	}

	var totalCount, totalErr = summaryFlowData(ctx, r.ClientSet, podList)
	ins.Status.TotalCount, ins.Status.SuccessCount = totalCount, totalCount-totalErr
	createTime, _ := time.ParseInLocation(v1alpha1.TimeFormat, ins.Status.CreateTime, time.Local)
	ins.Status.AvgRPS = totalCount / int(time.Now().Sub(createTime).Seconds())

	if job.Status.Active == 0 {
		ins.Status.Status = v1alpha1.SuccessStatus
		ins.Status.Message = "job finish"
	} else {
		if ins.Spec.Stopped {
			if err := deleteForce(ctx, r.ClientSet, job, podList); err != nil {
				ins.Status.Message = fmt.Sprintf("need to stop instance, but get error: %s", err.Error())
				return nil
			}
			ins.Status.Status = v1alpha1.SuccessStatus
			ins.Status.Message = "job stopped"
		}
	}

	return nil
}

func deleteForce(ctx context.Context, c *kubernetes.Clientset, job *batchv1.Job, podList *corev1.PodList) error {
	deletePolicy, gracePeriodSeconds := metav1.DeletePropagationForeground, int64(0)
	if err := c.BatchV1().Jobs(job.Namespace).Delete(ctx, job.Name, metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &gracePeriodSeconds,
	}); err != nil {
		return fmt.Errorf("failed to delete job %s/%s: %v", job.Namespace, job.Name, err)
	}
	//time.Sleep(3 * time.Second)

	for _, pod := range podList.Items {
		if err := c.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{
			PropagationPolicy:  &deletePolicy,
			GracePeriodSeconds: &gracePeriodSeconds,
		}); err != nil {
			return fmt.Errorf("failed to delete pod %s/%s: %v", job.Namespace, pod.Name, err)
		}
	}

	return nil
}

func summaryFlowData(ctx context.Context, c *kubernetes.Clientset, podList *corev1.PodList) (int, int) {
	logger := log.FromContext(ctx)
	var totalCount, totalErr int
	for _, unitPod := range podList.Items {
		logStr, err := getPodLog(ctx, c, unitPod.Namespace, unitPod.Name, unitPod.Spec.Containers[0].Name)
		if err != nil {
			logger.Error(err, fmt.Sprintf("get log of pod [%s/%s] error", unitPod.Namespace, unitPod.Name))
			continue
		}
		total, errCount := getFlowDataFromLog(logStr)
		totalCount += total
		totalErr += errCount
	}

	return totalCount, totalErr
}

func getFlowDataFromLog(logStr string) (int, int) {
	lines := strings.Split(logStr, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "summary") {
			fields := strings.Fields(logStr)
			summaryIndex, errIndex := findStrIndex(fields, "summary"), findStrIndex(fields, "Err:")
			summaryStr, errStr := fields[summaryIndex+2], fields[errIndex+1]
			summary, err := strconv.Atoi(summaryStr)
			if err != nil {
				return 0, 0
			}
			errCount, err := strconv.Atoi(errStr)
			if err != nil {
				return 0, 0
			}

			return summary, errCount
		}
	}

	return 0, 0
}

func findStrIndex(arr []string, target string) int {
	for i, unit := range arr {
		if target == unit {
			return i
		}
	}

	return -1
}

func getPodLog(ctx context.Context, c *kubernetes.Clientset, ns, name, cname string) (string, error) {
	req := c.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{Container: cname})
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening stream: %s", err.Error())
	}
	defer stream.Close()

	out := new(bytes.Buffer)
	_, err = io.Copy(out, stream)
	if err != nil {
		return "", fmt.Errorf("error copying container output: %s", err.Error())
	}

	return out.String(), nil
}

func loadConfig(ctx context.Context, ins *v1alpha1.LoadTest) string {
	durationSecond, _ := v1alpha1.ConvertDuration(ins.Spec.Duration)
	argsMap := v1alpha1.GetArgsMap(ins.Spec.Args)
	configFileStr := strings.ReplaceAll(v1alpha1.JmeterConfigStr, "@NUM_THREADS@", strconv.Itoa(ins.Spec.Parallelism/ins.Spec.Source))
	configFileStr = strings.ReplaceAll(configFileStr, "@DURATION@", strconv.Itoa(durationSecond))
	configFileStr = strings.ReplaceAll(configFileStr, "@HOST@", argsMap[v1alpha1.HostArgsKey])
	configFileStr = strings.ReplaceAll(configFileStr, "@PORT@", argsMap[v1alpha1.PortArgsKey])
	configFileStr = strings.ReplaceAll(configFileStr, "@PATH@", argsMap[v1alpha1.PathArgsKey])
	configFileStr = strings.ReplaceAll(configFileStr, "@METHOD@", argsMap[v1alpha1.MethodArgsKey])
	configFileStr = strings.ReplaceAll(configFileStr, "@BODY@", strings.ReplaceAll(argsMap[v1alpha1.BodyArgsKey], "\"", "&quot;"))
	headerMap, _ := v1alpha1.GetHeaderMap(argsMap[v1alpha1.HeaderArgsKey])

	var headerConfigStrList []string
	for k, v := range headerMap {
		headerConfigStrList = append(headerConfigStrList,
			fmt.Sprintf("\n                            <elementProp name=\"\" elementType=\"Header\">\n                              <stringProp name=\"Header.name\">%s</stringProp>\n                              <stringProp name=\"Header.value\">%s</stringProp>\n                            </elementProp>", k, v))
	}

	configFileStr = strings.ReplaceAll(configFileStr, "@ELEMENT_PROP@", strings.Join(headerConfigStrList, "\n"))
	return configFileStr
}

func loadJob(ctx context.Context, ins *v1alpha1.LoadTest, configFileStr string) (*batchv1.Job, error) {
	mainConfig := config.GetGlobalConfig()
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("read config: image: %s, cpu: %s, mem: %s", mainConfig.Executor.Image, mainConfig.Executor.Resource.CPU, mainConfig.Executor.Resource.Memory))

	yamlStr := strings.ReplaceAll(v1alpha1.JobYamlStr, "@INITIAL_CONFIG@", configFileStr)

	cpuCore := strconv.Itoa(ins.Spec.Parallelism / ins.Spec.Source / 2)
	if cpuCore == "0" {
		cpuCore = "0.5"
	}

	if mainConfig.Executor.Resource.CPU != "" && mainConfig.Executor.Resource.CPU != "0" {
		cpuCore = mainConfig.Executor.Resource.CPU
	}

	yamlStr = strings.ReplaceAll(yamlStr, "@CPU_REQ@", cpuCore)
	yamlStr = strings.ReplaceAll(yamlStr, "@MEM_REQ@", mainConfig.Executor.Resource.Memory)

	job := &batchv1.Job{}
	if err := yaml.Unmarshal([]byte(yamlStr), job); err != nil {
		return nil, fmt.Errorf("convert yaml to job instance error: %s", err.Error())
	}

	job.Spec.Template.Spec.Containers[0].Image = mainConfig.Executor.Image
	job.Name = ins.Name
	job.Spec.Template.Name = ins.Name
	job.Namespace = ins.Namespace
	job.Spec.Template.Namespace = ins.Namespace
	replicas := int32(ins.Spec.Source)
	job.Spec.Completions = &replicas
	job.Spec.Parallelism = &replicas

	return job, nil
}

func solveFinalizer(instance *v1alpha1.LoadTest) {
	for index := 0; index < len(instance.ObjectMeta.Finalizers); index++ {
		if instance.ObjectMeta.Finalizers[index] == v1alpha1.FinalizerName {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers[:index], instance.ObjectMeta.Finalizers[index+1:]...)
			return
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoadTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.LoadTest{}).
		Complete(r)
}
