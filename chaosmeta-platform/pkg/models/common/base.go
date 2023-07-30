package models

import "github.com/beego/beego/v2/client/orm"

var GlobalORM orm.Ormer

func GetORM() orm.Ormer {
	return GlobalORM
}
