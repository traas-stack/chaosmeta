package inject

import (
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

func (i *InjectService) ListFault(ctx context.Context, targetId int, orderBy string, page, pageSize int) (int64, []basic.Fault, error) {
	total, targets, err := basic.ListFaults(ctx, targetId, orderBy, page, pageSize)
	return total, targets, err
}
