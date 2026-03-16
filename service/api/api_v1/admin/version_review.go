package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// VersionReviewApi 版本审核 API
type VersionReviewApi struct{}

// GetList 获取版本列表
func (a *VersionReviewApi) GetList(c *gin.Context) {
	param := AdminVersionGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppVersion{}
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, &param.AppId, param.Status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 填充应用信息
	for i := range list {
		var app models.MicroApp
		if err := global.Db.Where("id = ?", list[i].AppId).First(&app).Error; err == nil {
			list[i].AppId = app.ID // 临时存储应用ID
		}
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetPendingList 获取待审核版本列表
func (a *VersionReviewApi) GetPendingList(c *gin.Context) {
	param := AdminVersionGetPendingListReq{}
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

	list, total, err := biz.GetPendingListWithAppInfo(global.Db, param.Page, param.Limit)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// Review 审核版本
func (a *VersionReviewApi) Review(c *gin.Context) {
	param := AdminVersionReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前管理员ID（从 token 中获取）
	// 这里需要根据实际的认证方式来获取管理员ID
	// 暂时使用 0，实际使用时需要从 context 中获取
	adminId := uint(0)
	if c.GetUint("userId") > 0 {
		adminId = c.GetUint("userId")
	}

	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, param.VersionId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 只有待审核状态可以审核
	if version.Status != 0 {
		apiReturn.Error(c, "当前版本不在待审核状态")
		return
	}

	if err := m.Review(global.Db, param.VersionId, param.Status, adminId, param.ReviewNote); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Offline 下架版本
func (a *VersionReviewApi) Offline(c *gin.Context) {
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
