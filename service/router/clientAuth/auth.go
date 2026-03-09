package clientAuth

import (
	v1 "sun-panel/apiClientApp/v1"
	"sun-panel/apiClientApp/v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuth(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.Api
	r := router.Group("", middleware.ParseVersionWithSecretKey)

	// 拦截失效秘钥版本
	interceptInvalidSecretKeyVersionRouter := r.Group("", middleware.InterceptInvalidSecretKeyVersion)
	{
		r.POST("register", api.Register)
		r.POST("checkVersion", api.CheckVersion)
		interceptInvalidSecretKeyVersionRouter.POST("refreshInfo", api.RefreshInfo)
		interceptInvalidSecretKeyVersionRouter.POST("login", api.Login)
		interceptInvalidSecretKeyVersionRouter.POST("autoLogin", api.AutoLogin)
	}

	// 中间件（验证token的中间件）
	{
		r.POST("renewTempAuth", api.RenewTempAuth)
		r.POST("ping", api.Ping)

	}
}
