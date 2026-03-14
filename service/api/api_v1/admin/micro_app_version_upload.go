package admin

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
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"time"

	"github.com/gin-gonic/gin"
)

// MicroAppVersionUploadApi 微应用版本上传 API
type MicroAppVersionUploadApi struct{}

// Upload 上传微应用版本包
func (a *MicroAppVersionUploadApi) Upload(c *gin.Context) {
	// 获取上传的文件
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.Error(c, "请选择要上传的文件")
		return
	}

	// 检查文件扩展名
	fileExt := strings.ToLower(path.Ext(f.Filename))
	if fileExt != ".zip" {
		apiReturn.Error(c, "只支持 .zip 格式的文件")
		return
	}

	// 获取保存路径
	configUpload := global.Config.GetValueString("base", "micro_app_source_path")
	if configUpload == "" {
		// 如果没有配置，使用默认路径
		configUpload = "./micro_app_upload"
	}

	// 创建日期目录
	dateDir := fmt.Sprintf("%s/%d/%02d/%02d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(dateDir)
	if !isExist {
		os.MkdirAll(dateDir, os.ModePerm)
	}

	// 生成唯一文件名
	originalName := strings.TrimSuffix(f.Filename, path.Ext(f.Filename))
	hashStr := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fileName := fmt.Sprintf("%s-%s%s", originalName, hashStr[:12], fileExt)
	filePath := dateDir + fileName

	// 保存文件
	if err := c.SaveUploadedFile(f, filePath); err != nil {
		apiReturn.Error(c, "文件保存失败: "+err.Error())
		return
	}

	// 计算文件 MD5 校验值
	fileHash, err := calculateFileMD5(filePath)
	if err != nil {
		apiReturn.Error(c, "计算文件校验值失败")
		return
	}

	// 创建临时解压目录
	tempDir := filepath.Join(os.TempDir(), "micro_app_extract", hashStr[:16])
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		apiReturn.Error(c, "创建临时目录失败")
		return
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 解压文件
	if err := unzipFile(filePath, tempDir); err != nil {
		apiReturn.Error(c, "解压文件失败: "+err.Error())
		return
	}

	// 查找并解析配置文件 (app.config.js 或 app.config.json)
	config := parseAppConfig(tempDir)

	// 提取图标并保存到静态资源目录
	iconURL := extractAndSaveIcon(tempDir, config, fileName)

	// ==================== 机器审核预留 ====================
	// TODO: 在这里触发机器审核
	// 可以使用 goroutine 异步执行：
	// go machineReview(tempDir, config)
	// ====================

	// 返回结果
	relativePath := filePath[len(configUpload)-1:]
	// 提取文件夹名（不含路径）
	folderName := fileName
	apiReturn.SuccessData(c, MicroAppVersionUploadResp{
		URL:        relativePath,
		Hash:       fileHash,
		Config:     config,
		FileName:   fileName,
		FileSize:   f.Size,
		FolderName: folderName,
		IconURL:    iconURL,
	})
}

// calculateFileMD5 计算文件的 MD5 校验值
func calculateFileMD5(filepath string) (string, error) {
	file, err := os.Open(filepath)
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
func unzipFile(zipPath, destDir string) error {
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
		targetFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		if _, err := io.Copy(targetFile, reader); err != nil {
			return err
		}
	}

	return nil
}

// parseAppConfig 解析应用配置文件
func parseAppConfig(tempDir string) *MicroAppConfig {
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
			// 尝试解析
			if config := parseConfigFile(configPath); config != nil {
				return config
			}
		}
	}

	// 尝试在子目录中查找
	filepath.Walk(tempDir, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), "app.config.js") || strings.HasSuffix(info.Name(), "app.config.json") || strings.HasSuffix(info.Name(), ".config.js")) {
			if config := parseConfigFile(walkPath); config != nil {
				return filepath.SkipDir // 找到后停止遍历
			}
		}
		return nil
	})

	return nil
}

// parseConfigFile 解析配置文件
func parseConfigFile(configPath string) *MicroAppConfig {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	// 尝试解析 JSON
	var config MicroAppConfig
	strContent := string(content)

	// 尝试移除 JavaScript 导出语法 (export default 或 module.exports)
	strContent = strings.TrimPrefix(strContent, "export default ")
	strContent = strings.TrimPrefix(strContent, "module.exports = ")
	strContent = strings.TrimSpace(strContent)

	// 如果是 JSONP 格式，提取 JSON 部分
	if strings.HasPrefix(strContent, "(") && strings.HasSuffix(strContent, ")") {
		strContent = strContent[1 : len(strContent)-1]
	}

	if err := json.Unmarshal([]byte(strContent), &config); err != nil {
		// 尝试解析简化的 JSON 格式（用户提供的那种）
		if tryParseSimpleJSON(strContent, &config) {
			return &config
		}
		return nil
	}

	return &config
}

// tryParseSimpleJSON 尝试解析简单 JSON 格式
func tryParseSimpleJSON(content string, config *MicroAppConfig) bool {
	// 移除可能的 JavaScript 注释
	content = removeJSComments(content)

	// 尝试直接解析
	if err := json.Unmarshal([]byte(content), config); err != nil {
		return false
	}

	// 检查必要字段
	if config.MicroAppId == "" || config.Version == "" {
		return false
	}

	return true
}

// removeJSComments 移除 JavaScript 注释
func removeJSComments(content string) string {
	// 移除单行注释 //
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
				// 跳过到行尾
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
func extractAndSaveIcon(tempDir string, config *MicroAppConfig, versionFileName string) string {
	if config == nil || config.Icon == "" {
		return ""
	}

	// 获取图标文件名
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

