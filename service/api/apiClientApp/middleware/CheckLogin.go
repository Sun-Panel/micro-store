package middleware

import (
	"sun-panel/api/apiClientApp/common/apiReturn"
	"sun-panel/api/apiClientApp/common/base"
	"sun-panel/biz"

	"github.com/gin-gonic/gin"
)

// ResolveCurrentUser 解析当前用户信息（不拦截）
// 尝试获取用户信息并存入上下文，获取失败不拦截请求
func ResolveCurrentUser(c *gin.Context) {
	username := c.GetString("username")
	if username != "" {
		user, err := biz.User.GetUser(username)
		if err == nil {
			c.Set("currentUser", user)
		}
	}
	c.Next()
}

// LoginInterceptor 登录拦截中间件
// 验证用户是否已登录，未登录则拦截请求
func LoginInterceptor(c *gin.Context) {
	// 获取用户（优先从上下文，否则查库）
	user, ok := base.GetCurrentUser(c)
	if !ok {
		username := c.GetString("username")
		if username == "" {
			apiReturn.ErrorNotLogin(c)
			c.Abort()
			return
		}
		var err error
		user, err = biz.User.GetUser(username)
		if err != nil {
			apiReturn.ErrorNotLogin(c)
			c.Abort()
			return
		}
	}

	c.Set("currentUser", user)
	c.Next()
}
