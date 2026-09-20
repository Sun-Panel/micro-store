package models

import (
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

// 代码片块数量与预览图限制
const (
	CustomCodeMaxBlocks         = 10        // 单帖最多块数
	CustomCodeMaxImagesPerBlock = 5         // 每块最多预览图
	CustomCodeMaxImageSize      = 512 * 1024 // 单张预览图大小上限（字节）
)

// CustomCodeBlock 代码片块
// review_id = 0 表示线上块；非 0 表示某次审核快照的块
type CustomCodeBlock struct {
	BaseModel
	CustomCodeId uint                 `gorm:"type:int(11);not null;index" json:"customCodeId"`
	ReviewId     uint                 `gorm:"type:int(11);not null;default:0;index" json:"reviewId"`
	Sort         int                  `gorm:"type:int(11);not null;default:0" json:"sort"`
	CodeType     int                  `gorm:"type:tinyint(2);not null" json:"codeType"` // 1-JS 2-CSS 3-页脚
	Title        string               `gorm:"type:varchar(100)" json:"title"`
	Note         string               `gorm:"type:text" json:"note"`
	Code         string               `gorm:"type:longtext;not null" json:"code"`
	Images       datatype.StringArray `gorm:"type:varchar(1000)" json:"images"`
	OnlyId       string               `gorm:"type:varchar(40);not null;default:''" json:"onlyId"` // 块唯一标识（开发者可改，跨版本保持稳定，用于复制粘贴去重）
}

// 表名
func (CustomCodeBlock) TableName() string {
	return "custom_code_block"
}

// GetListByCustomCodeId 获取指定版本的块列表（reviewId 为 0 时取线上块）
func (m *CustomCodeBlock) GetListByCustomCodeId(db *gorm.DB, customCodeId uint, reviewId uint) ([]CustomCodeBlock, error) {
	var list []CustomCodeBlock
	err := db.Where("custom_code_id = ? AND review_id = ?", customCodeId, reviewId).
		Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// ReplaceAll 全量替换某版本的块（调用方需自行放在事务中）
func (m *CustomCodeBlock) ReplaceAll(db *gorm.DB, customCodeId uint, reviewId uint, blocks []CustomCodeBlock) error {
	if err := db.Where("custom_code_id = ? AND review_id = ?", customCodeId, reviewId).
		Delete(&CustomCodeBlock{}).Error; err != nil {
		return err
	}

	for i := range blocks {
		block := blocks[i]
		block.ID = 0
		block.CustomCodeId = customCodeId
		block.ReviewId = reviewId
		block.Sort = i
		if err := db.Create(&block).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteByCustomCodeId 删除某帖子的块（reviewId 为 0 时删线上块）
func (m *CustomCodeBlock) DeleteByCustomCodeId(db *gorm.DB, customCodeId uint, reviewId uint) error {
	return db.Where("custom_code_id = ? AND review_id = ?", customCodeId, reviewId).
		Delete(&CustomCodeBlock{}).Error
}

// DeleteByReviewId 删除某次审核快照的块
func (m *CustomCodeBlock) DeleteByReviewId(db *gorm.DB, reviewId uint) error {
	return db.Where("review_id = ?", reviewId).Delete(&CustomCodeBlock{}).Error
}
