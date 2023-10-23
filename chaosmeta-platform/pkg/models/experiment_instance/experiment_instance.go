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

package experiment_instance

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/util/log"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
	"time"
)

type (
	ExperimentInstanceStatus string
	TimeType                 string
)

const (
	TablePrefix = "experiment_"

	Pending = ExperimentInstanceStatus("Pending") //待执行
	Running = ExperimentInstanceStatus("Running") //执行中

	TimeLayout = "2006-01-02 15:04:05"

	RecentDayType = TimeType("recent")
	RangeTimeType = TimeType("range")
)

type ExperimentInstance struct {
	UUID           string `json:"uuid,omitempty" orm:"column(uuid);size(128);pk"`
	Name           string `json:"name" orm:"index;column(name);size(255)"`
	NamespaceID    int    `json:"namespace_id" orm:"index;column(namespace_id)"`
	Description    string `json:"description" orm:"column(description);size(1024)"`
	ExperimentUUID string `json:"experiment_uuid,omitempty" orm:"column(experiment_uuid);size(128);index"`
	Creator        int    `json:"creator" orm:"index;column(creator)"`
	Status         string `json:"status" orm:"column(status);size(32);index"`
	Message        string `json:"message" orm:"column(message);size(1024)"`
	Version        int    `json:"-" orm:"column(version);default(0);version"`
	models.BaseTimeModel
}

func (e *ExperimentInstance) TableName() string {
	return "experiment_instance"
}

func CreateExperimentInstance(experiment *ExperimentInstance) error {
	_, err := models.GetORM().Insert(experiment)
	return err
}

func UpdateExperimentInstance(experiment *ExperimentInstance) error {
	o := models.GetORM()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	existing := ExperimentInstance{UUID: experiment.UUID}
	err = tx.Read(&existing)
	if err != nil {
		tx.Rollback()
		return err
	}

	if experiment.Version != existing.Version {
		tx.Rollback()
		log.Error("Concurrent modification detected")
		return errors.New("Concurrent modification detected")
	}

	experiment.Version = existing.Version + 1
	if _, err = tx.Update(experiment); err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Error(err)
		return err
	}

	return nil
}

func UpdateExperimentInstanceStatus(uuid string, status, message string) error {
	experimentInstance, err := GetExperimentInstanceByUUID(uuid)
	if err != nil || experimentInstance == nil {
		return fmt.Errorf("error:%v", err)
	}
	experimentInstance.Status = status
	if message != "" {
		experimentInstance.Message = message
	}
	return UpdateExperimentInstance(experimentInstance)
}

func GetExperimentInstanceByUUID(uuid string) (*ExperimentInstance, error) {
	var exp ExperimentInstance
	err := models.GetORM().QueryTable(new(ExperimentInstance).TableName()).Filter("uuid", uuid).One(&exp)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &exp, nil
}

func DeleteExperimentInstanceByUUID(uuid string) error {
	experiment := &ExperimentInstance{UUID: uuid}
	_, err := models.GetORM().Delete(experiment)
	return err
}

func SearchExperimentInstances(lastInstance string, experimentUUID string, namespaceId int, creator int, name string, timeType string, timeSearchField string, status string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []*ExperimentInstance, error) {
	o := models.GetORM()
	experiments := []*ExperimentInstance{}
	qs := o.QueryTable(new(ExperimentInstance).TableName())

	experimentQuery, err := models.NewDataSelectQuery(&qs)
	if err != nil {
		return 0, nil, err
	}

	if creator != 0 {
		experimentQuery.Filter("creator", models.NEGLECT, false, creator)
	}
	if experimentUUID != "" {
		experimentQuery.Filter("experiment_uuid", models.NEGLECT, false, experimentUUID)
	}
	if lastInstance != "" {
		experimentQuery.Filter("last_instance", models.NEGLECT, false, lastInstance)
	}
	if namespaceId > 0 {
		experimentQuery.Filter("namespace_id", models.NEGLECT, false, namespaceId)
	}
	if status != "" {
		if status == "null" {
			status = ""
		}
		experimentQuery.Filter("status", models.NEGLECT, false, status)
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
	if err == orm.ErrNoRows {
		return 0, nil, nil
	}
	return totalCount, experiments, err
}

func ListExperimentsInstancesByStatus(experimentStatus []ExperimentInstanceStatus) (int64, []*ExperimentInstance, error) {
	o := models.GetORM()
	experiments := []*ExperimentInstance{}
	qs := o.QueryTable(new(ExperimentInstance).TableName())

	experimentQuery, err := models.NewDataSelectQuery(&qs)
	if err != nil {
		return 0, nil, err
	}
	if experimentStatus != nil {
		experimentQuery.Filter("status", models.IN, false, experimentStatus)
	}
	var totalCount int64
	totalCount, err = experimentQuery.GetOamQuerySeter().Count()
	if err != nil {
		return 0, nil, err
	}

	experimentQuery.OrderBy("create_time")
	_, err = experimentQuery.GetOamQuerySeter().All(experiments)

	return totalCount, experiments, err
}

func CountExperimentInstance(namespaceId, day int) (map[string]int64, int64, error) {
	o := models.GetORM()

	var (
		counts []orm.Params
		result = make(map[string]int64)
		total  int64
		sql    string
		args   []interface{}
	)

	if day == 0 {
		sql = "SELECT status, COUNT(*) as count FROM experiment_instance WHERE namespace_id = ? GROUP BY status"
		args = []interface{}{namespaceId}
	} else {
		startTime := time.Now().AddDate(0, 0, -day).Format(TimeLayout)
		sql = "SELECT status, COUNT(*) as count FROM experiment_instance WHERE namespace_id = ? AND create_time > ? GROUP BY status"
		args = []interface{}{namespaceId, startTime}
	}

	_, err := o.Raw(sql, args...).Values(&counts)
	if err != nil {
		return nil, 0, err
	}

	for _, c := range counts {
		status := c["status"].(string)
		countStr := c["count"].(string)
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		result[status] = count
		total += count
	}

	return result, total, nil
}

type ExperimentInstanceStatusCount struct {
	PendingCount   int64
	RunningCount   int64
	SucceededCount int64
	FailedCount    int64
	ErrorCount     int64
	OtherCount     int64
}

func CountExperimentInstances(namespaceID int, experimentUUID string, status string, recentDays int) (int64, error) {
	o := models.GetORM()
	qs := o.QueryTable(new(ExperimentInstance).TableName())

	if namespaceID != 0 {
		qs = qs.Filter("namespace_id", namespaceID)
	}

	if experimentUUID != "" {
		qs = qs.Filter("experiment_uuid", experimentUUID)
	}
	if status != "" {
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
