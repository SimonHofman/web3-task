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
	if err := bsl.db.Where("number = ?", blockNumber).First(&search_block).Error; err != nil {
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
	})

	return nil
}
