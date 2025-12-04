package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	once sync.Once
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

func NewConfiguration(path string) *Configuration {
	config := &Configuration{}
	file, err := os.Open(path)
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil
	}
	if config.MySQL.Port == 0 {
		config.MySQL.Port = 3306
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	return config
}
