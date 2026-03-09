# License 授权验证库

企业级 License 授权验证库，提供完整的 License 生成、验证和管理功能。

## 功能特性

✅ **RSA 非对称加密签名验证**
- 使用 RSA 2048 位密钥对 License 进行数字签名
- 私钥签名，公钥验证，确保 License 不可篡改

✅ **多维度机器码绑定（支持 Docker 环境）**
- 组合容器 ID、MAC 地址、CPU 信息、主机名等多个特征
- 支持 Docker 容器环境，获取容器唯一标识
- 可自定义特征组合和盐值

✅ **时间防篡改检测**
- 通过 NTP 服务器校验真实时间
- 检测系统时间回拨攻击
- 累计运行时间验证

✅ **在线心跳验证 + 离线降级**
- 支持定期向授权服务器报告状态
- 离线环境下可运行一定时间（可配置）
- 防止同一 License 被多个实例同时使用

✅ **功能权限控制**
- 支持多种授权类型：试用版、标准版、专业版、企业版
- 细粒度的功能权限管理
- 用户数/节点数限制

✅ **Gin 框架中间件支持**
- 开箱即用的 Gin 中间件
- 功能权限验证中间件
- License 信息注入到请求上下文

## 快速开始

### 1. 生成密钥对（管理端）

```go
package main

import (
    "fmt"
    "os"
    "sun-panel/lib/license"
)

func main() {
    // 生成 RSA 密钥对
    privateKey, publicKey, err := license.GenerateKeyPair(2048)
    if err != nil {
        panic(err)
    }

    // 转换为 PEM 格式
    privateKeyPEM := license.PrivateKeyToPEM(privateKey)
    publicKeyPEM, _ := license.PublicKeyToPEM(publicKey)

    // 保存到文件（私钥务必保密！）
    os.WriteFile("private_key.pem", []byte(privateKeyPEM), 0600)
    os.WriteFile("public_key.pem", []byte(publicKeyPEM), 0644)

    fmt.Println("密钥对生成成功！")
    fmt.Println("private_key.pem - 私钥（管理端使用，请妥善保管）")
    fmt.Println("public_key.pem  - 公钥（应用端使用，可公开）")
}
```

### 2. 生成 License（管理端）

```go
package main

import (
    "fmt"
    "os"
    "sun-panel/lib/license"
)

func main() {
    // 读取私钥
    privateKeyPEM, _ := os.ReadFile("private_key.pem")

    // 创建生成器
    config := &license.Config{
        PrivateKey: string(privateKeyPEM),
    }
    generator, err := license.NewGenerator(config)
    if err != nil {
        panic(err)
    }

    // 获取机器码（可选，绑定机器）
    machineID, _ := license.GetMachineID("your-secret-salt")
    fmt.Printf("机器码: %s\n", machineID)

    // 生成 License
    licenseKey, lic, err := generator.GenerateAndEncode(license.LicenseOptions{
        IssuedTo:  "user@example.com",
        Duration:  365, // 365 天
        Features:  []string{"basic", "standard", "professional"},
        MachineID: machineID, // 可选，绑定机器
        MaxUsers:  100,       // 可选，最大用户数
        Type:      "professional",
    })
    if err != nil {
        panic(err)
    }

    // 保存 License
    os.WriteFile("license.lic", []byte(licenseKey), 0644)

    fmt.Printf("License 生成成功！\n")
    fmt.Printf("License ID: %s\n", lic.LicenseID)
    fmt.Printf("授权类型: %s\n", lic.Type)
    fmt.Printf("过期时间: %s\n", lic.ExpiresAt.Format("2006-01-02"))
}
```

### 3. 验证 License（应用端）

```go
package main

import (
    "fmt"
    "sun-panel/lib/license"
)

func main() {
    // 配置验证器
    config := &license.Config{
        PublicKey:       `-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----`,
        LicenseFilePath: "./license.lic",
        MachineBind:     true,               // 绑定机器码
        Salt:            "your-secret-salt", // 与生成时相同的盐值
    }

    // 创建验证器
    validator, err := license.NewValidator(config)
    if err != nil {
        panic(err)
    }

    // 初始化并验证
    result, err := validator.Init()
    if err != nil || !result.Success {
        fmt.Printf("License 验证失败: %s\n", result.Message)
        return
    }

    fmt.Println("License 验证成功！")
    fmt.Printf("License ID: %s\n", validator.GetLicenseID())
    fmt.Printf("授权类型: %s\n", validator.GetType())
    fmt.Printf("剩余天数: %d\n", validator.GetRemainingDays())

    // 检查功能权限
    if validator.HasFeature("professional") {
        fmt.Println("拥有专业版功能权限")
    }

    // 检查是否即将过期
    if validator.CheckExpiringSoon(30) {
        fmt.Println("警告：License 将在 30 天内过期！")
    }
}
```

### 4. 使用 Gin 中间件

```go
package main

import (
    "github.com/gin-gonic/gin"
    "sun-panel/lib/license"
)

func main() {
    // 创建验证器
    config := &license.Config{
        PublicKey:       `-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----`,
        LicenseFilePath: "./license.lic",
    }
    validator, _ := license.NewValidator(config)
    validator.Init()

    // 创建路由
    router := gin.Default()

    // 全局 License 验证中间件
    middleware := license.NewMiddleware(validator, config)
    router.Use(middleware.Handler())

    // 基本路由（所有授权用户可访问）
    router.GET("/api/basic", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "basic feature"})
    })

    // 专业版功能（需要专业版授权）
    router.GET("/api/professional",
        middleware.FeatureRequired("professional"),
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "professional feature"})
        },
    )

    // 企业版功能（需要企业版授权）
    router.GET("/api/enterprise",
        middleware.FeatureRequired("enterprise"),
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "enterprise feature"})
        },
    )

    // 获取 License 信息 API
    router.GET("/api/license/info", license.GetLicenseInfoHandler(validator))

    router.Run(":8080")
}
```

## Docker 环境防护

本库针对 Docker 环境提供了多层防护机制：

### 1. 多维度机器码

```go
// 使用默认特征生成机器码
machineID, _ := license.GetMachineID("your-salt")

// 使用自定义特征组合
machineID, _ := license.GetMachineIDWithFeatures(
    []string{"hostname", "container", "mac", "cpu"}, 
    "your-salt",
)
```

### 2. 心跳检测

```go
// 创建心跳管理器
heartbeat := license.NewHeartbeatManager(validator, &license.Config{
    HeartbeatURL:    "https://license-server.com/api/heartbeat",
    HeartbeatPeriod: 5 * time.Minute,
    OfflineLimit:    7, // 7 天离线期限
})

// 启动心跳
heartbeat.Start()
defer heartbeat.Stop()

// 检查心跳状态
result := heartbeat.CheckHeartbeat()
if !result.Success {
    fmt.Printf("心跳检测失败: %s\n", result.Message)
}
```

### 3. 时间防篡改

```go
// 创建时间验证器
timeValidator := license.NewTimeValidator(&license.Config{
    TimeDriftLimit: 300, // 允许 5 分钟偏差
})

// 验证时间
result := timeValidator.ValidateTime()
if !result.Valid {
    fmt.Printf("时间验证失败: %s\n", result.Message)
}
```

## 授权类型

### 试用版（Trial）
- 有效期：通常 30 天
- 功能：基本功能
- 限制：无机器绑定

### 标准版（Standard）
- 有效期：1 年或自定义
- 功能：基础 + 标准功能
- 限制：绑定机器码

### 专业版（Professional）
- 有效期：1 年或自定义
- 功能：基础 + 标准 + 专业功能
- 限制：绑定机器码 + 用户数限制

### 企业版（Enterprise）
- 有效期：1 年或自定义
- 功能：所有功能
- 限制：绑定机器码 + 用户数限制 + 节点数限制

## API 文档

### 主要类型

- `License`: 授权许可证结构
- `LicenseOptions`: 生成 License 的选项
- `ValidateResult`: 验证结果
- `Config`: 配置
- `Generator`: License 生成器
- `Validator`: License 验证器
- `HeartbeatManager`: 心跳管理器
- `TimeValidator`: 时间验证器
- `Middleware`: Gin 中间件

### 核心函数

```go
// 密钥管理
GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error)
PrivateKeyToPEM(privateKey *rsa.PrivateKey) string
PublicKeyToPEM(publicKey *rsa.PublicKey) (string, error)
ParsePrivateKey(pemKey string) (*rsa.PrivateKey, error)
ParsePublicKey(pemKey string) (*rsa.PublicKey, error)

// 生成器
NewGenerator(config *Config) (*Generator, error)
generator.Generate(options LicenseOptions) (*License, error)
generator.GenerateAndEncode(options LicenseOptions) (string, *License, error)

// 验证器
NewValidator(config *Config) (*Validator, error)
validator.Init() (*InitResult, error)
validator.Validate() *ValidateResult
validator.HasFeature(feature string) bool
validator.IsExpired() bool
validator.GetRemainingDays() int

// 机器码
GetMachineID(salt string) (string, error)
GetMachineInfo() (*MachineInfo, error)
VerifyMachineID(expectedMachineID string, salt string) (bool, error)

// 中间件
NewMiddleware(validator *Validator, config *Config, opts ...MiddlewareOption) *Middleware
middleware.Handler() gin.HandlerFunc
middleware.FeatureRequired(feature string) gin.HandlerFunc
```

## 文件结构

```
lib/license/
├── types.go              # 类型定义
├── crypto.go             # RSA 加密工具
├── machine_id.go         # 机器码获取（支持 Docker）
├── generator.go          # License 生成器
├── validator.go          # License 验证器
├── heartbeat.go          # 心跳检测
├── time_validator.go     # 时间防篡改验证
├── encode.go             # 编码工具
├── middleware.go         # Gin 中间件
├── keys.go               # 密钥管理
├── index.go              # 统一导出
├── doc.go                # 包文档
├── example_test.go       # 示例和测试
└── README.md             # 本文档
```

## 安全建议

1. **私钥管理**
   - 私钥必须妥善保管，切勿泄露
   - 不要将私钥打包到客户端应用
   - 建议使用密钥管理系统或硬件安全模块（HSM）

2. **公钥嵌入**
   - 公钥可以公开，建议编译时嵌入到应用中
   - 避免运行时从外部加载公钥

3. **盐值管理**
   - 使用强随机盐值
   - 盐值应该保密，不要硬编码在客户端代码中

4. **密钥长度**
   - 测试环境：2048 位
   - 生产环境：建议 4096 位

5. **定期检查**
   - 在关键功能处再次验证 License
   - 定期检查 License 状态
   - 及时提醒用户续费

## 测试

```bash
# 运行所有测试
go test -v ./lib/license/...

# 运行特定测试
go test -v ./lib/license/... -run TestLicenseGenerationAndValidation

# 运行性能测试
go test -bench=. ./lib/license/...
```

## License

MIT License

## 作者

Sun Panel Team
