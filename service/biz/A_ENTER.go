package biz

import "sun-panel/biz/clientCache"

var (
	Message                  = new(MessageType)
	Captcha                  = new(CaptchaType)
	PayPlatformConfig        = new(PayPlatformConfigType)
	SunStore                 = new(SunStoreType)
	ClientCache              = new(clientCache.ClientCacheType)
	RedeemCode               = new(RedeemCodeType)
	Developer                = new(DeveloperService)
	MicroAppDeveloper        = new(MicroAppDeveloperService)
	// AES          = new(AESType)
)

func InitBIZ() {
	ClientCache.Init()
}
