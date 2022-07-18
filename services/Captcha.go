package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, content, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "failed to generate captcha",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"content":   content,
	})
}

func CheckCaptcha(c *gin.Context) {
	id := c.DefaultQuery("captchaId", "")
	captchaVal := c.DefaultQuery("captchaVal", "")

	res := store.Verify(id, captchaVal, true)
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
	if !res {
		c.Abort()
	}
	c.Next()
}
