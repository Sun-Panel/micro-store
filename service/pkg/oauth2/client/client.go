package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidConfig        = errors.New("invalid configuration")
	ErrTokenRequestFailed   = errors.New("token request failed")
	ErrInvalidTokenResponse = errors.New("invalid token response")
	ErrRefreshTokenFailed   = errors.New("refresh token failed")
)

// OAuth2Client OAuth2 客户端
type OAuth2Client struct {
	config     *Config
	httpClient *http.Client
}

// NewOAuth2Client 创建 OAuth2 客户端
func NewOAuth2Client(config *Config) (*OAuth2Client, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("%w: client_id and client_secret are required", ErrInvalidConfig)
	}

	if config.Timeout == 0 {
		config.Timeout = 30
	}

	return &OAuth2Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}, nil
}

// GetAuthorizationURL 获取授权 URL
func (c *OAuth2Client) GetAuthorizationURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	if state != "" {
		params.Set("state", state)
	}

	return fmt.Sprintf("%s/oauth2/v1/authorize?%s", c.config.AuthServerURL, params.Encode())
}

// GetAccessTokenByCode 使用授权码获取 Access Token
func (c *OAuth2Client) GetAccessTokenByCode(ctx context.Context, code string) (*TokenResponse, error) {
	req := TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		RedirectURI:  c.config.RedirectURI,
	}

	return c.requestToken(ctx, "/oauth2/v1/token", req)
}

// GetAccessTokenByPassword 使用密码模式获取 Access Token
func (c *OAuth2Client) GetAccessTokenByPassword(ctx context.Context, username, password string) (*TokenResponse, error) {
	req := TokenRequest{
		GrantType:    "password",
		Username:     username,
		Password:     password,
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
	}

	return c.requestToken(ctx, "/oauth2/v1/token", req)
}

// GetClientCredentialsToken 使用客户端凭证模式获取 Access Token
func (c *OAuth2Client) GetClientCredentialsToken(ctx context.Context) (*TokenResponse, error) {
	req := TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
	}

	return c.requestToken(ctx, "/oauth2/v1/clientCredentials/token", req)
}

// RefreshAccessToken 刷新 Access Token
func (c *OAuth2Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	req := TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
	}

	return c.requestToken(ctx, "/oauth2/v1/clientCredentials/token", req)
}

// requestToken 请求 Token
func (c *OAuth2Client) requestToken(ctx context.Context, endpoint string, req TokenRequest) (*TokenResponse, error) {
	url := c.config.AuthServerURL + endpoint

	// 构造请求体
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenRequestFailed, err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("%w: %s - %s", ErrTokenRequestFailed, errResp.Error, errResp.ErrorDescription)
		}
		return nil, fmt.Errorf("%w: status code %d", ErrTokenRequestFailed, resp.StatusCode)
	}

	// ==================== 兼容性代码开始 ====================
	// TODO: 临时兼容方案 - 当授权中心迁移到标准 OAuth2 响应格式后删除此段代码
	// 
	// 背景：当前授权中心返回的是包装格式 {"code":0,"data":{...},"msg":"OK"}
	//       而标准 OAuth2 应该直接返回 {"access_token":"...","token_type":"Bearer",...}
	// 
	// 修复时间：2026-03-13
	// 修复原因：授权中心使用了 apiReturn 包装响应，导致客户端无法正确解析 token
	// 
	// 删除条件：当授权中心的 /oauth2/v1/token 端点改为直接返回标准 OAuth2 格式后
	//          可以删除此兼容代码，只保留直接解析部分

	// 尝试解析包装格式 {"code":0,"data":{...},"msg":"OK"}
	var wrappedResp struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &wrappedResp); err == nil && wrappedResp.Code == 0 && len(wrappedResp.Data) > 0 {
		// 如果是包装格式，解析 data 字段
		if err := json.Unmarshal(wrappedResp.Data, &tokenResp); err != nil {
			return nil, fmt.Errorf("%w: failed to parse data field: %v", ErrInvalidTokenResponse, err)
		}
	} else {
		// 如果不是包装格式，直接解析（标准 OAuth2 格式）
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidTokenResponse, err)
		}
	}
	// ==================== 兼容性代码结束 ====================

	return &tokenResp, nil
}

// APIClient API 客户端（使用 Access Token 调用 API）
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient 创建 API 客户端
func NewAPIClient(baseURL string, timeout int) *APIClient {
	if timeout == 0 {
		timeout = 30
	}

	return &APIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Call 调用 API
func (c *APIClient) Call(ctx context.Context, method, endpoint string, accessToken string, requestData interface{}, responseData interface{}) error {
	url := c.baseURL + endpoint

	// 序列化请求体
	var bodyReader io.Reader
	if requestData != nil {
		bodyBytes, err := json.Marshal(requestData)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	if responseData != nil {
		if err := json.Unmarshal(body, responseData); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get GET 请求
func (c *APIClient) Get(ctx context.Context, endpoint string, accessToken string, responseData interface{}) error {
	return c.Call(ctx, "GET", endpoint, accessToken, nil, responseData)
}

// Post POST 请求
func (c *APIClient) Post(ctx context.Context, endpoint string, accessToken string, requestData interface{}, responseData interface{}) error {
	return c.Call(ctx, "POST", endpoint, accessToken, requestData, responseData)
}

// GetUserInfo 获取用户信息
func (c *APIClient) GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponse, error) {
	var userInfo UserInfoResponse
	err := c.Get(ctx, "/api/v1/user/info", accessToken, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}
