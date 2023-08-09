package inject

import (
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

func (i *InjectService) ListArg(ctx context.Context, faultId int, orderBy string, page, pageSize int) (int64, []basic.Args, error) {
	total, targets, err := basic.ListArgs(ctx, "", faultId, orderBy, page, pageSize)
	return total, targets, err
}
