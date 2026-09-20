package datatype

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// scanJSONArray 解析数据库中的 JSON 数组文本
// 需要兼容历史数据：列可能为 NULL、空字符串，甚至非 []byte 类型
func scanJSONArray(target interface{}, value interface{}) error {
	if value == nil {
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New(fmt.Sprint("unsupported JSON array value:", value))
	}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		// 空值/NULL 视为空数组（而非 nil -> JSON null），避免前端对 null 做 .map/.length 报错
		switch t := target.(type) {
		case *StringArray:
			*t = StringArray{}
		case *IntArray:
			*t = IntArray{}
		}
		return nil
	}

	// 解析失败时不阻断查询，按空数组处理
	_ = json.Unmarshal(data, target)
	return nil
}

// StringArray 以 JSON 数组形式存储的字符串切片（如关键词、预览图）
type StringArray []string

// 查询的时候解析
func (a *StringArray) Scan(value interface{}) error {
	return scanJSONArray(a, value)
}

// 保存时的编译
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return string(mustMarshal(a)), nil
}

// IntArray 以 JSON 数组形式存储的整型切片（如适用版本、代码类型）
type IntArray []int

// 查询的时候解析
func (a *IntArray) Scan(value interface{}) error {
	return scanJSONArray(a, value)
}

// 保存时的编译
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return string(mustMarshal(a)), nil
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("[]")
	}
	return b
}
