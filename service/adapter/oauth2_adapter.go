package adapter

import (
	"context"
	"errors"

	"sun-panel/global"
	"github.com/sunjingliang/oauth2-go/client"
)

// ==================== 客户端适配器 ====================

// SunStoreClientAdapter 适配现有的 SunStore 客户端
// 用于替代 service/lib/sunStore/oAuth2.go 中的实现
type SunStoreClientAdapter struct {
	oauthClient *client.OAuth2Client
	apiClient   *client.APIClient
}

// NewSunStoreClientAdapter 创建 SunStore 客户端适配器
// 替代原来的 sunStore.NewSunStoreApi
func NewSunStoreClientAdapter(authServerURL, clientID, clientSecret string) (*SunStoreClientAdapter, error) {
	config := &client.Config{
		AuthServerURL: authServerURL,
		APIServerURL:  authServerURL,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		Timeout:       30,
	}

	oauthClient, err := client.NewOAuth2Client(config)
	if err != nil {
		return nil, err
	}

	apiClient := client.NewAPIClient(authServerURL, 30)

	return &SunStoreClientAdapter{
		oauthClient: oauthClient,
		apiClient:   apiClient,
	}, nil
}

// GetAccessToken 使用授权码获取 Token
// 替代原来的 GetAccessToken 方法
func (a *SunStoreClientAdapter) GetAccessToken(code string) (*client.TokenResponse, error) {
	return a.oauthClient.GetAccessTokenByCode(context.Background(), code)
}

// ClientCredentialsAuth 客户端凭证模式授权
// 替代原来的 ClientCredentialsAuth 方法
func (a *SunStoreClientAdapter) ClientCredentialsAuth() (*client.TokenResponse, error) {
	return a.oauthClient.GetClientCredentialsToken(context.Background())
}

// PasswordAuth 密码模式授权
// 替代原来的 PasswordAuth 方法
func (a *SunStoreClientAdapter) PasswordAuth(username, password string) (*client.TokenResponse, error) {
	return a.oauthClient.GetAccessTokenByPassword(context.Background(), username, password)
}

// RefreshToken 刷新 Token
func (a *SunStoreClientAdapter) RefreshToken(refreshToken string) (*client.TokenResponse, error) {
	return a.oauthClient.RefreshAccessToken(context.Background(), refreshToken)
}

// GetAuthorizationURL 获取授权 URL
func (a *SunStoreClientAdapter) GetAuthorizationURL(redirectURI, state string) string {
	return a.oauthClient.GetAuthorizationURL(redirectURI, state)
}

// APICall 调用 API（通用方法）
func (a *SunStoreClientAdapter) APICall(ctx context.Context, method, endpoint string, accessToken string, requestData interface{}, responseData interface{}) error {
	return a.apiClient.Call(ctx, method, endpoint, accessToken, requestData, responseData)
}

// ==================== 辅助函数 ====================

// GetSunStoreClient 获取 SunStore 客户端实例
// 这是一个便捷函数，直接从配置中读取参数
func GetSunStoreClient() (*SunStoreClientAdapter, error) {
	clientID := global.Config.GetValueString("sun_store", "client_id")
	clientSecret := global.Config.GetValueString("sun_store", "client_secret")
	apiHost := global.Config.GetValueString("sun_store", "api_host")

	if clientID == "" || clientSecret == "" || apiHost == "" {
		return nil, errors.New("sun_store config not found")
	}

	return NewSunStoreClientAdapter(apiHost, clientID, clientSecret)
}
