package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func VerifySignature(request *http.Request, secretKey string) bool {
	// 1. 获取Paddle-Signature头
	signatureHeader := request.Header.Get("SunStore-Signature")
	if signatureHeader == "" {
		return false
	}

	parts := strings.Split(signatureHeader, ";")
	if len(parts) != 2 {
		return false
	}

	timestamp, signature := parts[0], parts[1]
	if !strings.HasPrefix(timestamp, "ts=") {
		return false
	}

	if len(timestamp) <= 5 || len(signature) <= 5 {
		return false
	}

	timestamp = timestamp[3:]
	signature = signature[3:]

	body, err := io.ReadAll(request.Body)

	if err != nil {
		return false
	}
	request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置请求体，以便后续可以再次读取
	signedPayload := fmt.Sprintf("%s:%s", string(body), timestamp)

	calcSignature, err := GenerateHMACSHA256([]byte(signedPayload), []byte(secretKey))
	if err != nil {
		return false
	}

	return calcSignature == signature
}

func GenerateHMACSHA256(message, secret []byte) (string, error) {
	// 创建一个新的HMAC对象，使用sha256作为哈希函数
	hmac := hmac.New(sha256.New, secret)

	// 向HMAC写入要签名的消息
	_, err := hmac.Write(message)
	if err != nil {
		return "", err
	}

	// 计算签名
	signature := hmac.Sum(nil)
	signatureStr := hex.EncodeToString(signature)

	return signatureStr, nil
}
