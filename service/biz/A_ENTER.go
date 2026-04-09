package biz

import (
	"sun-panel/biz/clientCache"
	"sun-panel/global"
	"sun-panel/lib/debugWrap"

	"github.com/redis/go-redis/v9"
)

var (
	Config             = new(ConfigType)
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
	MicroAppDownload   = new(microAppDownload)
	MicroAppAudit      = new(microAppAudit)
	// AES          = new(AESType)
)

var (
	// 日志包装器,仅在debug模式下生效
	//	global.Logger.Debugln(LdWrap.Log("test", LdWrap.Json("jsonData",jsonData),LdWrap.Data("data",data)
	LdWrap *debugWrap.DebugLogger

	// 日志包装器,任意模式下都进行解析
	//	global.Logger.Infoln(LogWrap.Log("test", LogWrap.Json("jsonData",jsonData),LogWrap.Data("data",data)
	LogWrap *debugWrap.DebugLogger
)

func InitBIZ(redisClient *redis.Client) {
	ClientCache.Init()
	MicroAppPackage.Init()
	MicroAppStatistics.Init(redisClient)
	MicroAppDownload.Init()
	LdWrap = debugWrap.NewDebugLogger(global.RUNCODE == "debug") // 仅在debug模式下生效
	LogWrap = debugWrap.NewDebugLogger(true)                     // 任意模式下都进行解析
}
