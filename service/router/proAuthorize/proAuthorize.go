package proAuthorize

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitProAuthorizeRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiProAuthorize.ProAuthorizeApi

	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("proAuthorize", middleware.LoginInterceptor)
	{
		private.POST("/getAuthorize", api.GetAuthorize)
		private.POST("/getAuthorizeHistoryRecord", api.GetAuthorizeHistoryRecord)
		// private.POST("/getRedeemCodeInfo", api.GetRedeemCodeInfo) // 禁用此接口
		private.POST("/redeemCodeWriteOff", api.RedeemCodeWriteOff)
	}

}
