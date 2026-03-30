package admin

// ==================== 开发者端版本管理 DTO ====================

// MicroAppVersionGetListReq 获取版本列表请求
type MicroAppVersionGetListReq struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	AppRecordId uint `json:"appRecordId"`
	Status      *int `json:"status"`
}

// MicroAppVersionGetInfoReq 获取版本详情请求
type MicroAppVersionGetInfoReq struct {
	Id uint `json:"id" binding:"required"`
}

// MicroAppVersionCreateReq 创建版本请求
type MicroAppVersionCreateReq struct {
	// AppRecordId   uint                          `json:"appId" binding:"required"`
	// Version       string                        `json:"version" binding:"required"`
	// VersionCode   int                           `json:"versionCode" binding:"required"`
	// PackageUrl    string                        `json:"packageUrl" binding:"required"`
	// PackageHash   string                        `json:"packageHash"`
	VersionDesc string `json:"versionDesc"`
	// Config        *models.MicroAppVersionConfig `json:"config"` // 完整配置信息
	UploadCacheId string `json:"uploadCacheId"`
}

// MicroAppVersionUpdateReq 更新版本请求
type MicroAppVersionUpdateReq struct {
	Id          uint   `json:"id" binding:"required"`
	Version     string `json:"version"`
	VersionCode int    `json:"versionCode"`
	VersionDesc string `json:"versionDesc"`
}

// MicroAppVersionSubmitReviewReq 提交审核请求
type MicroAppVersionSubmitReviewReq struct {
	VersionId uint `json:"versionId" binding:"required"`
}

// MicroAppVersionCancelReviewReq 撤销审核请求
type MicroAppVersionCancelReviewReq struct {
	VersionId uint `json:"versionId" binding:"required"`
}

// MicroAppVersionDeleteReq 删除版本请求
type MicroAppVersionDeleteReq struct {
	Ids []uint `json:"ids" binding:"required"`
}

// MicroAppVersionOfflineReq 版本下架请求
type MicroAppVersionOfflineReq struct {
	Id     uint   `json:"id" binding:"required"`
	Type   int    `json:"type" binding:"required"` // 下架类型：1-作者下架 2-平台下架
	Reason string `json:"reason"`                  // 下架原因
}

// ==================== 管理员端版本审核 DTO ====================

// AdminVersionGetListReq 管理员获取版本列表请求
type AdminVersionGetListReq struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	AppRecordId uint `json:"appId"`
	Status      *int `json:"status"`
}

// AdminVersionGetPendingListReq 获取待审核列表请求
type AdminVersionGetPendingListReq struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	KeyWord string `json:"keyWord"`
}

// AdminVersionReviewReq 审核版本请求
type AdminVersionReviewReq struct {
	VersionId  uint   `json:"versionId" binding:"required"`
	Status     int    `json:"status" binding:"required"`
	ReviewNote string `json:"reviewNote"`
}
