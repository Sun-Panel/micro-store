package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMicroAppCategoryRouter(router *gin.RouterGroup) {
	categoryApi := api_v1.ApiGroupApp.ApiAdmin.MicroAppCategoryApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("microAppCategory/getList", categoryApi.GetList)
		r.POST("microAppCategory/getInfo", categoryApi.GetInfo)
		r.POST("microAppCategory/create", categoryApi.Create)
		r.POST("microAppCategory/update", categoryApi.Update)
		r.POST("microAppCategory/deletes", categoryApi.Deletes)
		r.POST("microAppCategory/updateStatus", categoryApi.UpdateStatus)
		r.POST("microAppCategory/getEnabledList", categoryApi.GetEnabledList)
	}
}
