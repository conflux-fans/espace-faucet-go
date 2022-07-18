package routers

import (
	"github.com/conflux-fans/espace-faucet-go/faucetErrors"
	"github.com/conflux-fans/espace-faucet-go/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", indexEndpoint)
	router.GET("/captcha/", services.GetCaptcha)

	apiV1 := router.Group("/v1")
	apiV1.Use(services.CheckCapcha)
	{
		apiV1.POST("/CFX", sendCFX)
		apiV1.POST("/ERC20", sendERC20)
	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(200, dataResponse("Ethereum-space-faucet"))
}

func dataResponse(data interface{}) gin.H {
	return gin.H{
		"code": 0,
		"data": data,
	}
}

func errorResponse(code int, err error) gin.H {
	return gin.H{
		"code":    code,
		"message": err.Error(),
	}
}

func buildResponse(data interface{}, err error) (int, gin.H) {
	if err != nil {
		return 200, errorResponse(500, err)
	}
	return 200, dataResponse(data)
}

func renderResponse(c *gin.Context, data interface{}, err error) {
	code, response := buildResponse(data, err)
	c.JSON(code, response)
}

func renderError(c *gin.Context, code int, err error) {
	c.JSON(code, errorResponse(code, err))
}

func renderBaseError(c *gin.Context, err faucetErrors.BaseError) {
	renderError(c, err.Code, err)
}
