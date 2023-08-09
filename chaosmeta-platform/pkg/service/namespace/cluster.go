package namespace

import (
	"chaosmeta-platform/pkg/models/cluster"
	"chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/util/sort"
	"context"
	"errors"
)

func (s *NamespaceService) SetAttackableCluster(ctx context.Context, namespaceId int, username string, clusterId int) error {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return errors.New("permission denied")
	}
	return namespace.SetClusterIDsForNamespace(namespaceId, []int{clusterId})
}

func (s *NamespaceService) ClearAttackableCluster(ctx context.Context, namespaceId int, username string) error {
	if !s.IsAdmin(ctx, namespaceId, username) {
		return errors.New("permission denied")
	}
	return namespace.ClearClusterIDsForNamespace(namespaceId)
}

func (s *NamespaceService) GetAttackableClusterIDsByNamespaceID(ctx context.Context, namespaceId int) ([]int, error) {
	return namespace.GetClusterIDsByNamespaceID(namespaceId)
}

func (s *NamespaceService) GetAttackableClustersByNamespaceID(ctx context.Context, namespaceId int, orderBy string, page, pageSize int) (int64, []*cluster.Cluster, error) {
	clusterIds, err := namespace.GetClusterIDsByNamespaceID(namespaceId)
	if err != nil {
		return 0, nil, err
	}
	return cluster.GetClustersByIdList(ctx, sort.RemoveDuplicates(clusterIds), orderBy, page, pageSize)
}
