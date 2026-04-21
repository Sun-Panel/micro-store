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
	opts := models.MicroAppListQueryOpts{
		Page:             param.Page,
		Limit:            param.Limit,
		Status:           param.Status,
		CategoryId:       param.CategoryId,
		KeyWord:          param.KeyWord,
		IncludeDeveloper: true,
		SortBy:           param.SortBy,
		SortOrder:        param.SortOrder,
		IncludeLangList:  true,
	}
	list, total, err := m.GetAppList(global.Db, opts)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// UpdateStatus 更新微应用状态
func (a *MicroAppAdminApi) GetInfo(c *gin.Context) {
	req := models.BaseModel{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	lang := base.GetCurrentUserLang(c)

	// 获取基础信息（有缓存优化）
	info, err := biz.MicroApp.GetByIdWithLang(global.Db.Debug(), req.ID, lang, "Developer", "LangList")
	if err != nil {
		handleBizError(c, err)
		return
	}

	// 获取扩展信息（审核状态、草稿）
	reviewStatus, draft, err := biz.MicroAppDeveloper.GetDeveloperAppExtendInfo(global.Db, info.MicroAppId)
	if err != nil {
		handleBizError(c, err)
		return
	}

	// 组合数据返回
	result := map[string]interface{}{
		"id":         info.ID,
		"microAppId": info.MicroAppId,
		// "appName":       info.AppName,
		"appIcon": info.AppIcon,
		// "appDesc":       info.AppDesc,
		"remark":        info.Remark,
		"categoryId":    info.CategoryId,
		"chargeType":    info.ChargeType,
		"points":        info.Points,
		"authorId":      info.DeveloperId,
		"status":        info.Status,
		"screenshots":   info.Screenshots,
		"langList":      info.LangList,
		"createTime":    info.CreatedAt,
		"updateTime":    info.UpdatedAt,
		"reviewStatus":  reviewStatus,       // 审核状态：0-已通过 1-审核中 2-已拒绝 3-草稿
		"draft":         draft,              // 草稿版本（如果存在）
		"offlineReason": info.OfflineReason, // 下线原因
	}

	apiReturn.SuccessData(c, result)
}

// // GetInfo 获取微应用详情（管理员专用）
// func (a *MicroAppAdminApi) GetInfo(c *gin.Context) {
// 	param := MicroAppGetInfoReq{}
// 	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	m := models.MicroApp{}
// 	// 管理员可以查看任何应用详情
// 	err := global.Db.Preload("LangList").Where("id = ?", param.Id).First(&m).Error
// 	if err != nil {
// 		apiReturn.ErrorDataNotFound(c)
// 		return
// 	}

// 	// 查询作者名字
// 	var authorName string
// 	if m.DeveloperId > 0 {
// 		var user models.User
// 		if err := global.Db.Select("name").Where("id = ?", m.DeveloperId).First(&user).Error; err == nil {
// 			authorName = user.Name
// 		}
// 	}

// 	// 返回数据，包含作者名字
// 	result := gin.H{
// 		"id":          m.ID,
// 		"microAppId":  m.MicroAppId,
// 		"appName":     m.AppName,
// 		"appIcon":     m.AppIcon,
// 		"appDesc":     m.AppDesc,
// 		"remark":      m.Remark,
// 		"categoryId":  m.CategoryId,
// 		"chargeType":  m.ChargeType,
// 		"points":      m.Points,
// 		"authorId":    m.DeveloperId,
// 		"authorName":  authorName,
// 		"status":      m.Status,
// 		"screenshots": m.Screenshots,
// 		"langList":    m.LangList,
// 		"createTime":  m.CreatedAt,
// 		"updateTime":  m.UpdatedAt,
// 	}

// 	apiReturn.SuccessData(c, result)
// }

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

// UpdateStatus 更新微应用状态
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
	req := MicroAppOfflineReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 平台下架时，原因必填
	if req.OfflineType == 2 && req.Reason == "" {
		apiReturn.ErrorParamFomat(c, "平台下架时，下架原因不能为空")
		return
	}

	m := models.MicroApp{}
	if err := m.Offline(global.Db, req.Id, req.OfflineType, req.Reason); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
