package license

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	ErrLicenseNotFound   = errors.New("license not found")
	ErrLicenseExpired    = errors.New("license expired")
	ErrLicenseInvalid    = errors.New("license invalid")
	ErrMachineIDMismatch = errors.New("machine id mismatch")
	ErrSignatureInvalid  = errors.New("signature invalid")
	ErrLicenseRevoked    = errors.New("license revoked")
)

// Validator License 验证器
type Validator struct {
	publicKey *rsa.PublicKey
	config    *Config
	license   *License
	mu        sync.RWMutex
	lastCheck time.Time
}

// NewValidator 创建 License 验证器
func NewValidator(config *Config) (*Validator, error) {
	publicKey, err := ParsePublicKey(config.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("parse public key failed: %w", err)
	}
	return &Validator{
		publicKey: publicKey,
		config:    config,
	}, nil
}

// Init 初始化验证器（从文件加载 License）
func (v *Validator) Init() (*InitResult, error) {
	// 尝试从文件加载 License
	if v.config.LicenseFilePath != "" {
		if err := v.LoadFromFile(v.config.LicenseFilePath); err != nil {
			return &InitResult{
				Success: false,
				Message: fmt.Sprintf("load license failed: %v", err),
			}, err
		}
	}

	// 验证 License
	result := v.Validate()
	if !result.Valid {
		return &InitResult{
			Success: false,
			Message: result.Message,
		}, ErrLicenseInvalid
	}

	// 获取机器信息
	machineID, _ := GetMachineID(v.config.Salt)

	return &InitResult{
		Success: true,
		Message: "license initialized successfully",
		Info: map[string]string{
			"licenseId": v.license.LicenseID,
			"issuedTo":  v.license.IssuedTo,
			"expiresAt": v.license.ExpiresAt.Format("2006-01-02 15:04:05"),
			"type":      v.license.Type,
			"machineId": machineID,
			"features":  fmt.Sprintf("%v", v.license.Features),
		},
	}, nil
}

// LoadFromFile 从文件加载 License
func (v *Validator) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read license file failed: %w", err)
	}
	return v.LoadFromBytes(data)
}

// LoadFromBytes 从字节数组加载 License
func (v *Validator) LoadFromBytes(data []byte) error {
	// 解码
	decoded, err := DecodeLicense(data)
	if err != nil {
		return fmt.Errorf("decode license failed: %w", err)
	}

	// 解析 JSON
	var license License
	if err := json.Unmarshal(decoded, &license); err != nil {
		return fmt.Errorf("parse license failed: %w", err)
	}

	v.mu.Lock()
	v.license = &license
	v.mu.Unlock()

	return nil
}

// LoadFromString 从字符串加载 License
func (v *Validator) LoadFromString(licenseStr string) error {
	return v.LoadFromBytes([]byte(licenseStr))
}

// Validate 验证 License
func (v *Validator) Validate() *ValidateResult {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return &ValidateResult{
			Valid:   false,
			Message: "license not loaded",
		}
	}

	// 1. 验证签名
	if err := VerifyLicense(v.publicKey, v.license); err != nil {
		return &ValidateResult{
			Valid:   false,
			Message: fmt.Sprintf("signature verification failed: %v", err),
		}
	}

	// 2. 检查状态
	if v.license.Status == "revoked" {
		return &ValidateResult{
			Valid:   false,
			Message: "license has been revoked",
		}
	}

	// 3. 检查过期时间
	if time.Now().After(v.license.ExpiresAt) {
		return &ValidateResult{
			Valid:     false,
			Message:   "license expired",
			ExpiresAt: v.license.ExpiresAt,
			Remaining: 0,
		}
	}

	// 4. 验证机器码（如果配置了绑定）
	if v.config.MachineBind && v.license.MachineID != "" {
		matched, err := VerifyMachineID(v.license.MachineID, v.config.Salt)
		if err != nil {
			return &ValidateResult{
				Valid:   false,
				Message: fmt.Sprintf("verify machine id failed: %v", err),
			}
		}
		if !matched {
			return &ValidateResult{
				Valid:   false,
				Message: "machine id mismatch",
			}
		}
	}

	// 计算剩余天数
	remaining := int(time.Until(v.license.ExpiresAt).Hours() / 24)
	if remaining < 0 {
		remaining = 0
	}

	v.lastCheck = time.Now()

	return &ValidateResult{
		Valid:       true,
		Message:     "license is valid",
		LicenseInfo: v.license,
		ExpiresAt:   v.license.ExpiresAt,
		Remaining:   remaining,
	}
}

// GetInfo 获取 License 信息
func (v *Validator) GetInfo() *License {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.license
}

// HasFeature 检查是否拥有指定功能权限
func (v *Validator) HasFeature(feature string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return false
	}

	for _, f := range v.license.Features {
		if f == feature || f == "enterprise" {
			return true
		}
	}
	return false
}

// HasAnyFeature 检查是否拥有任意一个指定功能权限
func (v *Validator) HasAnyFeature(features ...string) bool {
	for _, feature := range features {
		if v.HasFeature(feature) {
			return true
		}
	}
	return false
}

// HasAllFeatures 检查是否拥有所有指定功能权限
func (v *Validator) HasAllFeatures(features ...string) bool {
	for _, feature := range features {
		if !v.HasFeature(feature) {
			return false
		}
	}
	return true
}

// IsExpired 检查是否过期
func (v *Validator) IsExpired() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return true
	}
	return time.Now().After(v.license.ExpiresAt)
}

// GetRemainingDays 获取剩余有效天数
func (v *Validator) GetRemainingDays() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return 0
	}
	remaining := int(time.Until(v.license.ExpiresAt).Hours() / 24)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetMaxUsers 获取最大用户数
func (v *Validator) GetMaxUsers() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return 0
	}
	return v.license.MaxUsers
}

// GetMaxNodes 获取最大节点数
func (v *Validator) GetMaxNodes() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return 0
	}
	return v.license.MaxNodes
}

// GetType 获取授权类型
func (v *Validator) GetType() string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return ""
	}
	return v.license.Type
}

// GetLicenseID 获取 License ID
func (v *Validator) GetLicenseID() string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return ""
	}
	return v.license.LicenseID
}

// SaveToFile 将 License 保存到文件
func (v *Validator) SaveToFile(filePath string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return ErrLicenseNotFound
	}

	data, err := json.MarshalIndent(v.license, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal license failed: %w", err)
	}

	encoded := EncodeLicense(data)
	return os.WriteFile(filePath, []byte(encoded), 0644)
}

// CheckExpiringSoon 检查是否即将过期（指定天数内）
func (v *Validator) CheckExpiringSoon(days int) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.license == nil {
		return true
	}

	remaining := int(time.Until(v.license.ExpiresAt).Hours() / 24)
	return remaining <= days
}

// GetLicense 获取当前 License（用于内部访问）
func (v *Validator) GetLicense() *License {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.license
}
