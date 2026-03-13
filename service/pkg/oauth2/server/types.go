package server

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ThirdAppInfo 第三方应用信息接口
type ThirdAppInfo interface {
	GetClientID() string
	GetClientSecret() string
	IsEnabled() bool
	IsSSOLogout() bool
	GetSSOLogoutURL() string
}

// UserInfo 用户信息接口
type UserInfo interface {
	GetUserID() uint
	GetUsername() string
	GetPassword() string
	IsActive() bool
}

// OAuthCodeData OAuth 授权码数据
type OAuthCodeData struct {
	Code        string    `json:"code"`
	AccessToken string    `json:"access_token"`
	CToken      string    `json:"c_token"`
	ClientID    string    `json:"client_id"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// AccessTokenData Access Token 数据
type AccessTokenData struct {
	UserID   uint   `json:"user_id"`
	ClientID string `json:"client_id"`
	CToken   string `json:"c_token"`
	OpenID   string `json:"openid,omitempty"`
}

// ClientCredentialsTokenData 客户端凭证模式 Token 数据
type ClientCredentialsTokenData struct {
	ClientID string `json:"client_id"`
}

// RefreshTokenData 刷新令牌数据
type RefreshTokenData struct {
	AccessToken string `json:"access_token"`
	ClientID    string `json:"client_id"`
}

// SSOLogoutResult 单点登出结果
type SSOLogoutResult struct {
	ClientID string
	Error    error
}

// OAuthConfig OAuth2 服务端配置
type OAuthConfig struct {
	// Token 过期时间（秒）
	AccessTokenExpireTime  int
	RefreshTokenExpireTime int
	AuthCodeExpireTime     int

	// JWT 签名密钥
	JWTSecret []byte

	// 是否启用单点登出
	EnableSSOLogout bool
}

// DefaultOAuthConfig 默认配置
func DefaultOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		AccessTokenExpireTime:  7200,        // 2小时
		RefreshTokenExpireTime: 604800,     // 7天
		AuthCodeExpireTime:     600,        // 10分钟
		EnableSSOLogout:        true,
	}
}

// OAuthClaims 自定义 JWT Claims
type OAuthClaims struct {
	UserID   uint   `json:"user_id"`
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid,omitempty"`
	jwt.RegisteredClaims
}

// AuthCodeClaims 授权码 Claims
type AuthCodeClaims struct {
	Code     string `json:"code"`
	UserID   uint   `json:"user_id"`
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}
