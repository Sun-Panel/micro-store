package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitRedeemCodeRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.RedeemCodeApi

	r := router.Group("redeemCode", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("getList", api.GetList)
		r.POST("create", api.Create)
		r.POST("setInvalid", api.SetInvalid)
	}

}
