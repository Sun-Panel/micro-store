package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitCustomCodeRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.CustomCodeApi

	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("customCode/getList", api.GetList)
		rPublic.POST("customCode/get", api.Get)
	}
}
