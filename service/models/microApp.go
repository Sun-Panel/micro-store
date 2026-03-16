package models

import (
	"gorm.io/gorm"
	"time"
)

// 微应用列表
type MicroApp struct {
	BaseModel
	MicroAppId      string  `gorm:"type:varchar(50);not null;uniqueIndex" json:"microAppId"` // 关联微应用ID（唯一标识）
	AppName         string  `gorm:"type:varchar(100);not null" json:"appName"`                // 应用名称（默认语言）
	AppIcon         string  `gorm:"type:varchar(200);not null" json:"appIcon"`                // 应用图标URL
	AppDesc         string  `gorm:"type:varchar(500)" json:"appDesc"`                         // 应用简介（默认语言）
	Remark          string  `gorm:"type:varchar(500)" json:"remark"`                          // 应用备注
	CategoryId      int     `gorm:"type:int(11);not null;index" json:"categoryId"`            // 应用分类ID
	ChargeType      int     `gorm:"type:tinyint(1);not null;default:0" json:"chargeType"`     // 收费方式：0-免费 1-积分 2-订阅PRO免费
	Price           float64 `gorm:"type:decimal(10,2)" json:"price"`                          // 价格（积分时）
	AuthorId        uint    `gorm:"type:int(11);not null;index" json:"authorId"`              // 开发者ID
	PermissionLevel int     `gorm:"type:tinyint(1)" json:"permissionLevel"`                   // 应用权限等级
	Status          int     `gorm:"type:tinyint(1);not null;default:1" json:"status"`         // 状态：0-下架 1-上架 2-审核中
	Screenshots     string  `gorm:"type:varchar(2000)" json:"screenshots"`                    // 图集（多个图片URL用逗号分隔）

	// 审核相关字段
	ReviewStatus int    `gorm:"type:tinyint(2);not null;default:0;index" json:"reviewStatus"` // 审核状态：0-无审核 1-审核中 2-已通过 3-已拒绝
	ReviewId     uint   `gorm:"type:int(11)" json:"reviewId"`                                 // 当前审核记录ID
	ReviewTime   *time.Time `gorm:"type:datetime" json:"reviewTime"`                         // 最后审核时间

	// 下架相关字段
	OfflineType   int    `gorm:"type:tinyint(1);not null;default:0" json:"offlineType"`     // 下架类型：0-正常 1-作者下架 2-平台下架
	OfflineReason string `gorm:"type:varchar(500)" json:"offlineReason"`                    // 下架原因

	// 关联多语言信息
	LangList []MicroAppLang `gorm:"foreignKey:MicroAppId;references:MicroAppId" json:"langList,omitempty"`
}

// 表名
func (MicroApp) TableName() string {
	return "micro_app"
}

// 获取微应用列表（支持分页和筛选）
func (m *MicroApp) GetList(db *gorm.DB, page, limit int, status *int, categoryId *int, authorId *uint, keyWord string) ([]MicroApp, int64, error) {
	var list []MicroApp
	var total int64

	query := db.Model(&MicroApp{})

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 分类筛选
	if categoryId != nil {
		query = query.Where("category_id = ?", *categoryId)
	}

	// 开发者筛选
	if authorId != nil {
		query = query.Where("author_id = ?", *authorId)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("app_name LIKE ? OR app_desc LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
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

// 获取上架的应用列表（用户端）
func (m *MicroApp) GetPublishedList(db *gorm.DB, page, limit int, categoryId *int, keyWord string) ([]MicroApp, int64, error) {
	var list []MicroApp
	var total int64

	query := db.Model(&MicroApp{}).Where("status = ?", 1) // 只查询上架的应用

	// 分类筛选
	if categoryId != nil {
		query = query.Where("category_id = ?", *categoryId)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("app_name LIKE ? OR app_desc LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
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

// 根据ID获取应用详情
func (m *MicroApp) GetById(db *gorm.DB, id uint) (MicroApp, error) {
	var app MicroApp
	err := db.Where("id = ?", id).First(&app).Error
	return app, err
}

// 创建应用
func (m *MicroApp) Create(db *gorm.DB) error {
	return db.Create(m).Error
}

// 更新应用
func (m *MicroApp) Update(db *gorm.DB, id uint, updateData map[string]interface{}) error {
	return db.Model(&MicroApp{}).Where("id = ?", id).Updates(updateData).Error
}

// 删除应用
func (m *MicroApp) Delete(db *gorm.DB, ids []uint) error {
	return db.Where("id IN ?", ids).Delete(&MicroApp{}).Error
}

// 更新应用状态
func (m *MicroApp) UpdateStatus(db *gorm.DB, id uint, status int) error {
	return db.Model(&MicroApp{}).Where("id = ?", id).Update("status", status).Error
}

// 批量更新状态
func (m *MicroApp) BatchUpdateStatus(db *gorm.DB, ids []uint, status int) error {
	return db.Model(&MicroApp{}).Where("id IN ?", ids).Update("status", status).Error
}

// 下架应用
func (m *MicroApp) Offline(db *gorm.DB, id uint, offlineType int, reason string) error {
	return db.Model(&MicroApp{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           0, // 下架状态
		"offline_type":    offlineType,
		"offline_reason":  reason,
	}).Error
}

// 获取开发者的应用数量
func (m *MicroApp) GetCountByAuthorId(db *gorm.DB, authorId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroApp{}).Where("author_id = ?", authorId).Count(&count).Error
	return count, err
}

// 根据分类ID获取应用数量
func (m *MicroApp) GetCountByCategoryId(db *gorm.DB, categoryId int) (int64, error) {
	var count int64
	err := db.Model(&MicroApp{}).Where("category_id = ?", categoryId).Count(&count).Error
	return count, err
}

// ===================== 多语言相关查询方法 =====================

// 带多语言信息的微应用DTO（用于返回给前端）
type MicroAppWithLang struct {
	MicroApp
	Lang       string `json:"lang"`       // 当前语言
	AppName    string `json:"appName"`    // 当前语言的应用名称
	AppDesc    string `json:"appDesc"`    // 当前语言的应用简介
}

// 获取微应用列表（预加载所有多语言信息）
func (m *MicroApp) GetListWithAllLangs(db *gorm.DB, page, limit int, status *int, categoryId *int, authorId *uint, keyWord string) ([]MicroApp, int64, error) {
	var list []MicroApp
	var total int64

	query := db.Model(&MicroApp{})

	// 状态筛选
	if status != nil {
		query = query.Where("micro_app.status = ?", *status)
	}

	// 分类筛选
	if categoryId != nil {
		query = query.Where("micro_app.category_id = ?", *categoryId)
	}

	// 开发者筛选
	if authorId != nil {
		query = query.Where("micro_app.author_id = ?", *authorId)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询并预加载多语言信息
	offset, limitSize := calcPage(page, limit)
	err = query.Order("micro_app.id DESC").
		Offset(offset).
		Limit(limitSize).
		Preload("LangList").
		Find(&list).Error

	return list, total, err
}

// 获取微应用列表（使用JOIN查询，只返回指定语言的信息）- 推荐，性能最佳
func (m *MicroApp) GetListWithLang(db *gorm.DB, page, limit int, lang string, status *int, categoryId *int, authorId *uint, keyWord string) ([]MicroAppWithLang, int64, error) {
	var list []MicroAppWithLang
	var total int64

	query := db.Table("micro_app").
		Select(`
			micro_app.*,
			COALESCE(lang.lang, 'default') as lang,
			COALESCE(lang.app_name, micro_app.app_name) as app_name,
			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
		`).
		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang)

	// 状态筛选
	if status != nil {
		query = query.Where("micro_app.status = ?", *status)
	}

	// 分类筛选
	if categoryId != nil {
		query = query.Where("micro_app.category_id = ?", *categoryId)
	}

	// 开发者筛选
	if authorId != nil {
		query = query.Where("micro_app.author_id = ?", *authorId)
	}

	// 关键词搜索（搜索多语言内容或默认内容）
	if keyWord != "" {
		query = query.Where("lang.app_name LIKE ? OR lang.app_desc LIKE ? OR micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?",
			"%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("micro_app.id DESC").
		Offset(offset).
		Limit(limitSize).
		Scan(&list).Error

	return list, total, err
}

// 获取上架的应用列表（带指定语言）- 用户端推荐使用
func (m *MicroApp) GetPublishedListWithLang(db *gorm.DB, page, limit int, lang string, categoryId *int, keyWord string) ([]MicroAppWithLang, int64, error) {
	var list []MicroAppWithLang
	var total int64

	query := db.Table("micro_app").
		Select(`
			micro_app.*,
			COALESCE(lang.lang, 'default') as lang,
			COALESCE(lang.app_name, micro_app.app_name) as app_name,
			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
		`).
		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
		Where("micro_app.status = ?", 1) // 只查询上架的应用

	// 分类筛选
	if categoryId != nil {
		query = query.Where("micro_app.category_id = ?", *categoryId)
	}

	// 关键词搜索
	if keyWord != "" {
		query = query.Where("lang.app_name LIKE ? OR lang.app_desc LIKE ? OR micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?",
			"%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(page, limit)
	err = query.Order("micro_app.id DESC").
		Offset(offset).
		Limit(limitSize).
		Scan(&list).Error

	return list, total, err
}

// 根据ID获取应用详情（带指定语言）
func (m *MicroApp) GetByIdWithLang(db *gorm.DB, id uint, lang string) (MicroAppWithLang, error) {
	var result MicroAppWithLang
	err := db.Table("micro_app").
		Select(`
			micro_app.*,
			COALESCE(lang.lang, 'default') as lang,
			COALESCE(lang.app_name, micro_app.app_name) as app_name,
			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
		`).
		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
		Where("micro_app.id = ?", id).
		Scan(&result).Error

	return result, err
}

// 根据MicroAppId获取应用详情（带指定语言）
func (m *MicroApp) GetByMicroAppIdWithLang(db *gorm.DB, microAppId string, lang string) (MicroAppWithLang, error) {
	var result MicroAppWithLang
	err := db.Table("micro_app").
		Select(`
			micro_app.*,
			COALESCE(lang.lang, 'default') as lang,
			COALESCE(lang.app_name, micro_app.app_name) as app_name,
			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
		`).
		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
		Where("micro_app.micro_app_id = ?", microAppId).
		Scan(&result).Error

	return result, err
}

