package microApp

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitMicroAppRouter(router *gin.RouterGroup) {
	apiGroup := api_v1.ApiGroupApp
	versionApi := api_v1.ApiGroupApp.MicroAppApi

	// ==================== 公开接口（微应用） ====================
	router.POST("microApp/getInfo", apiGroup.MicroAppApi.GetInfo)
	router.POST("microApp/version/getList", versionApi.GetList)
}
