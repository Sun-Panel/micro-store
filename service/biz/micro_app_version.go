package biz

import (
	"os"
	"path/filepath"
	"strings"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sun-panel/models/datatype"
	"sync"
	"time"

	"gorm.io/gorm"
)

// MicroAppVersionService 微应用版本业务服务
type MicroAppVersionService struct{}

// CreateWithCheck 创建版本（包含业务检查）
func (s *MicroAppVersionService) CreateOrUpdateWithCheck(db *gorm.DB, version *models.MicroAppVersion, packageFolderName string) error {
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

	// 4. 设置 AppId
	version.AppRecordId = app.ID

	// 5. 创建版本或更新
	version.Status = -1 // 默认草稿状态

	if version.ID == 0 {
		if err := version.Create(db); err != nil {
			return err // 数据库错误，直接返回
		}
	} else {
		if err := version.Update(db); err != nil {
			return err
		}
	}

	// 6. 解压 zip 文件到临时目录
	extractPath, err := s.extractZipToTemp(packageFolderName)
	if err != nil {
		global.Logger.Errorln("解压文件失败:", err)
		return err
	}

	// 7. 创建 WaitGroup 用于等待异步操作
	var wg sync.WaitGroup
	wg.Add(1) // 为异步操作添加一个计数

	// 8. 同步操作（如果有）
	// 这里可以添加需要同步执行的解压后操作
	// 例如：解析配置文件、验证文件结构等
	// err = s.syncOperation(db, version, extractPath)
	// if err != nil {
	// 	os.RemoveAll(extractPath)
	// 	return err
	// }

	// 9. 异步审核
	go func() {
		defer wg.Done() // 完成时通知 WaitGroup
		s.rebootAudit(db, version, extractPath)
	}()

	// 10. 启动清理协程，等待所有操作完成后清理临时目录
	go func() {
		wg.Wait()
		// os.RemoveAll(extractPath) // 删除临时目录
		global.Logger.Infoln("临时目录已清理:", extractPath)
	}()

	return nil
}

func (s *MicroAppVersionService) rebootAudit(db *gorm.DB, version *models.MicroAppVersion, extractDir string) {

	// 获取默认配置，确保包含允许的文件扩展名
	defaultConfig := MicroAppAudit.GetSecurityAuditConfig()

	// 覆盖自定义配置
	securityAuditResult, err := MicroAppAudit.CodeSecurityAudit(extractDir, SecurityAuditConfig{
		PlatformURL:     "http://127.0.0.1:3025",
		APISecret:       "hYWxxDCCcM5Ma8Mt3h2H0RemTn9bTG6Q",
		Timeout:         60 * time.Second,
		MaxFileSize:     defaultConfig.MaxFileSize,
		AllowedFileExts: []string{".js"},
	})
	if err != nil {
		global.Logger.Errorln("安全审核失败:", err)
		return
	}
	version.SecurityAuditReport = (*datatype.SecurityAuditReport)(securityAuditResult)
	if err := version.Update(db); err != nil {
		global.Logger.Errorln("更新安全审核结果失败:", err)
		return
	}
}

// extractZipToTemp 解压 zip 文件到临时目录
func (s *MicroAppVersionService) extractZipToTemp(zipPath string) (string, error) {
	// 计算文件的哈希值，用于创建唯一的临时目录名
	// 去掉文件后缀作为目录名
	uniqueName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))

	tempPath := Config.GetTempPath()

	// 创建临时解压目录
	tempDir := filepath.Join(tempPath, "micro_app_extract", uniqueName+"_"+cmn.BuildRandCode(10, cmn.RAND_CODE_MODE1))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}

	// 解压文件
	err := s.unzipFile(zipPath, tempDir)
	if err != nil {
		os.RemoveAll(tempDir) // 解压失败，清理临时目录
		return "", err
	}

	return tempDir, nil
}

// unzipFile 解压 zip 文件到指定目录
func (s *MicroAppVersionService) unzipFile(zipPath, destDir string) error {
	// 使用 MicroAppPackageService 的公开方法
	packageService := &MicroAppPackageService{}
	return packageService.UnzipFile(zipPath, destDir)
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

	version.Status = 0
	if err := version.Update(db); err != nil {
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

	version.Status = -1
	if err := version.Update(db); err != nil {
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

	if version == "" && versionCode <= 0 {
		return NewBizError(ErrCodeNoUpdateContent)
	}

	if version != "" {
		v.Version = version
	}
	if versionCode > 0 {
		v.VersionCode = versionCode
	}

	if err := v.Update(db); err != nil {
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

	if err := version.Review(db, status, reviewerId, reviewNote); err != nil {
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
