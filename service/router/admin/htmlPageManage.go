package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitHtmlPageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.HtmlPageApi

	r := router.Group("htmlPageManage", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/getList", api.GetList)
		r.POST("/getInfo", api.GetInfo)
		r.POST("/edit", api.Edit)
		r.POST("/delete", api.Delete)
	}

}
