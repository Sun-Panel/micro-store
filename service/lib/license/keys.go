package license

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// Keys 密钥对
type Keys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// KeyManager 密钥管理器
type KeyManager struct {
	privateKeyPath string
	publicKeyPath  string
	keys           *Keys
}

// NewKeyManager 创建密钥管理器
func NewKeyManager(privateKeyPath, publicKeyPath string) *KeyManager {
	return &KeyManager{
		privateKeyPath: privateKeyPath,
		publicKeyPath:  publicKeyPath,
	}
}

// GenerateAndSave 生成并保存密钥对
func (km *KeyManager) GenerateAndSave(bits int) error {
	// 生成密钥对
	privateKey, publicKey, err := GenerateKeyPair(bits)
	if err != nil {
		return fmt.Errorf("generate key pair failed: %w", err)
	}

	// 保存私钥
	privateKeyPEM := PrivateKeyToPEM(privateKey)
	if err := km.saveKey(km.privateKeyPath, privateKeyPEM); err != nil {
		return fmt.Errorf("save private key failed: %w", err)
	}

	// 保存公钥
	publicKeyPEM, err := PublicKeyToPEM(publicKey)
	if err != nil {
		return fmt.Errorf("convert public key failed: %w", err)
	}
	if err := km.saveKey(km.publicKeyPath, publicKeyPEM); err != nil {
		return fmt.Errorf("save public key failed: %w", err)
	}

	km.keys = &Keys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return nil
}

// Load 加载密钥对
func (km *KeyManager) Load() error {
	privateKeyPEM, err := km.loadKey(km.privateKeyPath)
	if err != nil {
		return fmt.Errorf("load private key failed: %w", err)
	}

	publicKeyPEM, err := km.loadKey(km.publicKeyPath)
	if err != nil {
		return fmt.Errorf("load public key failed: %w", err)
	}

	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return fmt.Errorf("parse private key failed: %w", err)
	}

	publicKey, err := ParsePublicKey(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("parse public key failed: %w", err)
	}

	km.keys = &Keys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	return nil
}

// GetKeys 获取密钥对
func (km *KeyManager) GetKeys() *Keys {
	return km.keys
}

// GetPrivateKey 获取私钥
func (km *KeyManager) GetPrivateKey() *rsa.PrivateKey {
	if km.keys == nil {
		return nil
	}
	return km.keys.PrivateKey
}

// GetPublicKey 获取公钥
func (km *KeyManager) GetPublicKey() *rsa.PublicKey {
	if km.keys == nil {
		return nil
	}
	return km.keys.PublicKey
}

// saveKey 保存密钥到文件
func (km *KeyManager) saveKey(path string, keyPEM string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// 写入文件（权限 0600，仅所有者可读写）
	return os.WriteFile(path, []byte(keyPEM), 0600)
}

// loadKey 从文件加载密钥
func (km *KeyManager) loadKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadOrGenerate 加载或生成密钥对
func (km *KeyManager) LoadOrGenerate(bits int) error {
	// 尝试加载
	if err := km.Load(); err == nil {
		return nil
	}

	// 加载失败，生成新的
	return km.GenerateAndSave(bits)
}

// ExportPublicKeyPEM 导出公钥 PEM 格式
func (km *KeyManager) ExportPublicKeyPEM() (string, error) {
	if km.keys == nil || km.keys.PublicKey == nil {
		return "", fmt.Errorf("public key not loaded")
	}
	return PublicKeyToPEM(km.keys.PublicKey)
}

// ExportPrivateKeyPEM 导出私钥 PEM 格式
func (km *KeyManager) ExportPrivateKeyPEM() string {
	if km.keys == nil || km.keys.PrivateKey == nil {
		return ""
	}
	return PrivateKeyToPEM(km.keys.PrivateKey)
}

// ParseCertificate 从证书解析公钥
func ParseCertificate(certPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain RSA public key")
	}

	return publicKey, nil
}

// DefaultPrivateKeyPath 默认私钥路径
func DefaultPrivateKeyPath() string {
	return "./keys/private_key.pem"
}

// DefaultPublicKeyPath 默认公钥路径
func DefaultPublicKeyPath() string {
	return "./keys/public_key.pem"
}

// DefaultLicenseFilePath 默认 License 文件路径
func DefaultLicenseFilePath() string {
	return "./license.lic"
}

// EnsureKeysExist 确保密钥文件存在
func EnsureKeysExist(bits int) error {
	km := NewKeyManager(DefaultPrivateKeyPath(), DefaultPublicKeyPath())
	return km.LoadOrGenerate(bits)
}

// GetEmbeddedPublicKey 获取嵌入的公钥（用于编译时嵌入公钥）
// 使用方式：在编译时通过 go embed 嵌入公钥文件
func GetEmbeddedPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	return ParsePublicKey(publicKeyPEM)
}
