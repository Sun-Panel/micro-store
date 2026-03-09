// Package license 提供企业级 License 授权验证功能
//
// # License 授权验证库
//
// 本库提供完整的 License 生成、验证和管理功能，适用于 Go 后端应用。
//
// # 功能特性
//
//   - RSA 非对称加密签名验证
//   - 多维度机器码绑定（支持 Docker 环境）
//   - 时间防篡改检测（NTP 时间校验）
//   - 在线心跳验证 + 离线降级模式
//   - 功能权限控制
//   - 授权类型管理（试用版、标准版、专业版、企业版）
//   - Gin 框架中间件支持
//   - 用户数/节点数限制
//
// # 安全措施
//
//   - RSA 2048 位密钥加密
//   - 数字签名防篡改
//   - 机器码绑定防复制
//   - 时间戳防重放攻击
//   - 在线验证防滥用
//   - 离线有效期限制
//
// # 快速开始
//
// ## 1. 生成密钥对（管理端）
//
//	privateKey, publicKey, err := license.GenerateKeyPair(2048)
//	if err != nil {
//	    panic(err)
//	}
//
//	privateKeyPEM := license.PrivateKeyToPEM(privateKey)
//	publicKeyPEM, _ := license.PublicKeyToPEM(publicKey)
//
//	// 保存密钥到文件
//	os.WriteFile("private_key.pem", []byte(privateKeyPEM), 0600)
//	os.WriteFile("public_key.pem", []byte(publicKeyPEM), 0644)
//
// ## 2. 生成 License（管理端）
//
//	config := &license.Config{
//	    PrivateKey: privateKeyPEM,
//	}
//
//	generator, err := license.NewGenerator(config)
//	if err != nil {
//	    panic(err)
//	}
//
//	// 获取机器码
//	machineID, _ := license.GetMachineID("your-salt")
//
//	// 生成 License
//	licenseKey, lic, err := generator.GenerateAndEncode(license.LicenseOptions{
//	    IssuedTo:  "user@example.com",
//	    Duration:  365, // 365 天
//	    Features:  []string{"basic", "standard", "professional"},
//	    MachineID: machineID,
//	    MaxUsers:  100,
//	    Type:      "professional",
//	})
//
//	// 保存 License 到文件
//	os.WriteFile("license.lic", []byte(licenseKey), 0644)
//
// ## 3. 验证 License（应用端）
//
//	config := &license.Config{
//	    PublicKey:       publicKeyPEM,
//	    LicenseFilePath: "./license.lic",
//	    MachineBind:     true,
//	    Salt:            "your-salt",
//	}
//
//	validator, err := license.NewValidator(config)
//	if err != nil {
//	    panic(err)
//	}
//
//	// 初始化并验证
//	result, err := validator.Init()
//	if err != nil || !result.Success {
//	    panic("License validation failed: " + result.Message)
//	}
//
//	// 检查功能权限
//	if validator.HasFeature("professional") {
//	    // 允许使用专业版功能
//	}
//
//	// 检查是否即将过期
//	if validator.CheckExpiringSoon(30) {
//	    // 提示用户续费
//	}
//
// ## 4. 使用 Gin 中间件
//
//	router := gin.Default()
//
//	// 全局 License 验证
//	middleware := license.NewMiddleware(validator, config)
//	router.Use(middleware.Handler())
//
//	// 功能权限验证
//	router.GET("/api/professional",
//	    middleware.FeatureRequired("professional"),
//	    func(c *gin.Context) {
//	        c.JSON(200, gin.H{"message": "professional feature"})
//	    },
//	)
//
// # Docker 环境防护
//
// 本库针对 Docker 环境提供了多层防护机制：
//
// 1. 多维度机器码：组合容器 ID、MAC 地址、CPU 信息等多个特征
// 2. 服务实例绑定：首次启动时生成实例 ID 并持久化
// 3. 心跳检测：定期向授权服务器报告状态，防止多实例滥用
// 4. 时间防篡改：通过 NTP 服务器校验时间，防止修改系统时间绕过过期检测
// 5. 离线有效期：离线环境下仍可运行一定时间（如 7 天），超过限制则拒绝服务
//
// # 文件结构
//
//	lib/license/
//	├── types.go           # 类型定义
//	├── crypto.go          # RSA 加密工具
//	├── machine_id.go      # 机器码获取（支持 Docker）
//	├── generator.go       # License 生成器
//	├── validator.go       # License 验证器
//	├── heartbeat.go       # 心跳检测
//	├── time_validator.go  # 时间防篡改验证
//	├── encode.go          # 编码工具
//	├── middleware.go      # Gin 中间件
//	├── keys.go            # 密钥管理
//	├── index.go           # 统一导出
//	└── doc.go             # 文档
//
// # 注意事项
//
// 1. 私钥必须妥善保管，切勿泄露或打包到客户端应用
// 2. 公钥可以公开，建议编译时嵌入到应用中
// 3. 生产环境建议使用更长的密钥（如 4096 位）
// 4. 定期检查 License 状态，及时提醒用户续费
// 5. 对于重要功能，建议在关键代码处再次验证 License
//
// # License
//
// MIT License
package license
