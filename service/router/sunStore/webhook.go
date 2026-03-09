package sunstore

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitWebhook(routerGroup *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSunStore

	routerGroup.POST("goodsOrder", api.GoodsOrder.Enter)
	routerGroup.POST("user", api.User.Enter)
}
