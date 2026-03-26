package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// MicroAppVersionAdminApi 审核员/管理员专用版本管理 API
type MicroAppVersionAdminApi struct{}

// GetVersionList 获取版本列表（管理员专用）
func (a *MicroAppVersionAdminApi) GetVersionList(c *gin.Context) {
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
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, &param.AppRecordId, param.Status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetPendingList 获取待审核版本列表（审核员专用）
func (a *MicroAppVersionAdminApi) GetPendingList(c *gin.Context) {
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
	// 状态 0 表示审核中
	status := 0
	db := global.Db.Preload("MicroApp.LangList").Preload("MicroApp.DefaultLangInfo")
	list, total, err := m.GetList(db, param.Page, param.Limit, nil, &status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// Review 审核版本（审核员专用）
func (a *MicroAppVersionAdminApi) Review(c *gin.Context) {
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
	if err := biz.MicroAppVersion.Review(global.Db, param.VersionId, param.Status, reviewerId, param.ReviewNote); err != nil {
		handleBizError(c, err)
		return
	}

	apiReturn.Success(c)
}

// Offline 下架版本（管理员/审核员专用）
func (a *MicroAppVersionAdminApi) Offline(c *gin.Context) {
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

// GetLatestOnlineByAppModelId 获取最新在线版本（审核员专用）
func (this *MicroAppVersionAdminApi) GetLatestOnlineByAppModelId(c *gin.Context) {
	req := models.MicroApp{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	// biz.MicroApp.GetById(global.Db, req.ID)
	version, err := biz.MicroAppVersion.GetLatestOnlineByAppModelId(global.Db, req.ID)
	if err != nil {
		base.HandleBizErrorAndReturn(c, err)
		return
	}
	apiReturn.SuccessData(c, version)
}
