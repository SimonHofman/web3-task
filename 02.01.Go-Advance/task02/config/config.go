package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	config *Configuration
	once   sync.Once
)

type Configuration struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`

	Auth struct {
		Jwt_secret  string `yaml:"jwt_secret"`
		Toke_expiry int    `yaml:"token_expiry"`
	} `yaml:"auth"`

	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	} `yaml:"log"`

	EthClient struct {
		RpcUrl    string   `yaml:"rpc_url"`
		Accounts  []string `yaml:"accounts"`
		Contracts []string `yaml:"contracts"`
		Timeout   int      `yaml:"timeout"`
	} `yaml:"eth_client"`
}

func InitConfig(path string) {
	once.Do(func() {
		config = &Configuration{}
		if err := loadconfig(path); err != nil {
			log.Fatal("配置文件加载失败")
		}
	})
}

func GetConfig() *Configuration {
	return config
}

func loadconfig(path string) error {
	file, err := os.Open(path)
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	return nil
}
