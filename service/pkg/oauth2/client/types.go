package client

import "time"

// Config OAuth2 客户端配置
type Config struct {
	// 授权服务器地址
	AuthServerURL string

	// API 服务器地址
	APIServerURL string

	// 客户端 ID
	ClientID string

	// 客户端密钥
	ClientSecret string

	// 重定向 URI
	RedirectURI string

	// HTTP 超时时间（秒）
	Timeout int
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Timeout: 30,
	}
}

// ==================== 请求结构体 ====================

// TokenRequest Token 请求
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	State        string `json:"state,omitempty"`
}

// ==================== 响应结构体 ====================

// TokenResponse Token 响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
}

// ClientInfo 客户端信息缓存
type ClientInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
