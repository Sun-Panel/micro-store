package customCodeApiStructs

import (
	"time"

	"sun-panel/models/datatype"
)

// CustomCodeBlockReq 代码片块（提交用）
type CustomCodeBlockReq struct {
	CodeType int                  `json:"codeType"` // 1-JS 2-CSS 3-页脚
	Title    string               `json:"title"`    // 块标题（与说明至少填一个）
	Note     string               `json:"note"`     // 块说明（纯文本）
	Code     string               `json:"code"`     // 代码片段（必填）
	Images   datatype.StringArray `json:"images"`   // 预览图
	OnlyId   string               `json:"onlyId"`   // 块唯一标识（可空，由服务端生成；开发者可改）
}

// CustomCodeBlockResp 代码片块（返回用）
type CustomCodeBlockResp struct {
	Id       uint                 `json:"id"`
	CodeType int                  `json:"codeType"`
	Title    string               `json:"title"`
	Note     string               `json:"note"`
	Code     string               `json:"code"`
	Images   datatype.StringArray `json:"images"`
	Sort     int                  `json:"sort"`
	OnlyId   string               `json:"onlyId"`
}

// CustomCodeEditReq 新增/编辑（保存草稿或提交审核）
type CustomCodeEditReq struct {
	Id          uint                 `json:"id"` // 0 表示新增
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Keywords    datatype.StringArray `json:"keywords"`
	IsOriginal  bool                 `json:"isOriginal"`
	SourceNote  string               `json:"sourceNote"`
	Versions    datatype.IntArray    `json:"versions"`   // 适用客户端版本：1-v1 2-v2
	Blocks      []CustomCodeBlockReq `json:"blocks"`     // 代码片块（1~10）
	Submit      bool                 `json:"submit"`     // true-提交审核 false-保存草稿
	CustomName  string               `json:"customName"` // 唯一标识后缀（开发者标识-后缀 中的后缀；留空自动生成）
}

// CustomCodeIdReq 按ID操作的通用请求
type CustomCodeIdReq struct {
	Id uint `json:"id"`
}

// CustomCodeListReq 前台列表请求
type CustomCodeListReq struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Keyword   string `json:"keyword"`
	SortBy    string `json:"sortBy"`    // id, read_count, published_at, created_at
	SortOrder string `json:"sortOrder"` // asc, desc
	Version   int    `json:"version"`   // 适用版本过滤：0-全部 1-v1 2-v2
	CodeType  int    `json:"codeType"`  // 代码类型过滤：0-全部 1-JS 2-CSS 3-页脚
}

// CustomCodeOfflineReq 下架
type CustomCodeOfflineReq struct {
	Id     uint   `json:"id"`
	Reason string `json:"reason"`
}

// CustomCodeReviewReq 审核操作
type CustomCodeReviewReq struct {
	Id   uint   `json:"id"`
	Note string `json:"note"` // 拒绝时必填
}

// CustomCodeListItemResp 列表项（附带快照状态）
type CustomCodeListItemResp struct {
	Id          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Keywords    datatype.StringArray `json:"keywords"`
	IsOriginal  bool                 `json:"isOriginal"`
	SourceNote  string               `json:"sourceNote"`
	Versions    datatype.IntArray    `json:"versions"`
	CodeTypes   datatype.IntArray    `json:"codeTypes"`
	AuthorId    uint                 `json:"authorId"`
	AuthorName  string               `json:"authorName"`
	UniqueKey   string               `json:"uniqueKey"`
	Status      int                  `json:"status"`
	ReadCount   int                  `json:"readCount"`
	PublishedAt *time.Time           `json:"publishedAt"`
	CreatedAt   time.Time            `json:"createTime"`

	ReviewStatus *int   `json:"reviewStatus"`
	ReviewNote   string `json:"reviewNote"`
}

// CustomCodeInfoResp 编辑/详情信息
type CustomCodeInfoResp struct {
	Id            uint                  `json:"id"`
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	Keywords      datatype.StringArray  `json:"keywords"`
	IsOriginal    bool                  `json:"isOriginal"`
	SourceNote    string                `json:"sourceNote"`
	Versions      datatype.IntArray     `json:"versions"`
	CodeTypes     datatype.IntArray     `json:"codeTypes"`
	Blocks        []CustomCodeBlockResp `json:"blocks"`
	Status        int                   `json:"status"`
	ReadCount     int                   `json:"readCount"`
	PublishedAt   *time.Time            `json:"publishedAt"`
	AuthorId      uint                  `json:"authorId"`
	AuthorName    string                `json:"authorName"`
	DeveloperName string                `json:"developerName"` // 作者开发者标识（唯一标识前缀，固定不可改）
	UniqueKey     string                `json:"uniqueKey"`     // 完整唯一标识：开发者标识-后缀

	CanEdit         bool   `json:"canEdit"`
	ReviewStatus    *int   `json:"reviewStatus"`
	ReviewNote      string `json:"reviewNote"`
	PendingReviewId *uint  `json:"pendingReviewId"`
}

// CustomCodeDetailResp 前台详情
type CustomCodeDetailResp struct {
	Id                  uint                  `json:"id"`
	Title               string                `json:"title"`
	Description         string                `json:"description"`
	Keywords            datatype.StringArray  `json:"keywords"`
	IsOriginal          bool                  `json:"isOriginal"`
	SourceNote          string                `json:"sourceNote"`
	Versions            datatype.IntArray     `json:"versions"`
	CodeTypes           datatype.IntArray     `json:"codeTypes"`
	Blocks              []CustomCodeBlockResp `json:"blocks"`
	AuthorId            uint                  `json:"authorId"`
	AuthorName          string                `json:"authorName"`
	UniqueKey           string                `json:"uniqueKey"`
	AuthorRewardContent string                `json:"authorRewardContent"`
	ReadCount           int                   `json:"readCount"`
	PublishedAt         *time.Time            `json:"publishedAt"`
}
