package license

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"sun-panel/lib/AES"
)

// EncodeLicense 编码 License 数据
// 使用 Base64 编码，可选 AES 加密
func EncodeLicense(data []byte) string {
	// Base64 编码
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeLicense 解码 License 数据
func DecodeLicense(encoded []byte) ([]byte, error) {
	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	return decoded, nil
}

// EncodeLicenseWithAES 使用 AES 加密后编码 License
func EncodeLicenseWithAES(data []byte, aesKey string) (string, error) {
	// AES 加密
	encrypted, err := AES.Encrypt(aesKey, string(data))
	if err != nil {
		return "", fmt.Errorf("aes encrypt failed: %w", err)
	}

	// Base64 编码
	return base64.StdEncoding.EncodeToString([]byte(encrypted)), nil
}

// DecodeLicenseWithAES 解码并解密 License
func DecodeLicenseWithAES(encoded []byte, aesKey string) ([]byte, error) {
	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	// AES 解密
	decrypted, err := AES.Decrypt(aesKey, string(decoded))
	if err != nil {
		return nil, fmt.Errorf("aes decrypt failed: %w", err)
	}

	return []byte(decrypted), nil
}

// EncodeLicenseToJSON 将 License 编码为 JSON 字符串
func EncodeLicenseToJSON(license *License) (string, error) {
	data, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return string(data), nil
}

// DecodeLicenseFromJSON 从 JSON 字符串解码 License
func DecodeLicenseFromJSON(jsonStr string) (*License, error) {
	var license License
	if err := json.Unmarshal([]byte(jsonStr), &license); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	return &license, nil
}

// EncodeLicenseCompact 紧凑编码（无缩进）
func EncodeLicenseCompact(license *License) (string, error) {
	data, err := json.Marshal(license)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return EncodeLicense(data), nil
}

// DecodeLicenseFromCompact 从紧凑编码解码 License
func DecodeLicenseFromCompact(encoded string) (*License, error) {
	decoded, err := DecodeLicense([]byte(encoded))
	if err != nil {
		return nil, err
	}

	var license License
	if err := json.Unmarshal(decoded, &license); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return &license, nil
}
