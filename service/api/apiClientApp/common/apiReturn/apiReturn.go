package apiReturn

import (
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"time"

	"github.com/gin-gonic/gin"
)

const ERROR_CODE_SUCCESS = 0 // 错误码 无任何错误

func ApiReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	// 从context获取responseId（由ResponseLogger中间件设置）
	responseId, _ := ctx.Get("responseId")
	if responseId == nil {
		// 如果中间件未设置，则自动生成
		responseId = time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	}

	// 标记业务错误，供ResponseLogger中间件使用
	if code != 0 {
		ctx.Set("hasError", true)
		ctx.Set("bizCode", code)
	}

	returnData := map[string]interface{}{
		"code":       code,
		"msg":        msg,
		"responseId": responseId,
	}
	if data != nil {
		returnData["data"] = data
	}
	ctx.JSON(200, returnData)
}

// 返回成功
func SuccessData(ctx *gin.Context, data interface{}) {
	ApiReturn(ctx, 0, "OK", data)
}

// 返回列表
// func SuccessListData(ctx *gin.Context, list interface{}, count int64) {
// 	ApiReturn(ctx, 0, "OK", gin.H{
// 		"list":  list,
// 		"count": count,
// 	})
// }

// 返回成功，没有data数据
func Success(ctx *gin.Context) {
	ApiReturn(ctx, 0, "OK", nil)
}

// 返回列表数据
func ListData(ctx *gin.Context, list interface{}, count int64) {
	data := map[string]interface{}{
		"list":  list,
		"count": count,
	}
	ApiReturn(ctx, 0, "OK", data)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorCode(ctx *gin.Context, code int, errMsg string, data interface{}) {
	ApiReturn(ctx, code, errMsg, data)
}

// 返回错误 普通提示错误
func Error(ctx *gin.Context, errMsg string) {
	ErrorCode(ctx, -1, errMsg, nil)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorNoAccess(ctx *gin.Context) {
	ErrorCode(ctx, 1005, global.Lang.Get("common.no_access"), nil)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	Error(ctx, global.Lang.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
	// Error(ctx, "参数错误")
}

// // 返回错误 数据库
func ErrorDatabase(ctx *gin.Context, errMsg string) {
	// Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
	ErrorByCodeAndMsg(ctx, 1200, errMsg)

}

// 返回错误 数据记录未找到
func ErrorDataNotFound(ctx *gin.Context) {
	// ErrorCode(ctx, ERROR_CODE_DATA_RECORD_NOT_FOUND, "未找到数据记录", nil)
	ErrorByCode(ctx, -1)
}

func ErrorClientNotFound(ctx *gin.Context) {
	ErrorByCode(ctx, 1016)
}

func ErrorByCode(ctx *gin.Context, code int) {
	msg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, msg, nil)
}

// 使用错误码的错误并附加错误信息
func ErrorByCodeAndMsg(ctx *gin.Context, code int, msg string) {
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, msg, nil)
}

func ErrorNotLogin(ctx *gin.Context) {
	// msg := "Server error"
	// if v, ok := GetErrorMsgByCode(1001); ok {
	// 	msg = v
	// }
	// ErrorCode(ctx, 1001, msg, nil)
	// 从context获取responseId
	responseId, _ := ctx.Get("responseId")
	if responseId == nil {
		responseId = time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	}

	// 标记错误
	ctx.Set("hasError", true)
	ctx.Set("bizCode", ErrNotLogged)

	returnData := map[string]interface{}{
		"code":       ErrNotLogged,
		"msg":        "Unauthorized",
		"responseId": responseId,
	}
	ctx.JSON(401, returnData)
}

// 返回401未授权错误 (HTTP状态码401)
func ErrorUnauthorized(ctx *gin.Context) {
	// 从context获取responseId
	responseId, _ := ctx.Get("responseId")
	if responseId == nil {
		responseId = time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	}

	// 标记错误
	ctx.Set("hasError", true)
	ctx.Set("bizCode", ErrCodeTokenInvalid)

	returnData := map[string]interface{}{
		"code":       ErrCodeTokenInvalid,
		"msg":        "Unauthorized",
		"responseId": responseId,
	}
	ctx.JSON(401, returnData)
}

// 返回401未授权错误 (HTTP状态码401)
func ErrorUnauthorizedAndBizCode(ctx *gin.Context, bizCode int, msg string) {
	// 从context获取responseId
	responseId, _ := ctx.Get("responseId")
	if responseId == nil {
		responseId = time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	}

	// 标记错误
	ctx.Set("hasError", true)
	ctx.Set("bizCode", bizCode)

	returnData := map[string]interface{}{
		"code":       bizCode,
		"msg":        msg,
		"responseId": responseId,
	}
	ctx.JSON(401, returnData)
}

// 过期
func ErrorTokenExpires(ctx *gin.Context) {
	ErrorUnauthorizedAndBizCode(ctx, ErrCodeTokenExpires, "Token expired")
}

func GetErrorMsgByCode(code int) (string, bool) {
	if v, ok := ErrorCodeMap[code]; ok {
		return v.Error(), true
	} else {
		return "", false
	}
}

// [废弃] 使用 ErrorWithNumReturnErrStr 结合 logger 替代，这样可以打印出具体行号
// func ErrorWithNum(ctx *gin.Context, err error) {
// 	errNum := time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")

// 	returnData := map[string]interface{}{
// 		"code":     -1,
// 		"msg":      "errorNum:" + errNum,
// 		"errorNum": errNum,
// 	}
// 	global.Logger.Error("ErrNum:"+errNum, "-", err)
// 	ctx.JSON(200, returnData)
// }

func ErrorWithNumReturnErrStr(ctx *gin.Context, err error) string {
	// 从context获取responseId（由ResponseLogger中间件设置）
	responseId, _ := ctx.Get("responseId")
	if responseId == nil {
		responseId = time.Now().Format("20060102150405") + cmn.BuildRandCode(5, "0123456789")
	}

	msg := "responseId:" + responseId.(string)

	// 标记错误
	ctx.Set("hasError", true)
	ctx.Set("bizCode", -1)

	returnData := map[string]interface{}{
		"code":       -1,
		"msg":        msg,
		"responseId": responseId,
	}
	ctx.JSON(200, returnData)
	return "RESPONSE-ID_" + responseId.(string) + ":" + err.Error()
}
