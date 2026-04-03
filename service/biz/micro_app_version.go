package biz

import (
	"sun-panel/models"

	"gorm.io/gorm"
)

// MicroAppVersionService 微应用版本业务服务
type MicroAppVersionService struct{}

// CreateWithCheck 创建版本（包含业务检查）
func (s *MicroAppVersionService) CreateWithCheck(db *gorm.DB, version *models.MicroAppVersion) error {
	// 1. 检查应用是否存在并获取 microAppId
	app := models.MicroApp{}
	app, err := app.GetById(db, version.AppRecordId)
	if err != nil {
		return NewBizError(ErrCodeAppNotFound)
	}

	// 2. 检查版本号是否存在
	m := models.MicroAppVersion{}
	exists, err := m.CheckVersionExist(db, version.AppRecordId, version.Version, 0)
	if err != nil {
		return err // 数据库错误，直接返回
	}
	if exists {
		return NewBizError(ErrCodeVersionExists)
	}

	// // 3. 检查版本号数字是否存在
	// exists, err = m.CheckVersionCodeExist(db, version.AppRecordId, version.VersionCode, 0)
	// if err != nil {
	// 	return err // 数据库错误，直接返回
	// }
	// if exists {
	// 	return NewBizError(ErrCodeVersionCodeExists)
	// }

	// 4. 设置 AppId
	version.AppRecordId = app.ID

	// 5. 创建版本
	version.Status = -1 // 默认草稿状态
	if err := version.Create(db); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

// GetPendingListWithAppInfo 获取待审核版本列表（包含应用信息）
func (s *MicroAppVersionService) GetPendingListWithAppInfo(db *gorm.DB, page, limit int) ([]map[string]interface{}, int64, error) {
	m := models.MicroAppVersion{}
	list, total, err := m.GetPendingList(db, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// 收集所有 AppRecordId
	appIds := make([]uint, 0, len(list))
	for _, v := range list {
		appIds = append(appIds, v.AppRecordId)
	}

	// 批量查询应用信息（避免 N+1）
	var apps []models.MicroApp
	appMap := make(map[uint]models.MicroApp)
	if len(appIds) > 0 {
		if err := db.Where("id IN ?", appIds).Find(&apps).Error; err == nil {
			for _, app := range apps {
				appMap[app.ID] = app
			}
		}
	}

	// 组装结果
	result := make([]map[string]interface{}, len(list))
	for i, v := range list {
		appInfo := map[string]interface{}{
			"id":          v.ID,
			"appId":       v.AppRecordId,
			"version":     v.Version,
			"versionCode": v.VersionCode,
			"packageUrl":  v.PackageUrl,
			"status":      v.Status,
			"createTime":  v.CreatedAt,
			"reviewTime":  v.ReviewTime,
			"reviewNote":  v.ReviewNote,
		}
		if app, ok := appMap[v.AppRecordId]; ok {
			appInfo["appName"] = app.AppName
			appInfo["appIcon"] = app.AppIcon
		}
		result[i] = appInfo
	}

	return result, total, nil
}

// SubmitReview 提交审核
func (s *MicroAppVersionService) SubmitReview(db *gorm.DB, versionId uint) error {
	m := models.MicroAppVersion{}
	version, err := m.GetById(db, versionId)
	if err != nil {
		return NewBizError(ErrCodeVersionNotFound)
	}

	// 草稿(-1)、拒绝(2)、下架(3)状态可以提交审核
	if version.Status != -1 && version.Status != 2 && version.Status != 3 {
		return NewBizError(ErrCodeStatusNotAllowed)
	}

	if err := m.Update(db, versionId, map[string]interface{}{"status": 0}); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

// CancelReview 撤销审核
func (s *MicroAppVersionService) CancelReview(db *gorm.DB, versionId uint) error {
	m := models.MicroAppVersion{}
	version, err := m.GetById(db, versionId)
	if err != nil {
		return NewBizError(ErrCodeVersionNotFound)
	}

	if version.Status != 0 {
		return NewBizError(ErrCodeStatusNotAllowed)
	}

	if err := m.Update(db, versionId, map[string]interface{}{"status": -1}); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

// DeleteVersion 删除版本
func (s *MicroAppVersionService) DeleteVersion(db *gorm.DB, ids []uint) error {
	m := models.MicroAppVersion{}
	version, err := m.GetById(db, ids[0])
	if err != nil {
		return NewBizError(ErrCodeVersionNotFound)
	}

	if version.Status == 1 {
		return NewBizError(ErrCodeApprovedCannotDelete)
	}

	if err := m.Delete(db, ids); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

// UpdateVersion 更新版本
func (s *MicroAppVersionService) UpdateVersion(db *gorm.DB, id uint, version string, versionCode int) error {
	m := models.MicroAppVersion{}
	v, err := m.GetById(db, id)
	if err != nil {
		return NewBizError(ErrCodeVersionNotFound)
	}

	if v.Status != -1 {
		return NewBizError(ErrCodeStatusNotAllowed)
	}

	updateData := map[string]interface{}{}
	if version != "" {
		updateData["version"] = version
	}
	if versionCode > 0 {
		updateData["version_code"] = versionCode
	}

	if len(updateData) == 0 {
		return NewBizError(ErrCodeNoUpdateContent)
	}

	if err := m.Update(db, id, updateData); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

// Review 审核版本
func (s *MicroAppVersionService) Review(db *gorm.DB, versionId uint, status int, reviewerId uint, reviewNote string) error {
	m := models.MicroAppVersion{}
	version, err := m.GetById(db, versionId)
	if err != nil {
		return NewBizError(ErrCodeVersionNotFound)
	}

	if version.Status != 0 {
		return NewBizError(ErrCodeNotPendingReview)
	}

	if err := m.Review(db, versionId, status, reviewerId, reviewNote); err != nil {
		return err // 数据库错误，直接返回
	}

	return nil
}

func (s *MicroAppVersionService) GetLatestOnlineByAppModelId(db *gorm.DB, appModelId uint) (models.MicroAppVersion, error) {
	m := models.MicroAppVersion{}
	return m.GetLatestOnlineByAppId(db, appModelId)
}

// 获取指定版本信息
func (s *MicroAppVersionService) GetInfoOnLineByVersion(db *gorm.DB, version string) (models.MicroAppVersion, error) {
	m := models.MicroAppVersion{}
	info, err := m.GetByVersion(db, version)
	if err != nil {
		return models.MicroAppVersion{}, err
	}

	if info.Status != 1 || info.OfflineType != 0 {
		return models.MicroAppVersion{}, NewBizError(ErrCodeVersionNotFound)
	}

	return info, nil
}
