package system

import (
	"strings"
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type HtmlPageApi struct {
}

var HtmlPageApiVarPrefix = "htmlPage."

func (a *HtmlPageApi) Get(c *gin.Context) {

	userInfo, _ := base.GetCurrentUserInfo(c)
	req := systemApiStructs.HtmlPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	pageInfo := systemApiStructs.HtmlPageInfo{}
	req.PageName = strings.TrimSpace(req.PageName)
	if err := global.SystemSetting.GetValueByInterface(HtmlPageApiVarPrefix+req.PageName, &pageInfo); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	if pageInfo.IsLogin && userInfo.ID == 0 {
		apiReturn.ErrorNotLogin(c)
		return
	}

	apiReturn.SuccessData(c, pageInfo)
}
