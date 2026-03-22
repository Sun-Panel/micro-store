package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

func DeveloperInterceptor(c *gin.Context) {
	global.Logger.Debugln("DeveloperInterceptor", "开发者拦截器")
	userInfo, _ := base.GetCurrentUserInfo(c)
	developer, err := biz.Developer.GetByUserId(global.Db, userInfo.ID)
	global.Logger.Debugln("current developer", developer)
	// 当前账号并非开发者
	if err != nil || developer.DeveloperName == "" {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoCurrentPermission)
		c.Abort()
		return
	}

	c.Set("developerInfo", developer)
}
