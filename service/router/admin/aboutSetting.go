package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitAboutSetting(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.AboutSettingApi

	r := router.Group("aboutSetting", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("save", api.Save)
	}

}
