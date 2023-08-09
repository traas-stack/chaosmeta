package namespace

import (
	"chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
)

type Label struct {
	Id          int    `json:"id" orm:"pk;auto;column(id)"`
	Name        string `json:"name" orm:"column(name);size(255);unique;index"`
	Color       string `json:"color" orm:"column(color);size(255);index"`
	NamespaceId int    `json:"namespaceId" orm:"column(namespace_id);unique;index"`
	Creator     string `json:"creator" orm:"column(creator); size(255); index"`
	models.BaseTimeModel
}

func (l *Label) TableName() string {
	return "namespace_label"
}

func InsertLabel(ctx context.Context, label *Label) (int64, error) {
	if label == nil {
		return 0, errors.New("label is nil")
	}
	id, err := models.GetORM().Insert(label)
	return id, err
}

func UpdateLabel(ctx context.Context, label *Label) (int64, error) {
	if label == nil {
		return 0, errors.New("label is nil")
	}
	num, err := models.GetORM().Update(label)
	return num, err
}

func GetLabelById(ctx context.Context, label *Label) error {
	if label == nil {
		return errors.New("label is nil")
	}
	return models.GetORM().Read(label)
}

func GetLabelByName(ctx context.Context, label *Label) error {
	return models.GetORM().Read(label, "name", "namespace_id")
}

func DeleteLabel(ctx context.Context, id int) (int64, error) {
	num, err := models.GetORM().Delete(&Label{Id: id})
	return num, err
}

func QueryLabels(ctx context.Context, nameSpaceId int, name, creator string, orderBy string, page, pageSize int) (int64, []Label, error) {
	label, labelList := Label{}, new([]Label)
	querySeter := models.GetORM().QueryTable(label.TableName())
	labelQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	labelQuery.Filter("namespace_id", models.NEGLECT, false, nameSpaceId)
	if len(name) > 0 {
		labelQuery.Filter("name", models.CONTAINS, true, name)
	}

	if len(creator) > 0 {
		labelQuery.Filter("creator", models.CONTAINS, true, creator)
	}

	totalCount, err := labelQuery.GetOamQuerySeter().Count()
	if err != nil {
		return 0, nil, err
	}

	orderByList := []string{}
	if orderBy != "" {
		orderByList = append(orderByList, orderBy)
	} else {
		orderByList = append(orderByList, "id")
	}
	labelQuery.OrderBy(orderByList...)
	if err := labelQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = labelQuery.GetOamQuerySeter().All(labelList)
	return totalCount, *labelList, err
}
