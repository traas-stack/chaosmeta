package cluster

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models/cluster"
	"chaosmeta-platform/util/enc_dec"
	"context"
	"encoding/base64"
	"errors"
)

type ClusterService struct{}

func (c *ClusterService) Create(ctx context.Context, name, kubeConfig string) (int64, error) {
	kubeConfigByte, err := base64.StdEncoding.DecodeString(kubeConfig)
	if err != nil {
		return 0, err
	}

	encryptedkubeConfig, err := enc_dec.Encrypt(kubeConfigByte, []byte(config.DefaultRunOptIns.SecretKey))
	if err != nil {
		return 0, err
	}

	insertCluster := cluster.Cluster{
		Name:       name,
		KubeConfig: string(encryptedkubeConfig),
	}
	if err := cluster.GetClusterByName(ctx, &insertCluster); err == nil {
		return 0, errors.New("cluster already exists")
	}

	return cluster.InsertCluster(ctx, &insertCluster)
}

func (c *ClusterService) Get(ctx context.Context, id int) (*cluster.Cluster, error) {
	clusterGet := cluster.Cluster{ID: id}
	if err := cluster.GetClusterById(ctx, &clusterGet); err != nil {
		return nil, err
	}
	kubeConf, err := enc_dec.Decrypt([]byte(clusterGet.KubeConfig), []byte(config.DefaultRunOptIns.SecretKey))
	if err != nil {
		return nil, err
	}

	clusterGet.KubeConfig = base64.StdEncoding.EncodeToString(kubeConf)
	return &clusterGet, nil
}

func (c *ClusterService) Update(ctx context.Context, id int, name, kubeConfig string) error {
	insertCluster := cluster.Cluster{
		ID: id,
	}
	if err := cluster.GetClusterById(ctx, &insertCluster); err != nil {
		return err
	}

	if name != "" {
		insertCluster.Name = name
	}

	if kubeConfig != "" {
		kubeConfigByte, err := base64.StdEncoding.DecodeString(kubeConfig)
		if err != nil {
			return err
		}

		encryptedkubeConfig, err := enc_dec.Encrypt(kubeConfigByte, []byte(config.DefaultRunOptIns.SecretKey))
		if err != nil {
			return err
		}
		insertCluster.KubeConfig = string(encryptedkubeConfig)
	}
	_, err := cluster.UpdateCluster(ctx, &insertCluster)
	return err
}

func (c *ClusterService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return cluster.DeleteClustersByIdList(ctx, []int{id})
}

func (c *ClusterService) DeleteList(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return errors.New("invalid id list")
	}
	return cluster.DeleteClustersByIdList(ctx, ids)
}

func (c *ClusterService) GetList(ctx context.Context, name, orderBy string, page, pageSize int) (int64, []cluster.Cluster, error) {
	return cluster.QueryCluster(ctx, name, "", orderBy, page, pageSize)
}
