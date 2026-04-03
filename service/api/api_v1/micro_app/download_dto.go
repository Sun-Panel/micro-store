package microapp

type DownloadGetUrlReq struct {
	MicroAppId string `json:"microAppId" binding:"required"`
	Version    string `json:"version"`
}
