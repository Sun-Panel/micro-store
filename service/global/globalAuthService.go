package global

import (
	"sun-panel/lib/cache"
	"sun-panel/models"
)

type AuthServiceClientTokenUser struct {
	User   models.User
	Ctoken string
}

type ClientLoginOnlineCacheInfo struct {
	IP        string
	LanIP     string
	Timestamp int64
	Ctoken    string // 临时token
	ClientID  string
}

// 临时token - 设备信息
type ClientLoginOnlineCacheInfoKeyValue map[string]ClientLoginOnlineCacheInfo

var (
	// 授权客户端的用户授权
	UserAuthServiceClientToken  cache.Cacher[AuthServiceClientTokenUser]
	CUserAuthServiceClientToken cache.Cacher[string]

	// 用户登录尝试次数
	ClientLoginAttemptsCacheCache cache.Cacher[int]

	// 客户端账号在线信息 // key:userId - value:ClientLoginOnlineCacheInfoKeyValue
	ClientAccountOnlineCache cache.Cacher[ClientLoginOnlineCacheInfoKeyValue]
)
