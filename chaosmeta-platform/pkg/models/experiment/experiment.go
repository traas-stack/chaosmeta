package experiment

import (
	models "chaosmeta-platform/pkg/models/common"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

const (
	TablePrefix = "experiment_"
)

type Experiment struct {
	UUID         string    `json:"uuid,omitempty" orm:"column(uuid);size(128);pk"`
	Name         string    `json:"name" orm:"index;column(name);size(255)"`
	Description  string    `json:"description" orm:"column(description);size(1024)"`
	Creator      int       `json:"creator" orm:"index;column(creator)"`
	NamespaceID  int       `json:"namespace_id" orm:"index;column(namespace_id)"`
	ScheduleType string    `json:"schedule_type" orm:"column(schedule_type);size(32);default(manual)"`
	ScheduleRule string    `json:"schedule_rule" orm:"column(schedule_rule);size(64)"`
	NextExec     time.Time `json:"next_exec" orm:"column(next_exec);type(datetime)"`
	LastInstance string    `json:"last_instance" orm:"column(last_instance);size(64)"`
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
	_, err := models.GetORM().Update(experiment)
	return err
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

func SearchExperiments(lastInstance string, namespaceId int, name string, scheduleType string, orderBy string, page, pageSize int) (int64, []*Experiment, error) {
	o := models.GetORM()
	experiments := []*Experiment{}
	qs := o.QueryTable(new(Experiment).TableName())

	experimentQuery, err := models.NewDataSelectQuery(&qs)
	if err != nil {
		return 0, nil, err
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
