package thirdAppApiStructs

type Oauth2AccountLoginReq struct {
	AppID     string `json:"appid"`
	AppSecret string `json:"app_secret"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Oauth2AuthLoginReq struct {
	AppID        string `json:"appid"`
	ResponseType string `json:"response_type"`
	RedirectUri  string `json:"redirect_uri"`
}

type Oauth2GetApiTokenReq struct {
	Code      string `json:"code"`
	AppID     string `json:"appid"`
	AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

type Oauth2GetApiTokenResp struct {
	ApiToken string `json:"api_token"`
	OpenId   string `json:"openid"`
	UserId   uint   `json:"userid"`
}
