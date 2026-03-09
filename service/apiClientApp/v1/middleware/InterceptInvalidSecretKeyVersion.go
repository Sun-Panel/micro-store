package middleware

import (
	"sun-panel/apiClientApp/v1/common/apiReturn"

	"github.com/gin-gonic/gin"
)

// 拦截失效秘钥版本
func InterceptInvalidSecretKeyVersion(c *gin.Context) {
	if c.GetBool("secretKeyStatus") == false {
		apiReturn.ErrorCode(c, 1403, "Secret key is not enabled", nil)
		c.Abort()
	}
}
