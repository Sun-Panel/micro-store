package oauth2

import (
	v1 "sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.ApiOAuth2
	// 登录重定向

	r := router.Group("oAuth2/v1/")

	r.GET("login", api.Login)
}
