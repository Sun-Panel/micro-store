package models

import (
	"fmt"
	"strconv"
	"time"

	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

// 自定义代码-线上状态
const (
	CustomCodeStatusDraft    = -1 // 草稿（保存未提交）
	CustomCodeStatusPending  = 0  // 待审核（首次提交，从未上线）
	CustomCodeStatusOnline   = 1  // 已上线
	CustomCodeStatusRejected = 2  // 已拒绝（从未上线）
	CustomCodeStatusOffline  = 3  // 已下架
)

// 自定义代码-审核快照状态
const (
	CustomCodeReviewStatusDraft    = -1 // 草稿（保存未提交/已撤回）
	CustomCodeReviewStatusPending  = 0  // 待审核
	CustomCodeReviewStatusApproved = 1  // 已通过
	CustomCodeReviewStatusRejected = 2  // 已拒绝
)

// 代码类型
const (
	CodeTypeJS     = 1 // 自定义JS
	CodeTypeCSS    = 2 // 自定义CSS
	CodeTypeFooter = 3 // 自定义页脚
)

// 适用的客户端版本
const (
	ClientVersionV1 = 1
	ClientVersionV2 = 2
)

// CustomCodeBaseInfo 主表与快照表共用的内容字段
// 具体内容由 custom_code_block（代码片块）承载
type CustomCodeBaseInfo struct {
	Title       string               `gorm:"type:varchar(100);not null" json:"title"`       // 标题
	Description string               `gorm:"type:varchar(500);not null" json:"description"` // 说明
	Keywords    datatype.StringArray `gorm:"type:varchar(500)" json:"keywords"`             // 关键词
	IsOriginal  bool                 `gorm:"type:tinyint(1);not null" json:"isOriginal"`
	SourceNote  string               `gorm:"type:varchar(500)" json:"sourceNote"`        // 来源说明（非原创必填）
	Versions    datatype.IntArray    `gorm:"type:varchar(50);not null" json:"versions"`  // 适用客户端版本：1-v1 2-v2
	CodeTypes   datatype.IntArray    `gorm:"type:varchar(50);not null" json:"codeTypes"` // 代码类型，由块自动汇总
}

// CustomCode 自定义代码主表（只存储已通过并在线的内容）
type CustomCode struct {
	BaseModel
	CustomCodeBaseInfo `gorm:"embedded"`

	AuthorId        uint       `gorm:"type:int(11);not null;index" json:"authorId"`            // 作者userID
	Status          int        `gorm:"type:tinyint(2);not null;default:0;index" json:"status"` // 线上状态
	ReadCount       int        `gorm:"type:int(11);not null;default:0;index" json:"readCount"`
	PublishedAt     *time.Time `gorm:"type:datetime" json:"publishedAt"`       // 作者发布/修改时间
	CurrentReviewId uint       `gorm:"type:int(11)" json:"currentReviewId"`    // 当前生效的快照ID
	OfflineReason   string     `gorm:"type:varchar(500)" json:"offlineReason"` // 下架原因
	UniqueKey       string     `gorm:"type:varchar(120)" json:"uniqueKey"`     // 自定义代码唯一标识：开发者标识-后缀
}

// CustomCodeListItem 列表项（附带作者展示名）
type CustomCodeListItem struct {
	CustomCode
	AuthorName string `gorm:"-" json:"authorName"` // 开发者名优先，无则用户名
}

// 表名
func (CustomCode) TableName() string {
	return "custom_code"
}

// CustomCodeQueryOptions 前台/列表查询选项
type CustomCodeQueryOptions struct {
	Page      int
	Limit     int
	Keyword   string
	SortBy    string // 排序字段：id, read_count, published_at, created_at
	SortOrder string // 排序方式：asc, desc
	Version   int    // 适用版本过滤：0-全部 1-v1 2-v2
	CodeType  int    // 代码类型过滤：0-全部 1-JS 2-CSS 3-页脚
}

// 允许的排序字段，避免拼接 SQL 带来注入风险
var customCodeSortFields = map[string]string{
	"id":           "custom_code.id",
	"read_count":   "custom_code.read_count",
	"published_at": "custom_code.published_at",
	"created_at":   "custom_code.created_at",
}

func normalizePage(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 10
	}
	return page, limit
}

// 创建
func (m *CustomCode) Create(db *gorm.DB, item *CustomCode) error {
	return db.Create(item).Error
}

// BackfillCustomCodeDefaults 回填历史数据的默认值
// 迁移前已存在的行 versions / code_types 可能为空，默认补上 v2 与空类型
func BackfillCustomCodeDefaults(db *gorm.DB) {
	// 每个操作使用独立会话（NewDB 彻底隔离语句），避免复用同一 *gorm.DB 语句导致
	// WHERE 条件 / 表名跨操作累积、错乱（例如 GetDeveloperName 的 developer 表泄漏到 custom_code 的更新）
	backfillColumn(db, &CustomCode{}, "versions", "[2]")
	backfillColumn(db, &CustomCode{}, "code_types", "[]")
	backfillColumn(db, &CustomCodeReview{}, "versions", "[2]")
	backfillColumn(db, &CustomCodeReview{}, "code_types", "[]")

	// 历史代码片块 only_id 为空时，用主键生成稳定唯一标识（b + id）；单条批量更新即可，避免按行循环
	if err := db.Session(&gorm.Session{NewDB: true}).Table("custom_code_block").
		Where("only_id IS NULL OR only_id = ''").
		Update("only_id", gorm.Expr("CONCAT('b', id)")).Error; err != nil {
		fmt.Printf("[backfill] 回填 only_id 失败: %v\n", err)
	}

	// 历史自定义代码回填唯一标识并建立整表唯一索引
	backfillCustomCodeUniqueKeys(db)
}

// backfillColumn 对指定表指定列做历史空值回填（空或 ” 时写入默认值），失败仅记录不阻断启动
func backfillColumn(db *gorm.DB, model interface{}, column, value string) {
	if err := db.Session(&gorm.Session{NewDB: true}).Model(model).
		Where(column+" IS NULL OR "+column+" = ''").
		Update(column, value).Error; err != nil {
		fmt.Printf("[backfill] 回填 %s 失败: %v\n", column, err)
	}
}

// backfillCustomCodeUniqueKeys 历史自定义代码回填唯一标识（开发者标识-<id>），并幂等建立整表唯一索引
// 用单条 UPDATE ... LEFT JOIN developer 一次完成，避免按行循环查询开发者表（N+1）与逐行更新
func backfillCustomCodeUniqueKeys(db *gorm.DB) {
	// 逻辑等价于 GetDeveloperName：有 developer 记录取 developer_name，否则回退 u<author_id>
	// 仅处理未软删除的行（与 GORM 软删除过滤一致）；一条语句搞定全部空值行
	if err := db.Session(&gorm.Session{NewDB: true}).
		Exec(`UPDATE custom_code cc
			LEFT JOIN developer d ON d.user_id = cc.author_id
			SET cc.unique_key = CONCAT(COALESCE(NULLIF(d.developer_name, ''), CONCAT('u', cc.author_id)), '-', cc.id)
			WHERE (cc.unique_key IS NULL OR cc.unique_key = '') AND cc.deleted_at IS NULL`).Error; err != nil {
		fmt.Printf("[backfill] 回填 unique_key 失败: %v\n", err)
	}
	// 回填后再建唯一索引。MySQL 不支持 CREATE UNIQUE INDEX IF NOT EXISTS，
	// 故先查 information_schema 判断索引是否已存在；已存在则跳过，避免反复打印 Duplicate key name 报错
	if !uniqueIndexExists(db, "custom_code", "uq_custom_code_unique_key") {
		if err := db.Session(&gorm.Session{NewDB: true}).
			Exec("CREATE UNIQUE INDEX uq_custom_code_unique_key ON custom_code (unique_key)").Error; err != nil {
			fmt.Printf("[backfill] 创建唯一索引失败: %v\n", err)
		}
	}
}

// uniqueIndexExists 判断指定表是否已存在某索引（MySQL 通过 information_schema.STATISTICS 查询）
func uniqueIndexExists(db *gorm.DB, table, indexName string) bool {
	var count int64
	if err := db.Session(&gorm.Session{NewDB: true}).
		Raw("SELECT COUNT(1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND INDEX_NAME = ?", table, indexName).
		Scan(&count).Error; err != nil {
		// 查询失败不阻断启动，回退为“不存在”，尝试建索引（失败也仅记录日志）
		return false
	}
	return count > 0
}

// GetAuthorName 获取作者展示名：开发者名优先，其次开发者标识，最后用户名
func GetAuthorName(db *gorm.DB, userId uint) string {
	var developer Developer
	if err := db.Where("user_id = ?", userId).First(&developer).Error; err == nil {
		if developer.Name != "" {
			return developer.Name
		}
		if developer.DeveloperName != "" {
			return developer.DeveloperName
		}
	}

	var user User
	if err := db.Where("id = ?", userId).First(&user).Error; err == nil && user.Name != "" {
		return user.Name
	}

	return ""
}

// GetDeveloperName 获取开发者标识（纯英文，多词用-分割）；无开发者记录则回退 u<id>
func GetDeveloperName(db *gorm.DB, userId uint) string {
	var dev Developer
	if err := db.Where("user_id = ?", userId).First(&dev).Error; err == nil && dev.DeveloperName != "" {
		return dev.DeveloperName
	}
	return fmt.Sprintf("u%d", userId)
}

// 根据ID获取
func (m *CustomCode) GetById(db *gorm.DB, id uint) (CustomCode, error) {
	var info CustomCode
	err := db.Where("id = ?", id).First(&info).Error
	return info, err
}

// 获取前台列表（仅已上线）
func (m *CustomCode) GetOnlineList(db *gorm.DB, opts CustomCodeQueryOptions) ([]CustomCodeListItem, int64, error) {
	page, limit := normalizePage(opts.Page, opts.Limit)

	query := db.Model(&CustomCode{}).Where("custom_code.status = ?", CustomCodeStatusOnline)
	if opts.Keyword != "" {
		kw := "%" + opts.Keyword + "%"
		query = query.Where(
			"custom_code.title LIKE ? OR custom_code.description LIKE ? OR custom_code.keywords LIKE ?",
			kw, kw, kw,
		)
	}
	// versions / code_types 以 JSON 数组字符串存储，且取值均为个位数，用 LIKE 匹配即可兼容 MySQL 与 SQLite
	if opts.Version > 0 {
		query = query.Where("custom_code.versions LIKE ?", "%"+strconv.Itoa(opts.Version)+"%")
	}
	if opts.CodeType > 0 {
		query = query.Where("custom_code.code_types LIKE ?", "%"+strconv.Itoa(opts.CodeType)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "custom_code.published_at DESC"
	if field, ok := customCodeSortFields[opts.SortBy]; ok {
		sortOrder := "DESC"
		if opts.SortOrder == "asc" {
			sortOrder = "ASC"
		}
		order = field + " " + sortOrder
	}

	// 列表不需要正文，避免传输大字段
	var list []CustomCodeListItem
	err := query.Select("custom_code.id, custom_code.title, custom_code.description, custom_code.keywords, custom_code.is_original, custom_code.source_note, custom_code.versions, custom_code.code_types, custom_code.author_id, custom_code.status, custom_code.read_count, custom_code.published_at, custom_code.unique_key, custom_code.created_at").
		Order(order).
		Offset((page - 1) * limit).Limit(limit).
		Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	attachAuthorNames(db, list)
	return list, total, nil
}

// 获取作者的帖子列表（含待审/草稿标记由上层补充）
func (m *CustomCode) GetListByAuthorId(db *gorm.DB, authorId uint, page, limit int, keyword string) ([]CustomCodeListItem, int64, error) {
	page, limit = normalizePage(page, limit)

	query := db.Model(&CustomCode{}).Where("custom_code.author_id = ?", authorId)
	if keyword != "" {
		kw := "%" + keyword + "%"
		query = query.Where("custom_code.title LIKE ? OR custom_code.description LIKE ?", kw, kw)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []CustomCodeListItem
	err := query.Select("custom_code.id, custom_code.title, custom_code.description, custom_code.keywords, custom_code.is_original, custom_code.source_note, custom_code.versions, custom_code.code_types, custom_code.author_id, custom_code.status, custom_code.read_count, custom_code.published_at, custom_code.unique_key, custom_code.created_at").
		Order("custom_code.id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	attachAuthorNames(db, list)
	return list, total, nil
}

// 阅读量+1
func (m *CustomCode) IncrReadCount(db *gorm.DB, id uint) error {
	return db.Model(&CustomCode{}).Where("id = ?", id).
		UpdateColumn("read_count", gorm.Expr("read_count + 1")).Error
}

// 下架（作者或平台）
func (m *CustomCode) SetOffline(db *gorm.DB, id uint, reason string) error {
	return db.Model(&CustomCode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         CustomCodeStatusOffline,
		"offline_reason": reason,
	}).Error
}

// UpdateStatus 仅更新线上状态字段（不改动内容），供撤回/状态回退等纯状态变更复用
func (m *CustomCode) UpdateStatus(db *gorm.DB, id uint, status int) error {
	return db.Model(&CustomCode{}).Where("id = ?", id).Update("status", status).Error
}

// Submit 作者提交：进入待审并同步由代码块汇总出的内容字段（版本、代码类型）
func (m *CustomCode) Submit(db *gorm.DB, id uint, versions, codeTypes datatype.IntArray) error {
	return db.Model(&CustomCode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     CustomCodeStatusPending,
		"versions":   versions,
		"code_types": codeTypes,
	}).Error
}

// 删除（软删除）
func (m *CustomCode) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&CustomCode{}, id).Error
}

// GetAuthorNames 批量获取作者展示名：开发者名优先，其次开发者标识，最后用户名
func GetAuthorNames(db *gorm.DB, userIds []uint) map[uint]string {
	nameMap := make(map[uint]string, len(userIds))
	if len(userIds) == 0 {
		return nameMap
	}

	var developers []Developer
	if err := db.Where("user_id IN ?", userIds).Find(&developers).Error; err == nil {
		for _, d := range developers {
			name := d.Name
			if name == "" {
				name = d.DeveloperName
			}
			if name != "" {
				nameMap[d.UserId] = name
			}
		}
	}

	var users []User
	if err := db.Where("id IN ?", userIds).Find(&users).Error; err == nil {
		for _, u := range users {
			if _, ok := nameMap[u.ID]; !ok && u.Name != "" {
				nameMap[u.ID] = u.Name
			}
		}
	}

	return nameMap
}

// attachAuthorNames 统一填充列表的作者展示名
func attachAuthorNames(db *gorm.DB, list []CustomCodeListItem) {
	if len(list) == 0 {
		return
	}

	userIds := make([]uint, 0, len(list))
	seen := make(map[uint]bool, len(list))
	for _, item := range list {
		if !seen[item.AuthorId] {
			seen[item.AuthorId] = true
			userIds = append(userIds, item.AuthorId)
		}
	}

	nameMap := GetAuthorNames(db, userIds)
	for i := range list {
		list[i].AuthorName = nameMap[list[i].AuthorId]
	}
}
