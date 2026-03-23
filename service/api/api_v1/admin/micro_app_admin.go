package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// MicroAppAdminApi 管理员专用微应用 API
type MicroAppAdminApi struct{}

// GetList 获取微应用列表（管理员专用）
func (a *MicroAppAdminApi) GetList(c *gin.Context) {
	param := MicroAppGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	// 管理员可以查看所有应用，不限制 authorId
	opts := models.MicroAppQueryOptions{
		Page:             param.Page,
		Limit:            param.Limit,
		Status:           param.Status,
		CategoryId:       param.CategoryId,
		KeyWord:          param.KeyWord,
		IncludeDeveloper: true,
	}
	list, total, err := m.GetListWithAllLangs(global.Db, opts)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取微应用详情（管理员专用）
func (a *MicroAppAdminApi) GetInfo(c *gin.Context) {
	param := MicroAppGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	// 管理员可以查看任何应用详情
	err := global.Db.Preload("LangList").Where("id = ?", param.Id).First(&m).Error
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 查询作者名字
	var authorName string
	if m.AuthorId > 0 {
		var user models.User
		if err := global.Db.Select("name").Where("id = ?", m.AuthorId).First(&user).Error; err == nil {
			authorName = user.Name
		}
	}

	// 返回数据，包含作者名字
	result := gin.H{
		"id":              m.ID,
		"microAppId":      m.MicroAppId,
		"appName":         m.AppName,
		"appIcon":         m.AppIcon,
		"appDesc":         m.AppDesc,
		"remark":          m.Remark,
		"categoryId":      m.CategoryId,
		"chargeType":      m.ChargeType,
		"points":          m.Points,
		"authorId":        m.AuthorId,
		"authorName":      authorName,
		"permissionLevel": m.PermissionLevel,
		"status":          m.Status,
		"screenshots":     m.Screenshots,
		"langList":        m.LangList,
		"createTime":      m.CreatedAt,
		"updateTime":      m.UpdatedAt,
	}

	apiReturn.SuccessData(c, result)
}

// Deletes 删除微应用（管理员专用）
func (a *MicroAppAdminApi) Deletes(c *gin.Context) {
	param := MicroAppDeletesReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 开启事务删除应用及其多语言数据
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		// 先获取要删除的应用的 microAppId
		var apps []models.MicroApp
		if err := tx.Where("id IN ?", param.Ids).Find(&apps).Error; err != nil {
			return err
		}

		// 删除多语言数据
		for _, app := range apps {
			if err := tx.Where("micro_app_id = ?", app.MicroAppId).Delete(&models.MicroAppLang{}).Error; err != nil {
				return err
			}
		}

		// 删除应用
		if err := tx.Where("id IN ?", param.Ids).Delete(&models.MicroApp{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// UpdateStatus 更新微应用状态（管理员专用）
func (a *MicroAppAdminApi) UpdateStatus(c *gin.Context) {
	param := MicroAppUpdateStatusReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	if err := m.UpdateStatus(global.Db, param.Id, param.Status); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Offline 下架微应用（管理员专用）
func (a *MicroAppAdminApi) Offline(c *gin.Context) {
	param := MicroAppOfflineReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 平台下架时，原因必填
	if param.Type == 2 && param.Reason == "" {
		apiReturn.ErrorParamFomat(c, "平台下架时，下架原因不能为空")
		return
	}

	m := models.MicroApp{}
	if err := m.Offline(global.Db, param.Id, param.Type, param.Reason); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
