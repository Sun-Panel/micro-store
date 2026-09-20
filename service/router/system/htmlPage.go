package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitHtmlPageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.HtmlPageApi

	rPubilc := router.Group("", middleware.PublicModeInterceptor)
	{
		rPubilc.POST("htmlPage/get", api.Get)
	}

}
