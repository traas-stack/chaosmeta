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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var globalConfig *MainConfig

func init() {
	_, filename, _, _ := runtime.Caller(0)
	absPath, _ := filepath.Abs(filepath.Dir(filename))
	if err := LoadConfig(fmt.Sprintf("%s/config/chaosmeta-measure.json", filepath.Dir(filepath.Dir(absPath)))); err != nil {
		panic(any(fmt.Sprintf("load config error: %s", err.Error())))
	}
}

func GetGlobalConfig() *MainConfig {
	return globalConfig
}

func LoadConfig(path string) error {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load config file error: %s", err.Error())
	}
	var mainConfig MainConfig
	if err := json.Unmarshal(configBytes, &mainConfig); err != nil {
		return fmt.Errorf("convert config file error: %s", err.Error())
	}

	globalConfig = &mainConfig
	fmt.Println(string(configBytes))
	return nil
}

type MainConfig struct {
	Monitor   MonitorConfig `json:"monitor"`
	TaskLimit int           `json:"tasklimit"`
}

type MonitorConfig struct {
	Url    string `json:"url"`
	Engine string `json:"engine"`
}
