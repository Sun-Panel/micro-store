package compat

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// 系统变量名称常量（不带 "var." 前缀，SetVariable/GetVariableByInterface 会自动添加）
const (
	SystemVarAPIVersion     = "micro_app_compat_api_version"
	SystemVarAppJsonVersion = "micro_app_compat_app_json_version"
)

// VersionEntry 版本兼容性条目
// 语义：如果输入版本 >= Version，则最低主应用版本为 LowVersion
type VersionEntry struct {
	Version    string `json:"version"`
	LowVersion string `json:"lowVersion"`
}

// DefaultAPIVersionCompat 默认的 API 版本兼容性映射
var DefaultAPIVersionCompat = []VersionEntry{
	{Version: "1.0.18", LowVersion: "2.0.1"},
	{Version: "1.0.15", LowVersion: "1.8.0"},
}

// DefaultAppJsonVersionCompat 默认的 app.json 版本兼容性映射
var DefaultAppJsonVersionCompat = []VersionEntry{
	{Version: "1.1", LowVersion: "2.0.0"},
	{Version: "1.0", LowVersion: "1.8.0"},
}

// systemSettingReader 系统变量读取接口
type systemSettingReader interface {
	GetVariableByInterface(configName string, configValue interface{}) error
}

// systemSettingInstance 通过 Init 注入
var systemSettingInstance systemSettingReader

// Init 初始化兼容性模块，注入 SystemSetting 实例
// 在程序启动时调用：compat.Init(global.SystemSetting)
func Init(instance systemSettingReader) {
	systemSettingInstance = instance
	// 初始化时确保默认值的排序
	sortEntries(DefaultAPIVersionCompat)
	sortEntries(DefaultAppJsonVersionCompat)
}

// ResolveLowVersion 根据 apiVersion 和 appJsonVersion 推导最低主应用版本
// 两个约束是"且"关系，取较严格的值（max）
func ResolveLowVersion(apiVersion, appJsonVersion string) string {
	low := ""

	apiEntries := getAPIVersionEntries()
	appJsonEntries := getAppJsonVersionEntries()

	// 查 apiVersion 范围映射
	if v := lookupRange(apiEntries, apiVersion); v != "" {
		low = maxVersion(low, v)
	}

	// 查 appJsonVersion 范围映射
	if v := lookupRange(appJsonEntries, appJsonVersion); v != "" {
		low = maxVersion(low, v)
	}

	return low
}

// getAPIVersionEntries 获取 API 版本兼容性映射（优先系统变量，兜底默认值）
func getAPIVersionEntries() []VersionEntry {
	return getEntries(SystemVarAPIVersion, DefaultAPIVersionCompat)
}

// getAppJsonVersionEntries 获取 appJson 版本兼容性映射（优先系统变量，兜底默认值）
func getAppJsonVersionEntries() []VersionEntry {
	return getEntries(SystemVarAppJsonVersion, DefaultAppJsonVersionCompat)
}

// getEntries 从系统变量读取映射，失败时回退到默认值
func getEntries(systemVarName string, defaultEntries []VersionEntry) []VersionEntry {
	entries, err := readSystemVariableEntries(systemVarName)
	if err == nil && len(entries) > 0 {
		// 按 Version 降序排列（范围查询需要）
		sortEntries(entries)
		return entries
	}

	if err != nil {
		log.Printf("[compat] 读取系统变量 %s 失败: %v，使用默认值\n", systemVarName, err)
	}

	// 回退到默认值（拷贝一份，避免修改原始数据）
	sorted := make([]VersionEntry, len(defaultEntries))
	copy(sorted, defaultEntries)
	sortEntries(sorted)
	return sorted
}

// readSystemVariableEntries 从系统变量读取 JSON 数组
func readSystemVariableEntries(configName string) ([]VersionEntry, error) {
	if systemSettingInstance == nil {
		return nil, fmt.Errorf("SystemSetting not initialized, call compat.Init() first")
	}

	var entries []VersionEntry
	if err := systemSettingInstance.GetVariableByInterface(configName, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// InvalidateCache 清除兼容性映射缓存
// 在系统变量更新后调用，确保下次读取时重新从系统变量加载
func InvalidateCache() {
	// 由于每次 ResolveLowVersion 调用都会实时读取 SystemSetting 缓存，
	// 而 SystemSettingCache.SetVariable 会自动清除对应缓存，
	// 所以这里不需要额外操作。兼容性读取是无状态的。
}

// lookupRange 范围查询：在已排序（降序）的映射表中，找到 <= input 的最大定义版本
func lookupRange(entries []VersionEntry, input string) string {
	for _, entry := range entries {
		if compareVersions(input, entry.Version) >= 0 {
			return entry.LowVersion
		}
	}
	return ""
}

// sortEntries 将映射表按 Version 降序排列
func sortEntries(entries []VersionEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return compareVersions(entries[i].Version, entries[j].Version) > 0
	})
}

// maxVersion 取两个版本号中较大的那个
func maxVersion(v1, v2 string) string {
	if v1 == "" {
		return v2
	}
	if v2 == "" {
		return v1
	}
	if compareVersions(v1, v2) >= 0 {
		return v1
	}
	return v2
}

// parseVersion 将版本号字符串解析为数字数组
func parseVersion(version string) []int {
	if version == "" {
		return nil
	}
	parts := strings.Split(version, ".")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			num = 0
		}
		result = append(result, num)
	}
	return result
}

// compareVersions 比较两个版本号
// 返回：-1: v1 < v2, 0: v1 == v2, 1: v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		num1 := 0
		if i < len(parts1) {
			num1 = parts1[i]
		}
		num2 := 0
		if i < len(parts2) {
			num2 = parts2[i]
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
}

// IsVersionCompatible 检查当前版本是否满足最低版本要求
func IsVersionCompatible(currentVersion, minimumVersion string) bool {
	return compareVersions(currentVersion, minimumVersion) >= 0
}

// MarshalEntriesToJSON 将 VersionEntry 列表序列化为 JSON 字符串（用于初始化系统变量）
func MarshalEntriesToJSON(entries []VersionEntry) (string, error) {
	data, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalEntriesFromJSON 从 JSON 字符串反序列化为 VersionEntry 列表
func UnmarshalEntriesFromJSON(jsonStr string) ([]VersionEntry, error) {
	var entries []VersionEntry
	if err := json.Unmarshal([]byte(jsonStr), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
