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

// ns/pod

type PodObject struct {
	Namespace  string
	PodName    string
	PodUID     string
	PodIP      string
	NodeName   string
	NodeIP     string
	Containers []ContainerObject
}

type ContainerObject struct {
	ContainerRuntime string
	ContainerId      string
	ContainerName    string
}

func (c *ContainerObject) GetObjectName() string {
	return fmt.Sprintf("%s%s%s%s%s", c.ContainerRuntime, ObjectNameSplit, c.ContainerId, ObjectNameSplit, c.ContainerName)
}

func (p *PodObject) GetContainersObjectName() []string {
	podInfo := fmt.Sprintf("%s%s%s%s%s", "pod", ObjectNameSplit, p.Namespace, ObjectNameSplit, p.PodName)
	if len(p.Containers) == 0 {
		return []string{podInfo}
	} else {
		objects := make([]string, 0)
		for _, container := range p.Containers {
			objects = append(objects, container.GetObjectName())
		}
		return objects
	}
}

func (p *PodObject) GetObjectName() string {
	podInfo := fmt.Sprintf("%s%s%s%s%s", "pod", ObjectNameSplit, p.Namespace, ObjectNameSplit, p.PodName)
	if len(p.Containers) != 0 {
		for _, container := range p.Containers {
			podInfo = fmt.Sprintf("%s%s%s", podInfo, ObjectNameSplit, container)
		}
	}

	return podInfo
}

func ParsePodInfo(podStr string) (ns, podName, containerName string, err error) {
	tmpArr := strings.Split(podStr, ObjectNameSplit)
	if len(tmpArr) == 4 {
		ns, podName, containerName = tmpArr[1], tmpArr[2], tmpArr[3]
	} else if len(tmpArr) == 3 {
		ns, podName = tmpArr[1], tmpArr[2]
	} else {
		err = fmt.Errorf("unexpected format of pod string: %s", podStr)
	}

	return
}

//func ParsePodInfo(podStr string) (ns, podName string, err error) {
//	tmpArr := strings.Split(podStr, ObjectNameSplit)
//	if len(tmpArr) == 3 {
//		ns, podName = tmpArr[1], tmpArr[2]
//	} else {
//		err = fmt.Errorf("unexpected format of pod string: %s", podStr)
//	}
//
//	return
//}
