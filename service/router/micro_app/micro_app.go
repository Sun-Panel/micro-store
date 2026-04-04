package microApp

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitMicroAppRouter(router *gin.RouterGroup) {
	microAppApi := api_v1.ApiGroupApp.MicroAppApi
	downloadApi := api_v1.ApiGroupApp.DownloadApi

	// ==================== 公开接口（微应用） ====================
	router.POST("microApp/getList", microAppApi.GetList)
	router.POST("microApp/getInfo", microAppApi.GetInfo)
	router.POST("microApp/version/getList", microAppApi.GetVersionList)

	router.POST("microApp/download/getUrl", downloadApi.GetUrl) // 获取下载链接

	// ==================== 下载接口（自动统计） ====================
	router.GET("microApp/download/:microAppId", downloadApi.DownloadByVersionOrLatest)          // 下载最新版本
	router.GET("microApp/download/:microAppId/:version", downloadApi.DownloadByVersionOrLatest) // 下载指定版本
}
