package panel

import (
	"os"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type DeveloperApi struct {
}

// getUserId 从上下文获取当前用户ID
func getUserId(c *gin.Context) (uint, bool) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return 0, false
	}
	if u, ok := userInfo.(models.User); ok {
		return u.ID, true
	}
	return 0, false
}

// Register 申请成为开发者
func (a *DeveloperApi) Register(c *gin.Context) {
	param := DeveloperRegisterReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// 获取当前登录用户ID
	userId, exists := getUserId(c)
	if !exists {
		apiReturn.ErrorByCode(c, 1000) // 未登录
		return
	}

	m := models.Developer{}
	id, err := m.Register(global.Db, userId, param.DeveloperName, param.ContactMail, param.PaymentName, param.PaymentQrcode, param.PaymentMethod, param.Name)
	if err != nil {
		if err == gorm.ErrRegistered {
			apiReturn.Error(c, "您已是开发者或开发者标识已存在")
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"id": id})
}

// GetInfo 获取当前用户的开发者信息
func (a *DeveloperApi) GetInfo(c *gin.Context) {
	userId, exists := getUserId(c)
	if !exists {
		apiReturn.ErrorByCode(c, 1000)
		return
	}

	m := models.Developer{}
	info, err := m.GetByUserId(global.Db, userId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, info)
}

// Update 更新当前用户的开发者信息
func (a *DeveloperApi) Update(c *gin.Context) {
	param := DeveloperUpdateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	userId, exists := getUserId(c)
	if !exists {
		apiReturn.ErrorByCode(c, 1000)
		return
	}

	// 获取当前用户的开发者信息
	m := models.Developer{}
	info, err := m.GetByUserId(global.Db, userId)
	if err != nil {
		apiReturn.Error(c, "您还不是开发者")
		return
	}

	if info.PaymentQrcode != param.PaymentQrcode {
		// 删除旧的二维码
		if info.PaymentQrcode != "" {
			path := "." + info.PaymentQrcode
			if err := os.Remove(path); err != nil {
				global.Logger.Error("删除旧的支付二维码文件失败", "path", path, "error", err)
			}
		}
	}

	// 更新信息
	err = m.UpdateInfo(global.Db, info.ID, param.DeveloperName, param.ContactMail, param.PaymentName, param.PaymentQrcode, param.PaymentMethod, param.Name)
	if err != nil {
		if err == gorm.ErrRegistered {
			apiReturn.Error(c, "开发者标识已存在")
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// CheckIsDeveloper 检查当前用户是否是开发者
func (a *DeveloperApi) CheckIsDeveloper(c *gin.Context) {
	userId, exists := getUserId(c)
	if !exists {
		apiReturn.ErrorByCode(c, 1000)
		return
	}

	m := models.Developer{}
	isDeveloper, err := m.CheckUserIsDeveloper(global.Db, userId)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"isDeveloper": isDeveloper})
}
