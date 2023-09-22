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

package daemonsetexecutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/base"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/model"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

type DaemonsetRemoteExecutor struct {
	//ApiServer  rest.Interface
	RESTConfig *rest.Config
	Schema     *runtime.Scheme

	LocalExecPath string
	Executor      string
	Version       string
	//DaemonsetPolicy   DaemonsetPolicyType
	DaemonsetNs    string
	DaemonsetLabel map[string]string

	//AutoLabelNode     bool
	//NodeSelectorLabel map[string]string
}

func (r *DaemonsetRemoteExecutor) CheckAlive(ctx context.Context, injectObject string) error {
	agentPod, err := r.getAgentPod(ctx, injectObject)
	if err != nil {
		return fmt.Errorf("get agent pod of node[%s] error: %s", injectObject, err.Error())
	}

	executor := fmt.Sprintf("%s/%s-%s/%s", r.LocalExecPath, r.Executor, r.Version, r.Executor)
	executeCmd := fmt.Sprintf("nsenter -t 1 -m -u %s version", executor)

	var stdout []byte
	stdout, err = r.kubeExec(ctx, agentPod.Namespace, agentPod.PodName, executeCmd)
	if err != nil {
		return fmt.Errorf("kubectl exec error: %s", err.Error())
	}

	var res base.VersionInfo
	if err := json.Unmarshal(stdout, &res); err != nil {
		return fmt.Errorf("version output [%s] is not json format: %s", string(stdout), err.Error())
	}

	if res.Version != r.Version {
		return fmt.Errorf("expected version %s, but get %s", r.Version, res.Version)
	}

	return nil
}

// Init install agent
func (r *DaemonsetRemoteExecutor) Init(ctx context.Context, target string) error {
	return nil
}

func (r *DaemonsetRemoteExecutor) Inject(ctx context.Context, injectObject string, target, fault, uid, timeout, cID, cRuntime string, args []v1alpha1.ArgsUnit) error {
	//if err := r.CheckAlive(ctx, injectObject); err != nil {
	//	if !r.AutoLabelNode {
	//		return fmt.Errorf("check target's status error: %s", err.Error())
	//	} else {
	//	}
	//}

	agentPod, err := r.getAgentPod(ctx, injectObject)
	if err != nil {
		return fmt.Errorf("get agent pod of node[%s] error: %s", injectObject, err.Error())
	}

	executor := fmt.Sprintf("%s/%s-%s/%s", r.LocalExecPath, r.Executor, r.Version, r.Executor)
	executeCmd := fmt.Sprintf("nsenter -t 1 -m -u %s inject %s %s --uid %s", executor, target, fault, uid)
	for _, unitArgs := range args {
		if unitArgs.Key == v1alpha1.ContainerKey {
			continue
		}

		unitArgs.Key = strings.ReplaceAll(unitArgs.Key, "_", "-")
		executeCmd = fmt.Sprintf("%s --%s=%s", executeCmd, unitArgs.Key, unitArgs.Value)
	}

	if timeout != "" {
		executeCmd = fmt.Sprintf("%s --timeout %s", executeCmd, timeout)
	}

	if cRuntime != "" {
		executeCmd = fmt.Sprintf("%s --container-runtime %s --container-id %s", executeCmd, cRuntime, cID)
	}

	if _, err = r.kubeExec(ctx, agentPod.Namespace, agentPod.PodName, executeCmd); err != nil {
		return fmt.Errorf("kubectl exec error: %s", err.Error())
	}

	return nil
}

func (r *DaemonsetRemoteExecutor) Recover(ctx context.Context, injectObject string, uid string) error {
	agentPod, err := r.getAgentPod(ctx, injectObject)
	if err != nil {
		return fmt.Errorf("get agent pod of node[%s] error: %s", injectObject, err.Error())
	}

	executor := fmt.Sprintf("%s/%s-%s/%s", r.LocalExecPath, r.Executor, r.Version, r.Executor)
	executeCmd := fmt.Sprintf("nsenter -t 1 -m -u %s recover %s", executor, uid)

	if _, err = r.kubeExec(ctx, agentPod.Namespace, agentPod.PodName, executeCmd); err != nil {
		return fmt.Errorf("kubectl exec error: %s", err.Error())
	}

	return nil
}

func (r *DaemonsetRemoteExecutor) Query(ctx context.Context, injectObject string, uid string, phase v1alpha1.PhaseType) (*model.SubExpInfo, error) {
	agentPod, err := r.getAgentPod(ctx, injectObject)
	if err != nil {
		return nil, fmt.Errorf("get agent pod of node[%s] error: %s", injectObject, err.Error())
	}

	executor := fmt.Sprintf("%s/%s-%s/%s", r.LocalExecPath, r.Executor, r.Version, r.Executor)
	executeCmd := fmt.Sprintf("nsenter -t 1 -m -u %s query -u %s --format json", executor, uid)

	var stdout []byte
	stdout, err = r.kubeExec(ctx, agentPod.Namespace, agentPod.PodName, executeCmd)
	if err != nil {
		return nil, fmt.Errorf("kubectl exec error: %s", err.Error())
	}

	var res base.QueryResponseData
	if err := json.Unmarshal(stdout, &res); err != nil {
		return nil, fmt.Errorf("query output [%s] is not json format: %s", string(stdout), err.Error())
	}

	if res.Total != 1 {
		return nil, fmt.Errorf("query output expect 1 but get: %d", res.Total)
	}

	return &model.SubExpInfo{
		UID:        uid,
		CreateTime: res.Experiments[0].CreateTime,
		UpdateTime: res.Experiments[0].UpdateTime,
		Message:    res.Experiments[0].Error_,
		Status:     base.ConvertStatus(res.Experiments[0].Status, phase),
	}, nil
}

func (r *DaemonsetRemoteExecutor) kubeExec(ctx context.Context, ns, podName, cmd string) ([]byte, error) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("%s/%s,exec: %s", ns, podName, cmd))

	execReq := restclient.GetApiServerClientMap(v1alpha1.PodCloudTarget).Post().
		Namespace(ns).Resource("pods").Name(podName).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			//Container: cName,
			Command: []string{"/bin/bash", "-c", cmd},
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
		}, runtime.NewParameterCodec(r.Schema))

	exec, err := remotecommand.NewSPDYExecutor(r.RESTConfig, "POST", execReq.URL())
	if err != nil {
		return nil, fmt.Errorf("create remote cmd executor error: %s", err.Error())
	}

	var stdout, stderr bytes.Buffer
	if err := exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		// TODO: Think about how to solve the basic error information collection of parameter operation
		return nil, fmt.Errorf("exec remote cmd error: %s %s %s", err.Error(), stdout.String(), stderr.String())
	}

	if stderr.String() != "" {
		return stdout.Bytes(), fmt.Errorf("exec remote cmd get error message: %s", stderr.String())
	}

	return stdout.Bytes(), nil
}

func (r *DaemonsetRemoteExecutor) getAgentPod(ctx context.Context, nodeIp string) (*model.PodObject, error) {
	// 用目标injectObject获取所在的DaemonSet pod name即可
	podList, err := selector.GetAnalyzer().GetPodListByLabelInNode(ctx, r.DaemonsetNs, r.DaemonsetLabel, nodeIp)
	if err != nil {
		return nil, err
	}

	if len(podList) != 1 {
		return nil, fmt.Errorf("length of agent pod is not 1")
	}

	return podList[0], nil
}
