package logic

import (
	"fmt"
	"task02/models"

	"task02/services"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type AddressLogic struct {
	db  *gorm.DB
	eth *services.EthereumService
}

func NewAddressLogic(db *gorm.DB, eth *services.EthereumService) *AddressLogic {
	return &AddressLogic{
		db:  db,
		eth: eth,
	}
}

func (al *AddressLogic) GetAddressInfo(address string) (interface{}, bool, error) {
	isContract, err := al.isContractAddress(address)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check if address is contract: %v", err)
	}

	if isContract {
		var contractOperations []models.ContractOperation
		err := al.db.Where("contract_address = ?", address).Find(&contractOperations).Error
		if err != nil {
			return nil, isContract, fmt.Errorf("failed to get contract operations: %v", err)
		}
		return contractOperations, isContract, nil
	} else {
		var accountOperations []models.AccountOperation
		err := al.db.Where("from_address = ? OR to_address = ?", address, address).Find(&accountOperations).Error
		if err != nil {
			return nil, isContract, fmt.Errorf("failed to get account operations: %v", err)
		}
		return accountOperations, isContract, nil
	}
}

func (al *AddressLogic) isContractAddress(address string) (bool, error) {
	code, err := al.eth.GetCode(common.HexToAddress(address))
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}
