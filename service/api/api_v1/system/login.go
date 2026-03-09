package system

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/sunStore"
	"sun-panel/lib/sunStore/openApi"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginApi struct {
}

// 登录输入验证
type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

type OAuth2CodeLoginReq struct {
	Code string `json:"code"`
}

type OAuth2CodeLoginResq struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	HeadImage string `json:"headImage"`
	Role      int    `json:"role"`
	Mail      string `json:"mail"`
}

var (
	ErrOAuth2CodeNeedLogout = errors.New("need to withdraw from the main platform to log in")
	ErrOAuth2CodeRetry      = errors.New("need to Retry")
)

// @Summary 登录账号
// @Accept application/json
// @Produce application/json
// @Param LoginLoginVerify body LoginLoginVerify true "登陆验证信息"
// @Tags user
// @Router /login [post]
func (l LoginApi) Login(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	settings := systemSetting.ApplicationSetting{}
	global.SystemSetting.GetValueByInterface("system_application", &settings)

	// 验证验证码
	if settings.Login.LoginCaptcha {
		if err := biz.Captcha.CaptchaVerifyHandle(c, param.VCode, true); err != nil {
			apiReturn.Error(c, err.Error())
			return
		}
	}

	mUser := models.User{}
	var (
		err  error
		info models.User
	)
	bToken := ""
	param.Username = strings.TrimSpace(param.Username)
	if info, err = mUser.GetUserInfoByUsernameAndPassword(param.Username, cmn.PasswordEncryption(param.Password)); err != nil {
		// 未找到记录 账号或密码错误
		if err == gorm.ErrRecordNotFound {
			apiReturn.ErrorByCode(c, 1003)
			return
		} else {
			// 未知错误
			apiReturn.Error(c, err.Error())
			return
		}

	}

	// 停用或未激活
	if info.Status != 1 {
		apiReturn.ErrorByCode(c, 1004)
		return
	}

	bToken = info.Token
	if info.Token == "" {
		// 生成token
		buildTokenOver := false
		for !buildTokenOver {
			bToken = cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
			if _, err := mUser.GetUserInfoByToken(bToken); err != nil {
				// 保存token
				mUser.UpdateUserInfoByUserId(info.ID, map[string]interface{}{
					"token": bToken,
				})
				buildTokenOver = true
			}
		}
		info.Token = bToken
	}
	info.Password = ""
	info.ReferralCode = ""

	// global.UserToken.SetDefault(bToken, info)
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(info.ID))))
	global.CUserToken.SetDefault(cToken, bToken)
	global.Logger.Debug("token:", cToken, "|", bToken)
	global.Logger.Debug(global.CUserToken.Get(cToken))

	// 设置当前用户信息
	c.Set("userInfo", info)
	info.Token = cToken // 重要 采用cToken,隐藏真实token
	apiReturn.SuccessData(c, info)
}

// 授权登录端点
// 警告：此接口为半公开模式，获取登录信息请一定判断是否为登录状态
func (l *LoginApi) OAuth2CodeLogin(c *gin.Context) {
	req := OAuth2CodeLoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取基本的配置信息
	clientId, clientSecret := biz.SunStore.GetClientIdAndSecret()
	apiHost := biz.SunStore.ApiHost()

	// 获取 access_token 和 openid
	sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
	accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
		Code: req.Code,
	})
	if err != nil {
		// 获取access_token失败
		global.Logger.Errorln("获取access_token失败", err.Error())
		apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
		return
	}

	// 根据 access_token 获取主平台账号信息（包括邮箱等信息）
	// 如果已登录取出登录的用户信息
	userInfo, loginStatus := base.GetCurrentUserInfo(c)
	global.Logger.Debugln("loginStatus", loginStatus)
	if loginStatus {
		global.Logger.Debugln("userInfo", userInfo)

		userInfo, err = l.oAuth2CodeAlreadyLoginAuthProcess(apiHost, userInfo, accessToken)
		if err != nil {
			if errors.Is(err, ErrOAuth2CodeNeedLogout) {
				apiReturn.ErrorByCodeAndMsg(c, -3, "need logout")
				return
			}

			global.Logger.Errorln("cannot be authorized:", err)
			apiReturn.Error(c, err.Error())
		}

	} else {
		// 未登录 - 流程
		userInfo, err = l.oAuth2CodeNoLoggedAuthProcess(apiHost, accessToken)
		if err != nil {
			if errors.Is(ErrOAuth2CodeRetry, err) {
				global.Logger.Errorln("oAuth2CodeNoLoggedAuthProcess", err)
				apiReturn.Error(c, ErrOAuth2CodeRetry.Error())
				return
			}

			// 其他错误
			global.Logger.Errorln(apiReturn.ErrorNumAndReturnMsg(c, err))
			return
		}
	}

	// 登录流程
	// 生成数据库长期bToken和ctoken短期cToken
	mUser := models.User{}
	bToken := userInfo.Token
	if userInfo.Token == "" {
		// 生成token
		buildTokenOver := false
		for !buildTokenOver {
			bToken = cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
			if _, err := mUser.GetUserInfoByToken(bToken); err != nil {
				// 保存token
				mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
					"token": bToken,
					// "store_email":
				})
				buildTokenOver = true
			}
		}
		userInfo.Token = bToken
	}

	// 生成Token
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(userInfo.ID))))
	global.CUserToken.SetDefault(cToken, bToken)
	global.CUserApiTokenAccessToken.SetDefault(cToken, accessToken.AccessToken)
	global.CUserAccessTokenApiToken.SetDefault(accessToken.AccessToken, cToken)

	apiReturn.SuccessData(c, OAuth2CodeLoginResq{
		Token:    cToken,
		Username: userInfo.Username,
		Name:     userInfo.Name,
		Role:     userInfo.Role,
	})
}

func (l *LoginApi) GetMainPlatformUserInfo(apiHost, accessToken string) (openApi.UserInfoResp, error) {
	userOpenApi := openApi.NewUser(openApi.NewOpenApi(apiHost, accessToken))
	openUser, _, err := userOpenApi.GetCurrentUserInfo()
	return openUser, err
}

// 授权绑定接口
func (l *LoginApi) OAuth2CodeBind(c *gin.Context) {
	req := OAuth2CodeLoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取基本的配置信息
	clientId, clientSecret := biz.SunStore.GetClientIdAndSecret()
	apiHost := biz.SunStore.ApiHost()

	// 获取 access_token 和 openid
	sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
	accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
		Code: req.Code,
	})
	if err != nil {
		// 获取access_token失败
		global.Logger.Errorln("获取access_token失败", err.Error())
		apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
		return
	}

	// 根据 access_token 获取主平台账号信息（包括邮箱等信息）
	// 如果已登录取出登录的用户信息
	userInfo, loginStatus := base.GetCurrentUserInfo(c)
	global.Logger.Debugln("loginStatus", loginStatus)
	if loginStatus {
		global.Logger.Debugln("userInfo", userInfo)

		userInfo, err = l.oAuth2CodeAlreadyLoginAuthProcess(apiHost, userInfo, accessToken)
		if err != nil {
			if errors.Is(err, ErrOAuth2CodeNeedLogout) {
				apiReturn.ErrorByCodeAndMsg(c, -3, "need logout")
				return
			}

			global.Logger.Errorln("无法通过授权:", err)
			apiReturn.Error(c, err.Error())
		}

	} else {
		// 未登录 - 流程
		userInfo, err = l.oAuth2CodeNoLoggedAuthProcess(apiHost, accessToken)
		if err != nil {
			if errors.Is(ErrOAuth2CodeRetry, err) {
				apiReturn.Error(c, ErrOAuth2CodeRetry.Error())
				return
			}

			// 其他错误
			global.Logger.Errorln(apiReturn.ErrorNumAndReturnMsg(c, err))
			return
		}
	}
	_ = userInfo

	apiReturn.Success(c)

}

// 解除绑定(目前单向解除绑定)
func (l *LoginApi) OAuth2UnBind(c *gin.Context) {
	req := OAuth2CodeLoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mUser := models.User{}
	mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
		"store_email": "",
	})

	apiReturn.Success(c)

}

// 单点退出
func (l *LoginApi) SSOLogoutWebhook(c *gin.Context) {
	type Req struct {
		AccessToken string `json:"accessToken"`
	}

	req := Req{}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	apitoken, ok := global.CUserAccessTokenApiToken.Get(req.AccessToken)
	if !ok {
		c.String(http.StatusOK, "")
		return
	}

	global.CUserAccessTokenApiToken.Delete(req.AccessToken)
	global.CUserApiTokenAccessToken.Delete(apitoken)
	global.CUserToken.Delete(apitoken)
}

// 已经登录的授权流程
func (l *LoginApi) oAuth2CodeAlreadyLoginAuthProcess(apiHost string, userInfo models.User, accessToken sunStore.AccessTokenResponse) (models.User, error) {
	// 已经登录 - 流程

	// 未绑定已有账号,获取主平台的账号、邮箱等信息
	openUser, err := l.GetMainPlatformUserInfo(apiHost, accessToken.AccessToken)
	if err != nil {
		return userInfo, err
	}

	// 判断当前账号的邮箱是否与主平台账号邮箱一致
	// 一致就返回用户，否则通知前端退出登录状态
	if openUser.Mail != userInfo.Username {
		// 不允许绑定账号
		global.Logger.Errorln("The mail account is inconsistent, the notification platform withdraws from the login status:", err)
		return userInfo, ErrOAuth2CodeNeedLogout
	}

	global.Logger.Debugln("bind User", openUser.Mail)

	// 绑定商城平台账号
	mUser := models.User{}
	mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
		"store_email": openUser.Mail,
		"lang":        openUser.Lang,
		"system_lang": openUser.SystemLang,
		"time_zone":   openUser.TimeZone,
	})
	userInfo.StoreEmail = openUser.Mail

	return userInfo, nil
}

// 未登录的授权流程
func (l *LoginApi) oAuth2CodeNoLoggedAuthProcess(apiHost string, accessToken sunStore.AccessTokenResponse) (models.User, error) {
	userInfo := models.User{}
	openUser, err := biz.SunStore.GetMainPlatformUserInfo(apiHost, accessToken.AccessToken)
	if err != nil {
		global.Logger.Errorln("GetMainPlatformUserInfo", err)
		// 可能过期了
		return userInfo, ErrOAuth2CodeRetry
	}

	// 查询商城平台邮箱是否在本平台有同邮箱账号
	// 存在进行绑定，否则创建新账号
	if err := global.Db.First(&userInfo, "username=?", openUser.Mail).Error; err != nil {
		// 创建新账号
		userInfo.Username = openUser.Username
		userInfo.Mail = openUser.Mail
		userInfo.Name = openUser.Name
		userInfo.Role = 2
		userInfo.Status = 1
		userInfo.StoreEmail = openUser.Mail
		userInfo.Lang = openUser.Lang
		userInfo.SystemLang = openUser.SystemLang
		userInfo.TimeZone = openUser.TimeZone
		userInfo.Password = openUser.Password

		if errCreate := global.Db.Create(&userInfo).Error; errCreate != nil {
			return userInfo, errCreate
		}
	}

	// 绑定商城平台账号
	mUser := models.User{}
	mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
		"head_image":  openUser.HeadImage,
		"name":        openUser.Name,
		"store_email": openUser.Mail,
		"lang":        openUser.Lang,
		"system_lang": openUser.SystemLang,
		"time_zone":   openUser.TimeZone,
		"password":    openUser.Password,
	})
	return userInfo, err
}

// 安全退出
func (l *LoginApi) Logout(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	cToken := c.GetHeader("token")
	global.CUserToken.Delete(cToken)
	apiReturn.Success(c)
}
