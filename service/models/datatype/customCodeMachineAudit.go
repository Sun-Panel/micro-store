package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CustomCodeMachineAudit 自定义代码「机器预审」结果。
// 作者提交后由系统产出：信任判断 + 静态安全扫描，供人工审核参考「为何未自动通过」。
type CustomCodeMachineAudit struct {
	Auto     bool      `json:"auto"`     // 是否自动审核通过
	Trusted  bool      `json:"trusted"`  // 作者是否受信任（近一年内审核通过过）
	ScanPass bool      `json:"scanPass"` // 静态安全扫描是否通过
	Reasons  []string  `json:"reasons"`  // 未通过/降级原因（人工审核时展示给审核员）
	Time     time.Time `json:"time"`     // 预审时间
}

func (j *CustomCodeMachineAudit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j CustomCodeMachineAudit) Value() (driver.Value, error) {
	if j.Time.IsZero() {
		return nil, nil
	}
	str, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(str), nil
}
