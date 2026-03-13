package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cnb.cool/hslr-s/go-pkg/oauth2-go/client"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("OAuth2 客户端测试")
	fmt.Println("====================================")
	fmt.Println()

	// 配置
	config := &client.Config{
		AuthServerURL: "http://192.168.3.101:3088",
		APIServerURL:  "http://192.168.3.101:3088",
		ClientID:      "xs1ubmt25p",
		ClientSecret:  "pKcHMOKBUOvwsOmokvn82JMTC2zR4mIy",
		Timeout:       30,
	}

	fmt.Println("📋 配置信息:")
	fmt.Printf("   授权服务器: %s\n", config.AuthServerURL)
	fmt.Printf("   API 服务器: %s\n", config.APIServerURL)
	fmt.Printf("   客户端 ID: %s\n", config.ClientID)
	fmt.Println()

	// 创建客户端
	fmt.Println("🔧 创建 OAuth2 客户端...")
	oauthClient, err := client.NewOAuth2Client(config)
	if err != nil {
		log.Fatal("❌ 创建客户端失败:", err)
	}
	fmt.Println("✅ 客户端创建成功")
	fmt.Println()

	// 测试 1: 客户端凭证模式
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试 1: 客户端凭证模式")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokenResp, err := oauthClient.GetClientCredentialsToken(ctx)
	if err != nil {
		fmt.Printf("❌ 获取 Token 失败: %v\n", err)
		fmt.Println()
		fmt.Println("💡 可能的原因:")
		fmt.Println("   1. 授权服务器不可达")
		fmt.Println("   2. 客户端凭证错误")
		fmt.Println("   3. 网络连接问题")
		fmt.Println()
		fmt.Println("🔍 请检查:")
		fmt.Println("   - 授权服务器是否运行: curl http://192.168.3.101:3088")
		fmt.Println("   - 客户端 ID 和 Secret 是否正确")
		return
	}

	fmt.Println("✅ Token 获取成功:")
	fmt.Printf("   Access Token: %s...\n", truncate(tokenResp.AccessToken, 30))
	fmt.Printf("   Token Type: %s\n", tokenResp.TokenType)
	fmt.Printf("   Expires In: %d 秒\n", tokenResp.ExpiresIn)
	if tokenResp.RefreshToken != "" {
		fmt.Printf("   Refresh Token: %s...\n", truncate(tokenResp.RefreshToken, 30))
	}
	if tokenResp.Scope != "" {
		fmt.Printf("   Scope: %s\n", tokenResp.Scope)
	}
	fmt.Println()

	// 测试 2: API 调用
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("测试 2: 调用 API 获取用户信息")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	apiClient := client.NewAPIClient(config.APIServerURL, 30)

	type UserInfo struct {
		Username string `json:"username"`
		Mail     string `json:"mail"`
		Name     string `json:"name"`
		Role     int    `json:"role"`
	}

	var userInfo UserInfo
	err = apiClient.Call(ctx, "POST", "/openApi/v1/u/user/getCurrentUserInfo",
		tokenResp.AccessToken, nil, &userInfo)

	if err != nil {
		fmt.Printf("❌ API 调用失败: %v\n", err)
		fmt.Println()
		fmt.Println("💡 可能的原因:")
		fmt.Println("   1. Access Token 无效或过期")
		fmt.Println("   2. API 端点不存在或路径错误")
		fmt.Println("   3. 权限不足")
		fmt.Println()
		fmt.Println("🔍 这就是导致 'need to Retry' 错误的原因！")
		fmt.Println()
		fmt.Println("🛠️ 解决方案:")
		fmt.Println("   1. 检查 API 端点路径是否正确")
		fmt.Println("   2. 检查 Token 是否需要特定前缀（如 'Bearer '）")
		fmt.Println("   3. 检查主平台的 Token 验证逻辑")
		return
	}

	fmt.Println("✅ API 调用成功:")
	fmt.Printf("   用户名: %s\n", userInfo.Username)
	fmt.Printf("   邮箱: %s\n", userInfo.Mail)
	fmt.Printf("   姓名: %s\n", userInfo.Name)
	fmt.Printf("   角色: %d\n", userInfo.Role)
	fmt.Println()

	fmt.Println("====================================")
	fmt.Println("✅ 所有测试通过！")
	fmt.Println("====================================")
	fmt.Println()
	fmt.Println("如果这里测试通过，但登录时仍报错，请检查:")
	fmt.Println("1. 授权码模式 (code) 是否正确获取")
	fmt.Println("2. 重定向 URI 是否匹配")
	fmt.Println("3. 前端发送的 code 是否有效")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
