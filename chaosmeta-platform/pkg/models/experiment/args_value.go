package experiment

import (
	models "chaosmeta-platform/pkg/models/common"
	"github.com/beego/beego/v2/client/orm"
)

type ArgsValue struct {
	ArgsID           int    `json:"args_id" orm:"column(args_id);index"`
	WorkflowNodeUUID string `json:"workflow_node_uuid,omitempty" orm:"column(workflow_node_uuid);index"`
	Value            string `json:"value" orm:"column(value);size(1024)"`
	models.BaseTimeModel
}

func (av *ArgsValue) TableName() string {
	return TablePrefix + "args_value"
}

func BatchInsertArgsValues(workflowNodeUUID string, argsValues []*ArgsValue) error {
	o := models.GetORM()
	oldValues := []*ArgsValue{}
	_, err := o.QueryTable(new(ArgsValue)).Filter("workflow_node_uuid", workflowNodeUUID).All(&oldValues)
	if err != nil {
		return err
	}
	if len(oldValues) > 0 {
		if _, err = o.QueryTable(new(ArgsValue)).Filter("workflow_node_uuid", workflowNodeUUID).Delete(); err != nil {
			return err
		}
	}
	for _, argsValue := range argsValues {
		argsValue.WorkflowNodeUUID = workflowNodeUUID
		if _, err := o.Insert(argsValue); err != nil {
			return err
		}
	}
	return nil
}

func ClearArgsValuesByWorkflowNodeUUID(workflowNodeUUID string) error {
	o := models.GetORM()
	_, err := o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).Delete()
	return err
}

func GetArgsValuesByWorkflowNodeUUID(workflowNodeUUID string) ([]*ArgsValue, error) {
	o := models.GetORM()

	var argsValues []*ArgsValue
	_, err := o.QueryTable(new(ArgsValue).TableName()).Filter("workflow_node_uuid", workflowNodeUUID).OrderBy("-created_at").All(&argsValues)
	if err != nil {
		return nil, err
	}

	return argsValues, nil
}

func BatchSearchArgsValues(searchCriteria map[string]interface{}) ([]*ArgsValue, error) {
	o := models.GetORM()
	argsValues := []*ArgsValue{}
	qs := o.QueryTable(new(ArgsValue).TableName())
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&argsValues)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return argsValues, nil
}
