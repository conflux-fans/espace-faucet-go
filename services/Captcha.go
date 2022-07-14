package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	// height 高度 png 像素高度
	// width  宽度 png 像素高度
	// length 验证码默认位数
	// maxSkew 单个数字的最大绝对倾斜因子
	// dotCount 背景圆圈的数量
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, content, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"content":   content,
	})
}

func CheckCapcha(c *gin.Context) {
	id := c.DefaultQuery("captchaId", "")
	captchaVal := c.DefaultQuery("captchaVal", "")
	// id 验证码id
	// answer 需要校验的内容
	// clear 校验完是否清除
	res := store.Verify(id, captchaVal, true)
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
