package tse

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/config"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/middlewareexecutor/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor/middlewareexecutor/tse/auth"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

var (
	cmdPrefix = "cmd://"
)

type TseMiddleware struct {
	Config     config.MiddlewareConfig
	MistClient auth.MistClient
	tseUrl     string
}

type TseTaskRequest struct {
	Ip        string `json:"ip"`
	Exeurl    string `json:"exeurl"`
	Timestamp int64  `json:"timestamp"`
	Key       string `json:"key"`
	Sign      string `json:"sign"`
	User      string `json:"user"`
}

type TseTaskResponse struct {
	UID       string
	SUCCESS   bool
	ERRORCODE string
	ERRORMSG  string
	JOBNAME   string
	JOBRESULT string
	IP        string
}

func (r *TseMiddleware) ExecCmdTask(ctx context.Context, host string, cmd string) common.TaskResult {
	errResult := common.TaskResult{
		Success: false,
		Message: "fail to exec",
	}
	requestBody := &TseTaskRequest{}
	requestBody.Ip = host
	requestBody.Exeurl = cmdPrefix + cmd
	requestBody.Timestamp = time.Now().Unix()
	ak, sk := r.MistClient.MistConfig()
	requestBody.Key = ak
	requestBody.User = "root"
	requestBody.Sign = hmacSha1(requestBody, sk)
	requestBodyStr, err := json.Marshal(requestBody)
	if err != nil {
		return errResult
	}
	execTaskUrl := fmt.Sprintf("%s/api/task", r.tseUrl)
	req, err := http.NewRequest("POST", execTaskUrl, bytes.NewBuffer(requestBodyStr))
	if err != nil {
		return errResult
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return errResult
	}
	result := &common.TaskResult{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return errResult
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return errResult
	}
	return common.TaskResult{}
}

func hmacSha1(tseTaskRequest *TseTaskRequest, key string) string {
	hmacSha1Res := ""
	argsMap := make(map[string]string)
	argsMap["ip"] = tseTaskRequest.Ip
	argsMap["exeurl"] = tseTaskRequest.Exeurl
	argsMap["timestamp"] = fmt.Sprint(tseTaskRequest.Timestamp)
	argsMap["key"] = tseTaskRequest.Key
	argsMap["user"] = tseTaskRequest.User
	var keys []string
	for k := range argsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		hmacSha1Res += k + argsMap[k]
	}
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(hmacSha1Res))
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum)
}

func (r *TseMiddleware) QueryTaskStatus(ctx context.Context, taskId string, userKey string) common.TaskResult {
	return common.TaskResult{}
}
