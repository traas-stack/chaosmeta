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

package version

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
)

func PrintVersion(ctx context.Context) {
	logger := log.GetLogger(ctx)
	reBytes, _ := json.Marshal(GetVersion())
	if log.Path == "" {
		fmt.Println(string(reBytes))
	} else {
		logger.Info(string(reBytes))
	}
}

func GetVersion() *Info {
	return &Info{
		Version:   "@VERSION@",
		BuildDate: "@DATE@",
	}
}
