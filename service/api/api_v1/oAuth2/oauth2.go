package oAuth2

import (
	"fmt"
	"net/http"
	"net/url"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
)

type OAuth2 struct {
}

// 重定向到授权服务器地址
func (o *OAuth2) Login(c *gin.Context) {

	// host := c.Request.Host

	// 判断当前使用的协议
	// scheme := "http"
	// if c.Request.TLS != nil {
	// 	scheme = "https"
	// }
	// domain := fmt.Sprintf("%s://%s", scheme, host)
	authPlatFormUrlWithPath := biz.SunStore.AuthEndpointUrl()

	// 获取所有的 GET 参数
	// queryParams := c.Request.URL.Query()
	callbackUrl := url.QueryEscape(c.Query("callback"))
	isBindStr := c.Query("isBind")

	sysApplication := systemSetting.ApplicationSetting{}
	global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &sysApplication)

	webSiteUrl := sysApplication.WebSiteUrl

	// 前端登录授权或绑定授权页面地址
	frontUrl := url.QueryEscape(webSiteUrl + "/oAuth2/login?callback=" + callbackUrl)
	if isBindStr == "true" {
		frontUrl = url.QueryEscape(webSiteUrl + "/oAuth2/login?callback=" + callbackUrl + "&isBind=true")
	}

	// 授权平台授权端点地址
	clientId, _ := biz.SunStore.GetClientIdAndSecret()
	rawQuery := fmt.Sprintf("client_id=%s&redirect_uri=%s&response_type=%s", clientId, frontUrl, "code")
	authPlatFormUrl := fmt.Sprintf("%s?%s", authPlatFormUrlWithPath, rawQuery)
	global.Logger.Debugln("auth platform Url", authPlatFormUrl)

	c.Redirect(http.StatusFound, authPlatFormUrl) // HTTP 302

}
