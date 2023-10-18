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
	"chaosmeta-platform/pkg/models/experiment_instance"
	"chaosmeta-platform/pkg/models/inject/basic"
	"chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/models/user"
	"chaosmeta-platform/util/log"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Setup() {
	orm.RegisterModel(
		new(namespace.ClusterNamespace), new(namespace.Label), new(namespace.Namespace), new(namespace.UserNamespace), new(user.User),
		new(cluster.Cluster),
		new(agent.Agent),
		new(basic.Scope), new(basic.Target), new(basic.Fault), new(basic.FlowInject), new(basic.MeasureInject), new(basic.Args),
		new(experiment.WorkflowNode), new(experiment.LabelExperiment), new(experiment.FaultRange), new(experiment.FlowRange), new(experiment.MeasureRange), new(experiment.Experiment), new(experiment.ArgsValue),
		new(experiment_instance.WorkflowNodeInstance), new(experiment_instance.LabelExperimentInstance), new(experiment_instance.FaultRangeInstance), new(experiment_instance.FlowRangeInstance), new(experiment_instance.MeasureRangeInstance), new(experiment_instance.ExperimentInstance), new(experiment_instance.ArgsValueInstance),
	)

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		for range ticker.C {
			err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local", DefaultRunOptIns.DB.User, DefaultRunOptIns.DB.Passwd, DefaultRunOptIns.DB.Url, DefaultRunOptIns.DB.Name), orm.MaxIdleConnections(DefaultRunOptIns.DB.MaxIdle), orm.MaxOpenConnections(DefaultRunOptIns.DB.MaxConn))
			if err != nil {
				log.Error(err)
			}
			if err == nil {
				done <- true
				break
			}
		}
	}()

	select {
	case <-done:
		log.Info("successfully connected to the database")
	case <-time.After(5 * time.Minute):
		panic(any("connect database failed"))
	}

	ticker.Stop()

	orm.Debug = DefaultRunOptIns.DB.Debug
	modelCommon.GlobalORM = orm.NewOrm()

	if err := DropTablesBeforeCreate(); err != nil {
		log.Error("drop tables before create error: %s", err.Error())
	}

	if err := orm.RunSyncdb("default", false, true); err != nil {
		panic(any(fmt.Sprintf("sync database error: %s", err.Error())))
	}
}

func DropTablesBeforeCreate() error {
	var dropTables = []string{"args", "fault", "flow_inject", "measure_inject", "scope", "target"}
	for _, tableName := range dropTables {
		if _, err := modelCommon.GlobalORM.Raw(fmt.Sprintf("DROP TABLE inject_basic_%s", tableName)).Exec(); err != nil {
			log.Error("drop table failed, err:", err)
		}
	}
	return nil
}
