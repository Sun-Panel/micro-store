package biz

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"
)

// MicroAppPackageResult 微应用包处理结果
type MicroAppPackageResult struct {
	Src          string                       `json:"src"`          // 文件源路径
	Hash         string                       `json:"hash"`         // 文件 MD5
	Config       models.MicroAppVersionConfig `json:"config"`       // 解析的配置
	FileName     string                       `json:"fileName"`     // 文件名
	FileSize     int64                        `json:"fileSize"`     // 文件大小
	FullFilePath string                       `json:"fullFilePath"` // 完整的文件路径
	IconURL      string                       `json:"iconURL"`      // 图标 URL
}

type MicroAppPackageUploadCache struct {
	PackageResult MicroAppPackageResult `json:"microAppPackageResult"`
	AppRecordId   uint                  `json:"appRecordId"`
}

// MicroAppCfg 微应用配置
type MicroAppCfg struct {
	AppJsonVersion string                   `json:"appJsonVersion"`
	MicroAppId     string                   `json:"microAppId"`
	Version        string                   `json:"version"`
	APIVersion     string                   `json:"apiVersion"`
	Author         string                   `json:"author"`
	Entry          string                   `json:"entry"`
	Icon           string                   `json:"icon"`
	Debug          bool                     `json:"debug"`
	Components     map[string]interface{}   `json:"components"`
	Permissions    []string                 `json:"permissions"`
	DataNodes      map[string]interface{}   `json:"dataNodes"`
	NetworkDomains []string                 `json:"networkDomains"`
	AppInfo        map[string]AppInfoConfig `json:"appInfo"`
}

// AppInfoConfig 应用信息配置
type AppInfoConfig struct {
	AppName            string `json:"appName"`
	Description        string `json:"description"`
	NetworkDescription string `json:"networkDescription"`
}

// MicroAppPackageService 微应用包处理服务
type MicroAppPackageService struct {
	UploadCache cache.Cacher[MicroAppPackageUploadCache]
}

func (s *MicroAppPackageService) Init() {
	s.UploadCache = global.NewCache[MicroAppPackageUploadCache](24*time.Hour, 12*time.Hour, "micro_app_package_result")
}

func (s *MicroAppPackageService) uploadCacheKey(appRecordId uint, version string) string {
	return fmt.Sprintf("%d_%s_%s", appRecordId, version, cmn.BuildRandCode(12, cmn.RAND_CODE_MODE1))
}

func (s *MicroAppPackageService) SetUploadCache(appRecordId uint, version string, result MicroAppPackageUploadCache) string {
	key := s.uploadCacheKey(appRecordId, version)
	s.UploadCache.SetDefault(key, result)
	return key
}

func (s *MicroAppPackageService) GetUploadCache(key string) (MicroAppPackageUploadCache, bool) {
	return s.UploadCache.Get(key)
}

func (s *MicroAppPackageService) DelUploadCache(key string) {
	s.UploadCache.Delete(key)
}

// UploadMicroAppPackage 上传并处理微应用包
func (s *MicroAppPackageService) UploadMicroAppPackage(fileData []byte, fileName string) (MicroAppPackageResult, error) {
	// 获取保存路径
	configUpload := s.getSavePath()

	// 创建日期目录部分
	dateDir := fmt.Sprintf("%d/%02d/%02d/", time.Now().Year(), time.Now().Month(), time.Now().Day())
	fullDir := fmt.Sprintf("%s/%s", configUpload, dateDir)
	isExist, _ := cmn.PathExists(fullDir)
	if !isExist {
		os.MkdirAll(fullDir, os.ModePerm)
	}

	// 先计算文件 MD5 校验值（用于生成文件名）
	md5Hash := md5.Sum(fileData)
	fileHash := hex.EncodeToString(md5Hash[:])

	// 先解压文件获取配置
	// 创建临时解压目录
	tempDir := filepath.Join(os.TempDir(), "micro_app_extract", fileHash[:16])
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return MicroAppPackageResult{}, fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 先保存临时文件用于解压
	tempFilePath := filepath.Join(tempDir, "temp"+strings.ToLower(path.Ext(fileName)))
	if err := os.WriteFile(tempFilePath, fileData, 0644); err != nil {
		return MicroAppPackageResult{}, fmt.Errorf("临时文件保存失败: %w", err)
	}

	// 解压文件获取配置
	if err := s.UnzipFile(tempFilePath, tempDir); err != nil {
		return MicroAppPackageResult{}, fmt.Errorf("解压文件失败: %w", err)
	}

	// 查找并解析配置文件
	config := s.parseAppConfig(tempDir)

	// 使用配置信息生成最终文件名：[microappid]_[version]_[hash前16位]
	fileExt := strings.ToLower(path.Ext(fileName))
	newFileName := fmt.Sprintf("%s_%s_%s%s", config.MicroAppId, config.Version, fileHash[:16], fileExt)
	filePath := dateDir + newFileName
	fullFilePath := fmt.Sprintf("%s/%s", configUpload, filePath) // 完整的保存路径

	// global.Logger.Info("fullFilePath", zap.String("fullFilePath", fullFilePath), zap.String("newFileName", newFileName), zap.String("filePath", filePath))
	// global.Logger.Info("appConfig", zap.Any("config", config))

	// 保存文件（如果存在则覆盖）
	if err := os.WriteFile(fullFilePath, fileData, 0644); err != nil {
		return MicroAppPackageResult{}, fmt.Errorf("文件保存失败: %w", err)
	}

	// 提取图标并保存到静态资源目录
	iconURL := s.extractAndSaveIcon(tempDir, config, newFileName)

	// downloadUrl := s.GenerateDownloadURL(filePath)

	return MicroAppPackageResult{
		Src:          filePath,
		Hash:         fileHash,
		Config:       config,
		FileName:     newFileName,
		FileSize:     int64(len(fileData)),
		FullFilePath: fullFilePath,
		IconURL:      iconURL,
	}, nil
}

// calculateFileMD5 计算文件的 MD5 校验值
func (s *MicroAppPackageService) calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// unzipFile 解压 zip 文件
func (s *MicroAppPackageService) UnzipFile(zipPath, destDir string) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destDir, file.Name)

		// 检查是否是目录
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// 解压文件
		if err := s.extractFile(file, filePath); err != nil {
			return err
		}
	}

	return nil
}

// extractFile 解压单个文件
func (s *MicroAppPackageService) extractFile(file *zip.File, destPath string) error {
	targetFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(targetFile, reader)
	return err
}

// parseAppConfig 解析应用配置文件
func (s *MicroAppPackageService) parseAppConfig(tempDir string) models.MicroAppVersionConfig {
	// 尝试查找 app.config.js 或 app.config.json
	configFiles := []string{
		"app.config.js",
		"app.config.json",
		"app.config",
		"app.json",
	}

	for _, configFile := range configFiles {
		configPath := filepath.Join(tempDir, configFile)
		if _, err := os.Stat(configPath); err == nil {
			if config := s.parseConfigFile(configPath); config != nil {
				return *config
			}
		}
	}

	// 尝试在子目录中查找
	var foundConfig models.MicroAppVersionConfig
	filepath.Walk(tempDir, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), "app.config.js") || strings.HasSuffix(info.Name(), "app.config.json") || strings.HasSuffix(info.Name(), ".config.js")) {
			if config := s.parseConfigFile(walkPath); config != nil {
				foundConfig = *config
				return filepath.SkipDir
			}
		}
		return nil
	})

	return foundConfig
}

// parseConfigFile 解析配置文件
func (s *MicroAppPackageService) parseConfigFile(configPath string) *models.MicroAppVersionConfig {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	strContent := string(content)

	// 尝试移除 JavaScript 导出语法
	strContent = strings.TrimPrefix(strContent, "export default ")
	strContent = strings.TrimPrefix(strContent, "module.exports = ")
	strContent = strings.TrimSpace(strContent)

	// 如果是 JSONP 格式，提取 JSON 部分
	if strings.HasPrefix(strContent, "(") && strings.HasSuffix(strContent, ")") {
		strContent = strContent[1 : len(strContent)-1]
	}

	// 移除 JavaScript 注释
	strContent = s.removeJSComments(strContent)

	var config models.MicroAppVersionConfig
	if err := json.Unmarshal([]byte(strContent), &config); err != nil {
		return nil
	}

	// 检查必要字段
	if config.MicroAppId == "" {
		return nil
	}
	if config.Version == "" {
		return nil
	}
	// if config.APIVersion == "" {
	// 	return nil
	// }

	return &config
}

// removeJSComments 移除 JavaScript 注释
func (s *MicroAppPackageService) removeJSComments(content string) string {
	var result []rune
	inString := false
	stringChar := ' '
	escaped := false

	for i, char := range content {
		if escaped {
			result = append(result, char)
			escaped = false
			continue
		}

		if char == '\\' {
			result = append(result, char)
			escaped = true
			continue
		}

		if char == '"' || char == '\'' {
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				inString = false
			}
		}

		if !inString {
			// 检查单行注释
			if i+1 < len(content) && char == '/' && content[i+1] == '/' {
				for j := i; j < len(content); j++ {
					if content[j] == '\n' {
						result = append(result, '\n')
						content = content[j:]
						break
					}
				}
				continue
			}
			// 检查多行注释
			if i+1 < len(content) && char == '/' && content[i+1] == '*' {
				j := i + 2
				for j < len(content)-1 {
					if content[j] == '*' && content[j+1] == '/' {
						j += 2
						break
					}
					j++
				}
				content = content[:i] + content[j:]
				continue
			}
		}

		result = append(result, char)
	}

	return string(result)
}

// extractAndSaveIcon 从解压目录中提取图标并保存到静态资源目录
func (s *MicroAppPackageService) extractAndSaveIcon(tempDir string, config models.MicroAppVersionConfig, versionFileName string) string {

	iconFileName := config.Icon

	// 尝试在多个可能的位置查找图标
	possiblePaths := []string{
		filepath.Join(tempDir, iconFileName),
		filepath.Join(tempDir, "public", iconFileName),
		filepath.Join(tempDir, "assets", iconFileName),
	}

	var iconSourcePath string
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			iconSourcePath = p
			break
		}
	}

	if iconSourcePath == "" {
		return ""
	}

	// 获取静态资源目录
	sourcePath := global.Config.GetValueString("base", "source_path")
	if sourcePath == "" {
		sourcePath = "./uploads"
	}

	// 创建图标保存目录
	iconDir := filepath.Join(sourcePath, "micro_app_icon")
	if err := os.MkdirAll(iconDir, os.ModePerm); err != nil {
		return ""
	}

	// 生成唯一的图标文件名
	iconExt := filepath.Ext(iconFileName)
	iconNewName := versionFileName + iconExt
	iconDestPath := filepath.Join(iconDir, iconNewName)

	// 复制文件
	sourceFile, err := os.Open(iconSourcePath)
	if err != nil {
		return ""
	}
	defer sourceFile.Close()

	destFile, err := os.Create(iconDestPath)
	if err != nil {
		return ""
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return ""
	}

	// 返回图标的访问 URL
	return "/uploads/micro_app_icon/" + iconNewName
}

// GenerateDownloadURL 生成完整的下载 URL
// 参数示例: 相对路径: "2026/03/30/ce3b0a0b94ca-yTJJM77I.zip"，
// 完整真实路径: "./micro_app_upload/2026/03/30/ce3b0a0b94ca-yTJJM77I.zip"，
// 固定的下载地址路径："/micro_app_uploads/2026/03/30/ce3b0a0b94ca-yTJJM77I.zip"
// 返回完整的浏览器可访问下载地址
// func (s *MicroAppPackageService) GenerateDownloadURL(relativeSrcPath string) string {
// 	// return strings.TrimPrefix(relativePath, ".")
// 	// savePath := s.getSavePath()
// 	// if relativePath == "" {
// 	// 	return ""
// 	// }
// 	return "/micro_app_uploads/" + strings.Trim(relativeSrcPath, "/")
// }

func (s *MicroAppPackageService) BuildDownloadUrl(microAppId, version string) string {
	if len(version) == 0 || version == "" {
		// 下载最新版本
		return fmt.Sprintf("/api/microApp/download/%s", microAppId)
	}
	// 下载指定版本
	return fmt.Sprintf("/api/microApp/download/%s/%s", microAppId, version)
}

func (s *MicroAppPackageService) getSavePath() string {
	// 获取保存路径
	configUpload := global.Config.GetValueString("base", "micro_app_source_path")
	if configUpload == "" {
		configUpload = "./micro_app_upload"
	}
	return configUpload
}
