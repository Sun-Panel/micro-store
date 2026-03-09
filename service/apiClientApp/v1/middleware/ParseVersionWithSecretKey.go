package middleware

import (
	"sun-panel/apiClientApp/v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

// 解析版本和秘钥
func ParseVersionWithSecretKey(c *gin.Context) {
	version := c.GetHeader("SunPanel-Version")
	// global.Logger.Info("SunPanel-Version", version)
	if version == "" {
		// 默认为v1系列
		// vs, err := biz.VersionSecret.GetVersionSecret("v1")
		// if err != nil {
		// 	// 数据库没有存入v1系列的独立秘钥,使用旧版本秘钥
		// 	c.Set("secretKey", "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN")
		// 	c.Set("secretKeyStatus", true)
		// 	global.Logger.Infoln("555")
		// } else {
		// 	c.Set("secretKey", vs.SecretKey)
		// 	c.Set("secretKeyStatus", vs.Status)
		// }

		// 数据库没有存入v1系列的独立秘钥,使用旧版本秘钥 如果更新需要把上面注释的代码打开
		c.Set("secretKey", "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN")
		c.Set("secretKeyStatus", true)
	} else {
		// v2以上
		// 使用新版本独立秘钥
		vs, err := biz.VersionSecret.GetVersionSecret(version)
		if err != nil {
			c.Set("secretKey", "")
			c.Set("secretKeyStatus", false)
			global.Logger.Errorw("Failed to obtain version key", "error", err)
			apiReturn.ErrorCode(c, 500, "Failed to obtain version key", err)
			c.Abort()
		} else {
			c.Set("secretKey", vs.SecretKey)
			c.Set("secretKeyStatus", vs.Status)
		}
	}

}
