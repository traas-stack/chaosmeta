package models

import (
	"chaosmeta-platform/pkg/models/common"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/spf13/cast"
)

type Permission int

const (
	NormalPermission  = Permission(0)
	AdminPermission   = Permission(1)
	VisitorPermission = Permission(2)
)

type UserNamespace struct {
	ID          int `json:"id" orm:"pk;auto;column(id)"`
	UserId      int `json:"userId" orm:"column(user_id)"`
	NamespaceId int `json:"namespaceId" orm:"column(namespace_id)"`
	//User       *User      `orm:"rel(fk)"`
	//Namespace  *Namespace `orm:"rel(fk)"`
	Permission Permission `json:"permission" orm:"default(0)"`
	models.BaseTimeModel
}

func (u *UserNamespace) TableName() string {
	return "user_namespace"
}

func (u *UserNamespace) TableUnique() [][]string {
	return [][]string{{"user_id", "namespace_id"}}
}

func UsersOrNamespacesDelete(userIdList []int, namespaceId []int) error {
	u := UserNamespace{}
	q := models.GetORM().QueryTable(u.TableName())
	if userIdList != nil {
		q = q.Filter("user_id__in", userIdList)
	}
	if namespaceId != nil {
		q = q.Filter("namespace_id__in", namespaceId)
	}
	_, err := q.Delete()
	return err
}

// batch remove users from space
func RemoveUsersFromNamespace(namespaceId int, userIds []int) error {
	u := UserNamespace{}
	_, err := models.GetORM().QueryTable(u.TableName()).Filter("namespace_id", namespaceId).Filter("user_id__in", userIds).Delete()
	if err != nil {
		return err
	}
	return nil
}

// clear the user from the space
func ClearUsersFromNamespace(namespaceId int, userIds []int) error {
	u := UserNamespace{}
	_, err := models.GetORM().QueryTable(u.TableName()).Filter("namespace_id", namespaceId).Filter("user_id__in", userIds).Delete()
	if err != nil {
		return err
	}
	return nil
}

// add users in batches in the space
func AddUsersInNamespace(namespaceId int, userIds []int) error {
	var members []*UserNamespace
	for _, userId := range userIds {
		member := &UserNamespace{UserId: userId, NamespaceId: namespaceId}
		members = append(members, member)
	}
	_, err := models.GetORM().InsertMulti(len(members), members)
	if err != nil {
		return err
	}
	return nil
}

// change the permissions of members in the space
func UpdateUserPermissionInNamespace(namespaceId, userId int, permission Permission) error {
	o := models.GetORM()
	member := &UserNamespace{UserId: userId, NamespaceId: namespaceId}
	if o.Read(member) != nil {
		return fmt.Errorf("member not found")
	}
	member.Permission = permission
	if _, err := o.Update(member); err != nil {
		return err
	}
	return nil
}

// batch change permissions of members in a space
func UpdateUsersPermissionInNamespace(namespaceId int, userIds []int, permission Permission) error {
	o, u := models.GetORM(), UserNamespace{}
	_, err := o.QueryTable(u.TableName()).Filter("namespace_id", namespaceId).Filter("user_id__in", userIds).Update(orm.Params{
		"permission": permission,
	})
	if err != nil {
		return err
	}
	return nil
}

func QueryUsers(namespaceId int, userName string, permission int, orderBy string, offset, limit int) ([]*User, int64, error) {
	o, un, u := models.GetORM(), UserNamespace{}, User{}

	var list orm.ParamsList
	unQS := o.QueryTable(un.TableName()).Filter("namespace_id", namespaceId)
	if permission >= 0 {
		unQS = unQS.Filter("permission", permission)
	}

	_, err := unQS.ValuesFlat(&list, "user_id")
	if err != nil {
		return nil, 0, errors.New("no users")
	}
	if idList := cast.ToIntSlice(list); idList == nil {
		return nil, 0, nil
	}

	q := o.QueryTable(u.TableName()).Filter("id__in", cast.ToIntSlice(list))
	if userName != "" {
		q = q.Filter("email__icontains", userName)
	}

	if len(orderBy) > 0 {
		q = q.OrderBy(orderBy)
	}

	total, _ := q.Count()

	if offset < 0 || limit <= 0 {
		return nil, 0, errors.New("invalid offset")
	}

	q = q.Limit(limit, offset)
	var users []*User
	if _, err := q.All(&users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// user move out of space in batches
func UserRemoveNamespace(userId int, namespaceId []int) error {
	u := UserNamespace{}
	_, err := models.GetORM().QueryTable(u.TableName()).Filter("user_id", userId).Filter("namespace_id__in", namespaceId).Delete()
	if err != nil {
		return err
	}
	return nil
}

// user clears space
func UserClearNamespace(userId int) error {
	u := UserNamespace{}
	_, err := models.GetORM().QueryTable(u.TableName()).Filter("user_id", userId).Delete()
	if err != nil {
		return err
	}
	return nil
}
