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
