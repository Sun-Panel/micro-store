package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitProAuthorizeRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.ProAuthorizeApi

	r := router.Group("proAuthorize", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("updateUserExpiredTimeByDay", api.UpdateUserExpiredTimeByDay)
		r.POST("getUserProAuthorizeList", api.GetUserProAuthorizeList)
		r.POST("getAuthorizeHistoryRecordByUserId", api.GetAuthorizeHistoryRecordByUserId)
	}

}
