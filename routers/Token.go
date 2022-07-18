package routers

import (
	"errors"
	"github.com/conflux-fans/espace-faucet-go/faucetErrors"
	"github.com/conflux-fans/espace-faucet-go/models"
	"github.com/conflux-fans/espace-faucet-go/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

var lastCFXClaimCache = make(map[string]int64)
var lastERC20ClaimCache = make(map[string]map[string]int64)

func sendCFX(c *gin.Context) {
	addr := c.Query("address")
	if addr == "" {
		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
		return
	}

	value, ok := lastCFXClaimCache[addr]
	if ok {
		res := time.Now().Unix() - value
		if time.Duration(res) < 3600000 {
			renderBaseError(c, faucetErrors.TIME_ERROR)
			return
		}
	}
	lastCFXClaimCache[addr] = time.Now().Unix()

	res, err := services.SendCFX(addr)
	renderResponse(c, res, err)
}

func sendERC20(c *gin.Context) {
	var ERC20Data *models.ERC20
	if err := c.ShouldBind(&ERC20Data); err != nil {
		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
		return
	}
	if err := checkERC20(ERC20Data.Name); err != nil {
		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
		return
	}

	value, ok := lastERC20ClaimCache[ERC20Data.Address][ERC20Data.Name]
	if ok {
		res := time.Now().Unix() - value
		if time.Duration(res) < 3600000 {
			renderBaseError(c, faucetErrors.TIME_ERROR)
			return
		}
		lastERC20ClaimCache[ERC20Data.Address][ERC20Data.Name] = time.Now().Unix()
	} else {
		subMap := make(map[string]int64)
		subMap[ERC20Data.Name] = time.Now().Unix()
		lastERC20ClaimCache[ERC20Data.Address] = subMap
	}
	res, err := services.SendERC20(*ERC20Data)
	renderResponse(c, res, err)

}

func checkERC20(symbol string) error {
	erc20Data := viper.GetStringMap("erc20")
	for i := range erc20Data {
		if symbol == i {
			return nil
		}
	}
	return errors.New("Unsupported symbol")
}
