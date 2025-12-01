package client

import (
	"log"
	"sync"
	"task02/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	once sync.Once
)

func GetClient() *ethclient.Client {
	var client *ethclient.Client
	once.Do(func() {
		config := config.GetConfig().EthClient
		var err error
		client, err = ethclient.Dial(config.Url)
		if err != nil {
			log.Fatal("连接以太坊客户端失败", err)
		}
	})
	return client
}
