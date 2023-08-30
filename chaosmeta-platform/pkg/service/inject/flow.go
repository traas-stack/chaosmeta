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
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

func (i *InjectService) ListFlowInjects(ctx context.Context, orderBy string, page, pageSize int) (int64, []basic.FlowInject, error) {
	total, targets, err := basic.ListFlowInjects(orderBy, page, pageSize)
	return total, targets, err
}

func InitHttpFlow(ctx context.Context) error {
	var (
		httpFlow = basic.FlowInject{Name: "HTTP", NameCn: "HTTP", Description: "continuously inject http request traffic to the target http server", DescriptionCn: "对目标 http 服务器持续注入 http 请求流量"}
	)
	if err := basic.InsertFlowInject(&httpFlow); err != nil {
		return err
	}
	return InitHttpFlowArgs(ctx, httpFlow)
}

func InitHttpFlowArgs(ctx context.Context, flowInject basic.FlowInject) error {
	argsHost := basic.Args{InjectId: flowInject.Id, ExecType: ExecFlow, Key: "host", KeyCn: "目标机器", ValueType: "string", DefaultValue: "", Required: false, Description: "目标端口,可选值：ip、域名", DescriptionCn: "Destination port:optional values: ip, domain name"}
	argsPort := basic.Args{InjectId: flowInject.Id, ExecType: ExecFlow, Key: "port", KeyCn: "端口", ValueType: "string", DefaultValue: "目标端口", Required: true, DescriptionCn: "目标端口, 单个端口号", Description: "destination port, a single port number"}
	argsPath := basic.Args{InjectId: flowInject.Id, ExecType: ExecFlow, Key: "path", KeyCn: "请求path", ValueType: "string", DefaultValue: "", Required: true, DescriptionCn: "请求path", Description: "request path"}
	argsHeader := basic.Args{InjectId: flowInject.Id, ExecType: ExecFlow, Key: "header", KeyCn: "请求header", ValueType: "string", DefaultValue: "", Required: false, Description: "List of key-value pairs, format: 'k1:v1,k2:v2'", DescriptionCn: "键值对列表，格式：'k1:v1,k2:v2'"}
	argsMethod := basic.Args{InjectId: flowInject.Id, ExecType: ExecFlow, Key: "method", KeyCn: "方法", ValueType: "string", ValueRule: "GET,POST", DefaultValue: "", Required: true, Description: "request method", DescriptionCn: "请求方法"}
	argsBody := basic.Args{InjectId: flowInject.Id, ExecType: ExecInject, Key: "body", KeyCn: "请求数据", ValueType: "string", DefaultValue: "", Description: "request data", DescriptionCn: "请求数据"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsHost, &argsPort, &argsPath, &argsHeader, &argsMethod, &argsBody})
}
