package panel

import (
	"github.com/gin-gonic/gin"
)

func InitDeveloperVersionRouter(router *gin.RouterGroup) {
	// versionApi := api_v1.ApiGroupApp.ApiAdmin.DeveloperVersionApi
	// uploadApi := api_v1.ApiGroupApp.ApiAdmin.MicroAppVersionUploadApi

	// // 需要登录权限的接口（开发者使用）
	// r := router.Group("", middleware.LoginInterceptor)
	// {
	// 	r.POST("developer/version/getList", versionApi.GetList)
	// 	r.POST("developer/version/getInfo", versionApi.GetInfo)
	// 	r.POST("developer/version/create", versionApi.Create)
	// 	r.POST("developer/version/update", versionApi.Update)
	// 	r.POST("developer/version/submitReview", versionApi.SubmitReview)
	// 	r.POST("developer/version/cancelReview", versionApi.CancelReview)
	// 	r.POST("developer/version/delete", versionApi.Delete)
	// 	r.POST("developer/version/upload", uploadApi.Upload)
	// }
}
