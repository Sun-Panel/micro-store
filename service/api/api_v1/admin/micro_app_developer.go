package admin

import (
	"fmt"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// MicroAppDeveloperApi 开发者专用微应用 API
type MicroAppDeveloperApi struct{}

// GetMyList 获取开发者的微应用列表（开发者专用）
func (a *MicroAppDeveloperApi) GetMyList(c *gin.Context) {
	param := MicroAppGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	opts := models.MicroAppQueryOptions{
		Page:        param.Page,
		Limit:       param.Limit,
		Status:      param.Status,
		CategoryId:  param.CategoryId,
		DeveloperId: &developer.ID,
		KeyWord:     param.KeyWord,
	}

	list, total, err := biz.MicroAppDeveloper.GetDeveloperAppList(global.Db, opts)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取微应用详情（开发者专用）
func (a *MicroAppDeveloperApi) GetInfo(c *gin.Context) {
	req := models.BaseModel{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	lang := base.GetCurrentUserLang(c)
	developer := base.GetCurrentDeveloper(c)

	// 获取基础信息（有缓存优化）
	info, err := biz.MicroApp.GetByIdWithLang(global.Db.Debug(), req.ID, lang, "Developer", "LangList")
	if err != nil {
		handleBizError(c, err)
		return
	}

	// global.Logger.Debugln("current app info:", cmn.AnyToJsonStr(info))
	// global.Logger.Debugln("current developer ID:", developer.ID, "appinfo Developer ID:", info.DeveloperId)

	// 验证权限
	if info.Developer.ID != developer.ID {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoCurrentPermission)
		return
	}

	// 获取扩展信息（审核状态、草稿）
	reviewStatus, draft, err := biz.MicroAppDeveloper.GetDeveloperAppExtendInfo(global.Db, info.MicroAppId)
	if err != nil {
		handleBizError(c, err)
		return
	}

	// 组合数据返回
	result := map[string]interface{}{
		"id":           info.ID,
		"microAppId":   info.MicroAppId,
		"appName":      info.AppName,
		"appIcon":      info.AppIcon,
		"appDesc":      info.AppDesc,
		"remark":       info.Remark,
		"categoryId":   info.CategoryId,
		"chargeType":   info.ChargeType,
		"points":       info.Points,
		"authorId":     info.DeveloperId,
		"status":       info.Status,
		"screenshots":  info.Screenshots,
		"langList":     info.LangList,
		"createTime":   info.CreatedAt,
		"updateTime":   info.UpdatedAt,
		"reviewStatus": reviewStatus, // 审核状态：0-已通过 1-审核中 2-已拒绝 3-草稿
		"draft":        draft,        // 草稿版本（如果存在）
	}

	apiReturn.SuccessData(c, result)
}

// 根据微应用自增ID获取微应用信息和审核信息
func (a *MicroAppDeveloperApi) GetMicroInfoAndReviewInfoByMicroAppModelId(c *gin.Context) {
	req := MicroAppGetReviewInfoByModelIdReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	// 获取微应用的全部信息（包含开发者、语言列表）
	info, err := biz.MicroApp.GetById(global.Db, req.Id, "Developer", "LangList")
	if err != nil {
		handleBizError(c, err)
		return
	}

	// 验证权限
	if info.DeveloperId != developer.ID {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoCurrentPermission)
		return
	}

	review, err := biz.MicroApp.GetMicroInfoAndLatestReview(global.Db, req.Id)
	if err != nil {
		handleBizError(c, err)
		return
	}

	// 构建返回数据
	result := map[string]interface{}{
		"microApp":       info,   // 微应用全部数据
		"microAppReview": review, // 最新的审核记录
	}

	apiReturn.SuccessData(c, result)
}

// Create 创建微应用（开发者专用）
func (a *MicroAppDeveloperApi) Create(c *gin.Context) {
	param := MicroAppCreateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	developer := base.GetCurrentDeveloper(c)

	// 验证微应用ID格式
	if errMsg := a.validateMicroAppIdFormat(param.MicroAppId, developer.DeveloperName); errMsg != "" {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// // 转换 langMap
	// langMap := convertToBizLangMap(param.LangMap)

	opts := biz.DeveloperAppOptions{
		MicroAppId: param.MicroAppId,
		MicroAppBaseInfo: models.MicroAppBaseInfo{
			AdminName: param.AdminName,
		},
	}

	result, err := biz.MicroAppDeveloper.CreateAppAndReview(global.Db, opts)
	if err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.SuccessData(c, result)
}

// Update 更新微应用（开发者专用，不提交审核）
func (a *MicroAppDeveloperApi) Update(c *gin.Context) {
	param := MicroAppUpdateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	developer := base.GetCurrentDeveloper(c)

	opts := biz.DeveloperAppOptions{
		ID: param.ID,
		MicroAppBaseInfo: models.MicroAppBaseInfo{
			AppName:     param.AppName,
			AppIcon:     param.AppIcon,
			AppDesc:     param.AppDesc,
			Remark:      param.Remark,
			CategoryId:  param.CategoryId,
			ChargeType:  param.ChargeType,
			Points:      param.Points,
			Screenshots: param.Screenshots,
			AdminName:   param.AdminName,
		},
		LangMap:     param.LangMap,
		DeveloperId: developer.ID,
	}

	if err := biz.MicroAppDeveloper.UpdateDraftApp(global.Db, opts); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// UpdateLang 更新微应用语言（开发者专用，不提交审核）
func (a *MicroAppDeveloperApi) UpdateLang(c *gin.Context) {
	param := MicroAppUpdateLangReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	langMap := convertToBizLangMap(param.LangMap)

	if err := biz.MicroAppDeveloper.UpdateLang(global.Db, param.Id, developer.ID, langMap); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// SubmitReview 提交微应用审核（开发者专用）
func (a *MicroAppDeveloperApi) SubmitReview(c *gin.Context) {
	param := MicroAppSubmitReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	if err := biz.MicroAppDeveloper.SubmitAppReview(global.Db, param.ReviewId, developer.ID); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// CancelReview 撤销微应用主信息审核（开发者专用）
func (a *MicroAppDeveloperApi) CancelReview(c *gin.Context) {
	param := MicroAppCancelReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	if err := biz.MicroAppDeveloper.CancelAppReview(global.Db, param.ReviewId, developer.ID); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// GetReviewHistory 获取微应用审核历史（开发者专用）
func (a *MicroAppDeveloperApi) GetReviewHistory(c *gin.Context) {
	param := MicroAppGetReviewHistoryReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	list, total, err := biz.MicroAppDeveloper.GetAppReviewHistory(global.Db, param.AppId, developer.ID, param.Page, param.Limit)
	if err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// convertToBizLangMap 将 API 层的 MicroAppLangInfo 转换为 biz 层的 map[string]interface{}
func convertToBizLangMap(m map[string]MicroAppLangInfo) map[string]interface{} {
	if m == nil {
		return nil
	}
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = map[string]interface{}{"appName": v.AppName, "appDesc": v.AppDesc}
	}
	return result
}

// validateMicroAppIdFormat 验证微应用ID格式
// 规则：
// 1. 必须以开发者标识开头
// 2. 不能仅包含开发者标识（即不能是 "developerName-"）
// 3. 不能以 "-" 结尾
// 返回空字符串表示验证通过，否则返回错误信息
func (a *MicroAppDeveloperApi) validateMicroAppIdFormat(microAppId, developerName string) string {
	// 规则1: 必须以开发者标识开头
	prefix := developerName + "-"
	if !a.startsWith(microAppId, prefix) {
		return fmt.Sprintf("微应用ID必须以\"%s\"开头", developerName)
	}

	// 规则2: 不能仅包含开发者标识
	if microAppId == prefix {
		return "微应用ID不能仅包含开发者标识，请添加应用名称"
	}

	// 规则3: 不能以 "-" 结尾
	if a.endsWith(microAppId, "-") {
		return "微应用ID不能以\"-\"结尾"
	}

	return ""
}

// startsWith 检查字符串是否以指定前缀开头
func (a *MicroAppDeveloperApi) startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// endsWith 检查字符串是否以指定后缀结尾
func (a *MicroAppDeveloperApi) endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
