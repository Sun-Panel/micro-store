package thirdAppApiStructs

type GetThirdAppInfoReq struct {
	Appid string `json:"appid"`
}

type GetThirdAppInfoResp struct {
	Appid       string `json:"appid"`
	AppName     string `json:"appName"`
	IsAutoAuth  bool   `json:"isAutoAuth"`
	IsSsoLogout bool   `json:"isSsoLogout"`
	IsEnabled   bool   `json:"isEnabled"`
}
