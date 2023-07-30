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
	"chaosmeta-platform/pkg/models"
	namespaceModel "chaosmeta-platform/pkg/models/namespace"
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
	u := &models.User{Email: "admin"}
	if err := models.GetUser(ctx, u); err != nil {
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

func (s *NamespaceService) Create(ctx context.Context, name, description string, creatorName string) error {
	creator := models.User{Email: creatorName}
	if err := models.GetUser(ctx, &creator); err != nil {
		return err
	}

	namespace := &namespaceModel.Namespace{
		Name:        name,
		Description: description,
		Creator:     creator.ID,
	}
	namespaceId, err := namespaceModel.InsertNamespace(ctx, namespace)
	if err != nil {
		return err
	}
	return namespaceModel.AddUsersInNamespace(int(namespaceId), namespaceModel.AddUsersParam{
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

func (s *NamespaceService) GetUsers(ctx context.Context, namespaceId int, userName string, permission int, orderBy string, page, pageSize int) ([]*models.User, int64, error) {
	return namespaceModel.QueryUsers(namespaceId, userName, permission, orderBy, page, pageSize)
}

func (s *NamespaceService) IsAdmin(ctx context.Context, namespaceId int, userName string) bool {
	user := models.User{Email: userName}
	if err := models.GetUser(ctx, &user); err != nil {
		return false
	}
	un := namespaceModel.UserNamespace{
		NamespaceId: namespaceId,
		UserId:      user.ID,
	}
	if err := namespaceModel.GetUserNamespace(ctx, &un); err != nil {
		return false
	}
	if un.Permission == namespaceModel.AdminPermission {
		return true
	}
	return false
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
