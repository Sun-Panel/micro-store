package middleware

import (
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseLogger 生成responseId并记录错误日志
func ResponseLogger(c *gin.Context) {
	// 生成唯一的responseId
	responseId := time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	c.Set("responseId", responseId)

	c.Next()

	// 只记录错误请求：HTTP状态码>=400 或 业务code!=0
	status := c.Writer.Status()
	hasError := c.GetBool("hasError")

	if status >= 400 || hasError {
		bizCode, _ := c.Get("bizCode")
		global.Logger.Warnf("ResponseId: %s | %s %s | HTTP: %d | BizCode: %v",
			responseId, c.Request.Method, c.Request.URL.Path, status, bizCode)
	}
}
