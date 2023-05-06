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

type DeploymentObject struct {
	Namespace      string
	DeploymentName string
}

func (d *DeploymentObject) GetObjectName() string {
	return fmt.Sprintf("%s%s%s%s%s", "deployment", ObjectNameSplit, d.Namespace, ObjectNameSplit, d.DeploymentName)
}

func ParseDeploymentInfo(str string) (nodeName, nodeIP string, err error) {
	tmpArr := strings.Split(str, ObjectNameSplit)
	if len(tmpArr) == 3 {
		nodeName, nodeIP = tmpArr[1], tmpArr[2]
	} else {
		err = fmt.Errorf("unexpected format of deployment string: %s", str)
	}

	return
}
