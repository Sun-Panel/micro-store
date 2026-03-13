package models

import "gorm.io/gorm"

// 应用分类表
type MicroAppCategory struct {
	BaseModel
	Name   string `gorm:"type:varchar(50);not null" json:"name"`   // 分类名称
	Icon   string `gorm:"type:varchar(200)" json:"icon"`           // 分类图标
	Sort   int    `gorm:"type:int(11);default:0" json:"sort"`      // 排序（数字越大越靠前）
	Status int    `gorm:"type:tinyint(1);default:1" json:"status"` // 状态：0-禁用 1-正常
}

// 表名
func (MicroAppCategory) TableName() string {
	return "micro_app_category"
}

// 获取分类列表（支持分页和筛选）
func (m *MicroAppCategory) GetList(db *gorm.DB, page, limit int, status *int, keyWord string) ([]MicroAppCategory, int64, error) {
	var list []MicroAppCategory
	var total int64

	query := db.Model(&MicroAppCategory{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("name LIKE ?", "%"+keyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("sort DESC, id DESC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 获取所有启用的分类
func (m *MicroAppCategory) GetEnabledList(db *gorm.DB) ([]MicroAppCategory, error) {
	var list []MicroAppCategory
	err := db.Where("status = ?", 1).Order("sort DESC, id ASC").Find(&list).Error
	return list, err
}

// 根据ID获取分类详情
func (m *MicroAppCategory) GetById(db *gorm.DB, id uint) (MicroAppCategory, error) {
	var category MicroAppCategory
	err := db.Where("id = ?", id).First(&category).Error
	return category, err
}

// 创建分类
func (m *MicroAppCategory) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新分类
func (m *MicroAppCategory) Update(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	return db.Model(&MicroAppCategory{}).Where("id = ?", id).Updates(updateData).Error
}

// 删除分类（支持批量删除）
func (m *MicroAppCategory) Delete(db *gorm.DB, ids []uint) error {
	return db.Where("id IN ?", ids).Delete(&MicroAppCategory{}).Error
}

// 检查分类名称是否存在
func (m *MicroAppCategory) CheckNameExist(db *gorm.DB, name string, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&MicroAppCategory{}).Where("name = ?", name)

	// 排除指定ID（用于更新时检查）
	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// 更新分类状态
func (m *MicroAppCategory) UpdateStatus(db *gorm.DB, id uint, status int) error {
	return db.Model(&MicroAppCategory{}).Where("id = ?", id).Update("status", status).Error
}
