package license

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Generator License 生成器
type Generator struct {
	privateKey *rsa.PrivateKey
	config     *Config
}

// NewGenerator 创建 License 生成器
func NewGenerator(config *Config) (*Generator, error) {
	privateKey, err := ParsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}
	return &Generator{
		privateKey: privateKey,
		config:     config,
	}, nil
}

// Generate 生成 License
func (g *Generator) Generate(options LicenseOptions) (*License, error) {
	// 生成唯一 ID
	if options.LicenseID == "" {
		options.LicenseID = uuid.New().String()
	}

	// 设置默认值
	if options.Product == "" {
		options.Product = "sun-panel"
	}
	if options.Version == "" {
		options.Version = "1.0.0"
	}
	if options.Type == "" {
		options.Type = "standard"
	}
	if options.Duration <= 0 {
		options.Duration = 365 // 默认 1 年
	}

	now := time.Now()
	expiresAt := now.AddDate(0, 0, options.Duration)

	license := &License{
		LicenseID: options.LicenseID,
		Product:   options.Product,
		Version:   options.Version,
		IssuedTo:  options.IssuedTo,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
		Features:  options.Features,
		MachineID: options.MachineID,
		MaxUsers:  options.MaxUsers,
		MaxNodes:  options.MaxNodes,
		Extra:     options.Extra,
		Type:      options.Type,
		Status:    "active",
	}

	// 签名
	if err := SignLicense(g.privateKey, license); err != nil {
		return nil, fmt.Errorf("sign license failed: %w", err)
	}

	return license, nil
}

// GenerateTrial 生成试用 License
func (g *Generator) GenerateTrial(issuedTo string, duration int) (*License, error) {
	return g.Generate(LicenseOptions{
		IssuedTo: issuedTo,
		Duration: duration,
		Features: []string{"basic"},
		Type:     "trial",
	})
}

// GenerateStandard 生成标准版 License
func (g *Generator) GenerateStandard(issuedTo string, duration int, machineID string) (*License, error) {
	return g.Generate(LicenseOptions{
		IssuedTo:  issuedTo,
		Duration:  duration,
		Features:  []string{"basic", "standard"},
		MachineID: machineID,
		Type:      "standard",
	})
}

// GenerateProfessional 生成专业版 License
func (g *Generator) GenerateProfessional(issuedTo string, duration int, machineID string, maxUsers int) (*License, error) {
	return g.Generate(LicenseOptions{
		IssuedTo:  issuedTo,
		Duration:  duration,
		Features:  []string{"basic", "standard", "professional"},
		MachineID: machineID,
		MaxUsers:  maxUsers,
		Type:      "professional",
	})
}

// GenerateEnterprise 生成企业版 License
func (g *Generator) GenerateEnterprise(issuedTo string, duration int, machineID string, maxUsers int, maxNodes int) (*License, error) {
	return g.Generate(LicenseOptions{
		IssuedTo:  issuedTo,
		Duration:  duration,
		Features:  []string{"basic", "standard", "professional", "enterprise"},
		MachineID: machineID,
		MaxUsers:  maxUsers,
		MaxNodes:  maxNodes,
		Type:      "enterprise",
	})
}

// Encode 将 License 编码为字符串
func (g *Generator) Encode(license *License) (string, error) {
	data, err := json.Marshal(license)
	if err != nil {
		return "", fmt.Errorf("marshal license failed: %w", err)
	}
	return EncodeLicense(data), nil
}

// GenerateAndEncode 生成并编码 License
func (g *Generator) GenerateAndEncode(options LicenseOptions) (string, *License, error) {
	license, err := g.Generate(options)
	if err != nil {
		return "", nil, err
	}
	encoded, err := g.Encode(license)
	if err != nil {
		return "", nil, err
	}
	return encoded, license, nil
}

// UpdateLicense 更新 License（延长有效期、更新功能等）
func (g *Generator) UpdateLicense(oldLicense *License, options LicenseOptions) (*License, error) {
	// 保留原有的 ID 和颁发时间
	if options.LicenseID == "" {
		options.LicenseID = oldLicense.LicenseID
	}

	// 生成新的 License
	newLicense, err := g.Generate(options)
	if err != nil {
		return nil, err
	}

	// 保持原有的颁发时间
	newLicense.IssuedAt = oldLicense.IssuedAt
	newLicense.Status = "active"

	// 重新签名
	if err := SignLicense(g.privateKey, newLicense); err != nil {
		return nil, fmt.Errorf("re-sign license failed: %w", err)
	}

	return newLicense, nil
}

// ExtendLicense 延长 License 有效期
func (g *Generator) ExtendLicense(oldLicense *License, additionalDays int) (*License, error) {
	options := LicenseOptions{
		LicenseID: oldLicense.LicenseID,
		Product:   oldLicense.Product,
		Version:   oldLicense.Version,
		IssuedTo:  oldLicense.IssuedTo,
		Duration:  int(time.Until(oldLicense.ExpiresAt).Hours()/24) + additionalDays,
		Features:  oldLicense.Features,
		MachineID: oldLicense.MachineID,
		MaxUsers:  oldLicense.MaxUsers,
		MaxNodes:  oldLicense.MaxNodes,
		Extra:     oldLicense.Extra,
		Type:      oldLicense.Type,
	}
	return g.UpdateLicense(oldLicense, options)
}

// RevokeLicense 撤销 License（生成一个已撤销的记录）
func (g *Generator) RevokeLicense(licenseID string, publicKeyPEM string) (*LicenseStatus, error) {
	return &LicenseStatus{
		LicenseID: licenseID,
		Status:    "revoked",
	}, nil
}
