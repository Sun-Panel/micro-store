package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMicroAppRouter(router *gin.RouterGroup) {
	microAppApi := api_v1.ApiGroupApp.ApiAdmin.MicroAppApi

	// 需要管理员权限的接口
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("microApp/getList", microAppApi.GetList)
		r.POST("microApp/getInfo", microAppApi.GetInfo)
		r.POST("microApp/create", microAppApi.Create)
		r.POST("microApp/update", microAppApi.Update)
		r.POST("microApp/deletes", microAppApi.Deletes)
		r.POST("microApp/updateStatus", microAppApi.UpdateStatus)
	}

	// 需要登录权限的接口（开发者使用）
	rNoAdmin := router.Group("", middleware.LoginInterceptor)
	{
		rNoAdmin.POST("microApp/updateLang", microAppApi.UpdateLang)
	}
}
