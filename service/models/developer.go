package models

import "gorm.io/gorm"

// 开发者表
type Developer struct {
	BaseModel
	UserId        uint   `gorm:"type:int(11);not null;uniqueIndex" json:"userId"`   // 用户ID
	DeveloperName string `gorm:"type:varchar(50);not null" json:"developerName"`    // 开发者名称
	ContactMail   string `gorm:"type:varchar(50)" json:"contactMail"`               // 联系邮箱
	Status        int    `gorm:"type:tinyint(1);default:1" json:"status"`           // 状态：0-禁用 1-正常
}

// 表名
func (Developer) TableName() string {
	return "developer"
}

// 获取开发者列表（支持分页和筛选）
func (m *Developer) GetList(db *gorm.DB, page, limit int, status *int, keyWord string) ([]Developer, int64, error) {
	var list []Developer
	var total int64

	query := db.Model(&Developer{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("developer_name LIKE ? OR contact_mail LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
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

// 根据ID获取开发者详情
func (m *Developer) GetById(db *gorm.DB, id uint) (Developer, error) {
	var developer Developer
	err := db.Where("id = ?", id).First(&developer).Error
	return developer, err
}

// 根据用户ID获取开发者信息
func (m *Developer) GetByUserId(db *gorm.DB, userId uint) (Developer, error) {
	var developer Developer
	err := db.Where("user_id = ?", userId).First(&developer).Error
	return developer, err
}

// 创建开发者
func (m *Developer) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新开发者
func (m *Developer) Update(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	return db.Model(&Developer{}).Where("id = ?", id).Updates(updateData).Error
}

// 删除开发者
func (m *Developer) Delete(db *gorm.DB, ids []uint) error {
	return db.Where("id IN ?", ids).Delete(&Developer{}).Error
}

// 更新开发者状态
func (m *Developer) UpdateStatus(db *gorm.DB, id uint, status int) error {
	return db.Model(&Developer{}).Where("id = ?", id).Update("status", status).Error
}

// 检查开发者名称是否存在
func (m *Developer) CheckNameExist(db *gorm.DB, developerName string, excludeId uint) (bool, error) {
	var count int64
	query := db.Model(&Developer{}).Where("developer_name = ?", developerName)

	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// 检查用户是否已经是开发者
func (m *Developer) CheckUserIsDeveloper(db *gorm.DB, userId uint) (bool, error) {
	var count int64
	err := db.Model(&Developer{}).Where("user_id = ?", userId).Count(&count).Error
	return count > 0, err
}
