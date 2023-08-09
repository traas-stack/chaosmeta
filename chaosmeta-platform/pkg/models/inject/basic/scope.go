package basic

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"github.com/beego/beego/v2/client/orm"
)

type Scope struct {
	ID            int    `json:"id" orm:"pk;auto;column(id)"`
	Name          string `json:"name" orm:"unique;index;size(255);column(name)"`
	NameCn        string `json:"nameCn" orm:"size(255);column(name_cn)"`
	Description   string `json:"description" orm:"size(1024);column(description)"`
	DescriptionCn string `json:"descriptionCn" orm:"size(1024);column(description_cn)"`
	models.BaseTimeModel
}

func (s *Scope) TableName() string {
	return TablePrefix + "scope"
}

func InsertScope(ctx context.Context, s *Scope) (int64, error) {
	return models.GetORM().Insert(s)
}

func DeleteScope(ctx context.Context, id int) error {
	scope := &Scope{ID: id}
	_, err := models.GetORM().Delete(scope)
	return err
}

func UpdateScope(ctx context.Context, scope *Scope) error {
	if models.GetORM().Read(scope) == nil {
		_, err := models.GetORM().Update(scope)
		return err
	}
	return errors.New("scope Not Found")
}

func GetScopeById(ctx context.Context, id int) (*Scope, error) {
	o := models.GetORM()
	scope := &Scope{ID: id}
	err := o.Read(scope)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err == orm.ErrMissPK {
		return nil, nil
	} else {
		return scope, err
	}
}

func ListScopes(ctx context.Context, orderBy string, page, pageSize int) (int64, []Scope, error) {
	scope, scopes := Scope{}, new([]Scope)

	querySeter := models.GetORM().QueryTable(scope.TableName())
	scopeQuery, err := models.NewDataSelectQuery(&querySeter)
	if err != nil {
		return 0, nil, err
	}

	var totalCount int64
	totalCount, err = scopeQuery.GetOamQuerySeter().Count()

	orderByList := []string{"id"}
	if len(orderBy) > 0 {
		orderByList = append(orderByList, orderBy)
	}
	scopeQuery.OrderBy(orderByList...)
	if err := scopeQuery.Limit(pageSize, (page-1)*pageSize); err != nil {
		return 0, nil, err
	}

	_, err = scopeQuery.GetOamQuerySeter().All(scopes)
	return totalCount, *scopes, err
}
