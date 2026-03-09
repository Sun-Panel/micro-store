package apiResp

import "github.com/gin-gonic/gin"

// 业务错误码结构体
type ErrorCode struct {
	Code  int
	Error string
}

// 预定义业务错误码
var (
	ErrCodeSuccess = ErrorCode{0, "success"}

	// 通用错误,代表业务上否定的错误
	ErrCodeErrGeneric = ErrorCode{-1, "generic_error"}

	// 认证相关错误
	ErrCodeAuthFailed = ErrorCode{1001, "authentication_failed"}

	// 请求参数相关错误
	ErrCodeBadRequest = ErrorCode{2001, "bad_request"}

	// 用户相关错误
	ErrCodeErrUserNotFound = ErrorCode{3001, "user_not_found"}

	// 业务逻辑相关错误
	ErrCodeNoPurchaseQualification = ErrorCode{4001, "not_eligible_for_purchase"}
)

// 统一响应结构体
type Response struct {
	Code      int     `json:"code"`
	Error     string  `json:"error"`
	Data      any     `json:"data,omitempty"`
	Msg       *string `json:"detail,omitempty"`
	RequestID string  `json:"requestId"`
}

// 发送成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:  ErrCodeSuccess.Code,
		Error: ErrCodeSuccess.Error,
		Data:  data,
	})
}

// 发送错误响应
func ErrorResponse(c *gin.Context, errCode ErrorCode, detail string) {
	resp := Response{
		Code:  errCode.Code,
		Error: errCode.Error,
		Msg:   &detail,
	}

	c.JSON(200, resp) // API接口通常返回200，错误通过code体现
}

// 发送带数据的错误响应
func ErrorWithDataResponse(c *gin.Context, errCode ErrorCode, data any, detail string) {
	resp := Response{
		Code:  errCode.Code,
		Error: errCode.Error,
		Data:  data,
		Msg:   &detail,
	}

	c.JSON(200, resp)
}
