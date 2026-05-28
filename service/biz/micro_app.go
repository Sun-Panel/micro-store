package biz

import (
	"fmt"
	"strings"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

// MicroAppInfoCache 微应用信息缓存结构（对外暴露的精简信息）
type MicroAppInfoCache struct {
	MicroAppId     string `json:"microAppId"`
	AppIcon        string `json:"appIcon"`
	ChargeType     int    `json:"chargeType"`
	Points         int    `json:"points"`
	DeveloperId    uint   `json:"developerId"`
	DeveloperName  string `json:"developerName"`  // 开发者名称
	DeveloperName2 string `json:"developerName2"` // 开发者标识
}

// microApp 微应用业务层
type microApp struct {
	// 微应用信息缓存，key: microAppId，value: MicroAppInfoCache
	InfoCache cache.Cacher[MicroAppInfoCache]
}

// Init 初始化微应用服务（含缓存）
func (s *microApp) Init() {
	s.InfoCache = global.NewCache[MicroAppInfoCache](30*time.Minute, 10*time.Minute, "MicroAppInfoCache")
}

// invalidateCache 清除指定 microAppId 的缓存
func (s *microApp) invalidateCache(microAppId string) {
	if s.InfoCache != nil {
		s.InfoCache.Delete(microAppId)
	}
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
	Page            int    `json:"page"`            // 页码
	Limit           int    `json:"limit"`           // 每页数量
	Order           string `json:"order"`           // 排序（如 "download_count desc"）
	CategoryId      uint   `json:"categoryId"`      // 分类ID（0表示不筛选）
	Keyword         string `json:"keyword"`         // 关键词搜索
	Status          *int   `json:"status"`          // 状态（可选，默认为1-上架）
	Lang            string `json:"lang"`            // 语言（可选，暂不支持）
	OnlyWithVersion bool   `json:"onlyWithVersion"` // 是否只返回有最新审核通过版本的应用
}

// GetList 获取微应用列表（公开接口）
// 参数：
//   - db: 数据库连接
//   - opts: 查询选项
//
// 返回：
//   - 微应用列表、总数、错误
func (s *microApp) GetList(db *gorm.DB, opts GetListOptions) ([]models.MicroAppListItem, int64, error) {
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

	// 构建回退语言列表：en 开头的语言作为兜底
	fallbackLangs := []string{}
	if opts.Lang != "" && !strings.HasPrefix(opts.Lang, "en") {
		fallbackLangs = append(fallbackLangs, "en-US")
	}

	queryOpts := models.MicroAppListWithLangQueryOpts{
		MicroAppListQueryOpts: models.MicroAppListQueryOpts{
			Page:             opts.Page,
			Limit:            opts.Limit,
			Status:           &status,
			CategoryId:       categoryId,
			KeyWord:          opts.Keyword,
			SortBy:           sortBy,
			SortOrder:        sortOrder,
			IncludeDeveloper: true, // 包含开发者信息
			OnlyWithVersion:  opts.OnlyWithVersion,
		},
		Lang:          opts.Lang,
		FallbackLangs: fallbackLangs,
	}

	return m.GetAppListWithLang(db, queryOpts)
}

// BatchGetMicroAppInfo 批量获取微应用信息（带缓存）
// 1. 先从缓存获取已有的
// 2. 未命中的再批量查库（包括开发者信息）
// 3. 查库结果写入缓存
func (s *microApp) BatchGetMicroAppInfo(db *gorm.DB, microAppIds []string) (map[string]MicroAppInfoCache, error) {
	result := make(map[string]MicroAppInfoCache, len(microAppIds))

	if s.InfoCache == nil {
		// 缓存未初始化，降级为直接查库
		return s.batchGetFromDB(db, microAppIds)
	}

	// 1. 读缓存
	var missIds []string
	for _, id := range microAppIds {
		if info, ok := s.InfoCache.Get(id); ok {
			result[id] = info
		} else {
			missIds = append(missIds, id)
		}
	}

	// 2. 缓存未命中，批量查库
	if len(missIds) > 0 {
		dbResults, err := s.batchGetFromDB(db, missIds)
		if err != nil {
			return nil, err
		}
		for id, info := range dbResults {
			result[id] = info
			s.InfoCache.SetDefault(id, info) // 写入缓存
		}
	}

	return result, nil
}

// batchGetFromDB 从数据库批量获取微应用信息（含开发者信息）
func (s *microApp) batchGetFromDB(db *gorm.DB, microAppIds []string) (map[string]MicroAppInfoCache, error) {
	result := make(map[string]MicroAppInfoCache, len(microAppIds))

	// 批量查询微应用
	var apps []models.MicroApp
	if err := db.Where("micro_app_id IN ?", microAppIds).Find(&apps).Error; err != nil {
		return nil, err
	}

	if len(apps) == 0 {
		return result, nil
	}

	// 收集开发者ID，批量查询开发者信息
	developerIdSet := make(map[uint]bool)
	for _, app := range apps {
		developerIdSet[app.DeveloperId] = true
	}
	developerIds := make([]uint, 0, len(developerIdSet))
	for id := range developerIdSet {
		developerIds = append(developerIds, id)
	}

	var developers []models.Developer
	if err := db.Where("id IN ?", developerIds).Find(&developers).Error; err != nil {
		return nil, err
	}

	// 构建开发者信息map
	developerMap := make(map[uint]models.Developer, len(developers))
	for _, dev := range developers {
		developerMap[dev.ID] = dev
	}

	// 组装结果
	for _, app := range apps {
		dev := developerMap[app.DeveloperId]
		result[app.MicroAppId] = MicroAppInfoCache{
			MicroAppId:     app.MicroAppId,
			AppIcon:        app.AppIcon,
			ChargeType:     app.ChargeType,
			Points:         app.Points,
			DeveloperId:    app.DeveloperId,
			DeveloperName:  dev.Name,
			DeveloperName2: dev.DeveloperName,
		}
	}

	return result, nil
}
