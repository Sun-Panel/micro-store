package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitVersionReviewRouter(router *gin.RouterGroup) {
	reviewApi := api_v1.ApiGroupApp.ApiAdmin.VersionReviewApi

	// 需要管理员权限的接口
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("version/getList", reviewApi.GetList)
		r.POST("version/getPendingList", reviewApi.GetPendingList)
		r.POST("version/review", reviewApi.Review)
	}
}
