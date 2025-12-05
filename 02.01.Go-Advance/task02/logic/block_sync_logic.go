package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"task02/models"
	"task02/services"
	"task02/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BlockSyncLogic struct {
	db       *gorm.DB
	eth      *services.EthereumService
	analyzer *utils.TransactionAnalyzer
	mu       sync.Mutex
}

func NewBlockSyncLogic(db *gorm.DB, eth *services.EthereumService, analyzer *utils.TransactionAnalyzer) *BlockSyncLogic {
	return &BlockSyncLogic{
		db:       db,
		eth:      eth,
		analyzer: analyzer,
	}
}

func (bsl *BlockSyncLogic) SyncLatestBlock() (uint64, error) {
	latestBlockNumber, err := bsl.eth.GetLatestBlockNumber(context.Background())
	if err != nil {
		return uint64(0), fmt.Errorf("failed to get latest block number: %v", err)
	}
	return latestBlockNumber, bsl.SyncBlockByNumber(latestBlockNumber)
}

func (bsl *BlockSyncLogic) SyncBlockByNumber(blockNumber uint64) error {
	bsl.mu.Lock()
	defer bsl.mu.Unlock()

	var search_block models.Block
	if err := bsl.db.Where("block_number = ?", blockNumber).First(&search_block).Error; err == nil {
		log.Printf("Block %d already exits", blockNumber)
		return nil
	}

	// 获取区块数据
	block, err := bsl.eth.GetBlockByNumber(blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get block %d: %v", blockNumber, err)
	}

	// 在事务中保存数据
	return bsl.db.Transaction(func(tx *gorm.DB) error {
		// 保证区块信息
		blockModel := &models.Block{
			BlockNumber: block.Number().Uint64(),
			BlockHash:   block.Hash().Hex(),
			ParentHash:  block.ParentHash().Hex(),
			Timestamp:   block.Time(),
			GasLimit:    block.GasLimit(),
			GasUsed:     block.GasUsed(),
			Miner:       strings.ToLower(block.Coinbase().Hex()),
			Difficulty:  block.Difficulty().String(),
			Nonce:       block.Nonce(),
			Size:        uint64(block.Size()),
			ExtraData:   common.Bytes2Hex(block.Extra()),
			MixHash:     block.MixDigest().Hex(),
		}

		// 设置 BaseFeePerGas （EIP-1559）
		if block.BaseFee() != nil {
			baseFee := block.BaseFee().Uint64()
			blockModel.BaseFeePerGas = &baseFee
		}

		if err := tx.Create(blockModel).Error; err != nil {
			return fmt.Errorf("failed to save block: %v", err)
		}

		// 初始化操作记录切片
		accountOperations := make([]*models.AccountOperation, 0)
		contractOperations := make([]*models.ContractOperation, 0)
		contracts := make(map[string]*models.Contract) //使用map避免重复合约

		// 批量保存交易
		transactions := make([]*models.Transaction, 0, len(block.Transactions()))
		for i, ethTX := range block.Transactions() {
			receipt, err := bsl.eth.GetTransactionReceipt(ethTX.Hash())
			if err != nil {
				log.Printf("failed to get receipt for tx %s: %v", ethTX.Hash().Hex(), err)
				continue
			}

			txModel, err := bsl.createTransactionModel(ethTX, receipt, blockModel.ID, block.Number().Uint64(), uint64(i))
			if err != nil {
				log.Printf("Failed to create transaction model: %v", err)
				continue
			}

			transactions = append(transactions, txModel)

			// 收集账户操作信息
			accountOps := bsl.collectAccountOperations(ethTX, receipt, block.Time())
			accountOperations = append(accountOperations, accountOps...)

			// 收集合约相关信息
			contractOps, newContracts := bsl.collectContractInfo(ethTX, receipt, block.Number().Uint64())
			contractOperations = append(contractOperations, contractOps...)

			// 合并新发现的合约
			for addr, contract := range newContracts {
				if _, exits := contracts[addr]; !exits {
					contracts[addr] = contract
				}
			}
		}

		// 批量插入交易
		if len(transactions) > 0 {
			batchSize := 100
			for i := 0; i < len(transactions); i += batchSize {
				end := i + batchSize
				if end > len(transactions) {
					end = len(transactions)
				}
				if err := tx.Create(transactions[i:end]).Error; err != nil {
					return fmt.Errorf("failed to save transactions batch: %v", err)
				}
			}
		}

		// 保存账户操作记录
		if len(accountOperations) > 0 {
			batchSize := 100
			for i := 0; i < len(accountOperations); i += batchSize {
				end := i + batchSize
				if end > len(accountOperations) {
					end = len(accountOperations)
				}
				if err := tx.Create(accountOperations[i:end]).Error; err != nil {
					return fmt.Errorf("failed to save account operations batch: %v", err)
				}
			}
		}

		// 保存合约操作记录
		if len(contractOperations) > 0 {
			batchSize := 100
			for i := 0; i < len(contractOperations); i += batchSize {
				end := i + batchSize
				if end > len(contractOperations) {
					end = len(contractOperations)
				}
				if err := tx.Create(contractOperations[i:end]).Error; err != nil {
					return fmt.Errorf("failed to save contract operations batch: %v", err)
				}
			}
		}

		// 保存新发现合约
		if len(contracts) > 0 {
			contractList := make([]*models.Contract, 0, len(contracts))
			for _, contract := range contracts {
				contractList = append(contractList, contract)
			}
			batchSize := 100
			for i := 0; i < len(contractList); i += batchSize {
				end := i + batchSize
				if end > len(contractList) {
					end = len(contractList)
				}
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "contract_address"}},
					DoNothing: true,
				}).Create(contractList[i:end]).Error; err != nil {
					return fmt.Errorf("failed to save contracts batch: %v", err)
				}
			}
		}

		log.Printf("Synced block %d with %d transcations, %d account operations, %d contract operations and %d contracts",
			blockNumber, len(transactions), len(accountOperations), len(contractOperations), len(contracts))
		return nil
	})
}

func (bsl *BlockSyncLogic) createTransactionModel(
	ethTx *types.Transaction,
	receipt *types.Receipt,
	blockID uint64,
	blockNumber uint64,
	txIndex uint64,
) (*models.Transaction, error) {
	fromAddr := getFromAddress(ethTx)

	// 分析交易类型
	analysis := bsl.analyzer.AnalyzeTransactionType(ethTx, receipt)

	// 创建交易模型
	txModel := &models.Transaction{
		BlockID:          blockID,
		BlockNumber:      blockNumber,
		TransactionHash:  strings.ToLower(ethTx.Hash().Hex()),
		TransactionIndex: txIndex,
		FromAddress:      fromAddr,
		ToAddress:        getToAddress(ethTx),
		Value:            ethTx.Value().String(),
		GasPrice:         ethTx.GasPrice().Uint64(),
		GasLimit:         ethTx.Gas(),
		GasUsed:          receipt.GasUsed,
		Nonce:            ethTx.Nonce(),
		InputData:        common.Bytes2Hex(ethTx.Data()),
		Status:           getTxStatus(receipt),
		TxType:           analysis.BusinessType,
		IsContractDeploy: ethTx.To() == nil,
		ContractAddress:  analysis.ContractAddress,
	}

	// 设置 EIP-1559 相关字段
	if ethTx.Type() == types.DynamicFeeTxType {
		maxPriorityFee := ethTx.GasTipCap().Uint64()
		maxFee := ethTx.GasFeeCap().Uint64()

		txModel.MaxPriorityFeePerGas = &maxPriorityFee
		txModel.MaxFeePerGas = &maxFee

		if receipt.EffectiveGasPrice != nil {
			txModel.EffectiveGasPrice = receipt.EffectiveGasPrice.Uint64()
		}
	}

	// 处理 AccessList
	//if ethTx.Type() == types.AccessListTxType || ethTx.Type() == types.DynamicFeeTxType {
	//	accessList := ethTx.AccessList()
	//	if len(accessList) > 0 {
	//		accessListJson, err := json.Marshal(accessList)
	//		if err != nil {
	//			txModel.AccessList = "{}"
	//		} else {
	//			txModel.AccessList = string(accessListJson)
	//		}
	//	} else {
	//		txModel.AccessList = "{}}"
	//	}
	//}

	return txModel, nil
}

// 收集交易涉及的所有账户操作记录
// 记录发送方的"send"操作
// 记录接收方的"receive"操作
// 如果是合约创建交易，还会记录新合约地址的"contract_created"操作
func (bsl *BlockSyncLogic) collectAccountOperations(
	ethTx *types.Transaction,
	receipt *types.Receipt,
	timestamp uint64,
) []*models.AccountOperation {
	operations := make([]*models.AccountOperation, 0)
	// 获取交易发送方地址
	fromAddr := getFromAddress(ethTx)
	// 获取交易接收方地址
	toAddr := getToAddress(ethTx)
	blockNumber := receipt.BlockNumber.Uint64()
	txHash := strings.ToLower(ethTx.Hash().Hex())
	value := ethTx.Value().String()

	// 记录发送方操作
	if fromAddr != "" {
		operations = append(operations, &models.AccountOperation{
			Address:         fromAddr,
			OperationType:   "send",
			TransactionHash: txHash,
			BlockNumber:     blockNumber,
			FromAddress:     fromAddr,
			ToAddress:       toAddr,
			Value:           value,
			Timestamp:       timestamp,
			CreatedAt:       time.Now(),
		})
	}

	// 记录接收方操作
	if toAddr != "" {
		operations = append(operations, &models.AccountOperation{
			Address:         toAddr,
			OperationType:   "receive",
			TransactionHash: txHash,
			BlockNumber:     blockNumber,
			FromAddress:     fromAddr,
			ToAddress:       toAddr,
			Value:           value,
			Timestamp:       timestamp,
			CreatedAt:       time.Now(),
		})
	}

	// 如果是合约创建交易，记录新合约地址的操作
	if receipt.ContractAddress != (common.Address{}) {
		operations = append(operations, &models.AccountOperation{
			Address:         strings.ToLower(receipt.ContractAddress.Hex()),
			OperationType:   "contract_created",
			TransactionHash: txHash,
			BlockNumber:     blockNumber,
			FromAddress:     fromAddr,
			ToAddress:       toAddr,
			Value:           value,
			Timestamp:       timestamp,
			CreatedAt:       time.Now(),
		})
	}

	return operations
}

// collectContractInfo 收集合约相关信息
// 收集合约创建和交互信息
// 对于合约创建交易:
// 创建新的 Contract 记录
// 添加"created"类型的 ContractOperation 记录
// 对于合约交互交易:
// 添加"interaction"类型的 ContractOperation 记录
// 包含事件日志信息
// 处理事件日志:
// 为每个事件创建"event"类型的 ContractOperation 记录
func (bsl *BlockSyncLogic) collectContractInfo(
	ethTx *types.Transaction,
	receipt *types.Receipt,
	blockNumber uint64,
) ([]*models.ContractOperation, map[string]*models.Contract) {
	contractOperations := make([]*models.ContractOperation, 0)
	contracts := make(map[string]*models.Contract)

	fromAddr := getFromAddress(ethTx)
	toAddr := getToAddress(ethTx)
	txHash := strings.ToLower(ethTx.Hash().Hex())
	inputData := common.Bytes2Hex(ethTx.Data())

	// 检查是否是合约创建交易
	if receipt.ContractAddress != (common.Address{}) {
		contractAddr := strings.ToLower(receipt.ContractAddress.Hex())

		// 添加到新合约的映射中
		contracts[contractAddr] = &models.Contract{
			ContractAddress: contractAddr,
			IsProxy:         false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		// 添加合约创建操作记录
		contractOperations = append(contractOperations, &models.ContractOperation{
			ContractAddress: contractAddr,
			OperationType:   "created",
			TransactionHash: txHash,
			BlockNumber:     blockNumber,
			FromAddress:     fromAddr,
			InputData:       inputData,
			EventData:       "{}",
			CreatedAt:       time.Now(),
		})
	} else if toAddr != "" {
		contractOp := &models.ContractOperation{
			ContractAddress: toAddr,
			OperationType:   "interaction",
			TransactionHash: txHash,
			BlockNumber:     blockNumber,
			FromAddress:     fromAddr,
			InputData:       inputData,
			EventData:       "{}",
			CreatedAt:       time.Now(),
		}

		// 如果有事件日志，添加事件数据
		if len(receipt.Logs) > 0 {
			eventData := fmt.Sprintf(`{"log_count": %d}`, len(receipt.Logs))
			contractOp.EventData = eventData

			if len(receipt.Logs) > 0 && len(receipt.Logs[0].Topics) > 0 {
				contractOp.EventName = fmt.Sprintf("event_%s", receipt.Logs[0].Topics[0].Hex()[:10])
			}
		}
		contractOperations = append(contractOperations, contractOp)
	}

	// 处理事件日志
	for _, logEntry := range receipt.Logs {
		if len(logEntry.Topics) > 0 {
			eventTopic := logEntry.Topics[0].Hex()
			eventName := fmt.Sprintf("event_%s", eventTopic[:10])

			// 构架结构化的事件数据
			eventStruct := map[string]interface{}{
				"topics": make([]string, len(logEntry.Topics)),
				"data":   common.Bytes2Hex(logEntry.Data),
			}

			// 转换 topics
			for i, topic := range logEntry.Topics {
				eventStruct["topics"].([]string)[i] = topic.Hex()
			}

			// 序列化为 JSON 字符串
			eventDataBytes, err := json.Marshal(eventStruct)
			var eventData string
			if err != nil {
				eventData = "{}"
			} else {
				eventData = string(eventDataBytes)
			}

			eventOp := &models.ContractOperation{
				ContractAddress: strings.ToLower(logEntry.Address.Hex()),
				OperationType:   "event",
				TransactionHash: txHash,
				BlockNumber:     blockNumber,
				FromAddress:     fromAddr,
				EventName:       eventName,
				EventData:       eventData,
				CreatedAt:       time.Now(),
			}
			contractOperations = append(contractOperations, eventOp)
		}
	}

	return contractOperations, contracts
}

func getFromAddress(tx *types.Transaction) string {
	var fromAddress string

	// 安全地获取发送方地址
	var sender common.Address
	var err error

	chainID := tx.ChainId()
	if chainID != nil && chainID.Sign() > 0 {
		signer := types.LatestSignerForChainID(chainID)
		sender, err = types.Sender(signer, tx)
	} else {
		chainConfig := params.MainnetChainConfig
		signer := types.LatestSigner(chainConfig)
		sender, err = types.Sender(signer, tx)
	}

	if err != nil {
		fromAddress = "0x0000000000000000000000000000000000000000"
	} else {
		fromAddress = strings.ToLower(sender.Hex())
	}
	return fromAddress
}

func getToAddress(tx *types.Transaction) string {
	to := tx.To()
	if to == nil {
		return ""
	}
	return strings.ToLower(to.Hex())
}

func getTxStatus(receipt *types.Receipt) uint64 {
	if receipt.Status == 1 {
		return 1
	}
	return 0
}
