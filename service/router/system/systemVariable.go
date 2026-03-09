package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitSystemVariableRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.SystemVariableApi

	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("systemVariableApi/getMultiple", api.GetMultiple)
		// r.POST("systemVariableApi/set", api.Set)
		r.POST("systemVariableApi/get", api.Get)

	}

}
