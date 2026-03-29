package models

import (
	"time"

	"gorm.io/gorm"
)

// 下载记录表
type MicroAppDownload struct {
	BaseModel
	AppRecordId    uint      `gorm:"type:int(11);not null;index" json:"appRecordId"`                       // 微应用ID
	VersionId      uint      `gorm:"type:int(11);not null;index" json:"versionId"`                         // 版本ID
	UserId         uint      `gorm:"type:int(11);index" json:"userId"`                                     // 用户ID（匿名下载为空）
	ClientId       string    `gorm:"type:varchar(100);not null;index" json:"clientId"`                     // 客户端标识
	DownloadIp     string    `gorm:"type:varchar(50);not null" json:"downloadIp"`                          // 下载IP
	DownloadDevice string    `gorm:"type:varchar(200)" json:"downloadDevice"`                              // 下载设备信息
	DownloadClient string    `gorm:"type:varchar(50)" json:"downloadClient"`                               // 下载客户端类型
	DownloadTime   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"downloadTime"` // 下载时间
}

// 表名
func (MicroAppDownload) TableName() string {
	return "micro_app_download"
}

// 获取下载记录列表（支持分页和筛选）
func (m *MicroAppDownload) GetList(db *gorm.DB, page, limit int, appId *uint, versionId *uint, userId *uint) ([]MicroAppDownload, int64, error) {
	var list []MicroAppDownload
	var total int64

	query := db.Model(&MicroAppDownload{})

	// 应用ID筛选
	if appId != nil {
		query = query.Where("app_record_id = ?", *appId)
	}

	// 版本ID筛选
	if versionId != nil {
		query = query.Where("version_id = ?", *versionId)
	}

	// 用户ID筛选
	if userId != nil {
		query = query.Where("user_id = ?", *userId)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("id DESC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 创建下载记录
func (m *MicroAppDownload) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 获取应用的下载次数
func (m *MicroAppDownload) GetDownloadCountByAppId(db *gorm.DB, appId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroAppDownload{}).Where("app_record_id = ?", appId).Count(&count).Error
	return count, err
}

// 获取版本的下载次数
func (m *MicroAppDownload) GetDownloadCountByVersionId(db *gorm.DB, versionId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroAppDownload{}).Where("version_id = ?", versionId).Count(&count).Error
	return count, err
}

// 获取用户下载过的应用ID列表
func (m *MicroAppDownload) GetDownloadedAppIds(db *gorm.DB, userId uint) ([]uint, error) {
	var appIds []uint
	err := db.Model(&MicroAppDownload{}).
		Select("DISTINCT app_record_id").
		Where("user_id = ?", userId).
		Pluck("app_record_id", &appIds).Error
	return appIds, err
}

// 检查用户是否下载过该应用
func (m *MicroAppDownload) CheckUserDownloaded(db *gorm.DB, userId uint, appId uint) (bool, error) {
	var count int64
	err := db.Model(&MicroAppDownload{}).Where("user_id = ? AND app_record_id = ?", userId, appId).Count(&count).Error
	return count > 0, err
}
