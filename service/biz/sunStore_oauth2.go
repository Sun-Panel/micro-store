package biz

import (
	"context"

	"sun-panel/global"
	"github.com/sunjingliang/oauth2-go/client"
)

// SunStoreOAuth2Client 新的 OAuth2 客户端封装
type SunStoreOAuth2Client struct {
	oauthClient *client.OAuth2Client
	apiClient   *client.APIClient
}

// NewSunStoreOAuth2Client 创建新的 OAuth2 客户端
func NewSunStoreOAuth2Client() *SunStoreOAuth2Client {
	clientID := global.Config.GetValueString("sun_store", "client_id")
	clientSecret := global.Config.GetValueString("sun_store", "client_secret")
	apiHost := global.Config.GetValueString("sun_store", "api_host")

	config := &client.Config{
		AuthServerURL: apiHost,
		APIServerURL:  apiHost,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		Timeout:       30,
	}

	oauthClient, err := client.NewOAuth2Client(config)
	if err != nil {
		global.Logger.Errorln("创建 OAuth2 客户端失败:", err)
		return nil
	}

	apiClient := client.NewAPIClient(apiHost, 30)

	return &SunStoreOAuth2Client{
		oauthClient: oauthClient,
		apiClient:   apiClient,
	}
}

// GetAccessTokenByCode 使用授权码获取 Access Token
func (c *SunStoreOAuth2Client) GetAccessTokenByCode(code string) (*client.TokenResponse, error) {
	ctx := context.Background()
	return c.oauthClient.GetAccessTokenByCode(ctx, code)
}

// GetClientCredentialsToken 使用客户端凭证模式获取 Token
func (c *SunStoreOAuth2Client) GetClientCredentialsToken() (*client.TokenResponse, error) {
	ctx := context.Background()
	return c.oauthClient.GetClientCredentialsToken(ctx)
}

// PasswordAuth 密码模式授权
func (c *SunStoreOAuth2Client) PasswordAuth(username, password string) (*client.TokenResponse, error) {
	ctx := context.Background()
	return c.oauthClient.GetAccessTokenByPassword(ctx, username, password)
}

// RefreshAccessToken 刷新 Access Token
func (c *SunStoreOAuth2Client) RefreshAccessToken(refreshToken string) (*client.TokenResponse, error) {
	ctx := context.Background()
	return c.oauthClient.RefreshAccessToken(ctx, refreshToken)
}

// GetAuthorizationURL 获取授权 URL
func (c *SunStoreOAuth2Client) GetAuthorizationURL(redirectURI, state string) string {
	return c.oauthClient.GetAuthorizationURL(redirectURI, state)
}

// APICall 调用 API
func (c *SunStoreOAuth2Client) APICall(ctx context.Context, method, endpoint string, accessToken string, requestData interface{}, responseData interface{}) error {
	return c.apiClient.Call(ctx, method, endpoint, accessToken, requestData, responseData)
}
