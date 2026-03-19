package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func AdminInterceptor(c *gin.Context) {
	currentUser, _ := base.GetCurrentUserInfo(c)
	// 使用位运算判断是否包含管理员角色
	if !models.HasAdmin(currentUser.Role) {
		apiReturn.ErrorNoAccess(c)
		c.Abort()
		return
	}
}
