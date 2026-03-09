package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitSystemSettingRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.SystemSettingApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("systemSetting/getEmail", api.GetEmail)
		r.POST("systemSetting/setEmail", api.SetEmail)
		r.POST("systemSetting/getApplicationSetting", api.GetApplicationSetting)
		r.POST("systemSetting/setApplicationSetting", api.SetApplicationSetting)
	}

}
