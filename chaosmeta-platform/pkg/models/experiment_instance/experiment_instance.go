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
	"github.com/beego/beego/v2/client/orm"
	"time"
)

const (
	TablePrefix = "experiment_"
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
	_, err := models.GetORM().Update(experiment)
	return err
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

func SearchExperimentInstances(lastInstance string, namespaceId int, creator int, name string, scheduleType string, timeType string, recentDays int, startTime, endTime time.Time, orderBy string, page, pageSize int) (int64, []*ExperimentInstance, error) {
	o := models.GetORM()
	experiments := []*ExperimentInstance{}
	qs := o.QueryTable(new(ExperimentInstance).TableName())

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
	if timeType != "" {
		if recentDays > 0 {
			start := time.Now().Add(time.Duration(-recentDays*24) * time.Hour).Format("2006-01-02 15:04:05")
			experimentQuery.Filter("create_time", models.GTE, false, start)
		}

		if !startTime.IsZero() && !endTime.IsZero() {
			experimentQuery.Filter("create_time", models.GTE, false, startTime)
			experimentQuery.Filter("create_time", models.LTE, false, endTime)
		}
	}

	var totalCount int64
	totalCount, err = experimentQuery.GetOamQuerySeter().Count()

	orderByList := []string{"uuid"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	experimentQuery.OrderBy(orderByList...)
	if err := experimentQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = experimentQuery.GetOamQuerySeter().All(experiments)
	return totalCount, experiments, err
}
