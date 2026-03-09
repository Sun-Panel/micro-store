package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitEmail(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.EmailApi

	r := router.Group("email", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("send", api.Send)
		// r.POST("sendByTemplateId", api.SendByTemplateId)
	}

}
