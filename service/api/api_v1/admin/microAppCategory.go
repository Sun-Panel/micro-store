package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MicroAppCategoryApi struct {
}

// 获取分类列表
func (a *MicroAppCategoryApi) GetList(c *gin.Context) {
	type ParamsStruct struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		Status  *int   `json:"status"`
		KeyWord string `json:"keyWord"`
	}

	param := ParamsStruct{}
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

// 获取分类详情
func (a *MicroAppCategoryApi) GetInfo(c *gin.Context) {
	type ParamsStruct struct {
		Id uint `json:"id" binding:"required"`
	}

	param := ParamsStruct{}
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

// 创建分类
func (a *MicroAppCategoryApi) Create(c *gin.Context) {
	type ParamsStruct struct {
		Name   string `json:"name" binding:"required"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	param := ParamsStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	m := models.MicroAppCategory{}
	// 检查名称是否存在
	if exist, _ := m.CheckNameExist(global.Db, param.Name, 0); exist {
		apiReturn.Error(c, "分类名称已存在")
		return
	}

	m.Name = param.Name
	m.Icon = param.Icon
	m.Sort = param.Sort
	m.Status = param.Status

	if err := m.Create(global.Db); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"id": m.ID})
}

// 更新分类
func (a *MicroAppCategoryApi) Update(c *gin.Context) {
	type ParamsStruct struct {
		Id     uint   `json:"id" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	param := ParamsStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	m := models.MicroAppCategory{}
	// 检查名称是否存在
	if exist, _ := m.CheckNameExist(global.Db, param.Name, param.Id); exist {
		apiReturn.Error(c, "分类名称已存在")
		return
	}

	updateData := map[string]interface{}{
		"name":   param.Name,
		"icon":   param.Icon,
		"sort":   param.Sort,
		"status": param.Status,
	}

	if err := m.Update(global.Db, param.Id, updateData); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// 删除分类
func (a *MicroAppCategoryApi) Deletes(c *gin.Context) {
	type ParamsStruct struct {
		Ids []uint `json:"ids" binding:"required"`
	}

	param := ParamsStruct{}
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

// 更新分类状态
func (a *MicroAppCategoryApi) UpdateStatus(c *gin.Context) {
	type ParamsStruct struct {
		Id     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"required"`
	}

	param := ParamsStruct{}
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

// 获取所有启用的分类（下拉选择用）
func (a *MicroAppCategoryApi) GetEnabledList(c *gin.Context) {
	m := models.MicroAppCategory{}
	list, err := m.GetEnabledList(global.Db)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}
