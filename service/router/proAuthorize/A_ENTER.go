package proAuthorize

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitProAuthorizeRouter(routerGroup)
}
