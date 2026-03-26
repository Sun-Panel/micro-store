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
