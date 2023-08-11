/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package namespace

import (
	"chaosmeta-platform/pkg/models/common"
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/spf13/cast"
)

type Permission int

const (
	NormalPermission = Permission(0) //只读
	AdminPermission  = Permission(1) //管理员
)

type UserNamespace struct {
	ID          int        `json:"id" orm:"pk;auto;column(id)"`
	UserId      int        `json:"userId" orm:"column(user_id);index"`
	NamespaceId int        `json:"namespaceId" orm:"column(namespace_id);index"`
	Permission  Permission `json:"permission" orm:"column(permission);default(0);index"`
	models.BaseTimeModel
}

func (u *UserNamespace) TableName() string {
	return "user_namespace"
}

func (u *UserNamespace) TableUnique() [][]string {
	return [][]string{{"user_id", "namespace_id"}}
}

func GetUserNamespace(u *UserNamespace) error {
	return models.GetORM().Read(u, "user_id", "namespace_id")
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

type UserNamespaceData struct {
	NamespaceId int        `json:"namespaceId"`
	Permission  Permission `json:"permission"`
}

type UserDataNamespace struct {
	UserId     int        `json:"userId"`
	Permission Permission `json:"permission"`
}

func GetNamespacesFromUser(ctx context.Context, userId int, permission int, orderBy string, page, pageSize int) (int64, []UserNamespaceData, error) {
	u := UserNamespace{}
	var (
		namespaceDataList []UserNamespaceData
		userNamespaces    []*UserNamespace
	)
	qs := models.GetORM().QueryTable(u.TableName()).Filter("user_id", userId)
	if permission >= 0 {
		qs = qs.Filter("permission", permission)
	}
	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	totalCount, err := qs.Count()
	if err != nil {
		return totalCount, nil, err
	}

	qs = qs.Limit(pageSize, (page-1)*pageSize)

	if _, err := qs.All(&userNamespaces); err != nil {
		return 0, nil, err
	}
	for _, userNamespace := range userNamespaces {
		namespaceDataList = append(namespaceDataList, UserNamespaceData{
			NamespaceId: userNamespace.NamespaceId,
			Permission:  userNamespace.Permission,
		})
	}
	return totalCount, namespaceDataList, nil
}

func GetUsersFromNamespace(ctx context.Context, namespaceId int) (int64, []UserDataNamespace, error) {
	u := UserNamespace{}
	var (
		userNamespaces []*UserNamespace
		userDataList   []UserDataNamespace
	)
	qs := models.GetORM().QueryTable(u.TableName()).Filter("namespace_id", namespaceId)

	totalCount, err := qs.Count()
	if err != nil {
		return totalCount, nil, err
	}

	if _, err := qs.All(&userNamespaces); err != nil {
		return 0, nil, err
	}
	for _, userNamespace := range userNamespaces {
		userDataList = append(userDataList, UserDataNamespace{
			UserId:     userNamespace.UserId,
			Permission: userNamespace.Permission,
		})
	}
	return totalCount, userDataList, nil
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

type UserData struct {
	Id         int `json:"id"`
	Permission int `json:"permission"`
}

type AddUsersParam struct {
	Users []UserData `json:"users"`
}

// add users in batches in the space
func AddUsersInNamespace(namespaceId int, addUsersParam AddUsersParam) error {
	var members []*UserNamespace
	for _, user := range addUsersParam.Users {
		member := &UserNamespace{UserId: user.Id, NamespaceId: namespaceId, Permission: Permission(user.Permission)}
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
	if userIds == nil {
		return errors.New("user list is nil")
	}
	o, u := models.GetORM(), UserNamespace{}
	_, err := o.QueryTable(u.TableName()).Filter("namespace_id", namespaceId).Filter("user_id__in", userIds).Update(orm.Params{
		"permission": permission,
	})
	if err != nil {
		return err
	}
	return nil
}

type NamespaceData struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	CreateTime string `json:"create_time"`
}

func GroupedUserInNamespaces(namespaceId int, namespaceName string, userName string, permission int, orderBy string, page, pageSize int) (int64, []NamespaceData, error) {
	sql := `
		SELECT 
			un.create_time,
			un.permission, 
			un.user_id,
			u.email,
			ns.id AS namespace_id,
			ns.name AS namespace_name
		FROM 
			user_namespace un 
			INNER JOIN user u ON un.user_id = u.id 
			INNER JOIN namespace ns ON un.namespace_id = ns.id
	`
	where := " WHERE 1=1 "
	if namespaceId > 0 {
		where += fmt.Sprintf(" AND ns.id = %d", namespaceId)
	}
	if namespaceName != "" {
		where += fmt.Sprintf(" AND ns.name LIKE '%%%s%%'", namespaceName)
	}
	if userName != "" {
		where += fmt.Sprintf(" AND u.email LIKE '%%%s%%'", userName)
	}
	if permission >= 0 {
		where += fmt.Sprintf(" AND un.permission = %d", permission)
	}
	sql += where

	var order string
	if orderBy != "" {
		if orderBy[0] == '-' {
			orderBy = orderBy[1:]
			order = " DESC"
		} else {
			order = " ASC"
		}
		sql += fmt.Sprintf(" ORDER BY ns.id, un.user_id, un.%s %s", orderBy, order)
	} else {
		sql += " ORDER BY ns.id, un.user_id, un.create_time ASC "
	}

	countSql := "SELECT COUNT(*) FROM (" + sql + ") as t"
	var totalCount int64
	o := models.GetORM()
	if err := o.Raw(countSql).QueryRow(&totalCount); err != nil {
		return 0, nil, err
	}
	sql += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	var rows []orm.Params
	if _, err := o.Raw(sql).Values(&rows); err != nil {
		return 0, nil, err
	}

	groupedData := make(map[string][]orm.Params)
	for _, row := range rows {
		key := fmt.Sprintf("%s-%d-%s", row["namespace_name"], row["user_id"], row["create_time"])
		groupedData[key] = append(groupedData[key], row)
	}

	var namespaceDataList []NamespaceData
	for _, data := range groupedData {
		namespaceDataList = append(namespaceDataList, NamespaceData{
			Id:         cast.ToInt(data[0]["user_id"]),
			Name:       cast.ToString(data[0]["email"]),
			Permission: cast.ToString(data[0]["permission"]),
			CreateTime: cast.ToString(data[0]["create_time"]),
		})
	}
	return totalCount, namespaceDataList, nil
}

type NamespaceInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	CreateTime  string `json:"create_time"`
}

func GroupNamespacesUserNotIn(userId int, namespaceName string, userName string, orderBy string, page, pageSize int) (int64, []NamespaceInfo, error) {
	sqlSelect := "SELECT n.id, n.name, un.permission, n.create_time, u.email"
	sqlFrom := "FROM namespace n"
	sqlJoin := "LEFT JOIN user_namespace un ON n.id = un.namespace_id AND un.user_id = ? " + "JOIN user u ON n.creator = u.id"
	sqlWhere := "WHERE n.creator != ? AND n.id NOT IN " +
		"(SELECT un.namespace_id FROM user_namespace un WHERE un.user_id = ? AND un.permission >= 0)"
	var params []interface{}
	params = append(params, userId, userId, userId)
	if namespaceName != "" {
		sqlWhere += " AND n.name LIKE ?"
		params = append(params, "%"+namespaceName+"%")
	}
	if userName != "" {
		sqlWhere += " AND u.email LIKE ?"
		params = append(params, "%"+userName+"%")
	}

	sqlGroupBy := "GROUP BY n.create_time DESC"
	var order string
	if orderBy != "" {
		if orderBy[0] == '-' {
			orderBy = orderBy[1:]
			order = " DESC"
		} else {
			order = " ASC"
		}
		sqlGroupBy = fmt.Sprintf(" ORDER BY n.%s %s", orderBy, order)
	} else {
		sqlGroupBy = " ORDER BY n.create_time ASC "
	}

	o := models.GetORM()
	var results []orm.Params
	sqlCount := "SELECT COUNT(*) " + sqlFrom + " " + sqlJoin + " " + sqlWhere + " " + sqlGroupBy
	var total int64
	if err := o.Raw(sqlCount, params...).QueryRow(&total); err != nil {
		return 0, nil, err
	}

	sqlQuery := sqlSelect + " " + sqlFrom + " " + sqlJoin + " " + sqlWhere + " " + sqlGroupBy
	sqlQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)
	if _, err := o.Raw(sqlQuery, params...).Values(&results); err != nil {
		return 0, nil, err
	}
	var namespaceInfos []NamespaceInfo
	for _, result := range results {
		info := &NamespaceInfo{}
		info.ID = cast.ToInt(result["id"].(string))
		info.Name = result["name"].(string)
		info.CreateTime = result["create_time"].(string)
		namespaceInfos = append(namespaceInfos, *info)
	}
	return total, namespaceInfos, nil
}

func GroupAllNamespacesByUserName(userId int, namespaceName string, userName string, permission int, orderBy string, page, pageSize int) (int64, []NamespaceInfo, error) {
	sqlSelect := "SELECT DISTINCT n.id, n.name, n.description, n.create_time"
	sqlFrom := "FROM namespace n"
	sqlJoin := "LEFT JOIN user_namespace un ON n.id = un.namespace_id "
	sqlWhere := "WHERE 1 = 1"
	var params []interface{}
	if namespaceName != "" {
		sqlWhere += " AND n.name LIKE ?"
		params = append(params, "%"+namespaceName+"%")
	}
	if userId > 0 {
		sqlWhere += fmt.Sprintf(" AND n.id = %d", userId)
	}
	if userName != "" {
		sqlJoin += " JOIN user u ON un.user_id = u.id"
		sqlWhere += " AND un.permission >= 0 AND u.email LIKE ?"
		params = append(params, "%"+userName+"%")
	}
	var (
		order      string
		sqlGroupBy string
	)
	if orderBy != "" {
		if orderBy[0] == '-' {
			orderBy = orderBy[1:]
			order = " DESC"
		} else {
			order = " ASC"
		}
		sqlGroupBy = fmt.Sprintf(" ORDER BY n.%s %s", orderBy, order)
	} else {
		sqlGroupBy = " ORDER BY n.create_time ASC "
	}
	if permission >= 0 {
		sqlWhere += fmt.Sprintf(" AND un.permission = %d", permission)
	}

	o := models.GetORM()
	var results []orm.Params
	sqlCount := "SELECT COUNT(DISTINCT n.id) " + sqlFrom + " " + sqlJoin + " " + sqlWhere + " " + sqlGroupBy
	var total int64
	if err := o.Raw(sqlCount, params...).QueryRow(&total); err != nil {
		return 0, nil, err
	}

	sqlQuery := sqlSelect + " " + sqlFrom + " " + sqlJoin + " " + sqlWhere
	sqlQuery += sqlGroupBy + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	if _, err := o.Raw(sqlQuery, params...).Values(&results); err != nil {
		return 0, nil, err
	}
	var namespaceInfos []NamespaceInfo
	for _, result := range results {
		info := &NamespaceInfo{}
		info.ID = cast.ToInt(result["id"].(string))
		info.Name = result["name"].(string)
		info.CreateTime = result["create_time"].(string)
		namespaceInfos = append(namespaceInfos, *info)
	}
	return total, namespaceInfos, nil

}

func GroupAllNamespaces(namespaceName string, userName string, orderBy string, page, pageSize int) (int64, []NamespaceInfo, error) {
	sqlSelect := "SELECT n.id, n.name, n.description, u.email, n.create_time"
	sqlFrom := "FROM namespace n"
	sqlJoin := "JOIN user u ON n.creator = u.id"
	sqlWhere := "WHERE 1 = 1"
	var params []interface{}
	if namespaceName != "" {
		sqlWhere += " AND n.name LIKE ?"
		params = append(params, "%"+namespaceName+"%")
	}

	if userName != "" {
		return GroupAllNamespacesByUserName(-1, namespaceName, userName, -1, orderBy, page, pageSize)
	}

	sqlQuery := sqlSelect + " " + sqlFrom + " " + sqlJoin + " " + sqlWhere
	sqlCount := "SELECT COUNT(*) " + sqlFrom + " " + sqlJoin + " " + sqlWhere
	var total int64
	o := models.GetORM()
	var results []orm.Params
	if err := o.Raw(sqlCount, params...).QueryRow(&total); err != nil {
		return 0, nil, err
	}

	sqlGroupBy := "GROUP BY n.create_time DESC"
	var order string
	if orderBy != "" {
		if orderBy[0] == '-' {
			orderBy = orderBy[1:]
			order = " DESC"
		} else {
			order = " ASC"
		}
		sqlGroupBy = fmt.Sprintf(" ORDER BY n.%s %s", orderBy, order)
	} else {
		sqlGroupBy = " ORDER BY n.create_time ASC "
	}

	sqlQuery += sqlGroupBy + fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)

	if _, err := o.Raw(sqlQuery, params...).Values(&results); err != nil {
		return 0, nil, err
	}
	var namespaceInfos []NamespaceInfo
	for _, result := range results {
		info := &NamespaceInfo{}
		info.ID = cast.ToInt(result["id"].(string))
		info.Name = result["name"].(string)
		info.CreateTime = result["create_time"].(string)
		namespaceInfos = append(namespaceInfos, *info)
	}
	return total, namespaceInfos, nil
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
