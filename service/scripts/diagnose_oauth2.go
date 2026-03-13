package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("OAuth2 完整流程诊断")
	fmt.Println("====================================")
	fmt.Println()

	// 配置
	authServerURL := "http://192.168.3.101:3088"
	clientID := "xs1ubmt25p"
	clientSecret := "pKcHMOKBUOvwsOmokvn82JMTC2zR4mIy"
	redirectURI := "http://localhost/callback"

	// 测试 1: 检查授权端点
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试 1: 检查授权端点")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	authorizeURL := fmt.Sprintf("%s/oauth2/v1/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		authServerURL, clientID, url.QueryEscape(redirectURI))
	fmt.Printf("授权 URL: %s\n", authorizeURL)
	
	resp, err := http.Get(authorizeURL)
	if err != nil {
		fmt.Printf("❌ 授权端点不可达: %v\n", err)
	} else {
		fmt.Printf("✅ 授权端点可访问 (状态码: %d)\n", resp.StatusCode)
		resp.Body.Close()
	}
	fmt.Println()

	// 测试 2: 客户端凭证模式
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试 2: 客户端凭证模式")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	tokenURL := authServerURL + "/oauth2/v1/clientCredentials/token"
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	type TokenResponse struct {
		Code int `json:"code"`
		Data struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	var tokenResp TokenResponse
	json.Unmarshal(body, &tokenResp)

	if tokenResp.Code == 0 {
		fmt.Println("✅ Token 获取成功:")
		fmt.Printf("   Access Token: %s\n", tokenResp.Data.AccessToken)
		fmt.Printf("   Token Type: %s\n", tokenResp.Data.TokenType)
		fmt.Printf("   Expires In: %d 秒\n", tokenResp.Data.ExpiresIn)
		
		// 测试 3: 使用 Token 调用 API
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("测试 3: 使用 Token 调用用户信息 API")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		apiURL := authServerURL + "/openApi/v1/u/user/getCurrentUserInfo"
		req2, _ := http.NewRequest("POST", apiURL, nil)
		req2.Header.Set("Authorization", "Bearer "+tokenResp.Data.AccessToken)
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Printf("❌ API 调用失败: %v\n", err)
		} else {
			defer resp2.Body.Close()
			body2, _ := io.ReadAll(resp2.Body)

			var apiResp struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}
			json.Unmarshal(body2, &apiResp)

			if apiResp.Code == 0 {
				fmt.Println("✅ API 调用成功")
				fmt.Printf("   响应: %s\n", string(body2))
			} else {
				fmt.Printf("❌ API 调用失败: %s\n", apiResp.Msg)
				fmt.Println()
				fmt.Println("💡 这是预期行为！")
				fmt.Println("   客户端凭证模式的 Token 没有用户上下文，")
				fmt.Println("   无法调用用户信息接口。")
				fmt.Println()
				fmt.Println("   授权码登录流程应该使用授权码模式的 Token。")
			}
		}
	} else {
		fmt.Printf("❌ Token 获取失败: %s\n", tokenResp.Msg)
	}

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("诊断完成")
	fmt.Println("====================================")
	fmt.Println()
	fmt.Println("📋 总结:")
	fmt.Println("1. 客户端凭证模式用于服务间调用（无用户上下文）")
	fmt.Println("2. 授权码模式用于用户登录（包含用户信息）")
	fmt.Println("3. 如果登录时报 'need to Retry'，说明:")
	fmt.Println("   - 授权码换取的 Token 获取失败")
	fmt.Println("   - 或 Token 格式不正确")
	fmt.Println("   - 或主平台 API 不可达")
	fmt.Println()
	fmt.Println("🔍 下一步排查:")
	fmt.Println("1. 在实际登录时查看详细日志")
	fmt.Println("2. 检查授权码是否正确获取")
	fmt.Println("3. 检查 code 换取 token 的响应")
}
