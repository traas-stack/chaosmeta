package models

import "time"

var tableNamePrefix = "chaosmeta_platform_"

type BaseTimeModel struct {
	CreateTime time.Time `orm:"column(create_time);auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"column(update_time);auto_now;type(datetime)"`
}

func (bt *BaseTimeModel) SetCreateTime(time time.Time) {
	bt.CreateTime = time
}

func (bt *BaseTimeModel) SetUpdateTime(time time.Time) {
	bt.UpdateTime = time
}
