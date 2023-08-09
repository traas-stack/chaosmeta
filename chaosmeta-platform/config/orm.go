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

package config

import (
	"chaosmeta-platform/pkg/models/agent"
	"chaosmeta-platform/pkg/models/cluster"
	modelCommon "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/pkg/models/experiment"
	"chaosmeta-platform/pkg/models/inject/basic"
	"chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/models/user"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Setup() {
	orm.RegisterModel(
		new(namespace.ClusterNamespace), new(namespace.Label), new(namespace.Namespace), new(namespace.UserNamespace),
		new(user.User),
		new(cluster.Cluster),
		new(agent.Agent),
		new(basic.Scope), new(basic.Target), new(basic.Fault), new(basic.Args),
		new(experiment.WorkflowNode), new(experiment.LabelExperiment), new(experiment.FaultRange), new(experiment.Experiment), new(experiment.ArgsValue),
	)

	if err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local", DefaultRunOptIns.DB.User, DefaultRunOptIns.DB.Passwd, DefaultRunOptIns.DB.Url, DefaultRunOptIns.DB.Name), orm.MaxIdleConnections(DefaultRunOptIns.DB.MaxIdle), orm.MaxOpenConnections(DefaultRunOptIns.DB.MaxConn)); err != nil {
		panic(any(fmt.Sprintf("connect database error: %s", err.Error())))
	}

	orm.Debug = true // TODO: only open in dev, should not open in prod
	if err := orm.RunSyncdb("default", false, true); err != nil {
		panic(any(fmt.Sprintf("sync database error: %s", err.Error())))
	}
	modelCommon.GlobalORM = orm.NewOrm()
}
