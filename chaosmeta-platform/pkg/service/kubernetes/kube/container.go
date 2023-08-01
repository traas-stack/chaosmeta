package kube

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sort"
	"strings"
	"time"
)

type ContainerService interface {
	GetLogDetail(namespace, pod, container string, logSelector *Selection, usePreviousLogs bool) (*LogDetails, error)
	DownLogFile(namespace, pod, container string, usePreviousLogs bool) (io.ReadCloser, error)
	List(namespace, pod string) (*PodContainerList, error)
	ExecWithBase64Cmd(namespace, pod, container, command string, args string) (string, error)
}

var (
	LineIndexNotFound             = -1         // LineIndexNotFound is returned if requested line could not be found
	DefaultDisplayNumLogLines     = 100        // DefaultDisplayNumLogLines returns default number of lines in case of invalid request.
	MaxLogLines               int = 2000000000 // MaxLogLines is a number that will be certainly bigger than any number of logs. Here 2 billion logs is certainly much larger number of log lines than we can handle
)

// Load logs from the beginning or the end of the log file.
// This matters only if the log file is too large to be loaded completely.
const (
	NewestTimestamp = "newest"
	OldestTimestamp = "oldest"

	Beginning = "beginning"
	End       = "end"

	defaultPrefixLen = 20
)

// NewestLogLineId is the reference Id of the newest line.
var NewestLogLineId = LogLineId{LogTimestamp: NewestTimestamp}

// OldestLogLineId is the reference Id of the oldest line.
var OldestLogLineId = LogLineId{LogTimestamp: OldestTimestamp}

// DefaultSelection loads default log view selector that is used in case of invalid request
// Downloads newest DefaultDisplayNumLogLines lines.
var DefaultSelection = &Selection{
	OffsetFrom:      1 - DefaultDisplayNumLogLines,
	OffsetTo:        1,
	ReferencePoint:  NewestLogLineId,
	LogFilePosition: End,
}

// AllSelection returns all logs.
var AllSelection = &Selection{
	OffsetFrom:     -MaxLogLines,
	OffsetTo:       MaxLogLines,
	ReferencePoint: NewestLogLineId,
}

// LogDetails returns representation of log lines
type LogDetails struct {

	// Additional information of the logs e.g. container name, dates,...
	Info LogInfo `json:"info"`

	// Reference point to keep track of the position of all the logs
	Selection `json:"selection"`

	// Actual log lines of this page
	LogLines `json:"logs"`
}

// LogInfo returns meta information about the selected log lines
type LogInfo struct {

	// Pod name.
	PodName string `json:"podName"`

	// The name of the container the logs are for.
	ContainerName string `json:"containerName"`

	// The name of the init container the logs are for.
	InitContainerName string `json:"initContainerName"`

	// Date of the first log line
	FromDate LogTimestamp `json:"fromDate"`

	// Date of the last log line
	ToDate LogTimestamp `json:"toDate"`

	// Some log lines in the middle of the log file could not be loaded, because the log file is too large.
	Truncated bool `json:"truncated"`
}

// Selection of a slice of logs.
// It works just like normal slicing, but indices are referenced relatively to certain reference line.
// So for example if reference line has index n and we want to download first 10 elements in array we have to use
// from -n to -n+10. Setting ReferenceLogLineId the first line will result in standard slicing.
type Selection struct {
	// ReferencePoint is the ID of a line which should serve as a reference point for this selector.
	// You can set it to last or first line if needed. Setting to the first line will result in standard slicing.
	ReferencePoint LogLineId `json:"referencePoint"`
	// First index of the slice relatively to the reference line(this one will be included).
	OffsetFrom int `json:"offsetFrom"`
	// Last index of the slice relatively to the reference line (this one will not be included).
	OffsetTo int `json:"offsetTo"`
	// The log file is loaded either from the beginning or from the end. This matters only if the log file is too
	// large to be handled and must be truncated (to avoid oom)
	LogFilePosition string `json:"logFilePosition"`
}

// LogLineId uniquely identifies a line in logs - immune to log addition/deletion.
type LogLineId struct {
	// timestamp of this line.
	LogTimestamp `json:"timestamp"`
	// in case of timestamp duplicates (rather unlikely) it gives the index of the duplicate.
	// For example if this LogTimestamp appears 3 times in the logs and the line is 1nd line with this timestamp,
	// then line num will be 1 or -3 (1st from beginning or 3rd from the end).
	// If timestamp is unique then it will be simply 1 or -1 (first from the beginning or first from the end, both mean the same).
	LineNum int `json:"lineNum"`
}

// LogLines provides means of selecting log views. Problem with logs is that new logs are constantly added.
// Therefore the number of logs constantly changes and we cannot use normal indexing. For example
// if certain line has index N then it may not have index N anymore 1 second later as logs at the beginning of the list
// are being deleted. Therefore it is necessary to reference log indices relative to some line that we are certain will not be deleted.
// For example line in the middle of logs should have lifetime sufficiently long for the purposes of log visualisation. On average its lifetime
// is equal to half of the log retention time. Therefore line in the middle of logs would serve as a good reference point.
// LogLines allows to get ID of any line - this ID later allows to uniquely identify this line. Also it allows to get any
// slice of logs relatively to certain reference line ID.
type LogLines []LogLine

// LogLine is a single log line that split into timestamp and the actual content.
type LogLine struct {
	Timestamp LogTimestamp `json:"timestamp"`
	Content   string       `json:"content"`
}

// LogTimestamp is a timestamp that appears on the beginning of each log line.
type LogTimestamp string

// SelectLogs returns selected part of LogLines as required by logSelector, moreover it returns IDs of first and last
// of returned lines and the information of the resulting logView.
func (l LogLines) SelectLogs(logSelection *Selection) (LogLines, LogTimestamp, LogTimestamp, Selection, bool) {
	requestedNumItems := logSelection.OffsetTo - logSelection.OffsetFrom
	referenceLineIndex := l.getLineIndex(&logSelection.ReferencePoint)
	if referenceLineIndex == LineIndexNotFound || requestedNumItems <= 0 || len(l) == 0 {
		// Requested reference line could not be found, probably it's already gone or requested no logs. Return no logs.
		return LogLines{}, "", "", Selection{}, false
	}
	fromIndex := referenceLineIndex + logSelection.OffsetFrom
	toIndex := referenceLineIndex + logSelection.OffsetTo
	lastPage := false
	if requestedNumItems > len(l) {
		fromIndex = 0
		toIndex = len(l)
		lastPage = true
	} else if toIndex > len(l) {
		fromIndex -= toIndex - len(l)
		toIndex = len(l)
		lastPage = logSelection.LogFilePosition == Beginning
	} else if fromIndex < 0 {
		toIndex += -fromIndex
		fromIndex = 0
		lastPage = logSelection.LogFilePosition == End
	}

	// set the middle of log array as a reference point, this part of array should not be affected by log deletion/addition.
	newSelection := Selection{
		ReferencePoint:  *l.createLogLineId(len(l) / 2),
		OffsetFrom:      fromIndex - len(l)/2,
		OffsetTo:        toIndex - len(l)/2,
		LogFilePosition: logSelection.LogFilePosition,
	}
	return l[fromIndex:toIndex], l[fromIndex].Timestamp, l[toIndex-1].Timestamp, newSelection, lastPage
}

// getLineIndex returns the index of the line (referenced from beginning of log array) with provided logLineId.
func (l LogLines) getLineIndex(logLineId *LogLineId) int {
	if logLineId == nil || logLineId.LogTimestamp == NewestTimestamp || len(l) == 0 || logLineId.LogTimestamp == "" {
		// if no line id provided return index of last item.
		return len(l) - 1
	} else if logLineId.LogTimestamp == OldestTimestamp {
		return 0
	}
	logTimestamp := logLineId.LogTimestamp

	matchingStartedAt := 0
	matchingStartedAt = sort.Search(len(l), func(i int) bool {
		return l[i].Timestamp >= logTimestamp
	})

	linesMatched := 0
	if matchingStartedAt < len(l) && l[matchingStartedAt].Timestamp == logTimestamp { // match found
		for (matchingStartedAt+linesMatched) < len(l) && l[matchingStartedAt+linesMatched].Timestamp == logTimestamp {
			linesMatched += 1
		}
	}

	var offset int
	if logLineId.LineNum < 0 {
		offset = linesMatched + logLineId.LineNum
	} else {
		offset = logLineId.LineNum - 1
	}
	if 0 <= offset && offset < linesMatched {
		return matchingStartedAt + offset
	}
	return LineIndexNotFound
}

// createLogLineId returns ID of the line with provided lineIndex.
func (l LogLines) createLogLineId(lineIndex int) *LogLineId {
	logTimestamp := l[lineIndex].Timestamp
	// determine whether to use negative or positive indexing
	// check whether last line has the same index as requested line. If so, we can only use positive referencing
	// as more lines may appear at the end.
	// negative referencing is preferred as higher indices disappear later.
	var step int
	if l[len(l)-1].Timestamp == logTimestamp {
		// use positive referencing
		step = 1
	} else {
		step = -1
	}
	offset := step
	for ; 0 <= lineIndex-offset && lineIndex-offset < len(l); offset += step {
		if l[lineIndex-offset].Timestamp != logTimestamp {
			break
		}
	}
	return &LogLineId{
		LogTimestamp: logTimestamp,
		LineNum:      offset,
	}
}

// ToLogLines converts rawLogs (string) to LogLines. Proper log lines start with a timestamp which is chopped off.
// In error cases the server returns a message without a timestamp
func ToLogLines(rawLogs string) LogLines {
	logLines := LogLines{}
	for _, line := range strings.Split(rawLogs, "\n") {
		if line != "" {
			startsWithDate := ('0' <= line[0] && line[0] <= '9') //2017-...
			idx := strings.Index(line, " ")
			if idx > 0 && startsWithDate {
				timestamp := LogTimestamp(line[0:idx])
				content := line[idx+1:]
				logLines = append(logLines, LogLine{Timestamp: timestamp, Content: content})
			} else {
				logLines = append(logLines, LogLine{Timestamp: LogTimestamp("0"), Content: line})
			}
		}
	}
	return logLines
}

var lineReadLimit int64 = 5000

// maximum number of bytes loaded from the apiserver
var byteReadLimit int64 = 500000

// PodContainerList is a list of containers of a pod.
type PodContainerList struct {
	Containers []string `json:"containers"`
}

type containerService struct {
	kubeClient kubernetes.Interface
	config     *rest.Config
}

func NewContainerService(kubeClient kubernetes.Interface, config *rest.Config) ContainerService {
	return &containerService{kubeClient, config}
}

// List GetPodContainers returns containers that a pod has.
func (c *containerService) List(namespace, podID string) (*PodContainerList, error) {
	pod, err := c.kubeClient.CoreV1().Pods(namespace).Get(context.TODO(), podID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containers := &PodContainerList{Containers: make([]string, 0)}

	for _, container := range pod.Spec.Containers {
		containers.Containers = append(containers.Containers, container.Name)
	}

	return containers, nil
}

func (c *containerService) ExecWithBase64Cmd(namespace, pod, container, command string, args string) (string, error) {
	var prefix string
	if len(command) < defaultPrefixLen {
		prefix = command[:]
	} else {
		prefix = command[:defaultPrefixLen]
	}
	cmdBts, err := base64.StdEncoding.DecodeString(command)
	if err != nil {
		err = fmt.Errorf("fail to base64-decode command, caused by: %v", err)
		return "", err
	}
	stdin := bytes.NewBuffer(cmdBts)
	filename := fmt.Sprintf("/tmp/%s_%v.sh", prefix, time.Now().Unix())
	actCmd := fmt.Sprintf("touch %s && cat - > %s && sh %s", filename, filename, filename)
	if args != "" {
		actCmd = fmt.Sprintf("%s %s", actCmd, args)
	}

	podList, err := c.List(namespace, pod)
	if err != nil {
		return "", err
	}

	for _, ct := range podList.Containers {
		if ct == container {
			execOptions := ExecOptions{
				Command:            actCmd,
				Namespace:          namespace,
				PodName:            pod,
				ContainerName:      container,
				Stdin:              stdin,
				CaptureStdout:      true,
				CaptureStderr:      true,
				PreserveWhitespace: false,
				Quiet:              false,
				config:             c.config,
				kubeClient:         c.kubeClient,
			}

			stdout, stderr, err := ExecWithOptions(execOptions)
			if err != nil {
				return "", fmt.Errorf("exec error: %v, stderr: %v", err, stderr)
			}
			if stderr != "" {
				return "", fmt.Errorf(stderr)
			}
			return stdout, nil
		}
	}

	return "", fmt.Errorf("The container %s does not exist. ", container)
}

func (c *containerService) DownLogFile(namespace, podID, container string, usePreviousLogs bool) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   usePreviousLogs,
		Timestamps: false,
	}
	logStream, err := openStream(c.kubeClient, namespace, podID, logOptions)
	return logStream, err
}

func (c *containerService) GetLogDetail(namespace, pod, container string, logSelector *Selection, usePreviousLogs bool) (*LogDetails, error) {
	podInfo, err := c.kubeClient.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(container) == 0 {
		container = podInfo.Spec.Containers[0].Name
	}

	logOptions := mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := readRawLogs(c.kubeClient, namespace, pod, logOptions)
	if err != nil {
		return nil, err
	}
	details := ConstructLogDetails(pod, rawLogs, container, logSelector)
	return details, nil
}

func mapToLogOptions(container string, logSelector *Selection, previous bool) *v1.PodLogOptions {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   previous,
		Timestamps: true,
	}

	if logSelector.LogFilePosition == Beginning {
		logOptions.LimitBytes = &byteReadLimit
	} else {
		logOptions.TailLines = &lineReadLimit
	}

	return logOptions
}

// Construct a request for getting the logs for a pod and retrieves the logs.
func readRawLogs(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := openStream(client, namespace, podID, logOptions)
	if err != nil {
		return err.Error(), nil
	}

	defer readCloser.Close()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func openStream(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream(context.TODO())
}

// ConstructLogDetails creates a new log details structure for given parameters.
func ConstructLogDetails(podID string, rawLogs string, container string, logSelector *Selection) *LogDetails {
	parsedLines := ToLogLines(rawLogs)
	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)

	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
	truncated := readLimitReached && lastPage

	info := LogInfo{
		PodName:       podID,
		ContainerName: container,
		FromDate:      fromDate,
		ToDate:        toDate,
		Truncated:     truncated,
	}
	return &LogDetails{
		Info:      info,
		Selection: logSelection,
		LogLines:  logLines,
	}
}

// Checks if the amount of log file returned from the apiserver is equal to the read limits
func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == End && linesLoaded >= lineReadLimit)
}
