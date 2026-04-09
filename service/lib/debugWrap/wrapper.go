package debugWrap

import (
	"encoding/json"
	"fmt"
)

// DebugWrapper 调试数据包装器接口
type DebugWrapper interface {
	Unwrap() string
}

// JsonWrapper JSON 包装器
type JsonWrapper struct {
	Key  string
	Data any
}

// Unwrap 实现 DebugWrapper 接口
func (j *JsonWrapper) Unwrap() string {
	return fmt.Sprintf("%s: %s", j.Key, anyToJsonStr(j.Data))
}

// DebugJsonWrapper debug模式不解析
type DebugNoJsonWrapper struct {
	Key  string
	Data any
}

// Unwrap 实现 DebugWrapper 接口
func (j *DebugNoJsonWrapper) Unwrap() string {
	return fmt.Sprintf("%s: %s", j.Key, anyToJsonStr(j.Data))
}

// DataWrapper 普通数据包装器
type DataWrapper struct {
	Key  string
	Data any
}

// Unwrap 实现 DebugWrapper 接口
func (d *DataWrapper) Unwrap() string {
	return fmt.Sprintf("%s: %v", d.Key, d.Data)
}

// anyToJsonStr 将任意数据转换为 JSON 字符串，失败时返回错误信息
func anyToJsonStr(data any) string {
	bjson, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("<json-error: %v>", err)
	}
	return string(bjson)
}
