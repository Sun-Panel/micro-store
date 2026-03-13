package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type MicroAppCategoryApi struct {
}

// GetList 获取分类列表
func (a *MicroAppCategoryApi) GetList(c *gin.Context) {
	param := MicroAppCategoryGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppCategory{}
	list, total, err := m.GetList(global.Db, param.Page, param.Limit, param.Status, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取分类详情
func (a *MicroAppCategoryApi) GetInfo(c *gin.Context) {
	param := MicroAppCategoryGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppCategory{}
	info, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, info)
}

// Create 创建分类
func (a *MicroAppCategoryApi) Create(c *gin.Context) {
	param := MicroAppCategoryCreateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	m := models.MicroAppCategory{}
	id, err := m.CreateWithCheck(global.Db, param.Name, param.Icon, param.Sort, param.Status)
	if err != nil {
		if err == gorm.ErrRegistered {
			apiReturn.Error(c, "分类名称已存在")
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"id": id})
}

// Update 更新分类
func (a *MicroAppCategoryApi) Update(c *gin.Context) {
	param := MicroAppCategoryUpdateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	m := models.MicroAppCategory{}
	err := m.UpdateWithCheck(global.Db, param.Id, param.Name, param.Icon, param.Sort, param.Status)
	if err != nil {
		if err == gorm.ErrRegistered {
			apiReturn.Error(c, "分类名称已存在")
			return
		}
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Deletes 删除分类
func (a *MicroAppCategoryApi) Deletes(c *gin.Context) {
	param := MicroAppCategoryDeletesReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppCategory{}
	if err := m.Delete(global.Db, param.Ids); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// UpdateStatus 更新分类状态
func (a *MicroAppCategoryApi) UpdateStatus(c *gin.Context) {
	param := MicroAppCategoryUpdateStatusReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroAppCategory{}
	if err := m.UpdateStatus(global.Db, param.Id, param.Status); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// GetEnabledList 获取所有启用的分类（下拉选择用）
func (a *MicroAppCategoryApi) GetEnabledList(c *gin.Context) {
	m := models.MicroAppCategory{}
	list, err := m.GetEnabledList(global.Db)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}
