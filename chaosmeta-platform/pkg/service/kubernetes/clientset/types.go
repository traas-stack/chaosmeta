package clientset

import (
	cv1alpha1 "chaosmeta-platform/pkg/gateway/apis/chaosmetacluster/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterDashboardInfo struct {
	Name      string
	NodeCount int
}

type ChaosmetaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ChaosmetaClusterSpec   `json:"spec,omitempty"`
	Status            ChaosmetaClusterStatus `json:"status,omitempty"`
}

type ClusterListResponse struct {
	Total    int                          `json:"total"`
	Current  int                          `json:"current"`
	PageSize int                          `json:"pageSize"`
	List     []cv1alpha1.ChaosmetaCluster `json:"list"`
}

type DefaultParam struct {
	//DefaultClusterName string
	DefaultClusterEnv string
	RegionID          string
	//DefaultImagePullSecret    string
	//DefaultPrometheusEndpoint string
}

type ChaosmetaClusterType string

const (
	LocalChaosmetaCluster  ChaosmetaClusterType = "Local"
	RemoteChaosmetaCluster ChaosmetaClusterType = "Remote"
	ProxyChaosmetaCluster  ChaosmetaClusterType = "Proxy"
)

type ChaosmetaClusterPhase string

const (
	ReadyChaosmetaClusterPhase     ChaosmetaClusterPhase = "Ready"
	FailedChaosmetaClusterPhase    ChaosmetaClusterPhase = "Failed"
	ScalingChaosmetaClusterPhase   ChaosmetaClusterPhase = "Scaling"   // 扩容中
	ScaleFailChaosmetaClusterPhase ChaosmetaClusterPhase = "ScaleFail" // 扩容失败
)

type ChaosmetaClusterSpec struct {
	Type             ChaosmetaClusterType `json:"type,omitempty"`
	Description      string               `json:"description,omitempty"`
	RegionID         string               `json:"regionId,omitempty"`
	KubernetesOption *KubernetesOption    `json:"kubernetesOption,omitempty"`
	CloudOption      CloudConfig          `json:"cloudOption,omitempty"`
}

// CloudConfig
// @Description: 云资源信息,当前主要是阿里云
type CloudConfig struct {
	Provider  string `json:"provider"`
	ClusterID string `json:"clusterID"`
	VPC       string `json:"vpc"`
}

type KubernetesOption struct {
	LoadMode        string `json:"loadMode,omitempty"`
	KubeConf        string `json:"kubeConf,omitempty"`
	IngressEndpoint string `json:"ingressEndpoint,omitepmty"`
}

type ChaosmetaClusterStatus struct {
	Phase                ChaosmetaClusterPhase `json:"phase,omitempty"`
	Reason               string                `json:"reason,omitempty"`
	LastUpdatedTimestamp metav1.Time           `json:"lastUpdatedTimestamp,omitempty"`
}

type ChaosmetaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosmetaCluster `json:"items,omitempty"`
}
