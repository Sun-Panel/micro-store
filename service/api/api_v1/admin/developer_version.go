package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// DeveloperVersionApi 开发者版本管理 API
type DeveloperVersionApi struct{}

// GetList 获取版本列表
func (a *DeveloperVersionApi) GetList(c *gin.Context) {
	param := MicroAppVersionGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Limit < 1 {
		param.Limit = 10
	}

	m := models.MicroAppVersion{}
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, &param.AppId, param.Status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取版本详情
func (a *DeveloperVersionApi) GetInfo(c *gin.Context) {
	param := MicroAppVersionGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, version)
}

// Create 创建版本
func (a *DeveloperVersionApi) Create(c *gin.Context) {
	param := MicroAppVersionCreateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	version := &models.MicroAppVersion{
		AppId:       param.AppId,
		Version:     param.Version,
		VersionCode: param.VersionCode,
		PackageUrl:  param.PackageUrl,
		PackageHash: param.PackageHash,
		VersionDesc: param.VersionDesc,
		Config:      param.Config,
	}

	if err := biz.CreateVersionWithCheck(global.Db, version); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.SuccessData(c, version)
}

// Update 更新版本
func (a *DeveloperVersionApi) Update(c *gin.Context) {
	param := MicroAppVersionUpdateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.UpdateVersion(global.Db, param.Id, param.Version, param.VersionCode); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// SubmitReview 提交审核
func (a *DeveloperVersionApi) SubmitReview(c *gin.Context) {
	param := MicroAppVersionSubmitReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.SubmitReview(global.Db, param.VersionId); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// CancelReview 撤销审核
func (a *DeveloperVersionApi) CancelReview(c *gin.Context) {
	param := MicroAppVersionCancelReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.CancelReview(global.Db, param.VersionId); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// Delete 删除版本
func (a *DeveloperVersionApi) Delete(c *gin.Context) {
	param := MicroAppVersionDeleteReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.DeleteVersion(global.Db, param.Ids); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// Offline 下架版本
func (a *DeveloperVersionApi) Offline(c *gin.Context) {
	param := MicroAppVersionOfflineReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 平台下架时，原因必填
	if param.Type == 2 && param.Reason == "" {
		apiReturn.ErrorParamFomat(c, "平台下架时，下架原因不能为空")
		return
	}

	m := models.MicroAppVersion{}
	if err := m.Offline(global.Db, param.Id, param.Type, param.Reason); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// GetVersionList 获取版本列表（管理员专用）
func (a *DeveloperVersionApi) GetVersionList(c *gin.Context) {
	param := AdminVersionGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Limit < 1 {
		param.Limit = 10
	}

	m := models.MicroAppVersion{}
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, &param.AppId, param.Status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetPendingList 获取待审核版本列表（审核员专用）
func (a *DeveloperVersionApi) GetPendingList(c *gin.Context) {
	param := AdminVersionGetPendingListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if param.Page < 1 {
		param.Page = 1
	}
	if param.Limit < 1 {
		param.Limit = 10
	}

	m := models.MicroAppVersion{}
	// 状态 1 表示审核中
	status := 1
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, nil, &status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// Review 审核版本（审核员专用）
func (a *DeveloperVersionApi) Review(c *gin.Context) {
	param := AdminVersionReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 验证审核状态
	if param.Status != 1 && param.Status != 2 {
		apiReturn.ErrorParamFomat(c, "审核状态无效")
		return
	}

	// 获取当前审核员ID（从token中获取）
	reviewerId := c.GetUint("adminId")

	// 调用业务层审核方法
	if err := biz.ReviewVersionWithReviewer(global.Db, param.VersionId, param.Status, reviewerId, param.ReviewNote); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// handleBizError 统一处理业务错误
func handleBizError(c *gin.Context, err error) {
	// 业务错误：转换为数字错误码，前端统一处理
	if errCode, ok := biz.IsBizError(err); ok {
		intCode := bizCodeToInt(errCode) // API层负责转换
		apiReturn.ErrorByCode(c, intCode)
		return
	}
	// 其他错误：数据库错误
	apiReturn.ErrorDatabase(c, err.Error())
}

// bizCodeToInt 业务错误码转数字错误码（API层负责转换）
func bizCodeToInt(code string) int {
	codeMap := map[string]int{
		biz.ErrCodeAppNotFound:          2000,
		biz.ErrCodeVersionNotFound:      2001,
		biz.ErrCodeVersionExists:        2002,
		biz.ErrCodeVersionCodeExists:    2003,
		biz.ErrCodeStatusNotAllowed:     2004,
		biz.ErrCodeApprovedCannotDelete: 2005,
		biz.ErrCodeNotPendingReview:     2006,
		biz.ErrCodeNoUpdateContent:      2007,
		// 微应用开发者相关 3000-3099
		biz.ErrCodeAppIdExists:         3000,
		biz.ErrCodeNoPermission:        3001,
		biz.ErrCodePendingReviewExists: 3002,
		biz.ErrCodeNoPendingReviewApp:  3003,
	}

	if intCode, ok := codeMap[code]; ok {
		return intCode
	}
	return -1
}
