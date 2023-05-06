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

const (
	ObjectNameSplit      = "/"
	containerIDDelimiter = "://"
	defaultRuntime       = "docker"
)

type AtomicObject interface {
	GetObjectName() string
	//GetName() string
	//GetUID() string
	//GetMessage() string
	//GetStatus() string
}

func ParseContainerID(cID string) (r, id string, err error) {
	if cID == "" {
		return "", "", fmt.Errorf("container id is empty")
	}

	tmpArr := strings.Split(cID, containerIDDelimiter)
	if len(tmpArr) == 2 {
		r, id = tmpArr[0], tmpArr[1]
	} else if len(tmpArr) == 1 {
		r, id = defaultRuntime, tmpArr[0]
	} else {
		err = fmt.Errorf("container info format error: %s", cID)
	}

	return
}
