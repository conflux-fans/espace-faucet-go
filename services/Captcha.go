package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore
var driver = base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)

func GetCaptcha(c *gin.Context) {
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
	if (id == "" || captchaVal == "") || !res {
		c.AbortWithStatusJSON(402, gin.H{
			"code":    402,
			"message": "wrong captcha value",
		})
		return
	}
	c.Next()
}
