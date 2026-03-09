package browserExtension

import (
	v1 "sun-panel/apiBrowserExtensionClient/v1"

	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	InitStatistics(routerGroup)
}

func InitStatistics(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.Statistics
	router.POST("v1/statistics/push", api.Push)
}
