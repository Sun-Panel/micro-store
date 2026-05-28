package types

type BaseResponse[T any] struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       T      `json:"data"`
	ResponseId string `json:"responseId"`
}

// ==================== 批量获取作者昵称 ====================

// BatchGetDeveloperInfoReq 批量获取作者昵称 - 请求
type BatchGetDeveloperInfoReq struct {
	DeveloperNames []string `json:"developerNames" binding:"required,min=1,max=100"`
}

// DeveloperInfo 开发者信息（对外暴露）
type DeveloperInfo struct {
	DeveloperName string `json:"developerName"` // 开发者标识
	Name          string `json:"name"`          // 开发者昵称
}

// BatchGetDeveloperInfoResp 批量获取作者昵称 - 响应
type BatchGetDeveloperInfoResp struct {
	DeveloperInfos map[string]DeveloperInfo `json:"developerInfos"` // key: developerName
}

// ==================== 微应用安装触发记录 ====================

// MicroAppInstallRecordReq 微应用安装触发记录 - 请求
type MicroAppInstallRecord struct {
	MicroAppId       string `json:"microAppId" binding:"required"`
	MicroAppVersion  string `json:"microAppVersion" binding:"required"`
	InstallTimestamp int64  `json:"installTimestamp" binding:"required"`
}

type MicroAppInstallRecordReq struct {
	InstallList []MicroAppInstallRecord `json:"installList" binding:"required,min=1,max=100"`
}

// ==================== 批量查询微应用信息 ====================

// MicroAppInfo 微应用信息（对外暴露）
type MicroAppInfo struct {
	MicroAppId     string `json:"microAppId"`
	AppIcon        string `json:"appIcon"`
	ChargeType     int    `json:"chargeType"`
	Points         int    `json:"points"`
	DeveloperName  string `json:"developerName"` // 开发者名称
	DeveloperName2 string `json:"developer"`     // 开发者标识
}

// BatchGetMicroAppInfoReq 批量查询微应用信息 - 请求
type BatchGetMicroAppInfoReq struct {
	MicroAppIds []string `json:"microAppIds" binding:"required,min=1,max=100"`
}

// BatchGetMicroAppInfoResp 批量查询微应用信息 - 响应
type BatchGetMicroAppInfoResp struct {
	MicroAppInfos map[string]*MicroAppInfo `json:"microAppInfos"` // key: microAppId，不存在的不返回
}

// ==================== 查询当前账号是否为开发者 ====================

// CheckIsDeveloperResp 查询当前账号是否为开发者 - 响应
type CheckIsDeveloperResp struct {
	IsDeveloper bool `json:"isDeveloper"`
}
