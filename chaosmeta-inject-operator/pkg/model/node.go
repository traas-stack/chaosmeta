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

package model

import (
	"fmt"
	"strings"
)

type NodeObject struct {
	NodeName         string
	NodeInternalIP   string
	HostName         string
	ContainerID      string
	ContainerRuntime string
}

func (n *NodeObject) GetObjectName() string {
	nodeInfo := fmt.Sprintf("%s%s%s%s%s", "node", ObjectNameSplit, n.NodeName, ObjectNameSplit, n.NodeInternalIP)
	//if n.ContainerID != "" {
	//	cInfo := fmt.Sprintf("%s%s%s", n.ContainerRuntime, ContainerSplit, n.ContainerID)
	//	nodeInfo = fmt.Sprintf("%s%s%s", nodeInfo, ObjectNameSplit, cInfo)
	//}

	return nodeInfo
}

func ParseNodeInfo(nodeStr string) (nodeName, nodeIP string, err error) {
	tmpArr := strings.Split(nodeStr, ObjectNameSplit)
	if len(tmpArr) != 3 {
		err = fmt.Errorf("unexpected format of node string: %s", nodeStr)
	} else {
		nodeName, nodeIP = tmpArr[1], tmpArr[2]
	}

	return
}

//func ParseNodeInfo(nodeStr string) (nodeName, nodeIP, cRuntime, cId string, err error) {
//	tmpArr := strings.Split(nodeStr, ObjectNameSplit)
//	if len(tmpArr) == 2 {
//		nodeName, nodeIP = tmpArr[0], tmpArr[1]
//	} else if len(tmpArr) == 3 {
//		cTmp := strings.Split(tmpArr[2], ContainerSplit)
//		if len(cTmp) != 2 {
//			err = fmt.Errorf("unexpected container info: %s", tmpArr[2])
//			return
//		}
//
//		return tmpArr[0], tmpArr[1], cTmp[0], cTmp[1], nil
//	} else {
//		err = fmt.Errorf("unexpected format of node string: %s", nodeStr)
//	}
//
//	return
//}
