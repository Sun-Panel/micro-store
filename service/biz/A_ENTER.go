package biz

import (
	"sun-panel/biz/clientCache"

	"github.com/redis/go-redis/v9"
)

var (
	Message            = new(MessageType)
	Captcha            = new(CaptchaType)
	PayPlatformConfig  = new(PayPlatformConfigType)
	SunStore           = new(SunStoreType)
	ClientCache        = new(clientCache.ClientCacheType)
	RedeemCode         = new(RedeemCodeType)
	Developer          = new(DeveloperService)
	MicroAppDeveloper  = new(MicroAppDeveloperService)
	MicroApp           = new(microApp)
	MicroAppVersion    = new(MicroAppVersionService)
	MicroAppPackage    = new(MicroAppPackageService)
	MicroAppStatistics = new(microAppStatistics)
	// AES          = new(AESType)
)

func InitBIZ(redisClient *redis.Client) {
	ClientCache.Init()
	MicroAppPackage.Init()
	MicroAppStatistics.Init(redisClient)
}
