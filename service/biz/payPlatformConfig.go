package biz

import "sun-panel/global"

type PayPlatformConfigAliPay struct {
	Appid           string `json:"appid"`
	PrivateKey      string `json:"privateKey"`
	AlipayPublicKey string `json:"alipayPublicKey"`
}

type PayPlatformConfigPaddle struct {
	ApiKey                string `json:"apiKey"`                // 调用API所使用的key
	NotificationSecretKey string `json:"notificationSecretKey"` // 通知webhook的验证签名的key
	ClientSideToken       string `json:"clientSideToken"`       // paddle.js调用所使用的token
}

type PayPlatformConfigType struct {
	AliPay PayPlatformConfigAliPay `json:"aliPay"`
	Paddle PayPlatformConfigPaddle `json:"paddle"`
}

const PayPlatformSystemConfig = "pay_platform_config"

func (p *PayPlatformConfigType) GetPaddle() *PayPlatformConfigPaddle {
	c := PayPlatformConfigPaddle{
		ApiKey:                global.Config.GetValueString("paddle", "api_key"),
		NotificationSecretKey: global.Config.GetValueString("paddle", "notification_secret_key"),
		ClientSideToken:       global.Config.GetValueString("paddle", "client_side_token"),
	}
	return &c
	// return &p.Get().Paddle
}

// func (p *PayPlatformConfigType) GetAliPay() *PayPlatformConfigAliPay {
// 	return &p.Get().AliPay
// }

// func (p *PayPlatformConfigType) Get() *PayPlatformConfigType {
// 	v := PayPlatformConfigType{}
// 	err := global.SystemSetting.GetValueByInterface(PayPlatformSystemConfig, &v)
// 	if err != nil {
// 		return nil
// 	}
// 	return &v
// }

func (p *PayPlatformConfigType) Set(v PayPlatformConfigType) error {
	return global.SystemSetting.Set(PayPlatformSystemConfig, &v)
}
