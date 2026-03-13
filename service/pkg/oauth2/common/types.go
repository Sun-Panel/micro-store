package common

import "github.com/golang-jwt/jwt/v4"

// OAuth2 授权类型
const (
	GrantTypeAuthorizationCode = "authorization_code" // 授权码模式
	GrantTypePassword          = "password"           // 密码模式
	GrantTypeClientCredentials = "client_credentials" // 客户端模式
	GrantTypeRefreshToken      = "refresh_token"      // 刷新令牌
)

// 响应类型
const (
	ResponseTypeCode = "code" // 授权码模式
)

// ==================== 请求结构体 ====================

// AuthorizationRequest 授权请求
type AuthorizationRequest struct {
	ClientID     string `json:"client_id" form:"client_id"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
	ResponseType string `json:"response_type" form:"response_type"`
	Scope        string `json:"scope" form:"scope"`
	State        string `json:"state" form:"state"`
}

// TokenRequest Token 请求通用结构
type TokenRequest struct {
	GrantType string `json:"grant_type" form:"grant_type" binding:"required"`
}

// AuthCodeTokenRequest 授权码模式 Token 请求
type AuthCodeTokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	Code         string `json:"code" form:"code" binding:"required"`
	ClientID     string `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
}

// PasswordTokenRequest 密码模式 Token 请求
type PasswordTokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	Username     string `json:"username" form:"username" binding:"required"`
	Password     string `json:"password" form:"password" binding:"required"`
	ClientID     string `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
	Scope        string `json:"scope" form:"scope"`
	State        string `json:"state" form:"state"`
}

// ClientCredentialsTokenRequest 客户端模式 Token 请求
type ClientCredentialsTokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	ClientID     string `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
	Scope        string `json:"scope" form:"scope"`
	State        string `json:"state" form:"state"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" binding:"required"`
	ClientID     string `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
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

// ==================== JWT Claims ====================

// AccessTokenClaims Access Token 的 JWT Claims
type AccessTokenClaims struct {
	UserID   uint   `json:"user_id"`
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid,omitempty"`
	jwt.RegisteredClaims
}

// AuthCodeClaims 授权码的 JWT Claims
type AuthCodeClaims struct {
	Code   string `json:"code"`
	UserID uint   `json:"user_id"`
	jwt.RegisteredClaims
}

// ==================== 错误定义 ====================

const (
	ErrInvalidRequest         = "invalid_request"
	ErrInvalidClient          = "invalid_client"
	ErrInvalidGrant           = "invalid_grant"
	ErrUnauthorizedClient     = "unauthorized_client"
	ErrUnsupportedGrantType   = "unsupported_grment_type"
	ErrInvalidScope           = "invalid_scope"
	ErrAccessDenied           = "access_denied"
	ErrServerError            = "server_error"
	ErrTemporarilyUnavailable = "temporarily_unavailable"
)

// NewErrorResponse 创建错误响应
func NewErrorResponse(errType, description string) ErrorResponse {
	return ErrorResponse{
		Error:            errType,
		ErrorDescription: description,
	}
}
