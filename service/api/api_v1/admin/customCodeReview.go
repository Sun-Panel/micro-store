package admin

import (
	"strings"
	"time"

	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/customCodeApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"
	"sun-panel/models/datatype"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CustomCodeReviewApi struct{}

// CustomCodeReviewListItem 审核列表项
type CustomCodeReviewListItem struct {
	Id           uint                             `json:"id"`
	CustomCodeId uint                             `json:"customCodeId"`
	Title        string                           `json:"title"`
	Description  string                           `json:"description"`
	Keywords     datatype.StringArray             `json:"keywords"`
	IsOriginal   bool                             `json:"isOriginal"`
	SourceNote   string                           `json:"sourceNote"`
	CodeTypes    datatype.IntArray                `json:"codeTypes"`
	Status       int                              `json:"status"`
	AuthorName   string                           `json:"authorName"`
	CreatedAt    time.Time                        `json:"createTime"`
	OnlineStatus int                              `json:"onlineStatus"` // 主表当前状态（判断是否首次提交）
	MachineAudit *datatype.CustomCodeMachineAudit `json:"machineAudit"` // 机器预审结果（含未通过原因）
}

// GetList 待审核列表（审核员/管理员）
func (a *CustomCodeReviewApi) GetList(c *gin.Context) {
	req := commonApiStructs.RequestPage{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	list, total, err := (&models.CustomCodeReview{}).GetPendingList(global.Db, req.Page, req.Limit, req.Keyword)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 作者名与主表状态
	codeIds := make([]uint, 0, len(list))
	for _, item := range list {
		codeIds = append(codeIds, item.CustomCodeId)
	}

	authorMap := make(map[uint]string, len(codeIds))
	statusMap := make(map[uint]int, len(codeIds))
	if len(codeIds) > 0 {
		var mains []models.CustomCode
		global.Db.Select("id, author_id, status").Where("id IN ?", codeIds).Find(&mains)

		userIds := make([]uint, 0, len(mains))
		for _, m := range mains {
			statusMap[m.ID] = m.Status
			userIds = append(userIds, m.AuthorId)
		}
		nameMap := models.GetAuthorNames(global.Db, userIds)
		for _, m := range mains {
			authorMap[m.ID] = nameMap[m.AuthorId]
		}
	}

	respList := make([]CustomCodeReviewListItem, 0, len(list))
	for _, item := range list {
		respList = append(respList, CustomCodeReviewListItem{
			Id:           item.ID,
			CustomCodeId: item.CustomCodeId,
			Title:        item.Title,
			Description:  item.Description,
			Keywords:     item.Keywords,
			IsOriginal:   item.IsOriginal,
			SourceNote:   item.SourceNote,
			CodeTypes:    item.CodeTypes,
			Status:       item.Status,
			AuthorName:   authorMap[item.CustomCodeId],
			CreatedAt:    item.CreatedAt,
			OnlineStatus: statusMap[item.CustomCodeId],
			MachineAudit: item.MachineAudit,
		})
	}

	apiReturn.SuccessListData(c, respList, total)
}

// GetInfo 审核详情：待审快照 + 当前线上版本（对比用）
func (a *CustomCodeReviewApi) GetInfo(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mReview := &models.CustomCodeReview{}
	review, err := mReview.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "审核记录不存在")
		return
	}

	main, mainErr := (&models.CustomCode{}).GetById(global.Db, review.CustomCodeId)

	mBlock := &models.CustomCodeBlock{}
	reviewBlocks, _ := mBlock.GetListByCustomCodeId(global.Db, review.CustomCodeId, review.ID)

	resp := gin.H{
		"review": gin.H{
			"id":           review.ID,
			"customCodeId": review.CustomCodeId,
			"title":        review.Title,
			"description":  review.Description,
			"keywords":     review.Keywords,
			"isOriginal":   review.IsOriginal,
			"sourceNote":   review.SourceNote,
			"versions":     review.Versions,
			"codeTypes":    review.CodeTypes,
			"blocks":       reviewBlocks,
			"status":       review.Status,
			"reviewNote":   review.ReviewNote,
			"createTime":   review.CreatedAt,
			"machineAudit": review.MachineAudit,
		},
	}

	if mainErr == nil {
		onlineBlocks, _ := mBlock.GetListByCustomCodeId(global.Db, main.ID, 0)
		resp["online"] = gin.H{
			"id":          main.ID,
			"title":       main.Title,
			"description": main.Description,
			"keywords":    main.Keywords,
			"isOriginal":  main.IsOriginal,
			"sourceNote":  main.SourceNote,
			"versions":    main.Versions,
			"codeTypes":   main.CodeTypes,
			"blocks":      onlineBlocks,
			"status":      main.Status,
			"publishedAt": main.PublishedAt,
			"authorName":  models.GetAuthorName(global.Db, main.AuthorId),
		}
	}

	apiReturn.SuccessData(c, resp)
}

// Approve 审核通过：快照内容覆盖主表并上线
func (a *CustomCodeReviewApi) Approve(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeReviewReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mReview := &models.CustomCodeReview{}
	review, err := mReview.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "审核记录不存在")
		return
	}
	if review.Status != models.CustomCodeReviewStatusPending {
		apiReturn.Error(c, "该记录不是待审核状态")
		return
	}

	if err := PublishApprovedReview(global.Db, review.CustomCodeId, review.ID, userInfo.ID, req.Note); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Reject 审核拒绝：拒绝原因必填，已上线内容不受影响
func (a *CustomCodeReviewApi) Reject(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeReviewReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	if strings.TrimSpace(req.Note) == "" {
		apiReturn.Error(c, "请填写拒绝原因")
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mReview := &models.CustomCodeReview{}
	review, err := mReview.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "审核记录不存在")
		return
	}
	if review.Status != models.CustomCodeReviewStatusPending {
		apiReturn.Error(c, "该记录不是待审核状态")
		return
	}

	if err := mReview.UpdateStatus(global.Db, review.ID, models.CustomCodeReviewStatusRejected, userInfo.ID, req.Note); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 首次提交被拒：主表标记为已拒绝（从未上线，前台不可见）
	main, err := (&models.CustomCode{}).GetById(global.Db, review.CustomCodeId)
	if err == nil && main.Status == models.CustomCodeStatusPending {
		global.Db.Model(&models.CustomCode{}).Where("id = ?", main.ID).
			Update("status", models.CustomCodeStatusRejected)
	}

	apiReturn.Success(c)
}

// GetHistory 某条内容的审核历史
func (a *CustomCodeReviewApi) GetHistory(c *gin.Context) {
	req := commonApiStructs.RequestPage{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	idReq := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&idReq, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	list, total, err := (&models.CustomCodeReview{}).GetHistoryByCustomCodeId(global.Db, idReq.Id, req.Page, req.Limit)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}
