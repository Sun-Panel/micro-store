package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type VersionType string

const (
	VersionTypeRelease VersionType = "release"
	VersionTypeBeta    VersionType = "beta"
	VersionTypeAlpha   VersionType = "alpha"
	VersionTypeRC      VersionType = "rc"
	VersionTypeDev     VersionType = "dev"
)

type Version struct {
	BaseModel
	ID           uint        `gorm:"primary_key" json:"id"`
	Version      string      `gorm:"type:varchar(100);not null" json:"version"`
	Type         VersionType `gorm:"type:varchar(20);not null" json:"type"`
	ReleaseTime  time.Time   `gorm:"not null" json:"-"`
	Description  string      `gorm:"type:text" json:"description"`
	DownloadURL  string      `gorm:"type:varchar(2000)" json:"downloadURL"`
	PageUrl      string      `gorm:"type:varchar(2000)" json:"pageUrl"`
	IsActive     bool        `gorm:"default:false" json:"isActive"`
	IsRolledBack bool        `gorm:"default:false" json:"isRolledBack"`
	VersionSecret VersionSecret `gorm:"foreignKey:Version;references:Version" json:"-"`
}

// VersionRepository 定义了与Version模型相关的数据库操作
type VersionRepository struct {
	DB *gorm.DB
}

// NewVersionRepository 创建一个新的版本仓库
func NewVersionRepository(db *gorm.DB) *VersionRepository {
	return &VersionRepository{DB: db}
}

// FindNewestVersionsAfterDate 根据指定日期查询最新版本
func (repo *VersionRepository) FindNewestVersionsAfterCurrentVersion(currentVersion string, versionType VersionType) (Version, error) {
	// ReleaseTime 用于临时存储查询到的发布时间
	type ReleaseTime struct {
		ReleaseTime time.Time
	}
	// 首先，根据currentVersion查询对应的发布时间
	var current ReleaseTime
	var version Version

	result := repo.DB.Model(&Version{}).Where("version = ?", currentVersion).Select("release_time").Take(&current)
	if result.Error != nil {
		return version, result.Error
	}

	// 如果没有找到当前版本，返回错误
	if current.ReleaseTime.IsZero() {
		return version, errors.New("current version not found")
	}

	// 使用找到的发布时间作为since参数来查询后续的版本

	query := repo.DB.Model(&Version{}).Where("release_time > ?", current.ReleaseTime)

	// 如果提供了版本类型，则添加到查询条件中
	if versionType != "" {
		query = query.Where("type = ?", versionType)
	}

	// 执行查询，按照发布时间降序排列
	result = query.Order("release_time desc").First(&version)
	if result.Error != nil {
		return version, result.Error
	}

	return version, nil
}
