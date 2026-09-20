package models

import (
	"time"

	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

// CustomCodeReview 自定义代码审核快照表（每次提交的内容快照 + 审核记录）
type CustomCodeReview struct {
	BaseModel
	CustomCodeBaseInfo `gorm:"embedded"`

	CustomCodeId uint                             `gorm:"type:int(11);not null;index" json:"customCodeId"` // 关联 custom_code.id
	Status       int                              `gorm:"type:tinyint(2);not null;default:0;index" json:"status"`
	ReviewerId   uint                             `gorm:"type:int(11)" json:"reviewerId"`      // 审核人userID
	ReviewNote   string                           `gorm:"type:varchar(500)" json:"reviewNote"` // 审核备注（对作者可见）
	ReviewTime   *time.Time                       `gorm:"type:datetime" json:"reviewTime"`     // 审核时间
	MachineAudit *datatype.CustomCodeMachineAudit `gorm:"type:text" json:"machineAudit"`       // 机器预审结果（含未通过原因，供人工审核参考）
}

// 表名
func (CustomCodeReview) TableName() string {
	return "custom_code_review"
}

// 创建
func (m *CustomCodeReview) Create(db *gorm.DB, item *CustomCodeReview) error {
	return db.Create(item).Error
}

// 根据ID获取
func (m *CustomCodeReview) GetById(db *gorm.DB, id uint) (CustomCodeReview, error) {
	var info CustomCodeReview
	err := db.Where("id = ?", id).First(&info).Error
	return info, err
}

// 获取某帖子最新的快照
func (m *CustomCodeReview) GetLatestByCustomCodeId(db *gorm.DB, customCodeId uint) (CustomCodeReview, error) {
	var info CustomCodeReview
	err := db.Where("custom_code_id = ?", customCodeId).Order("id DESC").First(&info).Error
	return info, err
}

// 获取某帖子待审核的快照
func (m *CustomCodeReview) GetPendingByCustomCodeId(db *gorm.DB, customCodeId uint) (CustomCodeReview, error) {
	var info CustomCodeReview
	err := db.Where("custom_code_id = ? AND status = ?", customCodeId, CustomCodeReviewStatusPending).
		Order("id DESC").First(&info).Error
	return info, err
}

// 获取待审核列表（审核员）
func (m *CustomCodeReview) GetPendingList(db *gorm.DB, page, limit int, keyword string) ([]CustomCodeReview, int64, error) {
	page, limit = normalizePage(page, limit)

	query := db.Model(&CustomCodeReview{}).Where("status = ?", CustomCodeReviewStatusPending)
	if keyword != "" {
		kw := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", kw, kw)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []CustomCodeReview
	err := query.Select("id, custom_code_id, title, description, keywords, is_original, source_note, code_types, status, reviewer_id, review_note, review_time, machine_audit, created_at, updated_at").
		Order("id ASC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&list).Error

	return list, total, err
}

// 获取某帖子的审核历史
func (m *CustomCodeReview) GetHistoryByCustomCodeId(db *gorm.DB, customCodeId uint, page, limit int) ([]CustomCodeReview, int64, error) {
	page, limit = normalizePage(page, limit)

	query := db.Model(&CustomCodeReview{}).Where("custom_code_id = ?", customCodeId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []CustomCodeReview
	err := query.Select("id, custom_code_id, title, description, keywords, is_original, source_note, code_types, status, reviewer_id, review_note, review_time, machine_audit, created_at, updated_at").
		Order("id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&list).Error

	return list, total, err
}

// 更新快照内容（草稿反复编辑）
func (m *CustomCodeReview) UpdateContent(db *gorm.DB, id uint, info CustomCodeBaseInfo) error {
	return db.Model(&CustomCodeReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":       info.Title,
		"description": info.Description,
		"keywords":    info.Keywords,
		"is_original": info.IsOriginal,
		"source_note": info.SourceNote,
		"versions":    info.Versions,
		"code_types":  info.CodeTypes,
	}).Error
}

// 更新审核状态
func (m *CustomCodeReview) UpdateStatus(db *gorm.DB, id uint, status int, reviewerId uint, reviewNote string) error {
	now := time.Now()
	return db.Model(&CustomCodeReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerId,
		"review_note": reviewNote,
		"review_time": now,
	}).Error
}
