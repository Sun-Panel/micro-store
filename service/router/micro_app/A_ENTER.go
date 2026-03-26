package microApp

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitMicroAppRouter(routerGroup)
}
