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

package cluster

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models/cluster"
	"chaosmeta-platform/pkg/service/kubernetes/clientset"
	"chaosmeta-platform/util/enc_dec"
	"context"
	"encoding/base64"
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
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

func (c *ClusterService) GetKubeConf(ctx context.Context, id int) (string, error) {
	clusterGet := cluster.Cluster{ID: id}
	if err := cluster.GetClusterById(ctx, &clusterGet); err != nil {
		return "", err
	}
	kubeConf, err := enc_dec.Decrypt([]byte(clusterGet.KubeConfig), []byte(config.DefaultRunOptIns.SecretKey))
	if err != nil {
		return "", err
	}
	return string(kubeConf), nil
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

func (c *ClusterService) GetRestConfig(ctx context.Context, id int) (*kubernetes.Clientset, *rest.Config, error) {
	if config.DefaultRunOptIns.RunMode == "KubeConfig" {
		id = -1
	}
	if config.DefaultRunOptIns.RunMode == "ServiceAccount" {
		id = 0
	}
	if id == 0 {
		return c.getRestConfigInCluster()
	}
	if id == -1 {
		return c.getRestConfigFromKubeConfig("")
	}
	return c.getRestConfigFromClusterId(ctx, id)
}

func (c *ClusterService) getRestConfigInCluster() (*kubernetes.Clientset, *rest.Config, error) {
	restConfig, err := clientset.GetKubeRestConf(clientset.KubeLoadInCluster, "")
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, restConfig, err
	}
	return clientSet, restConfig, err
}

func (c *ClusterService) getRestConfigFromKubeConfig(kubeConfigPath string) (*kubernetes.Clientset, *rest.Config, error) {
	if kubeConfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, err
		}
		kubeConfigPath = filepath.Join(homeDir, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, config, err
	}
	return clientset, config, nil
}

func (c *ClusterService) getRestConfigFromClusterId(ctx context.Context, id int) (*kubernetes.Clientset, *rest.Config, error) {
	cluster, err := c.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	restConfig, err := clientset.GetKubeRestConf(clientset.KubeLoadFromBase64Stream, cluster.KubeConfig)
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, restConfig, err
	}
	return clientSet, restConfig, err
}
