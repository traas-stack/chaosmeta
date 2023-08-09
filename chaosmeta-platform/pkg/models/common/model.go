package models

import "time"

var tableNamePrefix = "chaosmeta_platform_"

type BaseTimeModel struct {
	CreateTime time.Time `json:"create_time" orm:"column(create_time);auto_now_add;type(datetime)"`
	UpdateTime time.Time `json:"update_time" orm:"column(update_time);auto_now;type(datetime)"`
}

func (bt *BaseTimeModel) SetCreateTime(time time.Time) {
	bt.CreateTime = time
}

func (bt *BaseTimeModel) SetUpdateTime(time time.Time) {
	bt.UpdateTime = time
}
