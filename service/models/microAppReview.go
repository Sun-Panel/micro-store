package models

import (
	"gorm.io/gorm"
	"time"
)

// 微应用审核快照表
type MicroAppReview struct {
	BaseModel
	AppId uint `gorm:"type:int(11);not null;index" json:"appId"` // 微应用ID

	// 审核数据快照
	AppName     string  `gorm:"type:varchar(100);not null" json:"appName"`     // 应用名称
	AppIcon     string  `gorm:"type:varchar(200);not null" json:"appIcon"`     // 应用图标URL
	AppDesc     string  `gorm:"type:varchar(500)" json:"appDesc"`              // 应用简介
	CategoryId  int     `gorm:"type:int(11);not null;index" json:"categoryId"` // 应用分类ID
	ChargeType  int     `gorm:"type:tinyint(1);not null;default:0" json:"chargeType"` // 收费方式：0-免费 1-积分 2-订阅PRO免费
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`                       // 价格（积分时）
	Screenshots string  `gorm:"type:varchar(2000)" json:"screenshots"`                 // 图集（多个图片URL用逗号分隔）
	LangMap     string  `gorm:"type:json" json:"langMap"`                              // 多语言信息JSON
	Remark      string  `gorm:"type:varchar(500)" json:"remark"`                       // 应用备注

	// 审核信息
	Status     int        `gorm:"type:tinyint(2);not null;default:0;index" json:"status"` // 审核状态：0-待审核 1-已通过 2-已拒绝
	ReviewerId uint       `gorm:"type:int(11)" json:"reviewerId"`                          // 审核人ID
	ReviewNote string     `gorm:"type:varchar(500)" json:"reviewNote"`                     // 审核备注
	ReviewTime *time.Time `gorm:"type:datetime" json:"reviewTime"`                         // 审核时间
}

// 表名
func (MicroAppReview) TableName() string {
	return "micro_app_review"
}

// 创建审核记录
func (m *MicroAppReview) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 根据ID获取审核记录
func (m *MicroAppReview) GetById(db *gorm.DB, id uint) (MicroAppReview, error) {
	var review MicroAppReview
	err := db.Where("id = ?", id).First(&review).Error
	return review, err
}

// 获取应用的待审核记录
func (m *MicroAppReview) GetPendingByAppId(db *gorm.DB, appId uint) (MicroAppReview, error) {
	var review MicroAppReview
	err := db.Where("app_id = ? AND status = ?", appId, 0).First(&review).Error
	return review, err
}

// 获取应用的审核历史记录
func (m *MicroAppReview) GetListByAppId(db *gorm.DB, appId uint, page, limit int) ([]MicroAppReview, int64, error) {
	var list []MicroAppReview
	var total int64

	query := db.Model(&MicroAppReview{}).Where("app_id = ?", appId)

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

// 获取待审核列表（管理员用）
func (m *MicroAppReview) GetPendingList(db *gorm.DB, page, limit int) ([]MicroAppReview, int64, error) {
	var list []MicroAppReview
	var total int64

	query := db.Model(&MicroAppReview{}).Where("status = ?", 0)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("id ASC").Offset(offset).Limit(limitSize).Find(&list).Error

	return list, total, err
}

// 更新审核状态
func (m *MicroAppReview) UpdateStatus(db *gorm.DB, id uint, status int, reviewerId uint, reviewNote string) error {
	now := time.Now()
	return db.Model(&MicroAppReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerId,
		"review_note": reviewNote,
		"review_time": now,
	}).Error
}

// 删除应用的审核记录
func (m *MicroAppReview) DeleteByAppId(db *gorm.DB, appId uint) error {
	return db.Where("app_id = ?", appId).Delete(&MicroAppReview{}).Error
}
