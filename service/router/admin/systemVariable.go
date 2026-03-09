package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitSystemVariableRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.SystemVariableApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("systemVariable/getList", api.GetList)
		r.POST("systemVariable/set", api.Set)
		r.POST("systemVariable/edit", api.Edit)
		r.POST("systemVariable/delete", api.Delete)
		r.POST("systemVariable/clearCache", api.ClearCache)
		r.POST("systemVariable/getByCache", api.GetByCache)
	}

}
