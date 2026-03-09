package middleware

import (
	"sun-panel/apiClientApp/v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

func LoginInterceptor(c *gin.Context) {

	// 继续执行后续的操作，再回来
	// c.Next()

	// 获得token
	cToken := c.GetHeader("token")

	// 没有token信息视为未登录
	if cToken == "" {
		apiReturn.ErrorByCode(c, apiReturn.ErrNotLogged)
		c.Abort() // 终止执行后续的操作，一般配合return使用
		return
	}

	token := ""
	{
		var ok bool
		token, ok = global.CUserAuthServiceClientToken.Get(cToken)
		// 可能已经安全退出或者很久没有使用已过期
		if !ok || token == "" {
			apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
			c.Abort() // 终止执行后续的操作，一般配合return使用
			return
		}
	}

	// 找到缓存直接返回缓存的用户信息
	if clientTokenInfo, success := global.UserAuthServiceClientToken.Get(token); success {
		c.Set("userInfo", clientTokenInfo.User)
		return
	}

	// 否则认定过期，需要重新密码和账号登录，一般可以用于修改密码后
	apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
	c.Abort()

	// global.Logger.Debug("准备查询数据库的用户资料", token)

	// mUser := models.User{}
	// // 去库中查询是否存在该用户；否则返回错误
	// if info, err := mUser.GetUserInfoByToken(token); err != nil || info.Token == "" || info.ID == 0 {

	// 	apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
	// 	c.Abort()
	// 	return
	// } else {
	// 	clientTokenInfo := global.AuthServiceClientTokenUser{
	// 		Ctoken: info.Token,
	// 		User:   info,
	// 	}
	// 	// 通过 设置当前用户信息
	// 	global.UserAuthServiceClientToken.SetDefault(info.Token, clientTokenInfo)
	// 	global.CUserAuthServiceClientToken.SetDefault(cToken, token)
	// 	c.Set("userInfo", clientTokenInfo.User)
	// }

}

// 不验证缓存直接验证库省去没有缓存每次都要手动登录的问题
func LoginInterceptorDev(c *gin.Context) {

	// 获得token
	token := c.GetHeader("token")
	mUser := models.User{}

	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.ID == 0 {
		apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}
