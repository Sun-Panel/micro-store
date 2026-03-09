package biz

import (
	"errors"
	"fmt"
	"sun-panel/global"
	"sun-panel/lib/sunStore"
	"sun-panel/lib/sunStore/openApi"
	"time"
)

type SunStoreType struct {
}

type SunStoreClientAPIAuthToken struct {
	AccessToken                  string
	RefreshToken                 string
	AccessTokenExpiresTimestamp  int64 // AccessToken 过期时间戳
	RefreshTokenExpiresTimestamp int64 // RefreshToken 过期时间戳
}

func (s *SunStoreType) GetClientIdAndSecret() (clientId string, clientSecret string) {

	clientId = global.Config.GetValueString("sun_store", "client_id")
	clientSecret = global.Config.GetValueString("sun_store", "client_secret")
	return
}

func (s *SunStoreType) ApiHost() string {
	return global.Config.GetValueString("sun_store", "api_host")
}

// 授权端点地址
func (s *SunStoreType) AuthEndpointUrl() string {
	return global.Config.GetValueString("sun_store", "auth_endpoint")
}

func (s *SunStoreType) GetClientApiToken(apiHost, clientId, clientSecret string) (string, error) {

	tokenData := SunStoreClientAPIAuthToken{}
	err := global.SystemSetting.GetValueByInterface("SunStoreClientAPIAuthToken", &tokenData)
	if err != nil || tokenData.AccessTokenExpiresTimestamp < time.Now().Unix() {
		global.Logger.Debug("重新登录到SunStore平台")
		// 重新登录
		tokenData, err = s.ClientApiAuthLogin(apiHost, clientId, clientSecret)
		if err != nil {
			return "", errors.New(fmt.Sprintf("clientApiAuthLogin:%s", err.Error()))
		}
		err := global.SystemSetting.Set("SunStoreClientAPIAuthToken", tokenData)
		if err != nil {
			global.Logger.Errorln("cache SunStoreClientAPIAuthToken fail", err)
			return "", err
		}
	}

	return tokenData.AccessToken, nil

}

// 此处暂时未实现refreshToken相关内容
func (s *SunStoreType) ClientApiAuthLogin(apiHost, clientId, clientSecret string) (SunStoreClientAPIAuthToken, error) {
	c := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)

	tokenData := SunStoreClientAPIAuthToken{}
	resp, err := c.ClientCredentialsAuth(sunStore.ClientCredentialsParam{}, false)
	if err != nil {
		return tokenData, err
	}

	tokenData.AccessTokenExpiresTimestamp = time.Now().Unix() + resp.ExpiresIn
	tokenData.AccessToken = resp.AccessToken

	return tokenData, nil
}

// 获取主平台的用户信息根据 accessToken
func (s *SunStoreType) GetMainPlatformUserInfo(apiHost, accessToken string) (openApi.UserInfoResp, error) {
	userOpenApi := openApi.NewUser(openApi.NewOpenApi(apiHost, accessToken))
	openUser, _, err := userOpenApi.GetCurrentUserInfo()
	return openUser, err
}

// 获取主平台的用户信息根据，如果用户信息不存在本地将创建用户
// func (s *SunStoreType) GetUserInfoByMainPlatformUserInfoOrCreateUser(apiHost, accessToken string) (models.User, error) {
// 	userInfo := models.User{}
// 	openUser, err := s.GetMainPlatformUserInfo(apiHost, accessToken)
// 	if err != nil {
// 		return userInfo, err
// 	}

// 	// 查询主平台邮箱是否在本平台有同邮箱账号
// 	// 存在进行绑定，否则创建新账号
// 	if err := global.Db.First(&userInfo, "username=?", openUser.Mail).Error; err != nil {
// 		// 创建新账号
// 		userInfo.Username = openUser.Username
// 		userInfo.Mail = openUser.Mail
// 		userInfo.Name = openUser.Name
// 		userInfo.Role = 2
// 		userInfo.Status = 1
// 		if errCreate := global.Db.Create(&userInfo).Error; errCreate != nil {
// 			return userInfo, errCreate
// 		}
// 	}
// 	return userInfo, err
// }
