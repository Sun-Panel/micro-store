package example

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"sun-panel-micro-store/pkg/oauth2/client"
	"sun-panel-micro-store/pkg/oauth2/server"
)

// ==================== 服务端使用示例 ====================

// ThirdApp 第三方应用模型
type ThirdApp struct {
	ClientID      string
	ClientSecret  string
	Name          string
	Enabled       bool
	SSOLogout     bool
	SSOLogoutURL  string
}

func (t *ThirdApp) GetClientID() string       { return t.ClientID }
func (t *ThirdApp) GetClientSecret() string   { return t.ClientSecret }
func (t *ThirdApp) IsEnabled() bool           { return t.Enabled }
func (t *ThirdApp) IsSSOLogout() bool         { return t.SSOLogout }
func (t *ThirdApp) GetSSOLogoutURL() string   { return t.SSOLogoutURL }

// User 用户模型
type User struct {
	ID       uint
	Username string
	Password string
	Name     string
	Email    string
	Status   int
}

func (u *User) GetUserID() uint      { return u.ID }
func (u *User) GetUsername() string  { return u.Username }
func (u *User) GetPassword() string  { return u.Password }
func (u *User) IsActive() bool       { return u.Status == 1 }

// ThirdAppStoreImpl 第三方应用存储实现（内存版）
type ThirdAppStoreImpl struct {
	apps map[string]*ThirdApp
	mu   sync.RWMutex
}

func NewThirdAppStore() *ThirdAppStoreImpl {
	return &ThirdAppStoreImpl{
		apps: make(map[string]*ThirdApp),
	}
}

func (s *ThirdAppStoreImpl) GetByClientID(clientID string) (server.ThirdAppInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	app, ok := s.apps[clientID]
	if !ok {
		return nil, errors.New("third app not found")
	}
	return app, nil
}

func (s *ThirdAppStoreImpl) GetByClientIDAndSecret(clientID, clientSecret string) (server.ThirdAppInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	app, ok := s.apps[clientID]
	if !ok || app.ClientSecret != clientSecret {
		return nil, errors.New("invalid client credentials")
	}
	return app, nil
}

func (s *ThirdAppStoreImpl) AddApp(app *ThirdApp) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apps[app.ClientID] = app
}

// UserStoreImpl 用户存储实现（内存版）
type UserStoreImpl struct {
	users map[uint]*User
	mu    sync.RWMutex
}

func NewUserStore() *UserStoreImpl {
	return &UserStoreImpl{
		users: make(map[uint]*User),
	}
}

func (s *UserStoreImpl) GetByUsernameAndPassword(username, password string) (server.UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, user := range s.users {
		if user.Username == username && user.Password == password {
			return user, nil
		}
	}
	return nil, errors.New("invalid credentials")
}

func (s *UserStoreImpl) GetByID(userID uint) (server.UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, ok := s.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserStoreImpl) AddUser(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}

// TokenStoreImpl Token 存储实现（内存版）
type TokenStoreImpl struct {
	tokens map[string]server.AccessTokenData
	mu     sync.RWMutex
}

func NewTokenStore() *TokenStoreImpl {
	return &TokenStoreImpl{
		tokens: make(map[string]server.AccessTokenData),
	}
}

func (s *TokenStoreImpl) SetAccessToken(token string, data server.AccessTokenData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = data
	return nil
}

func (s *TokenStoreImpl) GetAccessToken(token string) (server.AccessTokenData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, ok := s.tokens[token]
	if !ok {
		return server.AccessTokenData{}, errors.New("token not found")
	}
	return data, nil
}

func (s *TokenStoreImpl) DeleteAccessToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
	return nil
}

// AuthCodeStoreImpl 授权码存储实现（内存版）
type AuthCodeStoreImpl struct {
	codes map[string]server.OAuthCodeData
	mu    sync.RWMutex
}

func NewAuthCodeStore() *AuthCodeStoreImpl {
	return &AuthCodeStoreImpl{
		codes: make(map[string]server.OAuthCodeData),
	}
}

func (s *AuthCodeStoreImpl) SetAuthCode(code string, data server.OAuthCodeData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code] = data
	return nil
}

func (s *AuthCodeStoreImpl) GetAuthCode(code string) (server.OAuthCodeData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, ok := s.codes[code]
	if !ok {
		return server.OAuthCodeData{}, errors.New("auth code not found")
	}
	return data, nil
}

func (s *AuthCodeStoreImpl) DeleteAuthCode(code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, code)
	return nil
}

// RefreshTokenStoreImpl 刷新令牌存储实现（内存版）
type RefreshTokenStoreImpl struct {
	tokens map[string]server.RefreshTokenData
	mu     sync.RWMutex
}

func NewRefreshTokenStore() *RefreshTokenStoreImpl {
	return &RefreshTokenStoreImpl{
		tokens: make(map[string]server.RefreshTokenData),
	}
}

func (s *RefreshTokenStoreImpl) SetRefreshToken(token string, data server.RefreshTokenData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = data
	return nil
}

func (s *RefreshTokenStoreImpl) GetRefreshToken(token string) (server.RefreshTokenData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, ok := s.tokens[token]
	if !ok {
		return server.RefreshTokenData{}, errors.New("refresh token not found")
	}
	return data, nil
}

func (s *RefreshTokenStoreImpl) DeleteRefreshToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
	return nil
}

// SetupOAuthServer 创建 OAuth2 服务端
func SetupOAuthServer() *server.OAuthHandler {
	// 创建配置
	config := &server.OAuthConfig{
		AccessTokenExpireTime:  7200,       // 2小时
		RefreshTokenExpireTime: 604800,    // 7天
		AuthCodeExpireTime:     600,       // 10分钟
		EnableSSOLogout:        true,
	}

	// 创建 handler
	handler := server.NewOAuthHandler(config)

	// 创建存储实例
	thirdAppStore := NewThirdAppStore()
	userStore := NewUserStore()
	tokenStore := NewTokenStore()
	authCodeStore := NewAuthCodeStore()
	refreshTokenStore := NewRefreshTokenStore()

	// 添加测试数据
	thirdAppStore.AddApp(&ThirdApp{
		ClientID:     "test_client",
		ClientSecret: "test_secret",
		Name:         "测试应用",
		Enabled:      true,
		SSOLogout:    true,
		SSOLogoutURL: "http://localhost:8081/sso/logout",
	})

	userStore.AddUser(&User{
		ID:       1,
		Username: "admin",
		Password: "admin123",
		Name:     "管理员",
		Email:    "admin@example.com",
		Status:   1,
	})

	// 设置存储实现
	handler.SetStores(thirdAppStore, userStore, tokenStore, authCodeStore, refreshTokenStore)

	return handler
}

// StartOAuthServerExample 启动 OAuth2 服务端示例
func StartOAuthServerExample() {
	handler := SetupOAuthServer()

	// 创建 Gin 路由
	r := gin.Default()

	// 注册 OAuth2 路由
	handler.RegisterRoutes(r.Group("/api"))

	// 添加用户授权登录接口
	r.POST("/api/auth/login", func(c *gin.Context) {
		// 这里应该验证用户登录状态，并注入 user_id 到上下文
		// 示例中直接设置 user_id = 1
		c.Set("user_id", uint(1))
		handler.AuthLogin(c)
	})

	fmt.Println("OAuth2 Server running on :8080")
	r.Run(":8080")
}

// ==================== 客户端使用示例 ====================

// OAuth2ClientExample OAuth2 客户端使用示例
func OAuth2ClientExample() {
	// 创建客户端配置
	config := &client.Config{
		AuthServerURL: "http://localhost:8080",
		APIServerURL:  "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		RedirectURI:   "http://localhost:3000/callback",
		Timeout:       30,
	}

	// 创建客户端
	oauthClient, err := client.NewOAuth2Client(config)
	if err != nil {
		panic(err)
	}

	// 1. 获取授权 URL
	authURL := oauthClient.GetAuthorizationURL(config.RedirectURI, "random_state")
	fmt.Printf("授权 URL: %s\n", authURL)
	fmt.Println("请访问此 URL 进行授权，授权后会回调到 redirect_uri 并带上 code 参数")

	// 2. 使用授权码获取 Token（模拟）
	// code := "从回调中获取的授权码"
	// tokenResp, err := oauthClient.GetAccessTokenByCode(context.Background(), code)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
	// fmt.Printf("Refresh Token: %s\n", tokenResp.RefreshToken)

	// 3. 使用密码模式获取 Token
	fmt.Println("\n使用密码模式获取 Token:")
	tokenResp, err := oauthClient.GetAccessTokenByPassword(context.Background(), "admin", "admin123")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
		fmt.Printf("Refresh Token: %s\n", tokenResp.RefreshToken)
		fmt.Printf("Expires In: %d\n", tokenResp.ExpiresIn)
	}

	// 4. 使用客户端凭证模式获取 Token
	fmt.Println("\n使用客户端凭证模式获取 Token:")
	clientToken, err := oauthClient.GetClientCredentialsToken(context.Background())
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("Access Token: %s\n", clientToken.AccessToken)
		fmt.Printf("Expires In: %d\n", clientToken.ExpiresIn)
	}

	// 5. 刷新 Token
	// newToken, err := oauthClient.RefreshAccessToken(context.Background(), tokenResp.RefreshToken)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("New Access Token: %s\n", newToken.AccessToken)

	// 6. 调用 API
	// apiClient := client.NewAPIClient("http://localhost:8080", 30)
	// userInfo, err := apiClient.GetUserInfo(context.Background(), tokenResp.AccessToken)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("用户信息: %+v\n", userInfo)
}

// ==================== 完整流程示例 ====================

// CompleteOAuthFlowExample 完整 OAuth2 流程示例
func CompleteOAuthFlowExample() {
	fmt.Println("========== OAuth2 完整流程示例 ==========")

	// 场景 1: 客户端凭证模式
	fmt.Println("\n场景 1: 客户端凭证模式")
	fmt.Println("适用于服务端对服务端的调用，无需用户参与")

	config := &client.Config{
		AuthServerURL: "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
	}

	oauthClient, _ := client.NewOAuth2Client(config)
	token, err := oauthClient.GetClientCredentialsToken(context.Background())
	if err != nil {
		fmt.Printf("获取 Token 失败: %v\n", err)
	} else {
		fmt.Printf("成功获取 Access Token: %s\n", token.AccessToken)
	}

	// 场景 2: 授权码模式
	fmt.Println("\n场景 2: 授权码模式")
	fmt.Println("适用于第三方应用访问用户资源")

	config2 := &client.Config{
		AuthServerURL: "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		RedirectURI:   "http://localhost:3000/callback",
	}

	oauthClient2, _ := client.NewOAuth2Client(config2)
	authURL := oauthClient2.GetAuthorizationURL(config2.RedirectURI, "random_state")
	fmt.Printf("授权 URL: %s\n", authURL)
	fmt.Println("用户访问此 URL 进行授权")
	fmt.Println("授权后会重定向到: http://localhost:3000/callback?code=AUTHORIZATION_CODE")
	fmt.Println("然后使用 code 获取 access_token")

	// 场景 3: 密码模式
	fmt.Println("\n场景 3: 密码模式")
	fmt.Println("适用于受信任的第一方应用")

	config3 := &client.Config{
		AuthServerURL: "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
	}

	oauthClient3, _ := client.NewOAuth2Client(config3)
	token3, err := oauthClient3.GetAccessTokenByPassword(context.Background(), "admin", "admin123")
	if err != nil {
		fmt.Printf("获取 Token 失败: %v\n", err)
	} else {
		fmt.Printf("成功获取 Access Token: %s\n", token3.AccessToken)
		fmt.Printf("Refresh Token: %s\n", token3.RefreshToken)
	}

	// 场景 4: 刷新 Token
	fmt.Println("\n场景 4: 刷新 Token")
	fmt.Println("当 Access Token 过期时，使用 Refresh Token 获取新的 Access Token")
	// newToken, err := oauthClient.RefreshAccessToken(context.Background(), token3.RefreshToken)

	fmt.Println("\n========== 示例结束 ==========")
}
