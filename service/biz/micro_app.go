package biz

import (
	"fmt"
	"strings"
	"sun-panel/models"

	"gorm.io/gorm"
)

// microApp 微应用业务层
type microApp struct {
	// GetListCache
}

// func (s *microApp) Init() ([]models.MicroAppWithLang, int64, error) {

// 	global.NewCache[]()
// }

// GetById 根据id获取微应用
//   - extendField 扩展字段，用于预加载 Developer、LangList 字段
func (s *microApp) GetById(db *gorm.DB, id uint, extendField ...string) (models.MicroApp, error) {
	query := db
	for _, field := range extendField {
		query = query.Preload(field)
	}

	app := models.MicroApp{}
	var err error
	app, err = app.GetById(query, id)
	if err != nil {
		return models.MicroApp{}, NewBizError(ErrCodeAppNotFound)
	}

	return app, nil
}

// GetByIdWithLang 根据id获取微应用，并根据lang获取对应的语言信息
// 语言回退策略：
//   - 首先尝试指定的语言（如 zh-CN）
//   - 如果不存在，尝试 en 开头的语言（en-US, en-GB 等）
//   - 如果还不存在，使用第一个查询到的语言
//   - extendField 扩展字段，用于预加载 Developer、LangList 字段
func (s *microApp) GetByIdWithLang(db *gorm.DB, id uint, lang string, extendField ...string) (models.MicroApp, error) {
	query := db

	// 预加载所有多语言信息（用于回退逻辑）
	query = query.Preload("LangList")

	// 预加载其他扩展字段
	for _, field := range extendField {
		query = query.Preload(field)
	}

	app := models.MicroApp{}
	var err error
	app, err = app.GetById(query, id)
	if err != nil {
		return models.MicroApp{}, NewBizError(ErrCodeAppNotFound)
	}

	// 语言回退逻辑：从 LangList 中选择合适的语言填充 DefaultLangInfo
	selectedLang := s.selectBestLang(app.LangList, lang)
	if selectedLang != nil {
		app.DefaultLangInfo = *selectedLang
	}

	return app, nil
}

// 获取微应用的最新一条审核表的记录
func (s *microApp) GetMicroInfoAndLatestReview(db *gorm.DB, microAppModelid uint) (models.MicroAppReview, error) {
	mReview := models.MicroAppReview{}
	review, err := mReview.GetLatestByAppRecordId(db, microAppModelid)
	if err != nil {
		return models.MicroAppReview{}, err
	}
	return review, nil
}

// selectBestLang 根据语言回退策略选择最佳语言
// 1. 首选：完全匹配指定语言
// 2. 备选：en 开头的语言（如 en-US, en-GB）
// 3. 保底：第一个语言
func (s *microApp) selectBestLang(langList []models.MicroAppLang, preferredLang string) *models.MicroAppLang {
	if len(langList) == 0 {
		return nil
	}

	// 1. 首选：完全匹配指定语言
	for i := range langList {
		if langList[i].Lang == preferredLang {
			return &langList[i]
		}
	}

	// 2. 备选：en 开头的语言
	for i := range langList {
		if len(langList[i].Lang) >= 2 && langList[i].Lang[:2] == "en" {
			return &langList[i]
		}
	}

	// 3. 保底：返回第一个语言
	return &langList[0]
}

func (s *microApp) GetInfo(db *gorm.DB, microAppId string) (models.MicroApp, error) {
	var m models.MicroApp
	info, err := m.GetByMicroAppId(db, microAppId)
	if err != nil {
		return models.MicroApp{}, NewBizError(ErrCodeAppNotFound)
	}
	return info, nil
}

// BuildDownloadUrl 构建下载 URL
// 参数：
//   - appId: 微应用 ID
//   - version: 版本号（可选，为空时使用最新版本）
//
// 返回：
//   - 下载 URL
func (s *microApp) BuildDownloadUrl(microAppId string, version ...string) string {
	if len(version) == 0 || version[0] == "" {
		// 下载最新版本
		return fmt.Sprintf("/api/microApp/download/%s", microAppId)
	}
	// 下载指定版本
	return fmt.Sprintf("/api/microApp/download/%s/%s", microAppId, version)
}

// GetListOptions 微应用列表查询选项
type GetListOptions struct {
	Page       int    `json:"page"`       // 页码
	Limit      int    `json:"limit"`      // 每页数量
	Order      string `json:"order"`      // 排序（如 "download_count desc"）
	CategoryId uint   `json:"categoryId"` // 分类ID（0表示不筛选）
	Keyword    string `json:"keyword"`    // 关键词搜索
	Status     *int   `json:"status"`     // 状态（可选，默认为1-上架）
	Lang       string `json:"lang"`       // 语言（可选，暂不支持）
}

// GetList 获取微应用列表（公开接口）
// 参数：
//   - db: 数据库连接
//   - opts: 查询选项
//
// 返回：
//   - 微应用列表、总数、错误
func (s *microApp) GetList(db *gorm.DB, opts GetListOptions) ([]models.MicroAppWithLang, int64, error) {
	m := models.MicroApp{}
	status := 1 // 默认只查询上架的应用

	// 如果指定了状态则使用指定状态
	if opts.Status != nil {
		status = *opts.Status
	}

	// 处理分类参数（0 表示不筛选）
	var categoryId *int
	if opts.CategoryId > 0 {
		catId := int(opts.CategoryId)
		categoryId = &catId
	}

	// 解析排序参数
	sortBy := ""
	sortOrder := ""
	if opts.Order != "" {
		// order 格式: "field desc" 或 "field"
		parts := strings.Fields(opts.Order)
		if len(parts) > 0 {
			sortBy = parts[0]
			if len(parts) > 1 {
				sortOrder = strings.ToUpper(parts[1])
			}
		}
	}

	queryOpts := models.MicroAppQueryOptions{
		Page:             opts.Page,
		Limit:            opts.Limit,
		Status:           &status,
		CategoryId:       categoryId,
		KeyWord:          opts.Keyword,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
		Lang:             opts.Lang,
		IncludeDeveloper: true, // 包含开发者信息
	}

	return m.GetListWithLang(db, queryOpts)
}
