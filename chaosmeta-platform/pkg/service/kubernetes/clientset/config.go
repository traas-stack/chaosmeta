package clientset

import (
	"encoding/base64"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeLoadMode string

func (k KubeLoadMode) IsEmpty() bool {
	return k == ""
}

func (k KubeLoadMode) String() string {
	return string(k)
}

const (
	KubeLoadFromLocal        KubeLoadMode = "LoadFromLocal"
	KubeLoadInCluster        KubeLoadMode = "LoadInCluster"
	KubeLoadFromBase64Stream KubeLoadMode = "LoadFromAdmin"
	KubeLoadFromAliyun       KubeLoadMode = "LoadFromAliyun"
)

func GetKubeRestConf(mode KubeLoadMode, cfgInfo string) (rfg *rest.Config, err error) {
	switch mode {
	case KubeLoadFromLocal:
		rfg, err = clientcmd.BuildConfigFromFlags("", cfgInfo)
	case KubeLoadInCluster:
		rfg, err = rest.InClusterConfig()
	case KubeLoadFromBase64Stream:
		var bts []byte
		bts, err = base64.StdEncoding.DecodeString(cfgInfo)
		if err != nil {
			break
		}
		var ccfg clientcmd.ClientConfig
		ccfg, err = clientcmd.NewClientConfigFromBytes(bts)
		if err != nil {
			break
		}
		rfg, err = ccfg.ClientConfig()
	default:
		err = fmt.Errorf("invalid kube config load mode: %s or empty", mode)
	}
	if err != nil {
		return
	}

	return
}
