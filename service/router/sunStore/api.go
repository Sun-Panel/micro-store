package sunstore

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitApi(routerGroup *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSunStore

	routerGroup.POST("onlyProGoodsBuyQualificationCheck", api.Api.OnlyProGoodsBuyQualificationCheck)
}
