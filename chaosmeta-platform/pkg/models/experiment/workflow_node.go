package experiment

import (
	models "chaosmeta-platform/pkg/models/common"
	"github.com/beego/beego/v2/client/orm"
)

type WorkflowNode struct {
	UUID           string `json:"uuid,omitempty" orm:"column(uuid);pk"`
	ExperimentUUID string `json:"experiment_uuid" orm:"index;column(experiment_uuid);size(64)"`
	Row            int    `json:"row" orm:"column(row)"`
	Column         int    `json:"column" orm:"column(column)"`
	Duration       string `json:"duration" orm:"column(duration);size(32)"`
	ExecType       string `json:"exec_type" orm:"column(exec_type);size(32)"`
	ExecID         int    `json:"exec_id" orm:"column(exec_id)"`
	models.BaseTimeModel
}

func (wn *WorkflowNode) TableName() string {
	return TablePrefix + "workflow_node"
}

func GetWorkflowNodeByUUID(uuid string) (*WorkflowNode, error) {
	workflowNode := &WorkflowNode{UUID: uuid}
	err := models.GetORM().Read(workflowNode)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNode, nil
}

func GetWorkflowNodesByExperimentUUID(experimentUUID string) ([]*WorkflowNode, error) {
	workflowNodes := []*WorkflowNode{}
	_, err := models.GetORM().QueryTable(new(WorkflowNode).TableName()).Filter("experiment_uuid", experimentUUID).OrderBy("row", "column").All(&workflowNodes)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNodes, nil
}

func CreateWorkflowNode(workflowNode *WorkflowNode) error {
	_, err := models.GetORM().Insert(workflowNode)
	return err
}

func DeleteWorkflowNodeByUUID(uuid string) error {
	workflowNode := &WorkflowNode{UUID: uuid}
	_, err := models.GetORM().Delete(workflowNode)
	return err
}

// BatchSearchWorkflowNodes 批量搜索workflow_nodes
func BatchSearchWorkflowNodes(searchCriteria map[string]interface{}) ([]*WorkflowNode, error) {
	o := models.GetORM()
	workflowNodes := []*WorkflowNode{}
	qs := o.QueryTable(new(WorkflowNode).TableName())
	for key, value := range searchCriteria {
		qs = qs.Filter(key, value)
	}

	_, err := qs.All(&workflowNodes)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return workflowNodes, nil
}
