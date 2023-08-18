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

package experiment

import (
	"chaosmeta-platform/util/log"
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"time"
)

type ArgoWorkFlowService interface {
	Get(workflowName string) (*v1alpha1.Workflow, string, error)
	List() (*v1alpha1.WorkflowList, error)
	Create(wf v1alpha1.Workflow) (*v1alpha1.Workflow, error)
	Update(wf v1alpha1.Workflow) (*v1alpha1.Workflow, error)
	Delete(workflowName string) error
	Patch(name string, pt types.PatchType, data []byte) error
	DeleteExpiredList() error
	ListPendingAndRecentWorkflows() (*v1alpha1.WorkflowList, []*v1alpha1.Workflow, error)
}

type argoWorkFlowService struct {
	Config    *rest.Config
	Client    dynamic.Interface
	Namespace string
}

func NewArgoWorkFlowService(config *rest.Config, namespace string) (ArgoWorkFlowService, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &argoWorkFlowService{
		Config: config, Client: client, Namespace: namespace,
	}, nil
}

func (a *argoWorkFlowService) Get(workflowName string) (*v1alpha1.Workflow, string, error) {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	workflow, err := wfClient.Get(context.Background(), workflowName, metav1.GetOptions{})
	if err != nil {
		return nil, "", err
	}
	log.Errorf("Workflow %s status: %s", workflowName, workflow.Status.Phase)
	return workflow, string(workflow.Status.Phase), nil
}

func (a *argoWorkFlowService) List() (*v1alpha1.WorkflowList, error) {
	return versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace).List(context.Background(), metav1.ListOptions{})
}

func (a *argoWorkFlowService) Create(wf v1alpha1.Workflow) (*v1alpha1.Workflow, error) {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	createdWorkflow, err := wfClient.Create(context.Background(), &wf, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	log.Errorf("Workflow %s created", createdWorkflow.Name)
	return createdWorkflow, nil
}

func (a *argoWorkFlowService) Update(wf v1alpha1.Workflow) (*v1alpha1.Workflow, error) {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	updatedWorkflow, err := wfClient.Update(context.Background(), &wf, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	log.Errorf("Workflow %s updated", updatedWorkflow.Name)
	return updatedWorkflow, nil
}

func (a *argoWorkFlowService) Delete(workflowName string) error {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	if err := wfClient.Delete(context.Background(), workflowName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	fmt.Printf("Workflow %s deleted\n", workflowName)
	return nil
}

func (a *argoWorkFlowService) Patch(name string, pt types.PatchType, data []byte) error {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	workflow, err := wfClient.Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	updatedWorkflow, err := wfClient.Patch(context.Background(), workflow.Name, pt, data, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	log.Errorf("Workflow %s patched", updatedWorkflow.Name)
	return nil
}

func (a *argoWorkFlowService) GetAllRunningWorkflows() ([]v1alpha1.Workflow, error) {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)
	listOpts := metav1.ListOptions{FieldSelector: "status.phase = Running"}
	workflowList, err := wfClient.List(context.Background(), listOpts)
	if err != nil {
		return nil, err
	}
	return workflowList.Items, nil
}

func (a *argoWorkFlowService) DeleteExpiredList() error {
	expirationTime := time.Now().AddDate(0, 0, -1)
	workflowList, err := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, workflow := range workflowList.Items {
		if workflow.Status.Phase == v1alpha1.WorkflowSucceeded || workflow.Status.Phase == v1alpha1.WorkflowFailed || workflow.Status.Phase == v1alpha1.WorkflowError {
			if workflow.CreationTimestamp.Time.Before(expirationTime) {
				if err := a.Delete(workflow.Name); err != nil {
					log.Error(err)
					return err
				}
				log.Infof("Deleted expired workflow: %s", workflow.Name)
			}
		}
	}

	return nil
}

func (a *argoWorkFlowService) ListPendingAndRecentWorkflows() (*v1alpha1.WorkflowList, []*v1alpha1.Workflow, error) {
	wfClient := versioned.NewForConfigOrDie(a.Config).ArgoprojV1alpha1().Workflows(a.Namespace)

	pendingWorkflows, err := wfClient.List(context.Background(), metav1.ListOptions{FieldSelector: "status.phase=Pending"})
	if err != nil {
		return nil, nil, err
	}

	//List the Workflows whose status is WorkflowSucceeded and WorkflowFailed within 10 minutes
	//expirationTime := time.Now().Add(-5 * time.Minute)
	recentListOpts := metav1.ListOptions{FieldSelector: "status.phase=Succeeded,status.phase=Failed"}
	recentWorkflows, err := wfClient.List(context.Background(), recentListOpts)
	if err != nil {
		return nil, nil, err
	}

	filteredWorkflows := make([]*v1alpha1.Workflow, 0)
	for _, workflow := range recentWorkflows.Items {
		//if workflow.CreationTimestamp.Time.After(expirationTime) {
		filteredWorkflows = append(filteredWorkflows, &workflow)
		//}
	}

	return pendingWorkflows, filteredWorkflows, nil
}
