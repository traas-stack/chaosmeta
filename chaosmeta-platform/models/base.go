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
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	_ "github.com/go-sql-driver/mysql"
)

var globalORM orm.Ormer

// TODO: Whether the field is empty, and the validity check of the value, the default value, etc.
// Do it directly in the platform logic, not in the table design

func Setup(cfg config.Configer) {
	name, err := cfg.String("db::name")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::name] from config error: %s", err.Error())))
	}

	user, err := cfg.String("db::user")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::user] from config error: %s", err.Error())))
	}

	passwd, err := cfg.String("db::passwd")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::passwd] from config error: %s", err.Error())))
	}

	url, err := cfg.String("db::url")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::url] from config error: %s", err.Error())))
	}

	maxidle, err := cfg.Int("db::maxidle")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::maxidle] from config error: %s", err.Error())))
	}

	maxconn, err := cfg.Int("db::maxconn")
	if err != nil {
		panic(any(fmt.Sprintf("get [db::maxconn] from config error: %s", err.Error())))
	}

	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Namespace))

	if err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, passwd, url, name), orm.MaxIdleConnections(maxidle), orm.MaxOpenConnections(maxconn)); err != nil {
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
