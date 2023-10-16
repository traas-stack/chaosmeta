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

package inject

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/pkg/models/inject/basic"
	"chaosmeta-platform/util/log"
	"context"
	"fmt"
)

func InitMeasure() error {
	measure := basic.MeasureInject{}
	if _, err := models.GetORM().Raw(fmt.Sprintf("TRUNCATE TABLE %s", measure.TableName())).Exec(); err != nil {
		return err
	}
	ctx := context.Background()
	if err := InitMonitorMeasure(ctx); err != nil {
		log.Error(err)
		return err
	}
	if err := InitTcpMeasure(ctx); err != nil {
		log.Error(err)
		return err
	}
	if err := InitPodMeasure(ctx); err != nil {
		log.Error(err)
		return err
	}
	if err := InitHTTPMeasure(ctx); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func initMeasureCommon(ctx context.Context, injectId int, measureType string) error {
	argMeasureType := basic.Args{InjectId: injectId, ExecType: ExecMeasureCommon, Key: "measureType", KeyCn: "measureType", ValueType: "string", ValueRule: measureType, DescriptionCn: "度量操作类型", Description: "Metric operation type"}
	argDuration := basic.Args{InjectId: injectId, ExecType: ExecMeasureCommon, Key: "duration", KeyCn: "duration", ValueType: "string", Unit: "s", UnitCn: "s", DescriptionCn: "度量任务持续时间", Description: "Measure task duration"}
	argsInterval := basic.Args{InjectId: injectId, ExecType: ExecMeasureCommon, Key: "interval", KeyCn: "interval", ValueType: "string", Unit: "s,m,h", UnitCn: "s,m,h", DescriptionCn: "度量操作执行间隔", Description: "Measuring operation execution interval"}
	argsSuccessCount := basic.Args{InjectId: injectId, ExecType: ExecMeasureCommon, Key: "successCount", KeyCn: "successCount", ValueType: "int", ValueRule: ">0", DescriptionCn: "成功次数阈值,度量任务结束时,度量成功次数不小于successCount,则CR任务状态为success", Description: "Threshold of the number of successes. When the measurement task ends, if the number of successes is not less than successCount, then the CR task status is success."}
	argsFailedCount := basic.Args{InjectId: injectId, ExecType: ExecMeasureCommon, Key: "failedCount", KeyCn: "failedCount", ValueType: "int", ValueRule: ">=0", DescriptionCn: "失败次数阈值,度量任务结束时,度量失败次数不小于failedCount，则CR任务状态为failed", Description: "Failure count threshold. When the measurement task ends, if the number of measurement failures is not less than failedCount, then the CR task status is failed."}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argMeasureType, &argDuration, &argsInterval, &argsSuccessCount, &argsFailedCount})
}

func InitMonitorMeasure(ctx context.Context) error {
	var (
		monitorMeasure = basic.MeasureInject{MeasureType: "monitor", Name: "monitor", NameCn: "monitor", Description: "Make expected judgments on the values of monitoring items, such as whether the CPU usage monitoring value of a certain machine is greater than 90%. Prometheus is supported by default", DescriptionCn: "对监控项的值进行预期判断,比如某个机器的cpu使用率监控值是否大于90%,默认支持prometheus"}
	)
	if err := basic.InsertMeasureInject(&monitorMeasure); err != nil {
		return err
	}
	return InitMonitorMeasureArgs(ctx, monitorMeasure)
}

func initMonitorMeasureJudge(ctx context.Context, measureInject basic.MeasureInject) error {
	argsJudgeType := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeType", KeyCn: "预期判断方式", ValueType: "string", DescriptionCn: "1.absolutevalue(绝对值比较)\njudgeValue样例:\n\"75,85\"：75<=监控值<=85\n\"80,\"：80<=监控值\n\",80\"：监控值<=80\n\"80\"：监控值 == 80\n2.relativevalue(相对值比较,以度量任务开始时的第一个监控值作为初始值：initial)\njudgeValue样例:\n\"20,30\"：initial+20<=监控值<=initial+30\n\"-30,-20\"：initial-30<=监控值<=initial-20\n\",-20\"：监控值<=initial-20\n\"20\"：initial+20 == 监控值\n3.relativepercent(相对百分比比较,以度量任务开始时的第一个监控值作为初始值：initial)\njudgeValue样例:\n\"5,10\"：initial(1+5%)<=监控值<=initial(1+10%)\n\"-10,-5\"：initial(1-10%)<=监控值<=initial(1-5%)\n\"5,\"：initial(1+5%) <= 监控值\n\"-5\"：监控值 == initial(1-5%)", Description: "1.absolutevalue (absolute value comparison)\njudgeValue example:\n\"75,85\": 75<=monitoring value<=85\n\"80,\"：80<=monitoring value\n\",80\": monitoring value <=80\n\"80\": monitoring value == 80\n2.relativevalue (relative value comparison, taking the first monitoring value at the beginning of the measurement task as the initial value: initial)\njudgeValue example:\n\"20,30\": initial+20<=monitoring value<=initial+30\n\"-30,-20\": initial-30<=monitoring value<=initial-20\n\",-20\": monitoring value<=initial-20\n\"20\": initial+20 == monitoring value\n3.relativepercent (relative percentage comparison, taking the first monitoring value at the beginning of the measurement task as the initial value: initial)\njudgeValue example:\n\"5,10\": initial(1+5%)<=monitoring value<=initial(1+10%)\n\"-10,-5\": initial(1-10%)<=monitoring value<=initial(1-5%)\n\"5,\": initial(1+5%) <= monitoring value\n\"-5\": monitoring value == initial(1-5%)"}
	argsJudgeValue := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeValue", KeyCn: "预期判断值", ValueType: "string", DescriptionCn: "预期判断值", Description: "expected judgment value"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsJudgeType, &argsJudgeValue})
}

func InitMonitorMeasureArgs(ctx context.Context, measureInject basic.MeasureInject) error {
	if err := initMeasureCommon(ctx, measureInject.Id, "monitor"); err != nil {
		log.Error(err)
		return err
	}
	if err := initMonitorMeasureJudge(ctx, measureInject); err != nil {
		log.Error(err)
		return err
	}
	argsQuery := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "query", KeyCn: "监控查询语句", ValueType: "string", DescriptionCn: "比如:node_memory_MemAvailable_bytes{instance=\"192.168.2.189:9100\"}", Description: "For example: node_memory_MemAvailable_bytes{instance=\"192.168.2.189:9100\"}"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsQuery})
}

func InitTcpMeasure(ctx context.Context) error {
	var (
		tcpMeasure = basic.MeasureInject{MeasureType: "tcp", Name: "tcp", NameCn: "tcp", Description: "Make expected judgments on TCP requests, such as testing whether the 8080 port of a certain server can be passed", DescriptionCn: "对tcp请求进行预期判断，比如测试某个服务器的8080端口是否能通"}
	)
	if err := basic.InsertMeasureInject(&tcpMeasure); err != nil {
		return err
	}
	return InitTcpMeasureArgs(ctx, tcpMeasure)
}

func initTcpMeasureJudge(ctx context.Context, measureInject basic.MeasureInject) error {
	argsJudgeType := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeType", KeyCn: "预期判断方式", ValueType: "bool", ValueRule: "connectivity", DescriptionCn: "判断目标端口tcp服务是否能连通", Description: "Determine whether the target port tcp service can be connected"}
	argsJudgeValue := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeValue", KeyCn: "预期判断值", ValueType: "string", ValueRule: "true,false", DescriptionCn: "预期判断值", Description: "expected judgment value"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsJudgeType, &argsJudgeValue})
}

func InitTcpMeasureArgs(ctx context.Context, measureInject basic.MeasureInject) error {
	if err := initMeasureCommon(ctx, measureInject.Id, "tcp"); err != nil {
		log.Error(err)
		return err
	}
	if err := initTcpMeasureJudge(ctx, measureInject); err != nil {
		log.Error(err)
		return err
	}
	argsIp := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "ip", KeyCn: "目标机器ip", ValueType: "string", DescriptionCn: "比如:192.168.2.189", Description: "For example: 192.168.2.189"}
	argsPort := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "port", KeyCn: "目标机器端口", ValueType: "string", DescriptionCn: "单个端口号", Description: "Single port number"}
	argsTimeout := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "timeout", KeyCn: "尝试连接超时时间(s)", ValueType: "int", DescriptionCn: "", Description: ""}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsIp, &argsPort, &argsTimeout})
}

func InitPodMeasure(ctx context.Context) error {
	var (
		podMeasure = basic.MeasureInject{MeasureType: "pod", Name: "pod", NameCn: "pod", Description: "Make expected judgments on pod-related data, such as whether the number of pod instances of an application is greater than 3", DescriptionCn: "对pod相关数据进行预期判断,比如某个应用的pod实例数是否大于3"}
	)
	if err := basic.InsertMeasureInject(&podMeasure); err != nil {
		return err
	}
	return InitPodMeasureArgs(ctx, podMeasure)
}

func initPodMeasureJudge(ctx context.Context, measureInject basic.MeasureInject) error {
	argsJudgeType := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeType", KeyCn: "预期判断方式", ValueType: "string", ValueRule: "count", DescriptionCn: "判断pod的数量是否符合预期\njudgeValue样例：\n\"2,5\"：2<=pod数量<=5\n\"5\"：pod数量等于5", Description: "Determine whether the number of pods is as expected\njudgeValue example:\n\"2,5\": 2<=number of pods<=5\n\"5\": The number of pods is equal to 5"}
	argsJudgeValue := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeValue", KeyCn: "预期判断值", ValueType: "string", DescriptionCn: "预期判断值", Description: "expected judgment value"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsJudgeType, &argsJudgeValue})
}

func InitPodMeasureArgs(ctx context.Context, measureInject basic.MeasureInject) error {
	if err := initMeasureCommon(ctx, measureInject.Id, "pod"); err != nil {
		log.Error(err)
		return err
	}
	if err := initPodMeasureJudge(ctx, measureInject); err != nil {
		log.Error(err)
		return err
	}
	argsNamespace := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "namespace", KeyCn: "namespace", ValueType: "string", DescriptionCn: "集群中已存在的namespace", Description: "Namespace that already exists in the cluster"}
	argsLabel := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "label", KeyCn: "label", ValueType: "string", DescriptionCn: "k8s标签", Description: "K8S label"}
	argsNameprefix := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "nameprefix", KeyCn: "pod名称前缀", ValueType: "string", DescriptionCn: "通过名称前缀筛选pod", Description: "Filter pods by name prefix"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsNamespace, &argsLabel, &argsNameprefix})
}

func InitHTTPMeasure(ctx context.Context) error {
	var (
		httpMeasure = basic.MeasureInject{MeasureType: "http", Name: "http", NameCn: "http", Description: "Make expected judgments on http requests, such as whether the return status code is 200 when making a specified http request", DescriptionCn: "对http请求进行预期判断,比如进行指定的http请求时,返回状态码是否为200"}
	)
	if err := basic.InsertMeasureInject(&httpMeasure); err != nil {
		return err
	}
	return InitHTTPMeasureArgs(ctx, httpMeasure)
}

func initHTTPMeasureJudge(ctx context.Context, measureInject basic.MeasureInject) error {
	argsJudgeType := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeType", KeyCn: "预期判断方式", ValueType: "string", ValueRule: "connectivity,code,body", DescriptionCn: "1.connectivity(判断http请求是否能连通)\njudgeValue样例:\ntrue、false\n2.code(判断http返回状态码是否符合预期)\njudgeValue样例:\n200、404、502等\n3.body(判断http返回body中目标key的值是否一致，只支持json格式数据)\njudgeValue样例:\n{\"result\": {\"code\": 0}}\n表示返回体的result的值是json，这个json里面的code是0", Description: "1.Connectivity (determine whether the http request can be connected)\njudgeValue example:\ntrue, false\n2.code (determine whether the http return status code meets expectations)\njudgeValue example:\n200, 404, 502, etc.\n3.body (determine whether the value of the target key in the body returned by http is consistent, only json format data is supported)\njudgeValue example:\n{\"result\": {\"code\": 0}}\nThe result value representing the return body is json, and the code in this json is 0"}
	argsJudgeValue := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasureCommon, Key: "judgeValue", KeyCn: "预期判断值", ValueType: "string", DescriptionCn: "预期判断值", Description: "expected judgment value"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsJudgeType, &argsJudgeValue})
}

func InitHTTPMeasureArgs(ctx context.Context, measureInject basic.MeasureInject) error {
	if err := initMeasureCommon(ctx, measureInject.Id, "http"); err != nil {
		log.Error(err)
		return err
	}
	if err := initHTTPMeasureJudge(ctx, measureInject); err != nil {
		log.Error(err)
		return err
	}
	argsIp := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "ip", KeyCn: "目标机器ip", ValueType: "string", DescriptionCn: "比如:192.168.2.189", Description: "For example: 192.168.2.189"}
	argsPort := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "port", KeyCn: "目标机器端口", ValueType: "string", DescriptionCn: "单个端口号", Description: "Single port number"}
	argsPath := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "path", KeyCn: "请求path", ValueType: "string", DescriptionCn: "url路径", Description: "URL path"}
	argsHeader := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "header", KeyCn: "请求header", ValueType: "string", DescriptionCn: "键值对列表,格式:'k1:v1,k2:v2'", Description: "List of key-value pairs, format: 'k1:v1,k2:v2'"}
	argsScheme := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "scheme", KeyCn: "请求协议", ValueType: "string", ValueRule: "HTTP,HTTPS", DescriptionCn: "", Description: ""}
	argsMethod := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "method", KeyCn: "请求方法", ValueType: "string", ValueRule: "GET,POST", DescriptionCn: "", Description: ""}
	argsTimeout := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "timeout", KeyCn: "尝试连接超时时间(s)", ValueType: "int", DescriptionCn: "", Description: ""}
	argsBody := basic.Args{InjectId: measureInject.Id, ExecType: ExecMeasure, Key: "body", KeyCn: "请求数据", ValueType: "string", DescriptionCn: "任意数据,比如:{\"name\":\"a\",\"age\":201}", Description: "Any data, such as: {\"name\":\"a\",\"age\":201}"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsIp, &argsPort, &argsPath, &argsHeader, &argsScheme, &argsMethod, &argsTimeout, &argsBody})
}

func (i *InjectService) ListMeasures(ctx context.Context, orderBy string, page, pageSize int) (int64, []basic.MeasureInject, error) {
	total, measures, err := basic.ListMeasureInjects(orderBy, page, pageSize)
	return total, measures, err
}
