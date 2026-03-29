package models

import (
	"time"

	"gorm.io/gorm"
)

// 应用安装记录表
type MicroAppInstall struct {
	BaseModel
	AppRecordId     uint      `gorm:"type:int(11);not null;index" json:"appRecordId"`                      // 微应用ID
	VersionId       uint      `gorm:"type:int(11);not null;index" json:"versionId"`                        // 版本ID
	UserId          uint      `gorm:"type:int(11);index" json:"userId"`                                    // 用户ID
	ClientId        string    `gorm:"type:varchar(100);not null;index" json:"clientId"`                    // 客户端标识
	IntranetIp      string    `gorm:"type:varchar(50)" json:"intranetIp"`                                  // 内网IP
	PublicIp        string    `gorm:"type:varchar(50);not null" json:"publicIp"`                           // 公网IP
	InstallTime     time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"installTime"` // 安装时间
	UserIsPro       bool      `gorm:"type:tinyint(1);default:0" json:"userIsPro"`                          // 安装时用户是否为PRO
	PointValue      int       `gorm:"type:int(11)" json:"pointValue"`                                      // 本次积分值
	AuthorPointRule string    `gorm:"type:varchar(500)" json:"authorPointRule"`                            // 作者当前积分规则JSON
}

// 表名
func (MicroAppInstall) TableName() string {
	return "micro_app_install"
}

// 获取安装记录列表（支持分页和筛选）
func (m *MicroAppInstall) GetList(db *gorm.DB, page, limit int, appId *uint, versionId *uint, userId *uint, clientId *string) ([]MicroAppInstall, int64, error) {
	var list []MicroAppInstall
	var total int64

	query := db.Model(&MicroAppInstall{})

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

	// 客户端标识筛选
	if clientId != nil {
		query = query.Where("client_id = ?", *clientId)
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

// 创建安装记录
func (m *MicroAppInstall) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 获取应用的安装次数
func (m *MicroAppInstall) GetInstallCountByAppId(db *gorm.DB, appId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroAppInstall{}).Where("app_record_id = ?", appId).Count(&count).Error
	return count, err
}

// 获取版本的安装次数
func (m *MicroAppInstall) GetInstallCountByVersionId(db *gorm.DB, versionId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroAppInstall{}).Where("version_id = ?", versionId).Count(&count).Error
	return count, err
}

// 获取用户安装过的应用列表
func (m *MicroAppInstall) GetInstalledAppsByUserId(db *gorm.DB, userId uint) ([]MicroAppInstall, error) {
	var installs []MicroAppInstall
	err := db.Where("user_id = ?", userId).
		Order("install_time DESC").
		Find(&installs).Error
	return installs, err
}

// 获取客户端安装过的应用列表
func (m *MicroAppInstall) GetInstalledAppsByClientId(db *gorm.DB, clientId string) ([]MicroAppInstall, error) {
	var installs []MicroAppInstall
	err := db.Where("client_id = ?", clientId).
		Order("install_time DESC").
		Find(&installs).Error
	return installs, err
}

// 检查用户是否安装过该应用
func (m *MicroAppInstall) CheckUserInstalled(db *gorm.DB, userId uint, appId uint) (bool, error) {
	var count int64
	err := db.Model(&MicroAppInstall{}).Where("user_id = ? AND app_record_id = ?", userId, appId).Count(&count).Error
	return count > 0, err
}

// 获取应用在某客户端的安装记录
func (m *MicroAppInstall) GetInstallByClientIdAndAppId(db *gorm.DB, clientId string, appId uint) (MicroAppInstall, error) {
	var install MicroAppInstall
	err := db.Where("client_id = ? AND app_record_id = ?", clientId, appId).
		Order("install_time DESC").
		First(&install).Error
	return install, err
}
