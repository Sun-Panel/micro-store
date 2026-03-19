package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func InitMicroAppRouter(router *gin.RouterGroup) {
	microAppApi := api_v1.ApiGroupApp.ApiAdmin.MicroAppApi
	versionApi := api_v1.ApiGroupApp.ApiAdmin.DeveloperVersionApi
	loginRouter := router.Group("", middleware.LoginInterceptor)

	// ==================== 管理员专用接口 ====================
	loginRouter.POST("microApp/getList", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), microAppApi.GetList)
	loginRouter.POST("microApp/getInfo", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), microAppApi.GetInfo)
	loginRouter.POST("microApp/updateStatus", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), microAppApi.UpdateStatus)

	// 管理员和开发者共享接口
	loginRouter.POST("microApp/deletes", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_DEVELOPER), microAppApi.Deletes)
	loginRouter.POST("microApp/offline", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_DEVELOPER), microAppApi.Offline)

	// ==================== 审核员专用接口 ====================
	loginRouter.POST("review/getPendingList", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), microAppApi.GetReviewPendingList)
	loginRouter.POST("review/getInfo", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), microAppApi.GetReviewInfo)
	loginRouter.POST("review/approve", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), microAppApi.ReviewApp)

	// ==================== 开发者专用接口（微应用） ====================
	loginRouter.POST("developer/microApp/create", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER|models.ROLE_ADMIN), microAppApi.Create)
	loginRouter.POST("developer/microApp/update", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.Update)
	loginRouter.POST("developer/microApp/updateLang", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.UpdateLang)
	loginRouter.POST("developer/microApp/cancelReview", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.CancelReview)
	loginRouter.POST("developer/microApp/getReviewHistory", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.GetReviewHistory)
	loginRouter.POST("developer/microApp/list", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.GetMyList)
	loginRouter.POST("developer/microApp/info", middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER), microAppApi.GetMyInfo)

	// ==================== 开发者专用接口（版本管理） ====================
	loginRouter.POST("developer/version/getList", versionApi.GetList)
	loginRouter.POST("developer/version/getInfo", versionApi.GetInfo)
	loginRouter.POST("developer/version/create", versionApi.Create)
	loginRouter.POST("developer/version/update", versionApi.Update)
	loginRouter.POST("developer/version/submitReview", versionApi.SubmitReview)
	loginRouter.POST("developer/version/cancelReview", versionApi.CancelReview)
	loginRouter.POST("developer/version/delete", versionApi.Delete)

	// ==================== 管理员/审核员专用接口（版本审核） ====================
	loginRouter.POST("reviewVersion/getPendingList", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), versionApi.GetPendingList)
	loginRouter.POST("reviewVersion/review", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), versionApi.Review)
	loginRouter.POST("reviewVersion/offline", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_AUDITOR), versionApi.Offline)
	loginRouter.POST("version/getList", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), versionApi.GetVersionList)
}
