// Package license 提供企业级 License 授权验证功能
//
// License 授权验证库
// 提供完整的 License 生成、验证和管理功能
//
// 功能特性：
//   - RSA 非对称加密签名验证
//   - 多维度机器码绑定（支持 Docker 环境）
//   - 时间防篡改检测
//   - 在线心跳验证 + 离线降级
//   - 功能权限控制
//   - Gin 框架中间件支持
//
// 快速开始：
//
//	// 1. 生成密钥对（管理端）
//	privateKey, publicKey, _ := license.GenerateKeyPair(2048)
//	privateKeyPEM := license.PrivateKeyToPEM(privateKey)
//	publicKeyPEM, _ := license.PublicKeyToPEM(publicKey)
//
//	// 2. 创建 License 生成器
//	config := &license.Config{
//	    PrivateKey: privateKeyPEM,
//	}
//	generator, _ := license.NewGenerator(config)
//
//	// 3. 生成 License
//	licenseKey, lic, _ := generator.GenerateAndEncode(license.LicenseOptions{
//	    IssuedTo:  "user@example.com",
//	    Duration:  365,
//	    Features:  []string{"basic", "standard", "professional"},
//	    Type:      "professional",
//	})
//
//	// 4. 创建 License 验证器（应用端）
//	validatorConfig := &license.Config{
//	    PublicKey:       publicKeyPEM,
//	    LicenseFilePath: "./license.lic",
//	}
//	validator, _ := license.NewValidator(validatorConfig)
//
//	// 5. 验证 License
//	result := validator.Validate()
//	if result.Valid {
//	    fmt.Println("License is valid!")
//	}
//
//	// 6. 使用中间件（Gin 框架）
//	middleware := license.NewMiddleware(validator, validatorConfig)
//	router.Use(middleware.Handler())
//
// 主要类型和函数说明：
//
// 类型：
//   - License: 授权许可证结构
//   - LicenseOptions: 生成 License 的选项
//   - ValidateResult: 验证结果
//   - Config: 验证器配置
//   - Generator: License 生成器
//   - Validator: License 验证器
//   - HeartbeatManager: 心跳管理器
//   - TimeValidator: 时间验证器
//   - Middleware: Gin 中间件
//
// 核心函数：
//   - GenerateKeyPair(): 生成 RSA 密钥对
//   - NewGenerator(): 创建 License 生成器
//   - NewValidator(): 创建 License 验证器
//   - GetMachineID(): 获取机器码
//   - SignLicense(): 签名 License
//   - VerifyLicense(): 验证 License 签名
//
// 中间件函数：
//   - RequireLicense(): License 验证中间件
//   - RequireFeature(): 功能权限验证中间件
//   - GetLicenseFromContext(): 从上下文获取 License 信息
package license
