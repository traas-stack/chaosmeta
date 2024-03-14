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

package experiment

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/util/log"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type (
	ExperimentStatus int
	ScheduleType     string
	TimeType         string
)

const (
	TablePrefix = "experiment_"

	ToBeExecuted = ExperimentStatus(0) //待执行
	Executed     = ExperimentStatus(1)

	ManualMode = ScheduleType("manual") //手动模式
	OnceMode   = ScheduleType("once")   //自动模式
	CronMode   = ScheduleType("cron")

	RecentDayType = TimeType("recent")
	RangeTimeType = TimeType("range")

	TimeLayout = "2006-01-02 15:04:05"
)

type Experiment struct {
	UUID         string           `json:"uuid,omitempty" orm:"column(uuid);size(128);pk"`
	Name         string           `json:"name" orm:"index;column(name);size(255)"`
	Description  string           `json:"description" orm:"column(description);size(1024)"`
	Creator      int              `json:"creator" orm:"index;column(creator)"`
	NamespaceID  int              `json:"namespace_id" orm:"index;column(namespace_id)"`
	ClusterID    int              `json:"cluster_id" orm:"index;column(cluster_id)"`
	ScheduleType string           `json:"schedule_type" orm:"column(schedule_type);size(32);default(manual)"`
	ScheduleRule string           `json:"schedule_rule" orm:"column(schedule_rule);size(64)"`
	NextExec     time.Time        `json:"next_exec,omitempty" orm:"null;column(next_exec);type(datetime)"`
	Status       ExperimentStatus `json:"-" orm:"index;column(status);type:tinyint(1)"`
	LastInstance string           `json:"last_instance" orm:"column(last_instance);size(64)"`
	Version      int              `json:"-" orm:"column(version);default(0);index"`
	models.BaseTimeModel
}

func (e *Experiment) TableName() string {
	return "experiment"
}

func CreateExperiment(experiment *Experiment) error {
	_, err := models.GetORM().Insert(experiment)
	return err
}

func UpdateExperiment(experiment *Experiment) error {
	o := models.GetORM()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	existing := Experiment{UUID: experiment.UUID}
	err = tx.Read(&existing)
	if err != nil {
		tx.Rollback()
		return err
	}

	if experiment.Version != existing.Version {
		tx.Rollback()
		return errors.New("Concurrent modification detected")
	}

	experiment.Version = existing.Version + 1
	if _, err = tx.Update(experiment); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

func GetExperimentByUUID(uuid string) (*Experiment, error) {
	var exp Experiment
	err := models.GetORM().QueryTable(new(Experiment).TableName()).Filter("uuid", uuid).One(&exp)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &exp, nil
}

func DeleteExperimentByUUID(uuid string) error {
	experiment := &Experiment{UUID: uuid}
	_, err := models.GetORM().Delete(experiment)
	return err
}

func ListExperimentsByScheduleTypeAndStatus(scheduleType ScheduleType, experimentStatus ExperimentStatus) (int64, []*Experiment, error) {
	o := models.GetORM()
	experiments := []*Experiment{}
	qs := o.QueryTable(new(Experiment).TableName())

	experimentQuery, err := models.NewDataSelectQuery(&qs)
	if err != nil {
		return 0, nil, err
	}
	if scheduleType != "" {
		experimentQuery.Filter("schedule_type", models.NEGLECT, false, string(scheduleType))
	}

	if experimentStatus >= 0 {
		experimentQuery.Filter("status", models.NEGLECT, false, experimentStatus)
	}
	var totalCount int64
	totalCount, err = experimentQuery.GetOamQuerySeter().Count()
	if err != nil {
		return 0, nil, err
	}

	experimentQuery.OrderBy("create_time")
	_, err = experimentQuery.GetOamQuerySeter().All(&experiments)

	if err == orm.ErrNoRows {
		return 0, nil, nil
	}

	return totalCount, experiments, err
}

func SearchExperiments(lastInstance string, namespaceId int, creator int, name string, scheduleType string, timeType string, timeSearchField string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []*Experiment, error) {
	o := models.GetORM()
	experiments := []*Experiment{}
	qs := o.QueryTable(new(Experiment).TableName())

	experimentQuery, err := models.NewDataSelectQuery(&qs)
	if err != nil {
		return 0, nil, err
	}

	if creator > 0 {
		experimentQuery.Filter("creator", models.NEGLECT, false, creator)
	}
	if lastInstance != "" {
		experimentQuery.Filter("last_instance", models.NEGLECT, false, lastInstance)
	}
	if namespaceId > 0 {
		experimentQuery.Filter("namespace_id", models.NEGLECT, false, namespaceId)
	}
	if scheduleType != "" {
		experimentQuery.Filter("schedule_type", models.NEGLECT, false, scheduleType)
	}
	if name != "" {
		experimentQuery.Filter("name", models.CONTAINS, true, name)
	}
	if timeSearchField == "" {
		timeSearchField = "create_time"
	}
	if timeType == string(RecentDayType) {
		if recentDays > 0 {
			start := time.Now().Add(time.Duration(-recentDays*24) * time.Hour).Format(TimeLayout)
			experimentQuery.Filter(timeSearchField, models.GTE, false, start)
		}
	}

	if timeType == string(RangeTimeType) {
		if !startTime.IsZero() && !endTime.IsZero() {
			experimentQuery.Filter(timeSearchField, models.GTE, false, startTime.Format(TimeLayout))
			experimentQuery.Filter(timeSearchField, models.LTE, false, endTime.Format(TimeLayout))
		}
	}

	var totalCount int64
	totalCount, err = experimentQuery.GetOamQuerySeter().Count()
	if err != nil {
		return 0, nil, err
	}

	orderByStr := "-create_time"
	if orderBy != "" {
		orderByStr = orderBy
	}
	experimentQuery.OrderBy(orderByStr)
	if err := experimentQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = experimentQuery.GetOamQuerySeter().All(&experiments)
	return totalCount, experiments, err
}

func acquireLock(o orm.Ormer, uuid string) bool {
	sql := fmt.Sprintf("UPDATE %s SET version=version+1 WHERE uuid=? AND version=?", TablePrefix)
	res, err := o.Raw(sql, uuid, 0).Exec()
	if err != nil {
		fmt.Println("获取锁失败:", err)
		return false
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Error("获取锁失败:", err)
		return false
	}
	return affected > 0
}

func CountExperiments(namespaceID int, status int, recentDays int) (int64, error) {
	o := models.GetORM()
	qs := o.QueryTable(new(Experiment).TableName())

	if namespaceID != 0 {
		qs = qs.Filter("namespace_id", namespaceID)
	}

	if status >= 0 {
		qs = qs.Filter("status", status)
	}

	if recentDays > 0 {
		start := time.Now().Add(time.Duration(-recentDays*24) * time.Hour).Format(TimeLayout)
		qs = qs.Filter("create_time__gte", start)
	}

	total, err := qs.Count()
	if err == orm.ErrNoRows {
		return 0, nil
	}
	return total, err
}
