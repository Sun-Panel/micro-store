package license

import (
	"fmt"
	"testing"
	"time"
)

// Example_basicUsage 基本使用示例
func Example_basicUsage() {
	// 1. 生成密钥对
	privateKey, publicKey, err := GenerateKeyPair(2048)
	if err != nil {
		panic(err)
	}

	privateKeyPEM := PrivateKeyToPEM(privateKey)
	publicKeyPEM, _ := PublicKeyToPEM(publicKey)

	// 2. 创建生成器并生成 License
	genConfig := &Config{
		PrivateKey: privateKeyPEM,
	}
	generator, err := NewGenerator(genConfig)
	if err != nil {
		panic(err)
	}

	licenseKey, lic, err := generator.GenerateAndEncode(LicenseOptions{
		IssuedTo: "user@example.com",
		Duration: 365,
		Features: []string{"basic", "standard"},
		Type:     "standard",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("License ID: %s\n", lic.LicenseID)
	fmt.Printf("Expires at: %s\n", lic.ExpiresAt.Format("2006-01-02"))
	fmt.Printf("License key length: %d bytes\n", len(licenseKey))

	// 3. 创建验证器并验证
	valConfig := &Config{
		PublicKey: publicKeyPEM,
	}
	validator, err := NewValidator(valConfig)
	if err != nil {
		panic(err)
	}

	if err := validator.LoadFromString(licenseKey); err != nil {
		panic(err)
	}

	result := validator.Validate()
	fmt.Printf("Valid: %v\n", result.Valid)
	fmt.Printf("Remaining days: %d\n", result.Remaining)

	// 注意：剩余天数可能在 364-365 之间，取决于测试运行时间
	// Output:
	// License ID: <generated-uuid>
	// Expires at: <date>
	// License key length: <number>
	// Valid: true
	// Remaining days: <364-365>
}

// Example_machineID 机器码示例
func Example_machineID() {
	// 获取机器信息
	info, err := GetMachineInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Printf("Platform: %s\n", info.Platform)
	fmt.Printf("Container ID: %s\n", info.ContainerID)

	// 生成机器码
	machineID, err := GetMachineID("my-salt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Machine ID: %s\n", machineID)

	// 使用指定特征生成机器码
	machineIDCustom, err := GetMachineIDWithFeatures([]string{"hostname", "mac", "cpu"}, "my-salt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Custom Machine ID: %s\n", machineIDCustom)
}

// Example_differentLicenseTypes 不同类型 License 示例
func Example_differentLicenseTypes() {
	privateKey, _, _ := GenerateKeyPair(2048)
	privateKeyPEM := PrivateKeyToPEM(privateKey)

	config := &Config{PrivateKey: privateKeyPEM}
	generator, _ := NewGenerator(config)

	// 试用版（30天，基本功能）
	trial, _ := generator.GenerateTrial("trial@example.com", 30)
	fmt.Printf("Trial: %s, expires: %s\n", trial.Type, trial.ExpiresAt.Format("2006-01-02"))

	// 标准版（1年，标准功能，绑定机器）
	machineID, _ := GetMachineID("salt")
	standard, _ := generator.GenerateStandard("standard@example.com", 365, machineID)
	fmt.Printf("Standard: %s, machine: %s\n", standard.Type, standard.MachineID[:8])

	// 专业版（1年，专业功能，限制用户数）
	professional, _ := generator.GenerateProfessional("pro@example.com", 365, machineID, 100)
	fmt.Printf("Professional: %s, max users: %d\n", professional.Type, professional.MaxUsers)

	// 企业版（1年，所有功能，限制节点数）
	enterprise, _ := generator.GenerateEnterprise("enterprise@example.com", 365, machineID, 1000, 10)
	fmt.Printf("Enterprise: %s, max nodes: %d\n", enterprise.Type, enterprise.MaxNodes)
}

// Example_keyManagement 密钥管理示例
func Example_keyManagement() {
	// 使用 KeyManager 管理密钥
	km := NewKeyManager("./keys/private.pem", "./keys/public.pem")

	// 生成或加载密钥
	if err := km.LoadOrGenerate(2048); err != nil {
		panic(err)
	}

	// 导出公钥
	publicKeyPEM, _ := km.ExportPublicKeyPEM()
	fmt.Printf("Public key length: %d bytes\n", len(publicKeyPEM))

	// 使用密钥创建生成器
	config := &Config{
		PrivateKey: km.ExportPrivateKeyPEM(),
	}
	generator, _ := NewGenerator(config)

	// 生成 License...
	_ = generator
}

// TestLicenseGenerationAndValidation 完整流程测试
func TestLicenseGenerationAndValidation(t *testing.T) {
	// 1. 生成密钥对
	privateKey, publicKey, err := GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	privateKeyPEM := PrivateKeyToPEM(privateKey)
	publicKeyPEM, _ := PublicKeyToPEM(publicKey)

	// 2. 创建生成器
	genConfig := &Config{PrivateKey: privateKeyPEM}
	generator, err := NewGenerator(genConfig)
	if err != nil {
		t.Fatalf("NewGenerator failed: %v", err)
	}

	// 3. 生成 License
	licenseKey, lic, err := generator.GenerateAndEncode(LicenseOptions{
		IssuedTo: "test@example.com",
		Duration: 365,
		Features: []string{"basic", "standard", "professional"},
		Type:     "professional",
		MaxUsers: 100,
	})
	if err != nil {
		t.Fatalf("GenerateAndEncode failed: %v", err)
	}

	// 验证生成的 License
	if lic.LicenseID == "" {
		t.Error("LicenseID should not be empty")
	}
	if lic.Type != "professional" {
		t.Errorf("Expected type 'professional', got '%s'", lic.Type)
	}
	if len(lic.Features) != 3 {
		t.Errorf("Expected 3 features, got %d", len(lic.Features))
	}
	if lic.MaxUsers != 100 {
		t.Errorf("Expected max users 100, got %d", lic.MaxUsers)
	}

	// 4. 创建验证器
	valConfig := &Config{PublicKey: publicKeyPEM}
	validator, err := NewValidator(valConfig)
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	// 5. 加载 License
	if err := validator.LoadFromString(licenseKey); err != nil {
		t.Fatalf("LoadFromString failed: %v", err)
	}

	// 6. 验证
	result := validator.Validate()
	if !result.Valid {
		t.Errorf("License should be valid, but got: %s", result.Message)
	}
	// 允许剩余天数在 364-365 之间（测试运行时间可能导致差一天）
	if result.Remaining < 364 || result.Remaining > 365 {
		t.Errorf("Expected remaining days 364-365, got %d", result.Remaining)
	}

	// 7. 检查功能权限
	if !validator.HasFeature("professional") {
		t.Error("Should have 'professional' feature")
	}
	if !validator.HasFeature("basic") {
		t.Error("Should have 'basic' feature")
	}
	if validator.HasFeature("enterprise") {
		t.Error("Should not have 'enterprise' feature")
	}

	// 8. 检查过期状态
	if validator.IsExpired() {
		t.Error("License should not be expired")
	}
	if validator.CheckExpiringSoon(30) {
		t.Error("License should not be expiring soon (within 30 days)")
	}
}

// TestMachineIDGeneration 机器码生成测试
func TestMachineIDGeneration(t *testing.T) {
	info, err := GetMachineInfo()
	if err != nil {
		t.Fatalf("GetMachineInfo failed: %v", err)
	}

	if info.Hostname == "" {
		t.Error("Hostname should not be empty")
	}

	machineID, err := GetMachineID("test-salt")
	if err != nil {
		t.Fatalf("GetMachineID failed: %v", err)
	}

	if len(machineID) != 64 { // SHA256 produces 64 hex characters
		t.Errorf("Expected machine ID length 64, got %d", len(machineID))
	}

	// 使用相同参数应该生成相同的机器码
	machineID2, _ := GetMachineID("test-salt")
	if machineID != machineID2 {
		t.Error("Same parameters should produce same machine ID")
	}

	// 使用不同盐值应该生成不同的机器码
	machineID3, _ := GetMachineID("different-salt")
	if machineID == machineID3 {
		t.Error("Different salt should produce different machine ID")
	}
}

// TestSignatureVerification 签名验证测试
func TestSignatureVerification(t *testing.T) {
	privateKey, publicKey, _ := GenerateKeyPair(2048)

	// 测试签名和验证
	data := []byte("test data for signing")
	signature, err := Sign(privateKey, data)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if err := Verify(publicKey, data, signature); err != nil {
		t.Errorf("Verify failed: %v", err)
	}

	// 使用错误的数据验证应该失败
	wrongData := []byte("wrong data")
	if err := Verify(publicKey, wrongData, signature); err == nil {
		t.Error("Verify should fail with wrong data")
	}
}

// TestTimeValidation 时间验证测试
func TestTimeValidation(t *testing.T) {
	config := &Config{
		TimeDriftLimit: 300, // 5 分钟
	}
	tv := NewTimeValidator(config)

	result := tv.ValidateTime()
	// 在正常环境下，时间验证应该通过
	// 注意：这个测试可能会因为网络问题失败，所以不强制要求
	t.Logf("Time validation result: Valid=%v, Message=%s, Drift=%d",
		result.Valid, result.Message, result.Drift)
}

// BenchmarkLicenseGeneration License 生成性能测试
func BenchmarkLicenseGeneration(b *testing.B) {
	privateKey, _, _ := GenerateKeyPair(2048)
	privateKeyPEM := PrivateKeyToPEM(privateKey)

	config := &Config{PrivateKey: privateKeyPEM}
	generator, _ := NewGenerator(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(LicenseOptions{
			IssuedTo: "benchmark@example.com",
			Duration: 365,
			Features: []string{"basic"},
			Type:     "standard",
		})
	}
}

// BenchmarkLicenseValidation License 验证性能测试
func BenchmarkLicenseValidation(b *testing.B) {
	privateKey, publicKey, _ := GenerateKeyPair(2048)
	privateKeyPEM := PrivateKeyToPEM(privateKey)
	publicKeyPEM, _ := PublicKeyToPEM(publicKey)

	genConfig := &Config{PrivateKey: privateKeyPEM}
	generator, _ := NewGenerator(genConfig)
	licenseKey, _, _ := generator.GenerateAndEncode(LicenseOptions{
		IssuedTo: "benchmark@example.com",
		Duration: 365,
		Features: []string{"basic"},
		Type:     "standard",
	})

	valConfig := &Config{PublicKey: publicKeyPEM}
	validator, _ := NewValidator(valConfig)
	validator.LoadFromString(licenseKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate()
	}
}

// BenchmarkMachineIDGeneration 机器码生成性能测试
func BenchmarkMachineIDGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetMachineID("benchmark-salt")
	}
}

// BenchmarkSignatureVerification 签名验证性能测试
func BenchmarkSignatureVerification(b *testing.B) {
	privateKey, publicKey, _ := GenerateKeyPair(2048)
	data := []byte("test data for benchmark")
	signature, _ := Sign(privateKey, data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Verify(publicKey, data, signature)
	}
}

// TestExpiredLicense 过期 License 测试
func TestExpiredLicense(t *testing.T) {
	privateKey, publicKey, _ := GenerateKeyPair(2048)
	privateKeyPEM := PrivateKeyToPEM(privateKey)
	publicKeyPEM, _ := PublicKeyToPEM(publicKey)

	// 创建一个已过期的 License
	genConfig := &Config{PrivateKey: privateKeyPEM}
	generator, _ := NewGenerator(genConfig)

	license := &License{
		LicenseID: "test-expired",
		Product:   "test",
		Version:   "1.0",
		IssuedTo:  "test@example.com",
		IssuedAt:  time.Now().AddDate(-1, 0, 0),
		ExpiresAt: time.Now().AddDate(0, 0, -1), // 昨天过期
		Features:  []string{"basic"},
		Type:      "standard",
		Status:    "active",
	}

	SignLicense(generator.privateKey, license)

	// 验证过期 License
	valConfig := &Config{PublicKey: publicKeyPEM}
	validator, _ := NewValidator(valConfig)

	licenseData, _ := EncodeLicenseCompact(license)
	validator.LoadFromString(licenseData)

	result := validator.Validate()
	if result.Valid {
		t.Error("Expired license should be invalid")
	}
	if result.Message != "license expired" {
		t.Errorf("Expected 'license expired' message, got: %s", result.Message)
	}
	if !validator.IsExpired() {
		t.Error("IsExpired should return true")
	}
}
