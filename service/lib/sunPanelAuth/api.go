package sunpanelauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"sun-panel/lib/sunStore/request"

	sun_api "cnb.cool/hslr-s/go-pkg/sun-api"
)

// ==================== 请求/响应结构 ====================

type CaptchaLoginRequest struct {
	Captcha string `json:"captcha"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	HeadImage string `json:"headImage"`
	Status    int    `json:"status"`
	Mail      string `json:"mail"`
	Lang      string `json:"lang"`
	TimeZone  string `json:"timeZone"`
}

type CaptchaLoginResponse struct {
	UserInfo UserInfo `json:"userInfo"`
}

type VersionSecretMapResponse map[string]string

// ==================== 客户端 ====================

type Client struct {
	BaseURL   string
	SecretKey string

	mu          sync.RWMutex
	cachedToken string
	tokenExpiry time.Time
}

func NewClient(baseURL, secretKey string) *Client {
	return &Client{
		BaseURL:   strings.TrimRight(baseURL, "/"),
		SecretKey: secretKey,
	}
}

// ==================== JWT ====================

func (c *Client) getToken() (string, error) {
	c.mu.RLock()
	if c.cachedToken != "" && time.Now().Before(c.tokenExpiry) {
		defer c.mu.RUnlock()
		return c.cachedToken, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cachedToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.cachedToken, nil
	}

	claims := sun_api.NewMapClaims("sun-panel-auth", 24*time.Hour)
	claims.Issuer = "sun-panel-auth"

	token, err := sun_api.GetJwtTokenHS256(c.SecretKey, claims)
	if err != nil {
		return "", fmt.Errorf("generate jwt token: %w", err)
	}

	c.cachedToken = token
	c.tokenExpiry = time.Now().Add(23*time.Hour + 55*time.Minute)

	return token, nil
}

// ==================== 请求方法 ====================

// doPost 发送 POST 请求并解析响应
func (c *Client) doPost(path string, reqBody interface{}, respData interface{}) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}

	url := c.BaseURL + path
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	respContent, httpResp, err := request.SendPostHeaderRequest(url, headers, body)
	if err != nil {
		return err
	}
	if httpResp.StatusCode == http.StatusInternalServerError {
		return errors.New("server error (500)")
	}

	// 两阶段解析
	resp := request.RequestRespDebug{}
	if err := json.Unmarshal(respContent, &resp); err != nil {
		return fmt.Errorf("parse response: %w, body: %s", err, string(respContent))
	}
	if resp.Code != 0 {
		return fmt.Errorf("api error [%d]: %s", resp.Code, resp.Msg)
	}
	if respData != nil && resp.Data != nil {
		dataBytes, _ := json.Marshal(resp.Data)
		if err := json.Unmarshal(dataBytes, respData); err != nil {
			return fmt.Errorf("parse data: %w", err)
		}
	}

	return nil
}

// ==================== 公开接口 ====================

// CaptchaLogin POST /microAppStore/iframe/captchaLogin
func (c *Client) CaptchaLogin(captcha string) (*CaptchaLoginResponse, error) {
	var resp CaptchaLoginResponse
	err := c.doPost("/api/microAppStore/iframe/captchaLogin",
		&CaptchaLoginRequest{Captcha: captcha}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// VersionSecretMap POST /microAppStore/core/versionSecretMap
func (c *Client) VersionSecretMap() (VersionSecretMapResponse, error) {
	var resp VersionSecretMapResponse
	err := c.doPost("/api//microAppStore/core/versionSecretMap", struct{}{}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
