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

package models

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

var globalORM orm.Ormer

func Setup() {
	orm.RegisterModel(new(models.UserNamespace), new(models.User), new(models.Namespace))

	if err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.DefaultRunOptIns.DB.User, config.DefaultRunOptIns.DB.Passwd, config.DefaultRunOptIns.DB.Url, config.DefaultRunOptIns.DB.Name), orm.MaxIdleConnections(config.DefaultRunOptIns.DB.MaxIdle), orm.MaxOpenConnections(config.DefaultRunOptIns.DB.MaxConn)); err != nil {
		panic(any(fmt.Sprintf("connect database error: %s", err.Error())))
	}

	orm.Debug = true // TODO: only open in dev, should not open in prod
	if err := orm.RunSyncdb("default", false, true); err != nil {
		panic(any(fmt.Sprintf("sync database error: %s", err.Error())))
	}

	globalORM = orm.NewOrm()
}

func GetORM() orm.Ormer {
	return globalORM
}
