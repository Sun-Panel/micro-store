package datatype

import (
	"database/sql/driver"
	"fmt"
)

// AppJson 自定义类型，用于存储原始的 app.json 内容
// 保持字符串形式，在数据库中以文本存储
type AppJson string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 AppJson
func (j *AppJson) Scan(value any) error {
	// 处理 nil 值
	if value == nil {
		*j = ""
		return nil
	}

	// 处理字节数组
	bytes, ok := value.([]byte)
	if !ok {
		// 尝试其他可能类型
		if str, ok := value.(string); ok {
			*j = AppJson(str)
			return nil
		}
		*j = ""
		return fmt.Errorf("failed to scan AppJson value: %v", value)
	}

	*j = AppJson(string(bytes))
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 AppJson 的字符串值
func (j AppJson) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}