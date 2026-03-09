package authService

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"time"
)

// 记录用户登录尝试次数
func InitClientLoginAttemptsCache() cache.Cacher[int] {
	return global.NewCache[int](10*time.Minute, 20*time.Minute, "ClientLoginAttempt")
}
