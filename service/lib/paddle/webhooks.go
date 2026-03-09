package paddle

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type WebhookType struct{}

const WEBHOOK_TYPE_TRANSACTION_COMPLETED = "transaction.completed"

type Notification struct {
	EventID        string                 `json:"event_id"`
	EventType      string                 `json:"event_type"`
	OccurredAt     string                 `json:"occurred_at"`
	NotificationID string                 `json:"notification_id"`
	Data           map[string]interface{} `json:"data"`
}

func (w *WebhookType) GinContextParse(c *gin.Context) (webhookType string, data interface{}, err error) {
	req := Notification{}

	if err = c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return
	}
	return req.EventType, req.Data, nil
}

func (w *WebhookType) TransactionCompleted(data interface{}) (resp TransactionData, err error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(dataBytes, &resp)
	if err != nil {
		return
	}
	return
}

// verifyWebhookSignature 验证Paddle Webhook的签名
func (w *WebhookType) VerifySignature(request *http.Request, secretKey string) bool {
	// 1. 获取Paddle-Signature头
	signatureHeader := request.Header.Get("Paddle-Signature")
	if signatureHeader == "" {
		return false
	}

	// 2. 从Paddle-Signature头中提取时间戳和签名
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
	// fmt.Println("ts,h1", timestamp, signature)

	// 3. 构建签名的负载（signed payload）
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return false
	}
	request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置请求体，以便后续可以再次读取
	signedPayload := fmt.Sprintf("%s:%s", timestamp, body)

	// 4. 使用HMAC和SHA256算法对负载进行哈希处理
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	// 5. 比较计算出的签名和Paddle-Signature头中的签名
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// func TestSpeed(function func()) {
// 	startTime := time.Now()
// 	function()
// 	endTime := time.Now()
// 	elapsedTime := endTime.Sub(startTime)
// 	fmt.Printf("代码块执行时间：%v\n", elapsedTime.Nanoseconds())
// }
