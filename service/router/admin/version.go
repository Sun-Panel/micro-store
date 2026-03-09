package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitVersionRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.VersionApi
	secretApi := api_v1.ApiGroupApp.ApiAdmin.VersionSecretApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("version/edit", api.Edit)
		r.POST("version/setActive", api.SetActive)
		r.POST("version/getList", api.GetList)
		r.POST("version/deletes", api.Deletes)
	}
	// 版本密钥相关接口
	{
		r.POST("versionSecret/edit", secretApi.Edit)
		r.POST("versionSecret/getByVersion", secretApi.GetVersionSecretInfoByVersion)
	}
}
