package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var DefaultRunOptIns *Config

type Config struct {
	SecretKey string `yaml:"secretkey"`
	DB        struct {
		Name    string `yaml:"name"`
		User    string `yaml:"user"`
		Passwd  string `yaml:"passwd"`
		Url     string `yaml:"url"`
		MaxIdle int    `yaml:"maxidle"`
		MaxConn int    `yaml:"maxconn"`
	} `yaml:"db"`
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	} `yaml:"log"`
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
