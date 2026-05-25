package microapp

type MicroAppPublic struct {
}

type MicroAppVersionGetVersionListReq struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	AppRecordId uint `json:"appRecordId"`
}

type MicroAppVersionGetListReq struct {
	Page            int    `json:"page"`
	Limit           int    `json:"limit"`
	Order           string `json:"order"`
	CategoryId      uint   `json:"categoryId"`
	Keyword         string `json:"keyword"`
	OnlyWithVersion bool   `json:"onlyWithVersion"`
}

// GetInfoByMicroAppIdReq 根据 microAppId 获取详情的请求
type GetInfoByMicroAppIdReq struct {
	MicroAppId string `json:"microAppId" binding:"required"`
}
