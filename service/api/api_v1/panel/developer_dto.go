package panel

// 开发者 API 请求参数定义

// DeveloperRegisterReq 开发者注册请求
type DeveloperRegisterReq struct {
	DeveloperName string `json:"developerName" binding:"required"`
	ContactMail   string `json:"contactMail"`
	PaymentName   string `json:"paymentName"`
	PaymentQrcode string `json:"paymentQrcode"`
	PaymentMethod string `json:"paymentMethod"`
	Name          string `json:"name"`
}

// DeveloperUpdateReq 更新开发者信息请求
type DeveloperUpdateReq DeveloperRegisterReq
