package global

import (
	"sun-panel/lib/cache"
	"sun-panel/structs"
	"time"
)

// 缓存驱动
const (
	CACHE_DRIVE_REDIS  = "redis"
	CACHE_DRIVE_MEMORY = "memory"
)

type CacherOption struct {
	CacheAreaExpired time.Duration // 整个缓存区的默认过期时长(暂时仅支持redis)
}

// 创建一个缓存区(redis是hset)
// | defaultExpiration:每个key的默认过期时长
// | cleanupInterval:清理过期的key间隔 0.不清理
// | name:缓存名称
func NewCache[T any](defaultExpiration time.Duration, cleanupInterval time.Duration, name string, extensionOption ...CacherOption) cache.Cacher[T] {
	drive := Config.GetValueString("base", "cache_drive")
	if drive == "" {
		drive = CACHE_DRIVE_MEMORY
	}
	var cacher cache.Cacher[T]
	Logger.Debugln("缓存驱动:", drive)
	switch drive {
	case CACHE_DRIVE_MEMORY:
		cacher = cache.NewGoCache[T](defaultExpiration, cleanupInterval)
	case CACHE_DRIVE_REDIS:
		redisConfig := structs.IniConfigRedis{}
		if err := Config.GetSection("redis",
			&redisConfig); err != nil {
			redisConfig.Prefix = ""
		}

		var hashExpired time.Duration = 0

		if len(extensionOption) > 0 {
			hashExpired = extensionOption[0].CacheAreaExpired
		}

		opt := cache.RedisCacheOptions{
			Redis:             RedisDb,
			HashKey:           redisConfig.Prefix + name,
			DefaultExpiration: defaultExpiration,
			CleanupInterval:   cleanupInterval,
			HashKeyExpiration: hashExpired,
		}
		cacher = cache.NewRedisCache[T](opt)
	}

	return cacher
}
