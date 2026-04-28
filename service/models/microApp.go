package models

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// MicroAppBaseInfo 微应用基础信息（公共字段）
type MicroAppBaseInfo struct {
	// MicroAppId  string  `gorm:"type:varchar(120);not null;index" json:"microAppId"`   // 关联微应用ID（唯一标识）
	AdminName string `gorm:"type:varchar(100);" json:"adminName"` // 应用名称(开发者相关页面可见)
	// AppName     string `gorm:"type:varchar(100);not null" json:"appName"`             // 应用名称（默认语言） 废弃-改用多语言表
	AppIcon string `gorm:"type:varchar(200);not null" json:"appIcon"` // 应用图标URL
	// AppDesc     string `gorm:"type:varchar(500)" json:"appDesc"`                      // 应用简介（默认语言）废弃-改用多语言表
	Remark      string `gorm:"type:varchar(500)" json:"remark"`                       // 应用备注
	CategoryId  int    `gorm:"type:int(11);not null;index" json:"categoryId"`         // 应用分类ID
	ChargeType  int    `gorm:"type:tinyint(1);not null;default:0" json:"chargeType"`  // 收费方式：0-免费 1-积分 2-订阅PRO免费
	Points      int    `gorm:"type:int(11)" json:"points"`                            // 价格（积分数值）
	Screenshots string `gorm:"type:varchar(2000)" json:"screenshots"`                 // 图集（多个图片URL用逗号分隔）
	ThirdCharge int    `gorm:"type:tinyint(1);not null;default:0" json:"thirdCharge"` // 第三方收费方式：0-不含 1-付费才可用 2-基础功能免费
	HaveIframe  bool   `gorm:"type:tinyint(1);not null;default:0" json:"haveIframe"`  // 是否包含iframe
}

// 微应用表（只存储生效版本）
type MicroApp struct {
	BaseModel
	MicroAppBaseInfo `gorm:"embedded"` // 嵌入公共字段

	MicroAppId  string         `gorm:"column:micro_app_id;type:varchar(120);not null;uniqueIndex:idx_micro_app_id_deleted_at" json:"microAppId"` // 覆盖嵌入字段，添加复合唯一约束（支持软删除）
	DeletedAt   gorm.DeletedAt `gorm:"uniqueIndex:idx_micro_app_id_deleted_at"`                                                                  // 覆盖基类字段，加入复合唯一索引
	DeveloperId uint           `gorm:"type:int(11);not null;index" json:"developerId"`                                                           // 开发者ID                                                              // 应用权限等级
	Status      int            `gorm:"type:tinyint(1);not null;default:0" json:"status"`                                                         // 状态：0-下架 1-上架

	// 下架相关字段
	OfflineType   int    `gorm:"type:tinyint(1);not null;default:0" json:"offlineType"` // 下架类型：0-正常 1-作者下架 2-平台下架 3-首次创建
	OfflineReason string `gorm:"type:varchar(500)" json:"offlineReason"`                // 下架原因

	DownloadCount int `gorm:"type:int(11);not null;default:0" json:"downloadCount"`
	InstallCount  int `gorm:"type:int(11);not null;default:0" json:"installCount"`

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
	SortBy           string // 排序字段：id, download_count, install_count
	SortOrder        string // 排序方式：asc, desc
}

// 获取上架的应用列表（用户端）
// func (m *MicroApp) GetPublishedList(db *gorm.DB, page, limit int, categoryId *int, keyWord string) ([]MicroApp, int64, error) {
// 	var list []MicroApp
// 	var total int64

// 	query := db.Model(&MicroApp{}).Where("status = ?", 1) // 只查询上架的应用

// 	// 分类筛选
// 	if categoryId != nil {
// 		query = query.Where("category_id = ?", *categoryId)
// 	}

// 	// 关键词搜索
// 	if keyWord != "" {
// 		query = query.Where("app_name LIKE ? OR app_desc LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%")
// 	}

// 	// 获取总数
// 	err := query.Count(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// 分页查询
// 	offset, limitSize := calcPage(page, limit)
// 	err = query.Order("id DESC").Offset(offset).Limit(limitSize).Find(&list).Error

// 	return list, total, err
// }

// 根据ID获取应用详情
func (m *MicroApp) GetById(db *gorm.DB, id uint) (MicroApp, error) {
	var app MicroApp
	err := db.Where("id = ?", id).First(&app).Error
	return app, err
}

// 根据MicroAppId获取应用详情
func (m *MicroApp) GetByMicroAppId(db *gorm.DB, microAppId string) (MicroApp, error) {
	var app MicroApp
	err := db.Order("id DESC").Where("micro_app_id = ?", microAppId).First(&app).Error
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

// Deprecated: 弃用，请使用 GetAppList
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
// func (m *MicroApp) GetListWithLang(db *gorm.DB, opts MicroAppQueryOptions) ([]MicroAppWithLang, int64, error) {
// 	var list []MicroAppWithLang
// 	var total int64

// 	// 构建基础查询
// 	query := db.Table("micro_app").
// 		Select(`
// 			micro_app.*,
// 			COALESCE(lang.lang, 'default') as lang,
// 			COALESCE(lang.app_name, micro_app.app_name) as app_name,
// 			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
// 		`).
// 		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", opts.Lang)

// 	// 如果需要查询开发者信息,添加 LEFT JOIN 和 Select 开发者字段
// 	if opts.IncludeDeveloper {
// 		query = query.
// 			Select(`
// 				micro_app.*,
// 				COALESCE(lang.lang, 'default') as lang,
// 				COALESCE(lang.app_name, micro_app.app_name) as app_name,
// 				COALESCE(lang.app_desc, micro_app.app_desc) as app_desc,
// 				developer.name as developer_name
// 			`).
// 			Joins("LEFT JOIN developer ON micro_app.developer_id = developer.id")
// 	}

// 	// 状态筛选
// 	if opts.Status != nil {
// 		query = query.Where("micro_app.status = ?", *opts.Status)
// 	}

// 	// 分类筛选
// 	if opts.CategoryId != nil {
// 		query = query.Where("micro_app.category_id = ?", *opts.CategoryId)
// 	}

// 	// 开发者筛选
// 	if opts.DeveloperId != nil {
// 		query = query.Where("micro_app.developer_id = ?", *opts.DeveloperId)
// 	}

// 	// 关键词搜索（搜索多语言内容或默认内容）
// 	if opts.KeyWord != "" {
// 		query = query.Where("lang.app_name LIKE ? OR lang.app_desc LIKE ? OR micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?",
// 			"%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%", "%"+opts.KeyWord+"%")
// 	}

// 	// 获取总数
// 	err := query.Count(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// 分页查询
// 	offset, limitSize := calcPage(opts.Page, opts.Limit)
// 	err = query.Order("micro_app.id DESC").
// 		Offset(offset).
// 		Limit(limitSize).
// 		Scan(&list).Error

// 	return list, total, err
// }

// // 获取上架的应用列表（带指定语言）- 用户端推荐使用
// func (m *MicroApp) GetPublishedListWithLang(db *gorm.DB, page, limit int, lang string, categoryId *int, keyWord string) ([]MicroAppWithLang, int64, error) {
// 	var list []MicroAppWithLang
// 	var total int64

// 	query := db.Table("micro_app").
// 		Select(`
// 			micro_app.*,
// 			COALESCE(lang.lang, 'default') as lang,
// 			COALESCE(lang.app_name, micro_app.app_name) as app_name,
// 			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
// 		`).
// 		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
// 		Where("micro_app.status = ?", 1) // 只查询上架的应用

// 	// 分类筛选
// 	if categoryId != nil {
// 		query = query.Where("micro_app.category_id = ?", *categoryId)
// 	}

// 	// 关键词搜索
// 	if keyWord != "" {
// 		query = query.Where("lang.app_name LIKE ? OR lang.app_desc LIKE ? OR micro_app.app_name LIKE ? OR micro_app.app_desc LIKE ?",
// 			"%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%", "%"+keyWord+"%")
// 	}

// 	// 获取总数
// 	err := query.Count(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// 分页查询
// 	offset, limitSize := calcPage(page, limit)
// 	err = query.Order("micro_app.id DESC").
// 		Offset(offset).
// 		Limit(limitSize).
// 		Scan(&list).Error

// 	return list, total, err
// }

// 根据ID获取应用详情（带指定语言）
// func (m *MicroApp) GetByIdWithLang(db *gorm.DB, id uint, lang string) (MicroAppWithLang, error) {
// 	var result MicroAppWithLang
// 	err := db.Table("micro_app").
// 		Select(`
// 			micro_app.*,
// 			COALESCE(lang.lang, 'default') as lang,
// 			COALESCE(lang.app_name, micro_app.app_name) as app_name,
// 			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
// 		`).
// 		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
// 		Where("micro_app.id = ?", id).
// 		Scan(&result).Error

// 	return result, err
// }

// 根据MicroAppId获取应用详情（带指定语言）
// func (m *MicroApp) GetByMicroAppIdWithLang(db *gorm.DB, microAppId string, lang string) (MicroAppWithLang, error) {
// 	var result MicroAppWithLang
// 	err := db.Table("micro_app").
// 		Select(`
// 			micro_app.*,
// 			COALESCE(lang.lang, 'default') as lang,
// 			COALESCE(lang.app_name, micro_app.app_name) as app_name,
// 			COALESCE(lang.app_desc, micro_app.app_desc) as app_desc
// 		`).
// 		Joins("LEFT JOIN micro_app_lang lang ON micro_app.micro_app_id = lang.micro_app_id AND lang.lang = ?", lang).
// 		Where("micro_app.micro_app_id = ?", microAppId).
// 		Scan(&result).Error

// 	return result, err
// }

// ===================== 应用列表查询（v2）=====================

// MicroAppListQueryOpts 应用列表查询选项（基础）
type MicroAppListQueryOpts struct {
	// 查询条件
	Status      *int   // 状态筛选
	OfflineType *int   // 下架类型筛选
	CategoryId  *int   // 分类筛选
	DeveloperId *uint  // 开发者筛选
	KeyWord     string // 关键字搜索（搜索多语言表的 app_name 和 app_desc）

	// 分页
	Page  int // 页码
	Limit int // 每页数量

	// 排序
	SortBy    string // 排序字段：id, download_count, install_count, created_at
	SortOrder string // 排序方式：asc, desc

	// 可选加载
	IncludeLangList  bool // 是否预加载全部多语言列表（编辑场景）
	IncludeDeveloper bool // 是否预加载开发者信息
}

// MicroAppListWithLangQueryOpts 带语言的应用列表查询选项（继承基础选项）
type MicroAppListWithLangQueryOpts struct {
	MicroAppListQueryOpts
	Lang          string   // 目标语言，如 "zh-CN"
	FallbackLangs []string // 回退语言列表，如 ["en-US"]，按优先级排列
}

// MicroAppListItem 应用列表项
// 嵌入 MicroApp，继承其 LangList/Developer/DefaultLangInfo 关联字段，支持 Preload
// 额外的 AppName/AppDesc/LangLabel 由 GetAppListWithLang 的 JOIN 填充
type MicroAppListItem struct {
	MicroApp
	AppName   string `json:"appName"`             // 当前语言的应用名称（GetAppListWithLang JOIN 填充）
	AppDesc   string `json:"appDesc"`             // 当前语言的应用简介（GetAppListWithLang JOIN 填充）
	LangLabel string `json:"langLabel,omitempty"` // 实际命中的语言代码（GetAppListWithLang JOIN 填充）
}

// GetAppList 获取应用列表（基于 gorm 原生查询）
// 返回基础字段 + 可选 LangList + 可选 Developer 信息
func (m *MicroApp) GetAppList(db *gorm.DB, opts MicroAppListQueryOpts) ([]MicroAppListItem, int64, error) {
	var list []MicroAppListItem
	var total int64

	query := db.Table("micro_app").Where("micro_app.deleted_at IS NULL")

	// 状态筛选
	if opts.Status != nil {
		query = query.Where("micro_app.status = ?", *opts.Status)
	}

	// 下架类型筛选
	if opts.OfflineType != nil {
		query = query.Where("micro_app.offline_type = ?", *opts.OfflineType)
	}

	// 分类筛选
	if opts.CategoryId != nil {
		query = query.Where("micro_app.category_id = ?", *opts.CategoryId)
	}

	// 开发者筛选
	if opts.DeveloperId != nil {
		query = query.Where("micro_app.developer_id = ?", *opts.DeveloperId)
	}

	// 关键字搜索（匹配 microAppId 或多语言表的 app_name/app_desc）
	if opts.KeyWord != "" {
		like := "%" + opts.KeyWord + "%"
		subQuery := db.Model(&MicroAppLang{}).
			Select("DISTINCT micro_app_id").
			Where("app_name LIKE ? OR app_desc LIKE ?", like, like)
		query = query.Where("micro_app.micro_app_id LIKE ? OR micro_app.micro_app_id IN (?)", like, subQuery)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "micro_app.id DESC"
	if opts.SortBy != "" {
		sortOrder := "DESC"
		if opts.SortOrder != "" {
			sortOrder = opts.SortOrder
		}
		// 白名单校验，防止注入
		allowedSortFields := map[string]string{
			"id":             "micro_app.id",
			"download_count": "micro_app.download_count",
			"install_count":  "micro_app.install_count",
			"created_at":     "micro_app.created_at",
		}
		if field, ok := allowedSortFields[opts.SortBy]; ok {
			order = field + " " + sortOrder
		}
	}
	query = query.Order(order)

	// 分页
	offset, limitSize := calcPage(opts.Page, opts.Limit)
	query = query.Offset(offset).Limit(limitSize)

	// 预加载多语言列表
	if opts.IncludeLangList {
		query = query.Preload("LangList")
	}

	// 预加载开发者信息
	if opts.IncludeDeveloper {
		query = query.Preload("Developer")
	}

	if err := query.Select("micro_app.*").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetAppListWithLang 获取应用列表（JOIN 方式只取当前语言，带回退策略）
// 优先匹配 Lang，依次回退 FallbackLangs，最后回退到任意一条多语言记录
// 前端直接使用 appName/appDesc 即可，无需遍历匹配
func (m *MicroApp) GetAppListWithLang(db *gorm.DB, opts MicroAppListWithLangQueryOpts) ([]MicroAppListItem, int64, error) {
	var list []MicroAppListItem
	var total int64

	// ========== 第一步：构建公共 WHERE 条件 ==========
	// Count 查询不需要 JOIN，先基于基础表计数，避免 correlated subquery 开销
	countQuery := db.Model(&MicroApp{})

	// 状态筛选
	if opts.Status != nil {
		countQuery = countQuery.Where("micro_app.status = ?", *opts.Status)
	}

	// 下架类型筛选
	if opts.OfflineType != nil {
		countQuery = countQuery.Where("micro_app.offline_type = ?", *opts.OfflineType)
	}

	// 分类筛选
	if opts.CategoryId != nil {
		countQuery = countQuery.Where("micro_app.category_id = ?", *opts.CategoryId)
	}

	// 开发者筛选
	if opts.DeveloperId != nil {
		countQuery = countQuery.Where("micro_app.developer_id = ?", *opts.DeveloperId)
	}

	// 关键字搜索（在多语言表中匹配，或匹配 microAppId）
	if opts.KeyWord != "" {
		like := "%" + opts.KeyWord + "%"
		subQuery := db.Model(&MicroAppLang{}).
			Select("DISTINCT micro_app_id").
			Where("app_name LIKE ? OR app_desc LIKE ?", like, like)
		countQuery = countQuery.Where("micro_app.micro_app_id LIKE ? OR micro_app.micro_app_id IN (?)", like, subQuery)
	}

	// 获取总数（不走 JOIN）
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ========== 第二步：构建 JOIN 数据查询 ==========
	// 构建 CASE WHEN 优先级排序：目标语言 > 回退语言 > 任意其他
	// 注意：子查询中用 `lang`（表名），主查询 JOIN 后用 `lang.lang`（别名）
	var caseBuilderForSubQuery strings.Builder
	caseBuilderForSubQuery.WriteString("CASE WHEN lang = '")
	caseBuilderForSubQuery.WriteString(opts.Lang)
	caseBuilderForSubQuery.WriteString("' THEN 0 ")
	for i, fl := range opts.FallbackLangs {
		caseBuilderForSubQuery.WriteString("WHEN lang = '")
		caseBuilderForSubQuery.WriteString(fl)
		caseBuilderForSubQuery.WriteString("' THEN ")
		caseBuilderForSubQuery.WriteString(strconv.Itoa(i + 1))
		caseBuilderForSubQuery.WriteString(" ")
	}
	caseBuilderForSubQuery.WriteString("ELSE 100 END")
	caseExprForSubQuery := caseBuilderForSubQuery.String()

	// 子查询：每个 micro_app_id 只取优先级最高的一条多语言记录
	langSubQuery := db.Model(&MicroAppLang{}).
		Select("id").
		Where("micro_app_id = micro_app.micro_app_id").
		Where("deleted_at IS NULL").
		Order(caseExprForSubQuery).
		Limit(1)

	query := db.Table("micro_app").
		Select(`micro_app.*,
			COALESCE(lang.app_name, '') as app_name,
			COALESCE(lang.app_desc, '') as app_desc,
			COALESCE(lang.lang, '') as lang_label`).
		Joins("LEFT JOIN micro_app_lang lang ON lang.id = (?)", langSubQuery).
		Where("micro_app.deleted_at IS NULL")

	// 复用相同的 WHERE 条件
	if opts.Status != nil {
		query = query.Where("micro_app.status = ?", *opts.Status)
	}
	if opts.OfflineType != nil {
		query = query.Where("micro_app.offline_type = ?", *opts.OfflineType)
	}
	if opts.CategoryId != nil {
		query = query.Where("micro_app.category_id = ?", *opts.CategoryId)
	}
	if opts.DeveloperId != nil {
		query = query.Where("micro_app.developer_id = ?", *opts.DeveloperId)
	}
	if opts.KeyWord != "" {
		like := "%" + opts.KeyWord + "%"
		subQuery := db.Model(&MicroAppLang{}).
			Select("DISTINCT micro_app_id").
			Where("app_name LIKE ? OR app_desc LIKE ?", like, like)
		query = query.Where("micro_app.micro_app_id LIKE ? OR micro_app.micro_app_id IN (?)", like, subQuery)
	}

	// 排序
	order := "micro_app.id DESC"
	if opts.SortBy != "" {
		sortOrder := "DESC"
		if opts.SortOrder != "" {
			sortOrder = opts.SortOrder
		}
		allowedSortFields := map[string]string{
			"id":             "micro_app.id",
			"download_count": "micro_app.download_count",
			"install_count":  "micro_app.install_count",
			"created_at":     "micro_app.created_at",
		}
		if field, ok := allowedSortFields[opts.SortBy]; ok {
			order = field + " " + sortOrder
		}
	}
	query = query.Order(order)

	// 分页
	offset, limitSize := calcPage(opts.Page, opts.Limit)
	query = query.Offset(offset).Limit(limitSize)

	// 预加载全部多语言列表（编辑场景可选）
	if opts.IncludeLangList {
		// JOIN 模式下无法直接 Preload，需二次查询后填充
		if err := query.Scan(&list).Error; err != nil {
			return nil, 0, err
		}
		if err := m.fillLangList(db, list); err != nil {
			return nil, 0, err
		}
		if opts.IncludeDeveloper {
			if err := m.fillDeveloper(db, list); err != nil {
				return nil, 0, err
			}
		}
		return list, total, nil
	}

	// 预加载开发者信息
	if opts.IncludeDeveloper {
		// JOIN 模式下需二次查询填充
		if err := query.Scan(&list).Error; err != nil {
			return nil, 0, err
		}
		if err := m.fillDeveloper(db, list); err != nil {
			return nil, 0, err
		}
		return list, total, nil
	}

	if err := query.Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// fillLangList 填充多语言列表（JOIN 模式下二次查询）
func (m *MicroApp) fillLangList(db *gorm.DB, list []MicroAppListItem) error {
	if len(list) == 0 {
		return nil
	}
	ids := make([]string, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.MicroAppId)
	}
	var langList []MicroAppLang
	if err := db.Where("micro_app_id IN ?", ids).Find(&langList).Error; err != nil {
		return err
	}
	langMap := make(map[string][]MicroAppLang)
	for _, l := range langList {
		langMap[l.MicroAppId] = append(langMap[l.MicroAppId], l)
	}
	for i := range list {
		list[i].LangList = langMap[list[i].MicroAppId]
	}
	return nil
}

// fillDeveloper 填充开发者信息（JOIN 模式下二次查询）
func (m *MicroApp) fillDeveloper(db *gorm.DB, list []MicroAppListItem) error {
	if len(list) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(list))
	idSet := make(map[uint]bool)
	for _, item := range list {
		if !idSet[item.DeveloperId] {
			idSet[item.DeveloperId] = true
			ids = append(ids, item.DeveloperId)
		}
	}
	var developers []Developer
	if err := db.Where("id IN ?", ids).Find(&developers).Error; err != nil {
		return err
	}
	devMap := make(map[uint]*Developer)
	for i := range developers {
		devMap[developers[i].ID] = &developers[i]
	}
	for i := range list {
		if dev, ok := devMap[list[i].DeveloperId]; ok {
			list[i].Developer = *dev
		}
	}
	return nil
}
