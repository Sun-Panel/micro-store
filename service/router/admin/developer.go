package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitDeveloperRouter(router *gin.RouterGroup) {
	developerApi := api_v1.ApiGroupApp.ApiAdmin.DeveloperApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("developer/getList", developerApi.GetList)
		r.POST("developer/getInfo", developerApi.GetInfo)
		r.POST("developer/update", developerApi.Update)
		r.POST("developer/deletes", developerApi.Deletes)
		r.POST("developer/updateStatus", developerApi.UpdateStatus)
	}
}
