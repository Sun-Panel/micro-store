package license

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrInvalidKey       = errors.New("invalid key")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrSignFailed       = errors.New("sign failed")
	ErrVerifyFailed     = errors.New("verify failed")
)

// GenerateKeyPair 生成 RSA 密钥对
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key pair failed: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// PrivateKeyToPEM 将私钥转换为 PEM 格式
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// PublicKeyToPEM 将公钥转换为 PEM 格式
func PublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("marshal public key failed: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}

// ParsePrivateKey 从 PEM 格式解析私钥
func ParsePrivateKey(pemKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, ErrInvalidKey
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}
	return privateKey, nil
}

// ParsePublicKey 从 PEM 格式解析公钥
func ParsePublicKey(pemKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, ErrInvalidKey
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key failed: %w", err)
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	return publicKey, nil
}

// Sign 使用私钥对数据进行签名
func Sign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrSignFailed, err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify 使用公钥验证签名
func Verify(publicKey *rsa.PublicKey, data []byte, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature failed: %w", err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrVerifyFailed, err)
	}
	return nil
}

// SignLicense 对 License 进行签名
func SignLicense(privateKey *rsa.PrivateKey, license *License) error {
	// 创建签名数据（不包含签名字段）
	signData := licenseToSignData(license)
	signature, err := Sign(privateKey, signData)
	if err != nil {
		return err
	}
	license.Signature = signature
	return nil
}

// VerifyLicense 验证 License 签名
func VerifyLicense(publicKey *rsa.PublicKey, license *License) error {
	signature := license.Signature
	license.Signature = "" // 临时清空签名
	signData := licenseToSignData(license)
	license.Signature = signature // 恢复签名

	return Verify(publicKey, signData, signature)
}

// licenseToSignData 将 License 转换为待签名数据
func licenseToSignData(license *License) []byte {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%v|%s|%d|%d|%s|%s",
		license.LicenseID,
		license.Product,
		license.Version,
		license.IssuedTo,
		license.IssuedAt.Format("2006-01-02 15:04:05"),
		license.ExpiresAt.Format("2006-01-02 15:04:05"),
		license.Features,
		license.MachineID,
		license.MaxUsers,
		license.MaxNodes,
		license.Extra,
		license.Type,
	)
	return []byte(data)
}
