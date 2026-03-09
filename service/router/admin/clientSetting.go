package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitClientSetting(router *gin.RouterGroup) {
	r := router.Group("clientSetting")
	initClientBlacklistIP(r)
	initCilentAPIAuthCheckSetting(r)
}

func initClientBlacklistIP(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.ClientBlackListIPApi

	r := router.Group("blacklistIP", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("getList", api.GetList)
		r.POST("deletes", api.Deletes)
		r.POST("set", api.Set)
	}
}

func initCilentAPIAuthCheckSetting(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.ClientCreateOnlineCacheApi

	r := router.Group("createOnlineCache", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("getAll", api.GetAll)
		r.POST("setAll", api.SetAll)
	}
}
