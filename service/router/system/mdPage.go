package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMdPageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.MdPageApi

	rPubilc := router.Group("", middleware.PublicModeInterceptor)
	{
		rPubilc.POST("mdPage/get", api.Get)
	}

}
