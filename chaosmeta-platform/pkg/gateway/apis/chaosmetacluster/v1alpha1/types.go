package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:openapi-gen=true
// +genclient:nonNamespaced
// +kubebuilder:resource:scope=Cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ChaosmetaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ChaosmetaClusterSpec   `json:"spec,omitempty"`
	Status            ChaosmetaClusterStatus `json:"status,omitempty"`
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
	MonitorOption    *MonitorOption       `json:"monitorOption,omitempty"`
	EdasOption       *EdasOption          `json:"edasOption,omitempty"`
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

type MonitorOption struct {
	PrometheusEndpoint string `json:"prometheusEndpoint,omitempty"`
}

type EdasOption struct {
	ClusterID string `json:"clusterId"`
	//Workspace string `json:"workspace"`
	Endpoint string `json:"endpoint"`
}

type ChaosmetaClusterStatus struct {
	Phase                ChaosmetaClusterPhase `json:"phase,omitempty"`
	Reason               string                `json:"reason,omitempty"`
	LastUpdatedTimestamp metav1.Time           `json:"lastUpdatedTimestamp,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ChaosmetaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosmetaCluster `json:"items,omitempty"`
}
