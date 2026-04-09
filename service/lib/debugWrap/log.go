package debugWrap

import "fmt"

// DebugLogger 调试日志记录器
type DebugLogger struct {
	debug bool
}

// NewDebugLogger 创建调试日志记录器
func NewDebugLogger(debug bool) *DebugLogger {
	return &DebugLogger{
		debug: debug,
	}
}

// Wrap 包装调试消息和数据，返回格式化的日志数据
func (d *DebugLogger) wrapLog(message string, data ...any) []any {
	if !d.debug {
		return data
	}

	result := make([]any, 0, len(data)+1)
	result = append(result, fmt.Sprintf("[%s]:", message))

	for _, item := range data {
		if w, ok := item.(DebugWrapper); ok {
			result = append(result, w.Unwrap())
		} else {
			result = append(result, item)
		}
	}

	return result
}

// Log 包装调试消息和数据，返回格式化的日志数据
func (d *DebugLogger) Log(message string, data ...any) []any {
	return d.wrapLog(message, data...)
}

// 创建 JSON 包装器
func (d *DebugLogger) Json(key string, data any) DebugWrapper {
	return &JsonWrapper{
		Key:  key,
		Data: data,
	}
}

// 创建数据包装器
func (d *DebugLogger) Data(key string, data any) DebugWrapper {
	return &DataWrapper{
		Key:  key,
		Data: data,
	}
}
