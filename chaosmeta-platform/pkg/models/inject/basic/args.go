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

package basic

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type Args struct {
	ID            int    `json:"id" orm:"pk;auto;column(id)"`
	ExecType      string `json:"execType" orm:"index;column(exec_type)"`
	InjectId      int    `json:"injectId" orm:"index;column(inject_id)"`
	Key           string `json:"key" orm:"size(255);column(key)"`
	KeyCn         string `json:"keyCn"  orm:"size(255);column(key_cn)"`
	ValueType     string `json:"valueType" orm:"size(32);column(value_type)"`
	ValueRule     string `json:"valueRule" orm:"size(255);column(value_rule)"`
	Description   string `json:"description" orm:"size(1024);column(description)"`
	DescriptionCn string `json:"descriptionCn" orm:"size(1024);column(description_cn)"`
	Unit          string `json:"unit" orm:"size(1024);column(unit)"`
	UnitCn        string `json:"unitCn" orm:"size(1024);column(unit_cn)"`
	DefaultValue  string `json:"defaultValue" orm:"size(1024);column(default_value)"`
	Required      bool   `json:"required" orm:"column(required)"`
	models.BaseTimeModel
}

func (a *Args) TableName() string {
	return TablePrefix + "args"
}

func InsertArgsMulti(ctx context.Context, argsList []*Args) error {
	_, err := models.GetORM().InsertMulti(len(argsList), argsList)
	return err
}

func InsertArgs(ctx context.Context, args *Args) error {
	_, err := models.GetORM().Insert(args)
	return err
}

func DeleteArgs(ctx context.Context, id int) error {
	args := &Args{ID: id}
	_, err := models.GetORM().Delete(args)
	return err
}

func DeleteArgsMulti(ctx context.Context, injectId int, execType string) error {
	arg := Args{}
	_, err := models.GetORM().QueryTable(arg.TableName()).Filter("inject_id", injectId).Filter("exec_type", execType).Delete()
	return err
}

func UpdateArgs(ctx context.Context, args *Args) error {
	if models.GetORM().Read(args) == nil {
		_, err := models.GetORM().Update(args)
		return err
	}

	return errors.New("args not found")
}

func GetArgsById(ctx context.Context, id int) (*Args, error) {
	args := &Args{ID: id}
	err := models.GetORM().Read(args)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, nil
	} else {
		return args, err
	}
}

func ListArgs(ctx context.Context, execType []string, injectId int, orderBy string, page, pageSize int) (int64, []Args, error) {
	arg, args := Args{}, new([]Args)

	querySeter := models.GetORM().QueryTable(arg.TableName())
	argsQuery, err := models.NewDataSelectQuery(&querySeter)

	if injectId > 0 {
		argsQuery.Filter("inject_id", models.NEGLECT, false, injectId)
	}

	if len(execType) != 0 {
		argsQuery.Filter("exec_type", models.IN, false, execType)
	}

	var totalCount int64
	totalCount, err = argsQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	argsQuery.OrderBy(orderByList...)
	if err := argsQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = argsQuery.GetOamQuerySeter().All(args)
	return totalCount, *args, err
}
