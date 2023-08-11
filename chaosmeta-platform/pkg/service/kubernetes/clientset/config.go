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
