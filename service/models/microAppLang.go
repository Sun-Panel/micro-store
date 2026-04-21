package models

import "gorm.io/gorm"

// 微应用多语言信息表
type MicroAppLang struct {
	BaseModel
	MicroAppId string         `gorm:"type:varchar(50);not null" json:"microAppId"` // 关联微应用ID
	Lang       string         `gorm:"type:varchar(10);not null" json:"lang"`       // 语言代码：zh-CN, en-US, ja-JP 等
	AppName    string         `gorm:"type:varchar(100);not null" json:"appName"`   // 应用名称
	AppDesc    string         `gorm:"type:varchar(500)" json:"appDesc"`            // 应用简介
	DeletedAt  gorm.DeletedAt `gorm:"uniqueIndex:idx_app_lang" json:"deletedAt,omitempty"`
}

// 复合唯一索引：micro_app_id + lang + deleted_at
// 这样软删除后不会与新增记录冲突

// 表名
func (MicroAppLang) TableName() string {
	return "micro_app_lang"
}

// 获取微应用的多语言信息列表
func (m *MicroAppLang) GetListByAppId(db *gorm.DB, microAppId string) ([]MicroAppLang, error) {
	var list []MicroAppLang
	err := db.Where("micro_app_id = ?", microAppId).Find(&list).Error
	return list, err
}

// 根据微应用ID和语言获取信息
func (m *MicroAppLang) GetByAppIdAndLang(db *gorm.DB, microAppId string, lang string) (MicroAppLang, error) {
	var langInfo MicroAppLang
	err := db.Where("micro_app_id = ? AND lang = ?", microAppId, lang).First(&langInfo).Error
	return langInfo, err
}

// 创建或更新多语言信息
func (m *MicroAppLang) CreateOrUpdate(db *gorm.DB) error {
	// 先查询是否存在
	var existing MicroAppLang
	err := db.Where("micro_app_id = ? AND lang = ?", m.MicroAppId, m.Lang).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		return db.Create(m).Error
	} else if err != nil {
		// 其他错误
		return err
	}

	// 存在，更新记录
	return db.Model(&existing).Updates(map[string]interface{}{
		"app_name": m.AppName,
		"app_desc": m.AppDesc,
	}).Error
}

// 批量创建多语言信息
func (m *MicroAppLang) BatchCreate(db *gorm.DB, langList []MicroAppLang) error {
	return db.Create(&langList).Error
}

// 更新多语言信息
func (m *MicroAppLang) Update(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	return db.Model(&MicroAppLang{}).Where("id = ?", id).Updates(updateData).Error
}

// 删除某个微应用的所有多语言信息
func (m *MicroAppLang) DeleteByAppId(db *gorm.DB, microAppId string) error {
	return db.Where("micro_app_id = ?", microAppId).Delete(&MicroAppLang{}).Error
}

// 删除指定语言信息
func (m *MicroAppLang) DeleteByAppIdAndLang(db *gorm.DB, microAppId string, lang string) error {
	return db.Where("micro_app_id = ? AND lang = ?", microAppId, lang).Delete(&MicroAppLang{}).Error
}

// 检查语言是否存在
func (m *MicroAppLang) CheckLangExist(db *gorm.DB, microAppId string, lang string) (bool, error) {
	var count int64
	err := db.Model(&MicroAppLang{}).Where("micro_app_id = ? AND lang = ?", microAppId, lang).Count(&count).Error
	return count > 0, err
}

// 获取某个微应用支持的语言列表
func (m *MicroAppLang) GetSupportedLangs(db *gorm.DB, microAppId string) ([]string, error) {
	var langs []string
	err := db.Model(&MicroAppLang{}).
		Where("micro_app_id = ?", microAppId).
		Pluck("lang", &langs).Error
	return langs, err
}
