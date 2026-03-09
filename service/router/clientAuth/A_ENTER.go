package clientAuth

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitAuth(routerGroup)
}
