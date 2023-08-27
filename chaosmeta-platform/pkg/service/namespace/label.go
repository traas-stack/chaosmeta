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
	namespaceModel "chaosmeta-platform/pkg/models/namespace"
	"context"
	"errors"
)

func (s *NamespaceService) CreateLabel(ctx context.Context, namespaceId int, username, name, color string) (int64, error) {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return 0, errors.New("permission denied")
	}
	label := namespaceModel.Label{Name: name, NamespaceId: namespaceId, Color: color, Creator: username}
	if err := namespaceModel.GetLabelByName(ctx, &label); err == nil {
		return int64(label.Id), errors.New("label already exists")
	}
	return namespaceModel.InsertLabel(ctx, &label)
}

func (s *NamespaceService) DeleteLabel(ctx context.Context, namespaceId int, username string, labelId int) error {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return errors.New("permission denied")
	}
	label := namespaceModel.Label{Id: labelId}
	if err := namespaceModel.GetLabelById(ctx, &label); err != nil {
		return err
	}
	_, err := namespaceModel.DeleteLabel(ctx, labelId)
	return err
}

func (s *NamespaceService) ListLabel(ctx context.Context, nameSpaceId int, name, creator, orderBy string, page, pageSize int) (int64, []namespaceModel.Label, error) {
	return namespaceModel.QueryLabels(ctx, nameSpaceId, name, creator, orderBy, page, pageSize)
}

func (s *NamespaceService) GetLabelByName(ctx context.Context, nameSpaceId int, name string) (namespaceModel.Label, error) {
	label := namespaceModel.Label{
		Name:        name,
		NamespaceId: nameSpaceId,
	}
	return label, namespaceModel.GetLabelByName(ctx, &label)
}
