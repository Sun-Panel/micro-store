package structs

// 授权码缓存数据
type OAuthCodeCacheData struct {
	SsoSession string
	ApiToken   string
}
