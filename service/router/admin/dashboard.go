package admin

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitDashboard(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiAdmin.DashboardApi

	r := router.Group("dashboard", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("getStatistics", api.GetStatistics)
		r.POST("getActiveClientVersionStatistics", api.GetActiveClientVersionStatistics)
		r.POST("getVersions", api.GetVersions)
		r.POST("getVersionHistory", api.GetVersionHistory)

		// r.POST("getAmountLine", api.GetAmountLine)
		r.POST("getClientLine", api.GetClientLine)
		r.POST("getUserLine", api.GetUserLine)
	}

}
