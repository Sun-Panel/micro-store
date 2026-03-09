package sunStore

import (
	"encoding/json"
	"errors"
	"net/http"
	"sun-panel/global"
	"sun-panel/lib/sunStore/request"
)

type SunPayApi struct {
	AppId     string
	ApiHost   string
	AppSecret string
}

type ApiCommonResp struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type AccessTokenResquest struct {
	Code string `json:"code"`
	// AppID     string `json:"appid"`
	// AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	// OpenId      string `json:"openid"`
	// UserId      uint   `json:"userid"`
	Scope string `json:"scope"`
}

type PasswordAuthRequest struct {
	GrantType string `json:"grant_type" form:"grant_type" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	// ClientID     string `json:"client_id" form:"client_id" binding:"required"`
	// ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
	State string `json:"state" form:"state"`
	Scop  string `json:"scope" form:"scope"`
}

type ClientCredentialsParam struct {
	State        string `json:"state" form:"state"`
	Scope        string `json:"scope" form:"scope"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type ClientCredentialsRespone struct {
	State        string `json:"state" form:"state"`
	Scope        string `json:"scope" form:"scope"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type ClientCredentialsAuthResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	TokenType    string `json:"token_type" form:"token_type"`
	ExpiresIn    int64  `json:"expires_in" form:"expires_in"`
}

var (
	ErrorServer = errors.New("server error")
)

func NewSunStoreApi(apiHost, appid, appSecret string) *SunPayApi {
	sunpayApi := SunPayApi{
		AppId:     appid,
		ApiHost:   apiHost,
		AppSecret: appSecret,
	}
	return &sunpayApi
}

func (s *SunPayApi) GetTokenByCode(code string) {

}

func (s *SunPayApi) GetAccessToken(param AccessTokenResquest) (AccessTokenResponse, error) {

	requestData := map[string]interface{}{
		"code":          param.Code,
		"client_id":     s.AppId,
		"client_secret": s.AppSecret,
		"grant_type":    "authorization_code", // 固定
	}

	respData := AccessTokenResponse{}

	host := s.ApiHost + "/oauth2/v1/token"

	global.Logger.Debugln("HOST", host)

	errCode, err, _ := request.SendPostJsonRequest(host, requestData, &respData)
	if err != nil {
		return respData, err
	}

	resByte, _ := json.Marshal(respData)
	// global.Logger.Debugln("RESP", respData)
	global.Logger.Debugln("RESP-json", string(resByte))
	_ = errCode

	return respData, nil
}

func (s *SunPayApi) PasswordAuth(param PasswordAuthRequest) (AccessTokenResponse, error) {

	requestData := map[string]interface{}{
		"client_id":     s.AppId,
		"client_secret": s.AppSecret,
		"grant_type":    "password", // 固定
		"username":      param.Username,
		"password":      param.Password,
	}

	respData := AccessTokenResponse{}

	host := s.ApiHost + "/oauth2/v1/token"

	global.Logger.Debugln("HOST", host)

	errCode, err, httpResp := request.SendPostJsonRequest(host, requestData, &respData)

	// 服务器错误
	if httpResp != nil && httpResp.StatusCode == http.StatusInternalServerError {
		return respData, ErrorServer
	}

	if err != nil {
		return respData, err
	}

	resByte, _ := json.Marshal(respData)
	// global.Logger.Debugln("RESP", respData)
	global.Logger.Debugln("RESP-json", string(resByte))
	_ = errCode

	return respData, nil
}

func (s *SunPayApi) ClientCredentialsAuth(param ClientCredentialsParam, isRefreshToken bool) (ClientCredentialsAuthResp, error) {

	grantType := "client_credentials"
	if isRefreshToken {
		grantType = "refresh_token" // 刷新令牌
	}
	requestData := map[string]interface{}{
		"client_id":     s.AppId,
		"client_secret": s.AppSecret,
		"grant_type":    grantType,
		"state":         param.State,
		"scope":         param.Scope,
	}

	respData := ClientCredentialsAuthResp{}

	host := s.ApiHost + "/oauth2/v1/clientCredentials/token"

	global.Logger.Debugln("HOST", host)

	errCode, err, _ := request.SendPostJsonRequest(host, requestData, &respData)
	if err != nil {
		return respData, err
	}

	resByte, _ := json.Marshal(respData)
	// global.Logger.Debugln("RESP", respData)
	global.Logger.Debugln("RESP-json", string(resByte))
	_ = errCode

	return respData, nil
}
