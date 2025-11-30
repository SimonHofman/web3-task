package client

import (
	"log"
	"sync"
	"task02/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	once   sync.Once
	client *ethclient.Client
)

func GetClient() {
	once.Do(func() {
		config := config.GetConfig().EthClient
		var err error
		client, err = ethclient.Dial(config.Url)
		if err != nil {
			log.Fatal("连接以太坊客户端失败", err)
		}
	})
}
