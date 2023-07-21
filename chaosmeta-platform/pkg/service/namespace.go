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

package service

import (
	"chaosmeta-platform/pkg/models"
	"context"
	"errors"
)

type NamespaceService struct{}

func (s *NamespaceService) CreateNamespace(ctx context.Context, name, description string, creatorName string) error {
	creator := models.User{Email: creatorName}
	if err := models.GetUser(ctx, &creator); err != nil {
		return err
	}

	namespace := &models.Namespace{
		Name:        name,
		Description: description,
		Creator:     creator.ID,
	}
	if _, err := models.InsertNamespace(ctx, namespace); err != nil {
		return err
	}
	return nil
}

func (s *NamespaceService) UpdateNamespace(ctx context.Context, id int, name, description string) error {
	namespace := models.Namespace{Id: id}

	if err := models.GetNamespace(ctx, &namespace); err != nil {
		return err
	}

	namespace.Name = name
	namespace.Description = description
	if _, err := models.UpdateNamespace(ctx, &namespace); err != nil {
		return err
	}
	return nil
}

func (s *NamespaceService) GetNamespace(ctx context.Context, id int) (*models.Namespace, error) {
	namespace := models.Namespace{Id: id}
	if err := models.GetNamespace(ctx, &namespace); err != nil {
		return nil, err
	}
	return &namespace, nil
}

func (s *NamespaceService) DeleteNamespace(ctx context.Context, id int) error {
	namespace := models.Namespace{Id: id}
	if err := models.GetNamespace(ctx, &namespace); err != nil {
		return errors.New("namespace not found")
	}
	if _, err := models.DeleteNamespace(ctx, id); err != nil {
		return err
	}
	return models.UsersOrNamespacesDelete(nil, []int{id})
}

func (s *NamespaceService) GetAllNamespaces(ctx context.Context) ([]*models.Namespace, error) {
	namespaces, err := models.GetAllNamespaces()
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

func (s *NamespaceService) AddUsers(ctx context.Context, userIds []int, namespaceId int) error {
	return models.AddUsersInNamespace(namespaceId, userIds)
}

func (s *NamespaceService) RemoveUsers(ctx context.Context, userIds []int, namespaceId int) error {
	return models.RemoveUsersFromNamespace(namespaceId, userIds)
}

func (s *NamespaceService) ChangeUsersPermission(ctx context.Context, userIds []int, namespaceId int, permission models.Permission) error {
	return models.UpdateUsersPermissionInNamespace(namespaceId, userIds, permission)
}

func (s *NamespaceService) GetUsers(ctx context.Context, namespaceId int, userName string, permission int, orderBy string, offset, limit int) ([]*models.User, int64, error) {
	return models.QueryUsers(namespaceId, userName, permission, orderBy, offset, limit)
}
