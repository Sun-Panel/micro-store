package admin

import (
	"crypto/rand"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/customCodeApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz/event"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/file"
	"sun-panel/models"
	"sun-panel/models/datatype"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CustomCodeApi struct{}

// 最新快照状态（用于列表展示待审/草稿/拒绝）
type customCodeReviewState struct {
	CustomCodeId uint   `gorm:"column:custom_code_id" json:"customCodeId"`
	ReviewId     uint   `gorm:"column:review_id" json:"reviewId"`
	Status       int    `gorm:"column:status" json:"status"`
	ReviewNote   string `gorm:"column:review_note" json:"reviewNote"`
}

// getLatestReviewStates 批量查询帖子最新一次快照的状态
func getLatestReviewStates(ids []uint) map[uint]customCodeReviewState {
	result := make(map[uint]customCodeReviewState, len(ids))
	if len(ids) == 0 {
		return result
	}

	rows := []customCodeReviewState{}
	err := global.Db.Raw(`
		SELECT r.custom_code_id AS custom_code_id, r.id AS review_id, r.status AS status, r.review_note AS review_note
		FROM custom_code_review r
		INNER JOIN (
			SELECT custom_code_id, MAX(id) AS max_id FROM custom_code_review WHERE custom_code_id IN ? GROUP BY custom_code_id
		) t ON r.id = t.max_id`, ids).Scan(&rows).Error
	if err != nil {
		return result
	}

	for _, row := range rows {
		result[row.CustomCodeId] = row
	}
	return result
}

// summarizeCodeTypes 由块自动汇总代码类型（去重并排序）
func summarizeCodeTypes(blocks []customCodeApiStructs.CustomCodeBlockReq) datatype.IntArray {
	seen := make(map[int]bool, len(blocks))
	result := datatype.IntArray{}
	for _, block := range blocks {
		if !seen[block.CodeType] {
			seen[block.CodeType] = true
			result = append(result, block.CodeType)
		}
	}
	// 按 1-JS 2-CSS 3-页脚 的顺序输出
	sorted := datatype.IntArray{}
	for _, codeType := range []int{models.CodeTypeJS, models.CodeTypeCSS, models.CodeTypeFooter} {
		if seen[codeType] {
			sorted = append(sorted, codeType)
		}
	}
	return sorted
}

func blocksToModels(blocks []customCodeApiStructs.CustomCodeBlockReq) []models.CustomCodeBlock {
	list := make([]models.CustomCodeBlock, 0, len(blocks))
	for i, block := range blocks {
		list = append(list, models.CustomCodeBlock{
			Sort:     i,
			CodeType: block.CodeType,
			Title:    strings.TrimSpace(block.Title),
			Note:     block.Note,
			Code:     block.Code,
			Images:   block.Images,
			OnlyId:   block.OnlyId,
		})
	}
	return list
}

func blocksToResp(blocks []models.CustomCodeBlock) []customCodeApiStructs.CustomCodeBlockResp {
	list := make([]customCodeApiStructs.CustomCodeBlockResp, 0, len(blocks))
	for _, block := range blocks {
		list = append(list, customCodeApiStructs.CustomCodeBlockResp{
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
	return list
}

// onlyIdRegexp 唯一标识合法字符集（不含 -，避免与 {id}-{onlyId} 的连接符冲突）
var onlyIdRegexp = regexp.MustCompile(`^[A-Za-z0-9_]{1,40}$`)

// genOnlyId 生成块唯一标识（12 位十六进制）
func genOnlyId() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		// 极不可能失败，退化用时间戳
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", buf)
}

// normalizeBlockOnlyIds 为空补全、校验格式、校验同帖内唯一
func normalizeBlockOnlyIds(blocks []customCodeApiStructs.CustomCodeBlockReq) string {
	seen := make(map[string]int, len(blocks))
	for i, block := range blocks {
		id := strings.TrimSpace(block.OnlyId)
		if id == "" {
			id = genOnlyId()
		} else if !onlyIdRegexp.MatchString(id) {
			return fmt.Sprintf("第%d个片段的唯一标识格式不正确（仅限字母、数字、下划线，最长40位）", i+1)
		}
		if prev, ok := seen[id]; ok {
			return fmt.Sprintf("第%d个片段的唯一标识与第%d个重复：%s", i+1, prev+1, id)
		}
		seen[id] = i
		blocks[i].OnlyId = id
	}
	return ""
}

// validateCustomCodeReq 校验帖子与块
func validateCustomCodeReq(req *customCodeApiStructs.CustomCodeEditReq) string {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.SourceNote = strings.TrimSpace(req.SourceNote)

	switch {
	case req.Title == "":
		return "标题不能为空"
	case len([]rune(req.Title)) > 100:
		return "标题不能超过100个字符"
	case req.Description == "":
		return "简短描述不能为空"
	case len([]rune(req.Description)) > 500:
		return "说明不能超过500个字符"
	case len(req.Keywords) > 10:
		return "关键词不能超过10个"
	case !req.IsOriginal && req.SourceNote == "":
		return "非原创必须填写来源说明"
	case len(req.Versions) == 0:
		return "请至少选择一个适用版本"
	case len(req.Blocks) == 0:
		return "请至少添加一个代码片块"
	case len(req.Blocks) > models.CustomCodeMaxBlocks:
		return fmt.Sprintf("代码片块最多%d个", models.CustomCodeMaxBlocks)
	}

	validVersions := map[int]bool{models.ClientVersionV1: true, models.ClientVersionV2: true}
	for _, version := range req.Versions {
		if !validVersions[version] {
			return "适用版本不合法"
		}
	}

	validCodeTypes := map[int]bool{
		models.CodeTypeJS:     true,
		models.CodeTypeCSS:    true,
		models.CodeTypeFooter: true,
	}
	for i, block := range req.Blocks {
		switch {
		case !validCodeTypes[block.CodeType]:
			return fmt.Sprintf("第%d个块的代码类型不合法", i+1)
		case strings.TrimSpace(block.Title) == "" && strings.TrimSpace(block.Note) == "":
			return fmt.Sprintf("第%d个块的标题和说明至少填写一个", i+1)
		case strings.TrimSpace(block.Code) == "":
			return fmt.Sprintf("第%d个块的代码不能为空", i+1)
		case len(block.Images) > models.CustomCodeMaxImagesPerBlock:
			return fmt.Sprintf("第%d个块的预览图最多%d张", i+1, models.CustomCodeMaxImagesPerBlock)
		}
	}

	return ""
}

// GetMyList 作者：我的自定义代码
func (a *CustomCodeApi) GetMyList(c *gin.Context) {
	req := commonApiStructs.RequestPage{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	list, total, err := (&models.CustomCode{}).GetListByAuthorId(global.Db, userInfo.ID, req.Page, req.Limit, req.Keyword)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	ids := make([]uint, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.ID)
	}
	states := getLatestReviewStates(ids)

	respList := make([]customCodeApiStructs.CustomCodeListItemResp, 0, len(list))
	for _, item := range list {
		resp := customCodeApiStructs.CustomCodeListItemResp{
			Id:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Keywords:    item.Keywords,
			IsOriginal:  item.IsOriginal,
			SourceNote:  item.SourceNote,
			Versions:    item.Versions,
			CodeTypes:   item.CodeTypes,
			AuthorId:    item.AuthorId,
			AuthorName:  item.AuthorName,
			Status:      item.Status,
			ReadCount:   item.ReadCount,
			PublishedAt: item.PublishedAt,
			CreatedAt:   item.CreatedAt,
		}
		if state, ok := states[item.ID]; ok {
			status := state.Status
			resp.ReviewStatus = &status
			resp.ReviewNote = state.ReviewNote
		}
		respList = append(respList, resp)
	}

	apiReturn.SuccessListData(c, respList, total)
}

// GetInfo 作者：获取编辑信息（含块）
func (a *CustomCodeApi) GetInfo(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mCustomCode := &models.CustomCode{}
	info, err := mCustomCode.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "内容不存在")
		return
	}
	if info.AuthorId != userInfo.ID {
		apiReturn.Error(c, "无权访问该内容")
		return
	}

	resp := customCodeApiStructs.CustomCodeInfoResp{
		Id:          info.ID,
		Title:       info.Title,
		Description: info.Description,
		Keywords:    info.Keywords,
		IsOriginal:  info.IsOriginal,
		SourceNote:  info.SourceNote,
		Versions:    info.Versions,
		CodeTypes:   info.CodeTypes,
		Status:      info.Status,
		ReadCount:   info.ReadCount,
		PublishedAt: info.PublishedAt,
		AuthorId:    info.AuthorId,
		AuthorName:  models.GetAuthorName(global.Db, info.AuthorId),
		CanEdit:     true,
	}

	mBlock := &models.CustomCodeBlock{}
	mReview := &models.CustomCodeReview{}

	// 存在快照时以快照为准（草稿/待审/被拒的修改）
	if latest, err := mReview.GetLatestByCustomCodeId(global.Db, info.ID); err == nil && latest.ID != 0 {
		status := latest.Status
		resp.ReviewStatus = &status
		resp.ReviewNote = latest.ReviewNote
		resp.Title = latest.Title
		resp.Description = latest.Description
		resp.Keywords = latest.Keywords
		resp.IsOriginal = latest.IsOriginal
		resp.SourceNote = latest.SourceNote
		resp.Versions = latest.Versions

		if blocks, err := mBlock.GetListByCustomCodeId(global.Db, info.ID, latest.ID); err == nil {
			resp.Blocks = blocksToResp(blocks)
		}

		if latest.Status == models.CustomCodeReviewStatusPending {
			resp.CanEdit = false
			reviewId := latest.ID
			resp.PendingReviewId = &reviewId
		}
	} else if blocks, err := mBlock.GetListByCustomCodeId(global.Db, info.ID, 0); err == nil {
		resp.Blocks = blocksToResp(blocks)
	}

	if resp.Blocks == nil {
		resp.Blocks = []customCodeApiStructs.CustomCodeBlockResp{}
	}

	apiReturn.SuccessData(c, resp)
}

// Edit 新增/保存草稿/提交审核（含块）
func (a *CustomCodeApi) Edit(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeEditReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if msg := validateCustomCodeReq(&req); msg != "" {
		apiReturn.Error(c, msg)
		return
	}

	// 块唯一标识：空则补全、校验格式与同帖唯一
	if msg := normalizeBlockOnlyIds(req.Blocks); msg != "" {
		apiReturn.Error(c, msg)
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mReview := &models.CustomCodeReview{}
	mBlock := &models.CustomCodeBlock{}
	codeTypes := summarizeCodeTypes(req.Blocks)

	baseInfo := models.CustomCodeBaseInfo{
		Title:       req.Title,
		Description: req.Description,
		Keywords:    req.Keywords,
		IsOriginal:  req.IsOriginal,
		SourceNote:  req.SourceNote,
		Versions:    req.Versions,
		CodeTypes:   codeTypes,
	}

	// 审核模式：作者受信任 + 安全扫描通过 -> 自动审核，否则人工（均在提交时计算）
	modelBlocks := blocksToModels(req.Blocks)
	var machineAudit *datatype.CustomCodeMachineAudit
	auto := false
	if req.Submit {
		trusted, _ := ShouldAutoApprove(global.Db, userInfo.ID)
		scanPass, scanReasons := AutoReviewSecurityScan(modelBlocks)
		auto = DecideReviewMode(trusted, scanPass) == ReviewModeAuto
		// 机器预审结果：记录未通过原因，供人工审核快速定位风险点
		reasons := scanReasons
		if !trusted {
			reasons = append([]string{"作者未受信任（近一年内无审核通过记录），转人工审核"}, reasons...)
		}
		machineAudit = &datatype.CustomCodeMachineAudit{
			Auto:     auto,
			Trusted:  trusted,
			ScanPass: scanPass,
			Reasons:  reasons,
			Time:     time.Now(),
		}
	}

	// 新增
	if req.Id == 0 {
		mainStatus := models.CustomCodeStatusDraft
		reviewStatus := models.CustomCodeReviewStatusDraft
		if req.Submit {
			if auto {
				reviewStatus = models.CustomCodeReviewStatusApproved
			} else {
				mainStatus = models.CustomCodeStatusPending
				reviewStatus = models.CustomCodeReviewStatusPending
			}
		}

		main := models.CustomCode{
			CustomCodeBaseInfo: baseInfo,
			AuthorId:           userInfo.ID,
			Status:             mainStatus,
		}
		if err := (&models.CustomCode{}).Create(global.Db, &main); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		review := models.CustomCodeReview{
			CustomCodeBaseInfo: baseInfo,
			CustomCodeId:       main.ID,
			Status:             reviewStatus,
			MachineAudit:       machineAudit,
		}
		if err := mReview.Create(global.Db, &review); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		if err := mBlock.ReplaceAll(global.Db, main.ID, review.ID, modelBlocks); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		// 非首次提交（作者受信任）且通过安全扫描：自动审核通过，直接发布上线
		if req.Submit && auto {
			if err := PublishApprovedReview(global.Db, main.ID, review.ID, 0, "自动审核通过"); err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return
			}
		}

		// 提交完成：发出领域事件，由后台 worker 异步处理后续副作用（人工路径通知审核员/审计等）
		if req.Submit {
			event.Publish(event.CustomCodeSubmitted, event.CustomCodeSubmittedPayload{
				CustomCodeID: main.ID,
				ReviewID:     review.ID,
				AuthorID:     userInfo.ID,
				Auto:         auto,
			})
		}

		apiReturn.SuccessData(c, gin.H{"id": main.ID})
		return
	}

	// 编辑
	mCustomCode := &models.CustomCode{}
	main, err := mCustomCode.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "内容不存在")
		return
	}
	if main.AuthorId != userInfo.ID {
		apiReturn.Error(c, "无权修改该内容")
		return
	}

	// 待审核期间不可编辑
	if pending, err := mReview.GetPendingByCustomCodeId(global.Db, main.ID); err == nil && pending.ID != 0 {
		apiReturn.Error(c, "审核中不可修改，请先撤回审核")
		return
	}

	reviewStatus := models.CustomCodeReviewStatusDraft
	if req.Submit {
		if auto {
			reviewStatus = models.CustomCodeReviewStatusApproved
		} else {
			reviewStatus = models.CustomCodeReviewStatusPending
		}
	}

	// 复用已有草稿/被拒快照，没有则新建
	needCreate := true
	var reviewId uint
	if latest, err := mReview.GetLatestByCustomCodeId(global.Db, main.ID); err == nil && latest.ID != 0 {
		if latest.Status == models.CustomCodeReviewStatusDraft || latest.Status == models.CustomCodeReviewStatusRejected {
			reviewId = latest.ID
			needCreate = false
		}
	}

	if needCreate {
		review := models.CustomCodeReview{
			CustomCodeBaseInfo: baseInfo,
			CustomCodeId:       main.ID,
			Status:             reviewStatus,
			MachineAudit:       machineAudit,
		}
		if err := mReview.Create(global.Db, &review); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		reviewId = review.ID
	} else {
		if err := global.Db.Model(&models.CustomCodeReview{}).Where("id = ?", reviewId).Updates(map[string]interface{}{
			"title":         baseInfo.Title,
			"description":   baseInfo.Description,
			"keywords":      baseInfo.Keywords,
			"is_original":   baseInfo.IsOriginal,
			"source_note":   baseInfo.SourceNote,
			"versions":      baseInfo.Versions,
			"code_types":    baseInfo.CodeTypes,
			"status":        reviewStatus,
			"review_note":   "",
			"reviewer_id":   0,
			"review_time":   nil,
			"machine_audit": machineAudit,
		}).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	// 块写入快照
	if err := mBlock.ReplaceAll(global.Db, main.ID, reviewId, modelBlocks); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 提交：人工路径置待审并同步由块汇总出的版本与代码类型；
	// 自动路径由后续 PublishApprovedReview 统一发布（含标题/版本/代码类型/状态全量同步），此处无需前置写
	if req.Submit && !auto && (main.Status == models.CustomCodeStatusDraft ||
		main.Status == models.CustomCodeStatusRejected ||
		main.Status == models.CustomCodeStatusPending) {
		if err := mCustomCode.Submit(global.Db, main.ID, baseInfo.Versions, codeTypes); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	// 非首次提交（作者受信任）且通过安全扫描：自动审核通过，直接发布上线
	if req.Submit && auto {
		if err := PublishApprovedReview(global.Db, main.ID, reviewId, 0, "自动审核通过"); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	// 提交完成：发出领域事件，由后台 worker 异步处理后续副作用（人工路径通知审核员/审计等）
	if req.Submit {
		event.Publish(event.CustomCodeSubmitted, event.CustomCodeSubmittedPayload{
			CustomCodeID: main.ID,
			ReviewID:     reviewId,
			AuthorID:     userInfo.ID,
			Auto:         auto,
		})
	}

	apiReturn.SuccessData(c, gin.H{"id": main.ID, "reviewId": reviewId})
}

// Withdraw 作者撤回审核（待审 → 草稿）
func (a *CustomCodeApi) Withdraw(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mCustomCode := &models.CustomCode{}
	main, err := mCustomCode.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "内容不存在")
		return
	}
	if main.AuthorId != userInfo.ID {
		apiReturn.Error(c, "无权操作该内容")
		return
	}

	mReview := &models.CustomCodeReview{}
	pending, err := mReview.GetPendingByCustomCodeId(global.Db, main.ID)
	if err != nil || pending.ID == 0 {
		apiReturn.Error(c, "没有待审核的内容")
		return
	}

	if err := mReview.UpdateStatus(global.Db, pending.ID, models.CustomCodeReviewStatusDraft, 0, ""); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	if main.Status == models.CustomCodeStatusPending {
		if err := mCustomCode.UpdateStatus(global.Db, main.ID, models.CustomCodeStatusDraft); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	apiReturn.Success(c)
}

// Offline 作者下架
func (a *CustomCodeApi) Offline(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mCustomCode := &models.CustomCode{}
	main, err := mCustomCode.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "内容不存在")
		return
	}
	if main.AuthorId != userInfo.ID {
		apiReturn.Error(c, "无权操作该内容")
		return
	}

	if err := mCustomCode.SetOffline(global.Db, main.ID, ""); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// Delete 删除（仅未上线或已下架）
func (a *CustomCodeApi) Delete(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mCustomCode := &models.CustomCode{}
	main, err := mCustomCode.GetById(global.Db, req.Id)
	if err != nil {
		apiReturn.Error(c, "内容不存在")
		return
	}
	if main.AuthorId != userInfo.ID {
		apiReturn.Error(c, "无权操作该内容")
		return
	}
	if main.Status == models.CustomCodeStatusOnline {
		apiReturn.Error(c, "已上线的内容请先下架再删除")
		return
	}

	if err := mCustomCode.Delete(global.Db, main.ID); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	// 同时清理块与快照
	global.Db.Where("custom_code_id = ?", main.ID).Delete(&models.CustomCodeBlock{})
	global.Db.Where("custom_code_id = ?", main.ID).Delete(&models.CustomCodeReview{})

	apiReturn.Success(c)
}

// ForceOffline 管理员强制下架
func (a *CustomCodeApi) ForceOffline(c *gin.Context) {
	req := customCodeApiStructs.CustomCodeOfflineReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	if !models.HasAdmin(userInfo.Role) {
		apiReturn.Error(c, "需要管理员权限")
		return
	}
	if strings.TrimSpace(req.Reason) == "" {
		apiReturn.Error(c, "请填写下架原因")
		return
	}

	if err := (&models.CustomCode{}).SetOffline(global.Db, req.Id, req.Reason); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// UploadPreviewImage 块预览图上传（单张 ≤512K）
func (a *CustomCodeApi) UploadPreviewImage(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(fildDir)
	if !isExist {
		os.MkdirAll(fildDir, os.ModePerm)
	}

	uploadInfo, err := file.FormFileUpload(c, file.FormFileUploadOptions{
		FormFileName:  "imgfile",
		AgreeExtNames: []string{".png", ".jpg", ".gif", ".jpeg", ".webp"},
		MaxSize:       models.CustomCodeMaxImageSize,
		SaveDir:       fildDir,
	})
	if err != nil {
		switch err {
		case file.ErrUploadExceedMaxSize:
			apiReturn.Error(c, "预览图不能大于512K")
			return
		case file.ErrUploadExtensionNameNotAllowed:
			apiReturn.Error(c, "不支持的图片格式")
			return
		default:
			apiReturn.Error(c, "上传失败")
			return
		}
	}

	pureFilePath := uploadInfo.FileSavePath[len(configUpload):]
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, uploadInfo.FileOriginalName, uploadInfo.Ext, pureFilePath)

	apiReturn.SuccessData(c, gin.H{
		"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
	})
}
