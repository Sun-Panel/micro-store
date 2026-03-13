package models

import (
	"gorm.io/gorm"
	"time"
)

// 微应用版本列表
type MicroAppVersion struct {
	BaseModel
	AppId       uint      `gorm:"type:int(11);not null;index" json:"appId"`         // 微应用ID
	Version     string    `gorm:"type:varchar(20);not null" json:"version"`         // 版本号（如 1.0.0）
	VersionCode int       `gorm:"type:int(11);not null" json:"versionCode"`         // 版本号（数字）
	PackageUrl  string    `gorm:"type:varchar(500);not null" json:"packageUrl"`     // 应用包下载地址
	PackageHash string    `gorm:"type:varchar(100);not null" json:"packageHash"`    // 版本校验值（MD5/SHA）
	Status      int       `gorm:"type:tinyint(1);not null;default:0" json:"status"` // 审核状态：0-待审核 1-通过 2-拒绝
	ReviewTime  time.Time `gorm:"type:datetime" json:"reviewTime"`                  // 审核时间
	ReviewerId  uint      `gorm:"type:int(11)" json:"reviewerId"`                   // 审核人ID
	ReviewNote  string    `gorm:"type:varchar(500)" json:"reviewNote"`              // 审核备注
}

// 表名
func (MicroAppVersion) TableName() string {
	return "micro_app_version"
}

// 获取版本列表（支持分页和筛选）
func (m *MicroAppVersion) GetList(db *gorm.DB, page, limit int, appId *uint, status *int) ([]MicroAppVersion, int64, error) {
	var list []MicroAppVersion
	var total int64

	query := db.Model(&MicroAppVersion{})

	// 应用ID筛选
	if appId != nil {
		query = query.Where("app_id = ?", *appId)
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
	err = query.Order("version_code DESC, id DESC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 根据ID获取版本详情
func (m *MicroAppVersion) GetById(db *gorm.DB, id uint) (MicroAppVersion, error) {
	var version MicroAppVersion
	err := db.Where("id = ?", id).First(&version).Error
	return version, err
}

// 获取应用的最新版本
func (m *MicroAppVersion) GetLatestByAppId(db *gorm.DB, appId uint) (MicroAppVersion, error) {
	var version MicroAppVersion
	err := db.Where("app_id = ? AND status = ?", appId, 1). // 只获取已通过审核的版本
		Order("version_code DESC").
		First(&version).Error
	return version, err
}

// 获取应用的所有已发布版本
func (m *MicroAppVersion) GetPublishedVersions(db *gorm.DB, appId uint) ([]MicroAppVersion, error) {
	var versions []MicroAppVersion
	err := db.Where("app_id = ? AND status = ?", appId, 1).
		Order("version_code DESC").
		Find(&versions).Error
	return versions, err
}

// 创建版本
func (m *MicroAppVersion) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新版本
func (m *MicroAppVersion) Update(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	return db.Model(&MicroAppVersion{}).Where("id = ?", id).Updates(updateData).Error
}

// 删除版本
func (m *MicroAppVersion) Delete(db *gorm.DB, ids []uint) error {
	return db.Where("id IN ?", ids).Delete(&MicroAppVersion{}).Error
}

// 审核版本
func (m *MicroAppVersion) Review(db *gorm.DB, id uint, status int, reviewerId uint, reviewNote string) error {
	return db.Model(&MicroAppVersion{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerId,
		"review_note": reviewNote,
		"review_time": time.Now(),
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
func (m *MicroAppVersion) CheckVersionExist(db *gorm.DB, appId uint, version string, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&MicroAppVersion{}).Where("app_id = ? AND version = ?", appId, version)

	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// 检查版本号数字是否存在
func (m *MicroAppVersion) CheckVersionCodeExist(db *gorm.DB, appId uint, versionCode int, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&MicroAppVersion{}).Where("app_id = ? AND version_code = ?", appId, versionCode)

	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
