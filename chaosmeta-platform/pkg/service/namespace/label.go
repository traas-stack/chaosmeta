package namespace

import (
	namespaceModel "chaosmeta-platform/pkg/models/namespace"
	"context"
	"errors"
)

func (s *NamespaceService) CreateLabel(ctx context.Context, namespaceId int, username, name string) (int64, error) {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return 0, errors.New("permission denied")
	}
	label := namespaceModel.Label{Name: name, NamespaceId: namespaceId}
	if err := namespaceModel.GetLabelByName(ctx, &label); err != nil {
		return 0, err
	}
	return namespaceModel.InsertLabel(ctx, &label)
}

func (s *NamespaceService) DeleteLabel(ctx context.Context, namespaceId int, username string, labelId int) error {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return errors.New("permission denied")
	}
	label := namespaceModel.Label{Id: labelId}
	if err := namespaceModel.GetLabelByName(ctx, &label); err != nil {
		return err
	}
	_, err := namespaceModel.DeleteLabel(ctx, labelId)
	return err
}

func (s *NamespaceService) ListLabel(ctx context.Context, nameSpaceId int, name, orderBy string, page, pageSize int) (int64, []namespaceModel.Label, error) {
	return namespaceModel.QueryLabels(ctx, nameSpaceId, name, orderBy, page, pageSize)
}
