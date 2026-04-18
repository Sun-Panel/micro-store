package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type versionDescItem struct {
	Content string `json:"content"`
}
type VersionDesc map[string]versionDescItem

// 查询的时候解析（增加降级处理）
func (j *VersionDesc) Scan(value interface{}) error {
	// 处理 nil 值
	if value == nil {
		*j = make(VersionDesc)
		return nil
	}

	// 处理字节数组
	bytes, ok := value.([]byte)
	if !ok {
		// 尝试其他可能类型
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			*j = make(VersionDesc)
			return nil
		}
	}

	// 处理空字符串
	if len(bytes) == 0 {
		*j = make(VersionDesc)
		return nil
	}

	// 先尝试解析为当前结构
	err := json.Unmarshal(bytes, j)
	if err == nil {
		return nil
	}

	// 降级处理：尝试解析为旧格式（字符串数组或字符串）
	// 旧格式可能是 ["content1", "content2"] 或 "content"
	var legacyStr string
	if errLegacy := json.Unmarshal(bytes, &legacyStr); errLegacy == nil && legacyStr != "" {
		// 将旧格式字符串转换为新的 map 结构
		*j = VersionDesc{
			"default": {Content: legacyStr},
		}
		return nil
	}

	var legacyArr []string
	if errLegacy := json.Unmarshal(bytes, &legacyArr); errLegacy == nil {
		*j = make(VersionDesc)
		for i, content := range legacyArr {
			key := fmt.Sprintf("item_%d", i)
			(*j)[key] = versionDescItem{Content: content}
		}
		return nil
	}

	// 所有格式都无法解析，返回原始错误
	*j = make(VersionDesc)
	return fmt.Errorf("failed to unmarshal VersionDesc: %w", err)
}

// 保存时的编译
func (j VersionDesc) Value() (driver.Value, error) {
	str, err := json.Marshal(j)
	if err != nil {
		return string(str), err
	}
	return string(str), nil
}
