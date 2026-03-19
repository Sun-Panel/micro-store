package admin

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// ==================== 微应用 API ====================
// 本文件包含所有微应用相关的接口，按角色划分：
// - 管理员专用接口
// - 审核员专用接口
// - 开发者专用接口

type MicroAppApi struct {
}

// ==================== 管理员专用接口 ====================

// GetList 获取微应用列表（管理员专用）
func (a *MicroAppApi) GetList(c *gin.Context) {
	param := MicroAppGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	m := models.MicroApp{}
	// 管理员可以查看所有应用，不限制 authorId
	list, total, err := m.GetListWithAllLangs(global.Db, param.Page, param.Limit, param.Status, param.CategoryId, nil, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetInfo 获取微应用详情（管理员专用）
func (a *MicroAppApi) GetInfo(c *gin.Context) {
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
		"price":           m.Price,
		"authorId":        m.AuthorId,
		"authorName":      authorName,
		"permissionLevel": m.PermissionLevel,
		"status":          m.Status,
		"screenshots":     m.Screenshots,
		"reviewStatus":    m.ReviewStatus,
		"reviewId":        m.ReviewId,
		"reviewTime":      m.ReviewTime,
		"langList":        m.LangList,
		"createTime":      m.CreatedAt,
		"updateTime":      m.UpdatedAt,
	}

	apiReturn.SuccessData(c, result)
}

// Deletes 删除微应用（管理员专用）
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

// UpdateStatus 更新微应用状态（管理员专用）
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

// Offline 下架微应用（管理员专用）
func (a *MicroAppApi) Offline(c *gin.Context) {
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

// ==================== 审核员专用接口 ====================

// GetReviewPendingList 获取待审核应用列表（审核员专用）
func (a *MicroAppApi) GetReviewPendingList(c *gin.Context) {
	param := MicroAppGetPendingReviewListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reviewModel := models.MicroAppReview{}
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	list, total, err := reviewModel.GetPendingList(global.Db, param.Page, param.Limit)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetReviewInfo 获取审核详情（审核员专用）
func (a *MicroAppApi) GetReviewInfo(c *gin.Context) {
	param := MicroAppGetReviewInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reviewModel := models.MicroAppReview{}
	review, err := reviewModel.GetById(global.Db, param.ReviewId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, review)
}

// ReviewApp 审核微应用主信息（审核员专用）
func (a *MicroAppApi) ReviewApp(c *gin.Context) {
	param := MicroAppReviewAppReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 验证审核状态
	if param.Status != 1 && param.Status != 2 {
		apiReturn.ErrorParamFomat(c, "审核状态无效")
		return
	}

	// 获取审核记录
	reviewModel := models.MicroAppReview{}
	review, err := reviewModel.GetById(global.Db, param.ReviewId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 检查是否已审核
	if review.Status != 0 {
		apiReturn.Error(c, "该记录已审核")
		return
	}

	// 获取当前审核员ID（从token中获取）
	auditorId := c.GetUint("adminId")

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 更新审核记录状态
		now := time.Now()
		if err := tx.Model(&models.MicroAppReview{}).Where("id = ?", param.ReviewId).Updates(map[string]interface{}{
			"status":       param.Status,
			"reviewer_id":  auditorId,
			"review_note":  param.ReviewNote,
			"review_time":  now,
		}).Error; err != nil {
			return err
		}

		// 如果审核通过，将审核快照数据更新到主表
		if param.Status == 1 {
			// 反序列化多语言信息
			var langMap map[string]MicroAppLangInfo
			json.Unmarshal([]byte(review.LangMap), &langMap)

			// 更新主表
			if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppId).Updates(map[string]interface{}{
				"app_name":      review.AppName,
				"app_icon":      review.AppIcon,
				"app_desc":      review.AppDesc,
				"category_id":   review.CategoryId,
				"charge_type":   review.ChargeType,
				"price":         review.Price,
				"screenshots":   review.Screenshots,
				"remark":        review.Remark,
				"review_status": 2, // 已通过
				"review_time":   &now,
			}).Error; err != nil {
				return err
			}

			// 更新多语言信息
			m := models.MicroApp{}
			app, _ := m.GetById(tx, review.AppId)
			if app.MicroAppId != "" {
				// 先删除当前的多语言记录
				if err := tx.Where("micro_app_id = ?", app.MicroAppId).Delete(&models.MicroAppLang{}).Error; err != nil {
					return err
				}

				// 然后插入新的多语言记录
				for lang, langInfo := range langMap {
					// 先尝试查找已存在的记录（包括软删除的）
					var existingLang models.MicroAppLang
					err := tx.Unscoped().Where("micro_app_id = ? AND lang = ?", app.MicroAppId, lang).First(&existingLang).Error

					if err == nil {
						// 记录已存在，更新它（包括软删除的记录）
						existingLang.AppName = langInfo.AppName
						existingLang.AppDesc = langInfo.AppDesc
						existingLang.DeletedAt = gorm.DeletedAt{} // 取消软删除标记
						if err := tx.Save(&existingLang).Error; err != nil {
							return err
						}
					} else {
						// 记录不存在，创建新记录
						langModel := models.MicroAppLang{
							MicroAppId: app.MicroAppId,
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
		} else {
			// 审核拒绝
			if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppId).Updates(map[string]interface{}{
				"review_status": 3, // 已拒绝
				"review_time":   &now,
			}).Error; err != nil {
				return err
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

// ==================== 开发者专用接口 ====================

// generateMicroAppId 生成唯一的微应用ID
func generateMicroAppId() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetMyList 获取开发者的微应用列表（开发者专用）
func (a *MicroAppApi) GetMyList(c *gin.Context) {
	param := MicroAppGetListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	m := models.MicroApp{}
	// 只查询当前开发者的应用
	list, total, err := m.GetListWithAllLangs(global.Db, param.Page, param.Limit, param.Status, param.CategoryId, &developerId, param.KeyWord)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetMyInfo 获取开发者微应用详情（开发者专用）
func (a *MicroAppApi) GetMyInfo(c *gin.Context) {
	param := MicroAppGetInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	m := models.MicroApp{}
	err := global.Db.Preload("LangList").Where("id = ? AND author_id = ?", param.Id, developerId).First(&m).Error
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
		"price":           m.Price,
		"authorId":        m.AuthorId,
		"authorName":      authorName,
		"permissionLevel": m.PermissionLevel,
		"status":          m.Status,
		"screenshots":     m.Screenshots,
		"reviewStatus":    m.ReviewStatus,
		"reviewId":        m.ReviewId,
		"reviewTime":      m.ReviewTime,
		"langList":        m.LangList,
		"createTime":      m.CreatedAt,
		"updateTime":      m.UpdatedAt,
	}

	apiReturn.SuccessData(c, result)
}

// Create 创建微应用（开发者专用）
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

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")
	// 强制使用当前开发者ID
	param.AuthorId = developerId

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

// Update 更新微应用（开发者专用，修改即提交审核）
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

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	// 获取现有应用信息（并验证是否属于当前开发者）
	m := models.MicroApp{}
	existingApp, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 验证应用是否属于当前开发者
	if existingApp.AuthorId != developerId {
		apiReturn.Error(c, "无权操作此应用")
		return
	}

	// 检查是否已有待审核的记录
	reviewModel := models.MicroAppReview{}
	_, err = reviewModel.GetPendingByAppId(global.Db, param.Id)
	if err == nil {
		apiReturn.Error(c, "已有待审核的记录，请等待审核完成")
		return
	}

	// 获取现有多语言信息
	langList, err := (&models.MicroAppLang{}).GetListByAppId(global.Db, existingApp.MicroAppId)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 序列化多语言信息（新的数据）
	langMap := make(map[string]MicroAppLangInfo)
	// 合并现有的多语言信息和新的多语言信息
	for _, lang := range langList {
		langMap[lang.Lang] = MicroAppLangInfo{
			AppName: lang.AppName,
			AppDesc: lang.AppDesc,
		}
	}
	// 用新的数据覆盖
	if param.LangMap != nil {
		for lang, langInfo := range param.LangMap {
			langMap[lang] = langInfo
		}
	}
	langMapJson, _ := json.Marshal(langMap)

	// 创建审核快照记录
	review := models.MicroAppReview{
		AppId:      param.Id,
		AppName:    param.AppName,
		AppIcon:    param.AppIcon,
		AppDesc:    param.AppDesc,
		CategoryId: param.CategoryId,
		ChargeType: param.ChargeType,
		Price:      param.Price,
		Screenshots: param.Screenshots,
		LangMap:    string(langMapJson),
		Remark:     param.Remark,
		Status:     0, // 待审核
	}

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 创建审核记录
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		// 更新主表的审核状态（不更新其他字段）
		now := time.Now()
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
			"review_status": 1, // 审核中
			"review_id":     review.ID,
			"review_time":   &now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"reviewId": review.ID})
}

// UpdateLang 更新微应用语言（开发者专用，修改即提交审核）
func (a *MicroAppApi) UpdateLang(c *gin.Context) {
	param := MicroAppUpdateLangReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	// 获取应用信息（并验证是否属于当前开发者）
	m := models.MicroApp{}
	existingApp, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 验证应用是否属于当前开发者
	if existingApp.AuthorId != developerId {
		apiReturn.Error(c, "无权操作此应用")
		return
	}

	// 检查是否已有待审核的记录
	reviewModel := models.MicroAppReview{}
	_, err = reviewModel.GetPendingByAppId(global.Db, param.Id)
	if err == nil {
		apiReturn.Error(c, "已有待审核的记录，请等待审核完成")
		return
	}

	// 获取现有多语言信息
	langList, err := (&models.MicroAppLang{}).GetListByAppId(global.Db, existingApp.MicroAppId)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 序列化多语言信息（合并现有和新数据）
	langMap := make(map[string]MicroAppLangInfo)
	for _, lang := range langList {
		langMap[lang.Lang] = MicroAppLangInfo{
			AppName: lang.AppName,
			AppDesc: lang.AppDesc,
		}
	}
	// 用新的数据覆盖
	if param.LangMap != nil {
		for lang, langInfo := range param.LangMap {
			langMap[lang] = langInfo
		}
	}
	langMapJson, _ := json.Marshal(langMap)

	// 创建审核快照记录
	review := models.MicroAppReview{
		AppId:       param.Id,
		AppName:     existingApp.AppName,
		AppIcon:     existingApp.AppIcon,
		AppDesc:     existingApp.AppDesc,
		CategoryId:  existingApp.CategoryId,
		ChargeType:  existingApp.ChargeType,
		Price:       existingApp.Price,
		Screenshots: existingApp.Screenshots,
		LangMap:     string(langMapJson),
		Remark:      existingApp.Remark,
		Status:      0, // 待审核
	}

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 创建审核记录
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		// 更新主表的审核状态（不更新其他字段）
		now := time.Now()
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
			"review_status": 1, // 审核中
			"review_id":     review.ID,
			"review_time":   &now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"reviewId": review.ID})
}

// CancelReview 撤销微应用主信息审核（开发者专用）
func (a *MicroAppApi) CancelReview(c *gin.Context) {
	param := MicroAppCancelReviewReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	// 获取应用信息（并验证是否属于当前开发者）
	m := models.MicroApp{}
	app, err := m.GetById(global.Db, param.Id)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 验证应用是否属于当前开发者
	if app.AuthorId != developerId {
		apiReturn.Error(c, "无权操作此应用")
		return
	}

	// 检查是否有待审核的记录
	if app.ReviewId == 0 {
		apiReturn.Error(c, "没有待审核的记录")
		return
	}

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 删除审核记录
		if err := tx.Where("id = ?", app.ReviewId).Delete(&models.MicroAppReview{}).Error; err != nil {
			return err
		}

		// 更新主表审核状态
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
			"review_status": 0, // 无审核
			"review_id":     0,
			"review_time":   nil,
		}).Error; err != nil {
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

// GetReviewHistory 获取微应用审核历史（开发者专用）
func (a *MicroAppApi) GetReviewHistory(c *gin.Context) {
	param := MicroAppGetReviewHistoryReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 获取当前开发者ID
	developerId := c.GetUint("adminId")

	// 验证应用是否属于当前开发者
	m := models.MicroApp{}
	app, err := m.GetById(global.Db, param.AppId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	if app.AuthorId != developerId {
		apiReturn.Error(c, "无权查看此应用的审核历史")
		return
	}

	reviewModel := models.MicroAppReview{}
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	list, total, err := reviewModel.GetListByAppId(global.Db, param.AppId, param.Page, param.Limit)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}
