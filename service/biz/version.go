package biz

import (
	"sun-panel/global"
	"sun-panel/models"
	"time"
)

type VersionType struct {
}

// 获取当前版本的信息
func (s *VersionType) GetCurrentVersion(version string) (latestVersion models.Version, err error) {
	err = global.Db.Where("version=?", version).Order("release_time desc").First(&latestVersion).Error
	return
}

// 获取指定时间后的最新（活跃）版本
func (s *VersionType) GetLatestByAfterTime(afterTime time.Time, versionType models.VersionType) (latestVersion models.Version, err error) {
	err = global.Db.Where("is_active = ? AND release_time > ? AND type=?", true, afterTime, versionType).
		Order("release_time desc").
		First(&latestVersion).Error
	return
}

// 获取最新活跃的版本
func (s *VersionType) GetLatest(versionType models.VersionType) (latestVersion models.Version, err error) {
	err = global.Db.Where("is_active = ? AND type=?", true, versionType).
		Order("release_time desc").
		First(&latestVersion).Error
	return
}
