package inject

import (
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
)

func (i *InjectService) ListScopes(ctx context.Context, orderBy string, page, pageSize int) (int64, []basic.Scope, error) {
	total, scopes, err := basic.ListScopes(ctx, orderBy, page, pageSize)
	return total, scopes, err
}
