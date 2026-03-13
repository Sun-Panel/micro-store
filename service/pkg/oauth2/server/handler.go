package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"sun-panel-micro-store/pkg/oauth2/common"
)

// OAuthHandler OAuth2 处理器
type OAuthHandler struct {
	tokenManager *TokenManager
	config       *OAuthConfig

	// 数据获取接口
	thirdAppStore  ThirdAppStore
	userStore      UserStore
	tokenStore     TokenStore
	authCodeStore  AuthCodeStore
	refreshTokenStore RefreshTokenStore
}

// ThirdAppStore 第三方应用存储接口
type ThirdAppStore interface {
	GetByClientID(clientID string) (ThirdAppInfo, error)
	GetByClientIDAndSecret(clientID, clientSecret string) (ThirdAppInfo, error)
}

// UserStore 用户存储接口
type UserStore interface {
	GetByUsernameAndPassword(username, password string) (UserInfo, error)
	GetByID(userID uint) (UserInfo, error)
}

// TokenStore Token 存储接口
type TokenStore interface {
	SetAccessToken(token string, data AccessTokenData) error
	GetAccessToken(token string) (AccessTokenData, error)
	DeleteAccessToken(token string) error
}

// AuthCodeStore 授权码存储接口
type AuthCodeStore interface {
	SetAuthCode(code string, data OAuthCodeData) error
	GetAuthCode(code string) (OAuthCodeData, error)
	DeleteAuthCode(code string) error
}

// RefreshTokenStore 刷新令牌存储接口
type RefreshTokenStore interface {
	SetRefreshToken(token string, data RefreshTokenData) error
	GetRefreshToken(token string) (RefreshTokenData, error)
	DeleteRefreshToken(token string) error
}

// NewOAuthHandler 创建 OAuth 处理器
func NewOAuthHandler(config *OAuthConfig) *OAuthHandler {
	if config == nil {
		config = DefaultOAuthConfig()
	}
	return &OAuthHandler{
		tokenManager: NewTokenManager(config),
		config:       config,
	}
}

// SetStores 设置存储接口
func (h *OAuthHandler) SetStores(
	thirdAppStore ThirdAppStore,
	userStore UserStore,
	tokenStore TokenStore,
	authCodeStore AuthCodeStore,
	refreshTokenStore RefreshTokenStore,
) {
	h.thirdAppStore = thirdAppStore
	h.userStore = userStore
	h.tokenStore = tokenStore
	h.authCodeStore = authCodeStore
	h.refreshTokenStore = refreshTokenStore
}

// Authorize 授权端点 - 重定向到前端授权页面
func (h *OAuthHandler) Authorize(c *gin.Context) {
	var req common.AuthorizationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrInvalidRequest,
			"invalid request parameters",
		))
		return
	}

	// 验证 client_id
	_, err := h.thirdAppStore.GetByClientID(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client_id",
		))
		return
	}

	// 构造重定向 URL（这里需要根据实际前端路径配置）
	// 将请求参数传递给前端授权页面
	queryParams := c.Request.URL.Query().Encode()
	frontAuthURL := fmt.Sprintf("/auth?%s", queryParams)

	c.Redirect(http.StatusFound, frontAuthURL)
}

// AuthLogin 用户授权登录（已登录用户授权）
func (h *OAuthHandler) AuthLogin(c *gin.Context) {
	// 从上下文获取当前用户信息（需要通过中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrUnauthorizedClient,
			"user not logged in",
		))
		return
	}

	var req struct {
		ClientID string `json:"clientId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrInvalidRequest,
			"invalid request parameters",
		))
		return
	}

	// 获取第三方应用信息
	appInfo, err := h.thirdAppStore.GetByClientID(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client_id",
		))
		return
	}

	if !appInfo.IsEnabled() {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrUnauthorizedClient,
			"client is disabled",
		))
		return
	}

	// 生成 Access Token
	accessToken, err := h.tokenManager.GenerateAccessToken(
		req.ClientID,
		userID.(uint),
		appInfo.GetClientSecret(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			common.ErrServerError,
			"failed to generate access token",
		))
		return
	}

	// 生成授权码
	authCode, err := h.tokenManager.GenerateAuthCode(
		req.ClientID,
		userID.(uint),
		appInfo.GetClientSecret(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			common.ErrServerError,
			"failed to generate auth code",
		))
		return
	}

	// 缓存授权码
	_ = h.authCodeStore.SetAuthCode(authCode, OAuthCodeData{
		Code:        authCode,
		AccessToken: accessToken,
		ClientID:    req.ClientID,
		UserID:      userID.(uint),
		CreatedAt:   time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": authCode,
	})
}

// Token Token 端点
func (h *OAuthHandler) Token(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	// 读取请求体
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	var commonReq common.TokenRequest

	// 根据 Content-Type 绑定参数
	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&commonReq); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&commonReq); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 根据 grant_type 处理不同的授权模式
	switch commonReq.GrantType {
	case common.GrantTypeAuthorizationCode:
		h.handleAuthCodeGrant(c, bodyBytes, contentType)
	case common.GrantTypePassword:
		h.handlePasswordGrant(c, bodyBytes, contentType)
	case common.GrantTypeRefreshToken:
		h.handleRefreshTokenGrant(c, bodyBytes, contentType)
	default:
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrUnsupportedGrantType,
			"unsupported grant_type",
		))
	}
}

// handleAuthCodeGrant 处理授权码模式
func (h *OAuthHandler) handleAuthCodeGrant(c *gin.Context, bodyBytes []byte, contentType string) {
	var req common.AuthCodeTokenRequest

	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 验证客户端
	appInfo, err := h.thirdAppStore.GetByClientIDAndSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client credentials",
		))
		return
	}

	// 获取授权码数据
	codeData, err := h.authCodeStore.GetAuthCode(req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidGrant,
			"invalid authorization code",
		))
		return
	}

	// 删除授权码（一次性使用）
	defer h.authCodeStore.DeleteAuthCode(req.Code)

	// 验证 Access Token
	_, err = h.tokenManager.ValidateAccessToken(codeData.AccessToken, appInfo.GetClientSecret())
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidGrant,
			"invalid access token",
		))
		return
	}

	// 缓存 Access Token
	_ = h.tokenStore.SetAccessToken(codeData.AccessToken, AccessTokenData{
		UserID:   codeData.UserID,
		ClientID: codeData.ClientID,
		CToken:   codeData.CToken,
	})

	c.JSON(http.StatusOK, common.TokenResponse{
		AccessToken: codeData.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.tokenManager.GetExpireTime(),
	})
}

// handlePasswordGrant 处理密码模式
func (h *OAuthHandler) handlePasswordGrant(c *gin.Context, bodyBytes []byte, contentType string) {
	var req common.PasswordTokenRequest

	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 验证客户端
	appInfo, err := h.thirdAppStore.GetByClientIDAndSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client credentials",
		))
		return
	}

	// 验证用户
	userInfo, err := h.userStore.GetByUsernameAndPassword(req.Username, req.Password)
	if err != nil || !userInfo.IsActive() {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidGrant,
			"invalid user credentials",
		))
		return
	}

	// 生成 Access Token
	accessToken, err := h.tokenManager.GenerateAccessToken(
		req.ClientID,
		userInfo.GetUserID(),
		appInfo.GetClientSecret(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			common.ErrServerError,
			"failed to generate access token",
		))
		return
	}

	// 生成 Refresh Token
	refreshToken := h.tokenManager.GenerateRefreshToken()

	// 缓存 Access Token
	_ = h.tokenStore.SetAccessToken(accessToken, AccessTokenData{
		UserID:   userInfo.GetUserID(),
		ClientID: req.ClientID,
	})

	// 缓存 Refresh Token
	_ = h.refreshTokenStore.SetRefreshToken(refreshToken, RefreshTokenData{
		AccessToken: accessToken,
		ClientID:    req.ClientID,
	})

	c.JSON(http.StatusOK, common.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.tokenManager.GetExpireTime(),
		RefreshToken: refreshToken,
	})
}

// handleRefreshTokenGrant 处理刷新令牌
func (h *OAuthHandler) handleRefreshTokenGrant(c *gin.Context, bodyBytes []byte, contentType string) {
	var req common.RefreshTokenRequest

	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 验证客户端
	appInfo, err := h.thirdAppStore.GetByClientIDAndSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client credentials",
		))
		return
	}

	// 获取刷新令牌数据
	refreshData, err := h.refreshTokenStore.GetRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidGrant,
			"invalid refresh token",
		))
		return
	}

	// 删除旧的 Access Token
	_ = h.tokenStore.DeleteAccessToken(refreshData.AccessToken)

	// 生成新的 Access Token
	accessToken, err := h.tokenManager.GenerateAccessToken(
		req.ClientID,
		0, // 客户端模式下没有用户ID
		appInfo.GetClientSecret(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			common.ErrServerError,
			"failed to generate access token",
		))
		return
	}

	// 缓存新的 Access Token
	_ = h.tokenStore.SetAccessToken(accessToken, AccessTokenData{
		ClientID: req.ClientID,
	})

	// 更新刷新令牌
	_ = h.refreshTokenStore.SetRefreshToken(req.RefreshToken, RefreshTokenData{
		AccessToken: accessToken,
		ClientID:    req.ClientID,
	})

	c.JSON(http.StatusOK, common.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.tokenManager.GetExpireTime(),
		RefreshToken: req.RefreshToken,
	})
}

// ClientCredentialsToken 客户端凭证模式
func (h *OAuthHandler) ClientCredentialsToken(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	var commonReq common.TokenRequest

	// 根据 Content-Type 绑定参数
	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&commonReq); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&commonReq); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 根据 grant_type 处理
	switch commonReq.GrantType {
	case common.GrantTypeClientCredentials:
		h.handleClientCredentialsGrant(c, bodyBytes, contentType)
	case common.GrantTypeRefreshToken:
		h.handleRefreshTokenGrant(c, bodyBytes, contentType)
	default:
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			common.ErrUnsupportedGrantType,
			"unsupported grant_type",
		))
	}
}

// handleClientCredentialsGrant 处理客户端凭证模式
func (h *OAuthHandler) handleClientCredentialsGrant(c *gin.Context, bodyBytes []byte, contentType string) {
	var req common.ClientCredentialsTokenRequest

	if strings.Contains(contentType, "application/json") {
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				common.ErrInvalidRequest,
				"invalid request",
			))
			return
		}
	}

	// 验证客户端
	appInfo, err := h.thirdAppStore.GetByClientIDAndSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			common.ErrInvalidClient,
			"invalid client credentials",
		))
		return
	}

	// 生成 Access Token
	accessToken, err := h.tokenManager.GenerateAccessToken(
		req.ClientID,
		0, // 客户端模式下没有用户ID
		appInfo.GetClientSecret(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			common.ErrServerError,
			"failed to generate access token",
		))
		return
	}

	// 生成 Refresh Token
	refreshToken := h.tokenManager.GenerateRefreshToken()

	// 缓存 Access Token
	_ = h.tokenStore.SetAccessToken(accessToken, AccessTokenData{
		ClientID: req.ClientID,
	})

	// 缓存 Refresh Token
	_ = h.refreshTokenStore.SetRefreshToken(refreshToken, RefreshTokenData{
		AccessToken: accessToken,
		ClientID:    req.ClientID,
	})

	c.JSON(http.StatusOK, common.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.tokenManager.GetExpireTime(),
		RefreshToken: refreshToken,
	})
}

// SSOLogout 单点登出
func (h *OAuthHandler) SSOLogout(c *gin.Context) {
	var req struct {
		AccessToken string `json:"accessToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusOK, "")
		return
	}

	// 删除 Access Token
	_ = h.tokenStore.DeleteAccessToken(req.AccessToken)

	c.String(http.StatusOK, "")
}

// RegisterRoutes 注册路由
func (h *OAuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	oauth := router.Group("oauth2/v1")
	{
		// 授权端点
		oauth.GET("authorize", h.Authorize)

		// Token 端点
		oauth.POST("token", h.Token)

		// 客户端凭证模式 Token 端点
		oauth.POST("clientCredentials/token", h.ClientCredentialsToken)

		// 单点登出
		oauth.POST("sso/logout", h.SSOLogout)
	}
}
