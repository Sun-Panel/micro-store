package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitMdPageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.MdPageApi

	r := router.Group("mdPage", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/getList", api.GetList)
		r.POST("/edit", api.Edit)
		r.POST("/delete", api.Delete)
	}

}
