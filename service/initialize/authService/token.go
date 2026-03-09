package authService

import (
	"sun-panel/global"
	"sun-panel/lib/cache"

	"time"
)

// 记录每个用户最新唯一token
func InitUserAuthServiceClientToken() cache.Cacher[global.AuthServiceClientTokenUser] {
	return global.NewCache[global.AuthServiceClientTokenUser](720*time.Hour, 1*time.Hour, "UserAuthServiceClientToken")
}

// 记录用户登录返回的token（不唯一，参考 InitUserAuthServiceClientToken）
func InitCUserAuthServiceClientToken() cache.Cacher[string] {
	return global.NewCache[string](10*time.Hour, 10*time.Minute, "CUserAuthServiceClientToken")
}

// 客户端账号在线信息
func InitClientAccountOnlineCache() cache.Cacher[global.ClientLoginOnlineCacheInfoKeyValue] {
	return global.NewCache[global.ClientLoginOnlineCacheInfoKeyValue](1*time.Hour, 10*time.Minute, "ClientAccountOnlineCache")
}
