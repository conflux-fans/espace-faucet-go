package routers

import (
	"github.com/conflux-fans/espace-faucet-go/faucetErrors"
	"github.com/conflux-fans/espace-faucet-go/models"
	"github.com/conflux-fans/espace-faucet-go/services"
	"github.com/gin-gonic/gin"
	"time"
)

var lastCFXClaimCache = make(map[string]int64)
var lastERC20claimCache = make(map[string]int64)



func sendCFX(c *gin.Context)  {
	addr := c.Query("address")
	if addr == "" {
		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
		return
	}

	value,ok := lastCFXClaimCache[addr]
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

func sendERC20(c *gin.Context)  {
	var ERC20Data *models.ERC20
	if err := c.ShouldBind(&ERC20Data); err != nil {
		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
		return
	}
	value,ok := lastERC20claimCache[ERC20Data.Address]
	if ok {
		res := time.Now().Unix() - value
		if time.Duration(res) < 3600000 {
			renderBaseError(c, faucetErrors.TIME_ERROR)
			return
		}
	}
	lastERC20claimCache[ERC20Data.Address] = time.Now().Unix()
	res, err := services.SendERC20(*ERC20Data)
	renderResponse(c, res, err)

}

func deployERC20(c *gin.Context) {
	res, err := services.DeployERC20()
	renderResponse(c, res, err)
}

//func queryERC20(c *gin.Context) {
//	err := services.QueryERC20()
//	if err != nil {
//		renderBaseError(c, faucetErrors.INVALID_REQUEST_ERROR)
//		return
//	}
//
//}


