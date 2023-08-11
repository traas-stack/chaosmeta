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

package kube

import (
	"archive/tar"
	"chaosmeta-platform/pkg/models/common/page"
	"chaosmeta-platform/pkg/service/kubernetes"
	"chaosmeta-platform/util/json"
	"chaosmeta-platform/util/log"
	"context"
	"fmt"
	"github.com/panjf2000/ants"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	_ "unsafe"
)

type PodService interface {
	List(namespace string, dsQuery *page.DataSelectQuery) (*PodResponse, error)
	ListWithOptions(namespace string, labelSelector *metav1.LabelSelector) ([]corev1.Pod, error)
	Get(namespace, name string) (*PodDetail, error)
	GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error)
	Patch(originalObj, updatedObj *corev1.Pod) error
	Delete(namespace, name string) error
	Exec(namespace, name, command string) (string, string, error)
	GetByIP(ip string) (*corev1.Pod, error)
	PodBelongTo(item *corev1.Pod, kind, name string) bool
	Check(namespace, name string)
	ListPodsByNode(namespace string, nodeName string) ([]corev1.Pod, error)

	CopyToPod(namespace, name, container string, srcPath, destPath string) error
	CopyFromPod(namespace, name, container string, srcPath, destPath string) error

	ListPodByNames(ctx context.Context, space string, names []string) ([]*corev1.Pod, error)
	ListPodByIPs(ctx context.Context, ips []string) ([]*corev1.Pod, error)
	ListPodBySelector(ctx context.Context, namespace string, selectors metav1.ListOptions) ([]corev1.Pod, error)
	ArePodsExistByIPs(ctx context.Context, ips []string) (bool, error)

	LogFilePaths(namespace, name, container, filePath string) ([]string, error)
	LogPreview(namespace, name, container, file string, offset, size int) (string, error)
}

type podService struct {
	kubernetesParam *kubernetes.KubernetesParam
}

func NewPodService(kubernetesParam *kubernetes.KubernetesParam) PodService {
	return &podService{kubernetesParam: kubernetesParam}
}

type PodResponse struct {
	Total    int         `json:"total"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	List     []PodDetail `json:"list"`
}

type PodDetail struct {
	corev1.Pod `json:",inline"`
	PodPhase   string `json:"podPhase"`
}

type PodCell corev1.Pod

func (n PodCell) GetProperty(name page.PropertyName) page.ComparableValue {
	switch name {
	case page.NameProperty:
		return page.StdComparableString(n.ObjectMeta.Name)
	case page.CreationTimestampProperty:
		return page.StdComparableTime(n.ObjectMeta.CreationTimestamp.Time)
	case page.NamespaceProperty:
		return page.StdComparableString(n.ObjectMeta.Namespace)
	case page.PodIpProperty:
		return page.StdComparableString(n.Status.PodIP)
	default:
		return nil
	}
}

func ToCells(std []corev1.Pod) []page.DataCell {
	cells := make([]page.DataCell, len(std))
	for i := range std {
		cells[i] = PodCell(std[i])
	}
	return cells
}

func FromCells(cells []page.DataCell) []corev1.Pod {
	std := make([]corev1.Pod, len(cells))
	for i := range std {
		std[i] = corev1.Pod(cells[i].(PodCell))
	}
	return std
}

func hasPodReadyCondition(conditions []corev1.PodCondition) bool {
	for _, condition := range conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getPodStatus(pod corev1.Pod) string {
	restarts := 0
	readyContainers := 0

	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init: Signal %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init: ExitCode %d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = fmt.Sprintf("Init: %s", container.State.Waiting.Reason)
			initializing = true
		default:
			reason = fmt.Sprintf("Init: %d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		restarts = 0
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			restarts += int(container.RestartCount)
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal: %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode: %d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
				readyContainers++
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = string(corev1.PodRunning)
			} else {
				reason = "NotReady"
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		reason = string(corev1.PodUnknown)
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	if len(reason) == 0 {
		reason = string(corev1.PodUnknown)
	}

	return reason
}

func (p *podService) List(namespace string, dsQuery *page.DataSelectQuery) (*PodResponse, error) {
	var podResponse PodResponse
	pods, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	var podList []corev1.Pod
	for _, pp := range pods.Items {
		podList = append(podList, pp)
	}

	if err != nil {
		return nil, err
	}
	podCells, filteredTotal := page.GenericDataSelectWithFilter(ToCells(podList), dsQuery)
	ps := FromCells(podCells)

	var podDetailList []PodDetail

	for _, po := range ps {
		var detail PodDetail
		detail.Pod = po
		detail.PodPhase = getPodStatus(po)
		podDetailList = append(podDetailList, detail)
	}

	podResponse.List = podDetailList
	podResponse.Current = dsQuery.PaginationQuery.Page + 1
	podResponse.PageSize = dsQuery.PaginationQuery.ItemsPerPage
	podResponse.Total = filteredTotal
	return &podResponse, nil
}

func (p *podService) ListWithOptions(namespace string, labelSelector *metav1.LabelSelector) ([]corev1.Pod, error) {
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}
	pods, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (p *podService) ListPodsByNode(namespace string, nodeName string) ([]corev1.Pod, error) {
	var podList []corev1.Pod
	pods, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, tmp := range pods.Items {
		if tmp.Spec.NodeName == nodeName {
			podList = append(podList, tmp)
		}
	}

	return podList, nil
}

func (p *podService) Get(namespace, name string) (*PodDetail, error) {
	pod, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &PodDetail{PodPhase: getPodStatus(*pod), Pod: *pod}, nil
}

func (p *podService) Delete(namespace, name string) error {
	return p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (p *podService) GetEvents(namespace, name string, dsQuery *page.DataSelectQuery) (*EventResponse, error) {
	eventCtrl := NewEventService(p.kubernetesParam.KubernetesClient)
	eventResponse, err := eventCtrl.GetResourceEvents(namespace, name, dsQuery)
	if err != nil {
		return nil, err
	}
	return eventResponse, nil
}

func (p *podService) Patch(originalObj, updatedObj *corev1.Pod) error {
	updatedObj.ObjectMeta = originalObj.ObjectMeta

	data, err := json.MargePatch(originalObj, updatedObj)
	if err != nil {
		return err
	}

	_, err = p.kubernetesParam.KubernetesClient.CoreV1().Pods(originalObj.GetNamespace()).Patch(
		context.TODO(),
		originalObj.Name,
		types.MergePatchType,
		data,
		metav1.PatchOptions{},
	)
	return err
}

func (p *podService) GetByIP(ip string) (*corev1.Pod, error) {
	podList, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("status.podIP=%s", ip),
	})
	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, fmt.Errorf("The pod %s does not exist. ", ip)
	}

	return &podList.Items[0], nil
}

func (p *podService) Exec(namespace, name, command string) (string, string, error) {
	pod, err := p.Get(namespace, name)
	if err != nil {
		return "", "", err
	}

	execOptions := ExecOptions{
		Command:            command,
		Namespace:          namespace,
		PodName:            name,
		ContainerName:      pod.Spec.Containers[0].Name,
		Stdin:              nil,
		CaptureStdout:      true,
		CaptureStderr:      true,
		PreserveWhitespace: false,
		Quiet:              false,
		config:             p.kubernetesParam.RestConfig,
		kubeClient:         p.kubernetesParam.KubernetesClient,
	}

	stdout, _, err := ExecWithOptions(execOptions)
	if err != nil {
		return pod.Spec.Containers[0].Name, "", err
	}
	return pod.Spec.Containers[0].Name, stdout, nil
}

// Check whether the pod settings are compliant
func (p *podService) Check(namespace, name string) {
}

func (p *podService) ExecByContainer(namespace, name, container string, cmd string) (string, error) {
	_, err := p.Get(namespace, name)
	if err != nil {
		return "", err
	}

	execOptions := ExecOptions{
		Command:            cmd,
		Namespace:          namespace,
		PodName:            name,
		ContainerName:      container,
		Stdin:              nil,
		CaptureStdout:      true,
		CaptureStderr:      true,
		PreserveWhitespace: false,
		Quiet:              false,
		config:             p.kubernetesParam.RestConfig,
		kubeClient:         p.kubernetesParam.KubernetesClient,
	}

	stdout, stderr, err := ExecWithOptions(execOptions)
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf(stderr)
		}
		return "", err
	}
	if stderr != "" {
		return "", fmt.Errorf(stderr)
	}
	return stdout, nil
}

func (p *podService) checkDestinationIsDir(namespace, name, container, destPath string) error {
	_, err := p.ExecByContainer(namespace, name, container, fmt.Sprintf("test -d %s", destPath))
	return err
}

// copy from kubectl cp
func makeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}

// Upload files from local to pod
func (p *podService) CopyToPod(namespace, name, container string, srcPath, destPath string) error {

	reader, writer := io.Pipe()
	if destPath != "/" && strings.HasSuffix(string(destPath[len(destPath)-1]), "/") {
		destPath = destPath[:len(destPath)-1]
	}
	if err := p.checkDestinationIsDir(namespace, name, container, destPath); err == nil {
		destPath = destPath + "/" + path.Base(srcPath)
	}
	go func() {
		defer writer.Close()
		err := makeTar(srcPath, destPath, writer)
		cmdutil.CheckErr(err)
	}()
	var cmdArr []string

	cmdArr = []string{"tar", "-xf", "-"}
	destDir := path.Dir(destPath)
	if len(destDir) > 0 {
		cmdArr = append(cmdArr, "-C", destDir)
	}
	//remote shell.
	req := p.kubernetesParam.KubernetesClient.CoreV1().RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.kubernetesParam.RestConfig, "POST", req.URL())
	if err != nil {
		log.Fatalf("error %v
", err)
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		log.Fatalf("error %v
", err)
		return err
	}
	return nil
}

// Copy files from pod to local
func (p *podService) CopyFromPod(namespace, name, container string, srcPath, destPath string) error {
	// 从pod内copy文件到本地
	reader, outStream := io.Pipe()
	cmdArr := []string{"tar", "cf", "-", srcPath}
	req := p.kubernetesParam.KubernetesClient.CoreV1().RESTClient().
		Get().
		Namespace(namespace).
		Resource("pods").
		Name(name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.kubernetesParam.RestConfig, "POST", req.URL())
	if err != nil {
		log.Fatalf("error %s
", err)
		return err
	}
	go func() {
		defer outStream.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
		cmdutil.CheckErr(err)
	}()
	prefix := getPrefix(srcPath)
	prefix = path.Clean(prefix)
	prefix = stripPathShortcuts(prefix)
	destPath = path.Join(destPath, path.Base(prefix))
	err = untarAll(reader, destPath, prefix)
	return err
}

// unzip
func untarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// get prefix
func getPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

// stripPathShortcuts
// @Description: copy from kubectl
// @param p
// @return string
func stripPathShortcuts(p string) string {
	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}

func (p *podService) PodBelongTo(item *corev1.Pod, kind string, name string) bool {
	switch kind {
	case "Deployment":
		if p.podBelongToDeployment(item, name) {
			return true
		}
	case "ReplicaSet":
		if p.podBelongToReplicaSet(item, name) {
			return true
		}
	case "DaemonSet":
		if p.podBelongToDaemonSet(item, name) {
			return true
		}
	case "StatefulSet":
		if p.podBelongToStatefulSet(item, name) {
			return true
		}
	case "Job":
		if p.podBelongToJob(item, name) {
			return true
		}
	}
	return false
}

func (p *podService) podBelongToDeployment(item *corev1.Pod, deploymentName string) bool {
	replicas, err := p.kubernetesParam.KubernetesClient.AppsV1().ReplicaSets(item.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false
	}

	for _, r := range replicas.Items {
		if p.replicaSetBelongToDeployment(&r, deploymentName) && p.podBelongToReplicaSet(item, r.Name) {
			return true
		}
	}

	return false
}

func (p *podService) replicaSetBelongToDeployment(replicaSet *appsv1.ReplicaSet, deploymentName string) bool {
	for _, owner := range replicaSet.OwnerReferences {
		if owner.Kind == "Deployment" && owner.Name == deploymentName {
			return true
		}
	}
	return false
}

func (p *podService) podBelongToReplicaSet(item *corev1.Pod, replicaSetName string) bool {
	for _, owner := range item.OwnerReferences {
		if owner.Kind == "ReplicaSet" && owner.Name == replicaSetName {
			return true
		}
	}
	return false
}

func (p *podService) podBelongToStatefulSet(item *corev1.Pod, statefulSetName string) bool {
	for _, owner := range item.OwnerReferences {
		if owner.Kind == "StatefulSet" && owner.Name == statefulSetName {
			return true
		}
	}
	return false
}

func (p *podService) podBelongToDaemonSet(item *corev1.Pod, name string) bool {
	for _, owner := range item.OwnerReferences {
		if owner.Kind == "DaemonSet" && owner.Name == name {
			return true
		}
	}
	return false
}

func (p *podService) podBelongToJob(item *corev1.Pod, name string) bool {
	for _, owner := range item.OwnerReferences {
		if owner.Kind == "Job" && owner.Name == name {
			return true
		}
	}
	return false
}

func (p *podService) ListPodByNames(ctx context.Context, space string, names []string) (pods []*corev1.Pod, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error when list pods by names|%v", err)
			return
		}
	}()

	var (
		wg  sync.WaitGroup
		gp  *ants.Pool
		lth = len(names)
		ech = make(chan error, lth)
		pch = make(chan *corev1.Pod, lth)
	)
	gp, err = ants.NewPool(20)
	if err != nil {
		err = fmt.Errorf("fail to new goroutine pool, caused by: %v", err)
		return
	}
	defer gp.Release()

	for _, n := range names {
		wg.Add(1)
		n := n
		err = gp.Submit(func() {
			defer wg.Done()

			p, grr := p.Get(space, n)
			if grr != nil {
				grr = fmt.Errorf("fail to get pod: %v, caused by: %v", n, grr)
				ech <- grr
				return
			}
			pch <- &p.Pod
		})
		if err != nil {
			err = fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err)
			return
		}
	}
	wg.Wait()
	close(pch)
	select {
	case err = <-ech:
		return
	default:
		// do nothing
	}

	for p := range pch {
		pods = append(pods, p)
	}

	return
}

func (p *podService) ListPodByIPs(ctx context.Context, ips []string) (pods []*corev1.Pod, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("fail when listing pods by ips|%v", err)
			return
		}
	}()

	var (
		wg  sync.WaitGroup
		gp  *ants.Pool
		lth = len(ips)
		ech = make(chan error, lth)
		pch = make(chan *corev1.Pod, lth)
	)
	gp, err = ants.NewPool(20)
	if err != nil {
		err = fmt.Errorf("fail to new goroutine pool, caused by: %v", err)
		return
	}
	defer gp.Release()

	for _, ip := range ips {
		wg.Add(1)
		ip := ip
		err = gp.Submit(func() {
			defer wg.Done()

			p, e := p.GetByIP(ip)
			if e != nil {
				e = fmt.Errorf("fail to get pod: %v, caused by: %v", ip, e)
				ech <- e
				return
			}
			pch <- p
		})
		if err != nil {
			err = fmt.Errorf("fail to add task to goroutine pool, caused by: %v", err)
			return
		}
	}
	wg.Wait()
	close(pch)
	select {
	case err = <-ech:
		return
	default:
		// do nothing
	}

	for p := range pch {
		pods = append(pods, p)
	}

	return
}

func (p *podService) ListPodBySelector(ctx context.Context, namespace string, selector metav1.ListOptions) ([]corev1.Pod, error) {
	pods, err := p.kubernetesParam.KubernetesClient.CoreV1().Pods(namespace).List(ctx, selector)
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (p *podService) ArePodsExistByIPs(ctx context.Context, ips []string) (exist bool, err error) {
	if len(ips) == 0 {
		return
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("error when check pods exist by ips|%v", err)
			return
		}
	}()

	pods, err := p.ListPodByIPs(ctx, ips)
	if err != nil {
		return
	}
	if len(pods) != len(ips) {
		err = fmt.Errorf("some of the ips not exist")
		return
	}

	return
}

// Get all files and folders under a certain path of a pod
// similar to ls /
func (p *podService) LogFilePaths(namespace, name, container, filePath string) ([]string, error) {
	command := fmt.Sprintf("ls %s", filePath)
	out, err := p.ExecByContainer(namespace, name, container, command)
	if err != nil {
		return nil, err
	}
	return strings.Split(out, "
"), nil
}

func (p *podService) LogPreview(namespace, name, container, file string, offset, size int) (string, error) {
	if _, err := p.LogFilePaths(namespace, name, container, file); err != nil {
		return "", err
	}

	// Starting from the offset line, the actual size line
	cmd := fmt.Sprintf("cat %s |tail -n +%d|head -n %d", file, offset, size)
	out, err := p.ExecByContainer(namespace, name, container, cmd)
	if err != nil {
		return "", err
	}
	return out, nil
}
