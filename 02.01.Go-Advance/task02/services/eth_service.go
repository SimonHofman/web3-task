package services

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"task02/config"
)

type EthereumService struct {
	client *ethclient.Client
	rpcURL string
}

func NewEthereumService() *EthereumService {
	cfg := config.GetConfig().EthClient
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, cfg.RpcUrl)
	if err != nil {
		log.Fatal("RPC连接失败:", err)
	}

	return &EthereumService{
		client: client,
		rpcURL: cfg.RpcUrl,
	}
}

func (es *EthereumService) GetBlockByNumber(blockNumber uint64) (*types.Block, error) {
	cfg := config.GetConfig().EthClient
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	return es.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

func (es *EthereumService) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	header, err := es.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (es *EthereumService) GetBlockByHash(blockHash common.Hash) (*types.Block, error) {
	cfg := config.GetConfig().EthClient
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()
	return es.client.BlockByHash(ctx, blockHash)
}

func (es *EthereumService) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return es.client.TransactionReceipt(context.Background(), txHash)
}

func (es *EthereumService) GetCode(address common.Address) ([]byte, error) {
	return es.client.CodeAt(context.Background(), address, nil)
}

func (es *EthereumService) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	return es.client.CallContract(context.Background(), msg, nil)
}
