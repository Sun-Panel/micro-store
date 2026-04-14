package admin

// 开发者管理 API 请求参数定义

// GetListRequest 获取开发者列表请求
type DeveloperGetListReq struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Status  *int   `json:"status"`
	KeyWord string `json:"keyWord"`
}

// GetInfoRequest 获取开发者详情请求
type DeveloperGetInfoReq struct {
	Id uint `json:"id" binding:"required"`
}

// GetByDeveloperNameRequest 根据开发者标识获取开发者信息请求
type DeveloperGetByDeveloperNameReq struct {
	DeveloperName string `json:"developerName" binding:"required"`
}

// UpdateRequest 更新开发者请求
type DeveloperUpdateReq struct {
	Id            uint   `json:"id" binding:"required"`
	DeveloperName string `json:"developerName" binding:"required"`
	ContactMail   string `json:"contactMail"`
	PaymentName   string `json:"paymentName"`
	PaymentQrcode string `json:"paymentQrcode"`
	PaymentMethod string `json:"paymentMethod"`
	Name          string `json:"name"`
	Status        int    `json:"status"`
}

// DeletesRequest 删除开发者请求
type DeveloperDeletesReq struct {
	Ids []uint `json:"ids" binding:"required"`
}

// UpdateStatusRequest 更新状态请求
type DeveloperUpdateStatusReq struct {
	Id     uint `json:"id" binding:"required"`
	Status int  `json:"status" binding:"required"`
}
