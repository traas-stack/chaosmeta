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
	"chaosmeta-platform/pkg/models/experiment_instance"
	namespaceModel "chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/models/user"
	"chaosmeta-platform/util/log"
	"context"
	"errors"
)

func Init() {
	namespace := namespaceModel.Namespace{}
	ctx := context.Background()

	if err := namespaceModel.GetDefaultNamespace(ctx, &namespace); err == nil {
		return
	}

	defaultNamespace := &namespaceModel.Namespace{
		Name:        "default",
		Description: "This is the default namespace",
		Creator:     1,
		IsDefault:   true,
	}
	_, err := namespaceModel.InsertNamespace(ctx, defaultNamespace)
	if err != nil {
		log.Panic(err)
	}
	u := &user.User{Email: "admin"}
	if err := user.GetUser(ctx, u); err != nil {
		log.Panic(err)
	}

	if err := namespaceModel.AddUsersInNamespace(defaultNamespace.Id, namespaceModel.AddUsersParam{
		Users: []namespaceModel.UserData{{
			Id:         u.ID,
			Permission: int(namespaceModel.AdminPermission),
		}},
	}); err != nil {
		log.Error(err)
	}
}

type NamespaceService struct{}

func (s *NamespaceService) Create(ctx context.Context, name, description string, creatorName string) (int64, error) {
	creator := user.User{Email: creatorName}
	if err := user.GetUser(ctx, &creator); err != nil {
		return 0, err
	}

	namespace := &namespaceModel.Namespace{
		Name:        name,
		Description: description,
		Creator:     creator.ID,
	}
	namespaceId, err := namespaceModel.InsertNamespace(ctx, namespace)
	if err != nil {
		return 0, err
	}
	return namespaceId, namespaceModel.AddUsersInNamespace(int(namespaceId), namespaceModel.AddUsersParam{
		Users: []namespaceModel.UserData{{
			Id:         creator.ID,
			Permission: int(namespaceModel.AdminPermission),
		}},
	})
}

func (s *NamespaceService) Update(ctx context.Context, userName string, namespaceId int, namespaceName string, namespaceDescription string) error {
	if !s.IsAdmin(ctx, namespaceId, userName) {
		return errors.New("permission denied")
	}

	namespace := namespaceModel.Namespace{Id: namespaceId}
	if err := namespaceModel.GetNamespaceById(ctx, &namespace); err != nil {
		return err
	}

	if namespaceName != "" {
		namespace.Name = namespaceName
	}
	if namespaceDescription != "" {
		namespace.Description = namespaceDescription
	}

	if _, err := namespaceModel.UpdateNamespace(ctx, &namespace); err != nil {
		return err
	}
	return nil
}

func (s *NamespaceService) Get(ctx context.Context, id int) (*namespaceModel.Namespace, error) {
	namespace := namespaceModel.Namespace{Id: id}
	if err := namespaceModel.GetNamespaceById(ctx, &namespace); err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (s *NamespaceService) GetList(ctx context.Context, name, creator, orderBy string, page, pageSize int) (int64, []namespaceModel.Namespace, error) {
	return namespaceModel.QueryNamespaces(ctx, name, creator, orderBy, page, pageSize)
}

func (s *NamespaceService) Delete(ctx context.Context, userName string, namespaceId int) error {
	if s.IsDefault(ctx, namespaceId) {
		return errors.New("default namespace, remove users are not allowed")
	}
	if !s.IsAdmin(ctx, namespaceId, userName) {
		return errors.New("permission denied")
	}
	namespace := namespaceModel.Namespace{Id: namespaceId}
	if err := namespaceModel.GetNamespaceById(ctx, &namespace); err != nil {
		return errors.New("namespace not found")
	}

	if _, err := namespaceModel.DeleteNamespace(ctx, namespaceId); err != nil {
		return err
	}
	if err := namespaceModel.ClearClusterIDsForNamespace(namespaceId); err != nil {
		return err
	}
	return namespaceModel.UsersOrNamespacesDelete(nil, []int{namespaceId})
}

func (s *NamespaceService) GetAll(ctx context.Context) ([]*namespaceModel.Namespace, error) {
	namespaces, err := namespaceModel.GetAllNamespaces()
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

func (s *NamespaceService) AddUsers(ctx context.Context, userName string, namespaceId int, addUsersParam namespaceModel.AddUsersParam) error {
	if s.IsDefault(ctx, namespaceId) {
		return errors.New("default namespace, add users are not allowed")
	}
	if !s.IsAdmin(ctx, namespaceId, userName) {
		return errors.New("permission denied")
	}
	return namespaceModel.AddUsersInNamespace(namespaceId, addUsersParam)
}

func (s *NamespaceService) GroupedUserInNamespaces(ctx context.Context, namespaceId int, namespaceName string, userName string, permission int, orderBy string, page, pageSize int) (int64, []namespaceModel.NamespaceData, error) {
	return namespaceModel.GroupedUserInNamespaces(namespaceId, namespaceName, userName, permission, orderBy, page, pageSize)
}

type NamespaceData struct {
	NamespaceInfo  namespaceModel.NamespaceInfo      `json:"namespaceInfo"`
	users          []namespaceModel.NamespaceData    `json:"users"`
	UserData       UserDataInNamespace               `json:"userData"`
	ExperimentData ExperimentInstanceDataInNamespace `json:"experimentrData"`
}

type UserDataInNamespace struct {
	Permission int                                `json:"permission"`
	ToTal      int                                `json:"toTal"`
	Users      []namespaceModel.UserDataNamespace `json:"users"`
}

func (s *NamespaceService) getUserData(ctx context.Context, userId int, namespaceId int) (UserDataInNamespace, error) {
	userDataInNamespace := UserDataInNamespace{}
	permission, err := namespaceModel.IsUserInNamespace(userId, namespaceId)
	if err != nil {
		return userDataInNamespace, err
	}
	userDataInNamespace.Permission = int(permission)

	total, users, err := namespaceModel.GetUsersFromNamespace(ctx, namespaceId)
	userDataInNamespace.ToTal = int(total)
	userDataInNamespace.Users = append(userDataInNamespace.Users, users...)
	return userDataInNamespace, err
}

type ExperimentInstanceDataInNamespace struct {
	ToTal     int64            `json:"toTal"`
	StatusMap map[string]int64 `json:"statusCount"`
}

func (s *NamespaceService) getExperimentInstanceData(ctx context.Context, namespaceId int) (ExperimentInstanceDataInNamespace, error) {
	statusMap, total, err := experiment_instance.CountExperimentInstance(namespaceId, 0)
	return ExperimentInstanceDataInNamespace{ToTal: total, StatusMap: statusMap}, err
}

func (s *NamespaceService) getUserAndExperimentInstanceData(ctx context.Context, users []namespaceModel.NamespaceData, namespaceData *NamespaceData) error {
	if users == nil {
		return nil
	}
	if namespaceData == nil {
		return errors.New("namespaceData is nil")
	}
	var err error
	for _, user := range users {
		namespaceData.UserData, err = s.getUserData(ctx, user.Id, namespaceData.NamespaceInfo.ID)
		if err != nil {
			return err
		}
		namespaceData.ExperimentData, err = s.getExperimentInstanceData(ctx, namespaceData.NamespaceInfo.ID)
		if err != nil {
			return err
		}
	}
	return err
}

// 未加入
func (s *NamespaceService) GroupNamespacesUserNotIn(ctx context.Context, userId int, namespaceName string, userName string, orderBy string, page, pageSize int) (int64, []NamespaceData, error) {
	total, namesapceList, err := namespaceModel.GroupNamespacesUserNotIn(userId, namespaceName, userName, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	var namespaceDataList []NamespaceData
	for _, namespace := range namesapceList {
		namespaceData := NamespaceData{NamespaceInfo: namespace}
		_, userList, err := namespaceModel.GroupedUserInNamespaces(namespace.ID, "", "", -1, "", 1, 20)
		if err != nil {
			continue
		}
		if err := s.getUserAndExperimentInstanceData(ctx, userList, &namespaceData); err == nil {
			namespaceDataList = append(namespaceDataList, namespaceData)
		}
	}
	return total, namespaceDataList, nil
}

// 全部
func (s *NamespaceService) GroupAllNamespaces(ctx context.Context, namespaceName string, userName string, orderBy string, page, pageSize int) (int64, []NamespaceData, error) {
	total, namesapceList, err := namespaceModel.GroupAllNamespaces(namespaceName, userName, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	var namespaceDataList []NamespaceData
	for _, namespace := range namesapceList {
		namespaceData := NamespaceData{NamespaceInfo: namespace}
		var err error
		_, userList, err := namespaceModel.GroupedUserInNamespaces(namespace.ID, "", "", -1, "", 1, 20)
		if err != nil {
			continue
		}
		if err := s.getUserAndExperimentInstanceData(ctx, userList, &namespaceData); err == nil {
			namespaceDataList = append(namespaceDataList, namespaceData)
		}
	}
	return total, namespaceDataList, nil
}

// 已加入
func (s *NamespaceService) GroupNamespacesByUsername(ctx context.Context, userId int, namespace string, userName string, permission int, orderBy string, page, pageSize int) (int64, []NamespaceData, error) {
	total, namesapceList, err := namespaceModel.GroupAllNamespacesByUserName(userId, namespace, userName, permission, orderBy, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	var namespaceDataList []NamespaceData
	for _, namespace := range namesapceList {
		namespaceData := NamespaceData{NamespaceInfo: namespace}
		_, userList, err := namespaceModel.GroupedUserInNamespaces(namespace.ID, "", "", -1, "", 1, 20)
		if err != nil {
			continue
		}
		if err := s.getUserAndExperimentInstanceData(ctx, userList, &namespaceData); err == nil {
			namespaceDataList = append(namespaceDataList, namespaceData)
		}
	}
	return total, namespaceDataList, nil
}

func (s *NamespaceService) DefaultAddUsers(ctx context.Context, addUsersParam namespaceModel.AddUsersParam) error {
	namespace := namespaceModel.Namespace{}
	if err := namespaceModel.GetDefaultNamespace(ctx, &namespace); err != nil {
		return err
	}
	return namespaceModel.AddUsersInNamespace(namespace.Id, addUsersParam)
}

func (s *NamespaceService) RemoveUsers(ctx context.Context, userName string, userIds []int, namespaceId int) error {
	if s.IsDefault(ctx, namespaceId) {
		return errors.New("default namespace, remove users are not allowed")
	}
	if !s.IsAdmin(ctx, namespaceId, userName) {
		return errors.New("permission denied")
	}
	return namespaceModel.RemoveUsersFromNamespace(namespaceId, userIds)
}

func (s *NamespaceService) ChangeUsersPermission(ctx context.Context, userName string, userIds []int, namespaceId int, permission namespaceModel.Permission) error {
	if s.IsDefault(ctx, namespaceId) {
		return errors.New("default namespace, permission changes are not allowed")
	}
	if !s.IsAdmin(ctx, namespaceId, userName) {
		return errors.New("permission denied")
	}
	return namespaceModel.UpdateUsersPermissionInNamespace(namespaceId, userIds, permission)
}

func (s *NamespaceService) IsAdmin(ctx context.Context, namespaceId int, userName string) bool {
	userGet := user.User{Email: userName}
	if err := user.GetUser(ctx, &userGet); err != nil {
		return false
	}
	if userGet.Role == user.AdminRole {
		return true
	}
	un := namespaceModel.UserNamespace{
		NamespaceId: namespaceId,
		UserId:      userGet.ID,
	}
	if err := namespaceModel.GetUserNamespace(&un); err != nil {
		return false
	}
	if un.Permission == namespaceModel.AdminPermission {
		return true
	}
	return false
}

func (s *NamespaceService) IsUserJoin(ctx context.Context, namespaceId int, userId int) (bool, int) {
	userGet := user.User{ID: userId}
	if err := user.GetUserById(ctx, &userGet); err != nil {
		return false, -1
	}
	if userGet.Role == user.AdminRole {
		return true, 1
	}
	un := namespaceModel.UserNamespace{
		NamespaceId: namespaceId,
		UserId:      userGet.ID,
	}
	if err := namespaceModel.GetUserNamespace(&un); err != nil {
		return false, -1
	}
	return true, int(un.Permission)
}

type UserInfoInNamespace struct {
	User       user.User
	IsJoin     bool `json:"isJoin"`
	Permission int  `json:"permission"`
}

func (s *NamespaceService) GetUsersOfNamespacePermissions(ctx context.Context, users []user.User, namespaceId int) ([]UserInfoInNamespace, error) {
	var userInfoInNamespaces []UserInfoInNamespace
	for _, user := range users {
		userInfoInNamespace := UserInfoInNamespace{
			User:   user,
			IsJoin: false,
		}
		isJoin, permission := s.IsUserJoin(ctx, namespaceId, user.ID)
		if isJoin {
			userInfoInNamespace.IsJoin = isJoin
			userInfoInNamespace.Permission = permission
		}
		userInfoInNamespaces = append(userInfoInNamespaces, userInfoInNamespace)
	}
	return userInfoInNamespaces, nil
}

func (s *NamespaceService) IsDefault(ctx context.Context, namespaceId int) bool {
	namespace := namespaceModel.Namespace{}
	if err := namespaceModel.GetDefaultNamespace(ctx, &namespace); err != nil {
		return false
	}
	if namespaceId == namespace.Id {
		return true
	}
	return false
}
