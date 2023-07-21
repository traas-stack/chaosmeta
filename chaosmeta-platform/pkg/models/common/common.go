package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type OperatorType int

const (
	NEGLECT    OperatorType = iota // no operator
	GT                             // >
	GTE                            // >=
	LT                             // <
	LTE                            //<=
	IN                             // in
	EXACT                          //
	CONTAINS                       // like'%string%'
	STARTSWITH                     // like 'string%'
	ENDSWITH                       // like '%string'
	ISNULL
)

func (o OperatorType) String() string {
	return [...]string{"no-operator", "gt", "gte", "lt", "lte", "in", "exact", "contains", "startswith", "endswith", "isnull"}[o]
}

func (o OperatorType) IString() string {
	return "i" + o.String()
}

type DataSelectQuery struct {
	OamQuerySeter orm.QuerySeter
}

func NewDataSelectQuery(oamQuerySeter *orm.QuerySeter) (*DataSelectQuery, error) {
	if oamQuerySeter == nil {
		return nil, errors.New("param is wrong")
	}
	return &DataSelectQuery{*oamQuerySeter}, nil
}

func (d *DataSelectQuery) Limit(limit, offset int) error {
	if offset < 0 {
		return errors.New("invalid offset")
	}

	if limit <= 0 {
		return errors.New("invalid limit")
	}

	d.OamQuerySeter = d.OamQuerySeter.Limit(limit, offset)
	return nil
}

func (d *DataSelectQuery) Filter(propertyName string, operatorType OperatorType, ignoreCase bool, propertyValue interface{}) {
	if operatorType == NEGLECT {
		d.OamQuerySeter = d.OamQuerySeter.Filter(propertyName, propertyValue)
		return
	}

	operatorStr := operatorType.String()
	if ignoreCase {
		operatorStr = operatorType.IString()
	}
	d.OamQuerySeter = d.OamQuerySeter.Filter(fmt.Sprintf("%s__%s", propertyName, operatorStr), propertyValue)
}

func (d *DataSelectQuery) Exclude(propertyName string, operatorType OperatorType, ignoreCase bool, propertyValue interface{}) {
	if operatorType == NEGLECT {
		d.OamQuerySeter = d.OamQuerySeter.Exclude(propertyName, propertyValue)
		return
	}

	operatorStr := operatorType.String()
	if ignoreCase {
		operatorStr = operatorType.IString()
	}
	d.OamQuerySeter = d.OamQuerySeter.Exclude(fmt.Sprintf("%s__%s", propertyName, operatorStr), propertyValue)
}

func (d *DataSelectQuery) GroupBy(exprs ...string) {
	d.OamQuerySeter = d.OamQuerySeter.GroupBy(exprs...)
}

func (d *DataSelectQuery) OrderBy(exprs ...string) {
	d.OamQuerySeter = d.OamQuerySeter.OrderBy(exprs...)
}

func (d *DataSelectQuery) GetOamQuerySeter() orm.QuerySeter {
	return d.OamQuerySeter
}

func (d *DataSelectQuery) Delete() (int64, error) {
	return d.OamQuerySeter.Delete()
}

func (d *DataSelectQuery) Update(params orm.Params) (int64, error) {
	return d.OamQuerySeter.Update(params)
}

func (d *DataSelectQuery) Exist() bool {
	return d.OamQuerySeter.Exist()
}
