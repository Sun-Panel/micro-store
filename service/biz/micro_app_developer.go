package biz

import (
	"crypto/rand"
	"encoding/hex"
	"sun-panel/global"
	"sun-panel/models"
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

// MicroAppDeveloperService 微应用开发者业务服务
type MicroAppDeveloperService struct{}

// Biz 层直接接收 MicroAppQueryOptions
func (s *MicroAppDeveloperService) GetDeveloperAppList(db *gorm.DB, opts models.MicroAppQueryOptions) ([]models.MicroApp, int64, error) {
	return new(models.MicroApp).GetListWithAllLangs(db, opts)
}

// GetDeveloperAppInfo 获取开发者应用详情（含权限验证）
// 返回生效版本 + 草稿版本（如果存在）
func (s *MicroAppDeveloperService) GetDeveloperAppInfo(db *gorm.DB, appId uint, developerId uint) (map[string]interface{}, error) {
	// // 先查询指定记录，验证权限
	// m := models.MicroApp{}
	// err := db.Preload("LangList").Where("id = ? AND author_id = ?", appId, developerId).First(&m).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return nil, NewBizError(ErrCodeAppNotFound)
	// 	}
	// 	return nil, err
	// }

	// // 查询作者名字
	// authorName := ""
	// if m.AuthorId > 0 {
	// 	var user models.User
	// 	if err := db.Select("name").Where("id = ?", m.AuthorId).First(&user).Error; err == nil {
	// 		authorName = user.Name
	// 	}
	// }

	// // 查询该微应用的最新草稿版本
	// var draftReview models.MicroAppReview
	// err = db.Where("micro_app_id = ? AND status = -1", m.MicroAppId).First(&draftReview).Error

	// // 计算审核状态：0-已通过 1-审核中 2-已拒绝 3-草稿
	// reviewStatus := 0 // 默认已通过
	// if draftReview.ID > 0 {
	// 	// 有草稿，说明是草稿状态
	// 	reviewStatus = 3
	// } else {
	// 	// 查询是否有待审核记录
	// 	var pendingReview models.MicroAppReview
	// 	err = db.Where("micro_app_id = ? AND status = 0", m.MicroAppId).First(&pendingReview).Error
	// 	if err == nil {
	// 		// 有待审核记录，说明是审核中
	// 		reviewStatus = 1
	// 	} else {
	// 		// 查询是否有被拒绝的记录
	// 		var rejectedReview models.MicroAppReview
	// 		err = db.Where("micro_app_id = ? AND status = 2", m.MicroAppId).Order("created_at DESC").First(&rejectedReview).Error
	// 		if err == nil {
	// 			// 有被拒绝的记录，说明是已拒绝
	// 			reviewStatus = 2
	// 		} else {
	// 			// 没有任何审核记录，默认为已通过
	// 			reviewStatus = 0
	// 		}
	// 	}
	// }

	// // 构建返回数据
	// result := map[string]interface{}{
	// 	"id":              m.ID,
	// 	"microAppId":      m.MicroAppId,
	// 	"appName":         m.AppName,
	// 	"appIcon":         m.AppIcon,
	// 	"appDesc":         m.AppDesc,
	// 	"remark":          m.Remark,
	// 	"categoryId":      m.CategoryId,
	// 	"chargeType":      m.ChargeType,
	// 	"points":          m.Points,
	// 	"authorId":        m.AuthorId,
	// 	"authorName":      authorName,
	// 	"permissionLevel": m.PermissionLevel,
	// 	"status":          m.Status,
	// 	"screenshots":     m.Screenshots,
	// 	"langList":        m.LangList,
	// 	"createTime":      m.CreatedAt,
	// 	"updateTime":      m.UpdatedAt,
	// 	"reviewStatus":    reviewStatus, // 审核状态：0-已通过 1-审核中 2-已拒绝 3-草稿
	// 	"draft":           draftReview,  // 返回草稿版本（如果存在）
	// }

	// return result, nil
	return nil, nil
}

// GetDeveloperAppExtendInfo 获取开发者应用扩展信息（审核状态、草稿版本）
// 返回：审核状态、草稿版本（如果存在）
func (s *MicroAppDeveloperService) GetDeveloperAppExtendInfo(db *gorm.DB, microAppId string) (int, *models.MicroAppReview, error) {
	// 查询该微应用的最新草稿版本
	var draftReview models.MicroAppReview
	err := db.Where("micro_app_id = ? AND status = -1", microAppId).First(&draftReview).Error

	// 计算审核状态：0-已通过 1-审核中 2-已拒绝 3-草稿
	reviewStatus := 0 // 默认已通过
	if err == nil {
		// 有草稿，说明是草稿状态
		reviewStatus = 3
		return reviewStatus, &draftReview, nil
	}
	if err != gorm.ErrRecordNotFound {
		return 0, nil, err
	}

	// 查询是否有待审核记录
	var pendingReview models.MicroAppReview
	err = db.Where("micro_app_id = ? AND status = 0", microAppId).First(&pendingReview).Error
	if err == nil {
		// 有待审核记录，说明是审核中
		reviewStatus = 1
		return reviewStatus, nil, nil
	}

	// 查询是否有被拒绝的记录
	var rejectedReview models.MicroAppReview
	err = db.Where("micro_app_id = ? AND status = 2", microAppId).Order("created_at DESC").First(&rejectedReview).Error
	if err == nil {
		// 有被拒绝的记录，说明是已拒绝
		reviewStatus = 2
		return reviewStatus, nil, nil
	}

	// 没有任何审核记录，默认为已通过
	reviewStatus = 0
	return reviewStatus, nil, nil
}

type DeveloperAppOptions struct {
	ID               uint // 更新时使用（appId）
	MicroAppBaseInfo models.MicroAppBaseInfo
	MicroAppId       string           // 创建时使用
	LangMap          datatype.MapJson `json:"langMap"`
	DeveloperId      uint             `json:"developerId"`
}

// 同时在主表创建下架记录和审核表创建草稿记录
func (s *MicroAppDeveloperService) CreateAppAndReview(db *gorm.DB, opts DeveloperAppOptions) (map[string]interface{}, error) {
	// 检查 MicroAppId 是否已在主表中存在
	var existing models.MicroApp
	if err := db.Where("micro_app_id = ?", opts.MicroAppId).First(&existing).Error; err == nil {
		return nil, NewBizError(ErrCodeAppIdExists)
	}

	m := models.MicroApp{
		MicroAppBaseInfo: opts.MicroAppBaseInfo,
		MicroAppId:       opts.MicroAppId, // 同时设置外层字段
		DeveloperId:      opts.DeveloperId,
		Status:           0, // 默认下架，需要提交审核通过后才能上架
		OfflineType:      3,
	}

	// 事务保存主应用、审核记录和多语言信息
	err := db.Transaction(func(tx *gorm.DB) error {
		// 保存主应用记录
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		// 创建审核记录（草稿状态，status = -1）
		review := models.MicroAppReview{
			MicroAppBaseInfo: opts.MicroAppBaseInfo,
			MicroAppId:       opts.MicroAppId,
			AppRecordId:      m.ID, // 关联主表记录ID
			LangMap:          opts.LangMap,
			Status:           -1, // 草稿状态
		}

		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		// 保存多语言信息到 micro_app_lang 表
		for lang, langInfo := range opts.LangMap {
			// langInfo 的类型是 interface{}，需要类型断言
			infoMap, ok := langInfo.(map[string]interface{})
			if !ok {
				continue
			}
			appName, _ := infoMap["appName"].(string)
			appDesc, _ := infoMap["appDesc"].(string)
			if appName != "" || appDesc != "" {
				langModel := models.MicroAppLang{
					MicroAppId: m.MicroAppId,
					Lang:       lang,
					AppName:    appName,
					AppDesc:    appDesc,
				}
				if err := tx.Create(&langModel).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         m.ID,
		"microAppId": m.MicroAppId,
	}, nil
}

// SubmitAppReview 提交应用审核
// 提交 micro_app_review 表中的草稿版本
func (s *MicroAppDeveloperService) SubmitAppReview(db *gorm.DB, reviewId, developerId uint) error {
	// 根据审核ID获取数据
	mr := models.MicroAppReview{}
	review, err := mr.GetById(db, reviewId)
	if err != nil {
		return err
	}

	// 获取生效版本，验证权限
	app, err := s.getAppAndCheckPermission(db, review.AppRecordId, developerId)
	if err != nil {
		return err
	}

	// 检查是否已有待审核记录
	var count int64
	db.Model(&models.MicroAppReview{}).Where("micro_app_id = ? AND status = 0", app.MicroAppId).Count(&count)
	if count > 0 {
		return NewBizError(ErrCodePendingReviewExists)
	}

	// // 查找草稿版本
	// var draft models.MicroAppReview
	// err = db.Where("micro_app_id = ? AND status = -1", app.MicroAppId).First(&draft).Error
	// if err == gorm.ErrRecordNotFound {
	// 	// 没有草稿版本，说明没有修改
	// 	return NewBizError(ErrCodeNoUpdateContent)
	// }
	// if err != nil {
	// 	return err
	// }

	// // 检查必填项
	// if review.AppName == "" {
	// 	return NewBizError(ErrCodeInvalidParam)
	// }
	if review.AppIcon == "" {
		return NewBizError(ErrCodeInvalidParam)
	}
	if review.CategoryId == 0 {
		return NewBizError(ErrCodeInvalidParam)
	}

	// 更新草稿状态为待审核
	err = db.Model(&models.MicroAppReview{}).Where("id = ?", review.ID).Update("status", 0).Error
	if err != nil {
		return err
	}

	return nil
}

// GetOrCreateDraftApp 获取或创建草稿版本
// 在 micro_app_review 表中查找或创建草稿
func (s *MicroAppDeveloperService) GetOrCreateDraftApp(db *gorm.DB, reviewId, developerId uint) (*models.MicroAppReview, error) {
	var review models.MicroAppReview
	review, err := review.GetById(db, reviewId)
	if err != nil {
		return nil, err
	}

	// 获取生效版本，验证权限
	app, err := s.getAppAndCheckPermission(db, review.AppRecordId, developerId)
	if err != nil {
		return nil, err
	}

	// 查找已存在的草稿
	var draft models.MicroAppReview
	err = db.Where("micro_app_id = ? AND status = -1", app.MicroAppId).Order("created_at DESC").First(&draft).Error

	if err == nil {
		// 已有草稿，返回
		return &draft, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 合并多语言信息
	mergedLangMap, err := s.mergeLangMap(db, app.MicroAppId, nil)
	if err != nil {
		return nil, err
	}

	// 将 mergedLangMap (map[string]map[string]interface{}) 转换为 datatype.MapJson (map[string]interface{})
	langMap := make(map[string]interface{})
	for k, v := range mergedLangMap {
		langMap[k] = v
	}

	// 创建新的草稿记录
	draft = models.MicroAppReview{
		MicroAppBaseInfo: models.MicroAppBaseInfo{
			// AppName:     app.AppName,
			AppIcon: app.AppIcon,
			// AppDesc:     app.AppDesc,
			CategoryId:  app.CategoryId,
			ChargeType:  app.ChargeType,
			Points:      app.Points,
			Screenshots: app.Screenshots,
			Remark:      app.Remark,
		},
		MicroAppId:  app.MicroAppId, // 同时设置外层字段
		AppRecordId: app.ID,
		LangMap:     langMap,
		Status:      -1, // 草稿状态
	}

	if err := db.Create(&draft).Error; err != nil {
		return nil, err
	}

	return &draft, nil
}

// UpdateApp 更新应用信息（更新草稿版本）
// 自动获取或创建草稿版本，然后更新
func (s *MicroAppDeveloperService) UpdateDraftApp(db *gorm.DB, opts DeveloperAppOptions) error {
	// // 获取生效版本，验证权限
	// app, err := s.getAppAndCheckPermission(db, opts.ID, opts.DeveloperId)
	// if err != nil {
	// 	return err
	// }

	// // 检查是否已在审核中
	// var pendingReviewCount int64
	// db.Model(&models.MicroAppReview{}).Where("micro_app_id = ? AND status = 0", app.MicroAppId).Count(&pendingReviewCount)
	// if pendingReviewCount > 0 {
	// 	return NewBizError(ErrCodePendingReviewExists)
	// }

	// 获取或创建草稿
	draft, err := s.GetOrCreateDraftApp(db, opts.ID, opts.DeveloperId)
	if err != nil {
		return err
	}

	global.Logger.Debug("opts", opts)

	// 更新草稿记录
	err = db.Transaction(func(tx *gorm.DB) error {
		// 更新应用基本信息到审核表
		if err := tx.Model(&models.MicroAppReview{}).Where("id = ?", draft.ID).Updates(map[string]interface{}{
			// "app_name":    opts.MicroAppBaseInfo.AppName,
			"app_icon": opts.MicroAppBaseInfo.AppIcon,
			// "app_desc":    opts.MicroAppBaseInfo.AppDesc,
			"remark":      opts.MicroAppBaseInfo.Remark,
			"category_id": opts.MicroAppBaseInfo.CategoryId,
			"charge_type": opts.MicroAppBaseInfo.ChargeType,
			"points":      opts.MicroAppBaseInfo.Points,
			"screenshots": opts.MicroAppBaseInfo.Screenshots,
			"lang_map":    opts.LangMap,
			"admin_name":  opts.MicroAppBaseInfo.AdminName,
		}).Error; err != nil {
			return err
		}

		// 同时将 admin_name 和 remark 直接更新到 micro_app 表
		if err := tx.Model(&models.MicroApp{}).Where("micro_app_id = ?", draft.MicroAppId).Updates(map[string]interface{}{
			"admin_name": opts.MicroAppBaseInfo.AdminName,
			"remark":     opts.MicroAppBaseInfo.Remark,
		}).Error; err != nil {
			return err
		}

		// 更新多语言信息到 micro_app_lang 表（所有版本共享）
		for lang, langInfo := range opts.LangMap {
			// langInfo 的类型是 interface{}，需要类型断言
			infoMap, ok := langInfo.(map[string]interface{})
			if !ok {
				continue
			}
			appName, _ := infoMap["appName"].(string)
			appDesc, _ := infoMap["appDesc"].(string)

			// 查找是否已存在该语言的记录
			var existLang models.MicroAppLang
			err := tx.Where("micro_app_id = ? AND lang = ?", draft.MicroAppId, lang).First(&existLang).Error

			if err == gorm.ErrRecordNotFound {
				// 创建新的语言记录
				if appName != "" || appDesc != "" {
					langModel := models.MicroAppLang{
						MicroAppId: draft.MicroAppId,
						Lang:       lang,
						AppName:    appName,
						AppDesc:    appDesc,
					}
					if err := tx.Create(&langModel).Error; err != nil {
						return err
					}
				}
			} else if err == nil {
				// 更新已有的语言记录
				if err := tx.Model(&models.MicroAppLang{}).Where("id = ?", existLang.ID).Updates(map[string]interface{}{
					"app_name": appName,
					"app_desc": appDesc,
				}).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	return err
}

// UpdateLang 更新语言信息（不提交审核）
func (s *MicroAppDeveloperService) UpdateLang(db *gorm.DB, id, developerId uint, langMap map[string]interface{}) error {
	// 获取并验证权限
	app, err := s.getAppAndCheckPermission(db, id, developerId)
	if err != nil {
		return err
	}

	// 检查是否已在审核中
	var pendingReviewCount int64
	db.Model(&models.MicroAppReview{}).Where("micro_app_id = ? AND status = 0", app.MicroAppId).Count(&pendingReviewCount)
	if pendingReviewCount > 0 {
		return NewBizError(ErrCodePendingReviewExists)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 更新多语言信息
		for lang, langInfo := range langMap {
			// langInfo 的类型是 interface{}，需要类型断言
			infoMap, ok := langInfo.(map[string]interface{})
			if !ok {
				continue
			}
			appName, _ := infoMap["appName"].(string)
			appDesc, _ := infoMap["appDesc"].(string)

			// 查找是否已存在该语言的记录
			var existLang models.MicroAppLang
			err := tx.Where("micro_app_id = ? AND lang = ?", app.MicroAppId, lang).First(&existLang).Error

			if err == gorm.ErrRecordNotFound {
				// 创建新的语言记录
				if appName != "" || appDesc != "" {
					langModel := models.MicroAppLang{
						MicroAppId: app.MicroAppId,
						Lang:       lang,
						AppName:    appName,
						AppDesc:    appDesc,
					}
					if err := tx.Create(&langModel).Error; err != nil {
						return err
					}
				}
			} else if err == nil {
				// 更新已有的语言记录
				if err := tx.Model(&models.MicroAppLang{}).Where("id = ?", existLang.ID).Updates(map[string]interface{}{
					"app_name": appName,
					"app_desc": appDesc,
				}).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	return err
}

// CancelAppReview 撤销应用审核
// 将待审核的记录恢复为草稿状态
func (s *MicroAppDeveloperService) CancelAppReview(db *gorm.DB, reviewId, developerId uint) error {
	// 根据审核ID获取数据
	mr := models.MicroAppReview{}
	review, err := mr.GetById(db, reviewId)
	if err != nil {
		return err
	}

	// 获取生效版本，验证权限
	_, err = s.getAppAndCheckPermission(db, review.AppRecordId, developerId)
	if err != nil {
		return err
	}

	if review.Status != 0 {
		return NewBizError(ErrCodeStatusNotAllowed)
	}

	// // 查找待审核的记录
	// var review models.MicroAppReview
	// err = db.Where("micro_app_id = ? AND status = 0", app.MicroAppId).First(&review).Error
	// if err == gorm.ErrRecordNotFound {
	// 	return NewBizError(ErrCodeNoPendingReviewApp)
	// }
	// if err != nil {
	// 	return err
	// }

	// 将审核记录恢复为草稿状态
	return db.Model(&models.MicroAppReview{}).Where("id = ?", review.ID).Update("status", -1).Error
}

// GetAppReviewHistory 获取应用审核历史
func (s *MicroAppDeveloperService) GetAppReviewHistory(db *gorm.DB, appId, developerId uint, page, limit int) ([]models.MicroAppReview, int64, error) {
	// 验证权限
	if _, err := s.getAppAndCheckPermission(db, appId, developerId); err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	reviewModel := models.MicroAppReview{}
	return reviewModel.GetListByAppId(db, appId, page, limit)
}

// generateMicroAppId 生成唯一的微应用ID
func generateMicroAppId() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getAppAndCheckPermission 获取应用并验证开发者权限
func (s *MicroAppDeveloperService) getAppAndCheckPermission(db *gorm.DB, appId, developerId uint) (models.MicroApp, error) {
	m := models.MicroApp{}
	app, err := m.GetById(db, appId)
	if err != nil {
		return models.MicroApp{}, NewBizError(ErrCodeAppNotFound)
	}

	if app.DeveloperId != developerId {
		return models.MicroApp{}, NewBizError(ErrCodeNoPermission)
	}

	return app, nil
}

// checkNoPendingReview 检查是否没有待审核记录
func (s *MicroAppDeveloperService) checkNoPendingReview(db *gorm.DB, appId uint) error {
	reviewModel := models.MicroAppReview{}
	_, err := reviewModel.GetPendingByAppId(db, appId)
	if err == nil {
		return NewBizError(ErrCodePendingReviewExists)
	}
	return nil
}

// checkNoPendingReviewByMicroAppId 根据 microAppId 检查是否没有待审核记录
func (s *MicroAppDeveloperService) checkNoPendingReviewByMicroAppId(db *gorm.DB, microAppId string) error {
	var count int64
	db.Model(&models.MicroAppReview{}).Where("micro_app_id = ? AND status = 0", microAppId).Count(&count)
	if count > 0 {
		return NewBizError(ErrCodePendingReviewExists)
	}
	return nil
}

// mergeLangMap 合并现有多语言信息和新的多语言信息
func (s *MicroAppDeveloperService) mergeLangMap(db *gorm.DB, microAppId string, newLangMap map[string]map[string]interface{}) (map[string]map[string]interface{}, error) {
	langList, err := (&models.MicroAppLang{}).GetListByAppId(db, microAppId)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]map[string]interface{})
	for _, lang := range langList {
		merged[lang.Lang] = map[string]interface{}{
			"appName": lang.AppName,
			"appDesc": lang.AppDesc,
		}
	}

	for lang, langInfo := range newLangMap {
		merged[lang] = langInfo
	}

	return merged, nil
}
