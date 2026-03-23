package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"

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

	list, total, err := biz.MicroAppDeveloper.GetDeveloperAppList(global.Db, param.Page, param.Limit, param.Status, param.CategoryId, developer.ID, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetMyInfo 获取开发者微应用详情（开发者专用）
func (a *MicroAppDeveloperApi) GetMyInfo(c *gin.Context) {
	param := MicroAppGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer := base.GetCurrentDeveloper(c)

	result, err := biz.MicroAppDeveloper.GetDeveloperAppInfo(global.Db, param.Id, developer.ID)
	if err != nil {
		handleBizError(c, err)
		return
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

	// 转换 langMap
	langMap := convertToBizLangMap(param.LangMap)

	result, err := biz.MicroAppDeveloper.CreateApp(global.Db, param.MicroAppId, param.AppName, param.AppIcon, param.AppDesc, param.Remark, param.CategoryId, param.ChargeType, param.Points, developer.ID, param.Screenshots, langMap)
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

	langMap := convertToBizLangMap(param.LangMap)

	if err := biz.MicroAppDeveloper.UpdateApp(global.Db, param.Id, developer.ID, param.AppName, param.AppIcon, param.AppDesc, param.Remark, param.CategoryId, param.ChargeType, param.Price, param.Screenshots, langMap); err != nil {
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

	if err := biz.MicroAppDeveloper.SubmitAppReview(global.Db, param.Id, developer.ID); err != nil {
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

	if err := biz.MicroAppDeveloper.CancelAppReview(global.Db, param.Id, developer.ID); err != nil {
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

// convertToBizLangMap 将 API 层的 MicroAppLangInfo 转换为 biz 层的 map[string]map[string]string
func convertToBizLangMap(m map[string]MicroAppLangInfo) map[string]map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]map[string]string, len(m))
	for k, v := range m {
		result[k] = map[string]string{"appName": v.AppName, "appDesc": v.AppDesc}
	}
	return result
}
