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
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var DefaultRunOptIns *Config

type RunMode string

const (
	RunModeKubeConfig     RunMode = "KubeConfig"
	RunModeServiceAccount RunMode = "ServiceAccount"
)

func (r RunMode) Int() int {
	switch r {
	case RunModeKubeConfig:
		return -1
	case RunModeServiceAccount:
		return 0
	default:
		return 0
	}
}

type Config struct {
	SecretKey string `yaml:"secretkey"`
	DB        struct {
		Name    string `yaml:"name"`
		User    string `yaml:"user"`
		Passwd  string `yaml:"passwd"`
		Url     string `yaml:"url"`
		MaxIdle int    `yaml:"maxidle"`
		MaxConn int    `yaml:"maxconn"`
		Debug   bool   `yaml:"debug"`
	} `yaml:"db"`
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	} `yaml:"log"`
	RunMode RunMode `yaml:"runmode"`
}

func InitConfigWithFilePath(filePath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	viper.AddConfigPath(filepath.Join(home, "conf"))
	viper.AddConfigPath(filepath.Join(getCurrentPath(), "conf"))
	if len(filePath) > 0 {
		viper.AddConfigPath(filePath)
	}
	viper.SetConfigName("app")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	DefaultRunOptIns = &Config{}
	return viper.Unmarshal(DefaultRunOptIns)
}

func InitConfig() {
	DefaultRunOptIns = &Config{}
	if err := viper.Unmarshal(DefaultRunOptIns); err != nil {
		log.Panic(err)
	}
}

func getCurrentPath() string {
	if ex, err := os.Executable(); err == nil {
		return filepath.Dir(ex)
	}
	return "./"
}
