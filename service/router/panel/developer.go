package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitDeveloperRouter(router *gin.RouterGroup) {
	developerApi := api_v1.ApiGroupApp.ApiPanel.DeveloperApi

	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("developer/register", developerApi.Register)
		r.POST("developer/getInfo", developerApi.GetInfo)
		r.POST("developer/update", developerApi.Update)
		r.POST("developer/checkIsDeveloper", developerApi.CheckIsDeveloper)
	}
}
