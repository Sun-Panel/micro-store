package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userApi := api_v1.ApiGroupApp.ApiAdmin.UsersApi

	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("users/create", userApi.Create)
		r.POST("users/update", userApi.Update)
		r.POST("users/getList", userApi.GetList)
		r.POST("users/deletes", userApi.Deletes)
	}
}
