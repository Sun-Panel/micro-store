package models

import (
	"encoding/json"
	"sun-panel/models/datatype"
	"time"

	"gorm.io/gorm"
)

// MicroAppVersionConfig 完整的应用配置（对应 app.config.json）
type MicroAppVersionConfig struct {
	AppJsonVersion string                     `json:"appJsonVersion"` // JSON 版本
	MicroAppId     string                     `json:"microAppId"`     // 应用唯一标识
	Version        string                     `json:"version"`        // 版本号
	APIVersion     string                     `json:"apiVersion"`     // API 版本
	Author         string                     `json:"author"`         // 作者
	Entry          string                     `json:"entry"`          // 入口文件
	Icon           string                     `json:"icon"`           // 图标
	Debug          bool                       `json:"debug"`          // 调试模式
	Components     map[string]json.RawMessage `json:"components"`     // 组件配置
	Permissions    []string                   `json:"permissions"`    // 权限列表
	DataNodes      map[string]json.RawMessage `json:"dataNodes"`      // 数据节点配置
	NetworkDomains []string                   `json:"networkDomains"` // 网络域名白名单
	AppInfo        AppInfoField               `json:"appInfo"`        // 应用多语言信息

	// 临时字段：仅用于 v1.1 解析阶段，解析完成后清除
	Locales       map[string]string `json:"locales,omitempty"`       // 语言文件映射：语言代码 -> 文件名
	DefaultLocale string            `json:"defaultLocale,omitempty"` // 默认语言
}

// AppInfo 应用多语言信息
type AppInfo struct {
	AppName            string `json:"appName"`
	Description        string `json:"description"`
	NetworkDescription string `json:"networkDescription"`
}

// AppInfoField 兼容 v1.0 和 v1.1 格式的 AppInfo 字段
// v1.0: {"zh-CN": {...}, "en-US": {...}}
// v1.1: {"appName": "...", "description": "..."} → 转换为 {"default": {...}}
type AppInfoField struct {
	LangMap map[string]AppInfo
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (a *AppInfoField) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// 检查是否是 v1.1 格式（扁平对象，没有语言代码键）
	if len(raw) > 0 {
		// 获取第一个键，检查是否是语言代码格式
		firstKey := ""
		for k := range raw {
			firstKey = k
			break
		}

		// v1.1 格式：键是 appName, description 等
		if firstKey == "appName" || firstKey == "description" || firstKey == "networkDescription" {
			// 转换为 v1.0 格式：包装为 {"default": {...}}
			var appInfo AppInfo
			if err := json.Unmarshal(data, &appInfo); err != nil {
				return err
			}
			a.LangMap = map[string]AppInfo{
				"default": appInfo,
			}
			return nil
		}
	}

	// v1.0 格式：多语言映射
	a.LangMap = make(map[string]AppInfo)
	for lang, rawInfo := range raw {
		var appInfo AppInfo
		if err := json.Unmarshal(rawInfo, &appInfo); err != nil {
			return err
		}
		a.LangMap[lang] = appInfo
	}

	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (a AppInfoField) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.LangMap)
}

// GetAppInfo 获取指定语言的应用信息，如果没有则返回默认语言
func (a *AppInfoField) GetAppInfo(lang string) (AppInfo, bool) {
	if a.LangMap == nil {
		return AppInfo{}, false
	}

	// 尝试获取指定语言
	if info, ok := a.LangMap[lang]; ok {
		return info, true
	}

	// 尝试获取默认语言
	if info, ok := a.LangMap["default"]; ok {
		return info, true
	}

	// 返回第一个可用的语言
	for _, info := range a.LangMap {
		return info, true
	}

	return AppInfo{}, false
}

// 微应用版本列表
type MicroAppVersion struct {
	BaseModel
	AppRecordId       uint                          `gorm:"type:int(11);not null;index" json:"appRecordId"`    // 微应用ID
	Version           string                        `gorm:"type:varchar(20);not null" json:"version"`          // 版本号（如 1.0.0）
	VersionCode       int                           `gorm:"type:int(11);not null" json:"versionCode"`          // 版本号（数字）
	PackageUrl        string                        `gorm:"type:varchar(500);not null" json:"packageUrl"`      // 应用包下载地址
	PackageSrc        string                        `gorm:"type:varchar(500);not null" json:"packageSrc"`      // 应用包源
	PackageHash       string                        `gorm:"type:varchar(100);not null" json:"packageHash"`     // 版本校验值（MD5/SHA）
	IconUrl           string                        `gorm:"type:varchar(2000)" json:"iconUrl"`                 // 图标URL
	VersionDesc       datatype.VersionDesc          `gorm:"type:text" json:"versionDesc"`                      // 版本说明（多语言格式）
	Config            *MicroAppVersionConfig        `gorm:"type:json;serializer:json" json:"config"`           // 完整配置信息（JSON）
	AppJson           datatype.AppJson               `gorm:"type:text" json:"appJson"`                          // 原始 app.json 完整内容
	Status            int                           `gorm:"type:tinyint(2);not null;default:-1" json:"status"` // 审核状态：-1-草稿 0-待审核 1-通过 2-拒绝
	CodeSecurityAudit *datatype.SecurityAuditReport `gorm:"type:text" json:"codeSecurityAudit"`                // 代码安全审核报告
	ReviewTime        *time.Time                    `gorm:"type:datetime" json:"reviewTime"`                   // 审核时间
	ReviewerId        uint                          `gorm:"type:int(11)" json:"reviewerId"`                    // 审核人ID
	ReviewNote        string                        `gorm:"type:varchar(500)" json:"reviewNote"`               // 审核备注
	// 下架相关字段
	OfflineType   int       `gorm:"type:tinyint(1);not null;default:0" json:"offlineType"` // 下架类型：0-正常 1-作者下架 2-平台下架
	OfflineReason string    `gorm:"type:varchar(500)" json:"offlineReason"`                // 下架原因
	MicroApp      *MicroApp `gorm:"foreignKey:AppRecordId" json:"microApp"`
}

// 表名
func (MicroAppVersion) TableName() string {
	return "micro_app_version"
}

// 获取版本列表（支持分页和筛选）
func (m *MicroAppVersion) GetList(db *gorm.DB, page, limit int, appRecordId *uint, status *int) ([]MicroAppVersion, int64, error) {
	var list []MicroAppVersion
	var total int64

	query := db.Model(&MicroAppVersion{})

	// 应用ID筛选
	if appRecordId != nil {
		query = query.Where("app_record_id = ?", *appRecordId)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("created_at DESC, id DESC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 根据ID获取版本详情
func (m *MicroAppVersion) GetById(db *gorm.DB, id uint) (MicroAppVersion, error) {
	var version MicroAppVersion
	err := db.Where("id = ?", id).First(&version).Error
	return version, err
}

// 根据版本号获取版本详情
func (m *MicroAppVersion) GetByVersion(db *gorm.DB, version string) (MicroAppVersion, error) {
	var versionInfo MicroAppVersion
	err := db.Order("id DESC").Where("version = ?", version).First(&versionInfo).Error
	return versionInfo, err
}

// 获取应用的最新版本
func (m *MicroAppVersion) GetLatestOnlineByAppId(db *gorm.DB, appId uint) (MicroAppVersion, error) {
	var version MicroAppVersion
	err := db.Where("app_record_id = ? AND status = ? AND offline_type = ?", appId, 1, 0). // 只获取已通过审核的版本并且未下架的版本
												Order("created_at DESC").
												First(&version).Error
	return version, err
}

// 获取应用的所有已发布版本
func (m *MicroAppVersion) GetPublishedVersions(db *gorm.DB, appId uint) ([]MicroAppVersion, error) {
	var versions []MicroAppVersion
	err := db.Where("app_record_id = ? AND status = ?", appId, 1).
		Order("version_code DESC").
		Find(&versions).Error
	return versions, err
}

// 创建版本
func (m *MicroAppVersion) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新版本
func (m *MicroAppVersion) Update(db *gorm.DB) error {
	return db.Model(&MicroAppVersion{}).Where("id = ?", m.ID).Select("*").Updates(m).Error
}

// 删除版本
func (m *MicroAppVersion) Delete(db *gorm.DB, ids []uint) error {
	return db.Where("id IN ?", ids).Delete(&MicroAppVersion{}).Error
}

// 审核版本
func (m *MicroAppVersion) Review(db *gorm.DB, status int, reviewerId uint, reviewNote string) error {
	m.Status = status
	m.ReviewerId = reviewerId
	m.ReviewNote = reviewNote
	now := time.Now()
	m.ReviewTime = &now
	return db.Model(&MicroAppVersion{}).Where("id = ?", m.ID).Select("status", "reviewer_id", "review_note", "review_time").Updates(m).Error
}

// 下架版本
func (m *MicroAppVersion) Offline(db *gorm.DB, id uint, offlineType int, reason string) error {
	return db.Model(&MicroAppVersion{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         3, // 3-已下架
		"offline_type":   offlineType,
		"offline_reason": reason,
	}).Error
}

// 获取待审核的版本列表
func (m *MicroAppVersion) GetPendingList(db *gorm.DB, page, limit int) ([]MicroAppVersion, int64, error) {
	var list []MicroAppVersion
	var total int64

	query := db.Model(&MicroAppVersion{}).Where("status = ?", 0) // 待审核

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("id ASC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 检查版本号是否存在
func (m *MicroAppVersion) CheckVersionExist(db *gorm.DB, appRecordId uint, version string, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&MicroAppVersion{}).Where("app_record_id = ? AND version = ?", appRecordId, version)

	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// 检查版本号数字是否存在
func (m *MicroAppVersion) CheckVersionCodeExist(db *gorm.DB, appId uint, versionCode int, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&MicroAppVersion{}).Where("app_record_id = ? AND version_code = ?", appId, versionCode)

	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
