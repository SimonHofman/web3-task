package logic

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"task02/models"
	"task02/services"
	"task02/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"gorm.io/gorm"
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

func (bsl *BlockSyncLogic) BlockSyncByNumber(blockNumber uint64) error {
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

		log.Printf("Synced block %d with %d transcations", blockNumber, len(transactions))
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
	var fromAddr string

	// 安全地获取发送方地址
	var sender common.Address
	var err error

	chainID := ethTx.ChainId()
	if chainID != nil && chainID.Sign() > 0 {
		signer := types.LatestSignerForChainID(chainID)
		sender, err = types.Sender(signer, ethTx)
	} else {
		chainConfig := params.MainnetChainConfig
		signer := types.LatestSigner(chainConfig)
		sender, err = types.Sender(signer, ethTx)
	}

	if err != nil {
		fromAddr = "0x0000000000000000000000000000000000000000"
	} else {
		fromAddr = strings.ToLower(sender.Hex())
	}

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
