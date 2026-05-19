package biz

import (
	"time"

	"sun-panel/global"
	"sun-panel/lib/cache"
	sunpanelauth "sun-panel/lib/sunPanelAuth"
	"sun-panel/models"
)

type SunPanelAuthType struct {
	client             *sunpanelauth.Client
	versionSecretCache cache.Cacher[string]
}

func (s *SunPanelAuthType) Init() {
	s.versionSecretCache = global.NewCache[string](168*time.Hour, 7*time.Hour, "VersionSecretCache")
}

// getClient 获取或初始化 sunPanelAuth 客户端
func (s *SunPanelAuthType) getClient() *sunpanelauth.Client {
	if s.client != nil {
		return s.client
	}

	apiHost := global.Config.GetValueString("sun_panel_auth", "api_host")
	apiSecret := global.Config.GetValueString("sun_panel_auth", "api_secret")
	global.Logger.Debugln("sunPanelAuth api_host: %s, api_secret: %s", apiHost, apiSecret)
	s.client = sunpanelauth.NewClient(apiHost, apiSecret)
	return s.client
}

// CaptchaLogin 验证码登录
func (s *SunPanelAuthType) CaptchaLogin(captcha string) (user models.User, err error) {
	resp, err := s.getClient().CaptchaLogin(captcha)
	if err != nil {
		return
	}

	user, err = user.GetUserInfoByUsername(resp.UserInfo.Username)
	if err != nil {
		return
	}

	return
}

// VersionSecretMap 查询所有版本密钥映射，并将每个版本逐一缓存
func (s *SunPanelAuthType) VersionSecretMap() (sunpanelauth.VersionSecretMapResponse, error) {
	result, err := s.getClient().VersionSecretMap()
	if err != nil {
		return nil, err
	}

	// 逐一保存到缓存：key = "vs:{version}", value = secretKey
	for version, secretKey := range result {
		s.versionSecretCache.SetDefault(version, secretKey)
	}

	return result, nil
}

// GetVersionSecret 根据版本号获取密钥（优先从缓存读取）
func (s *SunPanelAuthType) GetVersionSecret(version string) (string, error) {
	// 1. 尝试从缓存获取
	cacheKey := version
	if secretKey, ok := s.versionSecretCache.Get(cacheKey); ok {
		return secretKey, nil
	}

	// 2. 缓存未命中，拉取全量数据（内部会逐条写入缓存）
	_, err := s.VersionSecretMap()
	if err != nil {
		return "", err
	}

	// 3. 再次从缓存获取
	if secretKey, ok := s.versionSecretCache.Get(cacheKey); ok {
		return secretKey, nil
	}

	return "", nil
}
