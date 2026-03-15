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
	}

	if intCode, ok := codeMap[code]; ok {
		return intCode
	}
	return -1
}
