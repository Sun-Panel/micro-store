package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMicroAppCategoryRouter(router *gin.RouterGroup) {
	categoryApi := api_v1.ApiGroupApp.ApiAdmin.MicroAppCategoryApi

	// 需要管理员权限的接口
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("microAppCategory/getList", categoryApi.GetList)
		r.POST("microAppCategory/getInfo", categoryApi.GetInfo)
		r.POST("microAppCategory/create", categoryApi.Create)
		r.POST("microAppCategory/update", categoryApi.Update)
		r.POST("microAppCategory/deletes", categoryApi.Deletes)
		r.POST("microAppCategory/updateStatus", categoryApi.UpdateStatus)
	}

	// 公开接口（微应用详情页使用，无需登录）
	router.POST("microAppCategory/getEnabledList", categoryApi.GetEnabledList)
}
