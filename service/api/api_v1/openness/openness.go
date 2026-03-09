package openness

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
)

type Openness struct {
}

func (a *Openness) LoginConfig(c *gin.Context) {
	cfg := systemSetting.ApplicationSetting{}
	if err := global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &cfg); err != nil {
		apiReturn.Error(c, "配置查询失败："+err.Error())
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"loginCaptcha": cfg.LoginCaptcha,
		"register":     cfg.Register,
	})
}

func (a *Openness) GetDisclaimer(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.DISCLAIMER); err != nil {
		global.SystemSetting.Set(systemSetting.DISCLAIMER, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}

func (a *Openness) GetAboutDescription(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.WEB_ABOUT_DESCRIPTION); err != nil {
		global.SystemSetting.Set(systemSetting.WEB_ABOUT_DESCRIPTION, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}

func (a *Openness) GetHomeBase(c *gin.Context) {
	result, _ := global.SystemSetting.GetMultipleVariableString("logo_text", "logo_url", "logo_click_to_link", "home_url")
	apiReturn.SuccessData(c, result)
}

func (a *Openness) GetProDescription(c *gin.Context) {
	result, _ := global.SystemSetting.GetMultipleVariableString("pro_before_buy_description", "pro_func_description")
	apiReturn.SuccessData(c, result)
}

func (a *Openness) GetRootPageDescription(c *gin.Context) {
	result, _ := global.SystemSetting.GetMultipleVariableString("root_page_description")
	apiReturn.SuccessData(c, result)
}
