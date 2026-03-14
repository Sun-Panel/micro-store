package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
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

	// 设置默认值
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

	// 检查应用是否存在
	app := models.MicroApp{}
	if _, err := app.GetById(global.Db, param.AppId); err != nil {
		apiReturn.Error(c, "微应用不存在")
		return
	}

	// 检查版本号是否存在
	m := models.MicroAppVersion{}
	exists, err := m.CheckVersionExist(global.Db, param.AppId, param.Version, 0)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if exists {
		apiReturn.Error(c, "版本号已存在")
		return
	}

	// 检查版本号数字是否存在
	exists, err = m.CheckVersionCodeExist(global.Db, param.AppId, param.VersionCode, 0)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if exists {
		apiReturn.Error(c, "版本号数字已存在")
		return
	}

	version := models.MicroAppVersion{
		AppId:       param.AppId,
		Version:     param.Version,
		VersionCode: param.VersionCode,
		PackageUrl:  param.PackageUrl,
		PackageHash: param.PackageHash,
		VersionDesc: param.VersionDesc,
		Config:      param.Config,
		Status:      -1, // 默认草稿状态
	}

	if err := version.Create(global.Db); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
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

	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 只有待审核状态可以更新
	if version.Status != 0 {
		apiReturn.Error(c, "只有待审核的版本可以更新")
		return
	}

	updateData := map[string]interface{}{}
	if param.Version != "" {
		updateData["version"] = param.Version
	}
	if param.VersionCode > 0 {
		updateData["version_code"] = param.VersionCode
	}

	if err := m.Update(global.Db, param.Id, updateData); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
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

	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, param.VersionId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 只有草稿或已拒绝状态可以提交审核
	if version.Status != -1 && version.Status != 2 {
		apiReturn.Error(c, "当前状态不允许提交审核")
		return
	}

	if err := m.Update(global.Db, param.VersionId, map[string]interface{}{"status": 0}); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
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

	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, param.VersionId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 只有待审核状态可以撤销
	if version.Status != 0 {
		apiReturn.Error(c, "只有待审核的版本可以撤销")
		return
	}

	// 撤销审核后变为草稿状态
	if err := m.Update(global.Db, param.VersionId, map[string]interface{}{"status": -1}); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
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

	m := models.MicroAppVersion{}
	// 检查版本是否存在
	version, err := m.GetById(global.Db, param.Ids[0])
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 只有待审核或已拒绝状态可以删除
	if version.Status == 1 {
		apiReturn.Error(c, "已通过的版本不能删除")
		return
	}

	if err := m.Delete(global.Db, param.Ids); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
