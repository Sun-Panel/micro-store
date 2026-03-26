package models

import (
	"sun-panel/models/datatype"
	"time"

	"gorm.io/gorm"
)

// 微应用审核快照表
type MicroAppReview struct {
	BaseModel
	MicroAppBaseInfo `gorm:"embedded"` // 嵌入公共字段
	MicroAppId       string            `gorm:"column:micro_app_id;type:varchar(120);not null;index" json:"microAppId"` // 覆盖嵌入字段
	AppRecordId      uint              `gorm:"type:int(11);not null;index" json:"appRecordId"`                         // 关联 micro_app.id（可选）
	LangMap          datatype.MapJson  `gorm:"type:json" json:"langMap"`                                               // 多语言信息JSON

	// 审核信息
	Status     int        `gorm:"type:tinyint(2);not null;default:0;index" json:"status"` // 审核状态：-1-草稿 0-待审核 1-已通过 2-已拒绝
	ReviewerId uint       `gorm:"type:int(11)" json:"reviewerId"`                         // 审核人ID，userID
	ReviewNote string     `gorm:"type:varchar(500)" json:"reviewNote"`                    // 审核备注
	ReviewTime *time.Time `gorm:"type:datetime" json:"reviewTime"`                        // 审核时间
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

// 根据应用自增ID获取最新记录
func (m *MicroAppReview) GetLatestByAppRecordId(db *gorm.DB, appRecordId uint) (MicroAppReview, error) {
	var review MicroAppReview
	err := db.Where("app_record_id = ?", appRecordId).Order("created_at DESC").First(&review).Error
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
	err = query.Order("updated_at ASC").Offset(offset).Limit(limitSize).Find(&list).Error

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

// 根据业务ID获取审核记录列表
func (m *MicroAppReview) GetListByMicroAppId(db *gorm.DB, microAppId string, page, limit int) ([]MicroAppReview, int64, error) {
	var list []MicroAppReview
	var total int64

	query := db.Model(&MicroAppReview{}).Where("micro_app_id = ?", microAppId)

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
