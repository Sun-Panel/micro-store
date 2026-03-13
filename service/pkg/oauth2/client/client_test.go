package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOAuth2Client(t *testing.T) {
	config := &Config{
		AuthServerURL: "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		Timeout:       30,
	}

	client, err := NewOAuth2Client(config)
	if err != nil {
		t.Fatalf("Failed to create OAuth2 client: %v", err)
	}

	if client == nil {
		t.Error("Client should not be nil")
	}

	if client.config.ClientID != config.ClientID {
		t.Errorf("Expected ClientID %s, got %s", config.ClientID, client.config.ClientID)
	}
}

func TestNewOAuth2Client_MissingConfig(t *testing.T) {
	_, err := NewOAuth2Client(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}
}

func TestNewOAuth2Client_MissingCredentials(t *testing.T) {
	config := &Config{
		AuthServerURL: "http://localhost:8080",
	}

	_, err := NewOAuth2Client(config)
	if err == nil {
		t.Error("Expected error for missing credentials")
	}
}

func TestOAuth2Client_GetAuthorizationURL(t *testing.T) {
	config := &Config{
		AuthServerURL: "http://localhost:8080",
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		RedirectURI:   "http://localhost:3000/callback",
	}

	client, _ := NewOAuth2Client(config)

	redirectURI := "http://localhost:3000/callback"
	state := "random_state"

	authURL := client.GetAuthorizationURL(redirectURI, state)

	if authURL == "" {
		t.Error("Authorization URL should not be empty")
	}

	// 检查 URL 包含必要参数
	if !containsAll(authURL, "client_id=test_client", "redirect_uri", "response_type=code", "state=random_state") {
		t.Errorf("Authorization URL missing required parameters: %s", authURL)
	}

	t.Logf("Authorization URL: %s", authURL)
}

func TestOAuth2Client_GetAccessTokenByCode(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 返回模拟响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "test_access_token",
			"token_type": "Bearer",
			"expires_in": 7200,
			"refresh_token": "test_refresh_token"
		}`))
	}))
	defer server.Close()

	config := &Config{
		AuthServerURL: server.URL,
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		Timeout:       30,
	}

	client, _ := NewOAuth2Client(config)

	ctx := context.Background()
	tokenResp, err := client.GetAccessTokenByCode(ctx, "test_code")
	if err != nil {
		t.Fatalf("Failed to get access token: %v", err)
	}

	if tokenResp.AccessToken != "test_access_token" {
		t.Errorf("Expected access token 'test_access_token', got '%s'", tokenResp.AccessToken)
	}

	if tokenResp.RefreshToken != "test_refresh_token" {
		t.Errorf("Expected refresh token 'test_refresh_token', got '%s'", tokenResp.RefreshToken)
	}

	if tokenResp.ExpiresIn != 7200 {
		t.Errorf("Expected expires_in 7200, got %d", tokenResp.ExpiresIn)
	}

	t.Logf("Token response: %+v", tokenResp)
}

func TestOAuth2Client_GetClientCredentialsToken(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "client_access_token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer server.Close()

	config := &Config{
		AuthServerURL: server.URL,
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		Timeout:       30,
	}

	client, _ := NewOAuth2Client(config)

	ctx := context.Background()
	tokenResp, err := client.GetClientCredentialsToken(ctx)
	if err != nil {
		t.Fatalf("Failed to get client credentials token: %v", err)
	}

	if tokenResp.AccessToken != "client_access_token" {
		t.Errorf("Expected access token 'client_access_token', got '%s'", tokenResp.AccessToken)
	}

	t.Logf("Client credentials token: %+v", tokenResp)
}

func TestOAuth2Client_GetAccessTokenByPassword(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "password_access_token",
			"token_type": "Bearer",
			"expires_in": 7200,
			"refresh_token": "password_refresh_token"
		}`))
	}))
	defer server.Close()

	config := &Config{
		AuthServerURL: server.URL,
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		Timeout:       30,
	}

	client, _ := NewOAuth2Client(config)

	ctx := context.Background()
	tokenResp, err := client.GetAccessTokenByPassword(ctx, "test_user", "test_password")
	if err != nil {
		t.Fatalf("Failed to get password token: %v", err)
	}

	if tokenResp.AccessToken != "password_access_token" {
		t.Errorf("Expected access token 'password_access_token', got '%s'", tokenResp.AccessToken)
	}

	t.Logf("Password token: %+v", tokenResp)
}

func TestOAuth2Client_ErrorHandling(t *testing.T) {
	// 创建返回错误的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": "invalid_client",
			"error_description": "Client authentication failed"
		}`))
	}))
	defer server.Close()

	config := &Config{
		AuthServerURL: server.URL,
		ClientID:      "test_client",
		ClientSecret:  "test_secret",
		Timeout:       30,
	}

	client, _ := NewOAuth2Client(config)

	ctx := context.Background()
	_, err := client.GetClientCredentialsToken(ctx)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	t.Logf("Expected error: %v", err)
}

func TestAPIClient_Get(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test_token" {
			t.Errorf("Expected Authorization header 'Bearer test_token', got '%s'", authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"user_id": 1,
			"username": "test_user",
			"name": "Test User",
			"email": "test@example.com"
		}`))
	}))
	defer server.Close()

	apiClient := NewAPIClient(server.URL, 30)

	ctx := context.Background()
	var userInfo UserInfoResponse
	err := apiClient.Get(ctx, "/api/v1/user/info", "test_token", &userInfo)
	if err != nil {
		t.Fatalf("Failed to call API: %v", err)
	}

	if userInfo.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", userInfo.UserID)
	}

	if userInfo.Username != "test_user" {
		t.Errorf("Expected username 'test_user', got '%s'", userInfo.Username)
	}

	t.Logf("User info: %+v", userInfo)
}

func TestAPIClient_Post(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	apiClient := NewAPIClient(server.URL, 30)

	ctx := context.Background()
	requestData := map[string]interface{}{
		"name": "test",
	}
	var responseData map[string]interface{}

	err := apiClient.Post(ctx, "/api/v1/test", "test_token", requestData, &responseData)
	if err != nil {
		t.Fatalf("Failed to call API: %v", err)
	}

	t.Logf("Response: %+v", responseData)
}

// 辅助函数
func containsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if !contains(s, substr) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
