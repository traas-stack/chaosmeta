package kube

import (
	"bytes"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/url"
	"strings"
)

type ExecOptions struct {
	Command string

	Namespace     string
	PodName       string
	ContainerName string

	Stdin              io.Reader
	CaptureStdout      bool
	CaptureStderr      bool
	PreserveWhitespace bool
	Quiet              bool

	config     *rest.Config
	kubeClient kubernetes.Interface
}

// ExecWithOptions executes a command in the specified container,
// returning stdout, stderr and error. `options` allowed for
// additional parameters to be passed.
func ExecWithOptions(options ExecOptions) (string, string, error) {
	const tty = false
	req := options.kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(options.PodName).
		Namespace(options.Namespace).
		SubResource("exec").
		Param("container", options.ContainerName)
	req.VersionedParams(&v1.PodExecOptions{
		Container: options.ContainerName,
		Command:   []string{"sh", "-c", options.Command},
		Stdin:     options.Stdin != nil,
		Stdout:    options.CaptureStdout,
		Stderr:    options.CaptureStderr,
		TTY:       tty,
	}, scheme.ParameterCodec)

	var stdout, stderr bytes.Buffer
	err := execute("POST", req.URL(), options.config, options.Stdin, &stdout, &stderr, tty)

	if options.PreserveWhitespace {
		return stdout.String(), stderr.String(), err
	}
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func execute(method string, url *url.URL, config *rest.Config, stdin io.Reader, stdout, stderr io.Writer, tty bool) error {
	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    tty,
	})
}
