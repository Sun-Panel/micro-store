package admin

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	adminGroup := routerGroup.Group("admin")
	InitUserRouter(adminGroup)
	InitSystemSettingRouter(adminGroup)
	// InitNoticeRouter(adminGroup)
	InitAboutSetting(adminGroup)
	InitEmail(adminGroup)
	// InitEmailTemplate(adminGroup)
	// InitMessage(adminGroup)
	InitProAuthorizeRouter(adminGroup)
	// InitGoodsManage(adminGroup)
	// InitOrderManageRouter(adminGroup)
	InitSystemVariableRouter(adminGroup)
	InitMdPageRouter(adminGroup)
	InitDashboard(adminGroup)
	InitVersionRouter(adminGroup)
	InitClientSetting(adminGroup)
	InitRedeemCodeRouter(adminGroup)
}
