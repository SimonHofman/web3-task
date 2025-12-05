package handlers

import (
	"log"
	"task02/logic"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	al *logic.AddressLogic
}

func NewAddressHandler(al *logic.AddressLogic) *AddressHandler {
	return &AddressHandler{
		al: al,
	}
}

func (ah *AddressHandler) GetAddressInfo(c *gin.Context) {
	address, exist := c.GetQuery("address")
	if !exist {
		c.JSON(400, gin.H{
			"message": "address is required",
		})
		return
	}
	info, isContract, err := ah.al.GetAddressInfo(address)
	if err != nil {
		log.Fatal("fail to get address info")
	}
	if isContract {
		c.JSON(200, gin.H{
			"message": "contract",
			"info":    info,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "account",
			"info":    info,
		})
	}
}
