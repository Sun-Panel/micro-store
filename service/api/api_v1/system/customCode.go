package system

import (
	"strconv"
	"sync"
	"time"

	"sun-panel/api/api_v1/common/apiData/customCodeApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type CustomCodeApi struct{}

// 阅读量去重缓存（24小时）
var (
	readCountCacheOnce   sync.Once
	readCountCache       cache.Cacher[bool]
	readCountExpiredTime = 24 * time.Hour
)

func getReadCountCache() cache.Cacher[bool] {
	readCountCacheOnce.Do(func() {
		readCountCache = global.NewCache[bool](readCountExpiredTime, time.Hour, "custom_code_read_count")
	})
	return readCountCache
}

// GetList 前台列表（仅已上线）
func (a *CustomCodeApi) GetList(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeListReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	list, total, err := (&models.CustomCode{}).GetOnlineList(global.Db, models.CustomCodeQueryOptions{
		Page:      req.Page,
		Limit:     req.Limit,
		Keyword:   req.Keyword,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Version:   req.Version,
		CodeType:  req.CodeType,
	})
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// Get 前台详情（含阅读量统计）
func (a *CustomCodeApi) Get(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	info, err := (&models.CustomCode{}).GetById(global.Db, req.Id)
	if err != nil || info.Status != models.CustomCodeStatusOnline {
		apiReturn.Error(c, "内容不存在或已下架")
		return
	}

	// 阅读量统计：同一访客24小时内只计一次
	go addReadCount(req.Id, c.ClientIP())

	// 线上块（review_id = 0）
	blocks, err := (&models.CustomCodeBlock{}).GetListByCustomCodeId(global.Db, info.ID, 0)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	respBlocks := make([]customCodeApiStructs.CustomCodeBlockResp, 0, len(blocks))
	for _, block := range blocks {
		respBlocks = append(respBlocks, customCodeApiStructs.CustomCodeBlockResp{
			Id:       block.ID,
			CodeType: block.CodeType,
			Title:    block.Title,
			Note:     block.Note,
			Code:     block.Code,
			Images:   block.Images,
			Sort:     block.Sort,
			OnlyId:   block.OnlyId,
		})
	}

	resp := customCodeApiStructs.CustomCodeDetailResp{
		Id:          info.ID,
		Title:       info.Title,
		Description: info.Description,
		Keywords:    info.Keywords,
		IsOriginal:  info.IsOriginal,
		SourceNote:  info.SourceNote,
		Versions:    info.Versions,
		CodeTypes:   info.CodeTypes,
		Blocks:      respBlocks,
		AuthorId:    info.AuthorId,
		AuthorName:  models.GetAuthorName(global.Db, info.AuthorId),
		ReadCount:   info.ReadCount,
		PublishedAt: info.PublishedAt,
	}

	apiReturn.SuccessData(c, resp)
}

// addReadCount 阅读量+1（带去重）
func addReadCount(id uint, clientIp string) {
	// 异步执行，兜底防止缓存/DB 异常 panic 拖垮整个进程
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorln("custom code read count panic recovered:", r)
		}
	}()

	cacheObj := getReadCountCache()
	if cacheObj == nil {
		// 缓存不可用时退化为直接计数
		global.Db.Model(&models.CustomCode{}).Where("id = ?", id).
			UpdateColumn("read_count", gorm.Expr("read_count + 1"))
		return
	}

	key := clientIp + ":" + strconv.FormatUint(uint64(id), 10)
	if _, exist := cacheObj.Get(key); exist {
		return
	}
	cacheObj.Set(key, true, readCountExpiredTime)

	if err := (&models.CustomCode{}).IncrReadCount(global.Db, id); err != nil {
		global.Logger.Errorln("custom code read count incr failed:", err.Error())
	}
}
