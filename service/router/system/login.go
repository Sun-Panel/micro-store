package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitLogin(router *gin.RouterGroup) {
	loginApi := api_v1.ApiGroupApp.ApiSystem.LoginApi

	router.POST("/oAuth2CodeBind", middleware.LoginInterceptor, loginApi.OAuth2CodeBind)
	router.POST("/oAuth2UnBind", middleware.LoginInterceptor, loginApi.OAuth2UnBind)
	router.POST("/web/logout", middleware.LoginInterceptor, loginApi.Logout)

	// 半公开模式
	router.POST("/oAuth2CodeLogin", middleware.PublicModeInterceptor, loginApi.OAuth2CodeLogin)

	// 公开模式
	router.POST("/web/login", loginApi.Login)
	router.POST("/ssoLogoutWebhook", loginApi.SSOLogoutWebhook)

}
