package microapp

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MicroAppApi struct{}

// GetInfo 获取微应用详情（管理员专用）
func (a *MicroAppApi) GetInfo(c *gin.Context) {
	req := models.BaseModel{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	lang := base.GetCurrentUserLang(c)

	info, err := biz.MicroApp.GetByIdWithLang(global.Db, req.ID, lang, "Developer")
	if err != nil {
		base.HandleBizErrorAndReturn(c, err)
		return
	}

	if info.Status == 0 {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeAppNotFound)
		return
	}

	apiReturn.SuccessData(c, info)
}

// GetList 获取版本列表
func (a *MicroAppApi) GetList(c *gin.Context) {
	req := MicroAppVersionGetListReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 1
	}

	m := models.MicroAppVersion{}
	status := 1
	list, total, err := m.GetList(global.Db.Debug(), req.Page, req.Limit, &req.AppRecordId, &status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for k, _ := range list {
		list[k].PackageSrc = "" // 置空
		// list[k].PackageUrl = biz.MicroAppPackage.GenerateDownloadURL(v.PackageSrc)
	}

	apiReturn.SuccessListData(c, list, total)
}
