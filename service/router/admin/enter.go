package admin

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	adminGroup := routerGroup.Group("admin")
	InitUserRouter(adminGroup)
	InitSystemSettingRouter(adminGroup)
	InitAboutSetting(adminGroup)
	InitEmail(adminGroup)
	InitSystemVariableRouter(adminGroup)
	InitMdPageRouter(adminGroup)
	InitDashboard(adminGroup)
	InitClientSetting(adminGroup)
	InitRedeemCodeRouter(adminGroup)
	InitMicroAppCategoryRouter(adminGroup)
	InitDeveloperRouter(adminGroup)
	InitMicroAppRouter(adminGroup)
	InitDeveloperVersionRouter(adminGroup)
	InitVersionReviewRouter(adminGroup)
}
