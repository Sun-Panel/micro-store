package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func InitMicroAppRouter(router *gin.RouterGroup) {
	apiGroup := api_v1.ApiGroupApp.ApiAdmin
	loginRouter := router.Group("", middleware.LoginInterceptor)

	// ==================== 管理员专用接口（微应用） ====================
	loginRouter.POST("microApp/getList", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.MicroAppAdminApi.GetList)
	// loginRouter.POST("microApp/getInfo", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.MicroAppAdminApi.GetInfo)
	loginRouter.POST("microApp/updateStatus", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.MicroAppAdminApi.UpdateStatus)

	// 管理员和开发者共享接口
	loginRouter.POST("microApp/deletes", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_DEVELOPER), apiGroup.MicroAppAdminApi.Deletes)
	loginRouter.POST("microApp/offline", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_DEVELOPER), apiGroup.MicroAppAdminApi.Offline)

	// ==================== 审核员专用接口（微应用） ====================
	auditorRouter := loginRouter.Group("", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR))
	{
		// 微应用相关
		auditorRouter.POST("review/microApp/getPendingList", apiGroup.MicroAppAuditorApi.GetReviewPendingList)

		// 微应用版本相关
		auditorRouter.POST("review/microApp/version/getList", apiGroup.MicroAppVersionAdminApi.GetVersionList)
		auditorRouter.POST("review/microApp/version/getPendingList", apiGroup.MicroAppVersionAdminApi.GetPendingList)
		auditorRouter.POST("review/microApp/version/getLatestOnlineByAppModelId", apiGroup.MicroAppVersionAdminApi.GetLatestOnlineByAppModelId)

		auditorRouter.POST("review/getReviewInfo", apiGroup.MicroAppAuditorApi.GetReviewInfo)
		auditorRouter.POST("review/getMicroAppInfo", apiGroup.MicroAppAuditorApi.GetMicroAppInfo)
		auditorRouter.POST("review/approve", apiGroup.MicroAppAuditorApi.ReviewApp)

		auditorRouter.POST("reviewVersion/review", apiGroup.MicroAppVersionAdminApi.Review)
		// auditorRouter.POST("reviewVersion/offline", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_AUDITOR), apiGroup.MicroAppVersionAdminApi.Offline)

	}

	// ==================== 开发者专用接口====================
	myMicroAppRouter := loginRouter.Group("",
		middleware.MultiRolesInterceptor(models.ROLE_DEVELOPER),
		middleware.DeveloperInterceptor)
	{
		// ==================== 微应用 ====================
		myMicroAppRouter.POST("developer/myMicroApp/create", apiGroup.MicroAppDeveloperApi.Create)
		myMicroAppRouter.POST("developer/myMicroApp/update", apiGroup.MicroAppDeveloperApi.Update)
		// myMicroAppRouter.POST("developer/myMicroApp/getInfo", apiGroup.MicroAppDeveloperApi.GetInfo)
		myMicroAppRouter.POST("developer/myMicroApp/updateLang", apiGroup.MicroAppDeveloperApi.UpdateLang)
		myMicroAppRouter.POST("developer/myMicroApp/submitReview", apiGroup.MicroAppDeveloperApi.SubmitReview)
		myMicroAppRouter.POST("developer/myMicroApp/cancelReview", apiGroup.MicroAppDeveloperApi.CancelReview)
		myMicroAppRouter.POST("developer/myMicroApp/getReviewHistory", apiGroup.MicroAppDeveloperApi.GetReviewHistory)
		myMicroAppRouter.POST("developer/myMicroApp/list", apiGroup.MicroAppDeveloperApi.GetMyList)
		myMicroAppRouter.POST("developer/myMicroApp/info", apiGroup.MicroAppDeveloperApi.GetInfo)
		myMicroAppRouter.POST("developer/myMicroApp/getMicroInfoAndReviewInfoByMicroAppModelId", apiGroup.MicroAppDeveloperApi.GetMicroInfoAndReviewInfoByMicroAppModelId)

		// ==================== 微应用版本管理 ====================
		myMicroAppRouter.POST("developer/version/getList", apiGroup.DeveloperVersionApi.GetList)
		myMicroAppRouter.POST("developer/version/getInfo", apiGroup.DeveloperVersionApi.GetInfo)
		myMicroAppRouter.POST("developer/version/create", apiGroup.DeveloperVersionApi.Create)
		myMicroAppRouter.POST("developer/version/update", apiGroup.DeveloperVersionApi.Update)
		myMicroAppRouter.POST("developer/version/submitReview", apiGroup.DeveloperVersionApi.SubmitReview)
		myMicroAppRouter.POST("developer/version/cancelReview", apiGroup.DeveloperVersionApi.CancelReview)
		myMicroAppRouter.POST("developer/version/delete", apiGroup.DeveloperVersionApi.Delete)
	}

	// // ==================== 管理员/审核员专用接口（版本审核） ====================
	// loginRouter.POST("reviewVersion/getPendingList", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), apiGroup.MicroAppVersionAdminApi.GetPendingList)
	// // loginRouter.POST("reviewVersion/getMicroAppInfo", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), apiGroup.MicroAppVersionAdminApi.GetMicroAppInfo)
	// loginRouter.POST("reviewVersion/review", middleware.MultiRolesInterceptor(models.ROLE_AUDITOR), apiGroup.MicroAppVersionAdminApi.Review)
	// loginRouter.POST("reviewVersion/offline", middleware.MultiRolesInterceptor(models.ROLE_ADMIN|models.ROLE_AUDITOR), apiGroup.MicroAppVersionAdminApi.Offline)
	// loginRouter.POST("version/getList", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.MicroAppVersionAdminApi.GetVersionList)

	// ==================== 管理员专用接口（开发者用户管理） ====================
	loginRouter.POST("developer/user/getList", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.GetList)
	loginRouter.POST("developer/user/getInfo", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.GetInfo)
	loginRouter.POST("developer/user/getByDeveloperName", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.GetByDeveloperName)
	loginRouter.POST("developer/user/update", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.Update)
	loginRouter.POST("developer/user/deletes", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.Deletes)
	loginRouter.POST("developer/user/updateStatus", middleware.MultiRolesInterceptor(models.ROLE_ADMIN), apiGroup.DeveloperApi.UpdateStatus)
}
