package models

import (
	"gorm.io/gorm"
)

// MicroAppBaseInfo 微应用基础信息（公共字段）
type MicroAppBaseInfo struct {
	// MicroAppId  string  `gorm:"type:varchar(120);not null;index" json:"microAppId"`   // 关联微应用ID（唯一标识）
	AppName     string `gorm:"type:varchar(100);not null" json:"appName"`            // 应用名称（默认语言）
	AppIcon     string `gorm:"type:varchar(200);not null" json:"appIcon"`            // 应用图标URL
	AppDesc     string `gorm:"type:varchar(500)" json:"appDesc"`                     // 应用简介（默认语言）
	Remark      string `gorm:"type:varchar(500)" json:"remark"`                      // 应用备注
	CategoryId  int    `gorm:"type:int(11);not null;index" json:"categoryId"`        // 应用分类ID
	ChargeType  int    `gorm:"type:tinyint(1);not null;default:0" json:"chargeType"` // 收费方式：0-免费 1-积分 2-订阅PRO免费
	Points      int    `gorm:"type:int(11)" json:"Points"`                           // 价格（积分数值）
	Screenshots string `gorm:"type:varchar(2000)" json:"screenshots"`                // 图集（多个图片URL用逗号分隔）
}

// 微应用表（只存储生效版本）
type MicroApp struct {
	BaseModel
	MicroAppBaseInfo `gorm:"embedded"` // 嵌入公共字段
	MicroAppId       string            `gorm:"column:micro_app_id;type:varchar(120);not null;uniqueIndex:idx_micro_app_id_deleted_at" json:"microAppId"` // 覆盖嵌入字段，添加复合唯一约束（支持软删除）
	DeletedAt        gorm.DeletedAt    `gorm:"uniqueIndex:idx_micro_app_id_deleted_at"`                                                                  // 覆盖基类字段，加入复合唯一索引
	DeveloperId      uint              `gorm:"type:int(11);not null;index" json:"developerId"`                                                           // 开发者ID                                                              // 应用权限等级
	Status           int               `gorm:"type:tinyint(1);not null;default:0" json:"status"`                                                         // 状态：0-下架 1-上架

	// 下架相关字段
	OfflineType   int    `gorm:"type:tinyint(1);not null;default:0" json:"offlineType"` // 下架类型：0-正常 1-作者下架 2-平台下架 3-首次创建
	OfflineReason string `gorm:"type:varchar(500)" json:"offlineReason"`                // 下架原因

	// 关联多语言信息
	LangList        []MicroAppLang `gorm:"foreignKey:MicroAppId;references:MicroAppId" json:"langList,omitempty"`
	DefaultLangInfo MicroAppLang   `gorm:"foreignKey:MicroAppId;references:MicroAppId" json:"defaultLangInfo"`
	Developer       Developer      `gorm:"foreignKey:DeveloperId;references:ID" json:"developer"`
}

// 表名
func (MicroApp) TableName() string {
	return "micro_app"
}

// 微应用查询选项
type MicroAppQueryOptions struct {
	Page             int
	Limit            int
	Status           *int
	CategoryId       *int
	DeveloperId      *uint
	KeyWord          string
	Lang             string // 可选,用于多语言查询
	IncludeDeveloper bool   // 是否查询开发者信息
}

// 获取微应用列表（支持分页和筛选）
func (m *MicroApp) GetList(db *gorm.DB, opts MicroAppQueryOptions) ([]MicroApp, int64, error) {
	var list []MicroApp
	var total int64

	query := db.Model(&MicroApp{})

	// 开发者列表：只显示最新版本（优先生效版本）
	if opts.DeveloperId != nil {
		// 子查询：获取每个 microAppId 的最新记录
		// 优先选择 review_status=0（生效版本），否则选择最新的记录
		subQuery := db.Model(&MicroApp{}).
			Select("COALESCE(MAX(CASE WHEN review_status = 0 THEN id END), MAX(id)) as id").
			Where("developer_id = ?", *opts.DeveloperId).
			Group("micro_app_id")

		query = query.Where("id IN (?)", subQuery)
	}

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	// 分类筛选
	if opts.CategoryId != nil {
		query = query.Where("category_id = ?", *opts.CategoryId)
	}

	// 关键词搜索
	if opts.KeyWord != "" {
		query = query.Where("app_name LIKE ? OR app_desc LIKE ?", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(opts.Page, opts.Limit)
	query = query.Order("id DESC").Offset(offset).Limit(limitSize)

	// 预加载开发者信息
	if opts.IncludeDeveloper {
		query = query.Preload("Developer")
	}

	err = query.Find(&list).Error

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
		"status":         0, // 下架状态
		"offline_type":   offlineType,
		"offline_reason": reason,
	}).Error
}

// 获取开发者的应用数量
func (m *MicroApp) GetCountByAuthorId(db *gorm.DB, authorId uint) (int64, error) {
	var count int64
	err := db.Model(&MicroApp{}).Where("developer_id = ?", authorId).Count(&count).Error
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
	Lang            string `json:"lang"`            // 当前语言
	AppName         string `json:"appName"`         // 当前语言的应用名称
	AppDesc         string `json:"appDesc"`         // 当前语言的应用简介
	DeveloperId     uint   `json:"developerId"`     // 开发者ID
	DeveloperName   string `json:"developerName"`   // 开发者名称
	DeveloperAvatar string `json:"developerAvatar"` // 开发者头像
}

// 获取微应用列表（预加载所有多语言信息）
func (m *MicroApp) GetListWithAllLangs(db *gorm.DB, opts MicroAppQueryOptions) ([]MicroApp, int64, error) {
	var list []MicroApp
	var total int64

	query := db.Model(&MicroApp{})

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("micro_app.status = ?", *opts.Status)
	}

	// 分类筛选
	if opts.CategoryId != nil {
		query = query.Where("micro_app.category_id = ?", *opts.CategoryId)
	}

	// 开发者筛选
	if opts.DeveloperId != nil {
		query = query.Where("micro_app.developer_id = ?", *opts.DeveloperId)
	}

	// 关键词搜索
	if opts.KeyWord != "" {
		query = query.Where("micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询并预加载多语言信息
	offset, limitSize := calcPage(opts.Page, opts.Limit)
	query = query.Order("micro_app.id DESC").
		Offset(offset).
		Limit(limitSize).
		Preload("LangList")

	// 预加载开发者信息
	if opts.IncludeDeveloper {
		query = query.Preload("Developer")
	}

	err = query.Find(&list).Error

	return list, total, err
}

// 获取微应用列表（使用JOIN查询，只返回指定语言的信息）- 推荐，性能最佳
func (m *MicroApp) GetListWithLang(db *gorm.DB, opts MicroAppQueryOptions) ([]MicroAppWithLang, int64, error) {
	var list []MicroAppWithLang
	var total int64

	// 构建基础查询
	query := db.Table("micro_app").
		Select(`
			micro_app.*,
			COALESCE(lang.lang, 'default') as lang,
			COALESCE(lang.app_name, micro_app.app_name) as app_name,
			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
		`).
		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", opts.Lang)

	// 如果需要查询开发者信息,添加 LEFT JOIN 和 Select 开发者字段
	if opts.IncludeDeveloper {
		query = query.
			Select(`
				micro_app.*,
				COALESCE(lang.lang, 'default') as lang,
				COALESCE(lang.app_name, micro_app.app_name) as app_name,
				COALESCE(lang.app_desc, micro_app.app_desc) as app_desc,
				user.id as developer_id,
				user.name as developer_name,
				user.avatar as developer_avatar
			`).
			Joins("LEFT JOIN user ON micro_app.developer_id = user.id")
	}

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("micro_app.status = ?", *opts.Status)
	}

	// 分类筛选
	if opts.CategoryId != nil {
		query = query.Where("micro_app.category_id = ?", *opts.CategoryId)
	}

	// 开发者筛选
	if opts.DeveloperId != nil {
		query = query.Where("micro_app.developer_id = ?", *opts.DeveloperId)
	}

	// 关键词搜索（搜索多语言内容或默认内容）
	if opts.KeyWord != "" {
		query = query.Where("lang.app_name LIKE ? OR lang.app_desc LIKE ? OR micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?",
			"%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset, limitSize := calcPage(opts.Page, opts.Limit)
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
