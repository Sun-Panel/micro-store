package admin

import (
	"sun-panel/api/api_v1/common/apiData/adminApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AboutSettingApi struct {
}

func (a *AboutSettingApi) Save(c *gin.Context) {
	req := adminApiStructs.AboutSettingRequest{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := global.SystemSetting.Set(systemSetting.WEB_ABOUT_DESCRIPTION, req.Content); err != nil {
		apiReturn.Error(c, "保存失败")
		return
	}

	apiReturn.Success(c)
}
