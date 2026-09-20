package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func InitCustomCodeRouter(router *gin.RouterGroup) {
	customCodeApi := api_v1.ApiGroupApp.ApiAdmin.CustomCodeApi
	reviewApi := api_v1.ApiGroupApp.ApiAdmin.CustomCodeReviewApi

	// 作者（开发者）
	author := router.Group("customCode", middleware.LoginInterceptor, middleware.DeveloperInterceptor)
	{
		author.POST("/getMyList", customCodeApi.GetMyList)
		author.POST("/getInfo", customCodeApi.GetInfo)
		author.POST("/edit", customCodeApi.Edit)
		author.POST("/withdraw", customCodeApi.Withdraw)
		author.POST("/offline", customCodeApi.Offline)
		author.POST("/delete", customCodeApi.Delete)
		author.POST("/uploadPreviewImage", customCodeApi.UploadPreviewImage)
	}

	// 管理员
	adminGroup := router.Group("customCode", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		adminGroup.POST("/forceOffline", customCodeApi.ForceOffline)
	}

	// 审核（审核员或管理员）
	auditor := router.Group("customCodeReview", middleware.LoginInterceptor,
		middleware.MultiRolesInterceptor(models.ROLE_AUDITOR|models.ROLE_ADMIN))
	{
		auditor.POST("/getList", reviewApi.GetList)
		auditor.POST("/getInfo", reviewApi.GetInfo)
		auditor.POST("/approve", reviewApi.Approve)
		auditor.POST("/reject", reviewApi.Reject)
		auditor.POST("/getHistory", reviewApi.GetHistory)
	}
}
