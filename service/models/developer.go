package models

import (
	"slices"
	"time"

	"gorm.io/gorm"
)

// 业务错误码定义
const (
	ErrCodeDeveloperNameExists   = "E_DEVELOPER_NAME_EXISTS"
	ErrCodeDeveloperNameCooldown = "E_DEVELOPER_NAME_COOLDOWN"
)

// 开发者表
type Developer struct {
	BaseModel
	UserId        uint       `gorm:"type:int(11);not null;uniqueIndex" json:"userId"`            // 用户ID
	DeveloperName string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"developerName"` // 开发者标识（纯英文，多词用-分割）
	Name          string     `gorm:"type:varchar(50)" json:"name"`                               // 开发者名称
	ContactMail   string     `gorm:"type:varchar(50)" json:"contactMail"`                        // 联系邮箱
	PaymentName   string     `gorm:"type:varchar(50)" json:"paymentName"`                        // 收款人真实姓名
	PaymentQrcode string     `gorm:"type:varchar(500)" json:"paymentQrcode"`                     // 收款二维码图片URL
	PaymentMethod string     `gorm:"type:varchar(200)" json:"paymentMethod"`                     // 收款方式描述
	Status        int        `gorm:"type:tinyint(1);default:1" json:"status"`                    // 状态：0-禁用 1-正常
	NameUpdatedAt *time.Time `gorm:"type:timestamp" json:"nameUpdatedAt"`                        // Name 上次修改时间
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

// 根据开发者标识获取开发者信息
func (m *Developer) GetByDeveloperName(db *gorm.DB, developerName string) (Developer, error) {
	var developer Developer
	err := db.Where("developer_name = ?", developerName).First(&developer).Error
	return developer, err
}

// 创建开发者
func (m *Developer) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新开发者
// field 为可选参数，指定要更新的字段名列表
func (m *Developer) Update(db *gorm.DB, id uint, updateData Developer, field ...string) error {
	// 如果提供了开发者名称，需要校验唯一性
	// if len(field) == 0 || containsField(field, "developer_name") {
	// 	if exist, err := m.CheckNameExist(db, updateData.DeveloperName, id); err != nil {
	// 		return err
	// 	} else if exist {
	// 		return gorm.ErrRegistered
	// 	}
	// }

	return db.Model(m).Where("id = ?", id).Select(field).Updates(updateData).Error
}

// containsField 检查字段列表是否包含指定字段
func containsField(fields []string, name string) bool {
	return slices.Contains(fields, name)
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

// ========== 业务逻辑方法 ==========

// Register 注册成为开发者（带校验）
func (m *Developer) Register(db *gorm.DB, userId uint, developerName, contactMail, paymentName, paymentQrcode, paymentMethod, name string) (uint, error) {
	// 检查用户是否已经是开发者
	if isDeveloper, err := m.CheckUserIsDeveloper(db, userId); err != nil {
		return 0, err
	} else if isDeveloper {
		return 0, gorm.ErrRegistered
	}

	// 检查开发者名称是否存在
	if exist, err := m.CheckNameExist(db, developerName, 0); err != nil {
		return 0, err
	} else if exist {
		return 0, gorm.ErrRegistered
	}

	m.UserId = userId
	m.DeveloperName = developerName
	m.ContactMail = contactMail
	m.PaymentName = paymentName
	m.PaymentQrcode = paymentQrcode
	m.PaymentMethod = paymentMethod
	m.Status = 1 // 默认启用
	m.Name = name

	if err := m.Create(db); err != nil {
		return 0, err
	}

	return m.ID, nil
}

// UpdateFields 需要更新的字段列表（指针类型）
// nil = 不更新该字段，非 nil = 更新为指定值（包括空字符串）
type DeveloperUpdateFields struct {
	DeveloperName *string // 开发者标识
	ContactMail   *string // 联系邮箱
	PaymentName   *string // 收款人真实姓名
	PaymentQrcode *string // 收款二维码图片URL
	PaymentMethod *string // 收款方式描述
	Name          *string // 开发者名称
	Status        *int
}

// UpdateInfo 更新开发者信息（仅数据库操作，业务校验由 biz 层负责）
func (m *Developer) UpdateInfo(db *gorm.DB, id uint, updateFields DeveloperUpdateFields) error {
	// 如果提供了开发者名称，需要校验唯一性
	if updateFields.DeveloperName != nil {
		if exist, err := m.CheckNameExist(db, *updateFields.DeveloperName, id); err != nil {
			return err
		} else if exist {
			return NewModelError("E_DEVELOPER_NAME_EXISTS")
		}
	}

	updateData := map[string]interface{}{}
	if updateFields.DeveloperName != nil {
		updateData["developer_name"] = *updateFields.DeveloperName
	}
	if updateFields.ContactMail != nil {
		updateData["contact_mail"] = *updateFields.ContactMail
	}
	if updateFields.PaymentName != nil {
		updateData["payment_name"] = *updateFields.PaymentName
	}
	if updateFields.PaymentQrcode != nil {
		updateData["payment_qrcode"] = *updateFields.PaymentQrcode
	}
	if updateFields.PaymentMethod != nil {
		updateData["payment_method"] = *updateFields.PaymentMethod
	}
	if updateFields.Name != nil {
		updateData["name"] = *updateFields.Name
		now := time.Now()
		updateData["name_updated_at"] = &now
	}
	if updateFields.Status != nil {
		updateData["status"] = *updateFields.Status
	}

	if len(updateData) == 0 {
		return nil // 没有需要更新的字段
	}

	return db.Model(&Developer{}).Where("id = ?", id).Updates(updateData).Error
}
