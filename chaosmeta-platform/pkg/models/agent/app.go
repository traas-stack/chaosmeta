package agent

import (
	models "chaosmeta-platform/pkg/models/common"
	"context"
)

type App struct {
	ID   int    `json:"id" orm:"pk;auto;column(id)"`
	Name string `json:"name" orm:"unique;index;column(name);size(255)"`
	models.BaseTimeModel
}

func (a *App) TableName() string {
	return "agent_app"
}

func InsertApp(ctx context.Context, a *App) (int64, error) {
	return models.GetORM().Insert(a)
}

func GetAppByName(ctx context.Context, a *App) error {
	return models.GetORM().Read(a, "name")
}

func GetAppById(ctx context.Context, a *App) error {
	return models.GetORM().Read(a)
}

//func DeleteUsersByIdList(ctx context.Context, userId []int) error {
//	user := User{}
//	querySeter := models.GetORM().QueryTable(user.TableName())
//	userQuery, err := models.NewDataSelectQuery(&querySeter)
//	if err != nil {
//		return err
//	}
//	userQuery.Filter("id", models.IN, false, userId)
//	_, err = userQuery.Update(orm.Params{
//		"is_deleted": true,
//	})
//	return err
//}
//
//func UpdateUser(ctx context.Context, u *User) error {
//	suc, err := models.GetORM().Update(u)
//	if suc == 0 {
//		return fmt.Errorf("record[email: %s] not found", u.Email)
//	}
//
//	return err
//}
