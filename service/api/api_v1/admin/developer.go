package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type DeveloperApi struct {
}

// GetList 获取开发者列表
func (a *DeveloperApi) GetList(c *gin.Context) {
	param := DeveloperGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	list, total, err := biz.Developer.GetDeveloperList(global.Db, param.Page, param.Limit, param.Status, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取开发者详情
func (a *DeveloperApi) GetInfo(c *gin.Context) {
	param := DeveloperGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	info, err := biz.Developer.GetDeveloperInfo(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, info)
}

// GetByDeveloperName 根据开发者标识获取开发者信息
func (a *DeveloperApi) GetByDeveloperName(c *gin.Context) {
	param := DeveloperGetByDeveloperNameReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	developer, err := biz.Developer.GetByDeveloperName(global.Db, param.DeveloperName)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, developer)
}

// Update 更新开发者
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

	m := models.Developer{}
	err := m.UpdateInfo(global.Db, param.Id, param.DeveloperName, param.ContactMail, param.PaymentName, param.PaymentQrcode, param.PaymentMethod, param.Name)
	if err != nil {
		if err == gorm.ErrRegistered {
			apiReturn.Error(c, "开发者标识已存在")
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 更新状态（独立更新）
	if err := m.UpdateStatus(global.Db, param.Id, param.Status); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Deletes 删除开发者
func (a *DeveloperApi) Deletes(c *gin.Context) {
	param := DeveloperDeletesReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.Developer{}
	if err := m.Delete(global.Db, param.Ids); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// UpdateStatus 更新开发者状态
func (a *DeveloperApi) UpdateStatus(c *gin.Context) {
	param := DeveloperUpdateStatusReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.Developer{}
	if err := m.UpdateStatus(global.Db, param.Id, param.Status); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
