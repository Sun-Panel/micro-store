package biz

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sun-panel/global"
	"sun-panel/models/datatype"

	"github.com/golang-jwt/jwt/v4"
)

// 微应用审核
type microAppAudit struct {
}
type SecurityAuditResult datatype.SecurityAuditReport

// SecurityAuditConfig 安全审核配置
type SecurityAuditConfig struct {
	PlatformURL     string        `json:"platformUrl"`     // 三方平台API地址
	APISecret       string        `json:"apiSecret"`       // API密钥
	Timeout         time.Duration `json:"timeout"`         // 请求超时时间
	MaxFileSize     int64         `json:"maxFileSize"`     // 最大文件大小（字节）
	AllowedFileExts []string      `json:"allowedFileExts"` // 允许扫描的文件扩展名
}

// 默认的安全审核配置
func (m *microAppAudit) getDefaultConfig() SecurityAuditConfig {
	return SecurityAuditConfig{
		PlatformURL: global.Config.GetValueString("security_audit", "platform_url"),
		APISecret:   global.Config.GetValueString("security_audit", "api_secret"),
		Timeout:     30 * time.Second,
		MaxFileSize: 10 * 1024 * 1024, // 10MB
		AllowedFileExts: []string{
			".js", ".ts", ".jsx", ".tsx",
			".vue", ".html", ".css",
			".json",
		},
	}
}

// 基本检查，包含微应用名称和作者Id是否匹配
func (m *microAppAudit) BasicCheck(microAppDir string) error {
	return nil
}

// 代码安全审核（基于三方平台）
//   - 推荐异步调用本函数来实现
//   - 参数：microAppDir - 微应用包的目录路径
//   - 返回：审核结果和错误
func (m *microAppAudit) CodeSecurityAudit(microAppDir string, config SecurityAuditConfig) (*SecurityAuditResult, error) {
	// 1. 验证目录是否存在
	if _, err := os.Stat(microAppDir); os.IsNotExist(err) {
		return &SecurityAuditResult{
			IsPassed: false,
			Error:    fmt.Sprintf("微应用目录不存在: %s", microAppDir),
		}, fmt.Errorf("微应用目录不存在: %w", err)
	}

	// 2. 扫描代码文件
	codeFiles, err := m.scanCodeFiles(microAppDir, config)
	if err != nil {
		return &SecurityAuditResult{
			IsPassed: false,
			Error:    fmt.Sprintf("扫描代码文件失败: %v", err),
		}, fmt.Errorf("扫描代码文件失败: %w", err)
	}

	if len(codeFiles) == 0 {
		return &SecurityAuditResult{
			IsPassed: false,
			Error:    "未找到可扫描的代码文件",
		}, fmt.Errorf("未找到可扫描的代码文件")
	}

	// 3. 调用三方平台 API 进行安全审核
	result, err := m.callSecurityAuditAPI(codeFiles, config)
	if err != nil {
		return &SecurityAuditResult{
			IsPassed: false,
			Error:    fmt.Sprintf("调用安全审核API失败: %v", err),
		}, fmt.Errorf("调用安全审核API失败: %w", err)
	}

	result.ScanTime = time.Now()

	// 4. 判断审核结果
	result.IsPassed = m.determineAuditResult(result)

	return result, nil
}

// 代码安全审核（基于默认配置）
//   - 推荐异步调用本函数来实现
//   - 参数：microAppDir - 微应用包的目录路径
//   - 返回：审核结果和错误
func (m *microAppAudit) CodeSecurityAuditByDefaultConfig(microAppDir string) (*SecurityAuditResult, error) {
	config := m.getDefaultConfig()
	return m.CodeSecurityAudit(microAppDir, config)
}

// scanCodeFiles 扫描指定目录下的代码文件
func (m *microAppAudit) scanCodeFiles(dir string, config SecurityAuditConfig) (map[string]string, error) {
	codeFiles := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和隐藏文件
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		allowed := false
		for _, allowedExt := range config.AllowedFileExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}

		if !allowed {
			return nil
		}

		// 检查文件大小
		if info.Size() > config.MaxFileSize {
			return fmt.Errorf("文件过大: %s (%d bytes)", path, info.Size())
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败: %s: %w", path, err)
		}

		// 使用相对路径作为键
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}

		codeFiles[relPath] = string(content)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("扫描目录失败: %w", err)
	}

	return codeFiles, nil
}

// callSecurityAuditAPI 调用三方平台的安全审核 API
// 使用 sun-api 客户端调用 sun-api 服务端的安全审核接口
func (m *microAppAudit) callSecurityAuditAPI(codeFiles map[string]string, config SecurityAuditConfig) (*SecurityAuditResult, error) {
	// 如果没有配置三方平台，则使用本地的静态代码分析
	if config.PlatformURL == "" {
		global.Logger.Debug("未配置三方平台API，使用本地静态代码分析")
		return m.localStaticAnalysis(codeFiles)
	}

	// 调试日志：打印配置信息
	global.Logger.Infof("开始调用安全审核API - PlatformURL: %s, Timeout: %v, APIKey长度: %d",
		config.PlatformURL, config.Timeout)
	global.Logger.Infof("代码文件数量: %d", len(codeFiles))
	for filename := range codeFiles {
		global.Logger.Infof("  - 文件: %s (大小: %d bytes)", filename, len(codeFiles[filename]))
	}

	// 生成 JWT token（使用 APISecret）
	claims := jwt.MapClaims{
		"sub": "micro-store",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.APISecret))
	if err != nil {
		global.Logger.Errorf("生成JWT token失败: %v", err)
		return nil, fmt.Errorf("生成JWT token失败: %w", err)
	}

	// 构造请求数据（根据 API 文档格式）
	files := []map[string]interface{}{}
	for filename, content := range codeFiles {
		files = append(files, map[string]interface{}{
			"filename": filename,
			"content":  base64.StdEncoding.EncodeToString([]byte(content)), // Base64 编码
		})
	}

	requestData := map[string]interface{}{
		"plugin_id": fmt.Sprintf("plugin_%d", time.Now().Unix()),
		"files":     files,
		"config": map[string]interface{}{
			"strict_mode": false,
		},
	}

	// 调试日志：打印请求数据（截断）
	requestJSON, _ := json.Marshal(requestData)
	if len(requestJSON) > 500 {
		global.Logger.Infof("请求数据: %s...(总长度: %d)", string(requestJSON[:500]), len(requestJSON))
	} else {
		global.Logger.Infof("请求数据: %s", string(requestJSON))
	}

	// 构造响应数据结构（API 的实际返回格式）
	type APILocation struct {
		File   string `json:"file"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	}

	type APIIssue struct {
		Type        string      `json:"type"`
		Severity    string      `json:"severity"`
		Location    APILocation `json:"location"`
		Description string      `json:"description"`
		CodeSnippet string      `json:"code_snippet"`
		Suggestion  string      `json:"suggestion"`
	}

	type APIAuditResponse struct {
		AuditID   string     `json:"audit_id"`
		Status    string     `json:"status"`
		RiskScore int        `json:"risk_score"`
		Issues    []APIIssue `json:"issues"`
		CreatedAt string     `json:"created_at"`
	}

	// 创建 HTTP 请求
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	url := config.PlatformURL + "/api/v1/audit"
	global.Logger.Infof("准备调用API - 完整URL: %s", url)

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// 发送请求
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}
	startTime := time.Now()
	resp, err := httpClient.Do(req)
	elapsed := time.Since(startTime)

	global.Logger.Infof("API调用耗时: %v", elapsed)

	if err != nil {
		global.Logger.Errorf("调用安全审核API失败: %v (耗时: %v)", err, elapsed)
		return nil, fmt.Errorf("调用安全审核API失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != 200 {
		global.Logger.Errorf("安全审核API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("安全审核API返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var apiResponse APIAuditResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 转换为我们的 SecurityAuditResult 格式
	result := &SecurityAuditResult{
		ScanID:          apiResponse.AuditID,
		Score:           apiResponse.RiskScore,
		HighRiskCount:   0,
		MediumRiskCount: 0,
		LowRiskCount:    0,
		Vulnerabilities: []datatype.VulnerabilityInfo{},
		ScanTime:        time.Now(),
	}

	// 根据 status 判断是否通过
	if apiResponse.Status == "PASS" {
		result.IsPassed = true
	} else {
		result.IsPassed = false
	}

	// 转换 issues 为 vulnerabilities
	for _, issue := range apiResponse.Issues {
		vuln := datatype.VulnerabilityInfo{
			ID:          issue.Type,
			Title:       issue.Type,
			Description: issue.Description,
			Severity:    issue.Severity,
			Location:    issue.Location.File,
			LineNumber:  issue.Location.Line,
			Remediation: issue.Suggestion,
		}
		result.Vulnerabilities = append(result.Vulnerabilities, vuln)

		// 统计风险数量
		switch issue.Severity {
		case "CRITICAL", "HIGH":
			result.HighRiskCount++
		case "MEDIUM":
			result.MediumRiskCount++
		case "LOW":
			result.LowRiskCount++
		}
	}

	global.Logger.Infof("审核完成 - ID: %s, 状态: %s, 分数: %d, 问题数: %d",
		result.ScanID, apiResponse.Status, result.Score, len(result.Vulnerabilities))

	return result, nil
}

// localStaticAnalysis 本地静态代码分析（备用方案）
// 当没有配置三方平台或三方平台调用失败时使用
func (m *microAppAudit) localStaticAnalysis(codeFiles map[string]string) (*SecurityAuditResult, error) {
	result := &SecurityAuditResult{
		ScanID:          fmt.Sprintf("local_%d", time.Now().Unix()),
		Score:           100,
		HighRiskCount:   0,
		MediumRiskCount: 0,
		LowRiskCount:    0,
		Vulnerabilities: []datatype.VulnerabilityInfo{},
	}

	// 检查常见的安全问题
	for filePath, content := range codeFiles {
		// 检查硬编码的密钥和密码
		if m.detectHardcodedSecrets(content) {
			result.Vulnerabilities = append(result.Vulnerabilities, datatype.VulnerabilityInfo{
				ID:          "HARD_CODED_SECRET",
				Title:       "检测到硬编码的密钥或密码",
				Description: "代码中包含疑似硬编码的密钥或密码，建议使用环境变量或配置文件管理敏感信息",
				Severity:    "high",
				Location:    filePath,
				Remediation: "将敏感信息移到环境变量或配置文件中",
			})
			result.HighRiskCount++
			result.Score -= 20
		}

		// 检查 eval 等危险函数的使用
		if m.detectDangerousFunctions(content) {
			result.Vulnerabilities = append(result.Vulnerabilities, datatype.VulnerabilityInfo{
				ID:          "DANGEROUS_FUNCTION",
				Title:       "检测到危险函数的使用",
				Description: "代码中使用了 eval 等危险函数，可能导致代码注入攻击",
				Severity:    "high",
				Location:    filePath,
				Remediation: "避免使用 eval 等危险函数，使用更安全的替代方案",
			})
			result.HighRiskCount++
			result.Score -= 15
		}

		// 检查 console.log 语句（开发时遗留）
		if strings.Contains(content, "console.log") {
			result.Vulnerabilities = append(result.Vulnerabilities, datatype.VulnerabilityInfo{
				ID:          "CONSOLE_LOG",
				Title:       "检测到 console.log 语句",
				Description: "代码中包含 console.log 调试语句，生产环境中应该移除",
				Severity:    "low",
				Location:    filePath,
				Remediation: "在发布到生产环境前移除所有调试语句",
			})
			result.LowRiskCount++
			result.Score -= 5
		}

		// 检查内联脚本（XSS 风险）
		if strings.Contains(content, "innerHTML") || strings.Contains(content, "document.write") {
			result.Vulnerabilities = append(result.Vulnerabilities, datatype.VulnerabilityInfo{
				ID:          "XSS_RISK",
				Title:       "潜在的 XSS 漏洞",
				Description: "使用了 innerHTML 或 document.write，可能导致跨站脚本攻击",
				Severity:    "medium",
				Location:    filePath,
				Remediation: "使用 textContent 或创建 DOM 节点来替代 innerHTML，对用户输入进行充分的转义和验证",
			})
			result.MediumRiskCount++
			result.Score -= 10
		}
	}

	// 确保分数在 0-100 范围内
	if result.Score < 0 {
		result.Score = 0
	}

	return result, nil
}

// detectHardcodedSecrets 检测硬编码的密钥和密码
func (m *microAppAudit) detectHardcodedSecrets(content string) bool {
	secretPatterns := []string{
		`password\s*[:=]\s*["\'][^"\']+["\']`,
		`api[_-]?key\s*[:=]\s*["\'][^"\']+["\']`,
		`secret\s*[:=]\s*["\'][^"\']+["\']`,
		`token\s*[:=]\s*["\'][^"\']+["\']`,
		`private[_-]?key\s*[:=]\s*["\'][^"\']+["\']`,
	}

	for _, pattern := range secretPatterns {
		if strings.Contains(strings.ToLower(content), strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// detectDangerousFunctions 检测危险函数的使用
func (m *microAppAudit) detectDangerousFunctions(content string) bool {
	dangerousFuncs := []string{
		"eval(",
		"Function(",
		"setTimeout(",
		"setInterval(",
	}

	for _, funcName := range dangerousFuncs {
		if strings.Contains(content, funcName) {
			return true
		}
	}

	return false
}

// determineAuditResult 判断审核是否通过
func (m *microAppAudit) determineAuditResult(result *SecurityAuditResult) bool {
	// 如果没有配置三方平台，则使用本地分析的结果
	if result.ScanID == "" {
		result.ScanID = fmt.Sprintf("local_%d", time.Now().Unix())
	}

	// 审核通过标准：
	// 1. 分数 >= 60
	// 2. 高风险漏洞数量 = 0
	// 3. 中风险漏洞数量 <= 2
	if result.Score >= 60 && result.HighRiskCount == 0 && result.MediumRiskCount <= 2 {
		return true
	}

	return false
}

// // ExtractMicroAppPackage 解压微应用包到指定目录
// func (m *microAppAudit) ExtractMicroAppPackage(zipPath string, extractDir string) error {
// 	// 创建解压目录
// 	if err := os.MkdirAll(extractDir, 0755); err != nil {
// 		return fmt.Errorf("创建解压目录失败: %w", err)
// 	}

// 	// 打开 zip 文件
// 	zipReader, err := zip.OpenReader(zipPath)
// 	if err != nil {
// 		return fmt.Errorf("打开zip文件失败: %w", err)
// 	}
// 	defer zipReader.Close()

// 	// 解压文件
// 	for _, file := range zipReader.File {
// 		// 构建目标文件路径
// 		destPath := filepath.Join(extractDir, file.Name)

// 		// 安全检查：防止路径遍历攻击
// 		if !strings.HasPrefix(destPath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
// 			return fmt.Errorf("非法的文件路径: %s", file.Name)
// 		}

// 		// 如果是目录，则创建
// 		if file.FileInfo().IsDir() {
// 			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
// 				return fmt.Errorf("创建目录失败: %w", err)
// 			}
// 			continue
// 		}

// 		// 确保父目录存在
// 		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
// 			return fmt.Errorf("创建父目录失败: %w", err)
// 		}

// 		// 解压文件
// 		srcFile, err := file.Open()
// 		if err != nil {
// 			return fmt.Errorf("打开zip中的文件失败: %w", err)
// 		}
// 		defer srcFile.Close()

// 		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
// 		if err != nil {
// 			return fmt.Errorf("创建目标文件失败: %w", err)
// 		}
// 		defer destFile.Close()

// 		if _, err := io.Copy(destFile, srcFile); err != nil {
// 			return fmt.Errorf("复制文件内容失败: %w", err)
// 		}
// 	}

// 	return nil
// }

// // GetPackageHash 计算包的哈希值（用于缓存和去重）
// func (m *microAppAudit) GetPackageHash(zipPath string) (string, error) {
// 	file, err := os.Open(zipPath)
// 	if err != nil {
// 		return "", fmt.Errorf("打开文件失败: %w", err)
// 	}
// 	defer file.Close()

// 	hash := md5.New()
// 	if _, err := io.Copy(hash, file); err != nil {
// 		return "", fmt.Errorf("计算哈希值失败: %w", err)
// 	}

// 	return fmt.Sprintf("%x", hash.Sum(nil)), nil
// }

// GetSecurityAuditConfig 获取当前的安全审核配置
func (m *microAppAudit) GetSecurityAuditConfig() SecurityAuditConfig {
	return m.getDefaultConfig()
}
