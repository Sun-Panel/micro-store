package global

import (
	"sun-panel/initialize/database"
	"sun-panel/lib/cache"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/iniConfig"
	"sun-panel/lib/language"
	"sun-panel/models"

	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ISDOCKER     = ""           // 是否为docker模式运行 true | [string]
	RUNCODE      = "debug"      // 运行模式：debug | release
	UPLOAD_ROUTE = "/uploads"   // 上传路由
	PAY_ENV      = "production" // 支付环境 production:生产环境 | sandbox:沙盒环境
	DB_DRIVER    = database.SQLITE
)

// var Log *cmn.LogStruct

var (
	Lang *language.LangStructObj

	UserToken                cache.Cacher[models.User]
	CUserToken               cache.Cacher[string] // 用户token
	CUserApiTokenAccessToken cache.Cacher[string] // 单点退出-主动 用户apitoken绑定AccessToken
	CUserAccessTokenApiToken cache.Cacher[string] // 单点退出-被动 用户AccessToken绑定apitoken

	Logger              *zap.SugaredLogger
	LoggerLevel         = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	VerifyCodeCachePool cache.Cacher[string]
	Config              *iniConfig.IniConfig
	Db                  *gorm.DB
	RedisDb             *redis.Client
	SystemSetting       *systemSetting.SystemSettingCache
	SystemMonitor       cache.Cacher[interface{}]
	RateLimit           *RateLimiter
)
