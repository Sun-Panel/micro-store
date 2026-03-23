package biz

import (
	"errors"
	"slices"
)

// 业务错误码常量（只定义有业务含义的错误）
const (
	// 应用相关
	ErrCodeAppNotFound = "APP_NOT_FOUND"

	// 版本相关 - 前端需要针对性提示
	ErrCodeVersionNotFound      = "VERSION_NOT_FOUND"
	ErrCodeVersionExists        = "VERSION_EXISTS"        // 提示用户换个版本号
	ErrCodeVersionCodeExists    = "VERSION_CODE_EXISTS"   // 提示用户换个版本编号
	ErrCodeStatusNotAllowed     = "STATUS_NOT_ALLOWED"    // 提示当前状态不允许
	ErrCodeApprovedCannotDelete = "APPROVED_CANNOT_DELETE" // 已审核版本不能删除
	ErrCodeNotPendingReview     = "NOT_PENDING_REVIEW"    // 非待审核状态
	ErrCodeNoUpdateContent      = "NO_UPDATE_CONTENT"     // 无更新内容

	// 微应用开发者相关 3000-3099
	ErrCodeAppIdExists         = "APP_ID_EXISTS"          // 应用ID已存在
	ErrCodeNoPermission        = "NO_PERMISSION"          // 无权操作此应用
	ErrCodePendingReviewExists = "PENDING_REVIEW_EXISTS"  // 已有待审核记录
	ErrCodeNoPendingReviewApp  = "NO_PENDING_REVIEW_APP"  // 没有待审核记录（应用审核）
	ErrCodeInvalidParam        = "INVALID_PARAM"          // 参数无效
)

// NewBizError 创建业务错误
func NewBizError(code string) error {
	return errors.New(code)
}

// IsBizError 判断是否为业务错误
func IsBizError(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	errMsg := err.Error()

	validCodes := []string{
		ErrCodeAppNotFound, ErrCodeVersionNotFound,
		ErrCodeVersionExists, ErrCodeVersionCodeExists, ErrCodeStatusNotAllowed,
		ErrCodeApprovedCannotDelete, ErrCodeNotPendingReview, ErrCodeNoUpdateContent,
		// 微应用开发者相关
		ErrCodeAppIdExists, ErrCodeNoPermission, ErrCodePendingReviewExists, ErrCodeNoPendingReviewApp, ErrCodeInvalidParam,
	}

	if slices.Contains(validCodes, errMsg) {
		return errMsg, true
	}
	return "", false
}
