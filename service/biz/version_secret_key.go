package biz

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"
)

type versionSecret struct {
	SecretKey cache.Cacher[models.VersionSecret]
}

func (v *versionSecret) Init() {
	v.SecretKey = global.NewCache[models.VersionSecret](5*time.Minute, 10*time.Hour, "version_secret_cache")
}

func (v *versionSecret) GetVersionSecret(version string) (models.VersionSecret, error) {
	record, exist := v.SecretKey.Get(version)
	if exist {
		return record, nil
	}
	record, err := v.GetVersionSecretFromDb(version)
	if err != nil {
		return record, err
	}
	v.SecretKey.SetDefault(version, record)
	return record, err
}

func (v *versionSecret) GetVersionSecretFromDb(version string) (models.VersionSecret, error) {
	mvs := models.VersionSecret{}
	err := mvs.Find(global.Db, version)
	return mvs, err
}

func (v *versionSecret) GetVersionSecretByIdFromDb(id uint) (models.VersionSecret, error) {
	mvs := models.VersionSecret{}
	err := mvs.FindById(global.Db, id)
	return mvs, err
}

// 创建
func (v *versionSecret) CreateVersionSecret(versionSecret models.VersionSecret) error {
	v.SecretKey.Delete(versionSecret.Version)
	return versionSecret.Create(global.Db)
}

// 更新
func (v *versionSecret) UpdateVersionSecret(versionSecret models.VersionSecret) error {
	v.SecretKey.Delete(versionSecret.Version)
	return versionSecret.Update(global.Db, "SecretKey", "Status")
}

// 删除
func (v *versionSecret) DeleteVersionSecret(versionSecret models.VersionSecret) error {
	v.SecretKey.Delete(versionSecret.Version)
	return versionSecret.Delete(global.Db)
}
