package middleware

import (
	"strings"
	"sun-panel/api/apiClientApp/common/apiReturn"
	"sun-panel/api/apiClientApp/common/base"
	"sun-panel/global"

	sunApi "cnb.cool/hslr-s/go-pkg/sun-api"
	"github.com/gin-gonic/gin"
)

// 解析数据 (2026-06-29 新增 / v3 API新增，最低支持客户端版本 2.0.0-dev-14)
func ParseData(c *gin.Context) {
	token, err := sunApi.ExtractToken(c.Request)
	if err != nil {
		apiReturn.ErrorUnauthorized(c)
		global.Logger.Warnln("token error:", err)
		c.Abort()
		return
	}

	// 从header中获取客户端的版本号 - 旧版
	// clientVersion := c.GetHeader("Client-App-Version")
	// if clientVersion == "" {
	// 	apiReturn.ErrorUnauthorized(c)
	// 	global.Logger.Warnln("missing Client-App-Version header")
	// 	c.Abort()
	// 	return
	// }

	// 根据版本号获取对应的密钥
	// secretKey, err := biz.SunPanelAuth.GetVersionSecret(clientVersion)
	// if err != nil || secretKey == "" {
	// 	apiReturn.ErrorUnauthorized(c)
	// 	global.Logger.Warnln("failed to get secret key for version:", clientVersion, "error:", err)
	// 	c.Abort()
	// 	return
	// }

	// 获取密钥 - 新版开始支持
	secretKey := global.Config.GetValueString("sun_panel_auth", "client_api_jwt_secret")

	// 创建claims对象用于接收解析后的payload
	claims := base.Claims{}

	// 使用获取到的密钥验证JWT
	if err := sunApi.ParseTokenHS256(token, secretKey, &claims); err != nil {
		// 判断是否为token过期
		if strings.Contains(err.Error(), "token is expired") {
			apiReturn.ErrorTokenExpires(c)
		} else {
			apiReturn.ErrorUnauthorized(c)
		}
		global.Logger.Warnln("invalid token:", err)
		c.Abort()
		return
	}

	// 将解析后的claims保存到context中，供后续使用
	c.Set("jwtClaims", claims)
	c.Set("username", claims.Username)
	c.Set("SunPanelClientInfo", claims.Client)
	c.Set("proStatus", claims.ProStatus)
	c.Next()

}
