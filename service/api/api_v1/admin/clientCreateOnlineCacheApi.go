package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClientCreateOnlineCacheApi struct {
}

func (a *ClientCreateOnlineCacheApi) GetAll(c *gin.Context) {
	res := systemSetting.ClientCreateOnlineClientCache{}

	global.SystemSetting.GetValueByInterface(systemSetting.ClientCreateOnlineClientCacheKey, &res)

	apiReturn.SuccessData(c, res)
}

func (a *ClientCreateOnlineCacheApi) SetAll(c *gin.Context) {
	req := systemSetting.ClientCreateOnlineClientCache{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	err := global.SystemSetting.Set(systemSetting.ClientCreateOnlineClientCacheKey, &req)
	if err != nil {
		global.Logger.Errorln(err)
	}

	apiReturn.Success(c)
}
