package admin

import (
	"crypto/rand"
	"encoding/hex"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// generateMicroAppId 生成唯一的微应用ID
func generateMicroAppId() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type MicroAppApi struct {
}

// GetList 获取微应用列表
func (a *MicroAppApi) GetList(c *gin.Context) {
	param := MicroAppGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	// 使用 GetListWithAllLangs 预加载多语言信息
	list, total, err := m.GetListWithAllLangs(global.Db, param.Page, param.Limit, param.Status, param.CategoryId, nil, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取微应用详情
func (a *MicroAppApi) GetInfo(c *gin.Context) {
	param := MicroAppGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	// 使用 Preload 预加载多语言信息
	err := global.Db.Preload("LangList").Where("id = ?", param.Id).First(&m).Error
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, m)
}

// Create 创建微应用
func (a *MicroAppApi) Create(c *gin.Context) {
	param := MicroAppCreateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// 如果没有提供 MicroAppId，则自动生成
	if param.MicroAppId == "" {
		newId, err := generateMicroAppId()
		if err != nil {
			apiReturn.ErrorDatabase(c, "生成应用ID失败")
			return
		}
		param.MicroAppId = newId
	}

	// 检查 MicroAppId 是否已存在
	m := models.MicroApp{}
	var existing models.MicroApp
	err := global.Db.Where("micro_app_id = ?", param.MicroAppId).First(&existing).Error
	if err == nil {
		apiReturn.Error(c, "应用ID已存在")
		return
	}

	m.MicroAppId = param.MicroAppId
	m.AppName = param.AppName
	m.AppIcon = param.AppIcon
	m.AppDesc = param.AppDesc
	m.Remark = param.Remark
	m.CategoryId = param.CategoryId
	m.ChargeType = param.ChargeType
	m.Price = param.Price
	m.AuthorId = param.AuthorId
	m.Status = 2 // 默认审核中
	m.Screenshots = param.Screenshots

	// 开启事务保存主应用和多语言信息
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 创建主应用
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		// 保存多语言信息
		if param.LangMap != nil {
			for lang, langInfo := range param.LangMap {
				if langInfo.AppName != "" || langInfo.AppDesc != "" {
					langModel := models.MicroAppLang{
						MicroAppId: m.MicroAppId,
						Lang:       lang,
						AppName:    langInfo.AppName,
						AppDesc:    langInfo.AppDesc,
					}
					if err := tx.Create(&langModel).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"id": m.ID, "microAppId": m.MicroAppId})
}

// Update 更新微应用
func (a *MicroAppApi) Update(c *gin.Context) {
	param := MicroAppUpdateReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// 先获取现有应用信息，用于获取 microAppId
	m := models.MicroApp{}
	existingApp, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	updateData := map[string]interface{}{
		"app_name":      param.AppName,
		"app_icon":      param.AppIcon,
		"app_desc":      param.AppDesc,
		"remark":       param.Remark,
		"category_id":   param.CategoryId,
		"charge_type":   param.ChargeType,
		"price":         param.Price,
		"screenshots":   param.Screenshots,
	}

	// 开启事务更新主应用和多语言信息
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 更新主应用
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", param.Id).Updates(updateData).Error; err != nil {
			return err
		}

		// 更新多语言信息
		if param.LangMap != nil {
			for lang, langInfo := range param.LangMap {
				var existing models.MicroAppLang
				err := tx.Where("micro_app_id = ? AND lang = ?", existingApp.MicroAppId, lang).First(&existing).Error

				if err == gorm.ErrRecordNotFound {
					// 不存在，创建
					if langInfo.AppName != "" || langInfo.AppDesc != "" {
						langModel := models.MicroAppLang{
							MicroAppId: existingApp.MicroAppId,
							Lang:       lang,
							AppName:    langInfo.AppName,
							AppDesc:    langInfo.AppDesc,
						}
						if err := tx.Create(&langModel).Error; err != nil {
							return err
						}
					}
				} else if err == nil {
					// 存在，更新
					tx.Model(&existing).Updates(map[string]interface{}{
						"app_name": langInfo.AppName,
						"app_desc": langInfo.AppDesc,
					})
				}
			}
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// Deletes 删除微应用
func (a *MicroAppApi) Deletes(c *gin.Context) {
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
func (a *MicroAppApi) UpdateStatus(c *gin.Context) {
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

// UpdateLang 更新微应用语言
func (a *MicroAppApi) UpdateLang(c *gin.Context) {
	param := MicroAppUpdateLangReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取应用信息
	m := models.MicroApp{}
	existingApp, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 开启事务更新语言信息
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		if param.LangMap != nil {
			for lang, langInfo := range param.LangMap {
				var existing models.MicroAppLang
				err := tx.Where("micro_app_id = ? AND lang = ?", existingApp.MicroAppId, lang).First(&existing).Error

				if err == gorm.ErrRecordNotFound {
					// 不存在，创建
					if langInfo.AppName != "" || langInfo.AppDesc != "" {
						langModel := models.MicroAppLang{
							MicroAppId: existingApp.MicroAppId,
							Lang:       lang,
							AppName:    langInfo.AppName,
							AppDesc:    langInfo.AppDesc,
						}
						if err := tx.Create(&langModel).Error; err != nil {
							return err
						}
					}
				} else if err == nil {
					// 存在，更新
					tx.Model(&existing).Updates(map[string]interface{}{
						"app_name": langInfo.AppName,
						"app_desc": langInfo.AppDesc,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
