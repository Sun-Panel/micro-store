package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SecurityAuditReport 代码安全审核结果
type SecurityAuditReport struct {
	IsPassed        bool                `json:"isPassed"`        // 是否通过审核
	ScanID          string              `json:"scanId"`          // 扫描任务ID
	Score           int                 `json:"score"`           // 安全评分（0-100）
	HighRiskCount   int                 `json:"highRiskCount"`   // 高风险漏洞数量
	MediumRiskCount int                 `json:"mediumRiskCount"` // 中风险漏洞数量
	LowRiskCount    int                 `json:"lowRiskCount"`    // 低风险漏洞数量
	ReportURL       string              `json:"reportUrl"`       // 审核报告URL
	Vulnerabilities []VulnerabilityInfo `json:"vulnerabilities"` // 漏洞详情列表
	ScanTime        time.Time           `json:"scanTime"`        // 扫描时间
	Error           string              `json:"error"`           // 错误信息
}

// VulnerabilityInfo 漏洞信息
type VulnerabilityInfo struct {
	ID          string `json:"id"`          // 漏洞ID
	Title       string `json:"title"`       // 漏洞标题
	Description string `json:"description"` // 漏洞描述
	Severity    string `json:"severity"`    // 严重程度：high, medium, low
	Location    string `json:"location"`    // 漏洞位置（文件路径）
	LineNumber  int    `json:"lineNumber"`  // 行号
	Remediation string `json:"remediation"` // 修复建议
}

// 查询的时候解析
func (j *SecurityAuditReport) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 保存时的编译
func (j SecurityAuditReport) Value() (driver.Value, error) {
	str, err := json.Marshal(j)
	if err != nil {
		return string(str), err
	}
	return string(str), nil
}
