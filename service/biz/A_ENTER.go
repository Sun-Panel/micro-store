package biz

import "sun-panel/biz/clientCache"

var (
	ProAuthorize      = new(ProAuthorizeType)
	Message           = new(MessageType)
	Captcha           = new(CaptchaType)
	PayPlatformConfig = new(PayPlatformConfigType)
	SunStore          = new(SunStoreType)
	Version           = new(VersionType)
	ClientCache       = new(clientCache.ClientCacheType)
	RedeemCode        = new(RedeemCodeType)
	VersionSecret     = versionSecret{}
	// AES          = new(AESType)
)

func InitBIZ() {
	ClientCache.Init()
	VersionSecret.Init()
	ProAuthorize.Init()
}
