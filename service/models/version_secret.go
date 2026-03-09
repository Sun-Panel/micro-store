package models

import "gorm.io/gorm"

type VersionSecret struct {
	BaseModel
	Version   string `gorm:"type:varchar(100);not null" json:"version"`
	SecretKey string `gorm:"type:text;not null" json:"secretKey"`
	Status    bool   `gorm:"type:tinyint;not null" json:"status"` // 是否启用
	// Encrypted bool   `gorm:"type:boolean;not null" json:"encrypted"` // 秘钥是否加密储存了
}

func (v *VersionSecret) Create(db *gorm.DB) error {
	return db.Create(v).Error
}

func (v *VersionSecret) Update(db *gorm.DB, fields ...string) error {
	return db.Select(fields).
		Where("id = ?", v.ID).
		Updates(v).Error
}

func (v *VersionSecret) Delete(db *gorm.DB) error {
	return db.Where("id = ?", v.ID).Delete(v).Error
}

func (v *VersionSecret) Find(db *gorm.DB, version string) error {
	return db.Where("version = ?", version).First(v).Error
}

func (v *VersionSecret) FindById(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).First(v).Error
}
