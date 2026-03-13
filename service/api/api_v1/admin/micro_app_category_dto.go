package admin

// 分类管理 API 请求参数定义

// GetListRequest 获取分类列表请求
type MicroAppCategoryGetListReq struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Status  *int   `json:"status"`
	KeyWord string `json:"keyWord"`
}

// GetInfoRequest 获取分类详情请求
type MicroAppCategoryGetInfoReq struct {
	Id uint `json:"id" binding:"required"`
}

// CreateRequest 创建分类请求
type MicroAppCategoryCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// UpdateRequest 更新分类请求
type MicroAppCategoryUpdateReq struct {
	Id     uint   `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// DeletesRequest 删除分类请求
type MicroAppCategoryDeletesReq struct {
	Ids []uint `json:"ids" binding:"required"`
}

// UpdateStatusRequest 更新状态请求
type MicroAppCategoryUpdateStatusReq struct {
	Id     uint `json:"id" binding:"required"`
	Status int  `json:"status" binding:"required"`
}
