package biz

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

// MicroAppDeveloperService 微应用开发者业务服务
type MicroAppDeveloperService struct{}

// GetDeveloperAppList 获取开发者应用列表
func (s *MicroAppDeveloperService) GetDeveloperAppList(db *gorm.DB, page, limit int, status *int, categoryId *int, developerId uint, keyword string) ([]models.MicroApp, int64, error) {
	m := models.MicroApp{}
	opts := models.MicroAppQueryOptions{
		Page:       page,
		Limit:      limit,
		Status:     status,
		CategoryId: categoryId,
		AuthorId:   &developerId,
		KeyWord:    keyword,
	}
	return m.GetListWithAllLangs(db, opts)
}

// GetDeveloperAppInfo 获取开发者应用详情（含权限验证）
func (s *MicroAppDeveloperService) GetDeveloperAppInfo(db *gorm.DB, appId uint, developerId uint) (map[string]interface{}, error) {
	m := models.MicroApp{}
	err := db.Preload("LangList").Where("id = ? AND author_id = ?", appId, developerId).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewBizError(ErrCodeAppNotFound)
		}
		return nil, err
	}

	// 查询作者名字
	authorName := ""
	if m.AuthorId > 0 {
		var user models.User
		if err := db.Select("name").Where("id = ?", m.AuthorId).First(&user).Error; err == nil {
			authorName = user.Name
		}
	}

	result := map[string]interface{}{
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

	return result, nil
}

// CreateApp 创建微应用
func (s *MicroAppDeveloperService) CreateApp(db *gorm.DB, microAppId, appName, appIcon, appDesc, remark string, categoryId int, chargeType int, price float64, developerId uint, screenshots string, langMap map[string]map[string]string) (map[string]interface{}, error) {
	// 自动生成 MicroAppId
	if microAppId == "" {
		newId, err := generateMicroAppId()
		if err != nil {
			return nil, err
		}
		microAppId = newId
	}

	// 检查 MicroAppId 是否已存在
	var existing models.MicroApp
	if err := db.Where("micro_app_id = ?", microAppId).First(&existing).Error; err == nil {
		return nil, NewBizError(ErrCodeAppIdExists)
	}

	m := models.MicroApp{
		MicroAppId: microAppId,
		AppName:    appName,
		AppIcon:    appIcon,
		AppDesc:    appDesc,
		Remark:     remark,
		CategoryId: categoryId,
		ChargeType: chargeType,
		Price:      price,
		AuthorId:   developerId,
		Status:     2, // 默认审核中
		Screenshots: screenshots,
	}

	// 事务保存主应用和多语言信息
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		for lang, langInfo := range langMap {
			if langInfo["appName"] != "" || langInfo["appDesc"] != "" {
				langModel := models.MicroAppLang{
					MicroAppId: m.MicroAppId,
					Lang:       lang,
					AppName:    langInfo["appName"],
					AppDesc:    langInfo["appDesc"],
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

// SubmitAppUpdate 提交应用信息更新（修改即提交审核）
func (s *MicroAppDeveloperService) SubmitAppUpdate(db *gorm.DB, id, developerId uint, appName, appIcon, appDesc, remark string, categoryId int, chargeType int, price float64, screenshots string, langMap map[string]map[string]string) (uint, error) {
	// 获取并验证权限
	app, err := s.getAppAndCheckPermission(db, id, developerId)
	if err != nil {
		return 0, err
	}

	// 检查是否已有待审核记录
	if err := s.checkNoPendingReview(db, id); err != nil {
		return 0, err
	}

	// 合并多语言信息
	mergedLangMap, err := s.mergeLangMap(db, app.MicroAppId, langMap)
	if err != nil {
		return 0, err
	}
	langMapJson, _ := json.Marshal(mergedLangMap)

	review := models.MicroAppReview{
		AppId:      id,
		AppName:    appName,
		AppIcon:    appIcon,
		AppDesc:    appDesc,
		CategoryId: categoryId,
		ChargeType: chargeType,
		Price:      price,
		Screenshots: screenshots,
		LangMap:    string(langMapJson),
		Remark:     remark,
		Status:     0, // 待审核
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", id).Updates(map[string]interface{}{
			"review_status": 1, // 审核中
			"review_id":     review.ID,
			"review_time":   &now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return review.ID, nil
}

// SubmitLangUpdate 提交语言信息更新（修改即提交审核）
func (s *MicroAppDeveloperService) SubmitLangUpdate(db *gorm.DB, id, developerId uint, langMap map[string]map[string]string) (uint, error) {
	// 获取并验证权限
	app, err := s.getAppAndCheckPermission(db, id, developerId)
	if err != nil {
		return 0, err
	}

	// 检查是否已有待审核记录
	if err := s.checkNoPendingReview(db, id); err != nil {
		return 0, err
	}

	// 合并多语言信息
	mergedLangMap, err := s.mergeLangMap(db, app.MicroAppId, langMap)
	if err != nil {
		return 0, err
	}
	langMapJson, _ := json.Marshal(mergedLangMap)

	review := models.MicroAppReview{
		AppId:       id,
		AppName:     app.AppName,
		AppIcon:     app.AppIcon,
		AppDesc:     app.AppDesc,
		CategoryId:  app.CategoryId,
		ChargeType:  app.ChargeType,
		Price:       app.Price,
		Screenshots: app.Screenshots,
		LangMap:     string(langMapJson),
		Remark:      app.Remark,
		Status:      0, // 待审核
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&models.MicroApp{}).Where("id = ?", id).Updates(map[string]interface{}{
			"review_status": 1,
			"review_id":     review.ID,
			"review_time":   &now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return review.ID, nil
}

// CancelAppReview 撤销应用审核
func (s *MicroAppDeveloperService) CancelAppReview(db *gorm.DB, id, developerId uint) error {
	// 获取并验证权限
	app, err := s.getAppAndCheckPermission(db, id, developerId)
	if err != nil {
		return err
	}

	// 检查是否有待审核记录
	if app.ReviewId == 0 {
		return NewBizError(ErrCodeNoPendingReviewApp)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", app.ReviewId).Delete(&models.MicroAppReview{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.MicroApp{}).Where("id = ?", id).Updates(map[string]interface{}{
			"review_status": 0,
			"review_id":     0,
			"review_time":   nil,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
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

	if app.AuthorId != developerId {
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

// mergeLangMap 合并现有多语言信息和新的多语言信息
func (s *MicroAppDeveloperService) mergeLangMap(db *gorm.DB, microAppId string, newLangMap map[string]map[string]string) (map[string]map[string]string, error) {
	langList, err := (&models.MicroAppLang{}).GetListByAppId(db, microAppId)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]map[string]string)
	for _, lang := range langList {
		merged[lang.Lang] = map[string]string{
			"appName": lang.AppName,
			"appDesc": lang.AppDesc,
		}
	}

	for lang, langInfo := range newLangMap {
		merged[lang] = langInfo
	}

	return merged, nil
}
