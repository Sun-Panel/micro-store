package admin

// 微应用管理 API 请求参数定义

// GetListRequest 获取微应用列表请求
type MicroAppGetListReq struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	CategoryId *int   `json:"categoryId"`
	Status     *int   `json:"status"`
	KeyWord    string `json:"keyWord"`
}

// GetInfoRequest 获取微应用详情请求
type MicroAppGetInfoReq struct {
	Id uint `json:"id" binding:"required"`
}

// LangInfo 多语言信息
type MicroAppLangInfo struct {
	AppName string `json:"appName"`
	AppDesc string `json:"appDesc"`
}

// CreateRequest 创建微应用请求
type MicroAppCreateReq struct {
	MicroAppId  string                      `json:"microAppId" binding:"required"`
	AppName     string                      `json:"appName" binding:"required"`
	AppIcon     string                      `json:"appIcon" binding:"required"`
	AppDesc     string                      `json:"appDesc"`
	Remark      string                      `json:"remark"`
	CategoryId  int                         `json:"categoryId" binding:"required"`
	ChargeType  int                         `json:"chargeType"`
	Price       float64                     `json:"price"`
	AuthorId    uint                        `json:"authorId" binding:"required"`
	Screenshots string                      `json:"screenshots"`
	LangMap     map[string]MicroAppLangInfo `json:"langMap"`
}

// UpdateRequest 更新微应用请求
type MicroAppUpdateReq struct {
	Id          uint                        `json:"id" binding:"required"`
	AppName     string                      `json:"appName" binding:"required"`
	AppIcon     string                      `json:"appIcon" binding:"required"`
	AppDesc     string                      `json:"appDesc"`
	Remark      string                      `json:"remark"`
	CategoryId  int                         `json:"categoryId" binding:"required"`
	ChargeType  int                         `json:"chargeType"`
	Price       float64                     `json:"price"`
	Screenshots string                      `json:"screenshots"`
	LangMap     map[string]MicroAppLangInfo `json:"langMap"`
}

// DeletesRequest 删除微应用请求
type MicroAppDeletesReq struct {
	Ids []uint `json:"ids" binding:"required"`
}

// UpdateStatusRequest 更新状态请求
type MicroAppUpdateStatusReq struct {
	Id     uint `json:"id" binding:"required"`
	Status int  `json:"status" binding:"required"`
}

// UpdateLangRequest 更新语言请求
type MicroAppUpdateLangReq struct {
	Id      uint                        `json:"id" binding:"required"`
	LangMap map[string]MicroAppLangInfo `json:"langMap"`
}
