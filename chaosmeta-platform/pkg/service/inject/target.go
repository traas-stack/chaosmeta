package inject

import (
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

func (i *InjectService) ListTargets(ctx context.Context, scopeId int, orderBy string, page, pageSize int) (int64, []basic.Target, error) {
	total, targets, err := basic.ListTargets(ctx, scopeId, orderBy, page, pageSize)
	return total, targets, err
}
